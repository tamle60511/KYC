package dto

import "time"

// 1. Request tạo đơn mới (ERP gửi sang)
type WorkflowInitiateReq struct {
	WorkflowID  uint64      `json:"workflow_id" validate:"required"`
	ServiceCode string      `json:"service_code"`
	DocNum      string      `json:"doc_num" validate:"required"`
	DocType     string      `json:"doc_type"`
	FactoryID   uint64      `json:"factory_id" validate:"required"`
	DeptID      uint64      `json:"dept_id" validate:"required"`
	RequestData interface{} `json:"request_data"`
}

// 2. Request Duyệt/Từ chối
type WorkflowActionReq struct {
	Action  string `json:"action" validate:"required,oneof=APPROVE REJECT"`
	Comment string `json:"comment"`
}

// 3. Response: Danh sách việc cần làm (Task List)
type PendingTaskRes struct {
	TaskID      uint64    `json:"task_id"`
	InstanceID  uint64    `json:"instance_id"`
	DocNum      string    `json:"doc_num"`
	DocType     string    `json:"doc_type"`
	ServiceCode string    `json:"service_code"`
	StepName    string    `json:"step_name"`
	Status      string    `json:"status"`
	ReceivedAt  time.Time `json:"received_at"`
	CreatorID   string    `json:"creator_id"`
}

// 4. Response: Chi tiết lịch sử (History)
type WorkflowLogRes struct {
	StepName  string    `json:"step_name"`
	Action    string    `json:"action"`
	ActorName string    `json:"actor_name"`
	Comment   string    `json:"comment"`
	Time      time.Time `json:"time"`
}
