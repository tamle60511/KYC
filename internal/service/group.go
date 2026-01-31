package service

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/model"
	"CQS-KYC/internal/repository"
	"context"
	"fmt"
)

type (
	groupService struct {
		repo repository.GroupRepo
	}
	GroupService interface {
		Create(ctx context.Context, g *dto.UserGroupCreateReq) error
		GetByID(ctx context.Context, id uint64) (*dto.UserGroupRes, error)
		Update(ctx context.Context, id uint64, g *dto.UserGroupUpdateReq) error
		Delete(ctx context.Context, id uint64) error
		GetAll(ctx context.Context) ([]dto.UserGroupRes, error)
	}
)

func NewGroupService(repo repository.GroupRepo) GroupService {
	return &groupService{
		repo: repo,
	}
}
func (r *groupService) Create(ctx context.Context, req *dto.UserGroupCreateReq) error {
	// 1. Check Code trùng
	_, err := r.repo.FindByCode(ctx, req.GroupCode)
	if err == nil {
		return fmt.Errorf("group code already exists")
	}

	// 2. Validate Users (Quan trọng)
	// Nếu danh sách user dài, nên dùng Where("id IN ?", ids) để query 1 lần
	for _, m := range req.Member {
		if _, err := r.repo.FindByID(ctx, m.UserID); err != nil {
			return fmt.Errorf("user id %d not found", m.UserID)
		}
	}

	// 3. Map Data
	group := model.UserGroup{
		GroupCode:   req.GroupCode,
		GroupName:   req.GroupName,
		Description: req.Description,
		IsActive:    true,
	}

	// Map Members
	for _, member := range req.Member {
		group.Members = append(group.Members, model.UserGroupMember{
			UserID: member.UserID,
			Role:   member.Role,
			// Không cần set GroupID ở đây, GORM tự set sau khi tạo Group cha
		})
	}

	if err := r.repo.Create(ctx, &group); err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}
	return nil
}

func (r *groupService) GetByID(ctx context.Context, id uint64) (*dto.UserGroupRes, error) {
	wf, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user group by id %w", err)
	}

	groupRes := &dto.UserGroupRes{
		ID:          wf.ID,
		GroupCode:   wf.GroupCode,
		GroupName:   wf.GroupName,
		Description: wf.Description,
		IsActive:    wf.IsActive,
		CreatedAt:   wf.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   wf.UpdatedAt.Format("2006-01-02 15:04:05"),
		Member:      make([]dto.UserGroupMemberRes, 0),
	}
	for _, member := range wf.Members {
		memberRes := dto.UserGroupMemberRes{
			ID:       member.ID,
			GroupID:  member.GroupID,
			UserID:   member.UserID,
			Role:     member.Role,
			JoinedAt: member.JoinedAt.Format("2006-01-02 15:04:05"),
		}
		groupRes.Member = append(groupRes.Member, memberRes)
	}
	return groupRes, nil
}
func (r *groupService) Update(ctx context.Context, id uint64, g *dto.UserGroupUpdateReq) error {
	existingGroup, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user group not found %w", err)
	}

	if g.GroupName != nil {
		existingGroup.GroupName = *g.GroupName
	}
	if g.Description != nil {
		existingGroup.Description = *g.Description
	}
	if g.IsActive != nil {
		existingGroup.IsActive = *g.IsActive
	}

	existingGroup.Members = make([]model.UserGroupMember, 0)
	for _, member := range g.Member {
		groupMember := model.UserGroupMember{
			UserID: member.UserID,
			Role:   member.Role,
		}
		existingGroup.Members = append(existingGroup.Members, groupMember)
	}

	if err := r.repo.Update(ctx, id, existingGroup); err != nil {
		return fmt.Errorf("failed to update user group %w", err)
	}
	return nil
}
func (r *groupService) Delete(ctx context.Context, id uint64) error {
	return r.repo.Delete(ctx, id)
}
func (r *groupService) GetAll(ctx context.Context) ([]dto.UserGroupRes, error) {
	groups, err := r.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all user groups %w", err)
	}

	groupResList := make([]dto.UserGroupRes, 0)
	for _, wf := range groups {
		groupRes := dto.UserGroupRes{
			ID:          wf.ID,
			GroupCode:   wf.GroupCode,
			GroupName:   wf.GroupName,
			Description: wf.Description,
			IsActive:    wf.IsActive,
			CreatedAt:   wf.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   wf.UpdatedAt.Format("2006-01-02 15:04:05"),
			Member:      make([]dto.UserGroupMemberRes, 0),
		}
		for _, member := range wf.Members {
			memberRes := dto.UserGroupMemberRes{
				ID:       member.ID,
				GroupID:  member.GroupID,
				UserID:   member.UserID,
				Role:     member.Role,
				JoinedAt: member.JoinedAt.Format("2006-01-02 15:04:05"),
			}
			groupRes.Member = append(groupRes.Member, memberRes)
		}
		groupResList = append(groupResList, groupRes)
	}
	return groupResList, nil
}
