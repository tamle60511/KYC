package model

import (
	"time"

	"gorm.io/gorm"
)

type Department struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Code string `gorm:"uniqueIndex;size:50;not null" json:"code"`

	// --- Names ---
	NameVN string `json:"name_vn"`
	NameEN string `json:"name_en"`
	NameTW string `json:"name_tw"`
	Desc   string `json:"desc"`

	// --- ERP/Workflow Core Fields ---
	Level    *uint8      `json:"level"`                  // Cấp bậc phòng ban trong hệ thống (Hierarchy)
	ParentID *uint64     `gorm:"index" json:"parent_id"` // Phòng ban cha (Hierarchy)
	Parent   *Department `gorm:"foreignKey:ParentID" json:"parent,omitempty"`

	ManagerID *uint64 `gorm:"index" json:"manager_id"` // Trưởng phòng (Quan trọng cho Workflow)
	// Lưu ý: Không preload User ở đây để tránh Circular Dependency (User->Dept->User)
	// Khi cần hiển thị tên trưởng phòng, ta sẽ join bảng User thủ công hoặc preload có kiểm soát.

	// --- System ---
	IsActive  bool           `gorm:"default:true" json:"is_active"` // Nên có để soft-disable phòng ban giải thể
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft Delete (ERP luôn cần cái này)
}

func (Department) TableName() string {
	return "departments"
}
