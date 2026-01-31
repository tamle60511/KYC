package service

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/model"
	"CQS-KYC/internal/repository"
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type (
	instanceService struct {
		repo repository.InstanceRepo
		db   *gorm.DB // Cần DB để mở Transaction
	}
	InstanceService interface {
		InitiateWorkflow(
			tx *gorm.DB, workflowID uint64,
			serviceCode string,
			docNum string,
			docType string,
			creatorID string,
			factoryID uint64,
			deptID uint64,
			requestData []byte,
			ip string,
			device string,
		) (*model.WorkflowInstance, error)
		Initiate(ctx context.Context, userID string, req dto.WorkflowInitiateReq) (*model.WorkflowInstance, error)
		ProcessAction(ctx context.Context, instanceID uint64, userID, userName string, req dto.WorkflowActionReq) error
		GetPendingTasks(ctx context.Context, userID string) ([]dto.PendingTaskRes, error)
		GetHistory(ctx context.Context, instanceID uint64) ([]dto.WorkflowLogRes, error)
	}
)

func NewInstanceService(repo repository.InstanceRepo, db *gorm.DB) InstanceService {
	return &instanceService{
		repo: repo,
		db:   db,
	}
}

func (s *instanceService) InitiateWorkflow(
	tx *gorm.DB, workflowID uint64,
	serviceCode string,
	docNum string,
	docType string,
	creatorID string,
	factoryID uint64,
	deptID uint64,
	requestData []byte,
	ip string,
	device string,
) (*model.WorkflowInstance, error) {
	return s.repo.InitiateWorkflow(tx, workflowID, serviceCode, docNum, docType, creatorID, factoryID, deptID, requestData, ip, device)
}

// 1. Tạo đơn mới
func (s *instanceService) Initiate(ctx context.Context, userID string, req dto.WorkflowInitiateReq) (*model.WorkflowInstance, error) {
	// Ép kiểu request_data sang JSON bytes
	reqDataBytes, err := json.Marshal(req.RequestData)
	if err != nil {
		return nil, fmt.Errorf("invalid request data json: %w", err)
	}

	// Mở Transaction (Vì hàm Repo yêu cầu tx)
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Gọi Repo
	// Lưu ý: IP và Device tạm thời để trống hoặc lấy từ Context nếu Middleware có set
	instance, err := s.repo.InitiateWorkflow(
		tx,
		req.WorkflowID,
		req.ServiceCode,
		req.DocNum,
		req.DocType,
		userID,
		req.FactoryID,
		req.DeptID,
		reqDataBytes,
		"127.0.0.1", // IP Mock
		"WebClient", // Device Mock
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return instance, nil
}

// 2. Xử lý Duyệt/Từ chối
func (s *instanceService) ProcessAction(ctx context.Context, instanceID uint64, userID, userName string, req dto.WorkflowActionReq) error {
	return s.repo.ProcessAction(ctx, instanceID, userID, userName, req.Action, req.Comment)
}

// 3. Lấy danh sách việc cần làm (Mapping Model -> DTO)
func (s *instanceService) GetPendingTasks(ctx context.Context, userID string) ([]dto.PendingTaskRes, error) {
	tasks, err := s.repo.GetPendingTasks(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Map data cho đẹp
	var res []dto.PendingTaskRes
	for _, t := range tasks {
		item := dto.PendingTaskRes{
			TaskID:     t.ID,
			InstanceID: t.InstanceID,
			StepName:   t.StepName,
			Status:     t.Status,
			ReceivedAt: t.CreatedAt,
		}
		// Lấy thông tin từ bảng cha (Instance) nhờ Preload
		if t.Instance != nil {
			item.DocNum = t.Instance.DocNum
			item.DocType = t.Instance.DocType
			item.ServiceCode = t.Instance.ServiceCode
			item.CreatorID = t.Instance.CreatorID
		}
		res = append(res, item)
	}
	return res, nil
}

// 4. Lấy lịch sử
func (s *instanceService) GetHistory(ctx context.Context, instanceID uint64) ([]dto.WorkflowLogRes, error) {
	logs, err := s.repo.GetHistory(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	var res []dto.WorkflowLogRes
	for _, l := range logs {
		res = append(res, dto.WorkflowLogRes{
			StepName:  l.StepName,
			Action:    l.Action,
			ActorName: l.ActorName,
			Comment:   l.Comment,
			Time:      l.CreatedAt,
		})
	}
	return res, nil
}
