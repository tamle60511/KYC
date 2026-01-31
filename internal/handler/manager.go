package handler

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/service"
	"CQS-KYC/utils"

	"github.com/gofiber/fiber/v3"
)

type ManagerHandler struct {
	service service.ManagerService
}

func NewManagerHandler(svc service.ManagerService) *ManagerHandler {
	return &ManagerHandler{service: svc}
}

func (h *ManagerHandler) Create(c fiber.Ctx) error {
	var manager dto.ManagerCreate
	if err := c.Bind().Body(&manager); err != nil {
		return utils.BadRequestResponse(c, "invalid request body", err)
	}
	if err := h.service.Create(c.Context(), manager); err != nil {
		return utils.InternalErrorResponse(c, "failed to create manager handler", err)
	}
	return utils.CreatedResponse(c, "create manager success", nil)
}

func (h *ManagerHandler) GetByID(c fiber.Ctx) error {
	var manager dto.ManagerRes
	if err := c.Bind().URI(&manager.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid request id", err)
	}
	managers, err := h.service.GetByID(c.Context(), manager.ID)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get manager by id", err)
	}
	return utils.SuccessResponse(c, "get manager by id success", managers)
}

func (h *ManagerHandler) Update(c fiber.Ctx) error {
	var manager dto.ManagerRes
	if err := c.Bind().URI(&manager.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid request id", err)
	}
	var req dto.ManagerUpdate
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "invalid request body", err)
	}
	if err := h.service.Update(c.Context(), manager.ID, req); err != nil {
		return utils.InternalErrorResponse(c, "failed to update manager", err)
	}
	return utils.SuccessResponse(c, "update manager success", nil)
}

func (h *ManagerHandler) Delete(c fiber.Ctx) error {
	var manager dto.ManagerRes
	if err := c.Bind().URI(&manager.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid request id", err)
	}
	if err := h.service.Delete(c.Context(), manager.ID); err != nil {
		return utils.InternalErrorResponse(c, "failed to delete manager", err)
	}
	return utils.SuccessResponse(c, "delete manager success", nil)
}

func (h *ManagerHandler) GetAll(c fiber.Ctx) error {
	managers, err := h.service.GetAll(c.Context())
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get all managers", err)
	}
	return utils.SuccessResponse(c, "get all managers success", managers)
}

func (h *ManagerHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	managerRouter := router.Group("/managers")
	for _, m := range ms {
		managerRouter.Use(m)
	}

	managerRouter.Post("/", h.Create)
	managerRouter.Get("/:id", h.GetByID)
	managerRouter.Put("/:id", h.Update)
	managerRouter.Delete("/:id", h.Delete)
	managerRouter.Get("/", h.GetAll)
}
