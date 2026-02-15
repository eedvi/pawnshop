package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/internal/repository/mocks"
)

func setupExpenseService() (ExpenseService, *mocks.MockExpenseRepository, *mocks.MockExpenseCategoryRepository, *mocks.MockBranchRepository) {
	expenseRepo := new(mocks.MockExpenseRepository)
	categoryRepo := new(mocks.MockExpenseCategoryRepository)
	branchRepo := new(mocks.MockBranchRepository)
	service := NewExpenseService(expenseRepo, categoryRepo, branchRepo)
	return service, expenseRepo, categoryRepo, branchRepo
}

// Category Tests

func TestExpenseService_CreateCategory_Success(t *testing.T) {
	service, _, categoryRepo, _ := setupExpenseService()
	ctx := context.Background()

	categoryRepo.On("GetByCode", ctx, "UTIL").Return(nil, nil)
	categoryRepo.On("Create", ctx, mock.AnythingOfType("*domain.ExpenseCategory")).Return(nil)

	req := CreateExpenseCategoryRequest{
		Name:        "Utilities",
		Code:        "UTIL",
		Description: "Utility expenses",
	}

	result, err := service.CreateCategory(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Utilities", result.Name)
	assert.Equal(t, "UTIL", result.Code)
	assert.True(t, result.IsActive)
	categoryRepo.AssertExpectations(t)
}

func TestExpenseService_CreateCategory_CodeExists(t *testing.T) {
	service, _, categoryRepo, _ := setupExpenseService()
	ctx := context.Background()

	existing := &domain.ExpenseCategory{ID: 1, Code: "UTIL"}
	categoryRepo.On("GetByCode", ctx, "UTIL").Return(existing, nil)

	req := CreateExpenseCategoryRequest{
		Name: "Utilities",
		Code: "UTIL",
	}

	result, err := service.CreateCategory(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "category code already exists", err.Error())
	categoryRepo.AssertExpectations(t)
}

func TestExpenseService_GetCategoryByID_Success(t *testing.T) {
	service, _, categoryRepo, _ := setupExpenseService()
	ctx := context.Background()

	category := &domain.ExpenseCategory{
		ID:       1,
		Name:     "Utilities",
		Code:     "UTIL",
		IsActive: true,
	}
	categoryRepo.On("GetByID", ctx, int64(1)).Return(category, nil)

	result, err := service.GetCategoryByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Utilities", result.Name)
	categoryRepo.AssertExpectations(t)
}

func TestExpenseService_GetCategoryByID_NotFound(t *testing.T) {
	service, _, categoryRepo, _ := setupExpenseService()
	ctx := context.Background()

	categoryRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	result, err := service.GetCategoryByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrCategoryNotFound, err)
	categoryRepo.AssertExpectations(t)
}

