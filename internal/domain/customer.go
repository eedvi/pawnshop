package domain

import (
	"time"
)

// Customer represents a pawnshop customer
type Customer struct {
	ID       int64 `json:"id"`
	BranchID int64 `json:"branch_id"`

	// Personal info
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	IdentityType   string     `json:"identity_type"` // dpi, passport, other
	IdentityNumber string     `json:"identity_number"`
	BirthDate      *time.Time `json:"birth_date,omitempty"`
	Gender         string     `json:"gender,omitempty"` // male, female, other

	// Contact info
	Phone          string `json:"phone"`
	PhoneSecondary string `json:"phone_secondary,omitempty"`
	Email          string `json:"email,omitempty"`
	Address        string `json:"address,omitempty"`
	City           string `json:"city,omitempty"`
	State          string `json:"state,omitempty"`
	PostalCode     string `json:"postal_code,omitempty"`

	// Emergency contact
	EmergencyContactName     string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone    string `json:"emergency_contact_phone,omitempty"`
	EmergencyContactRelation string `json:"emergency_contact_relation,omitempty"`

	// Business info
	Occupation    string  `json:"occupation,omitempty"`
	Workplace     string  `json:"workplace,omitempty"`
	MonthlyIncome float64 `json:"monthly_income,omitempty"`

	// Credit info
	CreditLimit    float64 `json:"credit_limit"`
	CreditScore    int     `json:"credit_score"` // 0-100
	TotalLoans     int     `json:"total_loans"`
	TotalPaid      float64 `json:"total_paid"`
	TotalDefaulted float64 `json:"total_defaulted"`

	// Loyalty program
	LoyaltyPoints     int        `json:"loyalty_points"`
	LoyaltyTier       string     `json:"loyalty_tier"` // standard, silver, gold, platinum
	LoyaltyEnrolledAt *time.Time `json:"loyalty_enrolled_at,omitempty"`

	// Status
	IsActive      bool   `json:"is_active"`
	IsBlocked     bool   `json:"is_blocked"`
	BlockedReason string `json:"blocked_reason,omitempty"`

	// Notes
	Notes string `json:"notes,omitempty"`

	// Photo
	PhotoURL string `json:"photo_url,omitempty"`

	// Audit
	CreatedBy int64 `json:"created_by,omitempty"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Relations
	Branch *Branch `json:"branch,omitempty"`
}

// TableName returns the database table name
func (Customer) TableName() string {
	return "customers"
}

// FullName returns the customer's full name
func (c *Customer) FullName() string {
	return c.FirstName + " " + c.LastName
}

// Age returns the customer's age in years
func (c *Customer) Age() int {
	if c.BirthDate == nil {
		return 0
	}
	now := time.Now()
	years := now.Year() - c.BirthDate.Year()
	if now.YearDay() < c.BirthDate.YearDay() {
		years--
	}
	return years
}

// IsAdult checks if the customer is at least 18 years old
func (c *Customer) IsAdult() bool {
	return c.Age() >= 18
}

// CanTakeLoan checks if the customer can take a new loan
func (c *Customer) CanTakeLoan() bool {
	return c.IsActive && !c.IsBlocked && c.IsAdult()
}

// Identity type constants
const (
	IdentityTypeDPI      = "dpi"
	IdentityTypePassport = "passport"
	IdentityTypeOther    = "other"
)

// Gender constants
const (
	GenderMale   = "male"
	GenderFemale = "female"
	GenderOther  = "other"
)

// Loyalty tier constants
const (
	LoyaltyTierStandard = "standard"
	LoyaltyTierSilver   = "silver"
	LoyaltyTierGold     = "gold"
	LoyaltyTierPlatinum = "platinum"
)

// Loyalty tier thresholds
const (
	LoyaltyTierSilverThreshold   = 1000
	LoyaltyTierGoldThreshold     = 5000
	LoyaltyTierPlatinumThreshold = 10000
)

// LoyaltyPointsHistory represents a loyalty points transaction
type LoyaltyPointsHistory struct {
	ID            int64   `json:"id"`
	CustomerID    int64   `json:"customer_id"`
	BranchID      *int64  `json:"branch_id,omitempty"`
	PointsChange  int     `json:"points_change"`
	PointsBalance int     `json:"points_balance"`
	ReferenceType string  `json:"reference_type,omitempty"` // loan, payment, redemption, bonus, adjustment
	ReferenceID   *int64  `json:"reference_id,omitempty"`
	Description   string  `json:"description,omitempty"`
	CreatedBy     *int64  `json:"created_by,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// CalculateLoyaltyTier determines the loyalty tier based on points
func CalculateLoyaltyTier(points int) string {
	if points >= LoyaltyTierPlatinumThreshold {
		return LoyaltyTierPlatinum
	}
	if points >= LoyaltyTierGoldThreshold {
		return LoyaltyTierGold
	}
	if points >= LoyaltyTierSilverThreshold {
		return LoyaltyTierSilver
	}
	return LoyaltyTierStandard
}

// GetLoyaltyDiscount returns the discount percentage for a tier
func GetLoyaltyDiscount(tier string) float64 {
	switch tier {
	case LoyaltyTierPlatinum:
		return 0.10 // 10% discount
	case LoyaltyTierGold:
		return 0.05 // 5% discount
	case LoyaltyTierSilver:
		return 0.02 // 2% discount
	default:
		return 0 // No discount
	}
}
