package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
)

// APIResponse represents standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends successful response
func SuccessResponse(c fiber.Ctx, message string, data interface{}) error {
	return c.JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c fiber.Ctx, statusCode int, message string, err interface{}) error {
	errorStr := ""
	if err != nil {
		if e, ok := err.(error); ok {
			errorStr = e.Error()
		} else if str, ok := err.(string); ok {
			errorStr = str
		} else {
			errorStr = fmt.Sprintf("%v", err)
		}
	}

	return c.Status(statusCode).JSON(APIResponse{
		Success: false,
		Message: message,
		Error:   errorStr,
	})
}

// CreatedResponse sends created response (201)
func CreatedResponse(c fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// NoContentResponse sends no content response (204)
func NoContentResponse(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// BadRequestResponse sends bad request response (400)
func BadRequestResponse(c fiber.Ctx, message string, error interface{}) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message, error)
}

// UnauthorizedResponse sends unauthorized response (401)
func UnauthorizedResponse(c fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, message, nil)
}

// ForbiddenResponse sends forbidden response (403)
func ForbiddenResponse(c fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, message, nil)
}

// NotFoundResponse sends not found response (404)
func NotFoundResponse(c fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, message, nil)
}

// InternalErrorResponse sends internal server error response (500)
func InternalErrorResponse(c fiber.Ctx, message string, error interface{}) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, message, error)
}
