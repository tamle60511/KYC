package service

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/model"
	"CQS-KYC/internal/repository"
	"context"
	"fmt"
)

type (
	workflowService struct {
		repo repository.WorkflowRepo
	}
	WorkflowService interface {
		Create(ctx context.Context, req dto.WorkflowDefinitionCreate) error
		GetByID(ctx context.Context, id uint64) (*dto.WorkflowDefinitionRes, error)
		Update(ctx context.Context, id uint64, req dto.WorkflowDefinitionUpdate) error
		Delete(ctx context.Context, id uint64) error
		GetByCode(ctx context.Context, code string) (*model.WorkflowDefinition, error)
		GetAll(ctx context.Context) ([]dto.WorkflowDefinitionRes, error)
	}
)

func NewWorkflowService(repo repository.WorkflowRepo) WorkflowService {
	return &workflowService{
		repo: repo,
	}
}

func (w *workflowService) Create(ctx context.Context, req dto.WorkflowDefinitionCreate) error {
	// Check exist
	_, err := w.repo.FindByServiceCode(ctx, req.ServiceCode)
	if err == nil {
		return fmt.Errorf("service code %s already exists", req.ServiceCode)
	}

	// Map DTO -> Model
	wf := model.WorkflowDefinition{
		ServiceCode:  req.ServiceCode,
		Operation:    req.Operation,
		WorkflowName: req.WorkflowName,
		Description:  req.Description,
		IsActive:     true, // Mặc định true khi tạo mới
		Version:      1,    // Version đầu tiên
		Steps:        make([]model.WorkflowStep, 0),
	}

	for _, step := range req.Steps {
		wfStep := model.WorkflowStep{
			StepCode:       step.StepCode,
			StepName:       step.StepName,
			StepOrder:      step.StepOrder,
			RequiredRole:   step.RequiredRole,
			Canskip:        step.Canskip,
			CanDelegate:    step.CanDelegate,
			RequireComment: step.RequireComment,
			TimeHours:      step.TimeHours,
			Assignments:    make([]model.WorkflowStepAssignment, 0), // Sửa Assignments
		}

		for _, assign := range step.Assisments { // DTO của em vẫn tên là Assisments (nếu chưa sửa DTO)
			wfAssign := model.WorkflowStepAssignment{
				DepartmentIDs:    assign.DepartmentIDs,
				AssignedType:     assign.AssignedType,
				AssignedIdentity: assign.AssignedIdentity,
				Priority:         assign.Priority,
				IsActive:         true,
			}
			wfStep.Assignments = append(wfStep.Assignments, wfAssign)
		}
		wf.Steps = append(wf.Steps, wfStep)
	}

	// Gọi Repo Create đơn giản
	return w.repo.Create(ctx, &wf)
}
func (w *workflowService) GetByID(ctx context.Context, id uint64) (*dto.WorkflowDefinitionRes, error) {
	wf, err := w.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow definition by id %w", err)
	}
	res := &dto.WorkflowDefinitionRes{
		ID:           wf.ID,
		ServiceCode:  wf.ServiceCode,
		Operation:    wf.Operation,
		WorkflowName: wf.WorkflowName,
		Description:  wf.Description,
		IsActive:     wf.IsActive,
		Steps:        make([]dto.WorkflowStepRes, 0),
	}
	for _, step := range wf.Steps {
		stepRes := dto.WorkflowStepRes{
			ID:                   step.ID,
			WorkflowDefinitionID: step.WorkflowDefinitionID,
			StepCode:             step.StepCode,
			StepName:             step.StepName,
			StepOrder:            step.StepOrder,
			RequiredRole:         step.RequiredRole,
			Canskip:              step.Canskip,
			CanDelegate:          step.CanDelegate,
			RequireComment:       step.RequireComment,
			TimeHours:            step.TimeHours,
			Assisments:           make([]dto.WorkflowStepAssignmentRes, 0),
		}
		for _, assign := range step.Assignments {
			assignRes := dto.WorkflowStepAssignmentRes{
				ID:               assign.ID,
				StepID:           assign.StepID,
				DepartmentIDs:    assign.DepartmentIDs,
				AssignedType:     assign.AssignedType,
				AssignedIdentity: assign.AssignedIdentity,
				Priority:         assign.Priority,
				IsActive:         assign.IsActive,
			}
			stepRes.Assisments = append(stepRes.Assisments, assignRes)
		}
		res.Steps = append(res.Steps, stepRes)
	}
	return res, nil
}
func (w *workflowService) Update(ctx context.Context, id uint64, req dto.WorkflowDefinitionUpdate) error {
	wf := &model.WorkflowDefinition{
		WorkflowName: *req.WorkflowName,
		Operation:    *req.Operation,
		Description:  *req.Description,
		IsActive:     *req.IsActive,
		Steps:        make([]model.WorkflowStep, 0),
	}
	for _, step := range req.Steps {
		wfStep := model.WorkflowStep{
			StepCode:       step.StepCode,
			StepName:       step.StepName,
			StepOrder:      step.StepOrder,
			RequiredRole:   step.RequiredRole,
			Canskip:        step.Canskip,
			CanDelegate:    step.CanDelegate,
			RequireComment: step.RequireComment,
			TimeHours:      step.TimeHours,
			Assignments:    make([]model.WorkflowStepAssignment, 0),
		}
		for _, assign := range step.Assisments {
			wfAssign := model.WorkflowStepAssignment{
				DepartmentIDs:    assign.DepartmentIDs,
				AssignedType:     assign.AssignedType,
				AssignedIdentity: assign.AssignedIdentity,
				Priority:         assign.Priority,
				IsActive:         assign.IsActive,
			}
			wfStep.Assignments = append(wfStep.Assignments, wfAssign)
		}
		wf.Steps = append(wf.Steps, wfStep)
	}
	return w.repo.Update(ctx, id, wf)
}
func (w *workflowService) Delete(ctx context.Context, id uint64) error {
	return w.repo.Delete(ctx, id)
}
func (w *workflowService) GetAll(ctx context.Context) ([]dto.WorkflowDefinitionRes, error) {
	wf, err := w.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all workflow definitions %w", err)
	}
	var res []dto.WorkflowDefinitionRes
	for _, item := range wf {
		wfRes := dto.WorkflowDefinitionRes{
			ID:           item.ID,
			ServiceCode:  item.ServiceCode,
			Operation:    item.Operation,
			WorkflowName: item.WorkflowName,
			Description:  item.Description,
			IsActive:     item.IsActive,
			Steps:        make([]dto.WorkflowStepRes, 0),
		}
		for _, step := range item.Steps {
			stepRes := dto.WorkflowStepRes{
				ID:                   step.ID,
				WorkflowDefinitionID: step.WorkflowDefinitionID,
				StepCode:             step.StepCode,
				StepName:             step.StepName,
				StepOrder:            step.StepOrder,
				RequiredRole:         step.RequiredRole,
				Canskip:              step.Canskip,
				CanDelegate:          step.CanDelegate,
				RequireComment:       step.RequireComment,
				TimeHours:            step.TimeHours,
				Assisments:           make([]dto.WorkflowStepAssignmentRes, 0),
			}
			for _, assign := range step.Assignments {
				assignRes := dto.WorkflowStepAssignmentRes{
					ID:               assign.ID,
					StepID:           assign.StepID,
					DepartmentIDs:    assign.DepartmentIDs,
					AssignedType:     assign.AssignedType,
					AssignedIdentity: assign.AssignedIdentity,
					Priority:         assign.Priority,
					IsActive:         assign.IsActive,
				}
				stepRes.Assisments = append(stepRes.Assisments, assignRes)
			}
			wfRes.Steps = append(wfRes.Steps, stepRes)
		}
		res = append(res, wfRes)
	}
	return res, nil
}

func (w *workflowService) GetByCode(ctx context.Context, code string) (*model.WorkflowDefinition, error) {
	wf, err := w.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow definition by code %w", err)
	}
	return wf, nil
}
