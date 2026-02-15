package domain

import (
	"time"
)

// User represents a system user
type User struct {
	ID       int64  `json:"id"`
	BranchID *int64 `json:"branch_id,omitempty"`
	RoleID   int64  `json:"role_id"`

	// Credentials
	Email        string `json:"email"`
	PasswordHash string `json:"-"` // Never expose password hash

	// Personal info
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`

	// Status
	IsActive      bool `json:"is_active"`
	EmailVerified bool `json:"email_verified"`

	// Security
	FailedLoginAttempts int        `json:"failed_login_attempts,omitempty"`
	LockedUntil         *time.Time `json:"locked_until,omitempty"`
	PasswordChangedAt   *time.Time `json:"password_changed_at,omitempty"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP         string     `json:"last_login_ip,omitempty"`

	// 2FA
	TwoFactorEnabled     bool       `json:"two_factor_enabled"`
	TwoFactorSecret      string     `json:"-"` // Never expose 2FA secret
	TwoFactorConfirmedAt *time.Time `json:"two_factor_confirmed_at,omitempty"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Relations (loaded when needed)
	Branch *Branch `json:"branch,omitempty"`
	Role   *Role   `json:"role,omitempty"`
}

// TableName returns the database table name
func (User) TableName() string {
	return "users"
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsLocked checks if the user account is locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// CanLogin checks if the user can log in
func (u *User) CanLogin() bool {
	return u.IsActive && !u.IsLocked()
}

// HasPermission checks if the user has a specific permission
func (u *User) HasPermission(permission string) bool {
	if u.Role == nil {
		return false
	}
	return u.Role.HasPermission(permission)
}

// UserPublic is a safe version of User for API responses
type UserPublic struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Phone     string `json:"phone,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	IsActive  bool   `json:"is_active"`

	BranchID *int64 `json:"branch_id,omitempty"`
	RoleID   int64  `json:"role_id"`

	Branch *Branch `json:"branch,omitempty"`
	Role   *Role   `json:"role,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// ToPublic converts User to UserPublic
func (u *User) ToPublic() *UserPublic {
	return &UserPublic{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		FullName:  u.FullName(),
		Phone:     u.Phone,
		AvatarURL: u.AvatarURL,
		IsActive:  u.IsActive,
		BranchID:  u.BranchID,
		RoleID:    u.RoleID,
		Branch:    u.Branch,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
