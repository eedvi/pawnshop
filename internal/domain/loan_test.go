package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoan_TableName(t *testing.T) {
	assert.Equal(t, "loans", Loan{}.TableName())
}

func TestLoan_RemainingBalance(t *testing.T) {
	loan := &Loan{
		PrincipalRemaining: 500.0,
		InterestRemaining:  50.0,
		LateFeeAmount:      10.0,
	}
	assert.Equal(t, 560.0, loan.RemainingBalance())
}

func TestLoan_RemainingBalance_Zero(t *testing.T) {
	loan := &Loan{}
	assert.Equal(t, 0.0, loan.RemainingBalance())
}

func TestLoan_IsOverdue_Active_PastDue(t *testing.T) {
	loan := &Loan{
		Status:  LoanStatusActive,
		DueDate: time.Now().AddDate(0, 0, -5),
	}
	assert.True(t, loan.IsOverdue())
}

func TestLoan_IsOverdue_Active_NotPastDue(t *testing.T) {
	loan := &Loan{
		Status:  LoanStatusActive,
		DueDate: time.Now().AddDate(0, 0, 5),
	}
	assert.False(t, loan.IsOverdue())
}

func TestLoan_IsOverdue_NotActive(t *testing.T) {
	loan := &Loan{
		Status:  LoanStatusPaid,
		DueDate: time.Now().AddDate(0, 0, -5),
	}
	assert.False(t, loan.IsOverdue())
}

func TestLoan_IsOverdue_OverdueStatus(t *testing.T) {
	loan := &Loan{
		Status:  LoanStatusOverdue,
		DueDate: time.Now().AddDate(0, 0, -5),
	}
	assert.False(t, loan.IsOverdue())
}

func TestLoan_IsInGracePeriod_WithinGrace(t *testing.T) {
	loan := &Loan{
		Status:          LoanStatusActive,
		DueDate:         time.Now().AddDate(0, 0, -5),
		GracePeriodDays: 15,
	}
	assert.True(t, loan.IsInGracePeriod())
}

func TestLoan_IsInGracePeriod_PastGrace(t *testing.T) {
	loan := &Loan{
		Status:          LoanStatusActive,
		DueDate:         time.Now().AddDate(0, 0, -20),
		GracePeriodDays: 15,
	}
	assert.False(t, loan.IsInGracePeriod())
}

func TestLoan_IsInGracePeriod_NotOverdue(t *testing.T) {
	loan := &Loan{
		Status:          LoanStatusActive,
		DueDate:         time.Now().AddDate(0, 0, 5),
		GracePeriodDays: 15,
	}
	assert.False(t, loan.IsInGracePeriod())
}

func TestLoan_DaysUntilDue_Future(t *testing.T) {
	loan := &Loan{
		DueDate: time.Now().Add(72 * time.Hour),
	}
	days := loan.DaysUntilDue()
	assert.True(t, days >= 2 && days <= 3)
}

func TestLoan_DaysUntilDue_Past(t *testing.T) {
	loan := &Loan{
		DueDate: time.Now().AddDate(0, 0, -5),
	}
	assert.Equal(t, 0, loan.DaysUntilDue())
}

func TestLoan_DaysUntilDue_Today(t *testing.T) {
	loan := &Loan{
		DueDate: time.Now().Add(1 * time.Hour),
	}
	assert.Equal(t, 0, loan.DaysUntilDue())
}

func TestLoan_CalculateDaysOverdue_Overdue(t *testing.T) {
	loan := &Loan{
		Status:  LoanStatusActive,
		DueDate: time.Now().Add(-72 * time.Hour),
	}
	days := loan.CalculateDaysOverdue()
	assert.True(t, days >= 2 && days <= 3)
}

func TestLoan_CalculateDaysOverdue_NotOverdue(t *testing.T) {
	loan := &Loan{
		Status:  LoanStatusActive,
		DueDate: time.Now().AddDate(0, 0, 5),
	}
	assert.Equal(t, 0, loan.CalculateDaysOverdue())
}

func TestLoan_CalculateDaysOverdue_PaidLoan(t *testing.T) {
	loan := &Loan{
		Status:  LoanStatusPaid,
		DueDate: time.Now().AddDate(0, 0, -5),
	}
	assert.Equal(t, 0, loan.CalculateDaysOverdue())
}

func TestLoanInstallment_TableName(t *testing.T) {
	assert.Equal(t, "loan_installments", LoanInstallment{}.TableName())
}

func TestLoanInstallment_RemainingAmount(t *testing.T) {
	li := &LoanInstallment{
		TotalAmount: 500.0,
		AmountPaid:  200.0,
	}
	assert.Equal(t, 300.0, li.RemainingAmount())
}

func TestLoanInstallment_RemainingAmount_FullyPaid(t *testing.T) {
	li := &LoanInstallment{
		TotalAmount: 500.0,
		AmountPaid:  500.0,
	}
	assert.Equal(t, 0.0, li.RemainingAmount())
}

func TestLoanInstallment_RemainingAmount_NoPaid(t *testing.T) {
	li := &LoanInstallment{
		TotalAmount: 500.0,
		AmountPaid:  0,
	}
	assert.Equal(t, 500.0, li.RemainingAmount())
}
