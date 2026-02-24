package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/pkg/auth"
	"pawnshop/pkg/logger"
	"pawnshop/pkg/response"
)

// AuthMiddleware handles authentication
type AuthMiddleware struct {
	jwtManager *auth.JWTManager
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	logger     zerolog.Logger
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(jwtManager *auth.JWTManager, userRepo repository.UserRepository, roleRepo repository.RoleRepository, log zerolog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		logger:     log.With().Str("middleware", "auth").Logger(),
	}
}

// Authenticate validates the JWT token and loads the user
func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.FromContext(c.UserContext(), m.logger)

		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Warn().
				Str("ip", c.IP()).
				Str("path", c.Path()).
				Str("method", c.Method()).
				Msg("Authentication failed: missing authorization header")
			return response.Unauthorized(c, "Missing authorization header")
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Warn().
				Str("ip", c.IP()).
				Str("path", c.Path()).
				Str("method", c.Method()).
				Str("auth_header_prefix", parts[0]).
				Msg("Authentication failed: invalid authorization header format")
			return response.Unauthorized(c, "Invalid authorization header format")
		}

		tokenString := parts[1]

		// Validate token (NEVER log the actual token)
		claims, err := m.jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			log.Warn().
				Err(err).
				Str("ip", c.IP()).
				Str("path", c.Path()).
				Str("method", c.Method()).
				Msg("Authentication failed: invalid or expired token")
			return response.Unauthorized(c, "Invalid or expired token")
		}

		// Get user from database
		user, err := m.userRepo.GetByID(c.Context(), claims.UserID)
		if err != nil {
			log.Warn().
				Err(err).
				Int64("user_id", claims.UserID).
				Str("ip", c.IP()).
				Str("path", c.Path()).
				Str("method", c.Method()).
				Msg("Authentication failed: user not found")
			return response.Unauthorized(c, "User not found")
		}

		// Check if user is active
		if !user.CanLogin() {
			log.Warn().
				Int64("user_id", user.ID).
				Str("email", user.Email).
				Bool("is_active", user.IsActive).
				Bool("is_locked", user.IsLocked()).
				Str("ip", c.IP()).
				Str("path", c.Path()).
				Str("method", c.Method()).
				Msg("Authentication failed: account is inactive or locked")
			return response.Unauthorized(c, "Account is inactive or locked")
		}

		// Load role
		role, err := m.roleRepo.GetByID(c.Context(), user.RoleID)
		if err != nil {
			log.Error().
				Err(err).
				Int64("user_id", user.ID).
				Int64("role_id", user.RoleID).
				Str("ip", c.IP()).
				Str("path", c.Path()).
				Msg("Failed to load user role")
			return response.InternalError(c, "Failed to load user role")
		}
		user.Role = role

		// Store user and claims in context
		c.Locals("user", user)
		c.Locals("claims", claims)

		// Inject user ID into context for services
		ctx := logger.WithUserID(c.UserContext(), user.ID)
		c.SetUserContext(ctx)

		log.Debug().
			Int64("user_id", user.ID).
			Str("email", user.Email).
			Str("role", role.Name).
			Str("path", c.Path()).
			Msg("Authentication successful")

		return c.Next()
	}
}

// RequirePermission checks if the user has the required permission
func (m *AuthMiddleware) RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.FromContext(c.UserContext(), m.logger)

		user, ok := c.Locals("user").(*domain.User)
		if !ok || user == nil {
			log.Warn().
				Str("required_permission", permission).
				Str("path", c.Path()).
				Msg("Permission check failed: user not found in context")
			return response.Unauthorized(c, "")
		}

		if !user.HasPermission(permission) {
			log.Warn().
				Int64("user_id", user.ID).
				Str("email", user.Email).
				Str("required_permission", permission).
				Str("path", c.Path()).
				Msg("Permission denied: user lacks required permission")
			return response.Forbidden(c, "")
		}

		return c.Next()
	}
}

// RequireAnyPermission checks if the user has any of the required permissions
func (m *AuthMiddleware) RequireAnyPermission(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.FromContext(c.UserContext(), m.logger)

		user, ok := c.Locals("user").(*domain.User)
		if !ok || user == nil {
			log.Warn().
				Strs("required_permissions", permissions).
				Str("path", c.Path()).
				Msg("Permission check failed: user not found in context")
			return response.Unauthorized(c, "")
		}

		for _, permission := range permissions {
			if user.HasPermission(permission) {
				return c.Next()
			}
		}

		log.Warn().
			Int64("user_id", user.ID).
			Str("email", user.Email).
			Strs("required_permissions", permissions).
			Str("path", c.Path()).
			Msg("Permission denied: user lacks any of the required permissions")
		return response.Forbidden(c, "")
	}
}

