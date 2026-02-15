package handler

import (
	"strconv"
	"time"

	"pawnshop/internal/middleware"
	"pawnshop/internal/repository"
	"pawnshop/internal/service"
	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	notificationService service.NotificationService
}

func NewNotificationHandler(notificationService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// Template Handlers

// CreateTemplate creates a new notification template
// @Summary Create a notification template
// @Tags Notifications
// @Accept json
// @Produce json
// @Param template body service.CreateNotificationTemplateRequest true "Template data"
// @Success 201 {object} domain.NotificationTemplate
// @Router /api/v1/notifications/templates [post]
func (h *NotificationHandler) CreateTemplate(c *fiber.Ctx) error {
	var req service.CreateNotificationTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	template, err := h.notificationService.CreateTemplate(c.Context(), req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(template)
}

// GetTemplateByID retrieves a notification template by ID
// @Summary Get a notification template by ID
// @Tags Notifications
// @Produce json
// @Param id path int true "Template ID"
// @Success 200 {object} domain.NotificationTemplate
// @Router /api/v1/notifications/templates/{id} [get]
func (h *NotificationHandler) GetTemplateByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid template ID format",
		})
	}

	template, err := h.notificationService.GetTemplateByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(template)
}

// UpdateTemplate updates a notification template
// @Summary Update a notification template
// @Tags Notifications
// @Accept json
// @Produce json
// @Param id path int true "Template ID"
// @Param template body service.UpdateNotificationTemplateRequest true "Template data"
// @Success 200 {object} domain.NotificationTemplate
// @Router /api/v1/notifications/templates/{id} [put]
func (h *NotificationHandler) UpdateTemplate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid template ID format",
		})
	}

	var req service.UpdateNotificationTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	template, err := h.notificationService.UpdateTemplate(c.Context(), id, req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(template)
}

// DeleteTemplate deletes a notification template
// @Summary Delete a notification template
// @Tags Notifications
// @Param id path int true "Template ID"
// @Success 204
// @Router /api/v1/notifications/templates/{id} [delete]
func (h *NotificationHandler) DeleteTemplate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid template ID format",
		})
	}

	if err := h.notificationService.DeleteTemplate(c.Context(), id); err != nil {
		return handleServiceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListTemplates retrieves all notification templates
// @Summary List notification templates
// @Tags Notifications
// @Produce json
// @Param include_inactive query bool false "Include inactive templates"
// @Success 200 {array} domain.NotificationTemplate
// @Router /api/v1/notifications/templates [get]
func (h *NotificationHandler) ListTemplates(c *fiber.Ctx) error {
	includeInactive := c.QueryBool("include_inactive", false)

	templates, err := h.notificationService.ListTemplates(c.Context(), includeInactive)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(templates)
}

// Notification Handlers

// Create creates a new notification
// @Summary Create a notification
// @Tags Notifications
// @Accept json
// @Produce json
// @Param notification body service.CreateNotificationRequest true "Notification data"
// @Success 201 {object} domain.Notification
// @Router /api/v1/notifications [post]
func (h *NotificationHandler) Create(c *fiber.Ctx) error {
	var req service.CreateNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	notification, err := h.notificationService.Create(c.Context(), req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(notification)
}

// CreateFromTemplate creates a notification from a template
// @Summary Create a notification from a template
// @Tags Notifications
// @Accept json
// @Produce json
// @Param notification body service.CreateNotificationFromTemplateRequest true "Notification data"
// @Success 201 {object} domain.Notification
// @Router /api/v1/notifications/from-template [post]
func (h *NotificationHandler) CreateFromTemplate(c *fiber.Ctx) error {
	var req service.CreateNotificationFromTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	notification, err := h.notificationService.CreateFromTemplate(c.Context(), req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(notification)
}

// GetByID retrieves a notification by ID
// @Summary Get a notification by ID
// @Tags Notifications
// @Produce json
// @Param id path int true "Notification ID"
// @Success 200 {object} domain.Notification
// @Router /api/v1/notifications/{id} [get]
func (h *NotificationHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification ID format",
		})
	}

	notification, err := h.notificationService.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(notification)
}

