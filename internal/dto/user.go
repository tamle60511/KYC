package dto

import "time"

type UserCreate struct {
	UserCode       string `json:"user_code" validate:"required,min=3,max=6"`
	FullName       string `json:"full_name" validate:"omitempty,min=2,max=100"`
	Email          string `json:"email" validate:"required"`
	Phone          string `json:"phone" validate:"omiempty"`
	Password       string `json:"-" validate:"required,min=6"`
	SignatureImage string `json:"signature_image"`
	DepartmentID   uint64 `json:"department_id" validate:"required"`
	ManagerID      uint64 `json:"manager_id" validate:"omitempty"`
	PositionID     uint64 `json:"position_id" validate:"required"`
	FactoryID      uint64 `json:"factory_id" validate:"required"`
}

type UserUpdate struct {
	FullName       *string `json:"full_name"`
	Email          *string `json:"email"`
	Phone          *string `json:"phone"`
	Password       *string `json:"-"`
	SignatureImage string  `json:"signature_image"`
	DepartmentID   *uint64 `json:"department_id"`
	ManagerID      *uint64 `json:"manager_id"`
	PositionID     *uint64 `json:"position_id"`
	FactoryID      *uint64 `json:"factory_id"`
	IsActive       *bool   `json:"is_active"`
}

type UserRes struct {
	ID             uint64    `json:"id" uri:"id"`
	UserCode       string    `json:"user_code"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Password       string    `json:"-"`
	SignatureImage string    `json:"signature_image"`
	DepartmentID   uint64    `json:"department_id"`
	ManagerID      uint64    `json:"manager_id"`
	PositionID     uint64    `json:"position_id"`
	FactoryID      uint64    `json:"factory_id"`
	Role           string    ` json:"role"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserID struct {
	ID uint64 `json:"id" uri:"id"`
}
