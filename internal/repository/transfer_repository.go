package repository

import (
	"context"

	"pawnshop/internal/domain"
)

// TransferRepository defines the interface for item transfer operations
type TransferRepository interface {
	// Create creates a new item transfer
	Create(ctx context.Context, transfer *domain.ItemTransfer) error

	// GetByID retrieves a transfer by its ID
	GetByID(ctx context.Context, id int64) (*domain.ItemTransfer, error)

	// GetByNumber retrieves a transfer by its number
	GetByNumber(ctx context.Context, number string) (*domain.ItemTransfer, error)

	// Update updates an existing transfer
	Update(ctx context.Context, transfer *domain.ItemTransfer) error

	// List retrieves transfers with filtering and pagination
	List(ctx context.Context, filter TransferFilter) ([]*domain.ItemTransfer, int64, error)

	// ListByBranch retrieves transfers for a specific branch (from or to)
	ListByBranch(ctx context.Context, branchID int64, filter TransferFilter) ([]*domain.ItemTransfer, int64, error)

	// ListByItem retrieves transfers for a specific item
	ListByItem(ctx context.Context, itemID int64) ([]*domain.ItemTransfer, error)

	// GetPendingForBranch retrieves pending transfers for a branch
	GetPendingForBranch(ctx context.Context, branchID int64) ([]*domain.ItemTransfer, error)

	// GetInTransitForBranch retrieves in-transit transfers for a branch
	GetInTransitForBranch(ctx context.Context, branchID int64) ([]*domain.ItemTransfer, error)

	// GenerateTransferNumber generates a unique transfer number
	GenerateTransferNumber(ctx context.Context) (string, error)
}

// TransferFilter contains filters for listing transfers
type TransferFilter struct {
	FromBranchID *int64
	ToBranchID   *int64
	Status       *string
	ItemID       *int64
	RequestedBy  *int64
	DateFrom     *string
	DateTo       *string
	Page         int
	PageSize     int
}
