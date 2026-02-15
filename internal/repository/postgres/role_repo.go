package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pawnshop/internal/domain"
)

// RoleRepository implements repository.RoleRepository
type RoleRepository struct {
	db *DB
}

// NewRoleRepository creates a new RoleRepository
func NewRoleRepository(db *DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// GetByID retrieves a role by ID
func (r *RoleRepository) GetByID(ctx context.Context, id int64) (*domain.Role, error) {
	query := `
		SELECT id, name, display_name, description, permissions, is_system, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	role := &domain.Role{}
	var description sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID, &role.Name, &role.DisplayName, &description,
		&role.Permissions, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	role.Description = StringPtr(description)

	return role, nil
}

// GetByName retrieves a role by name
func (r *RoleRepository) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	query := `
		SELECT id, name, display_name, description, permissions, is_system, created_at, updated_at
		FROM roles
		WHERE name = $1
	`

	role := &domain.Role{}
	var description sql.NullString

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID, &role.Name, &role.DisplayName, &description,
		&role.Permissions, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	role.Description = StringPtr(description)

	return role, nil
}

// List retrieves all roles
func (r *RoleRepository) List(ctx context.Context) ([]*domain.Role, error) {
	query := `
		SELECT id, name, display_name, description, permissions, is_system, created_at, updated_at
		FROM roles
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	roles := []*domain.Role{}
	for rows.Next() {
		role := &domain.Role{}
		var description sql.NullString

		err := rows.Scan(
			&role.ID, &role.Name, &role.DisplayName, &description,
			&role.Permissions, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}

		role.Description = StringPtr(description)
		roles = append(roles, role)
	}

	return roles, nil
}

// Create creates a new role
func (r *RoleRepository) Create(ctx context.Context, role *domain.Role) error {
	query := `
		INSERT INTO roles (name, display_name, description, permissions, is_system)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		role.Name,
		role.DisplayName,
		NullString(role.Description),
		role.Permissions,
		role.IsSystem,
	).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// Update updates an existing role
func (r *RoleRepository) Update(ctx context.Context, role *domain.Role) error {
	// Don't allow updating system roles
	existing, err := r.GetByID(ctx, role.ID)
	if err != nil {
		return err
	}
	if existing.IsSystem {
		return fmt.Errorf("cannot update system role")
	}

	query := `
		UPDATE roles SET
			name = $2,
			display_name = $3,
			description = $4,
			permissions = $5,
			updated_at = NOW()
		WHERE id = $1 AND is_system = false
	`

	result, err := r.db.ExecContext(ctx, query,
		role.ID,
		role.Name,
		role.DisplayName,
		NullString(role.Description),
		role.Permissions,
	)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("role not found or is system role")
	}

	return nil
}

// Delete deletes a role
func (r *RoleRepository) Delete(ctx context.Context, id int64) error {
	// Don't allow deleting system roles
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.IsSystem {
		return fmt.Errorf("cannot delete system role")
	}

	query := `DELETE FROM roles WHERE id = $1 AND is_system = false`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("role not found or is system role")
	}

	return nil
}
