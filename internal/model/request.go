package model

import (
	"time"

	"gorm.io/datatypes"
)

type Request struct {
	ID                 uint64         `gorm:"primaryKey"`
	ServiceName        string         `gorm:"size:100;index"`
	CompanyID          string         `gorm:"size:50;index"`
	Operation          string         `gorm:"size:50;index"`
	DocType            string         `gorm:"size:20;index:idx_doc_unique,unique"`
	DocNum             string         `gorm:"size:50;index:idx_doc_unique,unique"`
	CreatorID          string         `gorm:"index;size:50"`
	Status             string         `gorm:"index"`
	SSLProtocal        string         `gorm:"size:10"`
	WorkflowInstanceID uint64         `gorm:"index"`
	Detail             datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt          time.Time
}

func (Request) TableName() string {
	return "requests"
}
