package model

type WorkflowDefinition struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ServiceCode  string         `gorm:"size:100;not null;uniqueIndex:idx_service_code" json:"service_code"`
	Operation    string         `gorm:"size:50;index" json:"operation"` // Tên thao tác PURI05
	WorkflowName string         `json:"workflow_name"`
	Description  string         `json:"description"`
	Version      int            `gorm:"default:1" json:"version"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    int64          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    int64          `gorm:"autoUpdateTime" json:"updated_at"`
	Steps        []WorkflowStep `gorm:"foreignKey:WorkflowDefinitionID;constraint:OnDelete:CASCADE" json:"steps"`
}

func (WorkflowDefinition) TableName() string {
	return "workflow_definitions"
}

type WorkflowStep struct {
	ID                   uint64                   `gorm:"primaryKey;autoIncrement" json:"id"`
	WorkflowDefinitionID uint64                   `gorm:"not null;index:idx_workflow_step" json:"workflow_definition_id"`
	StepCode             string                   `gorm:"size:100;not null" json:"step_code"`
	StepName             string                   `gorm:"not null" json:"step_name"`
	StepOrder            int                      `gorm:"not null" json:"step_order"`
	RequiredRole         string                   `gorm:"size:100;not null" json:"required_role"`
	Canskip              bool                     `gorm:"default:false" json:"can_skip"`
	CanDelegate          bool                     `gorm:"default:false" json:"can_delegate"`
	RequireComment       bool                     `gorm:"default:false" json:"require_comment"`
	TimeHours            int                      `gorm:"default:0" json:"time_hours"`
	CreatedAt            int64                    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            int64                    `gorm:"autoUpdateTime" json:"updated_at"`
	Assignments          []WorkflowStepAssignment `gorm:"foreignKey:StepID;constraint:OnDelete:CASCADE" json:"assignments"`
}

func (WorkflowStep) TableName() string {
	return "workflow_steps"
}

type WorkflowStepAssignment struct {
	ID               uint64   `gorm:"primaryKey;autoIncrement" json:"id"`
	StepID           uint64   `gorm:"not null;index:idx_workflow_step" json:"step_id"`
	DepartmentIDs    []uint64 `gorm:"type:json;serializer:json" json:"department_ids"`
	FactoryID        *uint64  `gorm:"index;size:50" json:"factory_id"`
	AssignedType     string   `gorm:"size:50;not null" json:"assigned_type"`
	AssignedIdentity string   `gorm:"size:100;not null" json:"assigned_identity"`
	Priority         int      `gorm:"default:0" json:"priority"`
	IsActive         bool     `gorm:"default:true" json:"is_active"`
	CreatedAt        int64    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        int64    `gorm:"autoUpdateTime" json:"updated_at"`
}

func (WorkflowStepAssignment) TableName() string {
	return "workflow_step_assignments"
}
