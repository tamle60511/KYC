package handler

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/service"
	"CQS-KYC/utils"

	"github.com/gofiber/fiber/v3"
)

type FactoryHandler struct {
	service service.FactoryService
}

func NewFactoryHandler(svc service.FactoryService) *FactoryHandler {
	return &FactoryHandler{service: svc}
}

func (h *FactoryHandler) Create(c fiber.Ctx) error {
	var req dto.FactoryCreate
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "invalid request body", err)
	}
	if err := h.service.Create(c.Context(), req); err != nil {
		return utils.InternalErrorResponse(c, "failed to create factory", err)
	}
	return utils.CreatedResponse(c, "create factory success", nil)
}

func (h *FactoryHandler) GetByID(c fiber.Ctx) error {
	var factory dto.FactoryID
	if err := c.Bind().URI(&factory.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid request id", err)
	}
	f, err := h.service.GetByID(c.Context(), factory.ID)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get factory by id", err)
	}
	return utils.SuccessResponse(c, "get factory by id success", f)
}

func (h *FactoryHandler) Update(c fiber.Ctx) error {
	var factory dto.FactoryID
	if err := c.Bind().URI(&factory.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid request id", err)
	}
	var req dto.FactoryUpdate
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "invalid request body", err)
	}
	if err := h.service.Update(c.Context(), factory.ID, req); err != nil {
		return utils.InternalErrorResponse(c, "failed to update factory", err)
	}
	return utils.SuccessResponse(c, "update factory success", nil)
}

func (h *FactoryHandler) Delete(c fiber.Ctx) error {
	var factory dto.FactoryID
	if err := c.Bind().URI(&factory.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid request id", err)
	}
	if err := h.service.Delete(c.Context(), factory.ID); err != nil {
		return utils.InternalErrorResponse(c, "failed to delete factory", err)
	}
	return utils.SuccessResponse(c, "delete factory success", nil)
}

func (h *FactoryHandler) GetList(c fiber.Ctx) error {
	factories, err := h.service.GetList(c.Context())
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get factory list", err)
	}
	return utils.SuccessResponse(c, "get factory list success", factories)
}
func (h *FactoryHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	factoryRouter := router.Group("/factories")
	for _, m := range ms {
		factoryRouter.Use(m)
	}

	factoryRouter.Post("/", h.Create)
	factoryRouter.Get("/", h.GetList)
	factoryRouter.Get("/:id", h.GetByID)
	factoryRouter.Put("/:id", h.Update)
	factoryRouter.Delete("/:id", h.Delete)
}
