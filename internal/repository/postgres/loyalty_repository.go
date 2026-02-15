package postgres

import (
	"context"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

type loyaltyRepository struct {
	db *DB
}

// NewLoyaltyRepository creates a new loyalty repository
func NewLoyaltyRepository(db *DB) repository.LoyaltyRepository {
	return &loyaltyRepository{db: db}
}

func (r *loyaltyRepository) CreateHistory(ctx context.Context, history *domain.LoyaltyPointsHistory) error {
	query := `
		INSERT INTO loyalty_points_history (
			customer_id, branch_id, points_change, points_balance,
			reference_type, reference_id, description, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	return r.db.QueryRowContext(ctx, query,
		history.CustomerID,
		history.BranchID,
		history.PointsChange,
		history.PointsBalance,
		history.ReferenceType,
		history.ReferenceID,
		history.Description,
		history.CreatedBy,
	).Scan(&history.ID, &history.CreatedAt)
}

func (r *loyaltyRepository) GetHistoryByCustomer(ctx context.Context, customerID int64, page, pageSize int) ([]*domain.LoyaltyPointsHistory, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM loyalty_points_history WHERE customer_id = $1`
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, customerID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get history
	if pageSize <= 0 {
		pageSize = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query := `
		SELECT id, customer_id, branch_id, points_change, points_balance,
			   reference_type, reference_id, description, created_by, created_at
		FROM loyalty_points_history
		WHERE customer_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, customerID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var history []*domain.LoyaltyPointsHistory
	for rows.Next() {
		h := &domain.LoyaltyPointsHistory{}
		if err := rows.Scan(
			&h.ID,
			&h.CustomerID,
			&h.BranchID,
			&h.PointsChange,
			&h.PointsBalance,
			&h.ReferenceType,
			&h.ReferenceID,
			&h.Description,
			&h.CreatedBy,
			&h.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		history = append(history, h)
	}

	return history, total, rows.Err()
}

func (r *loyaltyRepository) GetHistoryByReference(ctx context.Context, referenceType string, referenceID int64) ([]*domain.LoyaltyPointsHistory, error) {
	query := `
		SELECT id, customer_id, branch_id, points_change, points_balance,
			   reference_type, reference_id, description, created_by, created_at
		FROM loyalty_points_history
		WHERE reference_type = $1 AND reference_id = $2
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, referenceType, referenceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*domain.LoyaltyPointsHistory
	for rows.Next() {
		h := &domain.LoyaltyPointsHistory{}
		if err := rows.Scan(
			&h.ID,
			&h.CustomerID,
			&h.BranchID,
			&h.PointsChange,
			&h.PointsBalance,
			&h.ReferenceType,
			&h.ReferenceID,
			&h.Description,
			&h.CreatedBy,
			&h.CreatedAt,
		); err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	return history, rows.Err()
}

func (r *loyaltyRepository) GetTotalPointsEarned(ctx context.Context, customerID int64) (int, error) {
	query := `
		SELECT COALESCE(SUM(points_change), 0)
		FROM loyalty_points_history
		WHERE customer_id = $1 AND points_change > 0`

	var total int
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(&total)
	return total, err
}

func (r *loyaltyRepository) GetTotalPointsRedeemed(ctx context.Context, customerID int64) (int, error) {
	query := `
		SELECT COALESCE(ABS(SUM(points_change)), 0)
		FROM loyalty_points_history
		WHERE customer_id = $1 AND points_change < 0`

	var total int
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(&total)
	return total, err
}
