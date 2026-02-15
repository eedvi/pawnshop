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

func setupTransferService() (TransferService, *mocks.MockTransferRepository, *mocks.MockItemRepository, *mocks.MockBranchRepository) {
	transferRepo := new(mocks.MockTransferRepository)
	itemRepo := new(mocks.MockItemRepository)
	branchRepo := new(mocks.MockBranchRepository)
	service := NewTransferService(transferRepo, itemRepo, branchRepo)
	return service, transferRepo, itemRepo, branchRepo
}

// --- Create tests ---

func TestTransferService_Create_Success(t *testing.T) {
	service, transferRepo, itemRepo, branchRepo := setupTransferService()
	ctx := context.Background()

	item := &domain.Item{
		ID:       1,
		BranchID: 1,
		Name:     "Test Item",
		Status:   domain.ItemStatusAvailable,
	}

	fromBranch := &domain.Branch{ID: 1, Name: "Branch A", IsActive: true}
	toBranch := &domain.Branch{ID: 2, Name: "Branch B", IsActive: true}

	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	branchRepo.On("GetByID", ctx, int64(1)).Return(fromBranch, nil)
	branchRepo.On("GetByID", ctx, int64(2)).Return(toBranch, nil)
	transferRepo.On("GenerateTransferNumber", ctx).Return("TR-001", nil)
	transferRepo.On("Create", ctx, mock.AnythingOfType("*domain.ItemTransfer")).Return(nil)

	req := CreateTransferRequest{
		ItemID:       1,
		FromBranchID: 1,
		ToBranchID:   2,
		RequestedBy:  100,
		Notes:        "Test transfer",
	}

	result, err := service.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "TR-001", result.TransferNumber)
	assert.Equal(t, domain.TransferStatusPending, result.Status)
	transferRepo.AssertExpectations(t)
	itemRepo.AssertExpectations(t)
	branchRepo.AssertExpectations(t)
}

func TestTransferService_Create_SameBranch(t *testing.T) {
	service, _, _, _ := setupTransferService()
	ctx := context.Background()

	req := CreateTransferRequest{
		ItemID:       1,
		FromBranchID: 1,
		ToBranchID:   1,
		RequestedBy:  100,
	}

	result, err := service.Create(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrSameBranch, err)
}

func TestTransferService_Create_ItemNotAvailable(t *testing.T) {
	service, _, itemRepo, _ := setupTransferService()
	ctx := context.Background()

	item := &domain.Item{
		ID:       1,
		BranchID: 1,
		Status:   domain.ItemStatusPawned,
	}

	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	req := CreateTransferRequest{
		ItemID:       1,
		FromBranchID: 1,
		ToBranchID:   2,
		RequestedBy:  100,
	}

	result, err := service.Create(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrItemInTransfer, err)
	itemRepo.AssertExpectations(t)
}

func TestTransferService_Create_ItemWrongBranch(t *testing.T) {
	service, _, itemRepo, _ := setupTransferService()
	ctx := context.Background()

	item := &domain.Item{
		ID:       1,
		BranchID: 3,
		Status:   domain.ItemStatusAvailable,
	}

	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	req := CreateTransferRequest{
		ItemID:       1,
		FromBranchID: 1,
		ToBranchID:   2,
		RequestedBy:  100,
	}

	result, err := service.Create(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "item does not belong to the source branch", err.Error())
	itemRepo.AssertExpectations(t)
}

