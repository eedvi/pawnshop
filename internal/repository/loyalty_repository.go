package repository

import (
	"context"

	"pawnshop/internal/domain"
)

// LoyaltyRepository defines the interface for loyalty operations
type LoyaltyRepository interface {
	// CreateHistory creates a new loyalty points history entry
	CreateHistory(ctx context.Context, history *domain.LoyaltyPointsHistory) error

	// GetHistoryByCustomer retrieves points history for a customer
	GetHistoryByCustomer(ctx context.Context, customerID int64, page, pageSize int) ([]*domain.LoyaltyPointsHistory, int64, error)

	// GetHistoryByReference retrieves history by reference
	GetHistoryByReference(ctx context.Context, referenceType string, referenceID int64) ([]*domain.LoyaltyPointsHistory, error)

	// GetTotalPointsEarned retrieves total points earned by a customer
	GetTotalPointsEarned(ctx context.Context, customerID int64) (int, error)

	// GetTotalPointsRedeemed retrieves total points redeemed by a customer
	GetTotalPointsRedeemed(ctx context.Context, customerID int64) (int, error)
}
