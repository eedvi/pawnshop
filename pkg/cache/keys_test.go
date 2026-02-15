package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettingKey(t *testing.T) {
	key := SettingKey("app_name")
	assert.Equal(t, "settings:app_name", key)
}

func TestRoleKey(t *testing.T) {
	key := RoleKey(1)
	assert.Equal(t, "roles:1", key)
}

func TestRolePermsKey(t *testing.T) {
	key := RolePermsKey(5)
	assert.Equal(t, "roles:5:permissions", key)
}

func TestUserPermsKey(t *testing.T) {
	key := UserPermsKey(100)
	assert.Equal(t, "users:100:permissions", key)
}

func TestUserKey(t *testing.T) {
	key := UserKey(42)
	assert.Equal(t, "users:42", key)
}

func TestUserEmailKey(t *testing.T) {
	key := UserEmailKey("test@example.com")
	assert.Equal(t, "users:email:test@example.com", key)
}

func TestBranchKey(t *testing.T) {
	key := BranchKey(3)
	assert.Equal(t, "branches:3", key)
}

func TestCategoryKey(t *testing.T) {
	key := CategoryKey(7)
	assert.Equal(t, "categories:7", key)
}

func TestCustomerKey(t *testing.T) {
	key := CustomerKey(123)
	assert.Equal(t, "customers:123", key)
}

func TestItemKey(t *testing.T) {
	key := ItemKey(456)
	assert.Equal(t, "items:456", key)
}

func TestLoanKey(t *testing.T) {
	key := LoanKey(789)
	assert.Equal(t, "loans:789", key)
}

func TestCustomerLoansKey(t *testing.T) {
	key := CustomerLoansKey(100)
	assert.Equal(t, "customers:100:loans", key)
}

func TestRateLimitKey(t *testing.T) {
	key := RateLimitKey("192.168.1.1", "login")
	assert.Equal(t, "ratelimit:192.168.1.1:login", key)
}

func TestSessionKey(t *testing.T) {
	key := SessionKey("abc123")
	assert.Equal(t, "sessions:abc123", key)
}

func TestLockKey(t *testing.T) {
	key := LockKey("loan-creation")
	assert.Equal(t, "lock:loan-creation", key)
}

func TestUserPattern(t *testing.T) {
	pattern := UserPattern(100)
	assert.Equal(t, "users:100:*", pattern)
}

func TestCustomerPattern(t *testing.T) {
	pattern := CustomerPattern(200)
	assert.Equal(t, "customers:200:*", pattern)
}

func TestRolePattern(t *testing.T) {
	pattern := RolePattern(5)
	assert.Equal(t, "roles:5:*", pattern)
}
