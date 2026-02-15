package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/internal/repository/mocks"
)

func setupNotificationService() (
	NotificationService,
	*mocks.MockNotificationRepository,
	*mocks.MockNotificationTemplateRepository,
	*mocks.MockCustomerNotificationPreferenceRepository,
	*mocks.MockInternalNotificationRepository,
	*mocks.MockCustomerRepository,
	*mocks.MockUserRepository,
) {
	notificationRepo := new(mocks.MockNotificationRepository)
	templateRepo := new(mocks.MockNotificationTemplateRepository)
	preferenceRepo := new(mocks.MockCustomerNotificationPreferenceRepository)
	internalRepo := new(mocks.MockInternalNotificationRepository)
	customerRepo := new(mocks.MockCustomerRepository)
	userRepo := new(mocks.MockUserRepository)

	service := NewNotificationService(
		notificationRepo,
		templateRepo,
		preferenceRepo,
		internalRepo,
		customerRepo,
		userRepo,
	)

	return service, notificationRepo, templateRepo, preferenceRepo, internalRepo, customerRepo, userRepo
}

// Template Tests

func TestNotificationService_CreateTemplate_Success(t *testing.T) {
	service, _, templateRepo, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	templateRepo.On("Create", ctx, mock.AnythingOfType("*domain.NotificationTemplate")).Return(nil)

	req := CreateNotificationTemplateRequest{
		NotificationType: domain.NotificationTypeLoanDueReminder,
		Channel:          domain.NotificationChannelSMS,
		Name:             "Loan Due Reminder SMS",
		Subject:          "",
		BodyTemplate:     "Dear {{customer_name}}, your loan {{loan_number}} is due on {{due_date}}.",
	}

	result, err := service.CreateTemplate(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.NotificationTypeLoanDueReminder, result.NotificationType)
	assert.Equal(t, domain.NotificationChannelSMS, result.Channel)
	assert.True(t, result.IsActive)
	templateRepo.AssertExpectations(t)
}

func TestNotificationService_GetTemplateByID_Success(t *testing.T) {
	service, _, templateRepo, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	template := &domain.NotificationTemplate{
		ID:               1,
		NotificationType: domain.NotificationTypeLoanDueReminder,
		Channel:          domain.NotificationChannelSMS,
		Name:             "Test Template",
		IsActive:         true,
	}

	templateRepo.On("GetByID", ctx, int64(1)).Return(template, nil)

	result, err := service.GetTemplateByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test Template", result.Name)
	templateRepo.AssertExpectations(t)
}

func TestNotificationService_GetTemplateByID_NotFound(t *testing.T) {
	service, _, templateRepo, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	templateRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	result, err := service.GetTemplateByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTemplateNotFound, err)
	templateRepo.AssertExpectations(t)
}

func TestNotificationService_UpdateTemplate_Success(t *testing.T) {
	service, _, templateRepo, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	template := &domain.NotificationTemplate{
		ID:       1,
		Name:     "Old Name",
		IsActive: true,
	}

	templateRepo.On("GetByID", ctx, int64(1)).Return(template, nil)
	templateRepo.On("Update", ctx, mock.AnythingOfType("*domain.NotificationTemplate")).Return(nil)

	isActive := false
	req := UpdateNotificationTemplateRequest{
		Name:     "New Name",
		IsActive: &isActive,
	}

	result, err := service.UpdateTemplate(ctx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Name", result.Name)
	assert.False(t, result.IsActive)
	templateRepo.AssertExpectations(t)
}

func TestNotificationService_DeleteTemplate_Success(t *testing.T) {
	service, _, templateRepo, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	templateRepo.On("Delete", ctx, int64(1)).Return(nil)

	err := service.DeleteTemplate(ctx, 1)

	assert.NoError(t, err)
	templateRepo.AssertExpectations(t)
}

func TestNotificationService_ListTemplates_Success(t *testing.T) {
	service, _, templateRepo, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	templates := []*domain.NotificationTemplate{
		{ID: 1, Name: "Template 1"},
		{ID: 2, Name: "Template 2"},
	}

	templateRepo.On("List", ctx, false).Return(templates, nil)

	result, err := service.ListTemplates(ctx, false)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	templateRepo.AssertExpectations(t)
}

// Notification Tests

func TestNotificationService_Create_Success(t *testing.T) {
	service, notificationRepo, _, preferenceRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	preferenceRepo.On("IsEnabled", ctx, int64(1), domain.NotificationTypeLoanDueReminder, domain.NotificationChannelSMS).Return(true, nil)
	notificationRepo.On("Create", ctx, mock.AnythingOfType("*domain.Notification")).Return(nil)

	req := CreateNotificationRequest{
		CustomerID:       1,
		NotificationType: domain.NotificationTypeLoanDueReminder,
		Channel:          domain.NotificationChannelSMS,
		Subject:          "Payment Reminder",
		Body:             "Your loan is due tomorrow.",
	}

	result, err := service.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.NotificationStatusPending, result.Status)
	notificationRepo.AssertExpectations(t)
	preferenceRepo.AssertExpectations(t)
}

