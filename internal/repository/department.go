package repository

import (
	"CQS-KYC/internal/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type (
	departmentRepo struct {
		db *gorm.DB
	}
	DepartmentRepo interface {
		Create(ctx context.Context, req *model.Department) error
		GetByID(ctx context.Context, id uint64) (*model.Department, error)
		GetByCode(ctx context.Context, code string) (*model.Department, error)
		Update(ctx context.Context, id uint64, req map[string]interface{}) error
		Delete(ctx context.Context, id uint64) error
		GetAll(ctx context.Context) ([]model.Department, error)
	}
)

func NewDepartmentRepo(db *gorm.DB) DepartmentRepo {
	return &departmentRepo{
		db: db,
	}
}

func (d *departmentRepo) Create(ctx context.Context, req *model.Department) error {
	return d.db.WithContext(ctx).Model(&model.Department{}).Create(req).Error
}
func (d *departmentRepo) GetByID(ctx context.Context, id uint64) (*model.Department, error) {
	var dept model.Department
	if err := d.db.WithContext(ctx).Model(&model.Department{}).Where("id = ?", id).First(&dept).Error; err != nil {
		return nil, fmt.Errorf("failed to get by department id %w", err)
	}
	return &dept, nil
}
func (d *departmentRepo) GetByCode(ctx context.Context, code string) (*model.Department, error) {
	var dept model.Department
	if err := d.db.WithContext(ctx).Model(&model.Department{}).Where("code = ?", code).First(&dept).Error; err != nil {
		return nil, fmt.Errorf("failed to get by department id %w", err)
	}
	return &dept, nil
}
func (d *departmentRepo) Update(ctx context.Context, id uint64, req map[string]interface{}) error {
	return d.db.WithContext(ctx).Model(&model.Department{}).Where("id = ?", id).Updates(req).Error
}
func (d *departmentRepo) Delete(ctx context.Context, id uint64) error {
	return d.db.WithContext(ctx).Delete(&model.Department{}, id).Error
}
func (d *departmentRepo) GetAll(ctx context.Context) ([]model.Department, error) {
	var depts []model.Department
	if err := d.db.WithContext(ctx).
		Preload("Parent").
		Find(&depts).Error; err != nil {
		return nil, fmt.Errorf("failed get all depts %w", err)
	}
	return depts, nil
}