func TestExpenseService_UpdateCategory_Success(t *testing.T) {
	service, _, categoryRepo, _ := setupExpenseService()
	ctx := context.Background()

	category := &domain.ExpenseCategory{
		ID:       1,
		Name:     "Utilities",
		Code:     "UTIL",
		IsActive: true,
	}
	categoryRepo.On("GetByID", ctx, int64(1)).Return(category, nil)
	categoryRepo.On("Update", ctx, mock.AnythingOfType("*domain.ExpenseCategory")).Return(nil)

	isActive := false
	req := UpdateExpenseCategoryRequest{
		Name:     "Updated Utilities",
		IsActive: &isActive,
	}

	result, err := service.UpdateCategory(ctx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Utilities", result.Name)
	assert.False(t, result.IsActive)
	categoryRepo.AssertExpectations(t)
}

func TestExpenseService_ListCategories_Success(t *testing.T) {
	service, _, categoryRepo, _ := setupExpenseService()
	ctx := context.Background()

	categories := []*domain.ExpenseCategory{
		{ID: 1, Name: "Utilities", Code: "UTIL", IsActive: true},
		{ID: 2, Name: "Rent", Code: "RENT", IsActive: true},
	}
	categoryRepo.On("List", ctx, false).Return(categories, nil)

	result, err := service.ListCategories(ctx, false)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	categoryRepo.AssertExpectations(t)
}

// Expense Tests

func TestExpenseService_Create_Success(t *testing.T) {
	service, expenseRepo, categoryRepo, branchRepo := setupExpenseService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main Branch", IsActive: true}
	category := &domain.ExpenseCategory{ID: 1, Name: "Utilities", IsActive: true}
	categoryID := int64(1)

	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	categoryRepo.On("GetByID", ctx, int64(1)).Return(category, nil)
	expenseRepo.On("GenerateExpenseNumber", ctx).Return("EXP-001", nil)
	expenseRepo.On("Create", ctx, mock.AnythingOfType("*domain.Expense")).Return(nil)

	req := CreateExpenseRequest{
		BranchID:      1,
		CategoryID:    &categoryID,
		Description:   "Electric bill",
		Amount:        150.00,
		ExpenseDate:   time.Now(),
		PaymentMethod: "cash",
		Vendor:        "Electric Company",
		CreatedBy:     100,
	}

	result, err := service.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "EXP-001", result.ExpenseNumber)
	assert.Equal(t, 150.00, result.Amount)
	expenseRepo.AssertExpectations(t)
	categoryRepo.AssertExpectations(t)
	branchRepo.AssertExpectations(t)
}

func TestExpenseService_Create_BranchNotFound(t *testing.T) {
	service, _, _, branchRepo := setupExpenseService()
	ctx := context.Background()

	branchRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	req := CreateExpenseRequest{
		BranchID:      999,
		Description:   "Electric bill",
		Amount:        150.00,
		ExpenseDate:   time.Now(),
		PaymentMethod: "cash",
		CreatedBy:     100,
	}

	result, err := service.Create(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrBranchNotFound, err)
	branchRepo.AssertExpectations(t)
}

func TestExpenseService_Create_CategoryNotFound(t *testing.T) {
	service, _, categoryRepo, branchRepo := setupExpenseService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main Branch", IsActive: true}
	categoryID := int64(999)

	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	categoryRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	req := CreateExpenseRequest{
		BranchID:      1,
		CategoryID:    &categoryID,
		Description:   "Electric bill",
		Amount:        150.00,
		ExpenseDate:   time.Now(),
		PaymentMethod: "cash",
		CreatedBy:     100,
	}

	result, err := service.Create(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrCategoryNotFound, err)
	branchRepo.AssertExpectations(t)
	categoryRepo.AssertExpectations(t)
}

