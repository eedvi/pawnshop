package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"pawnshop/internal/domain"
	"pawnshop/internal/repository"
	"pawnshop/internal/repository/mocks"
)

func setupItemService() (*ItemService, *mocks.MockItemRepository, *mocks.MockBranchRepository, *mocks.MockCategoryRepository, *mocks.MockCustomerRepository) {
	itemRepo := new(mocks.MockItemRepository)
	branchRepo := new(mocks.MockBranchRepository)
	categoryRepo := new(mocks.MockCategoryRepository)
	customerRepo := new(mocks.MockCustomerRepository)
	service := NewItemService(itemRepo, branchRepo, categoryRepo, customerRepo)
	return service, itemRepo, branchRepo, categoryRepo, customerRepo
}

// --- Create tests ---

func TestItemService_Create_Success_Pawn(t *testing.T) {
	service, itemRepo, branchRepo, _, _ := setupItemService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main", IsActive: true}
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	itemRepo.On("GenerateSKU", ctx, int64(1)).Return("MAIN-000001", nil)
	itemRepo.On("Create", ctx, mock.AnythingOfType("*domain.Item")).Return(nil)
	itemRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.ItemHistory")).Return(nil)

	input := CreateItemInput{
		BranchID:        1,
		Name:            "iPhone 15",
		Condition:       "good",
		AppraisedValue:  1000,
		LoanValue:       800,
		AcquisitionType: "pawn",
		CreatedBy:       1,
	}
	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "MAIN-000001", result.SKU)
	assert.Equal(t, domain.ItemStatusPawned, result.Status)
	itemRepo.AssertExpectations(t)
}

func TestItemService_Create_Success_Purchase(t *testing.T) {
	service, itemRepo, branchRepo, _, _ := setupItemService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main"}
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	itemRepo.On("GenerateSKU", ctx, int64(1)).Return("MAIN-000002", nil)
	itemRepo.On("Create", ctx, mock.AnythingOfType("*domain.Item")).Return(nil)
	itemRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.ItemHistory")).Return(nil)

	input := CreateItemInput{
		BranchID:        1,
		Name:            "MacBook Pro",
		Condition:       "excellent",
		AppraisedValue:  2000,
		LoanValue:       1500,
		AcquisitionType: "purchase",
		CreatedBy:       1,
	}
	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, domain.ItemStatusAvailable, result.Status)
}

func TestItemService_Create_InvalidBranch(t *testing.T) {
	service, _, branchRepo, _, _ := setupItemService()
	ctx := context.Background()

	branchRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := CreateItemInput{
		BranchID:        999,
		Name:            "Test Item",
		Condition:       "good",
		AppraisedValue:  100,
		LoanValue:       80,
		AcquisitionType: "pawn",
	}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid branch", err.Error())
}

func TestItemService_Create_InvalidCategory(t *testing.T) {
	service, _, branchRepo, categoryRepo, _ := setupItemService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main"}
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)

	catID := int64(999)
	categoryRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := CreateItemInput{
		BranchID:        1,
		CategoryID:      &catID,
		Name:            "Test Item",
		Condition:       "good",
		AppraisedValue:  100,
		LoanValue:       80,
		AcquisitionType: "pawn",
	}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid category", err.Error())
}

func TestItemService_Create_InvalidCustomer(t *testing.T) {
	service, _, branchRepo, _, customerRepo := setupItemService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main"}
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)

	custID := int64(999)
	customerRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := CreateItemInput{
		BranchID:        1,
		CustomerID:      &custID,
		Name:            "Test Item",
		Condition:       "good",
		AppraisedValue:  100,
		LoanValue:       80,
		AcquisitionType: "pawn",
	}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid customer", err.Error())
}

func TestItemService_Create_BlockedCustomer(t *testing.T) {
	service, _, branchRepo, _, customerRepo := setupItemService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main"}
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)

	custID := int64(1)
	customer := &domain.Customer{ID: 1, IsBlocked: true}
	customerRepo.On("GetByID", ctx, int64(1)).Return(customer, nil)

	input := CreateItemInput{
		BranchID:        1,
		CustomerID:      &custID,
		Name:            "Test Item",
		Condition:       "good",
		AppraisedValue:  100,
		LoanValue:       80,
		AcquisitionType: "pawn",
	}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "customer is blocked", err.Error())
}

