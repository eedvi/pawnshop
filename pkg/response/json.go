package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Response is the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []FieldError  `json:"details,omitempty"`
}

// FieldError contains field-level validation errors
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Meta contains response metadata
type Meta struct {
	RequestID   string     `json:"request_id,omitempty"`
	Timestamp   time.Time  `json:"timestamp"`
	Pagination  *Pagination `json:"pagination,omitempty"`
}

// Pagination contains pagination info
type Pagination struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
}

// Links for paginated responses
type Links struct {
	First string `json:"first,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
	Last  string `json:"last,omitempty"`
}

// PaginatedResponse is the response for paginated data
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta"`
	Links   *Links      `json:"links,omitempty"`
}

// newMeta creates a new Meta with request ID and timestamp
func newMeta(c *fiber.Ctx) *Meta {
	requestID := c.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	return &Meta{
		RequestID: requestID,
		Timestamp: time.Now(),
	}
}

// OK sends a successful response with status 200
func OK(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Data:    data,
		Meta:    newMeta(c),
	})
}

// Created sends a successful response with status 201
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Data:    data,
		Meta:    newMeta(c),
	})
}

// NoContent sends a successful response with status 204
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Paginated sends a paginated response
func Paginated(c *fiber.Ctx, data interface{}, page, perPage, total int) error {
	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}

	meta := newMeta(c)
	meta.Pagination = &Pagination{
		CurrentPage: page,
		PerPage:     perPage,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Error sends an error response
func Error(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
		Meta: newMeta(c),
	})
}

// ErrorWithDetails sends an error response with field errors
func ErrorWithDetails(c *fiber.Ctx, status int, code, message string, details []FieldError) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: newMeta(c),
	})
}

// Common error responses

// BadRequest sends a 400 error
func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, "BAD_REQUEST", message)
}

// Unauthorized sends a 401 error
func Unauthorized(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Authentication required"
	}
	return Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", message)
}

// Forbidden sends a 403 error
func Forbidden(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "You don't have permission to perform this action"
	}
	return Error(c, fiber.StatusForbidden, "FORBIDDEN", message)
}

// NotFound sends a 404 error
func NotFound(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Resource not found"
	}
	return Error(c, fiber.StatusNotFound, "NOT_FOUND", message)
}

// Conflict sends a 409 error
func Conflict(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusConflict, "CONFLICT", message)
}

// UnprocessableEntity sends a 422 error
func UnprocessableEntity(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", message)
}

// TooManyRequests sends a 429 error
func TooManyRequests(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Too many requests. Please try again later."
	}
	return Error(c, fiber.StatusTooManyRequests, "TOO_MANY_REQUESTS", message)
}

// InternalError sends a 500 error
func InternalError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "An internal server error occurred"
	}

	// Log the error
	requestID := c.Get("X-Request-ID")
	log.Error().
		Str("request_id", requestID).
		Str("method", c.Method()).
		Str("path", c.Path()).
		Str("error_message", message).
		Msg("Internal server error")

	return Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", message)
}

// InternalErrorWithErr sends a 500 error and logs the underlying error
func InternalErrorWithErr(c *fiber.Ctx, err error) error {
	message := "An internal server error occurred"

	// Log the detailed error
	requestID := c.Get("X-Request-ID")
	log.Error().
		Err(err).
		Str("request_id", requestID).
		Str("method", c.Method()).
		Str("path", c.Path()).
		Msg("Internal server error")

	return Error(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", message)
}

// ValidationError sends a 400 error with validation details
func ValidationError(c *fiber.Ctx, errors []FieldError) error {
	return ErrorWithDetails(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", errors)
}
