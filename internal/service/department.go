package service

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/model"
	"CQS-KYC/internal/repository"
	"context"
	"fmt"
)

type (
	departmentService struct {
		repo repository.DepartmentRepo
	}
	DepartmentService interface {
		Create(ctx context.Context, req dto.DepartmentCreate) error
		GetByID(ctx context.Context, id uint64) (*dto.DepartmentResponse, error)
		Update(ctx context.Context, id uint64, req dto.DepartmentUpdate) error
		Delete(ctx context.Context, id uint64) error
		GetAll(ctx context.Context) ([]dto.DepartmentResponse, error)
	}
)

func NewDepartmentSerivce(repo repository.DepartmentRepo) DepartmentService {
	return &departmentService{
		repo: repo,
	}
}

func (d *departmentService) Create(ctx context.Context, req dto.DepartmentCreate) error {
	// 1. Check Code trùng
	_, err := d.repo.GetByCode(ctx, req.Code)
	if err == nil {
		return fmt.Errorf("department code already exists")
	}

	// 2. Logic tính Level (Cấp bậc phòng ban)
	var level uint8 = 1 // Mặc định là cấp 1 (cao nhất)

	if req.ParentID != nil && *req.ParentID != 0 {
		// Tìm phòng cha để lấy level của cha
		parentDept, err := d.repo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return fmt.Errorf("invalid parent_id: %w", err)
		}
		if parentDept.Level != nil {
			level = *parentDept.Level + 1
		}
	}

	// 3. Map dữ liệu đầy đủ
	dept := model.Department{
		Code:      req.Code,
		NameVN:    req.NameVN,
		NameEN:    req.NameEN,
		NameTW:    req.NameTW,
		Desc:      req.Desc,
		ParentID:  req.ParentID,
		ManagerID: req.ManagerID,
		Level:     &level,
	}

	if err := d.repo.Create(ctx, &dept); err != nil {
		return fmt.Errorf("failed to create department %w", err)
	}
	return nil
}

func (d *departmentService) GetByID(ctx context.Context, id uint64) (*dto.DepartmentResponse, error) {
	dept, err := d.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get by department id %w", err)
	}
	return &dto.DepartmentResponse{
		ID:        dept.ID,
		Code:      dept.Code,
		NameVN:    dept.NameVN,
		NameEN:    dept.NameEN,
		NameTW:    dept.NameTW,
		Desc:      dept.Desc,
		CreatedAt: dept.CreatedAt,
		UpdatedAt: dept.UpdatedAt,
	}, nil
}
func (d *departmentService) Update(ctx context.Context, id uint64, req dto.DepartmentUpdate) error {
	updates := make(map[string]interface{})

	// 1. Mapping các field cơ bản
	if req.NameVN != nil {
		updates["name_vn"] = *req.NameVN
	}
	if req.NameEN != nil {
		updates["name_en"] = *req.NameEN
	}
	if req.NameTW != nil {
		updates["name_tw"] = *req.NameTW
	}
	if req.Desc != nil {
		updates["desc"] = *req.Desc
	}
	if req.ManagerID != nil {
		updates["manager_id"] = *req.ManagerID
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// 2. LOGIC PHỨC TẠP: Xử lý khi thay đổi Parent (Chuyển phòng ban cha)
	if req.ParentID != nil {
		newParentID := *req.ParentID

		// Check 1: Không được làm cha của chính mình
		if newParentID == id {
			return fmt.Errorf("cannot set department as its own parent")
		}

		updates["parent_id"] = req.ParentID

		// Check 2: Tính lại Level
		var newLevel uint8 = 1
		if newParentID != 0 {
			parentDept, err := d.repo.GetByID(ctx, newParentID)
			if err != nil {
				return fmt.Errorf("invalid parent id: %w", err)
			}
			if parentDept.Level != nil {
				newLevel = *parentDept.Level + 1
			}
		}
		updates["level"] = newLevel
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	if err := d.repo.Update(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update department service %w", err)
	}
	return nil
}
func (d *departmentService) Delete(ctx context.Context, id uint64) error {
	return d.repo.Delete(ctx, id)
}
func (d *departmentService) GetAll(ctx context.Context) ([]dto.DepartmentResponse, error) {
	depts, err := d.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all departments %w", err)
	}

	res := make([]dto.DepartmentResponse, 0, len(depts))
	for _, v := range depts {
		res = append(res, dto.DepartmentResponse{
			ID:        v.ID,
			Code:      v.Code,
			NameVN:    v.NameVN,
			NameEN:    v.NameEN,
			NameTW:    v.NameTW,
			Desc:      v.Desc,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}
	return res, nil
}
