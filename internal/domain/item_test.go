package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItem_TableName(t *testing.T) {
	assert.Equal(t, "items", Item{}.TableName())
}

func TestItem_IsAvailable_True(t *testing.T) {
	i := &Item{Status: ItemStatusAvailable}
	assert.True(t, i.IsAvailable())
}

func TestItem_IsAvailable_False(t *testing.T) {
	statuses := []ItemStatus{
		ItemStatusPawned, ItemStatusCollateral, ItemStatusForSale,
		ItemStatusSold, ItemStatusConfiscated, ItemStatusTransferred,
		ItemStatusInTransfer, ItemStatusDamaged, ItemStatusLost,
	}
	for _, status := range statuses {
		i := &Item{Status: status}
		assert.False(t, i.IsAvailable(), "expected false for status %s", status)
	}
}

func TestItem_CanBeSold_Available(t *testing.T) {
	i := &Item{Status: ItemStatusAvailable}
	assert.True(t, i.CanBeSold())
}

func TestItem_CanBeSold_Confiscated(t *testing.T) {
	i := &Item{Status: ItemStatusConfiscated}
	assert.True(t, i.CanBeSold())
}

func TestItem_CanBeSold_OtherStatuses(t *testing.T) {
	statuses := []ItemStatus{
		ItemStatusPawned, ItemStatusCollateral, ItemStatusForSale,
		ItemStatusSold, ItemStatusTransferred, ItemStatusInTransfer,
		ItemStatusDamaged, ItemStatusLost,
	}
	for _, status := range statuses {
		i := &Item{Status: status}
		assert.False(t, i.CanBeSold(), "expected false for status %s", status)
	}
}

func TestItemHistory_TableName(t *testing.T) {
	assert.Equal(t, "item_history", ItemHistory{}.TableName())
}
