package dto

type UserGroupRes struct {
	ID          uint64               `json:"id"`
	GroupCode   string               `json:"group_code"`
	GroupName   string               `json:"group_name"`
	Description string               `json:"description"`
	IsActive    bool                 `json:"is_active"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
	Member      []UserGroupMemberRes `json:"members,omitempty"`
}

type UserGroupMemberRes struct {
	ID       uint64 `json:"id"`
	GroupID  uint64 `json:"group_id"`
	UserID   uint64 `json:"user_id"`
	Role     string `json:"role"`
	JoinedAt string `json:"joined_at"`
}

type UserGroupCreateReq struct {
	GroupCode   string                     `json:"group_code" binding:"required"`
	GroupName   string                     `json:"group_name" binding:"required"`
	Description string                     `json:"description"`
	Member      []UserGroupMemberCreateReq `json:"members,omitempty"`
}

type UserGroupMemberCreateReq struct {
	GroupID uint64 `json:"group_id"`
	UserID  uint64 `json:"user_id"`
	Role    string `json:"role"`
}

type UserGroupUpdateReq struct {
	GroupName   *string                    `json:"group_name"`
	Description *string                    `json:"description"`
	IsActive    *bool                      `json:"is_active"`
	Member      []UserGroupMemberUpdateReq `json:"members,omitempty"`
}

type UserGroupMemberUpdateReq struct {
	GroupID uint64 `json:"group_id"`
	UserID  uint64 `json:"user_id"`
	Role    string `json:"role"`
}

type UserGroupID struct {
	ID uint64 `json:"id" uri:"id" binding:"required"`
}
