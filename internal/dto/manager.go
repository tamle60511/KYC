package dto

import "time"

type ManagerRes struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	UserID       uint64    `json:"user_id"`
	DepartmentID uint64    `json:"department_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ManagerCreate struct {
	UserID       uint64 `json:"user_id"`
	DepartmentID uint64 `json:"department_id"`
}

type ManagerUpdate struct {
	DepartmentID *uint64 `json:"department_id"`
}
