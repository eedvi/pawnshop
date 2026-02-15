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

func setupCategoryService() (*CategoryService, *mocks.MockCategoryRepository) {
	categoryRepo := new(mocks.MockCategoryRepository)
	service := NewCategoryService(categoryRepo)
	return service, categoryRepo
}

func TestCategoryService_Create_Success(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	categoryRepo.On("GetBySlug", ctx, "electronics").Return(nil, errors.New("not found"))
	categoryRepo.On("Create", ctx, mock.AnythingOfType("*domain.Category")).Return(nil)

	input := CreateCategoryInput{
		Name:                "Electronics",
		DefaultInterestRate: 5.0,
		LoanToValueRatio:    0.6,
	}
	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Electronics", result.Name)
	assert.Equal(t, "electronics", result.Slug)
	assert.True(t, result.IsActive)
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_Create_DuplicateSlug(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	existing := &domain.Category{ID: 1, Slug: "electronics"}
	categoryRepo.On("GetBySlug", ctx, "electronics").Return(existing, nil)

	input := CreateCategoryInput{Name: "Electronics"}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "category with this name already exists", err.Error())
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_Create_InvalidParent(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	categoryRepo.On("GetBySlug", ctx, "sub-cat").Return(nil, errors.New("not found"))

	parentID := int64(999)
	categoryRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := CreateCategoryInput{Name: "Sub Cat", ParentID: &parentID}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "parent category not found", err.Error())
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_Create_DeepNesting(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	categoryRepo.On("GetBySlug", ctx, "deep-cat").Return(nil, errors.New("not found"))

	grandParentID := int64(1)
	parentID := int64(2)
	parent := &domain.Category{ID: 2, ParentID: &grandParentID}
	categoryRepo.On("GetByID", ctx, int64(2)).Return(parent, nil)

	input := CreateCategoryInput{Name: "Deep Cat", ParentID: &parentID}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "categories can only be nested 2 levels deep", err.Error())
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_Create_MinExceedsMax(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	categoryRepo.On("GetBySlug", ctx, "test").Return(nil, errors.New("not found"))

	minAmt := 1000.0
	maxAmt := 500.0
	input := CreateCategoryInput{
		Name:          "Test",
		MinLoanAmount: &minAmt,
		MaxLoanAmount: &maxAmt,
	}
	result, err := service.Create(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "min loan amount cannot exceed max loan amount", err.Error())
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_Create_WithValidParent(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	categoryRepo.On("GetBySlug", ctx, "sub-electronics").Return(nil, errors.New("not found"))

	parentID := int64(1)
	parent := &domain.Category{ID: 1, ParentID: nil}
	categoryRepo.On("GetByID", ctx, int64(1)).Return(parent, nil)
	categoryRepo.On("Create", ctx, mock.AnythingOfType("*domain.Category")).Return(nil)

	input := CreateCategoryInput{Name: "Sub Electronics", ParentID: &parentID}
	result, err := service.Create(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_Update_Success(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	existing := &domain.Category{
		ID:                  1,
		Name:                "Old Name",
		Slug:                "old-name",
		DefaultInterestRate: 5.0,
		IsActive:            true,
	}

	categoryRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)
	categoryRepo.On("Update", ctx, mock.AnythingOfType("*domain.Category")).Return(nil)

	newRate := 10.0
	input := UpdateCategoryInput{
		Name:                "New Name",
		DefaultInterestRate: &newRate,
	}
	result, err := service.Update(ctx, 1, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Name", result.Name)
	assert.Equal(t, "new-name", result.Slug)
	assert.Equal(t, 10.0, result.DefaultInterestRate)
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_Update_NotFound(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	categoryRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	input := UpdateCategoryInput{Name: "New Name"}
	result, err := service.Update(ctx, 999, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "category not found", err.Error())
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_Update_MinExceedsMax(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	minAmt := 1000.0
	existing := &domain.Category{
		ID:            1,
		MinLoanAmount: &minAmt,
	}

	categoryRepo.On("GetByID", ctx, int64(1)).Return(existing, nil)

	newMax := 500.0
	input := UpdateCategoryInput{MaxLoanAmount: &newMax}
	result, err := service.Update(ctx, 1, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "min loan amount cannot exceed max loan amount", err.Error())
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_GetByID_Success(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	category := &domain.Category{ID: 1, Name: "Electronics"}
	categoryRepo.On("GetByID", ctx, int64(1)).Return(category, nil)

	result, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "Electronics", result.Name)
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_GetByID_NotFound(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	categoryRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	result, err := service.GetByID(ctx, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "category not found", err.Error())
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_GetBySlug_Success(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	category := &domain.Category{ID: 1, Slug: "electronics"}
	categoryRepo.On("GetBySlug", ctx, "electronics").Return(category, nil)

	result, err := service.GetBySlug(ctx, "electronics")

	assert.NoError(t, err)
	assert.Equal(t, "electronics", result.Slug)
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_GetBySlug_NotFound(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	categoryRepo.On("GetBySlug", ctx, "missing").Return(nil, errors.New("not found"))

	result, err := service.GetBySlug(ctx, "missing")

	assert.Error(t, err)
	assert.Nil(t, result)
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_List_Success(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	categories := []*domain.Category{
		{ID: 1, Name: "Cat 1"},
		{ID: 2, Name: "Cat 2"},
	}

	params := repository.CategoryListParams{}
	categoryRepo.On("List", ctx, params).Return(categories, nil)

	result, err := service.List(ctx, params)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_ListWithChildren_Success(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	categories := []*domain.Category{
		{ID: 1, Name: "Parent"},
	}

	categoryRepo.On("ListWithChildren", ctx).Return(categories, nil)

	result, err := service.ListWithChildren(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_Delete_Success(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	category := &domain.Category{ID: 1}
	categoryRepo.On("GetByID", ctx, int64(1)).Return(category, nil)
	categoryRepo.On("Delete", ctx, int64(1)).Return(nil)

	err := service.Delete(ctx, 1)

	assert.NoError(t, err)
	categoryRepo.AssertExpectations(t)
}

func TestCategoryService_Delete_NotFound(t *testing.T) {
	service, categoryRepo := setupCategoryService()
	ctx := context.Background()

	categoryRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("not found"))

	err := service.Delete(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, "category not found", err.Error())
	categoryRepo.AssertExpectations(t)
}

// generateSlug tests

func TestGenerateSlug_BasicName(t *testing.T) {
	assert.Equal(t, "electronics", generateSlug("Electronics"))
}

func TestGenerateSlug_WithSpaces(t *testing.T) {
	assert.Equal(t, "gold-and-silver", generateSlug("Gold and Silver"))
}

func TestGenerateSlug_SpecialCharacters(t *testing.T) {
	assert.Equal(t, "electronics-gadgets", generateSlug("Electronics & Gadgets!"))
}

func TestGenerateSlug_MultipleHyphens(t *testing.T) {
	assert.Equal(t, "a-b", generateSlug("a  -  b"))
}

func TestGenerateSlug_LeadingTrailingSpaces(t *testing.T) {
	assert.Equal(t, "test", generateSlug("  test  "))
}

func TestGenerateSlug_NumbersPreserved(t *testing.T) {
	assert.Equal(t, "category-123", generateSlug("Category 123"))
}