// RequireAllPermissions checks if the user has all required permissions
func (m *AuthMiddleware) RequireAllPermissions(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.FromContext(c.UserContext(), m.logger)

		user, ok := c.Locals("user").(*domain.User)
		if !ok || user == nil {
			log.Warn().
				Strs("required_permissions", permissions).
				Str("path", c.Path()).
				Msg("Permission check failed: user not found in context")
			return response.Unauthorized(c, "")
		}

		for _, permission := range permissions {
			if !user.HasPermission(permission) {
				log.Warn().
					Int64("user_id", user.ID).
					Str("email", user.Email).
					Str("missing_permission", permission).
					Strs("required_permissions", permissions).
					Str("path", c.Path()).
					Msg("Permission denied: user lacks one of the required permissions")
				return response.Forbidden(c, "")
			}
		}

		return c.Next()
	}
}

// RequireRole checks if the user has the required role
func (m *AuthMiddleware) RequireRole(roleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.FromContext(c.UserContext(), m.logger)

		user, ok := c.Locals("user").(*domain.User)
		if !ok || user == nil {
			log.Warn().
				Str("required_role", roleName).
				Str("path", c.Path()).
				Msg("Role check failed: user not found in context")
			return response.Unauthorized(c, "")
		}

		if user.Role == nil || user.Role.Name != roleName {
			actualRole := ""
			if user.Role != nil {
				actualRole = user.Role.Name
			}
			log.Warn().
				Int64("user_id", user.ID).
				Str("email", user.Email).
				Str("required_role", roleName).
				Str("actual_role", actualRole).
				Str("path", c.Path()).
				Msg("Role denied: user does not have required role")
			return response.Forbidden(c, "")
		}

		return c.Next()
	}
}

// RequireAnyRole checks if the user has any of the required roles
func (m *AuthMiddleware) RequireAnyRole(roleNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logger.FromContext(c.UserContext(), m.logger)

		user, ok := c.Locals("user").(*domain.User)
		if !ok || user == nil {
			log.Warn().
				Strs("required_roles", roleNames).
				Str("path", c.Path()).
				Msg("Role check failed: user not found in context")
			return response.Unauthorized(c, "")
		}

		if user.Role == nil {
			log.Warn().
				Int64("user_id", user.ID).
				Str("email", user.Email).
				Strs("required_roles", roleNames).
				Str("path", c.Path()).
				Msg("Role denied: user has no role assigned")
			return response.Forbidden(c, "")
		}

		for _, roleName := range roleNames {
			if user.Role.Name == roleName {
				return c.Next()
			}
		}

		log.Warn().
			Int64("user_id", user.ID).
			Str("email", user.Email).
			Strs("required_roles", roleNames).
			Str("actual_role", user.Role.Name).
			Str("path", c.Path()).
			Msg("Role denied: user does not have any of the required roles")
		return response.Forbidden(c, "")
	}
}

// OptionalAuth tries to authenticate but doesn't fail if no token
func (m *AuthMiddleware) OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Next()
		}

		tokenString := parts[1]

		claims, err := m.jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			return c.Next()
		}

		user, err := m.userRepo.GetByID(c.Context(), claims.UserID)
		if err != nil {
			return c.Next()
		}

		if user.CanLogin() {
			role, err := m.roleRepo.GetByID(c.Context(), user.RoleID)
			if err == nil {
				user.Role = role
			}
			c.Locals("user", user)
			c.Locals("claims", claims)
		}

		return c.Next()
	}
}

// GetUser returns the authenticated user from context
func GetUser(c *fiber.Ctx) *domain.User {
	user, ok := c.Locals("user").(*domain.User)
	if !ok {
		return nil
	}
	return user
}

// GetClaims returns the JWT claims from context
func GetClaims(c *fiber.Ctx) *auth.JWTClaims {
	claims, ok := c.Locals("claims").(*auth.JWTClaims)
	if !ok {
		return nil
	}
	return claims
}

// GetUserID returns the authenticated user's ID from context
func GetUserID(c *fiber.Ctx) int64 {
	user := GetUser(c)
	if user == nil {
		return 0
	}
	return user.ID
}
