package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRole_TableName(t *testing.T) {
	assert.Equal(t, "roles", Role{}.TableName())
}

func TestRole_GetPermissions_Success(t *testing.T) {
	r := &Role{
		Permissions: json.RawMessage(`["users.read", "users.create", "loans.read"]`),
	}
	perms, err := r.GetPermissions()
	assert.NoError(t, err)
	assert.Equal(t, []string{"users.read", "users.create", "loans.read"}, perms)
}

func TestRole_GetPermissions_Empty(t *testing.T) {
	r := &Role{
		Permissions: json.RawMessage(`[]`),
	}
	perms, err := r.GetPermissions()
	assert.NoError(t, err)
	assert.Empty(t, perms)
}

func TestRole_GetPermissions_InvalidJSON(t *testing.T) {
	r := &Role{
		Permissions: json.RawMessage(`invalid`),
	}
	perms, err := r.GetPermissions()
	assert.Error(t, err)
	assert.Nil(t, perms)
}

func TestRole_HasPermission_ExactMatch(t *testing.T) {
	r := &Role{
		Permissions: json.RawMessage(`["users.read", "users.create"]`),
	}
	assert.True(t, r.HasPermission("users.read"))
	assert.True(t, r.HasPermission("users.create"))
	assert.False(t, r.HasPermission("users.delete"))
}

func TestRole_HasPermission_SuperWildcard(t *testing.T) {
	r := &Role{
		Permissions: json.RawMessage(`["*"]`),
	}
	assert.True(t, r.HasPermission("users.read"))
	assert.True(t, r.HasPermission("anything"))
}

func TestRole_HasPermission_DomainWildcard(t *testing.T) {
	r := &Role{
		Permissions: json.RawMessage(`["customers.*", "loans.read"]`),
	}
	assert.True(t, r.HasPermission("customers.read"))
	assert.True(t, r.HasPermission("customers.create"))
	assert.True(t, r.HasPermission("customers.delete"))
	assert.False(t, r.HasPermission("users.read"))
	assert.True(t, r.HasPermission("loans.read"))
	assert.False(t, r.HasPermission("loans.create"))
}

func TestRole_HasPermission_InvalidJSON(t *testing.T) {
	r := &Role{
		Permissions: json.RawMessage(`invalid`),
	}
	assert.False(t, r.HasPermission("users.read"))
}

func TestRole_HasPermission_NoMatch(t *testing.T) {
	r := &Role{
		Permissions: json.RawMessage(`["users.read"]`),
	}
	assert.False(t, r.HasPermission("loans.read"))
}
