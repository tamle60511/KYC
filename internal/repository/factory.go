package repository

import (
	"CQS-KYC/internal/model"
	"context"

	"gorm.io/gorm"
)

type (
	factoryRepo struct {
		db *gorm.DB
	}
	FactoryRepo interface {
		Create(ctx context.Context, req *model.Factory) error
		GetByID(ctx context.Context, id uint64) (*model.Factory, error)
		GetByCode(ctx context.Context, code string) (*model.Factory, error)
		Update(ctx context.Context, id uint64, req map[string]interface{}) error
		Delete(ctx context.Context, id uint64) error
		GetList(ctx context.Context) ([]*model.Factory, error)
	}
)

func NewFactoryRepo(db *gorm.DB) FactoryRepo {
	return &factoryRepo{db: db}
}

func (r *factoryRepo) Create(ctx context.Context, req *model.Factory) error {
	return r.db.WithContext(ctx).Model(&model.Factory{}).Create(req).Error
}
func (r *factoryRepo) GetByID(ctx context.Context, id uint64) (*model.Factory, error) {
	var factory model.Factory
	if err := r.db.WithContext(ctx).First(&factory, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &factory, nil
}
func (r *factoryRepo) GetByCode(ctx context.Context, code string) (*model.Factory, error) {
	var factory model.Factory
	if err := r.db.WithContext(ctx).First(&factory, "code = ?", code).Error; err != nil {
		return nil, err
	}
	return &factory, nil
}
func (r *factoryRepo) Update(ctx context.Context, id uint64, req map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&model.Factory{}).Where("id = ?", id).Updates(req).Error
}
func (r *factoryRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.Factory{}, "id = ?", id).Error
}
func (r *factoryRepo) GetList(ctx context.Context) ([]*model.Factory, error) {
	var factories []*model.Factory
	if err := r.db.WithContext(ctx).Find(&factories).Error; err != nil {
		return nil, err
	}
	return factories, nil
}
