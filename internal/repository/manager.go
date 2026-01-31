package repository

import (
	"CQS-KYC/internal/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type (
	managerRepo struct {
		db *gorm.DB
	}
	ManagerRepo interface {
		Create(ctx context.Context, req *model.Manager) error
		GetByID(ctx context.Context, id uint64) (*model.Manager, error)
		GetByUserID(ctx context.Context, userid uint64) (*model.Manager, error)
		Update(ctx context.Context, id uint64, req map[string]interface{}) error
		Delete(ctx context.Context, id uint64) error
		GetAll(ctx context.Context) ([]model.Manager, error)
		GetByDepartmentID(ctx context.Context, deptID uint64) (*model.Manager, error)
	}
)

func NewManagerRepo(db *gorm.DB) ManagerRepo {
	return &managerRepo{
		db: db,
	}
}

func (m *managerRepo) Create(ctx context.Context, req *model.Manager) error {
	return m.db.WithContext(ctx).Model(&model.Manager{}).Create(req).Error
}
func (m *managerRepo) GetByID(ctx context.Context, id uint64) (*model.Manager, error) {
	var manager model.Manager
	if err := m.db.WithContext(ctx).Where("id = ?", id).First(&manager).Error; err != nil {
		return nil, fmt.Errorf("failed to get by manager id %w", err)
	}
	return &manager, nil
}
func (m *managerRepo) GetByUserID(ctx context.Context, userid uint64) (*model.Manager, error) {
	var manager model.Manager
	if err := m.db.WithContext(ctx).Where("user_id = ?", userid).First(&manager).Error; err != nil {
		return nil, fmt.Errorf("failed to get by manager id %w", err)
	}
	return &manager, nil
}
func (m *managerRepo) Update(ctx context.Context, id uint64, req map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&model.Manager{}).Where("id  = ?", id).Updates(req).Error
}
func (m *managerRepo) Delete(ctx context.Context, id uint64) error {
	return m.db.WithContext(ctx).Delete(&model.Manager{}, id).Error
}
func (m *managerRepo) GetAll(ctx context.Context) ([]model.Manager, error) {
	var managers []model.Manager
	if err := m.db.WithContext(ctx).Find(&managers).Error; err != nil {
		return nil, fmt.Errorf("failed to get all manager repo %w", err)
	}
	return managers, nil
}

func (m *managerRepo) GetByDepartmentID(ctx context.Context, deptID uint64) (*model.Manager, error) {
	var manager model.Manager
	if err := m.db.WithContext(ctx).Where("department_id = ?", deptID).First(&manager).Error; err != nil {
		return nil, fmt.Errorf("failed to get by department id %w", err)
	}
	return &manager, nil
}
