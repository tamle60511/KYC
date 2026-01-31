package handler

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/service"
	"CQS-KYC/utils"

	"github.com/gofiber/fiber/v3"
)

type DepartmentHandler struct {
	service service.DepartmentService
}

func NewDepartmentHandler(service service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		service: service,
	}
}
func (h *DepartmentHandler) Create(c fiber.Ctx) error {
	var dept dto.DepartmentCreate
	if err := c.Bind().Body(&dept); err != nil {
		return utils.BadRequestResponse(c, "invalid request body", err)
	}
	if err := h.service.Create(c.Context(), dept); err != nil {
		return utils.InternalErrorResponse(c, "failed to create department", err)
	}
	return utils.SuccessResponse(c, "department created successfully", nil)
}

func (h *DepartmentHandler) GetByID(c fiber.Ctx) error {
	var deptID dto.DepartmentID
	if err := c.Bind().URI(&deptID.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid department ID", err)
	}
	dept, err := h.service.GetByID(c.Context(), deptID.ID)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get department", err)
	}
	return utils.SuccessResponse(c, "department retrieved successfully", dept)
}

func (h *DepartmentHandler) Update(c fiber.Ctx) error {
	var deptID dto.DepartmentID
	if err := c.Bind().URI(&deptID.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid department ID", err)
	}
	var deptUpdate dto.DepartmentUpdate
	if err := c.Bind().Body(&deptUpdate); err != nil {
		return utils.BadRequestResponse(c, "invalid request body", err)
	}
	if err := h.service.Update(c.Context(), deptID.ID, deptUpdate); err != nil {
		return utils.InternalErrorResponse(c, "failed to update department", err)
	}
	return utils.SuccessResponse(c, "department updated successfully", nil)
}

func (h *DepartmentHandler) Delete(c fiber.Ctx) error {
	var deptID dto.DepartmentID
	if err := c.Bind().URI(&deptID.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid department ID", err)
	}
	if err := h.service.Delete(c.Context(), deptID.ID); err != nil {
		return utils.InternalErrorResponse(c, "failed to delete department", err)
	}
	return utils.SuccessResponse(c, "department deleted successfully", nil)
}

func (h *DepartmentHandler) GetAll(c fiber.Ctx) error {
	depts, err := h.service.GetAll(c.Context())
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get departments", err)
	}
	return utils.SuccessResponse(c, "departments retrieved successfully", depts)
}

func (h *DepartmentHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	dept := router.Group("/departments")
	for _, m := range ms {
		dept.Use(m)
	}
	dept.Post("/", h.Create)
	dept.Get("/", h.GetAll)
	dept.Get("/:id", h.GetByID)
	dept.Put("/:id", h.Update)
	dept.Delete("/:id", h.Delete)
}
