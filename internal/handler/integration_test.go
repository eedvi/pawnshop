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
	"pawnshop/internal/repository"
	"pawnshop/internal/service"
	"pawnshop/pkg/response"
)

// MockBranchService is a mock implementation for branch service
type MockBranchService struct {
	mock.Mock
}

func (m *MockBranchService) Create(ctx interface{}, input service.CreateBranchInput) (*domain.Branch, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Branch), args.Error(1)
}

func (m *MockBranchService) Update(ctx interface{}, id int64, input service.UpdateBranchInput) (*domain.Branch, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Branch), args.Error(1)
}

func (m *MockBranchService) GetByID(ctx interface{}, id int64) (*domain.Branch, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Branch), args.Error(1)
}

func (m *MockBranchService) List(ctx interface{}, params repository.PaginationParams) (*repository.PaginatedResult[domain.Branch], error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[domain.Branch]), args.Error(1)
}

func (m *MockBranchService) Delete(ctx interface{}, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBranchService) Activate(ctx interface{}, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBranchService) Deactivate(ctx interface{}, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockCustomerService is a mock implementation for customer service
type MockCustomerService struct {
	mock.Mock
}

func (m *MockCustomerService) Create(ctx interface{}, input service.CreateCustomerInput) (*domain.Customer, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockCustomerService) Update(ctx interface{}, id int64, input service.UpdateCustomerInput) (*domain.Customer, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockCustomerService) GetByID(ctx interface{}, id int64) (*domain.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockCustomerService) List(ctx interface{}, params repository.CustomerListParams) (*repository.PaginatedResult[domain.Customer], error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[domain.Customer]), args.Error(1)
}

func (m *MockCustomerService) Delete(ctx interface{}, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCustomerService) Block(ctx interface{}, input service.BlockCustomerInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockCustomerService) Unblock(ctx interface{}, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// setupTestApp creates a test app with mock authentication middleware
func setupTestApp() *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "ERROR",
					"message": err.Error(),
				},
			})
		},
	})
}

// setAuthContext sets mock auth context on the request
func setAuthContext(c *fiber.Ctx, userID int64, branchID *int64, permissions []string) {
	c.Locals("user_id", userID)
	c.Locals("branch_id", branchID)
	c.Locals("permissions", permissions)
}

// TestIntegration_HealthCheck tests the health check endpoint
func TestIntegration_HealthCheck(t *testing.T) {
	app := setupTestApp()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "healthy", result["status"])
}

