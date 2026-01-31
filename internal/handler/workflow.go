package handler

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/service"
	"CQS-KYC/utils"

	"github.com/gofiber/fiber/v3"
)

type WorkflowHandler struct {
	service service.WorkflowService
}

func NewWorkflowHandler(service service.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{
		service: service,
	}
}

func (h *WorkflowHandler) Create(c fiber.Ctx) error {
	var wfl dto.WorkflowDefinitionCreate
	if err := c.Bind().Body(&wfl); err != nil {
		return utils.BadRequestResponse(c, "invalid request body", err)
	}
	if err := h.service.Create(c.Context(), wfl); err != nil {
		return utils.InternalErrorResponse(c, "failed to create workflow definition handler", err)
	}
	return nil
}

func (h *WorkflowHandler) GetByID(c fiber.Ctx) error {
	var wflID dto.WorkflowDefinitionRes
	if err := c.Bind().URI(&wflID.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid request id", err)
	}
	wlf, err := h.service.GetByID(c.Context(), wflID.ID)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get by workflow definition id", err)
	}
	return utils.SuccessResponse(c, "get by id workflow definition success", wlf)
}

func (h *WorkflowHandler) Update(c fiber.Ctx) error {
	var wflID dto.WorkflowDefinitionRes
	if err := c.Bind().URI(&wflID.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid request id", err)
	}
	var wfl dto.WorkflowDefinitionUpdate
	if err := c.Bind().Body(&wfl); err != nil {
		return utils.BadRequestResponse(c, "invalid request body", err)
	}
	if err := h.service.Update(c.Context(), wflID.ID, wfl); err != nil {
		return utils.InternalErrorResponse(c, "failed to update workflow definition", err)
	}
	return utils.SuccessResponse(c, "update workflow definition success", nil)
}
func (h *WorkflowHandler) Delete(c fiber.Ctx) error {
	var wflID dto.WorkflowDefinitionRes
	if err := c.Bind().URI(&wflID.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid request id", err)
	}
	if err := h.service.Delete(c.Context(), wflID.ID); err != nil {
		return utils.InternalErrorResponse(c, "failed to delete workflow definition", err)
	}
	return utils.SuccessResponse(c, "delete workflow definition success", nil)
}

func (h *WorkflowHandler) GetAll(c fiber.Ctx) error {
	wfls, err := h.service.GetAll(c.Context())
	if err != nil {
		return utils.BadRequestResponse(c, "failed to get all workflow definition", err)
	}
	return utils.SuccessResponse(c, "get all workflow definition success", wfls)
}

func (h *WorkflowHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	wfls := router.Group("/workflow")
	for _, m := range ms {
		wfls.Use(m)
	}
	wfls.Post("/", h.Create)
	wfls.Get("/", h.GetAll)
	wfls.Get("/:id", h.GetByID)
	wfls.Put("/:id", h.Update)
	wfls.Delete("/:id", h.Delete)
}