func TestItemService_Create_LoanExceedsAppraised(t *testing.T) {
	service, _, branchRepo, _, _ := setupItemService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main"}
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)

	input := CreateItemInput{
		BranchID:        1,
		Name:            "Test Item",
		Condition:       "good",
		AppraisedValue:  100,
		LoanValue:       200, // Exceeds appraised value
		AcquisitionType: "pawn",
	}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "loan value cannot exceed appraised value", err.Error())
}

func TestItemService_Create_GenerateSKUError(t *testing.T) {
	service, itemRepo, branchRepo, _, _ := setupItemService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main"}
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	itemRepo.On("GenerateSKU", ctx, int64(1)).Return("", errors.New("db error"))

	input := CreateItemInput{
		BranchID:        1,
		Name:            "Test Item",
		Condition:       "good",
		AppraisedValue:  100,
		LoanValue:       80,
		AcquisitionType: "pawn",
	}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to generate SKU", err.Error())
}

func TestItemService_Create_RepoError(t *testing.T) {
	service, itemRepo, branchRepo, _, _ := setupItemService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main"}
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)
	itemRepo.On("GenerateSKU", ctx, int64(1)).Return("MAIN-000001", nil)
	itemRepo.On("Create", ctx, mock.AnythingOfType("*domain.Item")).Return(errors.New("db error"))

	input := CreateItemInput{
		BranchID:        1,
		Name:            "Test Item",
		Condition:       "good",
		AppraisedValue:  100,
		LoanValue:       80,
		AcquisitionType: "pawn",
	}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to create item", err.Error())
}

func TestItemService_Create_WithCategory(t *testing.T) {
	service, itemRepo, branchRepo, categoryRepo, _ := setupItemService()
	ctx := context.Background()

	branch := &domain.Branch{ID: 1, Name: "Main"}
	branchRepo.On("GetByID", ctx, int64(1)).Return(branch, nil)

	catID := int64(5)
	category := &domain.Category{ID: 5, Name: "Electronics"}
	categoryRepo.On("GetByID", ctx, int64(5)).Return(category, nil)

	itemRepo.On("GenerateSKU", ctx, int64(1)).Return("MAIN-000001", nil)
	itemRepo.On("Create", ctx, mock.AnythingOfType("*domain.Item")).Return(nil)
	itemRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.ItemHistory")).Return(nil)

	input := CreateItemInput{
		BranchID:        1,
		CategoryID:      &catID,
		Name:            "iPhone 15",
		Condition:       "good",
		AppraisedValue:  1000,
		LoanValue:       800,
		AcquisitionType: "purchase",
		CreatedBy:       1,
	}
	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, &catID, result.CategoryID)
}

// --- Update tests ---

func TestItemService_Update_Success(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:             1,
		Name:           "iPhone 15",
		Status:         domain.ItemStatusAvailable,
		AppraisedValue: 1000,
		LoanValue:      800,
	}

	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	itemRepo.On("Update", ctx, mock.AnythingOfType("*domain.Item")).Return(nil)

	input := UpdateItemInput{
		Name:      "iPhone 15 Pro",
		UpdatedBy: 1,
	}
	result, err := service.Update(ctx, 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "iPhone 15 Pro", result.Name)
	itemRepo.AssertExpectations(t)
}

func TestItemService_Update_NotFound(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	itemRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := UpdateItemInput{Name: "Test"}
	result, err := service.Update(ctx, 999, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "item not found", err.Error())
}

func TestItemService_Update_InvalidStatus(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:     1,
		Status: domain.ItemStatusSold,
	}
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	input := UpdateItemInput{Name: "Test"}
	result, err := service.Update(ctx, 1, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "cannot update item in current status", err.Error())
}

func TestItemService_Update_InvalidCategory(t *testing.T) {
	service, itemRepo, _, categoryRepo, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:             1,
		Status:         domain.ItemStatusAvailable,
		AppraisedValue: 1000,
		LoanValue:      800,
	}
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	catID := int64(999)
	categoryRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := UpdateItemInput{CategoryID: &catID}
	result, err := service.Update(ctx, 1, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid category", err.Error())
}

