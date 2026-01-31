package dto

import "time"

type DepartmentCreate struct {
	Code      string  `json:"code"`
	NameVN    string  `json:"name_vn"`
	NameEN    string  `json:"name_en"`
	NameTW    string  `json:"name_tw"`
	ParentID  *uint64 `json:"parent_id"`
	ManagerID *uint64 `json:"manager_id"`
	Desc      string  `json:"desc"`
}

type DepartmentUpdate struct {
	NameVN    *string `json:"name_vn"`
	NameEN    *string `json:"name_en"`
	NameTW    *string `json:"name_tw"`
	ParentID  *uint64 `json:"parent_id"`
	ManagerID *uint64 `json:"manager_id"`
	Desc      *string `json:"desc"`
	IsActive  *bool   `json:"is_active"`
}

type DepartmentResponse struct {
	ID        uint64    `json:"id"`
	Code      string    `json:"code"`
	NameVN    string    `json:"name_vn"`
	NameEN    string    `json:"name_en"`
	NameTW    string    `json:"name_tw"`
	ParentID  *uint64   `json:"parent_id"`
	ManagerID *uint64   `json:"manager_id"`
	Level     uint8     `json:"level"`
	Desc      string    `json:"desc"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DepartmentID struct {
	ID uint64 `json:"id" uri:"id" binding:"required"`
}
