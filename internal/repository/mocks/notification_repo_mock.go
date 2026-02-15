package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// MockNotificationTemplateRepository is a mock implementation of NotificationTemplateRepository
type MockNotificationTemplateRepository struct {
	mock.Mock
}

func (m *MockNotificationTemplateRepository) Create(ctx context.Context, template *domain.NotificationTemplate) error {
	args := m.Called(ctx, template)
	return args.Error(0)
}

func (m *MockNotificationTemplateRepository) GetByID(ctx context.Context, id int64) (*domain.NotificationTemplate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepository) GetByTypeAndChannel(ctx context.Context, notificationType, channel string) (*domain.NotificationTemplate, error) {
	args := m.Called(ctx, notificationType, channel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepository) Update(ctx context.Context, template *domain.NotificationTemplate) error {
	args := m.Called(ctx, template)
	return args.Error(0)
}

func (m *MockNotificationTemplateRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationTemplateRepository) List(ctx context.Context, includeInactive bool) ([]*domain.NotificationTemplate, error) {
	args := m.Called(ctx, includeInactive)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.NotificationTemplate), args.Error(1)
}

func (m *MockNotificationTemplateRepository) ListByType(ctx context.Context, notificationType string) ([]*domain.NotificationTemplate, error) {
	args := m.Called(ctx, notificationType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.NotificationTemplate), args.Error(1)
}

// MockNotificationRepository is a mock implementation of NotificationRepository
type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *MockNotificationRepository) GetByID(ctx context.Context, id int64) (*domain.Notification, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Notification), args.Error(1)
}

func (m *MockNotificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *MockNotificationRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationRepository) List(ctx context.Context, filter repository.NotificationFilter) ([]*domain.Notification, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Notification), args.Get(1).(int64), args.Error(2)
}

func (m *MockNotificationRepository) ListByCustomer(ctx context.Context, customerID int64, filter repository.NotificationFilter) ([]*domain.Notification, int64, error) {
	args := m.Called(ctx, customerID, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Notification), args.Get(1).(int64), args.Error(2)
}

func (m *MockNotificationRepository) ListPending(ctx context.Context, limit int) ([]*domain.Notification, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Notification), args.Error(1)
}

func (m *MockNotificationRepository) ListScheduled(ctx context.Context, before time.Time, limit int) ([]*domain.Notification, error) {
	args := m.Called(ctx, before, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Notification), args.Error(1)
}

func (m *MockNotificationRepository) ListFailed(ctx context.Context, maxRetries int, limit int) ([]*domain.Notification, error) {
	args := m.Called(ctx, maxRetries, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Notification), args.Error(1)
}

func (m *MockNotificationRepository) MarkAsSent(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationRepository) MarkAsDelivered(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationRepository) MarkAsFailed(ctx context.Context, id int64, reason string) error {
	args := m.Called(ctx, id, reason)
	return args.Error(0)
}

func (m *MockNotificationRepository) IncrementRetry(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationRepository) Cancel(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationRepository) GetStatsByCustomer(ctx context.Context, customerID int64) (*repository.NotificationStats, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.NotificationStats), args.Error(1)
}

func (m *MockNotificationRepository) GetStatsByBranch(ctx context.Context, branchID int64, dateFrom, dateTo time.Time) (*repository.NotificationStats, error) {
	args := m.Called(ctx, branchID, dateFrom, dateTo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.NotificationStats), args.Error(1)
}

// MockCustomerNotificationPreferenceRepository is a mock implementation
type MockCustomerNotificationPreferenceRepository struct {
	mock.Mock
}

func (m *MockCustomerNotificationPreferenceRepository) Create(ctx context.Context, pref *domain.CustomerNotificationPreference) error {
	args := m.Called(ctx, pref)
	return args.Error(0)
}

func (m *MockCustomerNotificationPreferenceRepository) GetByID(ctx context.Context, id int64) (*domain.CustomerNotificationPreference, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CustomerNotificationPreference), args.Error(1)
}

func (m *MockCustomerNotificationPreferenceRepository) GetByCustomerTypeAndChannel(ctx context.Context, customerID int64, notificationType, channel string) (*domain.CustomerNotificationPreference, error) {
	args := m.Called(ctx, customerID, notificationType, channel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CustomerNotificationPreference), args.Error(1)
}

func (m *MockCustomerNotificationPreferenceRepository) Update(ctx context.Context, pref *domain.CustomerNotificationPreference) error {
	args := m.Called(ctx, pref)
	return args.Error(0)
}

func (m *MockCustomerNotificationPreferenceRepository) Upsert(ctx context.Context, pref *domain.CustomerNotificationPreference) error {
	args := m.Called(ctx, pref)
	return args.Error(0)
}

func (m *MockCustomerNotificationPreferenceRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCustomerNotificationPreferenceRepository) ListByCustomer(ctx context.Context, customerID int64) ([]*domain.CustomerNotificationPreference, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.CustomerNotificationPreference), args.Error(1)
}

func (m *MockCustomerNotificationPreferenceRepository) IsEnabled(ctx context.Context, customerID int64, notificationType, channel string) (bool, error) {
	args := m.Called(ctx, customerID, notificationType, channel)
	return args.Bool(0), args.Error(1)
}

func (m *MockCustomerNotificationPreferenceRepository) BulkUpsert(ctx context.Context, customerID int64, prefs []*domain.CustomerNotificationPreference) error {
	args := m.Called(ctx, customerID, prefs)
	return args.Error(0)
}

// MockInternalNotificationRepository is a mock implementation
type MockInternalNotificationRepository struct {
	mock.Mock
}

func (m *MockInternalNotificationRepository) Create(ctx context.Context, notification *domain.InternalNotification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *MockInternalNotificationRepository) GetByID(ctx context.Context, id int64) (*domain.InternalNotification, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.InternalNotification), args.Error(1)
}

func (m *MockInternalNotificationRepository) Update(ctx context.Context, notification *domain.InternalNotification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *MockInternalNotificationRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInternalNotificationRepository) List(ctx context.Context, filter repository.InternalNotificationFilter) ([]*domain.InternalNotification, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.InternalNotification), args.Get(1).(int64), args.Error(2)
}

func (m *MockInternalNotificationRepository) ListByUser(ctx context.Context, userID int64, filter repository.InternalNotificationFilter) ([]*domain.InternalNotification, int64, error) {
	args := m.Called(ctx, userID, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.InternalNotification), args.Get(1).(int64), args.Error(2)
}

func (m *MockInternalNotificationRepository) ListUnreadByUser(ctx context.Context, userID int64, limit int) ([]*domain.InternalNotification, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.InternalNotification), args.Error(1)
}

func (m *MockInternalNotificationRepository) MarkAsRead(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInternalNotificationRepository) MarkAllAsRead(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockInternalNotificationRepository) GetUnreadCount(ctx context.Context, userID int64) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockInternalNotificationRepository) CreateBulk(ctx context.Context, notifications []*domain.InternalNotification) error {
	args := m.Called(ctx, notifications)
	return args.Error(0)
}

func (m *MockInternalNotificationRepository) DeleteOlderThan(ctx context.Context, olderThan time.Time) (int64, error) {
	args := m.Called(ctx, olderThan)
	return args.Get(0).(int64), args.Error(1)
}
