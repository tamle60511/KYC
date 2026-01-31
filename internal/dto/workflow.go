package dto

type WorkflowDefinitionRes struct {
	ID           uint64            `json:"id" uri:"id"`
	ServiceCode  string            `json:"service_code"`
	Operation    string            `json:"operation"`
	WorkflowName string            `json:"workflow_name"`
	Description  string            `json:"description"`
	IsActive     bool              `json:"is_active"`
	Steps        []WorkflowStepRes `json:"steps"`
}

type WorkflowStepRes struct {
	ID                   uint64                      `json:"id"`
	WorkflowDefinitionID uint64                      `json:"workflow_definition_id"`
	StepCode             string                      `json:"step_code"`
	StepName             string                      `json:"step_name"`
	StepOrder            int                         `json:"step_order"`
	RequiredRole         string                      `json:"required_role"`
	Canskip              bool                        `json:"can_skip"`
	CanDelegate          bool                        `json:"can_delegate"`
	RequireComment       bool                        `json:"require_comment"`
	TimeHours            int                         `json:"time_hours"`
	Assisments           []WorkflowStepAssignmentRes `json:"assignments"`
}

type WorkflowStepAssignmentRes struct {
	ID               uint64   `json:"id"`
	StepID           uint64   `json:"step_id"`
	DepartmentIDs    []uint64 `json:"department_ids"`
	AssignedType     string   `json:"assigned_type"`
	AssignedIdentity string   `json:"assigned_identity"`
	Priority         int      `json:"priority"`
	IsActive         bool     `json:"is_active"`
}
type WorkflowDefinitionCreate struct {
	ServiceCode  string                  `json:"service_code" binding:"required"`
	WorkflowName string                  `json:"workflow_name" binding:"required"`
	Operation    string                  `json:"operation"`
	Description  string                  `json:"description"`
	IsActive     bool                    `json:"is_active"`
	Steps        []WorkflowStepCreateReq `json:"steps" binding:"dive"`
}

type WorkflowStepCreateReq struct {
	StepCode       string                            `json:"step_code" binding:"required"`
	StepName       string                            `json:"step_name" binding:"required"`
	StepOrder      int                               `json:"step_order" binding:"required"`
	RequiredRole   string                            `json:"required_role" binding:"required"`
	Canskip        bool                              `json:"can_skip"`
	CanDelegate    bool                              `json:"can_delegate"`
	RequireComment bool                              `json:"require_comment"`
	TimeHours      int                               `json:"time_hours"`
	Assisments     []WorkflowStepAssignmentCreateReq `json:"assignments" binding:"dive"`
}

type WorkflowStepAssignmentCreateReq struct {
	DepartmentIDs    []uint64 `json:"department_ids"`
	AssignedType     string   `json:"assigned_type" binding:"required"`
	AssignedIdentity string   `json:"assigned_identity" binding:"required"`
	Priority         int      `json:"priority"`
	IsActive         bool     `json:"is_active"`
}

type WorkflowDefinitionUpdate struct {
	Operation    *string                  `json:"operation"`
	WorkflowName *string                  `json:"workflow_name" binding:"required"`
	Description  *string                  `json:"description"`
	IsActive     *bool                    `json:"is_active"`
	Steps        []*WorkflowStepCreateReq `json:"steps" binding:"dive"`
}

type WorkflowStepUpdateReq struct {
	StepName       *int                               `json:"step_name"`
	StepOrder      *int                               `json:"step_order"`
	RequiredRole   *string                            `json:"required_role"`
	Canskip        *bool                              `json:"can_skip"`
	CanDelegate    *bool                              `json:"can_delegate"`
	RequireComment *bool                              `json:"require_comment"`
	TimeHours      *int                               `json:"time_hours"`
	Assisments     []*WorkflowStepAssignmentCreateReq `json:"assignments" binding:"dive"`
}

type WorkflowStepAssignmentUpdateReq struct {
	DepartmentIDs    []uint64 `json:"department_ids"`
	AssignedType     *string  `json:"assigned_type"`
	AssignedIdentity *string  `json:"assigned_identity"`
	Priority         *int     `json:"priority"`
	IsActive         *bool    `json:"is_active"`
}
