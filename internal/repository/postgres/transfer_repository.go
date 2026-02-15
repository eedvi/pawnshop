package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

type transferRepository struct {
	db *DB
}

// NewTransferRepository creates a new transfer repository
func NewTransferRepository(db *DB) repository.TransferRepository {
	return &transferRepository{db: db}
}

func (r *transferRepository) Create(ctx context.Context, transfer *domain.ItemTransfer) error {
	query := `
		INSERT INTO item_transfers (
			transfer_number, item_id, from_branch_id, to_branch_id,
			status, requested_by, requested_at, request_notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		transfer.TransferNumber,
		transfer.ItemID,
		transfer.FromBranchID,
		transfer.ToBranchID,
		transfer.Status,
		transfer.RequestedBy,
		transfer.RequestedAt,
		transfer.RequestNotes,
	).Scan(&transfer.ID, &transfer.CreatedAt, &transfer.UpdatedAt)
}

func (r *transferRepository) GetByID(ctx context.Context, id int64) (*domain.ItemTransfer, error) {
	query := `
		SELECT id, transfer_number, item_id, from_branch_id, to_branch_id,
			   status, requested_by, approved_by, received_by,
			   requested_at, approved_at, shipped_at, received_at, cancelled_at,
			   request_notes, approval_notes, receipt_notes, cancellation_reason,
			   created_at, updated_at
		FROM item_transfers
		WHERE id = $1`

	transfer := &domain.ItemTransfer{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transfer.ID,
		&transfer.TransferNumber,
		&transfer.ItemID,
		&transfer.FromBranchID,
		&transfer.ToBranchID,
		&transfer.Status,
		&transfer.RequestedBy,
		&transfer.ApprovedBy,
		&transfer.ReceivedBy,
		&transfer.RequestedAt,
		&transfer.ApprovedAt,
		&transfer.ShippedAt,
		&transfer.ReceivedAt,
		&transfer.CancelledAt,
		&transfer.RequestNotes,
		&transfer.ApprovalNotes,
		&transfer.ReceiptNotes,
		&transfer.CancellationReason,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return transfer, nil
}

func (r *transferRepository) GetByNumber(ctx context.Context, number string) (*domain.ItemTransfer, error) {
	query := `
		SELECT id, transfer_number, item_id, from_branch_id, to_branch_id,
			   status, requested_by, approved_by, received_by,
			   requested_at, approved_at, shipped_at, received_at, cancelled_at,
			   request_notes, approval_notes, receipt_notes, cancellation_reason,
			   created_at, updated_at
		FROM item_transfers
		WHERE transfer_number = $1`

	transfer := &domain.ItemTransfer{}
	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&transfer.ID,
		&transfer.TransferNumber,
		&transfer.ItemID,
		&transfer.FromBranchID,
		&transfer.ToBranchID,
		&transfer.Status,
		&transfer.RequestedBy,
		&transfer.ApprovedBy,
		&transfer.ReceivedBy,
		&transfer.RequestedAt,
		&transfer.ApprovedAt,
		&transfer.ShippedAt,
		&transfer.ReceivedAt,
		&transfer.CancelledAt,
		&transfer.RequestNotes,
		&transfer.ApprovalNotes,
		&transfer.ReceiptNotes,
		&transfer.CancellationReason,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return transfer, nil
}

func (r *transferRepository) Update(ctx context.Context, transfer *domain.ItemTransfer) error {
	query := `
		UPDATE item_transfers SET
			status = $2,
			approved_by = $3,
			received_by = $4,
			approved_at = $5,
			shipped_at = $6,
			received_at = $7,
			cancelled_at = $8,
			approval_notes = $9,
			receipt_notes = $10,
			cancellation_reason = $11,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	return r.db.QueryRowContext(ctx, query,
		transfer.ID,
		transfer.Status,
		transfer.ApprovedBy,
		transfer.ReceivedBy,
		transfer.ApprovedAt,
		transfer.ShippedAt,
		transfer.ReceivedAt,
		transfer.CancelledAt,
		transfer.ApprovalNotes,
		transfer.ReceiptNotes,
		transfer.CancellationReason,
	).Scan(&transfer.UpdatedAt)
}

