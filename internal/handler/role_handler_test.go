package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/service"
)

// MockRoleService is a mock implementation of RoleService methods
type MockRoleService struct {
	mock.Mock
}

func (m *MockRoleService) Create(ctx interface{}, input service.CreateRoleInput) (*domain.Role, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleService) Update(ctx interface{}, id int64, input service.UpdateRoleInput) (*domain.Role, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleService) GetByID(ctx interface{}, id int64) (*domain.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleService) List(ctx interface{}) ([]*domain.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Role), args.Error(1)
}

func (m *MockRoleService) Delete(ctx interface{}, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleService) GetAvailablePermissions() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func TestRoleHandler_List_Success(t *testing.T) {
	app := fiber.New()
	mockService := new(MockRoleService)

	roles := []*domain.Role{
		{ID: 1, Name: "admin", DisplayName: "Administrator"},
		{ID: 2, Name: "manager", DisplayName: "Manager"},
	}

	mockService.On("List", mock.Anything).Return(roles, nil)

	app.Get("/api/v1/roles", func(c *fiber.Ctx) error {
		result, err := mockService.List(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"data":    result,
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/roles", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["success"].(bool))
	data := result["data"].([]interface{})
	assert.Len(t, data, 2)

	mockService.AssertExpectations(t)
}

func TestRoleHandler_GetByID_Success(t *testing.T) {
	app := fiber.New()
	mockService := new(MockRoleService)

	role := &domain.Role{
		ID:          1,
		Name:        "admin",
		DisplayName: "Administrator",
		Permissions: []byte(`["*"]`),
	}

	mockService.On("GetByID", mock.Anything, int64(1)).Return(role, nil)

	app.Get("/api/v1/roles/:id", func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")
		result, err := mockService.GetByID(c.Context(), int64(id))
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{
			"success": true,
			"data":    result,
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/roles/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["success"].(bool))

	mockService.AssertExpectations(t)
}

func TestRoleHandler_Create_Success(t *testing.T) {
	app := fiber.New()
	mockService := new(MockRoleService)

	createdRole := &domain.Role{
		ID:          3,
		Name:        "custom-role",
		DisplayName: "Custom Role",
		Permissions: []byte(`["users.read","users.create"]`),
	}

	mockService.On("Create", mock.Anything, mock.AnythingOfType("service.CreateRoleInput")).Return(createdRole, nil)

	app.Post("/api/v1/roles", func(c *fiber.Ctx) error {
		var input service.CreateRoleInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}

		result, err := mockService.Create(c.Context(), input)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(201).JSON(fiber.Map{
			"success": true,
			"data":    result,
		})
	})

	body := map[string]interface{}{
		"name":         "custom-role",
		"display_name": "Custom Role",
		"permissions":  []string{"users.read", "users.create"},
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/roles", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["success"].(bool))

	mockService.AssertExpectations(t)
}

func TestRoleHandler_Update_Success(t *testing.T) {
	app := fiber.New()
	mockService := new(MockRoleService)

	updatedRole := &domain.Role{
		ID:          1,
		Name:        "updated-role",
		DisplayName: "Updated Role",
	}

	mockService.On("Update", mock.Anything, int64(1), mock.AnythingOfType("service.UpdateRoleInput")).Return(updatedRole, nil)

	app.Put("/api/v1/roles/:id", func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")
		var input service.UpdateRoleInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
		}

		result, err := mockService.Update(c.Context(), int64(id), input)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    result,
		})
	})

	body := map[string]string{
		"name":         "updated-role",
		"display_name": "Updated Role",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/api/v1/roles/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestRoleHandler_Delete_Success(t *testing.T) {
	app := fiber.New()
	mockService := new(MockRoleService)

	mockService.On("Delete", mock.Anything, int64(1)).Return(nil)

	app.Delete("/api/v1/roles/:id", func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")
		err := mockService.Delete(c.Context(), int64(id))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "role deleted successfully",
		})
	})

	req := httptest.NewRequest("DELETE", "/api/v1/roles/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestRoleHandler_GetPermissions_Success(t *testing.T) {
	app := fiber.New()
	mockService := new(MockRoleService)

	permissions := []string{
		"users.read", "users.create", "users.update",
		"loans.read", "loans.create",
	}

	mockService.On("GetAvailablePermissions").Return(permissions)

	app.Get("/api/v1/roles/permissions", func(c *fiber.Ctx) error {
		result := mockService.GetAvailablePermissions()
		return c.JSON(fiber.Map{
			"success": true,
			"data":    result,
		})
	})

	req := httptest.NewRequest("GET", "/api/v1/roles/permissions", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["success"].(bool))
	data := result["data"].([]interface{})
	assert.Len(t, data, 5)

	mockService.AssertExpectations(t)
}
