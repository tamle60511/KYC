package service

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/model"
	"CQS-KYC/internal/repository"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type (
	managerService struct {
		repo repository.ManagerRepo
	}
	ManagerService interface {
		Create(ctx context.Context, req dto.ManagerCreate) error
		GetByID(ctx context.Context, id uint64) (*dto.ManagerRes, error)
		Update(ctx context.Context, id uint64, req dto.ManagerUpdate) error
		Delete(ctx context.Context, id uint64) error
		GetAll(ctx context.Context) ([]dto.ManagerRes, error)
	}
)

func NewManagerService(repo repository.ManagerRepo) ManagerService {
	return &managerService{
		repo: repo,
	}
}

func (m *managerService) Create(ctx context.Context, req dto.ManagerCreate) error {
	_, err := m.repo.GetByUserID(ctx, req.UserID)
	if err == nil {
		return fmt.Errorf("manager user already exits %w", err)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed get by manager user id %w", err)
	}
	_, err = m.repo.GetByDepartmentID(ctx, req.DepartmentID)
	if err == nil {
		return fmt.Errorf("department %d already has a manager", req.DepartmentID)
	}
	manager := model.Manager{
		UserID:       req.UserID,
		DepartmentID: req.DepartmentID,
	}
	if err := m.repo.Create(ctx, &manager); err != nil {
		return fmt.Errorf("failed to create manager %w", err)
	}
	return nil
}
func (m *managerService) GetByID(ctx context.Context, id uint64) (*dto.ManagerRes, error) {
	manager, err := m.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get by manager id %w", err)
	}
	return &dto.ManagerRes{
		ID:           manager.ID,
		UserID:       manager.UserID,
		DepartmentID: manager.DepartmentID,
		CreatedAt:    manager.CreatedAt,
		UpdatedAt:    manager.UpdatedAt,
	}, nil
}
func (m *managerService) Update(ctx context.Context, id uint64, req dto.ManagerUpdate) error {
	updates := make(map[string]interface{})
	if req.DepartmentID != nil {
		updates["department_id"] = *req.DepartmentID
	}
	if len(updates) == 0 {
		return fmt.Errorf("no failed record ")
	}
	if err := m.repo.Update(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update manager %w", err)
	}
	return nil
}
func (m *managerService) Delete(ctx context.Context, id uint64) error {
	return m.repo.Delete(ctx, id)
}
func (m *managerService) GetAll(ctx context.Context) ([]dto.ManagerRes, error) {
	managers, err := m.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]dto.ManagerRes, 0, len(managers))
	for _, v := range managers {
		// Map data sang DTO đầy đủ thông tin
		item := dto.ManagerRes{
			ID:           v.ID,
			UserID:       v.UserID,
			DepartmentID: v.DepartmentID,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		}

		// Map User Info nếu có Preload
		// if v.User.ID != 0 {
		// 	item.UserName = v.User.FullName
		// 	item.UserEmail = v.User.Email
		// }

		res = append(res, item)
	}
	return res, nil
}
