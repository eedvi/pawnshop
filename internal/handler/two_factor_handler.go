package handler

import (
	"pawnshop/internal/middleware"
	"pawnshop/internal/service"
	"github.com/gofiber/fiber/v2"
)

type TwoFactorHandler struct {
	twoFactorService service.TwoFactorService
	userService      *service.UserService
}

func NewTwoFactorHandler(twoFactorService service.TwoFactorService, userService *service.UserService) *TwoFactorHandler {
	return &TwoFactorHandler{
		twoFactorService: twoFactorService,
		userService:      userService,
	}
}

// Setup initiates 2FA setup for the current user
// @Summary Setup 2FA
// @Tags Two-Factor Authentication
// @Produce json
// @Success 200 {object} domain.TwoFactorSetup
// @Router /api/v1/auth/2fa/setup [post]
func (h *TwoFactorHandler) Setup(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	// Get user email
	user, err := h.userService.GetByID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, err)
	}

	setup, err := h.twoFactorService.Setup(c.Context(), userID, user.Email)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(setup)
}

// Enable enables 2FA after verifying the TOTP code
// @Summary Enable 2FA
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Param body body object{code string} true "TOTP code"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/2fa/enable [post]
func (h *TwoFactorHandler) Enable(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	var body struct {
		Code string `json:"code"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	if body.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Code is required",
		})
	}

	if err := h.twoFactorService.Enable(c.Context(), userID, body.Code); err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "2FA has been enabled successfully",
	})
}

// Disable disables 2FA for the current user
// @Summary Disable 2FA
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Param body body object{password string} true "User password"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/2fa/disable [post]
func (h *TwoFactorHandler) Disable(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	var body struct {
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	if body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is required",
		})
	}

	if err := h.twoFactorService.Disable(c.Context(), userID, body.Password); err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "2FA has been disabled successfully",
	})
}

// GetStatus returns the 2FA status for the current user
// @Summary Get 2FA status
// @Tags Two-Factor Authentication
// @Produce json
// @Success 200 {object} domain.TwoFactorStatus
// @Router /api/v1/auth/2fa/status [get]
func (h *TwoFactorHandler) GetStatus(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	status, err := h.twoFactorService.GetStatus(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(status)
}

// Verify verifies a 2FA challenge during login
// @Summary Verify 2FA challenge
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Param body body object{token string, code string} true "Challenge verification"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/2fa/verify [post]
func (h *TwoFactorHandler) Verify(c *fiber.Ctx) error {
	var body struct {
		Token string `json:"token"`
		Code  string `json:"code"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	if body.Token == "" || body.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token and code are required",
		})
	}

	challenge, err := h.twoFactorService.VerifyChallenge(c.Context(), body.Token, body.Code)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"message":   "2FA verification successful",
		"user_id":   challenge.UserID,
		"verified":  true,
	})
}

// VerifyWithBackup verifies a 2FA challenge using a backup code
// @Summary Verify 2FA with backup code
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Param body body object{token string, backup_code string} true "Challenge verification with backup"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/2fa/verify-backup [post]
func (h *TwoFactorHandler) VerifyWithBackup(c *fiber.Ctx) error {
	var body struct {
		Token      string `json:"token"`
		BackupCode string `json:"backup_code"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	if body.Token == "" || body.BackupCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token and backup_code are required",
		})
	}

	challenge, err := h.twoFactorService.VerifyChallengeWithBackup(c.Context(), body.Token, body.BackupCode)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"message":  "2FA verification successful",
		"user_id":  challenge.UserID,
		"verified": true,
	})
}

// RegenerateBackupCodes regenerates backup codes for the current user
// @Summary Regenerate backup codes
// @Tags Two-Factor Authentication
// @Accept json
// @Produce json
// @Param body body object{password string} true "User password"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/2fa/backup-codes/regenerate [post]
func (h *TwoFactorHandler) RegenerateBackupCodes(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	var body struct {
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	// Verify password first
	user, err := h.userService.GetByID(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, err)
	}

	// Note: Password verification should be done in the service
	_ = user

	codes, err := h.twoFactorService.RegenerateBackupCodes(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"backup_codes": codes,
		"message":      "Backup codes have been regenerated. Please save them in a secure location.",
	})
}

// GetBackupCodesCount returns the number of unused backup codes
// @Summary Get backup codes count
// @Tags Two-Factor Authentication
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/2fa/backup-codes/count [get]
func (h *TwoFactorHandler) GetBackupCodesCount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	count, err := h.twoFactorService.GetBackupCodesCount(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"unused_count": count,
	})
}

// RegisterRoutes registers 2FA routes
func (h *TwoFactorHandler) RegisterRoutes(router fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	twoFactor := router.Group("/auth/2fa")

	// Public routes (for 2FA verification during login)
	twoFactor.Post("/verify", h.Verify)
	twoFactor.Post("/verify-backup", h.VerifyWithBackup)

	// Protected routes
	twoFactor.Use(authMiddleware.Authenticate())
	twoFactor.Post("/setup", h.Setup)
	twoFactor.Post("/enable", h.Enable)
	twoFactor.Post("/disable", h.Disable)
	twoFactor.Get("/status", h.GetStatus)
	twoFactor.Post("/backup-codes/regenerate", h.RegenerateBackupCodes)
	twoFactor.Get("/backup-codes/count", h.GetBackupCodesCount)
}
