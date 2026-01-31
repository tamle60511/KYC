package service

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/model"
	"CQS-KYC/internal/repository"
	"CQS-KYC/utils"
	"context"
	"fmt"
)

type (
	userService struct {
		repo repository.UserRepo
	}
	UserService interface {
		Create(ctx context.Context, req dto.UserCreate) error
		GetByID(ctx context.Context, id uint64) (*dto.UserRes, error)
		Update(ctx context.Context, id uint64, req dto.UserUpdate) error
		Delete(ctx context.Context, id uint64) error
		GetAll(ctx context.Context) ([]dto.UserRes, error)
	}
)

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{
		repo: repo,
	}
}

func (u *userService) Create(ctx context.Context, req dto.UserCreate) error {

	hashpass, err := utils.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password %w", err)
	}
	user := model.User{
		UserCode:       req.UserCode,
		FullName:       req.FullName,
		Email:          req.Email,
		Phone:          req.Phone,
		Password:       hashpass,
		SignatureImage: req.SignatureImage,
		DepartmentID:   req.DepartmentID,
		ManagerID:      &req.ManagerID,
		PositionID:     req.PositionID,
		FactoryID:      req.FactoryID,
	}
	if req.ManagerID != 0 {
		user.ManagerID = &req.ManagerID
	} else {
		user.ManagerID = nil
	}

	if err := u.repo.Create(ctx, &user); err != nil {
		return fmt.Errorf("failed to create user %w", err)
	}
	return nil
}
func (u *userService) GetByID(ctx context.Context, id uint64) (*dto.UserRes, error) {
	user, err := u.repo.GetByID(ctx, uint64(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get by id service %w", err)
	}
	return &dto.UserRes{
		ID:             user.ID,
		UserCode:       user.UserCode,
		FullName:       user.FullName,
		Email:          user.Email,
		Phone:          user.Phone,
		DepartmentID:   user.DepartmentID,
		SignatureImage: user.SignatureImage,
		ManagerID:      safeUint64(user.ManagerID),
		PositionID:     user.PositionID,
		FactoryID:      user.FactoryID,
		Role:           user.Role,
		IsActive:       user.IsActive,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}
func (u *userService) Update(ctx context.Context, id uint64, req dto.UserUpdate) error {
	updates := make(map[string]interface{})

	if req.FullName != nil {
		updates["full_name"] = *req.FullName
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Password != nil && *req.Password != "" {
		hashPass, err := utils.HashPassword(*req.Password)
		if err != nil {
			return fmt.Errorf("failed to hash new password: %w", err)
		}
		updates["password"] = hashPass
	}
	if req.DepartmentID != nil {
		updates["department_id"] = *req.DepartmentID
	}
	if req.PositionID != nil {
		updates["position_id"] = *req.PositionID
	}
	if req.ManagerID != nil {
		updates["manager_id"] = *req.ManagerID
	}
	if req.FactoryID != nil {
		updates["factory_id"] = *req.FactoryID
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}
	if err := u.repo.Update(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update user service %w", err)
	}
	return nil
}
func (u *userService) Delete(ctx context.Context, id uint64) error {
	return u.repo.Delete(ctx, uint64(id))
}
func (u *userService) GetAll(ctx context.Context) ([]dto.UserRes, error) {
	users, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users %w", err)
	}

	res := make([]dto.UserRes, 0, len(users))
	for _, v := range users {
		res = append(res, dto.UserRes{
			ID:             v.ID,
			UserCode:       v.UserCode,
			FullName:       v.FullName,
			Email:          v.Email,
			Phone:          v.Phone,
			DepartmentID:   v.DepartmentID,
			SignatureImage: v.SignatureImage,
			ManagerID:      safeUint64(v.ManagerID),
			PositionID:     v.PositionID,
			FactoryID:      v.FactoryID,
			Role:           v.Role,
			IsActive:       v.IsActive,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
		})
	}
	return res, nil
}

func safeUint64(ptr *uint64) uint64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}
