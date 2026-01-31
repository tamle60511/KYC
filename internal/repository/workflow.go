package repository

import (
	"CQS-KYC/internal/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type (
	workflowRepo struct {
		db *gorm.DB
	}
	WorkflowRepo interface {
		Create(ctx context.Context, wf *model.WorkflowDefinition) error
		FindByID(ctx context.Context, id uint64) (*model.WorkflowDefinition, error)
		FindByServiceCode(ctx context.Context, serviceCode string) (*model.WorkflowDefinition, error)
		Update(ctx context.Context, id uint64, wf *model.WorkflowDefinition) error
		Delete(ctx context.Context, id uint64) error
		GetAll(ctx context.Context) ([]model.WorkflowDefinition, error)
		GetByCode(ctx context.Context, code string) (*model.WorkflowDefinition, error)
	}
)

func NewWorkflowRepo(db *gorm.DB) WorkflowRepo {
	return &workflowRepo{
		db: db,
	}
}

func (w *workflowRepo) Create(ctx context.Context, wf *model.WorkflowDefinition) error {
	return w.db.WithContext(ctx).Create(wf).Error
}
func (w *workflowRepo) FindByID(ctx context.Context, id uint64) (*model.WorkflowDefinition, error) {
	var wf model.WorkflowDefinition
	if err := w.db.WithContext(ctx).Preload("Steps.Assisments").First(&wf, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &wf, nil
}
func (w *workflowRepo) FindByServiceCode(ctx context.Context, serviceCode string) (*model.WorkflowDefinition, error) {
	var wf model.WorkflowDefinition
	if err := w.db.WithContext(ctx).
		Preload("Steps.Assignments").
		Where("service_code = ? AND is_active = ?", serviceCode, true).
		First(&wf).Error; err != nil {
		return nil, err
	}
	return &wf, nil
}
func (w *workflowRepo) Update(ctx context.Context, id uint64, req *model.WorkflowDefinition) error {
	return w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Tìm quy trình cũ
		var oldWf model.WorkflowDefinition
		if err := tx.Where("id = ?", id).First(&oldWf).Error; err != nil {
			return fmt.Errorf("workflow not found: %w", err)
		}

		// 2. Archive quy trình cũ (Rename ServiceCode để nhả Unique Index)
		// Ví dụ: PURI05 -> PURI05_v1_ARCHIVED_17000000
		archivedCode := fmt.Sprintf("%s_v%d_ARCHIVED_%d", oldWf.ServiceCode, oldWf.Version, oldWf.ID)
		if err := tx.Model(&oldWf).Updates(map[string]interface{}{
			"is_active":    false,
			"service_code": archivedCode,
		}).Error; err != nil {
			return fmt.Errorf("failed to archive old workflow: %w", err)
		}

		// 3. Chuẩn bị quy trình mới (Clone từ Request nhưng giữ lại ServiceCode gốc)
		// Lưu ý: req ở đây đã được Service map đầy đủ các Step và Assignment
		newWf := *req
		newWf.ID = 0 // Đảm bảo tạo mới
		newWf.Version = oldWf.Version + 1
		newWf.ServiceCode = oldWf.ServiceCode // Dùng lại code gốc PURI05
		newWf.IsActive = true

		// 4. Create Deep Insert
		if err := tx.Create(&newWf).Error; err != nil {
			return fmt.Errorf("failed to create new version: %w", err)
		}

		return nil
	})
}
func (w *workflowRepo) Delete(ctx context.Context, id uint64) error {
	if err := w.db.WithContext(ctx).Delete(&model.WorkflowDefinition{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
func (w *workflowRepo) GetAll(ctx context.Context) ([]model.WorkflowDefinition, error) {
	var wfs []model.WorkflowDefinition
	if err := w.db.WithContext(ctx).Preload("Steps.Assisments").Find(&wfs).Error; err != nil {
		return nil, err
	}
	return wfs, nil
}

func (w *workflowRepo) GetByCode(ctx context.Context, code string) (*model.WorkflowDefinition, error) {
	var wf model.WorkflowDefinition
	err := w.db.WithContext(ctx).
		Preload("Steps.Assignments"). // Sửa lại tên field đúng chính tả
		Where("operation = ? AND is_active = ?", code, true).
		Order("version desc"). // Lấy version cao nhất
		First(&wf).Error

	if err != nil {
		return nil, err
	}
	return &wf, nil
}
