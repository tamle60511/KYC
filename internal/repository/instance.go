package repository

import (
	"CQS-KYC/internal/model" // Import package utils chứa SignatureHelper
	"context"                // Cần để parse JSON DepartmentIDs
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type (
	instanceRepo struct {
		db              *gorm.DB
		groupRepo       GroupRepo
		signatureHelper SignatureHelper // Sửa lại đường dẫn import
	}
	InstanceRepo interface {
		// Core Flow
		InitiateWorkflow(tx *gorm.DB, workflowID uint64, serviceCode, docNum, docType, creatorID string, factoryID, deptID uint64, requestData []byte, ip, device string) (*model.WorkflowInstance, error)
		ProcessAction(ctx context.Context, instanceID uint64, actorID, actorName, action, comment string) error

		// View Data (CÁI EM ĐANG THIẾU)
		GetPendingTasks(ctx context.Context, userID string) ([]model.WorkflowTask, error)
		GetHistory(ctx context.Context, instanceID uint64) ([]model.WorkflowLog, error)
	}
)

func NewWorkflowEngine(db *gorm.DB, groupRepo GroupRepo, signatureHelper SignatureHelper) InstanceRepo {
	return &instanceRepo{db: db, groupRepo: groupRepo, signatureHelper: signatureHelper}
}

// =============================================================================
// 1. KHỞI TẠO (INITIATE)
// =============================================================================
func (e *instanceRepo) InitiateWorkflow(
	tx *gorm.DB,
	workflowID uint64,
	serviceCode, docNum, docType, creatorID string,
	factoryID, deptID uint64,
	requestData []byte, ip, device string,
) (*model.WorkflowInstance, error) {

	if tx == nil {
		return nil, errors.New("transaction is required")
	}

	// 1. Load Workflow & Đếm tổng số bước
	var workflow model.WorkflowDefinition
	if err := tx.Preload("Steps").First(&workflow, workflowID).Error; err != nil {
		return nil, fmt.Errorf("workflow definition not found: %d", workflowID)
	}
	if len(workflow.Steps) == 0 {
		return nil, errors.New("workflow definition has no steps")
	}

	// 2. Tạo Instance
	instance := model.WorkflowInstance{
		WorkflowID:   workflow.ID,
		ServiceCode:  serviceCode,
		DocNum:       docNum,
		DocType:      docType,
		FactoryID:    factoryID,
		DepartmentID: deptID,
		RequestData:  requestData,
		CurrentStep:  workflow.Steps[0].StepOrder, // Bước đầu tiên
		TotalSteps:   len(workflow.Steps),         // Tính tổng số bước
		Status:       model.STATUS_IN_PROGRESS,
		CreatorID:    creatorID,
		StartedAt:    time.Now(),
	}

	if err := tx.Create(&instance).Error; err != nil {
		return nil, err
	}

	// 3. Tạo Log Submit (Chữ ký người tạo)
	sigHash, dataHash, signedTime := e.signatureHelper.GenerateSignature(creatorID, docNum, model.ACTION_SUBMIT, 0, requestData)

	log := model.WorkflowLog{
		InstanceID:       instance.ID,
		StepOrder:        0,
		StepName:         "Submit",
		Action:           model.ACTION_SUBMIT,
		ActorID:          creatorID,
		ActorName:        "System Creator",
		Comment:          "Created via ERP System",
		SignatureHash:    sigHash,
		DataSnapshotHash: dataHash,
		SignedTimestamp:  signedTime,
		IPAddress:        ip,
		DeviceInfo:       device,
	}
	if err := tx.Create(&log).Error; err != nil {
		return nil, err
	}

	// 4. Phân bổ Task cho bước đầu tiên
	if err := e.distributeTasks(tx, &instance, &workflow.Steps[0]); err != nil {
		return nil, fmt.Errorf("failed to assign first task: %v", err)
	}

	return &instance, nil
}

// =============================================================================
// 2. XỬ LÝ DUYỆT (PROCESS ACTION)
// =============================================================================
func (e *instanceRepo) ProcessAction(
	ctx context.Context,
	instanceID uint64,
	actorID, actorName, action, comment string,
) error {
	return e.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Load Instance
		var instance model.WorkflowInstance
		if err := tx.First(&instance, instanceID).Error; err != nil {
			return err
		}

		if instance.Status != model.STATUS_IN_PROGRESS {
			return errors.New("request is not in progress")
		}

		// 2. Check quyền: Tìm task của User hoặc Group của User
		userGroups := e.getUserGroups(ctx, actorID) // Sửa: thêm ctx

		var myTask model.WorkflowTask
		// Logic: Task assigned to ME (is_group=false) OR assigned to MY GROUP (is_group=true)
		query := tx.Where("instance_id = ? AND status = ?", instanceID, "PENDING").
			Where(
				tx.Where("assigned_to = ? AND is_group = ?", actorID, false).
					Or("assigned_to IN ? AND is_group = ?", userGroups, true),
			)

		if err := query.First(&myTask).Error; err != nil {
			return errors.New("you do not have permission to approve this request")
		}

		// 3. Xử lý Action
		// 3.1 Xóa Task (Done task)
		if err := tx.Delete(&myTask).Error; err != nil {
			return err
		}

		// 3.2 Ghi Log
		sigHash, dataHash, signedTime := e.signatureHelper.GenerateSignature(actorID, instance.DocNum, action, instance.CurrentStep, instance.RequestData)

		log := model.WorkflowLog{
			InstanceID:       instance.ID,
			StepOrder:        instance.CurrentStep,
			StepName:         myTask.StepName, // Lấy tên từ Task, ko cần query lại Step
			Action:           action,
			ActorID:          actorID,
			ActorName:        actorName,
			Comment:          comment,
			SignatureHash:    sigHash,
			DataSnapshotHash: dataHash,
			SignedTimestamp:  signedTime,
		}
		if err := tx.Create(&log).Error; err != nil {
			return err
		}

		// 3.3 Điều hướng (Routing)
		if action == model.ACTION_REJECT {
			// REJECT: Hủy toàn bộ
			if err := tx.Where("instance_id = ?", instanceID).Delete(&model.WorkflowTask{}).Error; err != nil {
				return err
			}

			now := time.Now()
			instance.Status = model.STATUS_REJECTED
			instance.CompletedAt = &now
			return tx.Save(&instance).Error
		}

		if action == model.ACTION_APPROVE {
			// APPROVE: Tìm bước tiếp theo
			var nextStep model.WorkflowStep
			err := tx.Where("workflow_definition_id = ? AND step_order > ?", instance.WorkflowID, instance.CurrentStep).
				Order("step_order ASC").First(&nextStep).Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Hết bước -> SUCCESS
				now := time.Now()
				instance.Status = model.STATUS_APPROVED
				instance.CompletedAt = &now
				// Ở đây có thể bắn Webhook thông báo về ERP
				return tx.Save(&instance).Error
			}

			// Còn bước -> Update Instance & Tạo Task mới
			instance.CurrentStep = nextStep.StepOrder
			if err := tx.Save(&instance).Error; err != nil {
				return err
			}

			return e.distributeTasks(tx, &instance, &nextStep)
		}

		return nil
	})
}

