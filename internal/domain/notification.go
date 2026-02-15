package domain

import "time"

// Notification types
const (
	NotificationTypeLoanDueReminder   = "loan_due_reminder"
	NotificationTypeLoanOverdue       = "loan_overdue"
	NotificationTypeMinimumPaymentDue = "minimum_payment_due"
	NotificationTypePaymentReceived   = "payment_received"
	NotificationTypeLoanConfiscated   = "loan_confiscated"
	NotificationTypeItemForSale       = "item_for_sale"
	NotificationTypeItemSold          = "item_sold"
	NotificationTypePromotion         = "promotion"
	NotificationTypeLoyaltyPoints     = "loyalty_points"
	NotificationTypeGeneral           = "general"
)

// Notification channels
const (
	NotificationChannelEmail    = "email"
	NotificationChannelSMS      = "sms"
	NotificationChannelWhatsApp = "whatsapp"
	NotificationChannelPush     = "push"
	NotificationChannelInternal = "internal"
)

// Notification status
const (
	NotificationStatusPending   = "pending"
	NotificationStatusSent      = "sent"
	NotificationStatusDelivered = "delivered"
	NotificationStatusFailed    = "failed"
	NotificationStatusCancelled = "cancelled"
)

// NotificationTemplate represents a template for notifications
type NotificationTemplate struct {
	ID              int64  `json:"id"`
	NotificationType string `json:"notification_type"`
	Channel         string `json:"channel"`
	Name            string `json:"name"`
	Subject         string `json:"subject,omitempty"`
	BodyTemplate    string `json:"body_template"`
	IsActive        bool   `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Notification represents a notification to be sent to a customer
type Notification struct {
	ID               int64     `json:"id"`
	CustomerID       int64     `json:"customer_id"`
	Customer         *Customer `json:"customer,omitempty"`
	BranchID         *int64    `json:"branch_id,omitempty"`
	Branch           *Branch   `json:"branch,omitempty"`
	NotificationType string    `json:"notification_type"`
	Channel          string    `json:"channel"`

	// Content
	Subject string `json:"subject,omitempty"`
	Body    string `json:"body"`

	// Reference
	ReferenceType string `json:"reference_type,omitempty"`
	ReferenceID   *int64 `json:"reference_id,omitempty"`

	// Delivery
	Status       string     `json:"status"`
	ScheduledFor *time.Time `json:"scheduled_for,omitempty"`
	SentAt       *time.Time `json:"sent_at,omitempty"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty"`
	FailedAt     *time.Time `json:"failed_at,omitempty"`
	FailureReason string    `json:"failure_reason,omitempty"`
	RetryCount   int        `json:"retry_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsPending checks if notification is pending
func (n *Notification) IsPending() bool {
	return n.Status == NotificationStatusPending
}

// IsSent checks if notification was sent
func (n *Notification) IsSent() bool {
	return n.Status == NotificationStatusSent
}

// IsDelivered checks if notification was delivered
func (n *Notification) IsDelivered() bool {
	return n.Status == NotificationStatusDelivered
}

// IsFailed checks if notification failed
func (n *Notification) IsFailed() bool {
	return n.Status == NotificationStatusFailed
}

// CanRetry checks if notification can be retried
func (n *Notification) CanRetry() bool {
	return n.Status == NotificationStatusFailed && n.RetryCount < 3
}

// CustomerNotificationPreference represents customer preferences for notifications
type CustomerNotificationPreference struct {
	ID               int64  `json:"id"`
	CustomerID       int64  `json:"customer_id"`
	NotificationType string `json:"notification_type"`
	Channel          string `json:"channel"`
	IsEnabled        bool   `json:"is_enabled"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// InternalNotification represents a notification for internal users
type InternalNotification struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	User     *User  `json:"user,omitempty"`
	BranchID *int64 `json:"branch_id,omitempty"`
	Branch   *Branch `json:"branch,omitempty"`

	// Content
	Title   string `json:"title"`
	Message string `json:"message"`
	Type    string `json:"type"` // info, warning, error, success

	// Reference
	ReferenceType string `json:"reference_type,omitempty"`
	ReferenceID   *int64 `json:"reference_id,omitempty"`
	ActionURL     string `json:"action_url,omitempty"`

	// Status
	IsRead   bool       `json:"is_read"`
	ReadAt   *time.Time `json:"read_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// MarkAsRead marks the notification as read
func (n *InternalNotification) MarkAsRead() {
	n.IsRead = true
	now := time.Now()
	n.ReadAt = &now
}
