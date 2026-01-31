package model

import "time"

type Position struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"size:50;uniqueIndex;not null" json:"code"`
	NameVN    string    `gorm:"size:100;not null" json:"name_vn"`
	NameEN    string    `gorm:"size:100;not null" json:"name_en"`
	NameTW    string    `gorm:"size:100;not null" json:"name_tw"`
	Level     int       `gorm:"default:1" json:"level"`
	Desc      string    `gorm:"size:255" json:"desc"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Position) TableName() string {
	return "positions"
}