func TestItemService_Update_LoanExceedsAppraised(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:             1,
		Status:         domain.ItemStatusAvailable,
		AppraisedValue: 100,
		LoanValue:      80,
	}
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	loanVal := 200.0
	input := UpdateItemInput{LoanValue: &loanVal}
	result, err := service.Update(ctx, 1, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "loan value cannot exceed appraised value", err.Error())
}

func TestItemService_Update_PawnedItem(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:             1,
		Status:         domain.ItemStatusPawned,
		AppraisedValue: 1000,
		LoanValue:      800,
	}
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	itemRepo.On("Update", ctx, mock.AnythingOfType("*domain.Item")).Return(nil)

	input := UpdateItemInput{Name: "Updated Name", UpdatedBy: 1}
	result, err := service.Update(ctx, 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated Name", result.Name)
}

func TestItemService_Update_AllOptionalFields(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:             1,
		Name:           "iPhone",
		Status:         domain.ItemStatusAvailable,
		AppraisedValue: 1000,
		LoanValue:      800,
	}
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	itemRepo.On("Update", ctx, mock.AnythingOfType("*domain.Item")).Return(nil)

	desc := "A phone"
	brand := "Apple"
	model := "15 Pro"
	serial := "SN12345"
	color := "Black"
	appraised := 1200.0
	loanVal := 900.0
	salePrice := 1100.0
	weight := 0.5
	purity := "N/A"
	notes := "Good condition"
	input := UpdateItemInput{
		Name:           "iPhone 15 Pro",
		Description:    &desc,
		Brand:          &brand,
		Model:          &model,
		SerialNumber:   &serial,
		Color:          &color,
		Condition:      "excellent",
		AppraisedValue: &appraised,
		LoanValue:      &loanVal,
		SalePrice:      &salePrice,
		Weight:         &weight,
		Purity:         &purity,
		Notes:          &notes,
		Tags:           []string{"phone", "apple"},
		Photos:         []string{"photo1.jpg", "photo2.jpg"},
		UpdatedBy:      1,
	}
	result, err := service.Update(ctx, 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "iPhone 15 Pro", result.Name)
	assert.Equal(t, &desc, result.Description)
	assert.Equal(t, &brand, result.Brand)
	assert.Equal(t, &model, result.Model)
	assert.Equal(t, &serial, result.SerialNumber)
	assert.Equal(t, &color, result.Color)
	assert.Equal(t, "excellent", result.Condition)
	assert.Equal(t, 1200.0, result.AppraisedValue)
	assert.Equal(t, 900.0, result.LoanValue)
	assert.Equal(t, &salePrice, result.SalePrice)
	assert.Equal(t, 0.5, result.Weight)
	assert.Equal(t, &purity, result.Purity)
	assert.Equal(t, &notes, result.Notes)
	assert.Equal(t, []string{"phone", "apple"}, result.Tags)
	assert.Equal(t, []string{"photo1.jpg", "photo2.jpg"}, result.Photos)
	itemRepo.AssertExpectations(t)
}

func TestItemService_Update_ChangeCategoryWithExisting(t *testing.T) {
	service, itemRepo, _, categoryRepo, _ := setupItemService()
	ctx := context.Background()

	oldCatID := int64(1)
	item := &domain.Item{
		ID:             1,
		Status:         domain.ItemStatusAvailable,
		CategoryID:     &oldCatID,
		AppraisedValue: 1000,
		LoanValue:      800,
	}
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	newCatID := int64(2)
	newCategory := &domain.Category{ID: 2, Name: "Jewelry"}
	categoryRepo.On("GetByID", ctx, int64(2)).Return(newCategory, nil)
	itemRepo.On("Update", ctx, mock.AnythingOfType("*domain.Item")).Return(nil)

	input := UpdateItemInput{CategoryID: &newCatID, UpdatedBy: 1}
	result, err := service.Update(ctx, 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, &newCatID, result.CategoryID)
	categoryRepo.AssertExpectations(t)
	itemRepo.AssertExpectations(t)
}

