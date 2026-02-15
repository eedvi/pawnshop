package repository

import (
	"context"
	"time"

	"pawnshop/internal/domain"
)

// NotificationTemplateRepository defines the interface for notification template operations
type NotificationTemplateRepository interface {
	// Create creates a new notification template
	Create(ctx context.Context, template *domain.NotificationTemplate) error

	// GetByID retrieves a template by ID
	GetByID(ctx context.Context, id int64) (*domain.NotificationTemplate, error)

	// GetByTypeAndChannel retrieves a template by notification type and channel
	GetByTypeAndChannel(ctx context.Context, notificationType, channel string) (*domain.NotificationTemplate, error)

	// Update updates an existing template
	Update(ctx context.Context, template *domain.NotificationTemplate) error

	// Delete deletes a template
	Delete(ctx context.Context, id int64) error

	// List retrieves all templates
	List(ctx context.Context, includeInactive bool) ([]*domain.NotificationTemplate, error)

	// ListByType retrieves templates by notification type
	ListByType(ctx context.Context, notificationType string) ([]*domain.NotificationTemplate, error)
}

// NotificationRepository defines the interface for notification operations
type NotificationRepository interface {
	// Create creates a new notification
	Create(ctx context.Context, notification *domain.Notification) error

	// GetByID retrieves a notification by ID
	GetByID(ctx context.Context, id int64) (*domain.Notification, error)

	// Update updates an existing notification
	Update(ctx context.Context, notification *domain.Notification) error

	// Delete deletes a notification
	Delete(ctx context.Context, id int64) error

	// List retrieves notifications with filtering
	List(ctx context.Context, filter NotificationFilter) ([]*domain.Notification, int64, error)

	// ListByCustomer retrieves notifications for a customer
	ListByCustomer(ctx context.Context, customerID int64, filter NotificationFilter) ([]*domain.Notification, int64, error)

	// ListPending retrieves pending notifications ready to send
	ListPending(ctx context.Context, limit int) ([]*domain.Notification, error)

	// ListScheduled retrieves scheduled notifications ready to send
	ListScheduled(ctx context.Context, before time.Time, limit int) ([]*domain.Notification, error)

	// ListFailed retrieves failed notifications that can be retried
	ListFailed(ctx context.Context, maxRetries int, limit int) ([]*domain.Notification, error)

	// MarkAsSent marks a notification as sent
	MarkAsSent(ctx context.Context, id int64) error

	// MarkAsDelivered marks a notification as delivered
	MarkAsDelivered(ctx context.Context, id int64) error

	// MarkAsFailed marks a notification as failed
	MarkAsFailed(ctx context.Context, id int64, reason string) error

	// IncrementRetry increments the retry count
	IncrementRetry(ctx context.Context, id int64) error

	// Cancel cancels a pending notification
	Cancel(ctx context.Context, id int64) error

	// GetStatsByCustomer retrieves notification stats for a customer
	GetStatsByCustomer(ctx context.Context, customerID int64) (*NotificationStats, error)

	// GetStatsByBranch retrieves notification stats for a branch
	GetStatsByBranch(ctx context.Context, branchID int64, dateFrom, dateTo time.Time) (*NotificationStats, error)
}

// NotificationFilter contains filters for listing notifications
type NotificationFilter struct {
	CustomerID       *int64
	BranchID         *int64
	NotificationType *string
	Channel          *string
	Status           *string
	ReferenceType    *string
	ReferenceID      *int64
	DateFrom         *string
	DateTo           *string
	Page             int
	PageSize         int
}

// NotificationStats represents notification statistics
type NotificationStats struct {
	TotalSent      int64 `json:"total_sent"`
	TotalDelivered int64 `json:"total_delivered"`
	TotalFailed    int64 `json:"total_failed"`
	TotalPending   int64 `json:"total_pending"`
}

// CustomerNotificationPreferenceRepository defines the interface for customer notification preferences
type CustomerNotificationPreferenceRepository interface {
	// Create creates a new preference
	Create(ctx context.Context, pref *domain.CustomerNotificationPreference) error

	// GetByID retrieves a preference by ID
	GetByID(ctx context.Context, id int64) (*domain.CustomerNotificationPreference, error)

	// GetByCustomerTypeAndChannel retrieves a specific preference
	GetByCustomerTypeAndChannel(ctx context.Context, customerID int64, notificationType, channel string) (*domain.CustomerNotificationPreference, error)

	// Update updates an existing preference
	Update(ctx context.Context, pref *domain.CustomerNotificationPreference) error

	// Upsert creates or updates a preference
	Upsert(ctx context.Context, pref *domain.CustomerNotificationPreference) error

	// Delete deletes a preference
	Delete(ctx context.Context, id int64) error

	// ListByCustomer retrieves all preferences for a customer
	ListByCustomer(ctx context.Context, customerID int64) ([]*domain.CustomerNotificationPreference, error)

	// IsEnabled checks if a notification type/channel is enabled for a customer
	IsEnabled(ctx context.Context, customerID int64, notificationType, channel string) (bool, error)

	// BulkUpsert creates or updates multiple preferences
	BulkUpsert(ctx context.Context, customerID int64, prefs []*domain.CustomerNotificationPreference) error
}

// InternalNotificationRepository defines the interface for internal notification operations
type InternalNotificationRepository interface {
	// Create creates a new internal notification
	Create(ctx context.Context, notification *domain.InternalNotification) error

	// GetByID retrieves an internal notification by ID
	GetByID(ctx context.Context, id int64) (*domain.InternalNotification, error)

	// Update updates an existing notification
	Update(ctx context.Context, notification *domain.InternalNotification) error

	// Delete deletes an internal notification
	Delete(ctx context.Context, id int64) error

	// List retrieves internal notifications with filtering
	List(ctx context.Context, filter InternalNotificationFilter) ([]*domain.InternalNotification, int64, error)

	// ListByUser retrieves notifications for a user
	ListByUser(ctx context.Context, userID int64, filter InternalNotificationFilter) ([]*domain.InternalNotification, int64, error)

	// ListUnreadByUser retrieves unread notifications for a user
	ListUnreadByUser(ctx context.Context, userID int64, limit int) ([]*domain.InternalNotification, error)

	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, id int64) error

	// MarkAllAsRead marks all notifications as read for a user
	MarkAllAsRead(ctx context.Context, userID int64) error

	// GetUnreadCount retrieves the count of unread notifications for a user
	GetUnreadCount(ctx context.Context, userID int64) (int64, error)

	// CreateBulk creates notifications for multiple users
	CreateBulk(ctx context.Context, notifications []*domain.InternalNotification) error

	// DeleteOlderThan deletes notifications older than a date
	DeleteOlderThan(ctx context.Context, olderThan time.Time) (int64, error)
}

// InternalNotificationFilter contains filters for listing internal notifications
type InternalNotificationFilter struct {
	UserID        *int64
	BranchID      *int64
	Type          *string
	IsRead        *bool
	ReferenceType *string
	ReferenceID   *int64
	DateFrom      *string
	DateTo        *string
	Page          int
	PageSize      int
}
