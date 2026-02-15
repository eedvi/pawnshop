package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"pawnshop/internal/middleware"
	"pawnshop/internal/repository"
	"pawnshop/internal/service"
	"pawnshop/pkg/response"
	"pawnshop/pkg/validator"
)

// UserHandler handles user endpoints
type UserHandler struct {
	userService *service.UserService
	auditLogger *middleware.AuditLogger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *service.UserService, auditLogger *middleware.AuditLogger) *UserHandler {
	return &UserHandler{userService: userService, auditLogger: auditLogger}
}

// Create handles user creation
// @Summary Create User
// @Description Create a new user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body service.CreateUserInput true "User data"
// @Success 201 {object} response.Response{data=domain.UserPublic}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/users [post]
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var input service.CreateUserInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	// Validate input
	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	// Create user
	user, err := h.userService.Create(c.Context(), input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Log create event
	if h.auditLogger != nil {
		h.auditLogger.LogCreate(c, "user", user.ID, user)
	}

	return response.Created(c, user)
}

// GetByID handles getting a user by ID
// @Summary Get User
// @Description Get a user by ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response{data=domain.UserPublic}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID format")
	}

	user, err := h.userService.GetByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	return response.OK(c, user)
}

// List handles listing users
// @Summary List Users
// @Description Get a paginated list of users
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(20)
// @Param branch_id query int false "Filter by branch ID"
// @Param role_id query int false "Filter by role ID"
// @Param is_active query bool false "Filter by active status"
// @Param search query string false "Search by name or email"
// @Success 200 {object} response.PaginatedResponse{data=[]domain.UserPublic}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/users [get]
func (h *UserHandler) List(c *fiber.Ctx) error {
	params := repository.UserListParams{
		PaginationParams: repository.PaginationParams{
			Page:    c.QueryInt("page", 1),
			PerPage: c.QueryInt("per_page", 20),
			OrderBy: c.Query("order_by", "created_at"),
			Order:   c.Query("order", "desc"),
		},
		Search: c.Query("search"),
	}

	// Parse optional filters
	if branchID := c.Query("branch_id"); branchID != "" {
		id, _ := strconv.ParseInt(branchID, 10, 64)
		params.BranchID = &id
	}
	if roleID := c.Query("role_id"); roleID != "" {
		id, _ := strconv.ParseInt(roleID, 10, 64)
		params.RoleID = &id
	}
	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		params.IsActive = &active
	}

	result, err := h.userService.List(c.Context(), params)
	if err != nil {
		return response.InternalError(c, "")
	}

	return response.Paginated(c, result.Data, result.Page, result.PerPage, result.Total)
}

// Update handles user update
// @Summary Update User
// @Description Update an existing user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body service.UpdateUserInput true "User data"
// @Success 200 {object} response.Response{data=domain.UserPublic}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	var input service.UpdateUserInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	// Validate input
	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	// Update user
	user, err := h.userService.Update(c.Context(), id, input)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, user)
}

// Delete handles user deletion
// @Summary Delete User
// @Description Soft delete a user
// @Tags Users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 204
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	// Prevent self-deletion
	currentUser := middleware.GetUser(c)
	if currentUser != nil && currentUser.ID == id {
		return response.BadRequest(c, "Cannot delete your own account")
	}

	if err := h.userService.Delete(c.Context(), id); err != nil {
		return response.NotFound(c, "User not found")
	}

	return response.NoContent(c)
}

// ResetPasswordInput represents reset password request data
type ResetPasswordInput struct {
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ResetPassword handles admin password reset
// @Summary Reset User Password
// @Description Reset a user's password (admin only)
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body ResetPasswordInput true "New password"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/users/{id}/reset-password [post]
func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	var input ResetPasswordInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	// Validate input
	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	if err := h.userService.ResetPassword(c.Context(), id, input.NewPassword); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{"message": "Password reset successfully"})
}

// RegisterRoutes registers user routes
func (h *UserHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	users := app.Group("/users")
	users.Use(authMiddleware.Authenticate())

	// All routes require authentication and users.* permission
	users.Get("/", authMiddleware.RequirePermission("users.read"), h.List)
	users.Post("/", authMiddleware.RequirePermission("users.create"), h.Create)
	users.Get("/:id", authMiddleware.RequirePermission("users.read"), h.GetByID)
	users.Put("/:id", authMiddleware.RequirePermission("users.update"), h.Update)
	users.Delete("/:id", authMiddleware.RequirePermission("users.delete"), h.Delete)
	users.Post("/:id/reset-password", authMiddleware.RequirePermission("users.update"), h.ResetPassword)
}