func TestExpenseService_GetByID_Success(t *testing.T) {
	service, expenseRepo, categoryRepo, branchRepo := setupExpenseService()
	ctx := context.Background()

	categoryID := int64(1)
	expense := &domain.Expense{
		ID:            1,
		ExpenseNumber: "EXP-001",
		BranchID:      1,
		CategoryID:    &categoryID,
		Description:   "Electric bill",
		Amount:        150.00,
	}
	category := &domain.ExpenseCategory{ID: 1, Name: "Utilities"}
	branch := &domain.Branch{ID: 1, Name: "Main Branch"}

	expenseRepo.On("GetByID", ctx, int64(1)).Return(expense, nil)
	categoryRepo.On("GetByID", ctx, int64(1)).Return(category, nil)
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "EXP-001", result.ExpenseNumber)
	assert.NotNil(t, result.Category)
	assert.NotNil(t, result.Branch)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_GetByID_NotFound(t *testing.T) {
	service, expenseRepo, _, _ := setupExpenseService()
	ctx := context.Background()

	expenseRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	result, err := service.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrExpenseNotFound, err)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_Update_Success(t *testing.T) {
	service, expenseRepo, _, _ := setupExpenseService()
	ctx := context.Background()

	expense := &domain.Expense{
		ID:            1,
		ExpenseNumber: "EXP-001",
		BranchID:      1,
		Description:   "Electric bill",
		Amount:        150.00,
		ApprovedBy:    nil, // Not approved
	}

	expenseRepo.On("GetByID", ctx, int64(1)).Return(expense, nil)
	expenseRepo.On("Update", ctx, mock.AnythingOfType("*domain.Expense")).Return(nil)

	req := UpdateExpenseRequest{
		Description: "Updated electric bill",
		Amount:      175.00,
	}

	result, err := service.Update(ctx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated electric bill", result.Description)
	assert.Equal(t, 175.00, result.Amount)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_Update_AlreadyApproved(t *testing.T) {
	service, expenseRepo, _, _ := setupExpenseService()
	ctx := context.Background()

	approvedBy := int64(100)
	expense := &domain.Expense{
		ID:         1,
		ApprovedBy: &approvedBy, // Already approved
	}

	expenseRepo.On("GetByID", ctx, int64(1)).Return(expense, nil)

	req := UpdateExpenseRequest{
		Description: "Try to update",
	}

	result, err := service.Update(ctx, 1, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrExpenseAlreadyApproved, err)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_Delete_Success(t *testing.T) {
	service, expenseRepo, _, _ := setupExpenseService()
	ctx := context.Background()

	expense := &domain.Expense{
		ID:         1,
		ApprovedBy: nil, // Not approved
	}

	expenseRepo.On("GetByID", ctx, int64(1)).Return(expense, nil)
	expenseRepo.On("Delete", ctx, int64(1)).Return(nil)

	err := service.Delete(ctx, 1)

	assert.NoError(t, err)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_Delete_AlreadyApproved(t *testing.T) {
	service, expenseRepo, _, _ := setupExpenseService()
	ctx := context.Background()

	approvedBy := int64(100)
	expense := &domain.Expense{
		ID:         1,
		ApprovedBy: &approvedBy, // Already approved
	}

	expenseRepo.On("GetByID", ctx, int64(1)).Return(expense, nil)

	err := service.Delete(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, ErrExpenseAlreadyApproved, err)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_Delete_NotFound(t *testing.T) {
	service, expenseRepo, _, _ := setupExpenseService()
	ctx := context.Background()

	expenseRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	err := service.Delete(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, ErrExpenseNotFound, err)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_Approve_Success(t *testing.T) {
	service, expenseRepo, _, branchRepo := setupExpenseService()
	ctx := context.Background()

	expense := &domain.Expense{
		ID:            1,
		ExpenseNumber: "EXP-001",
		BranchID:      1,
		Description:   "Electric bill",
		Amount:        150.00,
		ApprovedBy:    nil,
	}

	approvedBy := int64(100)
	approvedExpense := &domain.Expense{
		ID:            1,
		ExpenseNumber: "EXP-001",
		BranchID:      1,
		Description:   "Electric bill",
		Amount:        150.00,
		ApprovedBy:    &approvedBy,
	}

	expenseRepo.On("GetByID", ctx, int64(1)).Return(expense, nil).Once()
	expenseRepo.On("Approve", ctx, int64(1), int64(100)).Return(nil)
	// GetByID called again for reload
	expenseRepo.On("GetByID", ctx, int64(1)).Return(approvedExpense, nil).Once()
	branchRepo.On("GetByID", ctx, int64(1)).Return(&domain.Branch{ID: 1}, nil)

	result, err := service.Approve(ctx, 1, 100)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.ApprovedBy)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_Approve_AlreadyApproved(t *testing.T) {
	service, expenseRepo, _, _ := setupExpenseService()
	ctx := context.Background()

	approvedBy := int64(100)
	expense := &domain.Expense{
		ID:         1,
		ApprovedBy: &approvedBy,
	}

	expenseRepo.On("GetByID", ctx, int64(1)).Return(expense, nil)

	result, err := service.Approve(ctx, 1, 200)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrExpenseAlreadyApproved, err)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_List_Success(t *testing.T) {
	service, expenseRepo, _, _ := setupExpenseService()
	ctx := context.Background()

	expenses := []*domain.Expense{
		{ID: 1, ExpenseNumber: "EXP-001", Amount: 100.00},
		{ID: 2, ExpenseNumber: "EXP-002", Amount: 200.00},
	}

	filter := repository.ExpenseFilter{Page: 1, PageSize: 10}
	expenseRepo.On("List", ctx, filter).Return(expenses, int64(2), nil)

	result, total, err := service.List(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_ListByBranch_Success(t *testing.T) {
	service, expenseRepo, _, _ := setupExpenseService()
	ctx := context.Background()

	expenses := []*domain.Expense{
		{ID: 1, ExpenseNumber: "EXP-001", BranchID: 1, Amount: 100.00},
	}

	filter := repository.ExpenseFilter{Page: 1, PageSize: 10}
	expenseRepo.On("ListByBranch", ctx, int64(1), filter).Return(expenses, int64(1), nil)

	result, total, err := service.ListByBranch(ctx, 1, filter)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_GetTotalByBranchAndDate_Success(t *testing.T) {
	service, expenseRepo, _, _ := setupExpenseService()
	ctx := context.Background()

	date := time.Now()
	expenseRepo.On("GetTotalByBranchAndDate", ctx, int64(1), date).Return(500.00, nil)

	result, err := service.GetTotalByBranchAndDate(ctx, 1, date)

	assert.NoError(t, err)
	assert.Equal(t, 500.00, result)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_GetTotalByBranchAndDate_Error(t *testing.T) {
	service, expenseRepo, _, _ := setupExpenseService()
	ctx := context.Background()

	date := time.Now()
	expenseRepo.On("GetTotalByBranchAndDate", ctx, int64(1), date).Return(0.0, errors.New("db error"))

	result, err := service.GetTotalByBranchAndDate(ctx, 1, date)

	assert.Error(t, err)
	assert.Equal(t, 0.0, result)
	expenseRepo.AssertExpectations(t)
}

// --- Missing coverage tests ---

func TestExpenseService_UpdateCategory_NotFound(t *testing.T) {
	service, _, categoryRepo, _ := setupExpenseService()
	ctx := context.Background()

	categoryRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	req := UpdateExpenseCategoryRequest{Name: "New"}
	result, err := service.UpdateCategory(ctx, 999, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrCategoryNotFound, err)
	categoryRepo.AssertExpectations(t)
}

func TestExpenseService_UpdateCategory_CodeAndDescription(t *testing.T) {
	service, _, categoryRepo, _ := setupExpenseService()
	ctx := context.Background()

	category := &domain.ExpenseCategory{ID: 1, Name: "Old", Code: "OLD", Description: "Old desc", IsActive: true}
	categoryRepo.On("GetByID", ctx, int64(1)).Return(category, nil)
	categoryRepo.On("Update", ctx, mock.AnythingOfType("*domain.ExpenseCategory")).Return(nil)

	req := UpdateExpenseCategoryRequest{
		Code:        "NEW",
		Description: "New desc",
	}
	result, err := service.UpdateCategory(ctx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "NEW", result.Code)
	assert.Equal(t, "New desc", result.Description)
	categoryRepo.AssertExpectations(t)
}

func TestExpenseService_Update_NotFound(t *testing.T) {
	service, expenseRepo, _, _ := setupExpenseService()
	ctx := context.Background()

	expenseRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	req := UpdateExpenseRequest{Description: "Test"}
	result, err := service.Update(ctx, 999, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrExpenseNotFound, err)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_Update_WithCategoryChange(t *testing.T) {
	service, expenseRepo, categoryRepo, _ := setupExpenseService()
	ctx := context.Background()

	expense := &domain.Expense{
		ID:         1,
		BranchID:   1,
		Amount:     100,
		ApprovedBy: nil,
	}
	expenseRepo.On("GetByID", ctx, int64(1)).Return(expense, nil)

	catID := int64(2)
	category := &domain.ExpenseCategory{ID: 2, Name: "Rent"}
	categoryRepo.On("GetByID", ctx, int64(2)).Return(category, nil)
	expenseRepo.On("Update", ctx, mock.AnythingOfType("*domain.Expense")).Return(nil)

	req := UpdateExpenseRequest{
		CategoryID: &catID,
	}
	result, err := service.Update(ctx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, &catID, result.CategoryID)
	categoryRepo.AssertExpectations(t)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_Update_InvalidCategory(t *testing.T) {
	service, expenseRepo, categoryRepo, _ := setupExpenseService()
	ctx := context.Background()

	expense := &domain.Expense{
		ID:         1,
		BranchID:   1,
		Amount:     100,
		ApprovedBy: nil,
	}
	expenseRepo.On("GetByID", ctx, int64(1)).Return(expense, nil)

	catID := int64(999)
	categoryRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	req := UpdateExpenseRequest{CategoryID: &catID}
	result, err := service.Update(ctx, 1, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrCategoryNotFound, err)
}

func TestExpenseService_Update_AllFields(t *testing.T) {
	service, expenseRepo, _, _ := setupExpenseService()
	ctx := context.Background()

	expense := &domain.Expense{
		ID:         1,
		BranchID:   1,
		Amount:     100,
		ApprovedBy: nil,
	}
	expenseRepo.On("GetByID", ctx, int64(1)).Return(expense, nil)
	expenseRepo.On("Update", ctx, mock.AnythingOfType("*domain.Expense")).Return(nil)

	expDate := time.Now()
	req := UpdateExpenseRequest{
		Description:   "Updated description",
		Amount:        200.0,
		ExpenseDate:   expDate,
		PaymentMethod: "card",
		ReceiptNumber: "REC-002",
		Vendor:        "New Vendor",
	}
	result, err := service.Update(ctx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated description", result.Description)
	assert.Equal(t, 200.0, result.Amount)
	assert.Equal(t, "card", result.PaymentMethod)
	assert.Equal(t, "REC-002", result.ReceiptNumber)
	assert.Equal(t, "New Vendor", result.Vendor)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_Approve_NotFound(t *testing.T) {
	service, expenseRepo, _, _ := setupExpenseService()
	ctx := context.Background()

	expenseRepo.On("GetByID", ctx, int64(999)).Return(nil, nil)

	result, err := service.Approve(ctx, 999, 100)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrExpenseNotFound, err)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_Create_NoCategoryProvided(t *testing.T) {
	service, expenseRepo, _, branchRepo := setupExpenseService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main", IsActive: true}
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	expenseRepo.On("GenerateExpenseNumber", ctx).Return("EXP-002", nil)
	expenseRepo.On("Create", ctx, mock.AnythingOfType("*domain.Expense")).Return(nil)

	req := CreateExpenseRequest{
		BranchID:      1,
		Description:   "Miscellaneous",
		Amount:        50.00,
		ExpenseDate:   time.Now(),
		PaymentMethod: "cash",
		CreatedBy:     100,
	}

	result, err := service.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.CategoryID)
	expenseRepo.AssertExpectations(t)
}

func TestExpenseService_GetByID_NoCategoryLoaded(t *testing.T) {
	service, expenseRepo, _, branchRepo := setupExpenseService()
	ctx := context.Background()

	expense := &domain.Expense{
		ID:            1,
		ExpenseNumber: "EXP-001",
		BranchID:      1,
		CategoryID:    nil,
		Amount:        100.00,
	}
	branch := &domain.Branch{ID: 1, Name: "Main"}

	expenseRepo.On("GetByID", ctx, int64(1)).Return(expense, nil)
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.Category)
	assert.NotNil(t, result.Branch)
	expenseRepo.AssertExpectations(t)
}
