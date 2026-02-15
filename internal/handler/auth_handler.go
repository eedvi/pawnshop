package handler

import (
	"github.com/gofiber/fiber/v2"
	"pawnshop/internal/middleware"
	"pawnshop/internal/service"
	"pawnshop/pkg/response"
	"pawnshop/pkg/validator"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *service.AuthService
	auditLogger *middleware.AuditLogger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *service.AuthService, auditLogger *middleware.AuditLogger) *AuthHandler {
	return &AuthHandler{authService: authService, auditLogger: auditLogger}
}

// Login handles user login
// @Summary Login
// @Description Authenticate user and return tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body service.LoginInput true "Login credentials"
// @Success 200 {object} response.Response{data=service.LoginOutput}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input service.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	// Validate input
	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	// Get client IP
	ip := c.IP()

	// Login
	output, err := h.authService.Login(c.Context(), input, ip)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	// Log login event
	if h.auditLogger != nil {
		h.auditLogger.LogLogin(c, output.User.ID, output.User.BranchID)
	}

	return response.OK(c, output)
}

// Refresh handles token refresh
// @Summary Refresh Token
// @Description Generate new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body service.RefreshInput true "Refresh token"
// @Success 200 {object} response.Response{data=service.LoginOutput}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var input service.RefreshInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	// Validate input
	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	// Refresh
	output, err := h.authService.Refresh(c.Context(), input)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.OK(c, output)
}

// Logout handles user logout
// @Summary Logout
// @Description Invalidate all user tokens
// @Tags Auth
// @Security BearerAuth
// @Success 204
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return response.Unauthorized(c, "")
	}

	if err := h.authService.Logout(c.Context(), user.ID); err != nil {
		return response.InternalError(c, "Failed to logout")
	}

	// Log logout event
	if h.auditLogger != nil {
		h.auditLogger.LogLogout(c, user.ID, user.BranchID)
	}

	return response.NoContent(c)
}

// Me returns the current authenticated user
// @Summary Get Current User
// @Description Get the currently authenticated user's profile
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=domain.UserPublic}
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return response.Unauthorized(c, "")
	}

	return response.OK(c, user.ToPublic())
}

// ChangePassword handles password change
// @Summary Change Password
// @Description Change the current user's password
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body service.ChangePasswordInput true "Password change data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return response.Unauthorized(c, "")
	}

	var input service.ChangePasswordInput
	if err := c.BodyParser(&input); err != nil {
		return response.BadRequest(c, "Error parsing request body: "+err.Error())
	}

	// Validate input
	if errors := validator.Validate(&input); errors != nil {
		return response.ValidationError(c, errors)
	}

	// Change password
	if err := h.authService.ChangePassword(c.Context(), user.ID, input); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.OK(c, fiber.Map{"message": "Password changed successfully"})
}

// RegisterRoutes registers auth routes
func (h *AuthHandler) RegisterRoutes(app fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	auth := app.Group("/auth")

	// Public routes
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.Refresh)

	// Protected routes
	auth.Use(authMiddleware.Authenticate())
	auth.Post("/logout", h.Logout)
	auth.Get("/me", h.Me)
	auth.Post("/change-password", h.ChangePassword)
}