func TestTransferService_Create_ItemNotFound(t *testing.T) {
	service, _, itemRepo, _ := setupTransferService()
	ctx := context.Background()

	itemRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	req := CreateTransferRequest{
		ItemID:       999,
		FromBranchID: 1,
		ToBranchID:   2,
		RequestedBy:  100,
	}

	result, err := service.Create(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestTransferService_Create_FromBranchNotFound(t *testing.T) {
	service, _, itemRepo, branchRepo := setupTransferService()
	ctx := context.Background()

	item := &domain.Item{
		ID:       1,
		BranchID: 1,
		Status:   domain.ItemStatusAvailable,
	}

	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	branchRepo.On("GetByID", ctx, int64(1)).Return(nil, errors.New("not found"))

	req := CreateTransferRequest{
		ItemID:       1,
		FromBranchID: 1,
		ToBranchID:   2,
		RequestedBy:  100,
	}

	result, err := service.Create(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrBranchNotFound, err)
}

func TestTransferService_Create_ToBranchNotFound(t *testing.T) {
	service, _, itemRepo, branchRepo := setupTransferService()
	ctx := context.Background()

	item := &domain.Item{
		ID:       1,
		BranchID: 1,
		Status:   domain.ItemStatusAvailable,
	}

	fromBranch := &domain.Branch{ID: 1, Name: "Branch A"}

	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	branchRepo.On("GetByID", ctx, int64(1)).Return(fromBranch, nil)
	branchRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	req := CreateTransferRequest{
		ItemID:       1,
		FromBranchID: 1,
		ToBranchID:   999,
		RequestedBy:  100,
	}

	result, err := service.Create(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrBranchNotFound, err)
}

// --- GetByID tests ---

func TestTransferService_GetByID_Success(t *testing.T) {
	service, transferRepo, itemRepo, branchRepo := setupTransferService()
	ctx := context.Background()

	transfer := &domain.ItemTransfer{
		ID:             1,
		TransferNumber: "TR-001",
		ItemID:         1,
		FromBranchID:   1,
		ToBranchID:     2,
		Status:         domain.TransferStatusPending,
	}

	item := &domain.Item{ID: 1, Name: "Test Item"}
	fromBranch := &domain.Branch{ID: 1, Name: "Branch A"}
	toBranch := &domain.Branch{ID: 2, Name: "Branch B"}

	transferRepo.On("GetByID", ctx, int64(1)).Return(transfer, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	branchRepo.On("GetByID", ctx, int64(1)).Return(fromBranch, nil)
	branchRepo.On("GetByID", ctx, int64(2)).Return(toBranch, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "TR-001", result.TransferNumber)
	assert.NotNil(t, result.Item)
	assert.NotNil(t, result.FromBranch)
	assert.NotNil(t, result.ToBranch)
	transferRepo.AssertExpectations(t)
}

func TestTransferService_GetByID_NotFound(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transferRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	transferRepo.AssertExpectations(t)
}

func TestTransferService_GetByID_NilTransfer(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transferRepo.On("GetByID", ctx, int64(1)).Return(nil, nil)

	result, err := service.GetByID(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTransferNotFound, err)
}

// --- GetByNumber tests ---

func TestTransferService_GetByNumber_Success(t *testing.T) {
	service, transferRepo, itemRepo, branchRepo := setupTransferService()
	ctx := context.Background()

	transfer := &domain.ItemTransfer{
		ID:             1,
		TransferNumber: "TR-001",
		ItemID:         1,
		FromBranchID:   1,
		ToBranchID:     2,
		Status:         domain.TransferStatusPending,
	}

	item := &domain.Item{ID: 1, Name: "Test Item"}
	fromBranch := &domain.Branch{ID: 1, Name: "Branch A"}
	toBranch := &domain.Branch{ID: 2, Name: "Branch B"}

	transferRepo.On("GetByNumber", ctx, "TR-001").Return(transfer, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	branchRepo.On("GetByID", ctx, int64(1)).Return(fromBranch, nil)
	branchRepo.On("GetByID", ctx, int64(2)).Return(toBranch, nil)

	result, err := service.GetByNumber(ctx, "TR-001")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "TR-001", result.TransferNumber)
}

func TestTransferService_GetByNumber_NotFound(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transferRepo.On("GetByNumber", ctx, "NONEXISTENT").Return(nil, nil)

	result, err := service.GetByNumber(ctx, "NONEXISTENT")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTransferNotFound, err)
}

// --- List tests ---

func TestTransferService_List_Success(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transfers := []*domain.ItemTransfer{
		{ID: 1, TransferNumber: "TR-001"},
		{ID: 2, TransferNumber: "TR-002"},
	}

	transferRepo.On("List", ctx, mock.AnythingOfType("repository.TransferFilter")).Return(transfers, int64(2), nil)

	filter := repository.TransferFilter{}
	result, total, err := service.List(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	transferRepo.AssertExpectations(t)
}

// --- ListByBranch tests ---

func TestTransferService_ListByBranch_Success(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transfers := []*domain.ItemTransfer{
		{ID: 1, TransferNumber: "TR-001", FromBranchID: 1},
	}

	transferRepo.On("ListByBranch", ctx, int64(1), mock.AnythingOfType("repository.TransferFilter")).Return(transfers, int64(1), nil)

	filter := repository.TransferFilter{}
	result, total, err := service.ListByBranch(ctx, 1, filter)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	transferRepo.AssertExpectations(t)
}

// --- Approve tests ---

func TestTransferService_Approve_Success(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transfer := &domain.ItemTransfer{
		ID:     1,
		Status: domain.TransferStatusPending,
	}

	transferRepo.On("GetByID", ctx, int64(1)).Return(transfer, nil)
	transferRepo.On("Update", ctx, mock.AnythingOfType("*domain.ItemTransfer")).Return(nil)

	result, err := service.Approve(ctx, 1, 100, "Approved for transfer")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.ApprovedBy)
	assert.Equal(t, int64(100), *result.ApprovedBy)
	assert.NotNil(t, result.ApprovedAt)
	assert.Equal(t, "Approved for transfer", result.ApprovalNotes)
	transferRepo.AssertExpectations(t)
}

func TestTransferService_Approve_NotFound(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transferRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	result, err := service.Approve(ctx, 999, 100, "notes")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTransferNotFound, err)
}

func TestTransferService_Approve_InvalidStatus(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transfer := &domain.ItemTransfer{
		ID:     1,
		Status: domain.TransferStatusCompleted,
	}

	transferRepo.On("GetByID", ctx, int64(1)).Return(transfer, nil)

	result, err := service.Approve(ctx, 1, 100, "notes")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTransferCannotApprove, err)
}

