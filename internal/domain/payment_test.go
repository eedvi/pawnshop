package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPayment_TableName(t *testing.T) {
	assert.Equal(t, "payments", Payment{}.TableName())
}

func TestPayment_IsReversed_True(t *testing.T) {
	p := &Payment{Status: PaymentStatusReversed}
	assert.True(t, p.IsReversed())
}

func TestPayment_IsReversed_False(t *testing.T) {
	statuses := []PaymentStatus{
		PaymentStatusCompleted, PaymentStatusPending, PaymentStatusFailed,
	}
	for _, status := range statuses {
		p := &Payment{Status: status}
		assert.False(t, p.IsReversed(), "expected false for status %s", status)
	}
}

func TestPayment_CanBeReversed_Completed(t *testing.T) {
	p := &Payment{Status: PaymentStatusCompleted}
	assert.True(t, p.CanBeReversed())
}

func TestPayment_CanBeReversed_OtherStatuses(t *testing.T) {
	statuses := []PaymentStatus{
		PaymentStatusPending, PaymentStatusReversed, PaymentStatusFailed,
	}
	for _, status := range statuses {
		p := &Payment{Status: status}
		assert.False(t, p.CanBeReversed(), "expected false for status %s", status)
	}
}
