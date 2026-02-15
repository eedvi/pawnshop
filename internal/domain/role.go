package domain

import (
	"encoding/json"
	"time"
)

// Role represents a user role with permissions
type Role struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	DisplayName string          `json:"display_name"`
	Description string          `json:"description,omitempty"`
	Permissions json.RawMessage `json:"permissions"`
	IsSystem    bool            `json:"is_system"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName returns the database table name
func (Role) TableName() string {
	return "roles"
}

// GetPermissions returns the permissions as a string slice
func (r *Role) GetPermissions() ([]string, error) {
	var permissions []string
	if err := json.Unmarshal(r.Permissions, &permissions); err != nil {
		return nil, err
	}
	return permissions, nil
}

// HasPermission checks if the role has a specific permission
func (r *Role) HasPermission(permission string) bool {
	permissions, err := r.GetPermissions()
	if err != nil {
		return false
	}

	for _, p := range permissions {
		if p == "*" || p == permission {
			return true
		}
		// Check wildcard permissions (e.g., "customers.*" matches "customers.read")
		if len(p) > 2 && p[len(p)-2:] == ".*" {
			prefix := p[:len(p)-2]
			if len(permission) > len(prefix) && permission[:len(prefix)] == prefix {
				return true
			}
		}
	}
	return false
}

// Predefined role names
const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleManager    = "manager"
	RoleCashier    = "cashier"
	RoleSeller     = "seller"
)
