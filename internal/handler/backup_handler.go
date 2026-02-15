package handler

import (
	"fmt"

	"pawnshop/internal/middleware"
	"pawnshop/internal/service"
	"github.com/gofiber/fiber/v2"
)

type BackupHandler struct {
	backupService service.BackupService
}

func NewBackupHandler(backupService service.BackupService) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
	}
}

// Create creates a new database backup
// @Summary Create database backup
// @Tags Backup
// @Accept json
// @Produce json
// @Param body body object{description string} false "Backup description"
// @Success 201 {object} service.BackupInfo
// @Router /api/v1/admin/backups [post]
func (h *BackupHandler) Create(c *fiber.Ctx) error {
	var body struct {
		Description string `json:"description"`
	}
	c.BodyParser(&body)

	backup, err := h.backupService.CreateBackup(c.Context(), body.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(backup)
}

// List lists all available backups
// @Summary List backups
// @Tags Backup
// @Produce json
// @Success 200 {array} service.BackupInfo
// @Router /api/v1/admin/backups [get]
func (h *BackupHandler) List(c *fiber.Ctx) error {
	backups, err := h.backupService.ListBackups(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(backups)
}

// Download downloads a backup file
// @Summary Download backup
// @Tags Backup
// @Produce application/octet-stream
// @Param filename path string true "Backup filename"
// @Success 200 {file} file
// @Router /api/v1/admin/backups/{filename}/download [get]
func (h *BackupHandler) Download(c *fiber.Ctx) error {
	filename := c.Params("filename")

	reader, info, err := h.backupService.GetBackup(c.Context(), filename)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer reader.Close()

	// Set headers for file download
	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", info.Filename))
	c.Set("Content-Length", fmt.Sprintf("%d", info.Size))

	return c.SendStream(reader)
}

// Restore restores a database from a backup
// @Summary Restore database from backup
// @Tags Backup
// @Param filename path string true "Backup filename"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/backups/{filename}/restore [post]
func (h *BackupHandler) Restore(c *fiber.Ctx) error {
	filename := c.Params("filename")

	if err := h.backupService.RestoreBackup(c.Context(), filename); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Database restored successfully",
		"filename": filename,
	})
}

// Delete deletes a backup file
// @Summary Delete backup
// @Tags Backup
// @Param filename path string true "Backup filename"
// @Success 204
// @Router /api/v1/admin/backups/{filename} [delete]
func (h *BackupHandler) Delete(c *fiber.Ctx) error {
	filename := c.Params("filename")

	if err := h.backupService.DeleteBackup(c.Context(), filename); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Cleanup removes old backups
// @Summary Cleanup old backups
// @Tags Backup
// @Accept json
// @Produce json
// @Param body body object{retention_days int} true "Retention period in days"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/backups/cleanup [post]
func (h *BackupHandler) Cleanup(c *fiber.Ctx) error {
	var body struct {
		RetentionDays int `json:"retention_days"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	if body.RetentionDays <= 0 {
		body.RetentionDays = 30 // Default to 30 days
	}

	deleted, err := h.backupService.CleanupOldBackups(c.Context(), body.RetentionDays)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":        "Cleanup completed",
		"deleted_count":  deleted,
		"retention_days": body.RetentionDays,
	})
}

// RegisterRoutes registers backup routes
func (h *BackupHandler) RegisterRoutes(router fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	backups := router.Group("/admin/backups")
	backups.Use(authMiddleware.Authenticate())
	backups.Use(authMiddleware.RequirePermission("admin:backup"))

	backups.Post("/", h.Create)
	backups.Get("/", h.List)
	backups.Get("/:filename/download", h.Download)
	backups.Post("/:filename/restore", h.Restore)
	backups.Delete("/:filename", h.Delete)
	backups.Post("/cleanup", h.Cleanup)
}
