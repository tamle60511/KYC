package model

import "time"

type Manager struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	UserID       uint64    `gorm:"index;not null" json:"user_id"`
	DepartmentID uint64    `gorm:"index;not null" json:"department_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Manager) TableName() string {
	return "managers"
}
