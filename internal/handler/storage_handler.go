package handler

import (
	"io"
	"strconv"

	"pawnshop/internal/middleware"
	"pawnshop/internal/service"
	"pawnshop/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type StorageHandler struct {
	storageService service.StorageService
	itemService    *service.ItemService
}

func NewStorageHandler(storageService service.StorageService, itemService *service.ItemService) *StorageHandler {
	return &StorageHandler{
		storageService: storageService,
		itemService:    itemService,
	}
}

// UploadItemImage uploads an image for an item
// @Summary Upload item image
// @Tags Storage
// @Accept multipart/form-data
// @Produce json
// @Param item_id path int true "Item ID"
// @Param image formance file true "Image file"
// @Success 201 {object} service.ImageInfo
// @Router /api/v1/items/{item_id}/images [post]
func (h *StorageHandler) UploadItemImage(c *fiber.Ctx) error {
	itemID, err := strconv.ParseInt(c.Params("item_id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid item ID format")
	}

	// Verify item exists
	_, err = h.itemService.GetByID(c.Context(), itemID)
	if err != nil {
		return response.NotFound(c, "Item not found")
	}

	// Get file from form
	file, err := c.FormFile("image")
	if err != nil {
		return response.BadRequest(c, "No image file provided")
	}

	// Upload image
	category := "items"
	imageInfo, err := h.storageService.UploadImage(c.Context(), file, category)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Add photo URL to item
	userID := middleware.GetUser(c).ID
	if err := h.itemService.AddPhoto(c.Context(), itemID, imageInfo.URL, userID); err != nil {
		// If we can't save to item, try to delete the uploaded image
		_ = h.storageService.DeleteImage(c.Context(), imageInfo.ID)
		return response.InternalError(c, "Failed to save photo to item")
	}

	return response.Created(c, imageInfo)
}

// GetItemImages lists images for an item
// @Summary List item images
// @Tags Storage
// @Produce json
// @Param item_id path int true "Item ID"
// @Success 200 {array} service.ImageInfo
// @Router /api/v1/items/{item_id}/images [get]
func (h *StorageHandler) GetItemImages(c *fiber.Ctx) error {
	itemID, err := strconv.ParseInt(c.Params("item_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid item ID",
		})
	}

	// Verify item exists
	_, err = h.itemService.GetByID(c.Context(), itemID)
	if err != nil {
		return handleServiceError(c, err)
	}

	// List images for the item category
	// In a real implementation, you'd filter by item ID
	images, err := h.storageService.ListImages(c.Context(), "items")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(images)
}

// DeleteItemImage deletes an image for an item
// @Summary Delete item image
// @Tags Storage
// @Param item_id path int true "Item ID"
// @Param url query string true "Photo URL to delete"
// @Success 204
// @Router /api/v1/items/{item_id}/images [delete]
func (h *StorageHandler) DeleteItemImage(c *fiber.Ctx) error {
	itemID, err := strconv.ParseInt(c.Params("item_id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid item ID format")
	}

	photoURL := c.Query("url")
	if photoURL == "" {
		return response.BadRequest(c, "Photo URL is required")
	}

	// Verify item exists
	_, err = h.itemService.GetByID(c.Context(), itemID)
	if err != nil {
		return response.NotFound(c, "Item not found")
	}

	// Remove photo from item
	userID := middleware.GetUser(c).ID
	if err := h.itemService.RemovePhoto(c.Context(), itemID, photoURL, userID); err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Try to delete the actual file (ignore errors as the URL removal is the important part)
	_ = h.storageService.DeleteImage(c.Context(), photoURL)

	return response.NoContent(c)
}

// ServeImage serves an image file
// @Summary Serve image
// @Tags Storage
// @Produce image/*
// @Param path path string true "Image path"
// @Success 200 {file} file
// @Router /storage/images/{path} [get]
func (h *StorageHandler) ServeImage(c *fiber.Ctx) error {
	path := c.Params("*")
	if path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid path",
		})
	}

	reader, info, err := h.storageService.GetImage(c.Context(), path)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Image not found",
		})
	}

	// Read all content before closing the reader
	content, err := io.ReadAll(reader)
	reader.Close()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read image",
		})
	}

	c.Set("Content-Type", info.MimeType)
	c.Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year
	return c.Send(content)
}

// ServeThumbnail serves a thumbnail file
// @Summary Serve thumbnail
// @Tags Storage
// @Produce image/*
// @Param path path string true "Thumbnail path"
// @Success 200 {file} file
// @Router /storage/thumbnails/{path} [get]
func (h *StorageHandler) ServeThumbnail(c *fiber.Ctx) error {
	path := c.Params("*")
	if path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid path",
		})
	}

	reader, info, err := h.storageService.GetThumbnail(c.Context(), path)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Thumbnail not found",
		})
	}

	// Read all content before closing the reader
	content, err := io.ReadAll(reader)
	reader.Close()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read thumbnail",
		})
	}

	c.Set("Content-Type", info.MimeType)
	c.Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year
	return c.Send(content)
}

// RegisterRoutes registers storage routes
func (h *StorageHandler) RegisterRoutes(app *fiber.App, apiRouter fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	// Public routes for serving images (no auth required)
	storage := app.Group("/storage")
	storage.Get("/images/*", h.ServeImage)
	storage.Get("/thumbnails/*", h.ServeThumbnail)

	// Protected routes for managing images
	items := apiRouter.Group("/items/:item_id/images")
	items.Use(authMiddleware.Authenticate())
	items.Post("/", authMiddleware.RequirePermission("items.update"), h.UploadItemImage)
	items.Get("/", authMiddleware.RequirePermission("items.read"), h.GetItemImages)
	items.Delete("/", authMiddleware.RequirePermission("items.update"), h.DeleteItemImage)
}
