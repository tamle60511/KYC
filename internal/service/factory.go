package service

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/model"
	"CQS-KYC/internal/repository"
	"context"
	"fmt"
)

type (
	factoryService struct {
		repo repository.FactoryRepo
	}
	FactoryService interface {
		Create(ctx context.Context, req dto.FactoryCreate) error
		GetByID(ctx context.Context, id uint64) (*dto.FactoryResponse, error)
		Update(ctx context.Context, id uint64, req dto.FactoryUpdate) error
		Delete(ctx context.Context, id uint64) error
		GetList(ctx context.Context) ([]*dto.FactoryResponse, error)
	}
)

func NewFactoryService(repo repository.FactoryRepo) FactoryService {
	return &factoryService{repo: repo}
}

func (s *factoryService) Create(ctx context.Context, req dto.FactoryCreate) error {
	f, err := s.repo.GetByCode(ctx, req.Code)
	if err == nil && f != nil {
		return fmt.Errorf("factory with code %s already exists", req.Code)
	}

	factory := &model.Factory{
		Code:    req.Code,
		Name:    req.Name,
		Address: req.Address,
		TaxCode: req.TaxCode,
	}
	if err := s.repo.Create(ctx, factory); err != nil {
		return err
	}
	return nil
}
func (s *factoryService) GetByID(ctx context.Context, id uint64) (*dto.FactoryResponse, error) {
	factory, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := &dto.FactoryResponse{
		ID:       factory.ID,
		Code:     factory.Code,
		Name:     factory.Name,
		Address:  factory.Address,
		TaxCode:  factory.TaxCode,
		IsActive: factory.IsActive,
	}
	return response, nil
}
func (s *factoryService) Update(ctx context.Context, id uint64, req dto.FactoryUpdate) error {
	updateData := make(map[string]interface{})
	if req.Name != nil {
		updateData["name"] = *req.Name
	}
	if req.Address != nil {
		updateData["address"] = *req.Address
	}
	if req.TaxCode != nil {
		updateData["tax_code"] = *req.TaxCode
	}
	if req.IsActive != nil {
		updateData["is_active"] = *req.IsActive
	}
	if len(updateData) == 0 {
		return fmt.Errorf("no data to update")
	}
	if err := s.repo.Update(ctx, id, updateData); err != nil {
		return err
	}
	return nil
}
func (s *factoryService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
func (s *factoryService) GetList(ctx context.Context) ([]*dto.FactoryResponse, error) {
	factories, err := s.repo.GetList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed go get all factory %w", err)
	}
	var responses []*dto.FactoryResponse
	for _, factory := range factories {
		responses = append(responses, &dto.FactoryResponse{
			ID:       factory.ID,
			Code:     factory.Code,
			Name:     factory.Name,
			Address:  factory.Address,
			TaxCode:  factory.TaxCode,
			IsActive: factory.IsActive,
		})
	}
	return responses, nil
}
