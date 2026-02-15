package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser_TableName(t *testing.T) {
	assert.Equal(t, "users", User{}.TableName())
}

func TestUser_FullName(t *testing.T) {
	u := &User{FirstName: "Jane", LastName: "Smith"}
	assert.Equal(t, "Jane Smith", u.FullName())
}

func TestUser_IsLocked_NilLockedUntil(t *testing.T) {
	u := &User{LockedUntil: nil}
	assert.False(t, u.IsLocked())
}

func TestUser_IsLocked_LockedInFuture(t *testing.T) {
	future := time.Now().Add(1 * time.Hour)
	u := &User{LockedUntil: &future}
	assert.True(t, u.IsLocked())
}

func TestUser_IsLocked_LockedInPast(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	u := &User{LockedUntil: &past}
	assert.False(t, u.IsLocked())
}

func TestUser_CanLogin_ActiveNotLocked(t *testing.T) {
	u := &User{IsActive: true, LockedUntil: nil}
	assert.True(t, u.CanLogin())
}

func TestUser_CanLogin_Inactive(t *testing.T) {
	u := &User{IsActive: false, LockedUntil: nil}
	assert.False(t, u.CanLogin())
}

func TestUser_CanLogin_Locked(t *testing.T) {
	future := time.Now().Add(1 * time.Hour)
	u := &User{IsActive: true, LockedUntil: &future}
	assert.False(t, u.CanLogin())
}

func TestUser_CanLogin_InactiveAndLocked(t *testing.T) {
	future := time.Now().Add(1 * time.Hour)
	u := &User{IsActive: false, LockedUntil: &future}
	assert.False(t, u.CanLogin())
}

func TestUser_HasPermission_NilRole(t *testing.T) {
	u := &User{Role: nil}
	assert.False(t, u.HasPermission("users.read"))
}

func TestUser_HasPermission_WithRole(t *testing.T) {
	u := &User{
		Role: &Role{
			Permissions: json.RawMessage(`["users.read", "users.create"]`),
		},
	}
	assert.True(t, u.HasPermission("users.read"))
	assert.False(t, u.HasPermission("users.delete"))
}

func TestUser_HasPermission_Wildcard(t *testing.T) {
	u := &User{
		Role: &Role{
			Permissions: json.RawMessage(`["*"]`),
		},
	}
	assert.True(t, u.HasPermission("anything"))
}

func TestUser_ToPublic(t *testing.T) {
	branchID := int64(5)
	u := &User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "12345678",
		AvatarURL: "http://example.com/avatar.png",
		IsActive:  true,
		BranchID:  &branchID,
		RoleID:    2,
		CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	pub := u.ToPublic()

	assert.Equal(t, int64(1), pub.ID)
	assert.Equal(t, "test@example.com", pub.Email)
	assert.Equal(t, "John", pub.FirstName)
	assert.Equal(t, "Doe", pub.LastName)
	assert.Equal(t, "John Doe", pub.FullName)
	assert.Equal(t, "12345678", pub.Phone)
	assert.Equal(t, "http://example.com/avatar.png", pub.AvatarURL)
	assert.True(t, pub.IsActive)
	assert.Equal(t, &branchID, pub.BranchID)
	assert.Equal(t, int64(2), pub.RoleID)
	assert.Equal(t, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), pub.CreatedAt)
}

func TestUser_ToPublic_NilBranch(t *testing.T) {
	u := &User{
		ID:        1,
		FirstName: "Jane",
		LastName:  "Doe",
	}
	pub := u.ToPublic()
	assert.Nil(t, pub.BranchID)
	assert.Equal(t, "Jane Doe", pub.FullName)
}
