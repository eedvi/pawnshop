package service

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"text/template"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

var (
	ErrNotificationNotFound     = errors.New("notification not found")
	ErrTemplateNotFound         = errors.New("notification template not found")
	ErrNotificationNotPending   = errors.New("notification is not pending")
	ErrNotificationCustomerNotFound = errors.New("customer not found for notification")
)

// NotificationService defines the interface for notification operations
type NotificationService interface {
	// Template operations
	CreateTemplate(ctx context.Context, req CreateNotificationTemplateRequest) (*domain.NotificationTemplate, error)
	GetTemplateByID(ctx context.Context, id int64) (*domain.NotificationTemplate, error)
	UpdateTemplate(ctx context.Context, id int64, req UpdateNotificationTemplateRequest) (*domain.NotificationTemplate, error)
	DeleteTemplate(ctx context.Context, id int64) error
	ListTemplates(ctx context.Context, includeInactive bool) ([]*domain.NotificationTemplate, error)

	// Notification operations
	Create(ctx context.Context, req CreateNotificationRequest) (*domain.Notification, error)
	CreateFromTemplate(ctx context.Context, req CreateNotificationFromTemplateRequest) (*domain.Notification, error)
	GetByID(ctx context.Context, id int64) (*domain.Notification, error)
	List(ctx context.Context, filter repository.NotificationFilter) ([]*domain.Notification, int64, error)
	ListByCustomer(ctx context.Context, customerID int64, filter repository.NotificationFilter) ([]*domain.Notification, int64, error)
	Cancel(ctx context.Context, id int64) error

	// Queue processing
	GetPendingNotifications(ctx context.Context, limit int) ([]*domain.Notification, error)
	GetScheduledNotifications(ctx context.Context, before time.Time, limit int) ([]*domain.Notification, error)
	GetFailedNotifications(ctx context.Context, limit int) ([]*domain.Notification, error)
	MarkAsSent(ctx context.Context, id int64) error
	MarkAsDelivered(ctx context.Context, id int64) error
	MarkAsFailed(ctx context.Context, id int64, reason string) error
	RetryNotification(ctx context.Context, id int64) error

	// Customer preferences
	GetCustomerPreferences(ctx context.Context, customerID int64) ([]*domain.CustomerNotificationPreference, error)
	UpdateCustomerPreferences(ctx context.Context, customerID int64, prefs []*domain.CustomerNotificationPreference) error
	IsChannelEnabled(ctx context.Context, customerID int64, notificationType, channel string) (bool, error)

	// Internal notifications
	CreateInternalNotification(ctx context.Context, req CreateInternalNotificationRequest) (*domain.InternalNotification, error)
	GetInternalNotificationByID(ctx context.Context, id int64) (*domain.InternalNotification, error)
	ListInternalNotificationsByUser(ctx context.Context, userID int64, filter repository.InternalNotificationFilter) ([]*domain.InternalNotification, int64, error)
	GetUnreadInternalNotifications(ctx context.Context, userID int64, limit int) ([]*domain.InternalNotification, error)
	MarkInternalNotificationAsRead(ctx context.Context, id int64) error
	MarkAllInternalNotificationsAsRead(ctx context.Context, userID int64) error
	GetUnreadCount(ctx context.Context, userID int64) (int64, error)

	// Bulk operations
	NotifyBranchUsers(ctx context.Context, branchID int64, title, message, notificationType string) error

	// Simple send operations
	SendToCustomer(ctx context.Context, req SendNotificationRequest) (*domain.Notification, error)

	// Stats
	GetStatsByCustomer(ctx context.Context, customerID int64) (*repository.NotificationStats, error)
	GetStatsByBranch(ctx context.Context, branchID int64, dateFrom, dateTo time.Time) (*repository.NotificationStats, error)
}

type notificationService struct {
	notificationRepo         repository.NotificationRepository
	templateRepo             repository.NotificationTemplateRepository
	preferenceRepo           repository.CustomerNotificationPreferenceRepository
	internalNotificationRepo repository.InternalNotificationRepository
	customerRepo             repository.CustomerRepository
	userRepo                 repository.UserRepository
}

// NewNotificationService creates a new notification service
func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	templateRepo repository.NotificationTemplateRepository,
	preferenceRepo repository.CustomerNotificationPreferenceRepository,
	internalNotificationRepo repository.InternalNotificationRepository,
	customerRepo repository.CustomerRepository,
	userRepo repository.UserRepository,
) NotificationService {
	return &notificationService{
		notificationRepo:         notificationRepo,
		templateRepo:             templateRepo,
		preferenceRepo:           preferenceRepo,
		internalNotificationRepo: internalNotificationRepo,
		customerRepo:             customerRepo,
		userRepo:                 userRepo,
	}
}

