package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDailyBalance_TotalIncome(t *testing.T) {
	d := &DailyBalance{
		InterestIncome: 100.0,
		LateFeeIncome:  25.0,
		SalesIncome:    500.0,
		OtherIncome:    50.0,
	}
	assert.Equal(t, 675.0, d.TotalIncome())
}

func TestDailyBalance_TotalIncome_Zero(t *testing.T) {
	d := &DailyBalance{}
	assert.Equal(t, 0.0, d.TotalIncome())
}

func TestDailyBalance_TotalExpenses(t *testing.T) {
	d := &DailyBalance{
		OperationalExpenses: 200.0,
		Refunds:             50.0,
		OtherExpenses:       30.0,
		LoanDisbursements:   1000.0,
	}
	assert.Equal(t, 1280.0, d.TotalExpenses())
}

func TestDailyBalance_TotalExpenses_Zero(t *testing.T) {
	d := &DailyBalance{}
	assert.Equal(t, 0.0, d.TotalExpenses())
}

func TestExpense_IsApproved_True(t *testing.T) {
	approver := int64(1)
	e := &Expense{ApprovedBy: &approver}
	assert.True(t, e.IsApproved())
}

func TestExpense_IsApproved_False(t *testing.T) {
	e := &Expense{ApprovedBy: nil}
	assert.False(t, e.IsApproved())
}