// --- Ship tests ---

func TestTransferService_Ship_Success(t *testing.T) {
	service, transferRepo, itemRepo, _ := setupTransferService()
	ctx := context.Background()

	transfer := &domain.ItemTransfer{
		ID:           1,
		ItemID:       1,
		FromBranchID: 1,
		ToBranchID:   2,
		Status:       domain.TransferStatusPending,
	}

	item := &domain.Item{
		ID:       1,
		BranchID: 1,
		Status:   domain.ItemStatusAvailable,
	}

	transferRepo.On("GetByID", ctx, int64(1)).Return(transfer, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	itemRepo.On("Update", ctx, mock.MatchedBy(func(i *domain.Item) bool {
		return i.Status == domain.ItemStatusInTransfer
	})).Return(nil)
	transferRepo.On("Update", ctx, mock.MatchedBy(func(tr *domain.ItemTransfer) bool {
		return tr.Status == domain.TransferStatusInTransit && tr.ShippedAt != nil
	})).Return(nil)

	result, err := service.Ship(ctx, 1, 100)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.TransferStatusInTransit, result.Status)
	transferRepo.AssertExpectations(t)
	itemRepo.AssertExpectations(t)
}

func TestTransferService_Ship_NotFound(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transferRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	result, err := service.Ship(ctx, 999, 100)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTransferNotFound, err)
}

func TestTransferService_Ship_InvalidStatus(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transfer := &domain.ItemTransfer{
		ID:     1,
		Status: domain.TransferStatusCompleted,
	}

	transferRepo.On("GetByID", ctx, int64(1)).Return(transfer, nil)

	result, err := service.Ship(ctx, 1, 100)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTransferCannotShip, err)
}

// --- Receive tests ---

func TestTransferService_Receive_Success(t *testing.T) {
	service, transferRepo, itemRepo, _ := setupTransferService()
	ctx := context.Background()

	shippedAt := time.Now().Add(-1 * time.Hour)
	transfer := &domain.ItemTransfer{
		ID:           1,
		ItemID:       1,
		FromBranchID: 1,
		ToBranchID:   2,
		Status:       domain.TransferStatusInTransit,
		ShippedAt:    &shippedAt,
	}

	item := &domain.Item{
		ID:       1,
		BranchID: 1,
		Status:   domain.ItemStatusInTransfer,
	}

	transferRepo.On("GetByID", ctx, int64(1)).Return(transfer, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	itemRepo.On("Update", ctx, mock.MatchedBy(func(i *domain.Item) bool {
		return i.BranchID == int64(2) && i.Status == domain.ItemStatusAvailable
	})).Return(nil)
	transferRepo.On("Update", ctx, mock.MatchedBy(func(tr *domain.ItemTransfer) bool {
		return tr.Status == domain.TransferStatusCompleted && tr.ReceivedAt != nil
	})).Return(nil)

	result, err := service.Receive(ctx, 1, 200, "Received in good condition")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.TransferStatusCompleted, result.Status)
	transferRepo.AssertExpectations(t)
	itemRepo.AssertExpectations(t)
}

