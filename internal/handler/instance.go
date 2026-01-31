package handler

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/service"
	"CQS-KYC/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type InstanceHandler struct {
	service service.InstanceService
}

func NewInstanceHandler(service service.InstanceService) *InstanceHandler {
	return &InstanceHandler{service: service}
}

// Helper: Lấy UserID từ Header (Giả lập Authentication)
// Trong thực tế, em sẽ lấy từ JWT Claims: c.Locals("user_id")
func getUserID(c fiber.Ctx) string {
	uid := c.Get("x-user-id")
	if uid == "" {
		return "ADMIN" // Fallback cho dễ test
	}
	return uid
}

// POST /api/workflow/initiate
func (h *InstanceHandler) Initiate(c fiber.Ctx) error {
	var req dto.WorkflowInitiateReq
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid body", err)
	}

	userID := getUserID(c) // Người tạo đơn

	result, err := h.service.Initiate(c.Context(), userID, req)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to initiate workflow", err)
	}

	return utils.SuccessResponse(c, "Workflow initiated", result)
}

// POST /api/workflow/:id/action (Approve/Reject)
func (h *InstanceHandler) ProcessAction(c fiber.Ctx) error {
	idStr := c.Params("id")
	instanceID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid instance ID", err)
	}

	var req dto.WorkflowActionReq
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid body", err)
	}

	userID := getUserID(c)
	// userName := getUserName(c) // Nếu có JWT thì lấy tên thật, tạm thời lấy ID làm tên

	if err := h.service.ProcessAction(c.Context(), instanceID, userID, userID, req); err != nil {
		return utils.InternalErrorResponse(c, "Action failed", err)
	}

	return utils.SuccessResponse(c, "Action processed successfully", nil)
}

// GET /api/workflow/tasks (My Tasks)
func (h *InstanceHandler) GetMyTasks(c fiber.Ctx) error {
	userID := getUserID(c)

	tasks, err := h.service.GetPendingTasks(c.Context(), userID)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get tasks", err)
	}

	return utils.SuccessResponse(c, "Pending tasks retrieved", tasks)
}

// GET /api/workflow/:id/history
func (h *InstanceHandler) GetHistory(c fiber.Ctx) error {
	idStr := c.Params("id")
	instanceID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid instance ID", err)
	}

	history, err := h.service.GetHistory(c.Context(), instanceID)
	if err != nil {
		return utils.InternalErrorResponse(c, "Failed to get history", err)
	}

	return utils.SuccessResponse(c, "History retrieved", history)
}

// Setup Routes
func (h *InstanceHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	instance := router.Group("/instance")
	for _, m := range ms {
		instance.Use(m)
	}
	instance.Post("/initiate", h.Initiate)        // Tạo đơn
	instance.Get("/tasks", h.GetMyTasks)          // Xem việc cần làm (Quan trọng)
	instance.Post("/:id/action", h.ProcessAction) // Duyệt/Hủy
	instance.Get("/:id/history", h.GetHistory)    // Xem lịch sử
}