func TestItemService_Update_SameCategoryNoValidation(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	catID := int64(1)
	item := &domain.Item{
		ID:             1,
		Status:         domain.ItemStatusAvailable,
		CategoryID:     &catID,
		AppraisedValue: 1000,
		LoanValue:      800,
	}
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	itemRepo.On("Update", ctx, mock.AnythingOfType("*domain.Item")).Return(nil)

	sameCatID := int64(1)
	input := UpdateItemInput{CategoryID: &sameCatID, UpdatedBy: 1}
	result, err := service.Update(ctx, 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// CategoryRepo.GetByID should NOT be called since category didn't change
	itemRepo.AssertExpectations(t)
}

func TestItemService_Update_RepoError(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:             1,
		Status:         domain.ItemStatusAvailable,
		AppraisedValue: 1000,
		LoanValue:      800,
	}
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	itemRepo.On("Update", ctx, mock.AnythingOfType("*domain.Item")).Return(errors.New("db error"))

	input := UpdateItemInput{Name: "Updated", UpdatedBy: 1}
	result, err := service.Update(ctx, 1, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to update item", err.Error())
	itemRepo.AssertExpectations(t)
}

// --- GetByID tests ---

func TestItemService_GetByID_Success(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:   1,
		Name: "iPhone 15",
		SKU:  "MAIN-000001",
	}

	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "iPhone 15", result.Name)
}

func TestItemService_GetByID_NotFound(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	itemRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestItemService_GetByID_WithCategory(t *testing.T) {
	service, itemRepo, _, categoryRepo, _ := setupItemService()
	ctx := context.Background()

	catID := int64(5)
	item := &domain.Item{
		ID:         1,
		Name:       "iPhone 15",
		CategoryID: &catID,
	}

	category := &domain.Category{ID: 5, Name: "Electronics"}

	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	categoryRepo.On("GetByID", ctx, int64(5)).Return(category, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Category)
	assert.Equal(t, "Electronics", result.Category.Name)
}

func TestItemService_GetByID_WithCustomer(t *testing.T) {
	service, itemRepo, _, _, customerRepo := setupItemService()
	ctx := context.Background()

	custID := int64(10)
	item := &domain.Item{
		ID:         1,
		Name:       "iPhone 15",
		CustomerID: &custID,
	}

	customer := &domain.Customer{ID: 10, FirstName: "John", LastName: "Doe"}

	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	customerRepo.On("GetByID", ctx, int64(10)).Return(customer, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Customer)
	assert.Equal(t, "John", result.Customer.FirstName)
}

// --- GetBySKU tests ---

func TestItemService_GetBySKU_Success(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:   1,
		Name: "iPhone 15",
		SKU:  "MAIN-000001",
	}

	itemRepo.On("GetBySKU", ctx, "MAIN-000001").Return(item, nil)

	result, err := service.GetBySKU(ctx, "MAIN-000001")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "MAIN-000001", result.SKU)
}

func TestItemService_GetBySKU_NotFound(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	itemRepo.On("GetBySKU", ctx, "NONEXISTENT").Return(nil, errors.New("not found"))

	result, err := service.GetBySKU(ctx, "NONEXISTENT")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "item not found", err.Error())
}

// --- List tests ---

func TestItemService_List_Success(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	items := []domain.Item{
		{ID: 1, Name: "iPhone 15"},
		{ID: 2, Name: "MacBook Pro"},
	}

	paginatedResult := &repository.PaginatedResult[domain.Item]{
		Data:       items,
		Total:      2,
		Page:       1,
		PerPage:    10,
		TotalPages: 1,
	}

	itemRepo.On("List", ctx, mock.AnythingOfType("repository.ItemListParams")).Return(paginatedResult, nil)

	params := repository.ItemListParams{
		BranchID: 1,
	}
	result, err := service.List(ctx, params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
}

// --- Delete tests ---

func TestItemService_Delete_Success(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:     1,
		Name:   "iPhone 15",
		Status: domain.ItemStatusAvailable,
	}

	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	itemRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.ItemHistory")).Return(nil)
	itemRepo.On("Delete", ctx, int64(1)).Return(nil)

	err := service.Delete(ctx, 1, 1)

	assert.NoError(t, err)
	itemRepo.AssertExpectations(t)
}