func TestTransferService_Receive_NotFound(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transferRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	result, err := service.Receive(ctx, 999, 200, "notes")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTransferNotFound, err)
}

func TestTransferService_Receive_InvalidStatus(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transfer := &domain.ItemTransfer{
		ID:     1,
		Status: domain.TransferStatusPending,
	}

	transferRepo.On("GetByID", ctx, int64(1)).Return(transfer, nil)

	result, err := service.Receive(ctx, 1, 200, "notes")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTransferCannotReceive, err)
}

// --- Cancel tests ---

func TestTransferService_Cancel_Success(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transfer := &domain.ItemTransfer{
		ID:           1,
		ItemID:       1,
		FromBranchID: 1,
		ToBranchID:   2,
		Status:       domain.TransferStatusPending,
	}

	transferRepo.On("GetByID", ctx, int64(1)).Return(transfer, nil)
	transferRepo.On("Update", ctx, mock.MatchedBy(func(tr *domain.ItemTransfer) bool {
		return tr.Status == domain.TransferStatusCancelled && tr.CancelledAt != nil
	})).Return(nil)

	result, err := service.Cancel(ctx, 1, 100, "No longer needed")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.TransferStatusCancelled, result.Status)
	assert.Equal(t, "No longer needed", result.CancellationReason)
	transferRepo.AssertExpectations(t)
}

func TestTransferService_Cancel_InTransit_RevertsItemStatus(t *testing.T) {
	service, transferRepo, itemRepo, _ := setupTransferService()
	ctx := context.Background()

	shippedAt := time.Now().Add(-1 * time.Hour)
	transfer := &domain.ItemTransfer{
		ID:           1,
		ItemID:       1,
		FromBranchID: 1,
		ToBranchID:   2,
		Status:       domain.TransferStatusInTransit,
		ShippedAt:    &shippedAt,
	}

	item := &domain.Item{
		ID:       1,
		BranchID: 1,
		Status:   domain.ItemStatusInTransfer,
	}

	transferRepo.On("GetByID", ctx, int64(1)).Return(transfer, nil)
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	itemRepo.On("Update", ctx, mock.MatchedBy(func(i *domain.Item) bool {
		return i.Status == domain.ItemStatusAvailable
	})).Return(nil)
	transferRepo.On("Update", ctx, mock.AnythingOfType("*domain.ItemTransfer")).Return(nil)

	result, err := service.Cancel(ctx, 1, 100, "Item damaged")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.TransferStatusCancelled, result.Status)
	transferRepo.AssertExpectations(t)
	itemRepo.AssertExpectations(t)
}

func TestTransferService_Cancel_NotFound(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transferRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	result, err := service.Cancel(ctx, 999, 100, "reason")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTransferNotFound, err)
}

func TestTransferService_Cancel_InvalidStatus(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transfer := &domain.ItemTransfer{
		ID:     1,
		Status: domain.TransferStatusCompleted,
	}

	transferRepo.On("GetByID", ctx, int64(1)).Return(transfer, nil)

	result, err := service.Cancel(ctx, 1, 100, "reason")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTransferCannotCancel, err)
}

// --- GetPendingForBranch tests ---

func TestTransferService_GetPendingForBranch_Success(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transfers := []*domain.ItemTransfer{
		{ID: 1, Status: domain.TransferStatusPending},
		{ID: 2, Status: domain.TransferStatusPending},
	}

	transferRepo.On("GetPendingForBranch", ctx, int64(1)).Return(transfers, nil)

	result, err := service.GetPendingForBranch(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	transferRepo.AssertExpectations(t)
}

// --- GetInTransitForBranch tests ---

func TestTransferService_GetInTransitForBranch_Success(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transfers := []*domain.ItemTransfer{
		{ID: 1, Status: domain.TransferStatusInTransit},
	}

	transferRepo.On("GetInTransitForBranch", ctx, int64(1)).Return(transfers, nil)

	result, err := service.GetInTransitForBranch(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	transferRepo.AssertExpectations(t)
}

func TestTransferService_GetInTransitForBranch_Error(t *testing.T) {
	service, transferRepo, _, _ := setupTransferService()
	ctx := context.Background()

	transferRepo.On("GetInTransitForBranch", ctx, int64(999)).Return(nil, errors.New("db error"))

	result, err := service.GetInTransitForBranch(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}
