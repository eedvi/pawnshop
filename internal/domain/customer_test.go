package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCustomer_TableName(t *testing.T) {
	assert.Equal(t, "customers", Customer{}.TableName())
}

func TestCustomer_FullName(t *testing.T) {
	c := &Customer{FirstName: "John", LastName: "Doe"}
	assert.Equal(t, "John Doe", c.FullName())
}

func TestCustomer_FullName_Empty(t *testing.T) {
	c := &Customer{}
	assert.Equal(t, " ", c.FullName())
}

func TestCustomer_Age_NilBirthDate(t *testing.T) {
	c := &Customer{BirthDate: nil}
	assert.Equal(t, 0, c.Age())
}

func TestCustomer_Age_Adult(t *testing.T) {
	bd := time.Now().AddDate(-30, 0, 0)
	c := &Customer{BirthDate: &bd}
	assert.Equal(t, 30, c.Age())
}

func TestCustomer_Age_Minor(t *testing.T) {
	bd := time.Now().AddDate(-15, 0, 0)
	c := &Customer{BirthDate: &bd}
	assert.Equal(t, 15, c.Age())
}

func TestCustomer_Age_BirthdayNotYetThisYear(t *testing.T) {
	// Set birthday to tomorrow's date but many years ago
	bd := time.Now().AddDate(-25, 0, 1)
	c := &Customer{BirthDate: &bd}
	assert.Equal(t, 24, c.Age())
}

func TestCustomer_IsAdult_True(t *testing.T) {
	bd := time.Now().AddDate(-20, 0, 0)
	c := &Customer{BirthDate: &bd}
	assert.True(t, c.IsAdult())
}

func TestCustomer_IsAdult_Exactly18(t *testing.T) {
	bd := time.Now().AddDate(-18, 0, 0)
	c := &Customer{BirthDate: &bd}
	assert.True(t, c.IsAdult())
}

func TestCustomer_IsAdult_False(t *testing.T) {
	bd := time.Now().AddDate(-17, 0, 0)
	c := &Customer{BirthDate: &bd}
	assert.False(t, c.IsAdult())
}

func TestCustomer_IsAdult_NilBirthDate(t *testing.T) {
	c := &Customer{BirthDate: nil}
	assert.False(t, c.IsAdult())
}

func TestCustomer_CanTakeLoan_AllConditionsMet(t *testing.T) {
	bd := time.Now().AddDate(-25, 0, 0)
	c := &Customer{
		IsActive:  true,
		IsBlocked: false,
		BirthDate: &bd,
	}
	assert.True(t, c.CanTakeLoan())
}

func TestCustomer_CanTakeLoan_Inactive(t *testing.T) {
	bd := time.Now().AddDate(-25, 0, 0)
	c := &Customer{
		IsActive:  false,
		IsBlocked: false,
		BirthDate: &bd,
	}
	assert.False(t, c.CanTakeLoan())
}

func TestCustomer_CanTakeLoan_Blocked(t *testing.T) {
	bd := time.Now().AddDate(-25, 0, 0)
	c := &Customer{
		IsActive:  true,
		IsBlocked: true,
		BirthDate: &bd,
	}
	assert.False(t, c.CanTakeLoan())
}

func TestCustomer_CanTakeLoan_Minor(t *testing.T) {
	bd := time.Now().AddDate(-16, 0, 0)
	c := &Customer{
		IsActive:  true,
		IsBlocked: false,
		BirthDate: &bd,
	}
	assert.False(t, c.CanTakeLoan())
}

func TestCalculateLoyaltyTier_Standard(t *testing.T) {
	assert.Equal(t, LoyaltyTierStandard, CalculateLoyaltyTier(0))
	assert.Equal(t, LoyaltyTierStandard, CalculateLoyaltyTier(500))
	assert.Equal(t, LoyaltyTierStandard, CalculateLoyaltyTier(999))
}

func TestCalculateLoyaltyTier_Silver(t *testing.T) {
	assert.Equal(t, LoyaltyTierSilver, CalculateLoyaltyTier(1000))
	assert.Equal(t, LoyaltyTierSilver, CalculateLoyaltyTier(3000))
	assert.Equal(t, LoyaltyTierSilver, CalculateLoyaltyTier(4999))
}

func TestCalculateLoyaltyTier_Gold(t *testing.T) {
	assert.Equal(t, LoyaltyTierGold, CalculateLoyaltyTier(5000))
	assert.Equal(t, LoyaltyTierGold, CalculateLoyaltyTier(7500))
	assert.Equal(t, LoyaltyTierGold, CalculateLoyaltyTier(9999))
}

func TestCalculateLoyaltyTier_Platinum(t *testing.T) {
	assert.Equal(t, LoyaltyTierPlatinum, CalculateLoyaltyTier(10000))
	assert.Equal(t, LoyaltyTierPlatinum, CalculateLoyaltyTier(50000))
}

func TestGetLoyaltyDiscount_Standard(t *testing.T) {
	assert.Equal(t, 0.0, GetLoyaltyDiscount(LoyaltyTierStandard))
}

func TestGetLoyaltyDiscount_Silver(t *testing.T) {
	assert.Equal(t, 0.02, GetLoyaltyDiscount(LoyaltyTierSilver))
}

func TestGetLoyaltyDiscount_Gold(t *testing.T) {
	assert.Equal(t, 0.05, GetLoyaltyDiscount(LoyaltyTierGold))
}

func TestGetLoyaltyDiscount_Platinum(t *testing.T) {
	assert.Equal(t, 0.10, GetLoyaltyDiscount(LoyaltyTierPlatinum))
}

func TestGetLoyaltyDiscount_Unknown(t *testing.T) {
	assert.Equal(t, 0.0, GetLoyaltyDiscount("unknown"))
}
