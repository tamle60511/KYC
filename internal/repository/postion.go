package repository

import (
	"CQS-KYC/internal/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type (
	positionRepo struct {
		db *gorm.DB
	}
	PositionRepo interface {
		Create(ctx context.Context, req *model.Position) error
		GetByID(ctx context.Context, id uint64) (*model.Position, error)
		GetByCode(ctx context.Context, code string) (*model.Position, error)
		Update(ctx context.Context, id uint64, req map[string]interface{}) error
		Delete(ctx context.Context, id uint64) error
		GetAll(ctx context.Context) ([]model.Position, error)
	}
)

func NewPositionRepo(db *gorm.DB) PositionRepo {
	return &positionRepo{
		db: db,
	}
}

func (p *positionRepo) Create(ctx context.Context, req *model.Position) error {
	return p.db.WithContext(ctx).Model(&model.Position{}).Create(req).Error
}
func (p *positionRepo) GetByID(ctx context.Context, id uint64) (*model.Position, error) {
	var pos model.Position
	if err := p.db.WithContext(ctx).Where("id = ?", id).First(&pos).Error; err != nil {
		return nil, fmt.Errorf("failed to gey by position id %w", err)
	}
	return &pos, nil
}
func (p *positionRepo) GetByCode(ctx context.Context, code string) (*model.Position, error) {
	var pos model.Position
	if err := p.db.WithContext(ctx).Where("code = ?", code).First(&pos).Error; err != nil {
		return nil, fmt.Errorf("failed to gey by position code %w", err)
	}
	return &pos, nil
}
func (p *positionRepo) Update(ctx context.Context, id uint64, req map[string]interface{}) error {
	return p.db.WithContext(ctx).Model(&model.Position{}).Where("id = ?", id).Updates(req).Error
}
func (p *positionRepo) Delete(ctx context.Context, id uint64) error {
	return p.db.WithContext(ctx).Delete(&model.Position{}, id).Error
}
func (p *positionRepo) GetAll(ctx context.Context) ([]model.Position, error) {
	var pos []model.Position
	if err := p.db.WithContext(ctx).Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("failed to get all position %w", err)
	}
	return pos, nil
}
