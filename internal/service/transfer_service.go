package service

import (
	"context"
	"errors"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

var (
	ErrTransferAlreadyExists  = errors.New("transfer already exists")
	ErrTransferCannotApprove  = errors.New("transfer cannot be approved")
	ErrTransferCannotShip     = errors.New("transfer cannot be shipped")
	ErrTransferCannotReceive  = errors.New("transfer cannot be received")
	ErrTransferCannotCancel   = errors.New("transfer cannot be cancelled")
	ErrItemInTransfer         = errors.New("item is already in transfer")
	ErrSameBranch             = errors.New("cannot transfer to the same branch")
)

// TransferService defines the interface for transfer operations
type TransferService interface {
	// Create creates a new item transfer request
	Create(ctx context.Context, req CreateTransferRequest) (*domain.ItemTransfer, error)

	// GetByID retrieves a transfer by ID
	GetByID(ctx context.Context, id int64) (*domain.ItemTransfer, error)

	// GetByNumber retrieves a transfer by number
	GetByNumber(ctx context.Context, number string) (*domain.ItemTransfer, error)

	// List retrieves transfers with filtering
	List(ctx context.Context, filter repository.TransferFilter) ([]*domain.ItemTransfer, int64, error)

	// ListByBranch retrieves transfers for a branch
	ListByBranch(ctx context.Context, branchID int64, filter repository.TransferFilter) ([]*domain.ItemTransfer, int64, error)

	// Approve approves a pending transfer
	Approve(ctx context.Context, id int64, approvedBy int64, notes string) (*domain.ItemTransfer, error)

	// Ship marks a transfer as shipped/in transit
	Ship(ctx context.Context, id int64, shippedBy int64) (*domain.ItemTransfer, error)

	// Receive marks a transfer as received/completed
	Receive(ctx context.Context, id int64, receivedBy int64, notes string) (*domain.ItemTransfer, error)

	// Cancel cancels a transfer
	Cancel(ctx context.Context, id int64, cancelledBy int64, reason string) (*domain.ItemTransfer, error)

	// GetPendingForBranch retrieves pending transfers for a branch
	GetPendingForBranch(ctx context.Context, branchID int64) ([]*domain.ItemTransfer, error)

	// GetInTransitForBranch retrieves in-transit transfers for a branch
	GetInTransitForBranch(ctx context.Context, branchID int64) ([]*domain.ItemTransfer, error)
}

type transferService struct {
	transferRepo repository.TransferRepository
	itemRepo     repository.ItemRepository
	branchRepo   repository.BranchRepository
}

// NewTransferService creates a new transfer service
func NewTransferService(
	transferRepo repository.TransferRepository,
	itemRepo repository.ItemRepository,
	branchRepo repository.BranchRepository,
) TransferService {
	return &transferService{
		transferRepo: transferRepo,
		itemRepo:     itemRepo,
		branchRepo:   branchRepo,
	}
}

// CreateTransferRequest represents a request to create a transfer
type CreateTransferRequest struct {
	ItemID       int64  `json:"item_id" validate:"required"`
	FromBranchID int64  `json:"from_branch_id" validate:"required"`
	ToBranchID   int64  `json:"to_branch_id" validate:"required"`
	RequestedBy  int64  `json:"requested_by" validate:"required"`
	Notes        string `json:"notes"`
}

func (s *transferService) Create(ctx context.Context, req CreateTransferRequest) (*domain.ItemTransfer, error) {
	// Validate branches are different
	if req.FromBranchID == req.ToBranchID {
		return nil, ErrSameBranch
	}

	// Check if item exists
	item, err := s.itemRepo.GetByID(ctx, req.ItemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}

	// Check if item is available for transfer (not in loan, not already in transfer)
	if item.Status != domain.ItemStatusAvailable {
		return nil, ErrItemInTransfer
	}

	// Check if item belongs to the source branch
	if item.BranchID != req.FromBranchID {
		return nil, errors.New("item does not belong to the source branch")
	}

	// Validate branches exist
	fromBranch, err := s.branchRepo.GetByID(ctx, req.FromBranchID)
	if err != nil || fromBranch == nil {
		return nil, ErrBranchNotFound
	}

	toBranch, err := s.branchRepo.GetByID(ctx, req.ToBranchID)
	if err != nil || toBranch == nil {
		return nil, ErrBranchNotFound
	}

	// Generate transfer number
	transferNumber, err := s.transferRepo.GenerateTransferNumber(ctx)
	if err != nil {
		return nil, err
	}

	// Create transfer
	transfer := &domain.ItemTransfer{
		TransferNumber: transferNumber,
		ItemID:         req.ItemID,
		FromBranchID:   req.FromBranchID,
		ToBranchID:     req.ToBranchID,
		Status:         domain.TransferStatusPending,
		RequestedBy:    req.RequestedBy,
		RequestedAt:    time.Now(),
		RequestNotes:   req.Notes,
	}

	if err := s.transferRepo.Create(ctx, transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *transferService) GetByID(ctx context.Context, id int64) (*domain.ItemTransfer, error) {
	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, ErrTransferNotFound
	}

	// Load related entities
	if err := s.loadTransferRelations(ctx, transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *transferService) GetByNumber(ctx context.Context, number string) (*domain.ItemTransfer, error) {
	transfer, err := s.transferRepo.GetByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, ErrTransferNotFound
	}

	if err := s.loadTransferRelations(ctx, transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *transferService) loadTransferRelations(ctx context.Context, transfer *domain.ItemTransfer) error {
	// Load item
	item, err := s.itemRepo.GetByID(ctx, transfer.ItemID)
	if err != nil {
		return err
	}
	transfer.Item = item

	// Load branches
	fromBranch, err := s.branchRepo.GetByID(ctx, transfer.FromBranchID)
	if err != nil {
		return err
	}
	transfer.FromBranch = fromBranch

	toBranch, err := s.branchRepo.GetByID(ctx, transfer.ToBranchID)
	if err != nil {
		return err
	}
	transfer.ToBranch = toBranch

	return nil
}

func (s *transferService) List(ctx context.Context, filter repository.TransferFilter) ([]*domain.ItemTransfer, int64, error) {
	return s.transferRepo.List(ctx, filter)
}

func (s *transferService) ListByBranch(ctx context.Context, branchID int64, filter repository.TransferFilter) ([]*domain.ItemTransfer, int64, error) {
	return s.transferRepo.ListByBranch(ctx, branchID, filter)
}

func (s *transferService) Approve(ctx context.Context, id int64, approvedBy int64, notes string) (*domain.ItemTransfer, error) {
	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, ErrTransferNotFound
	}

	if !transfer.CanApprove() {
		return nil, ErrTransferCannotApprove
	}

	now := time.Now()
	transfer.ApprovedBy = &approvedBy
	transfer.ApprovedAt = &now
	transfer.ApprovalNotes = notes

	if err := s.transferRepo.Update(ctx, transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *transferService) Ship(ctx context.Context, id int64, shippedBy int64) (*domain.ItemTransfer, error) {
	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, ErrTransferNotFound
	}

	if !transfer.CanShip() {
		return nil, ErrTransferCannotShip
	}

	// Mark item as in transfer
	item, err := s.itemRepo.GetByID(ctx, transfer.ItemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}

	// Update item status
	item.Status = domain.ItemStatusInTransfer
	if err := s.itemRepo.Update(ctx, item); err != nil {
		return nil, err
	}

	// Update transfer
	now := time.Now()
	transfer.Status = domain.TransferStatusInTransit
	transfer.ShippedAt = &now

	if err := s.transferRepo.Update(ctx, transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *transferService) Receive(ctx context.Context, id int64, receivedBy int64, notes string) (*domain.ItemTransfer, error) {
	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, ErrTransferNotFound
	}

	if !transfer.CanReceive() {
		return nil, ErrTransferCannotReceive
	}

	// Update item's branch
	item, err := s.itemRepo.GetByID(ctx, transfer.ItemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}

	// Update item
	item.BranchID = transfer.ToBranchID
	item.Status = domain.ItemStatusAvailable
	if err := s.itemRepo.Update(ctx, item); err != nil {
		return nil, err
	}

	// Update transfer
	now := time.Now()
	transfer.Status = domain.TransferStatusCompleted
	transfer.ReceivedBy = &receivedBy
	transfer.ReceivedAt = &now
	transfer.ReceiptNotes = notes

	if err := s.transferRepo.Update(ctx, transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *transferService) Cancel(ctx context.Context, id int64, cancelledBy int64, reason string) (*domain.ItemTransfer, error) {
	transfer, err := s.transferRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, ErrTransferNotFound
	}

	if !transfer.CanCancel() {
		return nil, ErrTransferCannotCancel
	}

	// If in transit, revert item status
	if transfer.IsInTransit() {
		item, err := s.itemRepo.GetByID(ctx, transfer.ItemID)
		if err != nil {
			return nil, err
		}
		if item != nil {
			item.Status = domain.ItemStatusAvailable
			if err := s.itemRepo.Update(ctx, item); err != nil {
				return nil, err
			}
		}
	}

	// Update transfer
	now := time.Now()
	transfer.Status = domain.TransferStatusCancelled
	transfer.CancelledAt = &now
	transfer.CancellationReason = reason

	if err := s.transferRepo.Update(ctx, transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *transferService) GetPendingForBranch(ctx context.Context, branchID int64) ([]*domain.ItemTransfer, error) {
	return s.transferRepo.GetPendingForBranch(ctx, branchID)
}

func (s *transferService) GetInTransitForBranch(ctx context.Context, branchID int64) ([]*domain.ItemTransfer, error) {
	return s.transferRepo.GetInTransitForBranch(ctx, branchID)
}
