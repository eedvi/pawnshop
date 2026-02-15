package service

import (
	"context"
	"errors"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

var (
	ErrLoyaltyNotEnrolled   = errors.New("customer not enrolled in loyalty program")
	ErrInsufficientPoints   = errors.New("insufficient loyalty points")
	ErrAlreadyEnrolled      = errors.New("customer already enrolled in loyalty program")
)

// Points earning rates
const (
	PointsPerDollarLoan    = 1  // 1 point per dollar of loan principal
	PointsPerDollarPayment = 2  // 2 points per dollar paid
	PointsRedemptionRate   = 100 // 100 points = 1 dollar discount
)

// LoyaltyService defines the interface for loyalty program operations
type LoyaltyService interface {
	// EnrollCustomer enrolls a customer in the loyalty program
	EnrollCustomer(ctx context.Context, customerID int64) error

	// GetCustomerLoyalty gets loyalty information for a customer
	GetCustomerLoyalty(ctx context.Context, customerID int64) (*CustomerLoyaltyInfo, error)

	// AddPoints adds points to a customer's account
	AddPoints(ctx context.Context, req AddPointsRequest) (*domain.LoyaltyPointsHistory, error)

	// RedeemPoints redeems points for a discount
	RedeemPoints(ctx context.Context, req RedeemPointsRequest) (*domain.LoyaltyPointsHistory, error)

	// GetPointsHistory gets the points history for a customer
	GetPointsHistory(ctx context.Context, customerID int64, page, pageSize int) ([]*domain.LoyaltyPointsHistory, int64, error)

	// AwardLoanPoints awards points for a new loan
	AwardLoanPoints(ctx context.Context, customerID, loanID int64, principal float64, awardedBy *int64) error

	// AwardPaymentPoints awards points for a payment
	AwardPaymentPoints(ctx context.Context, customerID, paymentID int64, amount float64, awardedBy *int64) error

	// CalculateDiscount calculates the discount amount for a customer
	CalculateDiscount(ctx context.Context, customerID int64, amount float64) (float64, error)
}

type loyaltyService struct {
	customerRepo repository.CustomerRepository
	loyaltyRepo  repository.LoyaltyRepository
}

// NewLoyaltyService creates a new loyalty service
func NewLoyaltyService(
	customerRepo repository.CustomerRepository,
	loyaltyRepo repository.LoyaltyRepository,
) LoyaltyService {
	return &loyaltyService{
		customerRepo: customerRepo,
		loyaltyRepo:  loyaltyRepo,
	}
}

// CustomerLoyaltyInfo contains loyalty information for a customer
type CustomerLoyaltyInfo struct {
	CustomerID     int64      `json:"customer_id"`
	Points         int        `json:"points"`
	Tier           string     `json:"tier"`
	TierDiscount   float64    `json:"tier_discount"`
	EnrolledAt     *time.Time `json:"enrolled_at,omitempty"`
	IsEnrolled     bool       `json:"is_enrolled"`
	PointsToNextTier int      `json:"points_to_next_tier,omitempty"`
	NextTier       string     `json:"next_tier,omitempty"`
}

// AddPointsRequest represents a request to add loyalty points
type AddPointsRequest struct {
	CustomerID    int64  `json:"customer_id" validate:"required"`
	Points        int    `json:"points" validate:"required,gt=0"`
	ReferenceType string `json:"reference_type"`
	ReferenceID   *int64 `json:"reference_id"`
	Description   string `json:"description"`
	BranchID      *int64 `json:"branch_id"`
	CreatedBy     *int64 `json:"created_by"`
}

// RedeemPointsRequest represents a request to redeem loyalty points
type RedeemPointsRequest struct {
	CustomerID    int64  `json:"customer_id" validate:"required"`
	Points        int    `json:"points" validate:"required,gt=0"`
	ReferenceType string `json:"reference_type"`
	ReferenceID   *int64 `json:"reference_id"`
	Description   string `json:"description"`
	BranchID      *int64 `json:"branch_id"`
	CreatedBy     *int64 `json:"created_by"`
}

func (s *loyaltyService) EnrollCustomer(ctx context.Context, customerID int64) error {
	customer, err := s.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return err
	}
	if customer == nil {
		return ErrCustomerNotFound
	}

	// Check if already enrolled
	if customer.LoyaltyEnrolledAt != nil {
		return ErrAlreadyEnrolled
	}

	// Enroll customer
	now := time.Now()
	customer.LoyaltyEnrolledAt = &now
	customer.LoyaltyTier = domain.LoyaltyTierStandard
	customer.LoyaltyPoints = 0

	return s.customerRepo.Update(ctx, customer)
}

