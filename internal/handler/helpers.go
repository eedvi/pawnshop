package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"pawnshop/internal/service"
)

// handleServiceError handles common service errors and returns appropriate HTTP responses
func handleServiceError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// Check for common service errors and return appropriate status codes
	switch {
	case errors.Is(err, service.ErrBranchNotFound),
		errors.Is(err, service.ErrCustomerNotFound),
		errors.Is(err, service.ErrItemNotFound),
		errors.Is(err, service.ErrUserNotFound),
		errors.Is(err, service.ErrLoanNotFound),
		errors.Is(err, service.ErrPaymentNotFound),
		errors.Is(err, service.ErrSaleNotFound),
		errors.Is(err, service.ErrCategoryNotFound),
		errors.Is(err, service.ErrRoleNotFound),
		errors.Is(err, service.ErrSettingNotFound),
		errors.Is(err, service.ErrTransferNotFound),
		errors.Is(err, service.ErrExpenseNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})

	case errors.Is(err, service.ErrUnauthorized):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})

	case errors.Is(err, service.ErrForbidden):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})

	case errors.Is(err, service.ErrInvalidInput),
		errors.Is(err, service.ErrInvalidStatus),
		errors.Is(err, service.ErrInvalidAmount),
		errors.Is(err, service.ErrInsufficientFunds):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})

	case errors.Is(err, service.ErrDuplicateEntry),
		errors.Is(err, service.ErrConflict):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})

	default:
		// Check for common error message patterns
		errMsg := err.Error()
		if strings.Contains(strings.ToLower(errMsg), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if strings.Contains(strings.ToLower(errMsg), "already exists") ||
			strings.Contains(strings.ToLower(errMsg), "duplicate") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if strings.Contains(strings.ToLower(errMsg), "invalid") ||
			strings.Contains(strings.ToLower(errMsg), "required") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Default to internal server error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
