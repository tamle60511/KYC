package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type WorkflowInstance struct {
	ID         uint64 `gorm:"primaryKey" json:"id"`
	WorkflowID uint64 `gorm:"index;not null" json:"workflow_id"` // ID của quy trình mẹ

	// --- THÔNG TIN TỪ ERP (BUSINESS KEY) ---
	DocNum      string `gorm:"index;size:50;not null" json:"doc_num"`  // Mã đơn: PO-2024-001
	DocType     string `gorm:"index;size:50;not null" json:"doc_type"` // Loại đơn: 310
	ServiceCode string `gorm:"index;size:50" json:"service_code"`      // FormId: WEITWPURI05MDT1

	// --- SNAPSHOT CONTEXT (QUAN TRỌNG ĐỂ ROUTING) ---
	// Lưu chết giá trị lúc tạo đơn. Dù sau này user có chuyển phòng ban, đơn cũ vẫn phải giữ nguyên context.
	FactoryID    uint64 `gorm:"index;size:50" json:"factory_id"`    // VD: VN01
	DepartmentID uint64 `gorm:"index;size:50" json:"department_id"` // VD: IT, ACC

	// --- QUẢN LÝ TRẠNG THÁI ---
	CurrentStep int `gorm:"default:1" json:"current_step"` // Đang ở bước mấy (1, 2, 3...)
	TotalSteps  int `gorm:"default:0" json:"total_steps"`  // Tổng số bước của quy trình này

	// Trạng thái nội bộ (System Status): :NEW, IN_PROGRESS, APPROVED, REJECTED
	Status      string         `gorm:"index;size:20;default:'IN_PROGRESS'" json:"status"`
	RequestData datatypes.JSON `gorm:"type:jsonb" json:"request_data"`
	// --- METADATA ---
	CreatorID   string     `gorm:"size:50;not null" json:"creator_id"` // UserID người tạo trên ERP
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"` // Null nếu chưa xong

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Workflow *WorkflowDefinition `gorm:"foreignKey:WorkflowID" json:"workflow,omitempty"`
	Tasks    []WorkflowTask      `gorm:"foreignKey:InstanceID;constraint:OnDelete:CASCADE" json:"tasks,omitempty"`
	Logs     []WorkflowLog       `gorm:"foreignKey:InstanceID;constraint:OnDelete:CASCADE" json:"logs,omitempty"`
}

type WorkflowTask struct {
	ID         uint64 `gorm:"primaryKey" json:"id"`
	InstanceID uint64 `gorm:"index;not null" json:"instance_id"`

	StepID    uint64 `gorm:"not null" json:"step_id"`
	StepOrder int    `json:"step_order"`
	StepName  string `gorm:"size:100" json:"step_name"` // Cache tên bước để hiển thị cho nhanh

	// --- NGƯỜI ĐƯỢC GIAO VIỆC ---
	// Logic: Hệ thống resolve từ Rule -> Ra Group hoặc User cụ thể -> Lưu vào đây
	AssignedTo string `gorm:"index;size:50;not null" json:"assigned_to"` // UserID hoặc GroupCode
	IsGroup    bool   `gorm:"default:false" json:"is_group"`             // True = Gán cho cả nhóm

	Status  string     `gorm:"size:20;default:'PENDING';index" json:"status"` // PENDING, DONE
	DueDate *time.Time `json:"due_date"`                                      // Tính toán từ TimeoutHours

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relation để join lấy thông tin đơn hàng
	Instance *WorkflowInstance `gorm:"foreignKey:InstanceID" json:"instance,omitempty"`
}

type WorkflowLog struct {
	ID         uint64 `gorm:"primaryKey" json:"id"`
	InstanceID uint64 `gorm:"index;not null" json:"instance_id"`

	// --- THÔNG TIN BƯỚC ---
	StepOrder int    `json:"step_order"`
	StepName  string `json:"step_name"`

	// --- HÀNH ĐỘNG & NGƯỜI DÙNG ---
	Action    string `gorm:"size:50;not null" json:"action"` // APPROVE, REJECT, SUBMIT
	ActorID   string `gorm:"index;size:50;not null" json:"actor_id"`
	ActorName string `gorm:"size:100" json:"actor_name"`
	Comment   string `gorm:"type:text" json:"comment"`

	// --- CÁC TRƯỜNG CHỮ KÝ SỐ (TÍCH HỢP VÀO ĐÂY) ---
	// Thay vì bảng riêng, ta lưu thẳng Hash vào Log
	SignatureHash    string `gorm:"size:255" json:"signature_hash"`     // HMAC Hash
	DataSnapshotHash string `gorm:"size:255" json:"data_snapshot_hash"` // Hash nội dung đơn lúc ký
	SignedTimestamp  int64  `gorm:"not null" json:"signed_timestamp"`   // UnixNano time
	IPAddress        string `gorm:"size:50" json:"ip_address"`
	DeviceInfo       string `gorm:"size:255" json:"device_info"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

const (
	STATUS_NEW         = "NEW"
	STATUS_IN_PROGRESS = "IN_PROGRESS"
	STATUS_APPROVED    = "APPROVED"
	STATUS_REJECTED    = "REJECTED"
	STATUS_CANCELLED   = "CANCELLED"
)

// Hành động trong Log
const (
	ACTION_SUBMIT  = "SUBMIT"  // Gửi đơn
	ACTION_APPROVE = "APPROVE" // Duyệt
	ACTION_REJECT  = "REJECT"  // Từ chối (Kết thúc luôn)
	ACTION_RETURN  = "RETURN"  // Trả về bước trước (Hoặc trả về đầu)
	ACTION_CANCEL  = "CANCEL"  // Hủy đơn
)
