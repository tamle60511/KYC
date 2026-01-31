package handler

import (
	"CQS-KYC/internal/dto"
	"CQS-KYC/internal/service"
	"CQS-KYC/utils"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{service: svc}
}

func (h *UserHandler) Create(c fiber.Ctx) error {
	var user dto.UserCreate
	if err := c.Bind().Body(&user); err != nil {
		return utils.BadRequestResponse(c, "invalid request body", err)
	}
	file, err := c.FormFile("signature_image")
	if err == nil {
		uploadPath, err := utils.Uploadfile(c, file, "signatures")
		if err != nil {
			return utils.InternalErrorResponse(c, "failed to upload signature image", err)
		}
		user.SignatureImage = uploadPath
	}
	if err := h.service.Create(c.Context(), user); err != nil {
		return utils.InternalErrorResponse(c, "failed to create user handler", err)
	}
	return utils.CreatedResponse(c, "create user success", nil)
}

func (h *UserHandler) GetByID(c fiber.Ctx) error {
	var user dto.UserID
	if err := c.Bind().URI(&user.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid request id", err)
	}
	users, err := h.service.GetByID(c.Context(), user.ID)
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get user by id", err)
	}
	return utils.SuccessResponse(c, "get user by id success", users)
}

func (h *UserHandler) Update(c fiber.Ctx) error {
	var user dto.UserID
	if err := c.Bind().URI(&user.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid request id", err)
	}
	var req dto.UserUpdate
	if err := c.Bind().Body(&req); err != nil {
		return utils.BadRequestResponse(c, "invalid request body", err)
	}
	if err := h.service.Update(c.Context(), user.ID, req); err != nil {
		return utils.InternalErrorResponse(c, "failed to update user", err)
	}
	return utils.SuccessResponse(c, "update user success", nil)
}

func (h *UserHandler) Delete(c fiber.Ctx) error {
	var user dto.UserID
	if err := c.Bind().URI(&user.ID); err != nil {
		return utils.BadRequestResponse(c, "invalid request id", err)
	}
	if err := h.service.Delete(c.Context(), user.ID); err != nil {
		return utils.InternalErrorResponse(c, "failed to delete user", err)
	}
	return utils.SuccessResponse(c, "delete user success", nil)
}

func (h *UserHandler) GetAll(c fiber.Ctx) error {
	users, err := h.service.GetAll(c.Context())
	if err != nil {
		return utils.InternalErrorResponse(c, "failed to get all users", err)
	}
	return utils.SuccessResponse(c, "get all users success", users)
}

func (h *UserHandler) SetupRoutes(router fiber.Router, ms ...fiber.Handler) {
	userRouter := router.Group("/users")
	for _, m := range ms {
		userRouter.Use(m)
	}
	userRouter.Post("/", h.Create)
	userRouter.Get("/", h.GetAll)
	userRouter.Get("/:id", h.GetByID)
	userRouter.Put("/:id", h.Update)
	userRouter.Delete("/:id", h.Delete)
}
