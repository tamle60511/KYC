package repository

import (
	"CQS-KYC/internal/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type (
	groupRepo struct {
		db *gorm.DB
	}
	GroupRepo interface {
		Create(ctx context.Context, g *model.UserGroup) error
		FindByID(ctx context.Context, id uint64) (*model.UserGroup, error)
		FindByCode(ctx context.Context, code string) (*model.UserGroup, error)
		Update(ctx context.Context, id uint64, g *model.UserGroup) error
		Delete(ctx context.Context, id uint64) error
		GetAll(ctx context.Context) ([]model.UserGroup, error)
		GetGroupsByUserID(ctx context.Context, userID string) ([]string, error)
	}
)

func NewGroupRepo(db *gorm.DB) GroupRepo {
	return &groupRepo{
		db: db,
	}
}

func (r *groupRepo) Create(ctx context.Context, g *model.UserGroup) error {
	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		group := model.UserGroup{
			GroupCode:   g.GroupCode,
			GroupName:   g.GroupName,
			Description: g.Description,
			IsActive:    g.IsActive,
			Members:     make([]model.UserGroupMember, 0),
		}
		for _, member := range g.Members {
			groupMember := model.UserGroupMember{
				GroupID: member.GroupID,
				UserID:  member.UserID,
				Role:    member.Role,
			}
			group.Members = append(group.Members, groupMember)
		}
		if err := tx.Create(&group).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
func (r *groupRepo) FindByID(ctx context.Context, id uint64) (*model.UserGroup, error) {
	var group model.UserGroup
	if err := r.db.WithContext(ctx).Preload("Members").First(&group, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get by id %w", err)
	}
	return &group, nil
}
func (r *groupRepo) FindByCode(ctx context.Context, code string) (*model.UserGroup, error) {
	var group model.UserGroup
	if err := r.db.WithContext(ctx).Preload("Members").First(&group, "group_code = ?", code).Error; err != nil {
		return nil, fmt.Errorf("failed to get by cgroup code %w", err)
	}
	return &group, nil
}
func (r *groupRepo) Update(ctx context.Context, id uint64, g *model.UserGroup) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Update thông tin chung của Group
		if err := tx.Model(&model.UserGroup{}).Where("id = ?", id).Updates(map[string]interface{}{
			"group_name":  g.GroupName,
			"description": g.Description,
			"is_active":   g.IsActive,
		}).Error; err != nil {
			return fmt.Errorf("failed to update group info: %w", err)
		}

		// 2. Xử lý Members (Cách an toàn: Sync)
		// Nếu danh sách member gửi lên rỗng -> Không làm gì (hoặc xóa hết tùy business)
		// Ở đây giữ logic cũ của em: Replace All (Chấp nhận mất JoinedAt)

		if len(g.Members) > 0 {
			// Xóa cũ
			if err := tx.Where("group_id = ?", id).Delete(&model.UserGroupMember{}).Error; err != nil {
				return err
			}

			for i := range g.Members {
				g.Members[i].GroupID = id
			}

			if err := tx.Create(&g.Members).Error; err != nil {
				return fmt.Errorf("failed to update members: %w", err)
			}
		}

		return nil
	})
}
func (r *groupRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&model.UserGroup{}, "id = ?", id).Error
}
func (r *groupRepo) GetAll(ctx context.Context) ([]model.UserGroup, error) {
	var groups []model.UserGroup
	if err := r.db.WithContext(ctx).Preload("Members").Find(&groups).Error; err != nil {
		return nil, fmt.Errorf("failed to get all groups: %w", err)
	}
	return groups, nil
}

func (r *groupRepo) GetGroupsByUserID(ctx context.Context, userID string) ([]string, error) {
	var groupCodes []string

	// Logic:
	// 1. Tìm trong bảng user_group_members những dòng có user_id = ?
	// 2. Join với bảng user_groups để lấy group_code
	// 3. Chỉ lấy nhóm đang Active

	err := r.db.WithContext(ctx).
		Table("user_group_members").
		Select("user_groups.group_code").
		Joins("JOIN user_groups ON user_groups.id = user_group_members.group_id").
		Where("user_group_members.user_id = ? AND user_groups.is_active = ?", userID, true).
		Find(&groupCodes).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get groups for user %s: %w", userID, err)
	}

	return groupCodes, nil
}
