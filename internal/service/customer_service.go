package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
)

// CustomerService handles customer business logic
type CustomerService struct {
	customerRepo repository.CustomerRepository
	branchRepo   repository.BranchRepository
}

// NewCustomerService creates a new CustomerService
func NewCustomerService(
	customerRepo repository.CustomerRepository,
	branchRepo repository.BranchRepository,
) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
		branchRepo:   branchRepo,
	}
}

// CreateCustomerInput represents create customer request data
type CreateCustomerInput struct {
	BranchID                 int64  `json:"branch_id" validate:"required"`
	FirstName                string `json:"first_name" validate:"required,min=2"`
	LastName                 string `json:"last_name" validate:"required,min=2"`
	IdentityType             string `json:"identity_type" validate:"required,oneof=dpi passport other"`
	IdentityNumber           string `json:"identity_number" validate:"required"`
	BirthDate                string `json:"birth_date"`
	Gender                   string `json:"gender" validate:"omitempty,oneof=male female other"`
	Phone                    string `json:"phone" validate:"required"`
	PhoneSecondary           string `json:"phone_secondary"`
	Email                    string `json:"email" validate:"omitempty,email"`
	Address                  string `json:"address"`
	City                     string `json:"city"`
	State                    string `json:"state"`
	PostalCode               string `json:"postal_code"`
	EmergencyContactName     string `json:"emergency_contact_name"`
	EmergencyContactPhone    string `json:"emergency_contact_phone"`
	EmergencyContactRelation string `json:"emergency_contact_relation"`
	Occupation               string `json:"occupation"`
	Workplace                string `json:"workplace"`
	MonthlyIncome            float64 `json:"monthly_income"`
	CreditLimit              float64 `json:"credit_limit"`
	Notes                    string  `json:"notes"`
	PhotoURL                 string  `json:"photo_url"`
	CreatedBy                int64   `json:"-"`
}

