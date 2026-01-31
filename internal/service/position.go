package service

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/model"
	"CQS-KYC/internal/repository"
	"context"
	"fmt"
)

type (
	positionService struct {
		repo repository.PositionRepo
	}
	PositionService interface {
		Create(ctx context.Context, req dto.PositionCreate) error
		GetByID(ctx context.Context, id uint64) (*dto.PositionResponse, error)
		Update(ctx context.Context, id uint64, req dto.PositionUpdate) error
		Delete(ctx context.Context, id uint64) error
		GetAll(ctx context.Context) ([]dto.PositionResponse, error)
	}
)

func NewPositionService(repo repository.PositionRepo) PositionService {
	return &positionService{
		repo: repo,
	}
}

func (p *positionService) Create(ctx context.Context, req dto.PositionCreate) error {
	_, err := p.repo.GetByCode(ctx, req.Code)
	if err == nil {
		return fmt.Errorf("position code already exists")
	}

	level := 1
	if req.Level > 0 {
		level = req.Level
	}

	pos := model.Position{
		Code:   req.Code,
		NameVN: req.NameVN,
		NameEN: req.NameEN,
		NameTW: req.NameTW,
		Desc:   req.Desc,
		Level:  level,
	}

	// Lưu ý: repo.Create đã sửa bỏ dấu &
	if err := p.repo.Create(ctx, &pos); err != nil {
		return fmt.Errorf("failed to create position: %w", err)
	}
	return nil
}
func (p *positionService) GetByID(ctx context.Context, id uint64) (*dto.PositionResponse, error) {
	pos, err := p.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get by position id %w", err)
	}
	return &dto.PositionResponse{
		ID:        pos.ID,
		Code:      pos.Code,
		NameVN:    pos.NameVN,
		NameEN:    pos.NameEN,
		NameTW:    pos.NameTW,
		CreatedAt: pos.CreatedAt,
		UpdatedAt: pos.UpdatedAt,
	}, nil
}
func (p *positionService) Update(ctx context.Context, id uint64, req dto.PositionUpdate) error {
	updates := make(map[string]interface{})
	if req.NameVN != nil {
		updates["name_vn"] = *req.NameVN
	}
	if req.NameEN != nil {
		updates["name_en"] = *req.NameEN
	}
	if req.NameTW != nil {
		updates["name_tw"] = *req.NameTW
	}
	if req.Level != nil {
		updates["level"] = *req.Level
	}
	if req.Desc != nil {
		updates["desc"] = *req.Desc
	}
	if len(updates) == 0 {
		return fmt.Errorf("no record update")
	}
	if err := p.repo.Update(ctx, id, updates); err != nil {
		return fmt.Errorf("failed to update position %w", err)
	}
	return nil
}
func (p *positionService) Delete(ctx context.Context, id uint64) error {
	return p.repo.Delete(ctx, id)
}
func (p *positionService) GetAll(ctx context.Context) ([]dto.PositionResponse, error) {
	positions, err := p.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all positions %w", err)
	}
	var res []dto.PositionResponse
	for _, pos := range positions {
		res = append(res, dto.PositionResponse{
			ID:        pos.ID,
			Code:      pos.Code,
			NameVN:    pos.NameVN,
			NameEN:    pos.NameEN,
			NameTW:    pos.NameTW,
			CreatedAt: pos.CreatedAt,
			UpdatedAt: pos.UpdatedAt,
		})
	}
	return res, nil

}
