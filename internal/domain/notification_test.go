package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNotification_IsPending(t *testing.T) {
	n := &Notification{Status: NotificationStatusPending}
	assert.True(t, n.IsPending())
	n.Status = NotificationStatusSent
	assert.False(t, n.IsPending())
}

func TestNotification_IsSent(t *testing.T) {
	n := &Notification{Status: NotificationStatusSent}
	assert.True(t, n.IsSent())
	n.Status = NotificationStatusPending
	assert.False(t, n.IsSent())
}

func TestNotification_IsDelivered(t *testing.T) {
	n := &Notification{Status: NotificationStatusDelivered}
	assert.True(t, n.IsDelivered())
	n.Status = NotificationStatusSent
	assert.False(t, n.IsDelivered())
}

func TestNotification_IsFailed(t *testing.T) {
	n := &Notification{Status: NotificationStatusFailed}
	assert.True(t, n.IsFailed())
	n.Status = NotificationStatusSent
	assert.False(t, n.IsFailed())
}

func TestNotification_CanRetry_FailedLowRetry(t *testing.T) {
	n := &Notification{Status: NotificationStatusFailed, RetryCount: 0}
	assert.True(t, n.CanRetry())

	n.RetryCount = 2
	assert.True(t, n.CanRetry())
}

func TestNotification_CanRetry_FailedMaxRetry(t *testing.T) {
	n := &Notification{Status: NotificationStatusFailed, RetryCount: 3}
	assert.False(t, n.CanRetry())

	n.RetryCount = 5
	assert.False(t, n.CanRetry())
}

func TestNotification_CanRetry_NotFailed(t *testing.T) {
	n := &Notification{Status: NotificationStatusSent, RetryCount: 0}
	assert.False(t, n.CanRetry())

	n.Status = NotificationStatusPending
	assert.False(t, n.CanRetry())
}

func TestInternalNotification_MarkAsRead(t *testing.T) {
	n := &InternalNotification{
		IsRead: false,
		ReadAt: nil,
	}

	n.MarkAsRead()

	assert.True(t, n.IsRead)
	assert.NotNil(t, n.ReadAt)
	assert.WithinDuration(t, time.Now(), *n.ReadAt, 2*time.Second)
}

func TestInternalNotification_MarkAsRead_AlreadyRead(t *testing.T) {
	earlier := time.Now().Add(-1 * time.Hour)
	n := &InternalNotification{
		IsRead: true,
		ReadAt: &earlier,
	}

	n.MarkAsRead()

	assert.True(t, n.IsRead)
	assert.NotNil(t, n.ReadAt)
	// ReadAt should be updated to now
	assert.WithinDuration(t, time.Now(), *n.ReadAt, 2*time.Second)
}