func (s *loyaltyService) GetCustomerLoyalty(ctx context.Context, customerID int64) (*CustomerLoyaltyInfo, error) {
	customer, err := s.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	info := &CustomerLoyaltyInfo{
		CustomerID:   customerID,
		Points:       customer.LoyaltyPoints,
		Tier:         customer.LoyaltyTier,
		TierDiscount: domain.GetLoyaltyDiscount(customer.LoyaltyTier),
		EnrolledAt:   customer.LoyaltyEnrolledAt,
		IsEnrolled:   customer.LoyaltyEnrolledAt != nil,
	}

	// Calculate points to next tier
	if info.Tier == "" {
		info.Tier = domain.LoyaltyTierStandard
	}

	switch info.Tier {
	case domain.LoyaltyTierStandard:
		info.NextTier = domain.LoyaltyTierSilver
		info.PointsToNextTier = domain.LoyaltyTierSilverThreshold - customer.LoyaltyPoints
	case domain.LoyaltyTierSilver:
		info.NextTier = domain.LoyaltyTierGold
		info.PointsToNextTier = domain.LoyaltyTierGoldThreshold - customer.LoyaltyPoints
	case domain.LoyaltyTierGold:
		info.NextTier = domain.LoyaltyTierPlatinum
		info.PointsToNextTier = domain.LoyaltyTierPlatinumThreshold - customer.LoyaltyPoints
	case domain.LoyaltyTierPlatinum:
		info.NextTier = ""
		info.PointsToNextTier = 0
	}

	if info.PointsToNextTier < 0 {
		info.PointsToNextTier = 0
	}

	return info, nil
}

func (s *loyaltyService) AddPoints(ctx context.Context, req AddPointsRequest) (*domain.LoyaltyPointsHistory, error) {
	customer, err := s.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	// Check if enrolled
	if customer.LoyaltyEnrolledAt == nil {
		return nil, ErrLoyaltyNotEnrolled
	}

	// Update points
	newBalance := customer.LoyaltyPoints + req.Points
	customer.LoyaltyPoints = newBalance

	// Update tier if needed
	newTier := domain.CalculateLoyaltyTier(newBalance)
	if newTier != customer.LoyaltyTier {
		customer.LoyaltyTier = newTier
	}

	// Update customer
	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, err
	}

	// Create history record
	history := &domain.LoyaltyPointsHistory{
		CustomerID:    req.CustomerID,
		BranchID:      req.BranchID,
		PointsChange:  req.Points,
		PointsBalance: newBalance,
		ReferenceType: req.ReferenceType,
		ReferenceID:   req.ReferenceID,
		Description:   req.Description,
		CreatedBy:     req.CreatedBy,
		CreatedAt:     time.Now(),
	}

	if err := s.loyaltyRepo.CreateHistory(ctx, history); err != nil {
		return nil, err
	}

	return history, nil
}

func (s *loyaltyService) RedeemPoints(ctx context.Context, req RedeemPointsRequest) (*domain.LoyaltyPointsHistory, error) {
	customer, err := s.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	// Check if enrolled
	if customer.LoyaltyEnrolledAt == nil {
		return nil, ErrLoyaltyNotEnrolled
	}

	// Check if enough points
	if customer.LoyaltyPoints < req.Points {
		return nil, ErrInsufficientPoints
	}

	// Update points (subtract)
	newBalance := customer.LoyaltyPoints - req.Points
	customer.LoyaltyPoints = newBalance

	// Update tier if needed
	newTier := domain.CalculateLoyaltyTier(newBalance)
	if newTier != customer.LoyaltyTier {
		customer.LoyaltyTier = newTier
	}

	// Update customer
	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, err
	}

	// Create history record
	history := &domain.LoyaltyPointsHistory{
		CustomerID:    req.CustomerID,
		BranchID:      req.BranchID,
		PointsChange:  -req.Points, // Negative for redemption
		PointsBalance: newBalance,
		ReferenceType: req.ReferenceType,
		ReferenceID:   req.ReferenceID,
		Description:   req.Description,
		CreatedBy:     req.CreatedBy,
		CreatedAt:     time.Now(),
	}

	if err := s.loyaltyRepo.CreateHistory(ctx, history); err != nil {
		return nil, err
	}

	return history, nil
}

func (s *loyaltyService) GetPointsHistory(ctx context.Context, customerID int64, page, pageSize int) ([]*domain.LoyaltyPointsHistory, int64, error) {
	return s.loyaltyRepo.GetHistoryByCustomer(ctx, customerID, page, pageSize)
}

func (s *loyaltyService) AwardLoanPoints(ctx context.Context, customerID, loanID int64, principal float64, awardedBy *int64) error {
	points := int(principal * PointsPerDollarLoan)
	if points <= 0 {
		return nil
	}

	_, err := s.AddPoints(ctx, AddPointsRequest{
		CustomerID:    customerID,
		Points:        points,
		ReferenceType: "loan",
		ReferenceID:   &loanID,
		Description:   "Points earned from loan",
		CreatedBy:     awardedBy,
	})
	return err
}

func (s *loyaltyService) AwardPaymentPoints(ctx context.Context, customerID, paymentID int64, amount float64, awardedBy *int64) error {
	points := int(amount * PointsPerDollarPayment)
	if points <= 0 {
		return nil
	}

	_, err := s.AddPoints(ctx, AddPointsRequest{
		CustomerID:    customerID,
		Points:        points,
		ReferenceType: "payment",
		ReferenceID:   &paymentID,
		Description:   "Points earned from payment",
		CreatedBy:     awardedBy,
	})
	return err
}

func (s *loyaltyService) CalculateDiscount(ctx context.Context, customerID int64, amount float64) (float64, error) {
	customer, err := s.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return 0, err
	}
	if customer == nil {
		return 0, ErrCustomerNotFound
	}

	discount := domain.GetLoyaltyDiscount(customer.LoyaltyTier)
	return amount * discount, nil
}