func TestNotificationService_Create_ChannelDisabled(t *testing.T) {
	service, _, _, preferenceRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	preferenceRepo.On("IsEnabled", ctx, int64(1), domain.NotificationTypeLoanDueReminder, domain.NotificationChannelSMS).Return(false, nil)

	req := CreateNotificationRequest{
		CustomerID:       1,
		NotificationType: domain.NotificationTypeLoanDueReminder,
		Channel:          domain.NotificationChannelSMS,
		Body:             "Test message",
	}

	result, err := service.Create(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "disabled")
	preferenceRepo.AssertExpectations(t)
}

func TestNotificationService_CreateFromTemplate_Success(t *testing.T) {
	service, notificationRepo, templateRepo, preferenceRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	template := &domain.NotificationTemplate{
		ID:               1,
		NotificationType: domain.NotificationTypeLoanDueReminder,
		Channel:          domain.NotificationChannelSMS,
		Subject:          "Loan Due: {{loan_number}}",
		BodyTemplate:     "Dear {{customer_name}}, your loan {{loan_number}} is due.",
		IsActive:         true,
	}

	templateRepo.On("GetByTypeAndChannel", ctx, domain.NotificationTypeLoanDueReminder, domain.NotificationChannelSMS).Return(template, nil)
	preferenceRepo.On("IsEnabled", ctx, int64(1), domain.NotificationTypeLoanDueReminder, domain.NotificationChannelSMS).Return(true, nil)
	notificationRepo.On("Create", ctx, mock.AnythingOfType("*domain.Notification")).Return(nil)

	req := CreateNotificationFromTemplateRequest{
		CustomerID:       1,
		NotificationType: domain.NotificationTypeLoanDueReminder,
		Channel:          domain.NotificationChannelSMS,
		TemplateData: map[string]string{
			"customer_name": "John Doe",
			"loan_number":   "LN-001",
		},
	}

	result, err := service.CreateFromTemplate(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.Body, "John Doe")
	assert.Contains(t, result.Body, "LN-001")
	templateRepo.AssertExpectations(t)
	preferenceRepo.AssertExpectations(t)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_GetByID_Success(t *testing.T) {
	service, notificationRepo, _, _, _, customerRepo, _ := setupNotificationService()
	ctx := context.Background()

	notification := &domain.Notification{
		ID:               1,
		CustomerID:       1,
		NotificationType: domain.NotificationTypeLoanDueReminder,
		Status:           domain.NotificationStatusPending,
	}
	customer := &domain.Customer{ID: 1, FirstName: "John"}

	notificationRepo.On("GetByID", ctx, int64(1)).Return(notification, nil)
	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Customer)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_GetByID_NotFound(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notificationRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	result, err := service.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrNotificationNotFound, err)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_Cancel_Success(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notification := &domain.Notification{
		ID:     1,
		Status: domain.NotificationStatusPending,
	}

	notificationRepo.On("GetByID", ctx, int64(1)).Return(notification, nil)
	notificationRepo.On("Cancel", ctx, int64(1)).Return(nil)

	err := service.Cancel(ctx, 1)

	assert.NoError(t, err)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_Cancel_NotPending(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notification := &domain.Notification{
		ID:     1,
		Status: domain.NotificationStatusSent, // Already sent
	}

	notificationRepo.On("GetByID", ctx, int64(1)).Return(notification, nil)

	err := service.Cancel(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, ErrNotificationNotPending, err)
	notificationRepo.AssertExpectations(t)
}

// Queue Processing Tests

func TestNotificationService_GetPendingNotifications_Success(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notifications := []*domain.Notification{
		{ID: 1, Status: domain.NotificationStatusPending},
		{ID: 2, Status: domain.NotificationStatusPending},
	}

	notificationRepo.On("ListPending", ctx, 10).Return(notifications, nil)

	result, err := service.GetPendingNotifications(ctx, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_MarkAsSent_Success(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notificationRepo.On("MarkAsSent", ctx, int64(1)).Return(nil)

	err := service.MarkAsSent(ctx, 1)

	assert.NoError(t, err)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_MarkAsDelivered_Success(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notificationRepo.On("MarkAsDelivered", ctx, int64(1)).Return(nil)

	err := service.MarkAsDelivered(ctx, 1)

	assert.NoError(t, err)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_MarkAsFailed_Success(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notificationRepo.On("MarkAsFailed", ctx, int64(1), "Connection timeout").Return(nil)

	err := service.MarkAsFailed(ctx, 1, "Connection timeout")

	assert.NoError(t, err)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_RetryNotification_Success(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notification := &domain.Notification{
		ID:         1,
		Status:     domain.NotificationStatusFailed,
		RetryCount: 1,
	}

	notificationRepo.On("GetByID", ctx, int64(1)).Return(notification, nil)
	notificationRepo.On("IncrementRetry", ctx, int64(1)).Return(nil)

	err := service.RetryNotification(ctx, 1)

	assert.NoError(t, err)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_RetryNotification_MaxRetriesReached(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notification := &domain.Notification{
		ID:         1,
		Status:     domain.NotificationStatusFailed,
		RetryCount: 3, // Max retries reached
	}

	notificationRepo.On("GetByID", ctx, int64(1)).Return(notification, nil)

	err := service.RetryNotification(ctx, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be retried")
	notificationRepo.AssertExpectations(t)
}

// Customer Preferences Tests

func TestNotificationService_GetCustomerPreferences_Success(t *testing.T) {
	service, _, _, preferenceRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	prefs := []*domain.CustomerNotificationPreference{
		{ID: 1, CustomerID: 1, NotificationType: "loan_due_reminder", Channel: "sms", IsEnabled: true},
		{ID: 2, CustomerID: 1, NotificationType: "loan_due_reminder", Channel: "email", IsEnabled: false},
	}

	preferenceRepo.On("ListByCustomer", ctx, int64(1)).Return(prefs, nil)

	result, err := service.GetCustomerPreferences(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	preferenceRepo.AssertExpectations(t)
}

func TestNotificationService_IsChannelEnabled_Success(t *testing.T) {
	service, _, _, preferenceRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	preferenceRepo.On("IsEnabled", ctx, int64(1), "loan_due_reminder", "sms").Return(true, nil)

	enabled, err := service.IsChannelEnabled(ctx, 1, "loan_due_reminder", "sms")

	assert.NoError(t, err)
	assert.True(t, enabled)
	preferenceRepo.AssertExpectations(t)
}

// Internal Notification Tests

func TestNotificationService_CreateInternalNotification_Success(t *testing.T) {
	service, _, _, _, internalRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	internalRepo.On("Create", ctx, mock.AnythingOfType("*domain.InternalNotification")).Return(nil)

	req := CreateInternalNotificationRequest{
		UserID:  100,
		Title:   "New Loan Created",
		Message: "A new loan has been created for customer John Doe.",
		Type:    "info",
	}

	result, err := service.CreateInternalNotification(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(100), result.UserID)
	assert.Equal(t, "info", result.Type)
	internalRepo.AssertExpectations(t)
}

func TestNotificationService_GetUnreadInternalNotifications_Success(t *testing.T) {
	service, _, _, _, internalRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	notifications := []*domain.InternalNotification{
		{ID: 1, UserID: 100, Title: "Notification 1", IsRead: false},
		{ID: 2, UserID: 100, Title: "Notification 2", IsRead: false},
	}

	internalRepo.On("ListUnreadByUser", ctx, int64(100), 10).Return(notifications, nil)

	result, err := service.GetUnreadInternalNotifications(ctx, 100, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	internalRepo.AssertExpectations(t)
}

func TestNotificationService_MarkInternalNotificationAsRead_Success(t *testing.T) {
	service, _, _, _, internalRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	internalRepo.On("MarkAsRead", ctx, int64(1)).Return(nil)

	err := service.MarkInternalNotificationAsRead(ctx, 1)

	assert.NoError(t, err)
	internalRepo.AssertExpectations(t)
}

func TestNotificationService_MarkAllInternalNotificationsAsRead_Success(t *testing.T) {
	service, _, _, _, internalRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	internalRepo.On("MarkAllAsRead", ctx, int64(100)).Return(nil)

	err := service.MarkAllInternalNotificationsAsRead(ctx, 100)

	assert.NoError(t, err)
	internalRepo.AssertExpectations(t)
}

func TestNotificationService_GetUnreadCount_Success(t *testing.T) {
	service, _, _, _, internalRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	internalRepo.On("GetUnreadCount", ctx, int64(100)).Return(int64(5), nil)

	count, err := service.GetUnreadCount(ctx, 100)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	internalRepo.AssertExpectations(t)
}

// SendToCustomer Tests

func TestNotificationService_SendToCustomer_Success(t *testing.T) {
	service, notificationRepo, _, preferenceRepo, _, customerRepo, _ := setupNotificationService()
	ctx := context.Background()

	customer := &domain.Customer{ID: 1, BranchID: 1, FirstName: "John"}

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	preferenceRepo.On("IsEnabled", ctx, int64(1), "general", "sms").Return(true, nil)
	notificationRepo.On("Create", ctx, mock.AnythingOfType("*domain.Notification")).Return(nil)

	req := SendNotificationRequest{
		CustomerID: 1,
		Type:       "general",
		Title:      "Test Notification",
		Message:    "This is a test message.",
		Channel:    "sms",
	}

	result, err := service.SendToCustomer(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.NotificationStatusPending, result.Status)
	customerRepo.AssertExpectations(t)
	preferenceRepo.AssertExpectations(t)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_SendToCustomer_CustomerNotFound(t *testing.T) {
	service, _, _, _, _, customerRepo, _ := setupNotificationService()
	ctx := context.Background()

	customerRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	req := SendNotificationRequest{
		CustomerID: 999,
		Type:       "general",
		Title:      "Test",
		Message:    "Test",
		Channel:    "sms",
	}

	result, err := service.SendToCustomer(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrNotificationCustomerNotFound, err)
	customerRepo.AssertExpectations(t)
}

// Stats Tests

func TestNotificationService_GetStatsByCustomer_Success(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	stats := &repository.NotificationStats{
		TotalSent:      10,
		TotalDelivered: 8,
		TotalFailed:    2,
		TotalPending:   0,
	}

	notificationRepo.On("GetStatsByCustomer", ctx, int64(1)).Return(stats, nil)

	result, err := service.GetStatsByCustomer(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(10), result.TotalSent)
	assert.Equal(t, int64(8), result.TotalDelivered)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_GetStatsByBranch_Success(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	dateFrom := time.Now().AddDate(0, -1, 0)
	dateTo := time.Now()

	stats := &repository.NotificationStats{
		TotalSent:      100,
		TotalDelivered: 90,
		TotalFailed:    10,
		TotalPending:   5,
	}

	notificationRepo.On("GetStatsByBranch", ctx, int64(1), dateFrom, dateTo).Return(stats, nil)

	result, err := service.GetStatsByBranch(ctx, 1, dateFrom, dateTo)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(100), result.TotalSent)
	notificationRepo.AssertExpectations(t)
}

// --- Missing coverage tests ---

func TestNotificationService_List_Success(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notifications := []*domain.Notification{
		{ID: 1, Status: domain.NotificationStatusPending},
		{ID: 2, Status: domain.NotificationStatusSent},
	}

	filter := repository.NotificationFilter{}
	notificationRepo.On("List", ctx, filter).Return(notifications, int64(2), nil)

	result, total, err := service.List(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_ListByCustomer_Success(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notifications := []*domain.Notification{
		{ID: 1, CustomerID: 1, Status: domain.NotificationStatusPending},
	}

	filter := repository.NotificationFilter{}
	notificationRepo.On("ListByCustomer", ctx, int64(1), filter).Return(notifications, int64(1), nil)

	result, total, err := service.ListByCustomer(ctx, 1, filter)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_GetScheduledNotifications_Success(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	before := time.Now()
	notifications := []*domain.Notification{
		{ID: 1, Status: domain.NotificationStatusPending, ScheduledFor: &before},
	}

	notificationRepo.On("ListScheduled", ctx, before, 10).Return(notifications, nil)

	result, err := service.GetScheduledNotifications(ctx, before, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_GetFailedNotifications_Success(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notifications := []*domain.Notification{
		{ID: 1, Status: domain.NotificationStatusFailed, RetryCount: 1},
	}

	notificationRepo.On("ListFailed", ctx, 3, 10).Return(notifications, nil)

	result, err := service.GetFailedNotifications(ctx, 10)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_UpdateCustomerPreferences_Success(t *testing.T) {
	service, _, _, preferenceRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	prefs := []*domain.CustomerNotificationPreference{
		{CustomerID: 1, NotificationType: "loan_due_reminder", Channel: "sms", IsEnabled: true},
		{CustomerID: 1, NotificationType: "loan_due_reminder", Channel: "email", IsEnabled: false},
	}

	preferenceRepo.On("BulkUpsert", ctx, int64(1), prefs).Return(nil)

	err := service.UpdateCustomerPreferences(ctx, 1, prefs)

	assert.NoError(t, err)
	preferenceRepo.AssertExpectations(t)
}

func TestNotificationService_GetInternalNotificationByID_Success(t *testing.T) {
	service, _, _, _, internalRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	notification := &domain.InternalNotification{
		ID:     1,
		UserID: 100,
		Title:  "Test",
	}

	internalRepo.On("GetByID", ctx, int64(1)).Return(notification, nil)

	result, err := service.GetInternalNotificationByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test", result.Title)
	internalRepo.AssertExpectations(t)
}

func TestNotificationService_GetInternalNotificationByID_NotFound(t *testing.T) {
	service, _, _, _, internalRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	internalRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	result, err := service.GetInternalNotificationByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrNotificationNotFound, err)
	internalRepo.AssertExpectations(t)
}

func TestNotificationService_ListInternalNotificationsByUser_Success(t *testing.T) {
	service, _, _, _, internalRepo, _, _ := setupNotificationService()
	ctx := context.Background()

	notifications := []*domain.InternalNotification{
		{ID: 1, UserID: 100, Title: "Notification 1"},
		{ID: 2, UserID: 100, Title: "Notification 2"},
	}

	filter := repository.InternalNotificationFilter{}
	internalRepo.On("ListByUser", ctx, int64(100), filter).Return(notifications, int64(2), nil)

	result, total, err := service.ListInternalNotificationsByUser(ctx, 100, filter)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	internalRepo.AssertExpectations(t)
}

func TestNotificationService_NotifyBranchUsers_Success(t *testing.T) {
	service, _, _, _, internalRepo, _, userRepo := setupNotificationService()
	ctx := context.Background()

	branchID := int64(1)
	users := &repository.PaginatedResult[domain.User]{
		Data: []domain.User{
			{ID: 1, FirstName: "John"},
			{ID: 2, FirstName: "Jane"},
		},
		Total: 2,
	}

	userRepo.On("List", ctx, repository.UserListParams{BranchID: &branchID}).Return(users, nil)
	internalRepo.On("CreateBulk", ctx, mock.AnythingOfType("[]*domain.InternalNotification")).Return(nil)

	err := service.NotifyBranchUsers(ctx, 1, "Alert", "Test message", "warning")

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	internalRepo.AssertExpectations(t)
}

func TestNotificationService_NotifyBranchUsers_UserListError(t *testing.T) {
	service, _, _, _, _, _, userRepo := setupNotificationService()
	ctx := context.Background()

	branchID := int64(1)
	userRepo.On("List", ctx, repository.UserListParams{BranchID: &branchID}).Return(nil, errors.New("db error"))

	err := service.NotifyBranchUsers(ctx, 1, "Alert", "Test", "info")

	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestNotificationService_CreateFromTemplate_TemplateNotFound(t *testing.T) {
	service, _, templateRepo, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	templateRepo.On("GetByTypeAndChannel", ctx, "loan_due_reminder", "sms").Return(nil, nil)

	req := CreateNotificationFromTemplateRequest{
		CustomerID:       1,
		NotificationType: "loan_due_reminder",
		Channel:          "sms",
	}

	result, err := service.CreateFromTemplate(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTemplateNotFound, err)
	templateRepo.AssertExpectations(t)
}

func TestNotificationService_CreateFromTemplate_ChannelDisabled(t *testing.T) {
	service, _, templateRepo, preferenceRepo, _, _, _ := setupNotificationService()
	ctx := context.Background()

	tmpl := &domain.NotificationTemplate{
		ID:               1,
		NotificationType: "loan_due_reminder",
		Channel:          "sms",
		BodyTemplate:     "Hello {{name}}",
		IsActive:         true,
	}

	templateRepo.On("GetByTypeAndChannel", ctx, "loan_due_reminder", "sms").Return(tmpl, nil)
	preferenceRepo.On("IsEnabled", ctx, int64(1), "loan_due_reminder", "sms").Return(false, nil)

	req := CreateNotificationFromTemplateRequest{
		CustomerID:       1,
		NotificationType: "loan_due_reminder",
		Channel:          "sms",
	}

	result, err := service.CreateFromTemplate(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "disabled")
	templateRepo.AssertExpectations(t)
	preferenceRepo.AssertExpectations(t)
}

func TestNotificationService_UpdateTemplate_NotFound(t *testing.T) {
	service, _, templateRepo, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	templateRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	req := UpdateNotificationTemplateRequest{Name: "New"}
	result, err := service.UpdateTemplate(ctx, 999, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTemplateNotFound, err)
	templateRepo.AssertExpectations(t)
}

func TestNotificationService_Cancel_NotFound(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notificationRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	err := service.Cancel(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, ErrNotificationNotFound, err)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_RetryNotification_NotFound(t *testing.T) {
	service, notificationRepo, _, _, _, _, _ := setupNotificationService()
	ctx := context.Background()

	notificationRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	err := service.RetryNotification(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, ErrNotificationNotFound, err)
	notificationRepo.AssertExpectations(t)
}

func TestNotificationService_SendToCustomer_ChannelDisabled(t *testing.T) {
	service, _, _, preferenceRepo, _, customerRepo, _ := setupNotificationService()
	ctx := context.Background()

	customer := &domain.Customer{ID: 1, BranchID: 1}
	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	preferenceRepo.On("IsEnabled", ctx, int64(1), "general", "sms").Return(false, nil)

	req := SendNotificationRequest{
		CustomerID: 1,
		Type:       "general",
		Title:      "Test",
		Message:    "Test",
		Channel:    "sms",
	}

	result, err := service.SendToCustomer(ctx, req)

	assert.NoError(t, err)
	assert.Nil(t, result) // Returns nil when disabled
	customerRepo.AssertExpectations(t)
	preferenceRepo.AssertExpectations(t)
}

func TestNotificationService_SendToCustomer_WithReferenceType(t *testing.T) {
	service, notificationRepo, _, preferenceRepo, _, customerRepo, _ := setupNotificationService()
	ctx := context.Background()

	customer := &domain.Customer{ID: 1, BranchID: 1}
	refType := "loan"
	refID := int64(100)

	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)
	preferenceRepo.On("IsEnabled", ctx, int64(1), "payment", "email").Return(true, nil)
	notificationRepo.On("Create", ctx, mock.AnythingOfType("*domain.Notification")).Return(nil)

	req := SendNotificationRequest{
		CustomerID:    1,
		Type:          "payment",
		Title:         "Payment Received",
		Message:       "Your payment has been received.",
		Channel:       "email",
		ReferenceType: &refType,
		ReferenceID:   &refID,
	}

	result, err := service.SendToCustomer(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "loan", result.ReferenceType)
	customerRepo.AssertExpectations(t)
	preferenceRepo.AssertExpectations(t)
	notificationRepo.AssertExpectations(t)
}
