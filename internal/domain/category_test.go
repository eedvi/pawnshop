package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategory_TableName(t *testing.T) {
	assert.Equal(t, "categories", Category{}.TableName())
}

func TestCategory_IsRoot_True(t *testing.T) {
	c := &Category{ParentID: nil}
	assert.True(t, c.IsRoot())
}

func TestCategory_IsRoot_False(t *testing.T) {
	parentID := int64(1)
	c := &Category{ParentID: &parentID}
	assert.False(t, c.IsRoot())
}