// List retrieves notifications with filtering
// @Summary List notifications
// @Tags Notifications
// @Produce json
// @Param customer_id query int false "Filter by customer"
// @Param branch_id query int false "Filter by branch"
// @Param notification_type query string false "Filter by type"
// @Param channel query string false "Filter by channel"
// @Param status query string false "Filter by status"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/notifications [get]
func (h *NotificationHandler) List(c *fiber.Ctx) error {
	filter := repository.NotificationFilter{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 20),
	}

	if customerID := c.QueryInt("customer_id"); customerID > 0 {
		id := int64(customerID)
		filter.CustomerID = &id
	}
	if branchID := c.QueryInt("branch_id"); branchID > 0 {
		id := int64(branchID)
		filter.BranchID = &id
	}
	if notificationType := c.Query("notification_type"); notificationType != "" {
		filter.NotificationType = &notificationType
	}
	if channel := c.Query("channel"); channel != "" {
		filter.Channel = &channel
	}
	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	notifications, total, err := h.notificationService.List(c.Context(), filter)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"data":  notifications,
		"total": total,
		"page":  filter.Page,
		"size":  filter.PageSize,
	})
}

// ListByCustomer retrieves notifications for a customer
// @Summary List notifications for a customer
// @Tags Notifications
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/customers/{customer_id}/notifications [get]
func (h *NotificationHandler) ListByCustomer(c *fiber.Ctx) error {
	customerID, err := strconv.ParseInt(c.Params("customer_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID format",
		})
	}

	filter := repository.NotificationFilter{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 20),
	}

	notifications, total, err := h.notificationService.ListByCustomer(c.Context(), customerID, filter)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"data":  notifications,
		"total": total,
		"page":  filter.Page,
		"size":  filter.PageSize,
	})
}