func (r *transferRepository) List(ctx context.Context, filter repository.TransferFilter) ([]*domain.ItemTransfer, int64, error) {
	var conditions []string
	var args []interface{}
	argPos := 1

	if filter.FromBranchID != nil {
		conditions = append(conditions, fmt.Sprintf("from_branch_id = $%d", argPos))
		args = append(args, *filter.FromBranchID)
		argPos++
	}
	if filter.ToBranchID != nil {
		conditions = append(conditions, fmt.Sprintf("to_branch_id = $%d", argPos))
		args = append(args, *filter.ToBranchID)
		argPos++
	}
	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *filter.Status)
		argPos++
	}
	if filter.ItemID != nil {
		conditions = append(conditions, fmt.Sprintf("item_id = $%d", argPos))
		args = append(args, *filter.ItemID)
		argPos++
	}
	if filter.RequestedBy != nil {
		conditions = append(conditions, fmt.Sprintf("requested_by = $%d", argPos))
		args = append(args, *filter.RequestedBy)
		argPos++
	}
	if filter.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("requested_at >= $%d", argPos))
		args = append(args, *filter.DateFrom)
		argPos++
	}
	if filter.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("requested_at <= $%d", argPos))
		args = append(args, *filter.DateTo)
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM item_transfers %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT id, transfer_number, item_id, from_branch_id, to_branch_id,
			   status, requested_by, approved_by, received_by,
			   requested_at, approved_at, shipped_at, received_at, cancelled_at,
			   request_notes, approval_notes, receipt_notes, cancellation_reason,
			   created_at, updated_at
		FROM item_transfers
		%s
		ORDER BY requested_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argPos, argPos+1)

	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transfers []*domain.ItemTransfer
	for rows.Next() {
		transfer := &domain.ItemTransfer{}
		if err := rows.Scan(
			&transfer.ID,
			&transfer.TransferNumber,
			&transfer.ItemID,
			&transfer.FromBranchID,
			&transfer.ToBranchID,
			&transfer.Status,
			&transfer.RequestedBy,
			&transfer.ApprovedBy,
			&transfer.ReceivedBy,
			&transfer.RequestedAt,
			&transfer.ApprovedAt,
			&transfer.ShippedAt,
			&transfer.ReceivedAt,
			&transfer.CancelledAt,
			&transfer.RequestNotes,
			&transfer.ApprovalNotes,
			&transfer.ReceiptNotes,
			&transfer.CancellationReason,
			&transfer.CreatedAt,
			&transfer.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		transfers = append(transfers, transfer)
	}

	return transfers, total, rows.Err()
}

