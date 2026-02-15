package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCashRegister_TableName(t *testing.T) {
	assert.Equal(t, "cash_registers", CashRegister{}.TableName())
}

func TestCashSession_TableName(t *testing.T) {
	assert.Equal(t, "cash_sessions", CashSession{}.TableName())
}

func TestCashSession_IsOpen_Open(t *testing.T) {
	cs := &CashSession{Status: CashSessionStatusOpen}
	assert.True(t, cs.IsOpen())
}

func TestCashSession_IsOpen_Closed(t *testing.T) {
	cs := &CashSession{Status: CashSessionStatusClosed}
	assert.False(t, cs.IsOpen())
}

func TestCashMovement_TableName(t *testing.T) {
	assert.Equal(t, "cash_movements", CashMovement{}.TableName())
}

func TestCashMovement_IsIncome_True(t *testing.T) {
	cm := &CashMovement{MovementType: CashMovementTypeIncome}
	assert.True(t, cm.IsIncome())
}

func TestCashMovement_IsIncome_False(t *testing.T) {
	cm := &CashMovement{MovementType: CashMovementTypeExpense}
	assert.False(t, cm.IsIncome())
}

func TestCashMovement_IsExpense_True(t *testing.T) {
	cm := &CashMovement{MovementType: CashMovementTypeExpense}
	assert.True(t, cm.IsExpense())
}

func TestCashMovement_IsExpense_False(t *testing.T) {
	cm := &CashMovement{MovementType: CashMovementTypeIncome}
	assert.False(t, cm.IsExpense())
}
