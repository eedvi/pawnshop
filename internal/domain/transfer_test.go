package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemTransfer_IsPending(t *testing.T) {
	tr := &ItemTransfer{Status: TransferStatusPending}
	assert.True(t, tr.IsPending())
	tr.Status = TransferStatusInTransit
	assert.False(t, tr.IsPending())
}

func TestItemTransfer_IsInTransit(t *testing.T) {
	tr := &ItemTransfer{Status: TransferStatusInTransit}
	assert.True(t, tr.IsInTransit())
	tr.Status = TransferStatusPending
	assert.False(t, tr.IsInTransit())
}

func TestItemTransfer_IsCompleted(t *testing.T) {
	tr := &ItemTransfer{Status: TransferStatusCompleted}
	assert.True(t, tr.IsCompleted())
	tr.Status = TransferStatusPending
	assert.False(t, tr.IsCompleted())
}

func TestItemTransfer_IsCancelled(t *testing.T) {
	tr := &ItemTransfer{Status: TransferStatusCancelled}
	assert.True(t, tr.IsCancelled())
	tr.Status = TransferStatusPending
	assert.False(t, tr.IsCancelled())
}

func TestItemTransfer_CanApprove_Pending(t *testing.T) {
	tr := &ItemTransfer{Status: TransferStatusPending}
	assert.True(t, tr.CanApprove())
}

func TestItemTransfer_CanApprove_NotPending(t *testing.T) {
	statuses := []string{TransferStatusInTransit, TransferStatusCompleted, TransferStatusCancelled}
	for _, s := range statuses {
		tr := &ItemTransfer{Status: s}
		assert.False(t, tr.CanApprove(), "expected false for status %s", s)
	}
}

func TestItemTransfer_CanShip_Pending(t *testing.T) {
	tr := &ItemTransfer{Status: TransferStatusPending}
	assert.True(t, tr.CanShip())
}

func TestItemTransfer_CanShip_NotPending(t *testing.T) {
	tr := &ItemTransfer{Status: TransferStatusInTransit}
	assert.False(t, tr.CanShip())
}

func TestItemTransfer_CanReceive_InTransit(t *testing.T) {
	tr := &ItemTransfer{Status: TransferStatusInTransit}
	assert.True(t, tr.CanReceive())
}

func TestItemTransfer_CanReceive_NotInTransit(t *testing.T) {
	statuses := []string{TransferStatusPending, TransferStatusCompleted, TransferStatusCancelled}
	for _, s := range statuses {
		tr := &ItemTransfer{Status: s}
		assert.False(t, tr.CanReceive(), "expected false for status %s", s)
	}
}

func TestItemTransfer_CanCancel_Pending(t *testing.T) {
	tr := &ItemTransfer{Status: TransferStatusPending}
	assert.True(t, tr.CanCancel())
}

func TestItemTransfer_CanCancel_InTransit(t *testing.T) {
	tr := &ItemTransfer{Status: TransferStatusInTransit}
	assert.True(t, tr.CanCancel())
}

func TestItemTransfer_CanCancel_CompletedOrCancelled(t *testing.T) {
	tr := &ItemTransfer{Status: TransferStatusCompleted}
	assert.False(t, tr.CanCancel())
	tr.Status = TransferStatusCancelled
	assert.False(t, tr.CanCancel())
}
