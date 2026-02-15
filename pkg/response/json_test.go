package response

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp() *fiber.App {
	return fiber.New()
}

func TestOK(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return OK(c, map[string]string{"key": "value"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.True(t, result.Success)
	assert.NotNil(t, result.Data)
	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Meta)
	assert.NotEmpty(t, result.Meta.RequestID)
}

func TestOK_WithRequestID(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return OK(c, nil)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "test-request-123")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "test-request-123", result.Meta.RequestID)
}

func TestCreated(t *testing.T) {
	app := setupTestApp()
	app.Post("/test", func(c *fiber.Ctx) error {
		return Created(c, map[string]int{"id": 1})
	})

	req := httptest.NewRequest("POST", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.True(t, result.Success)
	assert.NotNil(t, result.Data)
}

func TestNoContent(t *testing.T) {
	app := setupTestApp()
	app.Delete("/test", func(c *fiber.Ctx) error {
		return NoContent(c)
	})

	req := httptest.NewRequest("DELETE", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestPaginated(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		data := []string{"a", "b", "c"}
		return Paginated(c, data, 1, 10, 25)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result PaginatedResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.True(t, result.Success)
	assert.NotNil(t, result.Meta)
	assert.NotNil(t, result.Meta.Pagination)
	assert.Equal(t, 1, result.Meta.Pagination.CurrentPage)
	assert.Equal(t, 10, result.Meta.Pagination.PerPage)
	assert.Equal(t, 25, result.Meta.Pagination.TotalItems)
	assert.Equal(t, 3, result.Meta.Pagination.TotalPages)
}

func TestPaginated_ExactDivision(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return Paginated(c, []string{}, 1, 10, 20)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	var result PaginatedResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, 2, result.Meta.Pagination.TotalPages)
}

func TestError(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return Error(c, 500, "INTERNAL_ERROR", "something went wrong")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.False(t, result.Success)
	assert.NotNil(t, result.Error)
	assert.Equal(t, "INTERNAL_ERROR", result.Error.Code)
	assert.Equal(t, "something went wrong", result.Error.Message)
}

func TestErrorWithDetails(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		details := []FieldError{
			{Field: "email", Message: "Invalid email"},
			{Field: "name", Message: "Required"},
		}
		return ErrorWithDetails(c, 400, "VALIDATION_ERROR", "Validation failed", details)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.False(t, result.Success)
	assert.Len(t, result.Error.Details, 2)
	assert.Equal(t, "email", result.Error.Details[0].Field)
}

func TestBadRequest(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return BadRequest(c, "invalid input")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "BAD_REQUEST", result.Error.Code)
	assert.Equal(t, "invalid input", result.Error.Message)
}

func TestUnauthorized_CustomMessage(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return Unauthorized(c, "token expired")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "UNAUTHORIZED", result.Error.Code)
	assert.Equal(t, "token expired", result.Error.Message)
}

func TestUnauthorized_DefaultMessage(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return Unauthorized(c, "")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "Authentication required", result.Error.Message)
}

func TestForbidden_CustomMessage(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return Forbidden(c, "admin only")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "admin only", result.Error.Message)
}

func TestForbidden_DefaultMessage(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return Forbidden(c, "")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "You don't have permission to perform this action", result.Error.Message)
}

func TestNotFound_CustomMessage(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return NotFound(c, "user not found")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "NOT_FOUND", result.Error.Code)
	assert.Equal(t, "user not found", result.Error.Message)
}

func TestNotFound_DefaultMessage(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return NotFound(c, "")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "Resource not found", result.Error.Message)
}

func TestConflict(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return Conflict(c, "already exists")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 409, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "CONFLICT", result.Error.Code)
}

func TestUnprocessableEntity(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return UnprocessableEntity(c, "cannot process")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 422, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "UNPROCESSABLE_ENTITY", result.Error.Code)
}

func TestTooManyRequests_CustomMessage(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return TooManyRequests(c, "rate limited")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 429, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "rate limited", result.Error.Message)
}

func TestTooManyRequests_DefaultMessage(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return TooManyRequests(c, "")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "Too many requests. Please try again later.", result.Error.Message)
}

func TestInternalError_CustomMessage(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return InternalError(c, "db connection failed")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "INTERNAL_ERROR", result.Error.Code)
	assert.Equal(t, "db connection failed", result.Error.Message)
}

func TestInternalError_DefaultMessage(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return InternalError(c, "")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "An internal server error occurred", result.Error.Message)
}

func TestValidationError(t *testing.T) {
	app := setupTestApp()
	app.Get("/test", func(c *fiber.Ctx) error {
		return ValidationError(c, []FieldError{
			{Field: "name", Message: "required"},
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	var result Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	assert.Equal(t, "VALIDATION_ERROR", result.Error.Code)
	assert.Equal(t, "Validation failed", result.Error.Message)
	assert.Len(t, result.Error.Details, 1)
}
