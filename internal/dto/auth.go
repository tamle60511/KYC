package dto

type LoginRequest struct {
	UserCode string `json:"user_code"`
	Password string `json:"password"`
}

type LoginRes struct {
	Token    string `json:"token"`
	UserID   uint64 `json:"user_id"`
	UserCode string `json:"user"`
	Role     string `json:"role"`
}
