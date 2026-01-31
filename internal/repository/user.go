package repository

import (
	"CQS-KYC/internal/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type (
	userRepo struct {
		db *gorm.DB
	}
	UserRepo interface {
		Create(ctx context.Context, user *model.User) error
		GetByID(ctx context.Context, id uint64) (*model.User, error)
		GetByCode(ctx context.Context, code string) (*model.User, error)
		Update(ctx context.Context, id uint64, req map[string]interface{}) error
		Delete(ctx context.Context, id uint64) error
		GetAll(ctx context.Context) ([]model.User, error)
	}
)

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Create(ctx context.Context, user *model.User) error {
	return u.db.WithContext(ctx).Model(&model.User{}).Create(user).Error
}
func (u *userRepo) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	if err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get by user id repo %w", err)
	}
	return &user, nil
}
func (u *userRepo) GetByCode(ctx context.Context, code string) (*model.User, error) {
	var user model.User
	if err := u.db.WithContext(ctx).Where("user_code = ?", code).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get by user code repo %w", err)
	}
	return &user, nil
}
func (u *userRepo) Update(ctx context.Context, id uint64, req map[string]interface{}) error {
	return u.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(req).Error
}
func (u *userRepo) Delete(ctx context.Context, id uint64) error {
	return u.db.WithContext(ctx).Delete(&model.User{}, id).Error
}
func (u *userRepo) GetAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := u.db.WithContext(ctx).
		Preload("Department").
		Preload("Position").
		Preload("Factory").
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get all user repo %w", err)
	}
	return users, nil
}