func TestItemService_Delete_NotFound(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	itemRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := service.Delete(ctx, 999, 1)

	assert.Error(t, err)
	assert.Equal(t, "item not found", err.Error())
}

func TestItemService_Delete_NotAvailable(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:     1,
		Status: domain.ItemStatusPawned,
	}
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	err := service.Delete(ctx, 1, 1)

	assert.Error(t, err)
	assert.Equal(t, "can only delete available items", err.Error())
}

// --- UpdateStatus tests ---

func TestItemService_UpdateStatus_Success(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:     1,
		Name:   "iPhone 15",
		Status: domain.ItemStatusAvailable,
	}

	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	itemRepo.On("UpdateStatus", ctx, int64(1), domain.ItemStatusPawned).Return(nil)
	itemRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.ItemHistory")).Return(nil)

	input := UpdateStatusInput{
		Status:    domain.ItemStatusPawned,
		Notes:     "Pawned to loan",
		UpdatedBy: 1,
	}
	err := service.UpdateStatus(ctx, 1, input)

	assert.NoError(t, err)
	itemRepo.AssertExpectations(t)
}

func TestItemService_UpdateStatus_NotFound(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	itemRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := UpdateStatusInput{
		Status:    domain.ItemStatusPawned,
		UpdatedBy: 1,
	}
	err := service.UpdateStatus(ctx, 999, input)

	assert.Error(t, err)
}

func TestItemService_UpdateStatus_InvalidTransition(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:     1,
		Status: domain.ItemStatusSold,
	}
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	input := UpdateStatusInput{
		Status:    domain.ItemStatusAvailable,
		UpdatedBy: 1,
	}
	err := service.UpdateStatus(ctx, 1, input)

	assert.Error(t, err)
	assert.Equal(t, "invalid status transition", err.Error())
}

func TestItemService_UpdateStatus_RepoError(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:     1,
		Status: domain.ItemStatusAvailable,
	}
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	itemRepo.On("UpdateStatus", ctx, int64(1), domain.ItemStatusPawned).Return(errors.New("db error"))

	input := UpdateStatusInput{
		Status:    domain.ItemStatusPawned,
		UpdatedBy: 1,
	}
	err := service.UpdateStatus(ctx, 1, input)

	assert.Error(t, err)
	assert.Equal(t, "failed to update item status", err.Error())
}

// --- MarkForSale tests ---

func TestItemService_MarkForSale_Success_Confiscated(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:     1,
		Status: domain.ItemStatusConfiscated,
	}

	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	itemRepo.On("Update", ctx, mock.AnythingOfType("*domain.Item")).Return(nil)
	itemRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.ItemHistory")).Return(nil)

	err := service.MarkForSale(ctx, 1, 500.00, 1)

	assert.NoError(t, err)
	itemRepo.AssertExpectations(t)
}

func TestItemService_MarkForSale_Success_Available(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:     1,
		Status: domain.ItemStatusAvailable,
	}

	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	itemRepo.On("Update", ctx, mock.AnythingOfType("*domain.Item")).Return(nil)
	itemRepo.On("CreateHistory", ctx, mock.AnythingOfType("*domain.ItemHistory")).Return(nil)

	err := service.MarkForSale(ctx, 1, 500.00, 1)

	assert.NoError(t, err)
	itemRepo.AssertExpectations(t)
}

func TestItemService_MarkForSale_NotFound(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	itemRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := service.MarkForSale(ctx, 999, 500.00, 1)

	assert.Error(t, err)
	assert.Equal(t, "item not found", err.Error())
}

func TestItemService_MarkForSale_InvalidStatus(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:     1,
		Status: domain.ItemStatusPawned,
	}
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)

	err := service.MarkForSale(ctx, 1, 500.00, 1)

	assert.Error(t, err)
	assert.Equal(t, "can only mark confiscated or available items for sale", err.Error())
}

func TestItemService_MarkForSale_UpdateError(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	item := &domain.Item{
		ID:     1,
		Status: domain.ItemStatusConfiscated,
	}
	itemRepo.On("GetByID", ctx, int64(1)).Return(item, nil)
	itemRepo.On("Update", ctx, mock.AnythingOfType("*domain.Item")).Return(errors.New("db error"))

	err := service.MarkForSale(ctx, 1, 500.00, 1)

	assert.Error(t, err)
	assert.Equal(t, "failed to update item", err.Error())
}

