package handler

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/service"
	"CQS-KYC/utils"

	"github.com/gofiber/fiber/v3"
)

type PositionHandler struct {
	service service.PositionService
}

func NewPositionHandler(service service.PositionService) *PositionHandler {
	return &PositionHandler{
		service: service,
	}
}
func (h *PositionHandler) Create(c fiber.Ctx) error {
	var pos dto.PositionCreate
	if err := c.Bind().Body(&pos); err != nil {
		return utils.BadRequestResponse(c, "invalid request body", err)
	}
	if err := h.service.Create(c.Context(), pos); err != nil {
		return utils.InternalErrorResponse(c, "failed to create position", err)
	}
	return utils.SuccessResponse(c, "position created successfully", nil)
}

func (h *PositionHandler) GetByID(c fiber.Ctx) error {
	var posID dto.PositionID
	if err := c.Bind().URI(&posID.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid position ID", err)
	}
	pos, err := h.service.GetByID(c.Context(), posID.ID)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get position", err)
	}
	return utils.SuccessResponse(c, "position retrieved successfully", pos)
}

func (h *PositionHandler) Update(c fiber.Ctx) error {
	var posID dto.PositionID
	if err := c.Bind().URI(&posID.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid position ID", err)
	}
	var posUpdate dto.PositionUpdate
	if err := c.Bind().Body(&posUpdate); err != nil {
		return utils.BadRequestResponse(c, "invalid request body", err)
	}
	if err := h.service.Update(c.Context(), posID.ID, posUpdate); err != nil {
		return utils.InternalErrorResponse(c, "failed to update position", err)
	}
	return utils.SuccessResponse(c, "position updated successfully", nil)
}

func (h *PositionHandler) Delete(c fiber.Ctx) error {
	var posID dto.PositionID
	if err := c.Bind().URI(&posID.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid position ID", err)
	}
	if err := h.service.Delete(c.Context(), posID.ID); err != nil {
		return utils.InternalErrorResponse(c, "failed to delete position", err)
	}
	return utils.SuccessResponse(c, "position deleted successfully", nil)
}

func (h *PositionHandler) GetAll(c fiber.Ctx) error {
	positions, err := h.service.GetAll(c.Context())
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get positions", err)
	}
	return utils.SuccessResponse(c, "positions retrieved successfully", positions)
}

func (h *PositionHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	positionRouter := router.Group("/positions")
	for _, m := range ms {
		positionRouter.Use(m)
	}
	positionRouter.Post("/", h.Create)
	positionRouter.Get("/:id", h.GetByID)
	positionRouter.Put("/:id", h.Update)
	positionRouter.Delete("/:id", h.Delete)
	positionRouter.Get("/", h.GetAll)
}