func (r *transferRepository) ListByBranch(ctx context.Context, branchID int64, filter repository.TransferFilter) ([]*domain.ItemTransfer, int64, error) {
	var conditions []string
	var args []interface{}
	argPos := 1

	// Either from or to this branch
	conditions = append(conditions, fmt.Sprintf("(from_branch_id = $%d OR to_branch_id = $%d)", argPos, argPos+1))
	args = append(args, branchID, branchID)
	argPos += 2

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *filter.Status)
		argPos++
	}
	if filter.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("requested_at >= $%d", argPos))
		args = append(args, *filter.DateFrom)
		argPos++
	}
	if filter.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("requested_at <= $%d", argPos))
		args = append(args, *filter.DateTo)
		argPos++
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM item_transfers %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Main query
	query := fmt.Sprintf(`
		SELECT id, transfer_number, item_id, from_branch_id, to_branch_id,
			   status, requested_by, approved_by, received_by,
			   requested_at, approved_at, shipped_at, received_at, cancelled_at,
			   request_notes, approval_notes, receipt_notes, cancellation_reason,
			   created_at, updated_at
		FROM item_transfers
		%s
		ORDER BY requested_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argPos, argPos+1)

	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transfers []*domain.ItemTransfer
	for rows.Next() {
		transfer := &domain.ItemTransfer{}
		if err := rows.Scan(
			&transfer.ID,
			&transfer.TransferNumber,
			&transfer.ItemID,
			&transfer.FromBranchID,
			&transfer.ToBranchID,
			&transfer.Status,
			&transfer.RequestedBy,
			&transfer.ApprovedBy,
			&transfer.ReceivedBy,
			&transfer.RequestedAt,
			&transfer.ApprovedAt,
			&transfer.ShippedAt,
			&transfer.ReceivedAt,
			&transfer.CancelledAt,
			&transfer.RequestNotes,
			&transfer.ApprovalNotes,
			&transfer.ReceiptNotes,
			&transfer.CancellationReason,
			&transfer.CreatedAt,
			&transfer.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		transfers = append(transfers, transfer)
	}

	return transfers, total, rows.Err()
}

func (r *transferRepository) ListByItem(ctx context.Context, itemID int64) ([]*domain.ItemTransfer, error) {
	query := `
		SELECT id, transfer_number, item_id, from_branch_id, to_branch_id,
			   status, requested_by, approved_by, received_by,
			   requested_at, approved_at, shipped_at, received_at, cancelled_at,
			   request_notes, approval_notes, receipt_notes, cancellation_reason,
			   created_at, updated_at
		FROM item_transfers
		WHERE item_id = $1
		ORDER BY requested_at DESC`

	rows, err := r.db.QueryContext(ctx, query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []*domain.ItemTransfer
	for rows.Next() {
		transfer := &domain.ItemTransfer{}
		if err := rows.Scan(
			&transfer.ID,
			&transfer.TransferNumber,
			&transfer.ItemID,
			&transfer.FromBranchID,
			&transfer.ToBranchID,
			&transfer.Status,
			&transfer.RequestedBy,
			&transfer.ApprovedBy,
			&transfer.ReceivedBy,
			&transfer.RequestedAt,
			&transfer.ApprovedAt,
			&transfer.ShippedAt,
			&transfer.ReceivedAt,
			&transfer.CancelledAt,
			&transfer.RequestNotes,
			&transfer.ApprovalNotes,
			&transfer.ReceiptNotes,
			&transfer.CancellationReason,
			&transfer.CreatedAt,
			&transfer.UpdatedAt,
		); err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}

	return transfers, rows.Err()
}

func (r *transferRepository) GetPendingForBranch(ctx context.Context, branchID int64) ([]*domain.ItemTransfer, error) {
	query := `
		SELECT id, transfer_number, item_id, from_branch_id, to_branch_id,
			   status, requested_by, approved_by, received_by,
			   requested_at, approved_at, shipped_at, received_at, cancelled_at,
			   request_notes, approval_notes, receipt_notes, cancellation_reason,
			   created_at, updated_at
		FROM item_transfers
		WHERE (from_branch_id = $1 OR to_branch_id = $1) AND status = 'pending'
		ORDER BY requested_at ASC`

	rows, err := r.db.QueryContext(ctx, query, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []*domain.ItemTransfer
	for rows.Next() {
		transfer := &domain.ItemTransfer{}
		if err := rows.Scan(
			&transfer.ID,
			&transfer.TransferNumber,
			&transfer.ItemID,
			&transfer.FromBranchID,
			&transfer.ToBranchID,
			&transfer.Status,
			&transfer.RequestedBy,
			&transfer.ApprovedBy,
			&transfer.ReceivedBy,
			&transfer.RequestedAt,
			&transfer.ApprovedAt,
			&transfer.ShippedAt,
			&transfer.ReceivedAt,
			&transfer.CancelledAt,
			&transfer.RequestNotes,
			&transfer.ApprovalNotes,
			&transfer.ReceiptNotes,
			&transfer.CancellationReason,
			&transfer.CreatedAt,
			&transfer.UpdatedAt,
		); err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}

	return transfers, rows.Err()
}

func (r *transferRepository) GetInTransitForBranch(ctx context.Context, branchID int64) ([]*domain.ItemTransfer, error) {
	query := `
		SELECT id, transfer_number, item_id, from_branch_id, to_branch_id,
			   status, requested_by, approved_by, received_by,
			   requested_at, approved_at, shipped_at, received_at, cancelled_at,
			   request_notes, approval_notes, receipt_notes, cancellation_reason,
			   created_at, updated_at
		FROM item_transfers
		WHERE to_branch_id = $1 AND status = 'in_transit'
		ORDER BY shipped_at ASC`

	rows, err := r.db.QueryContext(ctx, query, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []*domain.ItemTransfer
	for rows.Next() {
		transfer := &domain.ItemTransfer{}
		if err := rows.Scan(
			&transfer.ID,
			&transfer.TransferNumber,
			&transfer.ItemID,
			&transfer.FromBranchID,
			&transfer.ToBranchID,
			&transfer.Status,
			&transfer.RequestedBy,
			&transfer.ApprovedBy,
			&transfer.ReceivedBy,
			&transfer.RequestedAt,
			&transfer.ApprovedAt,
			&transfer.ShippedAt,
			&transfer.ReceivedAt,
			&transfer.CancelledAt,
			&transfer.RequestNotes,
			&transfer.ApprovalNotes,
			&transfer.ReceiptNotes,
			&transfer.CancellationReason,
			&transfer.CreatedAt,
			&transfer.UpdatedAt,
		); err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}

	return transfers, rows.Err()
}

func (r *transferRepository) GenerateTransferNumber(ctx context.Context) (string, error) {
	now := time.Now()
	prefix := fmt.Sprintf("TR-%s-", now.Format("20060102"))

	query := `
		SELECT COUNT(*) + 1 FROM item_transfers
		WHERE transfer_number LIKE $1`

	var seq int
	if err := r.db.QueryRowContext(ctx, query, prefix+"%").Scan(&seq); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%04d", prefix, seq), nil
}