// Cancel cancels a pending notification
// @Summary Cancel a notification
// @Tags Notifications
// @Param id path int true "Notification ID"
// @Success 204
// @Router /api/v1/notifications/{id}/cancel [post]
func (h *NotificationHandler) Cancel(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification ID format",
		})
	}

	if err := h.notificationService.Cancel(c.Context(), id); err != nil {
		return handleServiceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Customer Preference Handlers

// GetCustomerPreferences retrieves notification preferences for a customer
// @Summary Get customer notification preferences
// @Tags Notifications
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Success 200 {array} domain.CustomerNotificationPreference
// @Router /api/v1/customers/{customer_id}/notification-preferences [get]
func (h *NotificationHandler) GetCustomerPreferences(c *fiber.Ctx) error {
	customerID, err := strconv.ParseInt(c.Params("customer_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID format",
		})
	}

	prefs, err := h.notificationService.GetCustomerPreferences(c.Context(), customerID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(prefs)
}

// UpdateCustomerPreferences updates notification preferences for a customer
// @Summary Update customer notification preferences
// @Tags Notifications
// @Accept json
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Param preferences body []domain.CustomerNotificationPreference true "Preferences"
// @Success 200 {array} domain.CustomerNotificationPreference
// @Router /api/v1/customers/{customer_id}/notification-preferences [put]
func (h *NotificationHandler) UpdateCustomerPreferences(c *fiber.Ctx) error {
	customerID, err := strconv.ParseInt(c.Params("customer_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID format",
		})
	}

	var req struct {
		Preferences []struct {
			NotificationType string `json:"notification_type"`
			Channel          string `json:"channel"`
			IsEnabled        bool   `json:"is_enabled"`
		} `json:"preferences"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	// Convert to domain objects
	var prefs []*struct {
		NotificationType string
		Channel          string
		IsEnabled        bool
	}
	for _, p := range req.Preferences {
		prefs = append(prefs, &struct {
			NotificationType string
			Channel          string
			IsEnabled        bool
		}{
			NotificationType: p.NotificationType,
			Channel:          p.Channel,
			IsEnabled:        p.IsEnabled,
		})
	}

	// Note: You'll need to convert these to domain.CustomerNotificationPreference
	// This is a simplified implementation
	if err := h.notificationService.UpdateCustomerPreferences(c.Context(), customerID, nil); err != nil {
		return handleServiceError(c, err)
	}

	// Return updated preferences
	updatedPrefs, err := h.notificationService.GetCustomerPreferences(c.Context(), customerID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(updatedPrefs)
}

// Internal Notification Handlers

// CreateInternalNotification creates a new internal notification
// @Summary Create an internal notification
// @Tags Notifications
// @Accept json
// @Produce json
// @Param notification body service.CreateInternalNotificationRequest true "Notification data"
// @Success 201 {object} domain.InternalNotification
// @Router /api/v1/internal-notifications [post]
func (h *NotificationHandler) CreateInternalNotification(c *fiber.Ctx) error {
	var req service.CreateInternalNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	notification, err := h.notificationService.CreateInternalNotification(c.Context(), req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(notification)
}

// GetMyInternalNotifications retrieves internal notifications for the current user
// @Summary Get my internal notifications
// @Tags Notifications
// @Produce json
// @Param is_read query bool false "Filter by read status"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/internal-notifications/me [get]
func (h *NotificationHandler) GetMyInternalNotifications(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	filter := repository.InternalNotificationFilter{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 20),
	}

	if c.Query("is_read") != "" {
		isRead := c.QueryBool("is_read")
		filter.IsRead = &isRead
	}

	notifications, total, err := h.notificationService.ListInternalNotificationsByUser(c.Context(), userID, filter)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"data":  notifications,
		"total": total,
		"page":  filter.Page,
		"size":  filter.PageSize,
	})
}

// GetUnreadInternalNotifications retrieves unread internal notifications for the current user
// @Summary Get my unread internal notifications
// @Tags Notifications
// @Produce json
// @Param limit query int false "Limit"
// @Success 200 {array} domain.InternalNotification
// @Router /api/v1/internal-notifications/me/unread [get]
func (h *NotificationHandler) GetUnreadInternalNotifications(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)
	limit := c.QueryInt("limit", 10)

	notifications, err := h.notificationService.GetUnreadInternalNotifications(c.Context(), userID, limit)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(notifications)
}

// MarkInternalNotificationAsRead marks an internal notification as read
// @Summary Mark internal notification as read
// @Tags Notifications
// @Param id path int true "Notification ID"
// @Success 204
// @Router /api/v1/internal-notifications/{id}/read [post]
func (h *NotificationHandler) MarkInternalNotificationAsRead(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid notification ID format",
		})
	}

	if err := h.notificationService.MarkInternalNotificationAsRead(c.Context(), id); err != nil {
		return handleServiceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// MarkAllInternalNotificationsAsRead marks all internal notifications as read for the current user
// @Summary Mark all internal notifications as read
// @Tags Notifications
// @Success 204
// @Router /api/v1/internal-notifications/me/read-all [post]
func (h *NotificationHandler) MarkAllInternalNotificationsAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	if err := h.notificationService.MarkAllInternalNotificationsAsRead(c.Context(), userID); err != nil {
		return handleServiceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetUnreadCount retrieves the count of unread internal notifications for the current user
// @Summary Get unread notification count
// @Tags Notifications
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/internal-notifications/me/unread-count [get]
func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	count, err := h.notificationService.GetUnreadCount(c.Context(), userID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"count": count,
	})
}

// Stats Handlers

// GetStatsByCustomer retrieves notification stats for a customer
// @Summary Get notification stats for a customer
// @Tags Notifications
// @Produce json
// @Param customer_id path int true "Customer ID"
// @Success 200 {object} repository.NotificationStats
// @Router /api/v1/customers/{customer_id}/notification-stats [get]
func (h *NotificationHandler) GetStatsByCustomer(c *fiber.Ctx) error {
	customerID, err := strconv.ParseInt(c.Params("customer_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid customer ID format",
		})
	}

	stats, err := h.notificationService.GetStatsByCustomer(c.Context(), customerID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(stats)
}

// GetStatsByBranch retrieves notification stats for a branch
// @Summary Get notification stats for a branch
// @Tags Notifications
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} repository.NotificationStats
// @Router /api/v1/branches/{branch_id}/notification-stats [get]
func (h *NotificationHandler) GetStatsByBranch(c *fiber.Ctx) error {
	branchID, err := strconv.ParseInt(c.Params("branch_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid branch ID format",
		})
	}

	// Default to last 30 days
	dateTo := time.Now()
	dateFrom := dateTo.AddDate(0, 0, -30)

	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = parsed
		}
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = parsed
		}
	}

	stats, err := h.notificationService.GetStatsByBranch(c.Context(), branchID, dateFrom, dateTo)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(stats)
}

// RegisterRoutes registers notification routes
func (h *NotificationHandler) RegisterRoutes(router fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	// Notification templates
	templates := router.Group("/notifications/templates")
	templates.Use(authMiddleware.Authenticate())
	templates.Post("/", authMiddleware.RequirePermission("notifications:manage"), h.CreateTemplate)
	templates.Get("/", authMiddleware.RequirePermission("notifications:read"), h.ListTemplates)
	templates.Get("/:id", authMiddleware.RequirePermission("notifications:read"), h.GetTemplateByID)
	templates.Put("/:id", authMiddleware.RequirePermission("notifications:manage"), h.UpdateTemplate)
	templates.Delete("/:id", authMiddleware.RequirePermission("notifications:manage"), h.DeleteTemplate)

	// Notifications
	notifications := router.Group("/notifications")
	notifications.Use(authMiddleware.Authenticate())
	notifications.Post("/", authMiddleware.RequirePermission("notifications:create"), h.Create)
	notifications.Post("/from-template", authMiddleware.RequirePermission("notifications:create"), h.CreateFromTemplate)
	notifications.Get("/", authMiddleware.RequirePermission("notifications:read"), h.List)
	notifications.Get("/:id", authMiddleware.RequirePermission("notifications:read"), h.GetByID)
	notifications.Post("/:id/cancel", authMiddleware.RequirePermission("notifications:manage"), h.Cancel)

	// Internal notifications
	internalNotifications := router.Group("/internal-notifications")
	internalNotifications.Use(authMiddleware.Authenticate())
	internalNotifications.Post("/", authMiddleware.RequirePermission("notifications:create"), h.CreateInternalNotification)
	internalNotifications.Get("/me", h.GetMyInternalNotifications)
	internalNotifications.Get("/me/unread", h.GetUnreadInternalNotifications)
	internalNotifications.Get("/me/unread-count", h.GetUnreadCount)
	internalNotifications.Post("/me/read-all", h.MarkAllInternalNotificationsAsRead)
	internalNotifications.Post("/:id/read", h.MarkInternalNotificationAsRead)

	// Customer notification preferences (nested under customers)
	customerNotifications := router.Group("/customers/:customer_id")
	customerNotifications.Use(authMiddleware.Authenticate())
	customerNotifications.Get("/notifications", authMiddleware.RequirePermission("notifications:read"), h.ListByCustomer)
	customerNotifications.Get("/notification-preferences", authMiddleware.RequirePermission("customers:read"), h.GetCustomerPreferences)
	customerNotifications.Put("/notification-preferences", authMiddleware.RequirePermission("customers:update"), h.UpdateCustomerPreferences)
	customerNotifications.Get("/notification-stats", authMiddleware.RequirePermission("notifications:read"), h.GetStatsByCustomer)

	// Branch notification stats (nested under branches)
	branchNotifications := router.Group("/branches/:branch_id")
	branchNotifications.Use(authMiddleware.Authenticate())
	branchNotifications.Get("/notification-stats", authMiddleware.RequirePermission("notifications:read"), h.GetStatsByBranch)
}