// Create creates a new customer
func (s *CustomerService) Create(ctx context.Context, input CreateCustomerInput) (*domain.Customer, error) {
	// Validate branch exists
	_, err := s.branchRepo.GetByID(ctx, input.BranchID)
	if err != nil {
		return nil, errors.New("invalid branch")
	}

	// Check for duplicate identity in the same branch
	existing, _ := s.customerRepo.GetByIdentity(ctx, input.BranchID, input.IdentityType, input.IdentityNumber)
	if existing != nil {
		return nil, errors.New("customer with this identity already exists")
	}

	// Parse birth date if provided
	var birthDate *time.Time
	if input.BirthDate != "" {
		parsed, err := time.Parse("2006-01-02", input.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("invalid birth date format, expected YYYY-MM-DD: %w", err)
		}
		birthDate = &parsed

		// Validate age
		age := calculateAge(parsed)
		if age < 18 {
			return nil, errors.New("customer must be at least 18 years old")
		}
	}

	// Create customer
	customer := &domain.Customer{
		BranchID:                 input.BranchID,
		FirstName:                input.FirstName,
		LastName:                 input.LastName,
		IdentityType:             input.IdentityType,
		IdentityNumber:           input.IdentityNumber,
		BirthDate:                birthDate,
		Gender:                   input.Gender,
		Phone:                    input.Phone,
		PhoneSecondary:           input.PhoneSecondary,
		Email:                    input.Email,
		Address:                  input.Address,
		City:                     input.City,
		State:                    input.State,
		PostalCode:               input.PostalCode,
		EmergencyContactName:     input.EmergencyContactName,
		EmergencyContactPhone:    input.EmergencyContactPhone,
		EmergencyContactRelation: input.EmergencyContactRelation,
		Occupation:               input.Occupation,
		Workplace:                input.Workplace,
		MonthlyIncome:            input.MonthlyIncome,
		CreditLimit:              input.CreditLimit,
		CreditScore:              50, // Default credit score
		IsActive:                 true,
		Notes:                    input.Notes,
		PhotoURL:                 input.PhotoURL,
		CreatedBy:                input.CreatedBy,
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

// UpdateCustomerInput represents update customer request data
type UpdateCustomerInput struct {
	FirstName                string     `json:"first_name" validate:"omitempty,min=2"`
	LastName                 string     `json:"last_name" validate:"omitempty,min=2"`
	BirthDate                *time.Time `json:"birth_date"`
	Gender                   string     `json:"gender" validate:"omitempty,oneof=male female other"`
	Phone                    string     `json:"phone"`
	PhoneSecondary           string     `json:"phone_secondary"`
	Email                    string     `json:"email" validate:"omitempty,email"`
	Address                  string     `json:"address"`
	City                     string     `json:"city"`
	State                    string     `json:"state"`
	PostalCode               string     `json:"postal_code"`
	EmergencyContactName     string     `json:"emergency_contact_name"`
	EmergencyContactPhone    string     `json:"emergency_contact_phone"`
	EmergencyContactRelation string     `json:"emergency_contact_relation"`
	Occupation               string     `json:"occupation"`
	Workplace                string     `json:"workplace"`
	MonthlyIncome            *float64   `json:"monthly_income"`
	CreditLimit              *float64   `json:"credit_limit"`
	IsActive                 *bool      `json:"is_active"`
	Notes                    string     `json:"notes"`
	PhotoURL                 string     `json:"photo_url"`
}

// Update updates an existing customer
func (s *CustomerService) Update(ctx context.Context, id int64, input UpdateCustomerInput) (*domain.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Update fields if provided
	if input.FirstName != "" {
		customer.FirstName = input.FirstName
	}
	if input.LastName != "" {
		customer.LastName = input.LastName
	}
	if input.BirthDate != nil {
		age := calculateAge(*input.BirthDate)
		if age < 18 {
			return nil, errors.New("customer must be at least 18 years old")
		}
		customer.BirthDate = input.BirthDate
	}
	if input.Gender != "" {
		customer.Gender = input.Gender
	}
	if input.Phone != "" {
		customer.Phone = input.Phone
	}
	customer.PhoneSecondary = input.PhoneSecondary
	customer.Email = input.Email
	customer.Address = input.Address
	customer.City = input.City
	customer.State = input.State
	customer.PostalCode = input.PostalCode
	customer.EmergencyContactName = input.EmergencyContactName
	customer.EmergencyContactPhone = input.EmergencyContactPhone
	customer.EmergencyContactRelation = input.EmergencyContactRelation
	customer.Occupation = input.Occupation
	customer.Workplace = input.Workplace
	if input.MonthlyIncome != nil {
		customer.MonthlyIncome = *input.MonthlyIncome
	}
	if input.CreditLimit != nil {
		customer.CreditLimit = *input.CreditLimit
	}
	if input.IsActive != nil {
		customer.IsActive = *input.IsActive
	}
	customer.Notes = input.Notes
	customer.PhotoURL = input.PhotoURL

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return customer, nil
}

// GetByID retrieves a customer by ID
func (s *CustomerService) GetByID(ctx context.Context, id int64) (*domain.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Load branch
	customer.Branch, _ = s.branchRepo.GetByID(ctx, customer.BranchID)

	return customer, nil
}

// List retrieves customers with pagination and filters
func (s *CustomerService) List(ctx context.Context, params repository.CustomerListParams) (*repository.PaginatedResult[domain.Customer], error) {
	return s.customerRepo.List(ctx, params)
}

// Delete soft deletes a customer
func (s *CustomerService) Delete(ctx context.Context, id int64) error {
	_, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}

	return s.customerRepo.Delete(ctx, id)
}

// BlockCustomerInput represents block customer request data
type BlockCustomerInput struct {
	CustomerID int64  `json:"customer_id" validate:"required"`
	Reason     string `json:"reason" validate:"required"`
}

// Block blocks a customer
func (s *CustomerService) Block(ctx context.Context, input BlockCustomerInput) error {
	customer, err := s.customerRepo.GetByID(ctx, input.CustomerID)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}

	customer.IsBlocked = true
	customer.BlockedReason = input.Reason

	return s.customerRepo.Update(ctx, customer)
}

// Unblock unblocks a customer
func (s *CustomerService) Unblock(ctx context.Context, id int64) error {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}

	customer.IsBlocked = false
	customer.BlockedReason = ""

	return s.customerRepo.Update(ctx, customer)
}

// UpdateCreditScore updates customer's credit score
func (s *CustomerService) UpdateCreditScore(ctx context.Context, id int64, score int) error {
	if score < 0 || score > 100 {
		return errors.New("credit score must be between 0 and 100")
	}

	return s.customerRepo.UpdateCreditInfo(ctx, id, repository.CustomerCreditUpdate{
		CreditScore: &score,
	})
}

// Helper function to calculate age
func calculateAge(birthDate time.Time) int {
	now := time.Now()
	years := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		years--
	}
	return years
}
