package model

import "time"

type User struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	UserCode       string    `gorm:"uniqueIndex;size:50;not null" json:"user_code"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Password       string    `json:"-"`
	DepartmentID   uint64    `gorm:"index" json:"department_id"`
	ManagerID      *uint64   `json:"manager_id"`
	FactoryID      uint64    `gorm:"index" json:"factory_id"`
	PositionID     uint64    `json:"position_id"`
	Role           string    `gorm:"default:user" json:"role"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	SignatureImage string    `json:"signature_image"`
	Subordinates   []User    `gorm:"foreignKey:ManagerID" json:"subordinates"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Position   *Position   `gorm:"foreignKey:PositionID" json:"position,omitempty"`
	Factory    *Factory    `gorm:"foreignKey:FactoryID" json:"factory,omitempty"`
}

func (User) TableName() string {
	return "users"
}