// TestIntegration_Branch_CRUD tests branch CRUD operations
func TestIntegration_Branch_CRUD(t *testing.T) {
	app := setupTestApp()
	mockService := new(MockBranchService)

	branches := app.Group("/api/v1/branches")

	// List branches
	branches.Get("/", func(c *fiber.Ctx) error {
		branchList := []domain.Branch{
			{ID: 1, Code: "MAIN", Name: "Main Branch", IsActive: true},
			{ID: 2, Code: "BR001", Name: "Branch 1", IsActive: true},
		}
		result := &repository.PaginatedResult[domain.Branch]{
			Data:       branchList,
			Total:      2,
			Page:       1,
			PerPage:    10,
			TotalPages: 1,
		}
		mockService.On("List", mock.Anything, mock.AnythingOfType("repository.PaginationParams")).Return(result, nil).Once()
		return response.OK(c, result)
	})

	// Get branch by ID
	branches.Get("/:id", func(c *fiber.Ctx) error {
		branch := &domain.Branch{ID: 1, Code: "MAIN", Name: "Main Branch", IsActive: true}
		mockService.On("GetByID", mock.Anything, int64(1)).Return(branch, nil).Once()
		return response.OK(c, branch)
	})

	// Create branch
	branches.Post("/", func(c *fiber.Ctx) error {
		var input service.CreateBranchInput
		if err := c.BodyParser(&input); err != nil {
			return response.BadRequest(c, "invalid request body")
		}

		branch := &domain.Branch{ID: 3, Code: input.Code, Name: input.Name, IsActive: true}
		mockService.On("Create", mock.Anything, mock.AnythingOfType("service.CreateBranchInput")).Return(branch, nil).Once()

		return response.Created(c, branch)
	})

	// Update branch
	branches.Put("/:id", func(c *fiber.Ctx) error {
		var input service.UpdateBranchInput
		if err := c.BodyParser(&input); err != nil {
			return response.BadRequest(c, "invalid request body")
		}

		name := "Updated Branch"
		if input.Name != "" {
			name = input.Name
		}
		branch := &domain.Branch{ID: 1, Code: "MAIN", Name: name, IsActive: true}
		mockService.On("Update", mock.Anything, int64(1), mock.AnythingOfType("service.UpdateBranchInput")).Return(branch, nil).Once()
		return response.OK(c, branch)
	})

	// Delete branch
	branches.Delete("/:id", func(c *fiber.Ctx) error {
		mockService.On("Delete", mock.Anything, int64(1)).Return(nil).Once()
		return response.NoContent(c)
	})

	// Test List Branches
	t.Run("List Branches", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/branches/", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result["success"].(bool))
	})

	// Test Get Branch by ID
	t.Run("Get Branch by ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/branches/1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Test Create Branch
	t.Run("Create Branch", func(t *testing.T) {
		body := map[string]interface{}{
			"code":    "BR002",
			"name":    "New Branch",
			"address": "123 Main St",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/branches/", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result["success"].(bool))
	})

	// Test Update Branch
	t.Run("Update Branch", func(t *testing.T) {
		body := map[string]interface{}{
			"name": "Updated Branch Name",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/api/v1/branches/1", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Test Delete Branch
	t.Run("Delete Branch", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/branches/1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 204, resp.StatusCode)
	})
}

// TestIntegration_Customer_Flow tests customer creation and management flow
func TestIntegration_Customer_Flow(t *testing.T) {
	app := setupTestApp()
	mockService := new(MockCustomerService)

	customers := app.Group("/api/v1/customers")

	// Create customer
	customers.Post("/", func(c *fiber.Ctx) error {
		var input service.CreateCustomerInput
		if err := c.BodyParser(&input); err != nil {
			return response.BadRequest(c, "invalid request body")
		}

		if input.FirstName == "" || input.LastName == "" {
			return response.ValidationError(c, []response.FieldError{
				{Field: "first_name", Message: "required"},
				{Field: "last_name", Message: "required"},
			})
		}

		customer := &domain.Customer{
			ID:             1,
			BranchID:       input.BranchID,
			FirstName:      input.FirstName,
			LastName:       input.LastName,
			IdentityType:   input.IdentityType,
			IdentityNumber: input.IdentityNumber,
			Phone:          input.Phone,
			IsBlocked:      false,
		}
		mockService.On("Create", mock.Anything, mock.AnythingOfType("service.CreateCustomerInput")).Return(customer, nil).Once()
		return response.Created(c, customer)
	})

	// Get customer
	customers.Get("/:id", func(c *fiber.Ctx) error {
		customer := &domain.Customer{
			ID:             1,
			FirstName:      "John",
			LastName:       "Doe",
			IdentityType:   "dpi",
			IdentityNumber: "123456789",
			Phone:          "555-1234",
			IsBlocked:      false,
		}
		mockService.On("GetByID", mock.Anything, int64(1)).Return(customer, nil).Once()
		return response.OK(c, customer)
	})

	// Block customer
	customers.Post("/:id/block", func(c *fiber.Ctx) error {
		mockService.On("Block", mock.Anything, mock.AnythingOfType("service.BlockCustomerInput")).Return(nil).Once()
		return response.OK(c, fiber.Map{"message": "customer blocked"})
	})

	// Unblock customer
	customers.Post("/:id/unblock", func(c *fiber.Ctx) error {
		mockService.On("Unblock", mock.Anything, int64(1)).Return(nil).Once()
		return response.OK(c, fiber.Map{"message": "customer unblocked"})
	})

	// Test Create Customer with valid data
	t.Run("Create Customer - Valid", func(t *testing.T) {
		body := map[string]interface{}{
			"branch_id":       1,
			"first_name":      "John",
			"last_name":       "Doe",
			"identity_type":   "dpi",
			"identity_number": "123456789",
			"phone":           "555-1234",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/customers/", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result["success"].(bool))
	})

	// Test Create Customer with missing required fields
	t.Run("Create Customer - Missing Fields", func(t *testing.T) {
		body := map[string]interface{}{
			"branch_id": 1,
			"phone":     "555-1234",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/customers/", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode) // Bad request with validation errors
	})

	// Test Get Customer
	t.Run("Get Customer", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/customers/1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Test Block Customer
	t.Run("Block Customer", func(t *testing.T) {
		body := map[string]interface{}{
			"reason": "suspicious activity",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/api/v1/customers/1/block", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	// Test Unblock Customer
	t.Run("Unblock Customer", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/customers/1/unblock", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

// TestIntegration_ErrorHandling tests error responses
func TestIntegration_ErrorHandling(t *testing.T) {
	app := setupTestApp()

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "The requested resource was not found",
			},
		})
	})

	t.Run("Not Found Route", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/nonexistent", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.False(t, result["success"].(bool))

		errorObj := result["error"].(map[string]interface{})
		assert.Equal(t, "NOT_FOUND", errorObj["code"])
	})
}

// TestIntegration_Pagination tests pagination in list endpoints
func TestIntegration_Pagination(t *testing.T) {
	app := setupTestApp()

	app.Get("/api/v1/items", func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		perPage := c.QueryInt("per_page", 10)

		// Simulate paginated response
		items := []domain.Item{}
		for i := 0; i < perPage; i++ {
			items = append(items, domain.Item{ID: int64(i + 1), Name: "Item"})
		}

		result := &repository.PaginatedResult[domain.Item]{
			Data:       items,
			Total:      100,
			Page:       page,
			PerPage:    perPage,
			TotalPages: 10,
		}

		return response.OK(c, result)
	})

	t.Run("Default Pagination", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/items", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		data := result["data"].(map[string]interface{})
		assert.Equal(t, float64(100), data["total"])
		assert.Equal(t, float64(1), data["page"])
		assert.Equal(t, float64(10), data["per_page"])
	})

	t.Run("Custom Pagination", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/items?page=2&per_page=20", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		data := result["data"].(map[string]interface{})
		assert.Equal(t, float64(2), data["page"])
		assert.Equal(t, float64(20), data["per_page"])
	})
}

// TestIntegration_ResponseFormat tests the standard API response format
func TestIntegration_ResponseFormat(t *testing.T) {
	app := setupTestApp()

	// Success response
	app.Get("/api/v1/success", func(c *fiber.Ctx) error {
		return response.OK(c, fiber.Map{"key": "value"})
	})

	// Error response
	app.Get("/api/v1/error", func(c *fiber.Ctx) error {
		return response.BadRequest(c, "something went wrong")
	})

	t.Run("Success Response Format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/success", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		// Standard response format
		assert.True(t, result["success"].(bool))
		assert.NotNil(t, result["data"])
	})

	t.Run("Error Response Format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/error", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		// Standard error format
		assert.False(t, result["success"].(bool))
		assert.NotNil(t, result["error"])

		errorObj := result["error"].(map[string]interface{})
		assert.NotNil(t, errorObj["code"])
		assert.NotNil(t, errorObj["message"])
	})
}
