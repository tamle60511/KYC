package model

import "time"

type UserGroup struct {
	ID          uint64            `gorm:"primaryKey" json:"id"`
	GroupCode   string            `gorm:"uniqueIndex;size:50;not null" json:"group_code"`
	GroupName   string            `json:"group_name"`
	Description string            `json:"description"`
	IsActive    bool              `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	Members     []UserGroupMember `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"members"`
}

func (UserGroup) TableName() string {
	return "user_groups"
}

type UserGroupMember struct {
	ID       uint64    `gorm:"primaryKey" json:"id"`
	GroupID  uint64    `gorm:"index;not null" json:"group_id"`
	UserID   uint64    `gorm:"index;not null;uniqueIndex:idx_group_user" json:"user_id"`
	User     User      `gorm:"foreignKey:UserID" json:"user"`
	Role     string    `gorm:"default:member" json:"role"`
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"`
}

func (UserGroupMember) TableName() string {
	return "user_group_members"
}