// Request types
type CreateNotificationTemplateRequest struct {
	NotificationType string `json:"notification_type" validate:"required"`
	Channel          string `json:"channel" validate:"required"`
	Name             string `json:"name" validate:"required"`
	Subject          string `json:"subject"`
	BodyTemplate     string `json:"body_template" validate:"required"`
}

type UpdateNotificationTemplateRequest struct {
	Name         string `json:"name"`
	Subject      string `json:"subject"`
	BodyTemplate string `json:"body_template"`
	IsActive     *bool  `json:"is_active"`
}

type CreateNotificationRequest struct {
	CustomerID       int64      `json:"customer_id" validate:"required"`
	BranchID         *int64     `json:"branch_id"`
	NotificationType string     `json:"notification_type" validate:"required"`
	Channel          string     `json:"channel" validate:"required"`
	Subject          string     `json:"subject"`
	Body             string     `json:"body" validate:"required"`
	ReferenceType    string     `json:"reference_type"`
	ReferenceID      *int64     `json:"reference_id"`
	ScheduledFor     *time.Time `json:"scheduled_for"`
}

type CreateNotificationFromTemplateRequest struct {
	CustomerID       int64             `json:"customer_id" validate:"required"`
	BranchID         *int64            `json:"branch_id"`
	NotificationType string            `json:"notification_type" validate:"required"`
	Channel          string            `json:"channel" validate:"required"`
	TemplateData     map[string]string `json:"template_data"`
	ReferenceType    string            `json:"reference_type"`
	ReferenceID      *int64            `json:"reference_id"`
	ScheduledFor     *time.Time        `json:"scheduled_for"`
}

type CreateInternalNotificationRequest struct {
	UserID        int64  `json:"user_id" validate:"required"`
	BranchID      *int64 `json:"branch_id"`
	Title         string `json:"title" validate:"required"`
	Message       string `json:"message" validate:"required"`
	Type          string `json:"type" validate:"required"` // info, warning, error, success
	ReferenceType string `json:"reference_type"`
	ReferenceID   *int64 `json:"reference_id"`
	ActionURL     string `json:"action_url"`
}

// SendNotificationRequest is a simplified request for sending notifications
type SendNotificationRequest struct {
	CustomerID    int64   `json:"customer_id" validate:"required"`
	Type          string  `json:"type" validate:"required"`
	Title         string  `json:"title" validate:"required"`
	Message       string  `json:"message" validate:"required"`
	Channel       string  `json:"channel" validate:"required"` // sms, email, whatsapp
	ReferenceType *string `json:"reference_type"`
	ReferenceID   *int64  `json:"reference_id"`
}

// Template operations
func (s *notificationService) CreateTemplate(ctx context.Context, req CreateNotificationTemplateRequest) (*domain.NotificationTemplate, error) {
	template := &domain.NotificationTemplate{
		NotificationType: req.NotificationType,
		Channel:          req.Channel,
		Name:             req.Name,
		Subject:          req.Subject,
		BodyTemplate:     req.BodyTemplate,
		IsActive:         true,
	}

	if err := s.templateRepo.Create(ctx, template); err != nil {
		return nil, err
	}

	return template, nil
}

func (s *notificationService) GetTemplateByID(ctx context.Context, id int64) (*domain.NotificationTemplate, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, ErrTemplateNotFound
	}
	return template, nil
}

func (s *notificationService) UpdateTemplate(ctx context.Context, id int64, req UpdateNotificationTemplateRequest) (*domain.NotificationTemplate, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, ErrTemplateNotFound
	}

	if req.Name != "" {
		template.Name = req.Name
	}
	if req.Subject != "" {
		template.Subject = req.Subject
	}
	if req.BodyTemplate != "" {
		template.BodyTemplate = req.BodyTemplate
	}
	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	}

	if err := s.templateRepo.Update(ctx, template); err != nil {
		return nil, err
	}

	return template, nil
}

func (s *notificationService) DeleteTemplate(ctx context.Context, id int64) error {
	return s.templateRepo.Delete(ctx, id)
}

func (s *notificationService) ListTemplates(ctx context.Context, includeInactive bool) ([]*domain.NotificationTemplate, error) {
	return s.templateRepo.List(ctx, includeInactive)
}

