package cache

import (
	"fmt"
	"time"
)

// Cache key patterns and TTLs
const (
	// Settings
	SettingsKey    = "settings:all"
	SettingKeyFmt  = "settings:%s"
	SettingsTTL    = 5 * time.Minute

	// Roles and Permissions
	RolesAllKey      = "roles:all"
	RoleKeyFmt       = "roles:%d"
	RolePermsFmt     = "roles:%d:permissions"
	UserPermsFmt     = "users:%d:permissions"
	RolesTTL         = 10 * time.Minute
	PermissionsTTL   = 10 * time.Minute

	// Users
	UserKeyFmt     = "users:%d"
	UserEmailFmt   = "users:email:%s"
	UsersTTL       = 5 * time.Minute

	// Branches
	BranchesAllKey = "branches:all"
	BranchKeyFmt   = "branches:%d"
	BranchesTTL    = 15 * time.Minute

	// Categories
	CategoriesAllKey = "categories:all"
	CategoryKeyFmt   = "categories:%d"
	CategoriesTTL    = 15 * time.Minute

	// Customers
	CustomerKeyFmt = "customers:%d"
	CustomersTTL   = 5 * time.Minute

	// Items
	ItemKeyFmt = "items:%d"
	ItemsTTL   = 2 * time.Minute

	// Loans
	LoanKeyFmt         = "loans:%d"
	CustomerLoansFmt   = "customers:%d:loans"
	LoansTTL           = 2 * time.Minute

	// Rate Limiting
	RateLimitKeyFmt = "ratelimit:%s:%s"
	RateLimitTTL    = 1 * time.Minute

	// Sessions
	SessionKeyFmt  = "sessions:%s"
	SessionsTTL    = 24 * time.Hour

	// Locks (for distributed locking)
	LockKeyFmt = "lock:%s"
	LockTTL    = 30 * time.Second
)

// Key generation functions

// SettingKey returns the cache key for a specific setting
func SettingKey(key string) string {
	return fmt.Sprintf(SettingKeyFmt, key)
}

// RoleKey returns the cache key for a role by ID
func RoleKey(roleID int64) string {
	return fmt.Sprintf(RoleKeyFmt, roleID)
}

// RolePermsKey returns the cache key for role permissions
func RolePermsKey(roleID int64) string {
	return fmt.Sprintf(RolePermsFmt, roleID)
}

// UserPermsKey returns the cache key for user permissions
func UserPermsKey(userID int64) string {
	return fmt.Sprintf(UserPermsFmt, userID)
}

// UserKey returns the cache key for a user by ID
func UserKey(userID int64) string {
	return fmt.Sprintf(UserKeyFmt, userID)
}

// UserEmailKey returns the cache key for a user by email
func UserEmailKey(email string) string {
	return fmt.Sprintf(UserEmailFmt, email)
}

// BranchKey returns the cache key for a branch by ID
func BranchKey(branchID int64) string {
	return fmt.Sprintf(BranchKeyFmt, branchID)
}

// CategoryKey returns the cache key for a category by ID
func CategoryKey(categoryID int64) string {
	return fmt.Sprintf(CategoryKeyFmt, categoryID)
}

// CustomerKey returns the cache key for a customer by ID
func CustomerKey(customerID int64) string {
	return fmt.Sprintf(CustomerKeyFmt, customerID)
}

// ItemKey returns the cache key for an item by ID
func ItemKey(itemID int64) string {
	return fmt.Sprintf(ItemKeyFmt, itemID)
}

// LoanKey returns the cache key for a loan by ID
func LoanKey(loanID int64) string {
	return fmt.Sprintf(LoanKeyFmt, loanID)
}

// CustomerLoansKey returns the cache key for a customer's loans
func CustomerLoansKey(customerID int64) string {
	return fmt.Sprintf(CustomerLoansFmt, customerID)
}

// RateLimitKey returns the cache key for rate limiting
func RateLimitKey(identifier, action string) string {
	return fmt.Sprintf(RateLimitKeyFmt, identifier, action)
}

// SessionKey returns the cache key for a session
func SessionKey(sessionID string) string {
	return fmt.Sprintf(SessionKeyFmt, sessionID)
}

// LockKey returns the cache key for a distributed lock
func LockKey(resource string) string {
	return fmt.Sprintf(LockKeyFmt, resource)
}

// Invalidation pattern helpers

// UserPattern returns a pattern to match all user-related keys
func UserPattern(userID int64) string {
	return fmt.Sprintf("users:%d:*", userID)
}

// CustomerPattern returns a pattern to match all customer-related keys
func CustomerPattern(customerID int64) string {
	return fmt.Sprintf("customers:%d:*", customerID)
}

// RolePattern returns a pattern to match all role-related keys
func RolePattern(roleID int64) string {
	return fmt.Sprintf("roles:%d:*", roleID)
}