// =============================================================================
// 3. HELPER LOGIC (QUAN TRỌNG)
// =============================================================================
func (e *instanceRepo) distributeTasks(tx *gorm.DB, instance *model.WorkflowInstance, step *model.WorkflowStep) error {
	// Lấy tất cả rule gán của bước này
	var assignments []model.WorkflowStepAssignment
	// Lưu ý: step.ID phải đúng là ID của bảng steps
	if err := tx.Where("step_id = ? AND is_active = ?", step.ID, true).Find(&assignments).Error; err != nil {
		return err
	}

	tasksCreated := 0
	for _, assign := range assignments {
		// 1. Lọc theo Factory (Nếu rule có set Factory)
		if assign.FactoryID != nil && *assign.FactoryID != instance.FactoryID {
			continue // Khác nhà máy -> Bỏ qua
		}

		// 2. Lọc theo Department (JSON parsing)
		if len(assign.DepartmentIDs) > 0 { // assign.DepartmentIDs giờ là []uint64 nhờ GORM serializer
			isMatch := false
			for _, dID := range assign.DepartmentIDs {
				if dID == instance.DepartmentID {
					isMatch = true
					break
				}
			}
			if !isMatch {
				continue // Khác phòng ban -> Bỏ qua
			}
		}

		// 3. Tạo Task
		task := model.WorkflowTask{
			InstanceID: instance.ID,
			StepID:     step.ID,
			StepOrder:  step.StepOrder,
			StepName:   step.StepName,
			Status:     "PENDING",
			AssignedTo: assign.AssignedIdentity,
			IsGroup:    assign.AssignedType == "GROUP", // Cần đảm bảo enum đúng
		}
		if err := tx.Create(&task).Error; err != nil {
			return err
		}
		tasksCreated++
	}

	if tasksCreated == 0 {
		return fmt.Errorf("configuration error: step %d has no valid assignment for factory %d dept %d", step.StepOrder, instance.FactoryID, instance.DepartmentID)
	}

	return nil
}

// =============================================================================
// 4. VIEW DATA (CÁI EM THIẾU)
// =============================================================================

// Lấy danh sách việc cần làm của User
func (e *instanceRepo) GetPendingTasks(ctx context.Context, userID string) ([]model.WorkflowTask, error) {
	userGroups := e.getUserGroups(ctx, userID)

	var tasks []model.WorkflowTask
	err := e.db.WithContext(ctx).
		Preload("Instance"). // Join để lấy thông tin đơn hàng (DocNum, ServiceCode)
		Where("status = ?", "PENDING").
		Where(
			e.db.Where("assigned_to = ? AND is_group = ?", userID, false).
				Or("assigned_to IN ? AND is_group = ?", userGroups, true),
		).
		Order("created_at DESC").
		Find(&tasks).Error

	return tasks, err
}

// Lấy lịch sử duyệt của 1 đơn
func (e *instanceRepo) GetHistory(ctx context.Context, instanceID uint64) ([]model.WorkflowLog, error) {
	var logs []model.WorkflowLog
	err := e.db.WithContext(ctx).
		Where("instance_id = ?", instanceID).
		Order("step_order ASC").
		Find(&logs).Error
	return logs, err
}

// Helper lấy group (Wrapper lại repo cũ)
func (e *instanceRepo) getUserGroups(ctx context.Context, userID string) []string {
	groups, err := e.groupRepo.GetGroupsByUserID(ctx, userID)
	if err != nil {
		return []string{}
	}
	return groups
}