// --- GetAvailableForSale tests ---

func TestItemService_GetAvailableForSale_Success(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	items := []domain.Item{
		{ID: 1, Name: "iPhone 15", Status: domain.ItemStatusForSale},
		{ID: 2, Name: "MacBook Pro", Status: domain.ItemStatusForSale},
	}

	paginatedResult := &repository.PaginatedResult[domain.Item]{
		Data:       items,
		Total:      2,
		Page:       1,
		PerPage:    1000,
		TotalPages: 1,
	}

	itemRepo.On("List", ctx, mock.AnythingOfType("repository.ItemListParams")).Return(paginatedResult, nil)

	result, err := service.GetAvailableForSale(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	itemRepo.AssertExpectations(t)
}

func TestItemService_GetAvailableForSale_Empty(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	paginatedResult := &repository.PaginatedResult[domain.Item]{
		Data:       []domain.Item{},
		Total:      0,
		Page:       1,
		PerPage:    1000,
		TotalPages: 0,
	}

	itemRepo.On("List", ctx, mock.AnythingOfType("repository.ItemListParams")).Return(paginatedResult, nil)

	result, err := service.GetAvailableForSale(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestItemService_GetAvailableForSale_Error(t *testing.T) {
	service, itemRepo, _, _, _ := setupItemService()
	ctx := context.Background()

	itemRepo.On("List", ctx, mock.AnythingOfType("repository.ItemListParams")).Return(nil, errors.New("db error"))

	result, err := service.GetAvailableForSale(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// --- isValidStatusTransition tests ---

func TestIsValidStatusTransition_ValidTransitions(t *testing.T) {
	tests := []struct {
		from domain.ItemStatus
		to   domain.ItemStatus
	}{
		{domain.ItemStatusAvailable, domain.ItemStatusPawned},
		{domain.ItemStatusAvailable, domain.ItemStatusForSale},
		{domain.ItemStatusAvailable, domain.ItemStatusSold},
		{domain.ItemStatusAvailable, domain.ItemStatusTransferred},
		{domain.ItemStatusPawned, domain.ItemStatusAvailable},
		{domain.ItemStatusPawned, domain.ItemStatusConfiscated},
		{domain.ItemStatusForSale, domain.ItemStatusSold},
		{domain.ItemStatusForSale, domain.ItemStatusAvailable},
		{domain.ItemStatusConfiscated, domain.ItemStatusForSale},
		{domain.ItemStatusConfiscated, domain.ItemStatusAvailable},
		{domain.ItemStatusTransferred, domain.ItemStatusAvailable},
	}

	for _, tt := range tests {
		assert.True(t, isValidStatusTransition(tt.from, tt.to), "transition %s -> %s should be valid", tt.from, tt.to)
	}
}

func TestIsValidStatusTransition_InvalidTransitions(t *testing.T) {
	tests := []struct {
		from domain.ItemStatus
		to   domain.ItemStatus
	}{
		{domain.ItemStatusSold, domain.ItemStatusAvailable},
		{domain.ItemStatusSold, domain.ItemStatusPawned},
		{domain.ItemStatusPawned, domain.ItemStatusForSale},
		{domain.ItemStatusForSale, domain.ItemStatusPawned},
		{domain.ItemStatusConfiscated, domain.ItemStatusPawned},
	}

	for _, tt := range tests {
		assert.False(t, isValidStatusTransition(tt.from, tt.to), "transition %s -> %s should be invalid", tt.from, tt.to)
	}
}

func TestIsValidStatusTransition_UnknownStatus(t *testing.T) {
	assert.False(t, isValidStatusTransition(domain.ItemStatus("unknown"), domain.ItemStatusAvailable))
}

// --- Helper function tests ---

func TestStringVal_Nil(t *testing.T) {
	assert.Equal(t, "", stringVal(nil))
}

func TestStringVal_NonNil(t *testing.T) {
	s := "hello"
	assert.Equal(t, "hello", stringVal(&s))
}

func TestFormatCurrency(t *testing.T) {
	result := formatCurrency(100.50)
	assert.Contains(t, result, "Q")
}
