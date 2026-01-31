package handler

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/service"
	"CQS-KYC/utils"

	"github.com/gofiber/fiber/v3"
)

type GroupHandler struct {
	service service.GroupService
}

func NewGroupHandler(service service.GroupService) *GroupHandler {
	return &GroupHandler{
		service: service,
	}
}

func (h *GroupHandler) Create(c fiber.Ctx) error {
	var group dto.UserGroupCreateReq
	if err := c.Bind().Body(&group); err != nil {
		return utils.BadRequestResponse(c, "invalid request body", err)
	}
	if err := h.service.Create(c.Context(), &group); err != nil {
		return utils.InternalErrorResponse(c, "failed to create group", err)
	}
	return utils.SuccessResponse(c, "group created successfully", nil)
}

func (h *GroupHandler) GetByID(c fiber.Ctx) error {
	var groupID dto.UserGroupID
	if err := c.Bind().URI(&groupID.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid group ID", err)
	}
	group, err := h.service.GetByID(c.Context(), groupID.ID)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get group", err)
	}
	return utils.SuccessResponse(c, "group retrieved successfully", group)
}

func (h *GroupHandler) Update(c fiber.Ctx) error {
	var groupID dto.UserGroupID
	if err := c.Bind().URI(&groupID.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid group ID", err)
	}
	var groupUpdate dto.UserGroupUpdateReq
	if err := c.Bind().Body(&groupUpdate); err != nil {
		return utils.BadRequestResponse(c, "invalid request body", err)
	}
	if err := h.service.Update(c.Context(), groupID.ID, &groupUpdate); err != nil {
		return utils.InternalErrorResponse(c, "failed to update group", err)
	}
	return utils.SuccessResponse(c, "group updated successfully", nil)
}

func (h *GroupHandler) Delete(c fiber.Ctx) error {
	var groupID dto.UserGroupID
	if err := c.Bind().URI(&groupID.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid group ID", err)
	}
	if err := h.service.Delete(c.Context(), groupID.ID); err != nil {
		return utils.InternalErrorResponse(c, "failed to delete group", err)
	}
	return utils.SuccessResponse(c, "group deleted successfully", nil)
}

func (h *GroupHandler) GetAll(c fiber.Ctx) error {
	groups, err := h.service.GetAll(c.Context())
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get all groups", err)
	}
	return utils.SuccessResponse(c, "get all groups success", groups)
}
func (h *GroupHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	groupRouter := router.Group("/groups")
	for _, err := range ms {
		groupRouter.Use(err)
	}
	groupRouter.Post("/", h.Create)
	groupRouter.Get("/:id", h.GetByID)
	groupRouter.Put("/:id", h.Update)
	groupRouter.Delete("/:id", h.Delete)
	groupRouter.Get("/", h.GetAll)
}
