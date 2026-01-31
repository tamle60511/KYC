package model

import "time"

type Factory struct {
	ID       uint64 `gorm:"primaryKey" json:"id"`
	Code     string `gorm:"size:50;not null;uniqueIndex" json:"code"` // uniqueIndex tốt hơn unique thường
	Name     string `gorm:"size:255;not null" json:"name"`
	Address  string `gorm:"size:255" json:"address"`       // Bổ sung
	TaxCode  string `gorm:"size:50" json:"tax_code"`       // Bổ sung
	IsActive bool   `gorm:"default:true" json:"is_active"` // Nên có để tạm khóa nhà máy

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Factory) TableName() string {
	return "factories"
}
