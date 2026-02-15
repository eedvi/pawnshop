package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSale_TableName(t *testing.T) {
	assert.Equal(t, "sales", Sale{}.TableName())
}

func TestSale_IsRefunded_True(t *testing.T) {
	s := &Sale{Status: SaleStatusRefunded}
	assert.True(t, s.IsRefunded())
}

func TestSale_IsRefunded_False(t *testing.T) {
	statuses := []SaleStatus{
		SaleStatusCompleted, SaleStatusPending, SaleStatusCancelled, SaleStatusPartialRefund,
	}
	for _, status := range statuses {
		s := &Sale{Status: status}
		assert.False(t, s.IsRefunded(), "expected false for status %s", status)
	}
}

func TestSale_CanBeRefunded_Completed(t *testing.T) {
	s := &Sale{Status: SaleStatusCompleted}
	assert.True(t, s.CanBeRefunded())
}

func TestSale_CanBeRefunded_OtherStatuses(t *testing.T) {
	statuses := []SaleStatus{
		SaleStatusPending, SaleStatusCancelled, SaleStatusRefunded, SaleStatusPartialRefund,
	}
	for _, status := range statuses {
		s := &Sale{Status: status}
		assert.False(t, s.CanBeRefunded(), "expected false for status %s", status)
	}
}