// Notification operations
func (s *notificationService) Create(ctx context.Context, req CreateNotificationRequest) (*domain.Notification, error) {
	// Check customer preferences
	enabled, err := s.preferenceRepo.IsEnabled(ctx, req.CustomerID, req.NotificationType, req.Channel)
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, errors.New("notification channel is disabled for this customer")
	}

	notification := &domain.Notification{
		CustomerID:       req.CustomerID,
		BranchID:         req.BranchID,
		NotificationType: req.NotificationType,
		Channel:          req.Channel,
		Subject:          req.Subject,
		Body:             req.Body,
		ReferenceType:    req.ReferenceType,
		ReferenceID:      req.ReferenceID,
		Status:           domain.NotificationStatusPending,
		ScheduledFor:     req.ScheduledFor,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *notificationService) CreateFromTemplate(ctx context.Context, req CreateNotificationFromTemplateRequest) (*domain.Notification, error) {
	// Get template
	tmpl, err := s.templateRepo.GetByTypeAndChannel(ctx, req.NotificationType, req.Channel)
	if err != nil {
		return nil, err
	}
	if tmpl == nil {
		return nil, ErrTemplateNotFound
	}

	// Check customer preferences
	enabled, err := s.preferenceRepo.IsEnabled(ctx, req.CustomerID, req.NotificationType, req.Channel)
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, errors.New("notification channel is disabled for this customer")
	}

	// Render template
	subject := s.renderTemplate(tmpl.Subject, req.TemplateData)
	body := s.renderTemplate(tmpl.BodyTemplate, req.TemplateData)

	notification := &domain.Notification{
		CustomerID:       req.CustomerID,
		BranchID:         req.BranchID,
		NotificationType: req.NotificationType,
		Channel:          req.Channel,
		Subject:          subject,
		Body:             body,
		ReferenceType:    req.ReferenceType,
		ReferenceID:      req.ReferenceID,
		Status:           domain.NotificationStatusPending,
		ScheduledFor:     req.ScheduledFor,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *notificationService) renderTemplate(tmpl string, data map[string]string) string {
	if data == nil {
		return tmpl
	}

	result := tmpl
	for key, value := range data {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

func (s *notificationService) renderGoTemplate(tmpl string, data map[string]interface{}) (string, error) {
	t, err := template.New("notification").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *notificationService) GetByID(ctx context.Context, id int64) (*domain.Notification, error) {
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if notification == nil {
		return nil, ErrNotificationNotFound
	}

	// Load customer
	customer, _ := s.customerRepo.GetByID(ctx, notification.CustomerID)
	notification.Customer = customer

	return notification, nil
}

func (s *notificationService) List(ctx context.Context, filter repository.NotificationFilter) ([]*domain.Notification, int64, error) {
	return s.notificationRepo.List(ctx, filter)
}

func (s *notificationService) ListByCustomer(ctx context.Context, customerID int64, filter repository.NotificationFilter) ([]*domain.Notification, int64, error) {
	return s.notificationRepo.ListByCustomer(ctx, customerID, filter)
}

func (s *notificationService) Cancel(ctx context.Context, id int64) error {
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if notification == nil {
		return ErrNotificationNotFound
	}
	if !notification.IsPending() {
		return ErrNotificationNotPending
	}

	return s.notificationRepo.Cancel(ctx, id)
}

// Queue processing
func (s *notificationService) GetPendingNotifications(ctx context.Context, limit int) ([]*domain.Notification, error) {
	return s.notificationRepo.ListPending(ctx, limit)
}

func (s *notificationService) GetScheduledNotifications(ctx context.Context, before time.Time, limit int) ([]*domain.Notification, error) {
	return s.notificationRepo.ListScheduled(ctx, before, limit)
}

func (s *notificationService) GetFailedNotifications(ctx context.Context, limit int) ([]*domain.Notification, error) {
	return s.notificationRepo.ListFailed(ctx, 3, limit) // Max 3 retries
}

func (s *notificationService) MarkAsSent(ctx context.Context, id int64) error {
	return s.notificationRepo.MarkAsSent(ctx, id)
}

func (s *notificationService) MarkAsDelivered(ctx context.Context, id int64) error {
	return s.notificationRepo.MarkAsDelivered(ctx, id)
}

func (s *notificationService) MarkAsFailed(ctx context.Context, id int64, reason string) error {
	return s.notificationRepo.MarkAsFailed(ctx, id, reason)
}

func (s *notificationService) RetryNotification(ctx context.Context, id int64) error {
	notification, err := s.notificationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if notification == nil {
		return ErrNotificationNotFound
	}
	if !notification.CanRetry() {
		return errors.New("notification cannot be retried")
	}

	return s.notificationRepo.IncrementRetry(ctx, id)
}

// Customer preferences
func (s *notificationService) GetCustomerPreferences(ctx context.Context, customerID int64) ([]*domain.CustomerNotificationPreference, error) {
	return s.preferenceRepo.ListByCustomer(ctx, customerID)
}

func (s *notificationService) UpdateCustomerPreferences(ctx context.Context, customerID int64, prefs []*domain.CustomerNotificationPreference) error {
	return s.preferenceRepo.BulkUpsert(ctx, customerID, prefs)
}

func (s *notificationService) IsChannelEnabled(ctx context.Context, customerID int64, notificationType, channel string) (bool, error) {
	return s.preferenceRepo.IsEnabled(ctx, customerID, notificationType, channel)
}

// Internal notifications
func (s *notificationService) CreateInternalNotification(ctx context.Context, req CreateInternalNotificationRequest) (*domain.InternalNotification, error) {
	notification := &domain.InternalNotification{
		UserID:        req.UserID,
		BranchID:      req.BranchID,
		Title:         req.Title,
		Message:       req.Message,
		Type:          req.Type,
		ReferenceType: req.ReferenceType,
		ReferenceID:   req.ReferenceID,
		ActionURL:     req.ActionURL,
	}

	if err := s.internalNotificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *notificationService) GetInternalNotificationByID(ctx context.Context, id int64) (*domain.InternalNotification, error) {
	notification, err := s.internalNotificationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if notification == nil {
		return nil, ErrNotificationNotFound
	}
	return notification, nil
}

func (s *notificationService) ListInternalNotificationsByUser(ctx context.Context, userID int64, filter repository.InternalNotificationFilter) ([]*domain.InternalNotification, int64, error) {
	return s.internalNotificationRepo.ListByUser(ctx, userID, filter)
}

func (s *notificationService) GetUnreadInternalNotifications(ctx context.Context, userID int64, limit int) ([]*domain.InternalNotification, error) {
	return s.internalNotificationRepo.ListUnreadByUser(ctx, userID, limit)
}

func (s *notificationService) MarkInternalNotificationAsRead(ctx context.Context, id int64) error {
	return s.internalNotificationRepo.MarkAsRead(ctx, id)
}

func (s *notificationService) MarkAllInternalNotificationsAsRead(ctx context.Context, userID int64) error {
	return s.internalNotificationRepo.MarkAllAsRead(ctx, userID)
}

func (s *notificationService) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	return s.internalNotificationRepo.GetUnreadCount(ctx, userID)
}

// Bulk operations
func (s *notificationService) NotifyBranchUsers(ctx context.Context, branchID int64, title, message, notificationType string) error {
	// Get all users for the branch
	result, err := s.userRepo.List(ctx, repository.UserListParams{BranchID: &branchID})
	if err != nil {
		return err
	}

	var notifications []*domain.InternalNotification
	for _, user := range result.Data {
		notifications = append(notifications, &domain.InternalNotification{
			UserID:   user.ID,
			BranchID: &branchID,
			Title:    title,
			Message:  message,
			Type:     notificationType,
		})
	}

	return s.internalNotificationRepo.CreateBulk(ctx, notifications)
}

// Simple send operations
func (s *notificationService) SendToCustomer(ctx context.Context, req SendNotificationRequest) (*domain.Notification, error) {
	// Check if customer exists
	customer, err := s.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrNotificationCustomerNotFound
	}

	// Check customer preferences - if preferences don't exist, allow notification
	enabled, err := s.preferenceRepo.IsEnabled(ctx, req.CustomerID, req.Type, req.Channel)
	if err == nil && !enabled {
		// Preferences exist and disabled - skip but don't error
		return nil, nil
	}

	var refType, body string
	if req.ReferenceType != nil {
		refType = *req.ReferenceType
	}

	// Use message as body
	body = req.Message

	notification := &domain.Notification{
		CustomerID:       req.CustomerID,
		BranchID:         &customer.BranchID,
		NotificationType: req.Type,
		Channel:          req.Channel,
		Subject:          req.Title,
		Body:             body,
		ReferenceType:    refType,
		ReferenceID:      req.ReferenceID,
		Status:           domain.NotificationStatusPending,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

// Stats
func (s *notificationService) GetStatsByCustomer(ctx context.Context, customerID int64) (*repository.NotificationStats, error) {
	return s.notificationRepo.GetStatsByCustomer(ctx, customerID)
}

func (s *notificationService) GetStatsByBranch(ctx context.Context, branchID int64, dateFrom, dateTo time.Time) (*repository.NotificationStats, error) {
	return s.notificationRepo.GetStatsByBranch(ctx, branchID, dateFrom, dateTo)
}
