package service

import "errors"

// Common service errors
var (
	// Entity not found errors
	ErrBranchNotFound   = errors.New("branch not found")
	ErrCustomerNotFound = errors.New("customer not found")
	ErrItemNotFound     = errors.New("item not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrLoanNotFound     = errors.New("loan not found")
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrSaleNotFound     = errors.New("sale not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrRoleNotFound     = errors.New("role not found")
	ErrSettingNotFound  = errors.New("setting not found")
	ErrTransferNotFound = errors.New("transfer not found")
	ErrExpenseNotFound  = errors.New("expense not found")

	// Validation errors
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidStatus     = errors.New("invalid status for this operation")
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrInsufficientFunds = errors.New("insufficient funds")

	// Authorization errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	// Operation errors
	ErrOperationFailed = errors.New("operation failed")
	ErrDuplicateEntry  = errors.New("duplicate entry")
	ErrConflict        = errors.New("conflict with existing data")
)
