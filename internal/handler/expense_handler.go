package handler

import (
	"strconv"
	"time"

	"pawnshop/internal/middleware"
	"pawnshop/internal/repository"
	"pawnshop/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ExpenseHandler struct {
	expenseService service.ExpenseService
}

func NewExpenseHandler(expenseService service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
	}
}

// CreateCategory creates a new expense category
// @Summary Create a new expense category
// @Tags Expenses
// @Accept json
// @Produce json
// @Param category body service.CreateExpenseCategoryRequest true "Category data"
// @Success 201 {object} domain.ExpenseCategory
// @Router /api/v1/expenses/categories [post]
func (h *ExpenseHandler) CreateCategory(c *fiber.Ctx) error {
	var req service.CreateExpenseCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	category, err := h.expenseService.CreateCategory(c.Context(), req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// GetCategoryByID retrieves an expense category by ID
// @Summary Get an expense category by ID
// @Tags Expenses
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} domain.ExpenseCategory
// @Router /api/v1/expenses/categories/{id} [get]
func (h *ExpenseHandler) GetCategoryByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID format",
		})
	}

	category, err := h.expenseService.GetCategoryByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(category)
}

// UpdateCategory updates an expense category
// @Summary Update an expense category
// @Tags Expenses
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body service.UpdateExpenseCategoryRequest true "Category data"
// @Success 200 {object} domain.ExpenseCategory
// @Router /api/v1/expenses/categories/{id} [put]
func (h *ExpenseHandler) UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID format",
		})
	}

	var req service.UpdateExpenseCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	category, err := h.expenseService.UpdateCategory(c.Context(), id, req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(category)
}

// ListCategories retrieves all expense categories
// @Summary List expense categories
// @Tags Expenses
// @Produce json
// @Param include_inactive query bool false "Include inactive categories"
// @Success 200 {array} domain.ExpenseCategory
// @Router /api/v1/expenses/categories [get]
func (h *ExpenseHandler) ListCategories(c *fiber.Ctx) error {
	includeInactive := c.QueryBool("include_inactive", false)

	categories, err := h.expenseService.ListCategories(c.Context(), includeInactive)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(categories)
}

// Create creates a new expense
// @Summary Create a new expense
// @Tags Expenses
// @Accept json
// @Produce json
// @Param expense body service.CreateExpenseRequest true "Expense data"
// @Success 201 {object} domain.Expense
// @Router /api/v1/expenses [post]
func (h *ExpenseHandler) Create(c *fiber.Ctx) error {
	var req service.CreateExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	// Get user ID from context
	userID := c.Locals("userID").(int64)
	req.CreatedBy = userID

	// Default to today if not specified
	if req.ExpenseDate.IsZero() {
		req.ExpenseDate = time.Now()
	}

	expense, err := h.expenseService.Create(c.Context(), req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(expense)
}

// GetByID retrieves an expense by ID
// @Summary Get an expense by ID
// @Tags Expenses
// @Produce json
// @Param id path int true "Expense ID"
// @Success 200 {object} domain.Expense
// @Router /api/v1/expenses/{id} [get]
func (h *ExpenseHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid expense ID format",
		})
	}

	expense, err := h.expenseService.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(expense)
}

// Update updates an expense
// @Summary Update an expense
// @Tags Expenses
// @Accept json
// @Produce json
// @Param id path int true "Expense ID"
// @Param expense body service.UpdateExpenseRequest true "Expense data"
// @Success 200 {object} domain.Expense
// @Router /api/v1/expenses/{id} [put]
func (h *ExpenseHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid expense ID format",
		})
	}

	var req service.UpdateExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body: " + err.Error(),
		})
	}

	expense, err := h.expenseService.Update(c.Context(), id, req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(expense)
}

// Delete deletes an expense
// @Summary Delete an expense
// @Tags Expenses
// @Param id path int true "Expense ID"
// @Success 204
// @Router /api/v1/expenses/{id} [delete]
func (h *ExpenseHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid expense ID format",
		})
	}

	if err := h.expenseService.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// List retrieves expenses with filtering
// @Summary List expenses
// @Tags Expenses
// @Produce json
// @Param branch_id query int false "Filter by branch"
// @Param category_id query int false "Filter by category"
// @Param is_approved query bool false "Filter by approval status"
// @Param date_from query string false "Filter by start date"
// @Param date_to query string false "Filter by end date"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/expenses [get]
func (h *ExpenseHandler) List(c *fiber.Ctx) error {
	filter := repository.ExpenseFilter{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 20),
	}

	if branchID := c.QueryInt("branch_id"); branchID > 0 {
		id := int64(branchID)
		filter.BranchID = &id
	}
	if categoryID := c.QueryInt("category_id"); categoryID > 0 {
		id := int64(categoryID)
		filter.CategoryID = &id
	}
	if c.Query("is_approved") != "" {
		isApproved := c.QueryBool("is_approved")
		filter.IsApproved = &isApproved
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filter.DateFrom = &dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		filter.DateTo = &dateTo
	}

	expenses, total, err := h.expenseService.List(c.Context(), filter)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"data":  expenses,
		"total": total,
		"page":  filter.Page,
		"size":  filter.PageSize,
	})
}

// ListByBranch retrieves expenses for a branch
// @Summary List expenses for a branch
// @Tags Expenses
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Param category_id query int false "Filter by category"
// @Param is_approved query bool false "Filter by approval status"
// @Param date_from query string false "Filter by start date"
// @Param date_to query string false "Filter by end date"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/branches/{branch_id}/expenses [get]
func (h *ExpenseHandler) ListByBranch(c *fiber.Ctx) error {
	branchID, err := strconv.ParseInt(c.Params("branch_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid branch ID format",
		})
	}

	filter := repository.ExpenseFilter{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 20),
	}

	if categoryID := c.QueryInt("category_id"); categoryID > 0 {
		id := int64(categoryID)
		filter.CategoryID = &id
	}
	if c.Query("is_approved") != "" {
		isApproved := c.QueryBool("is_approved")
		filter.IsApproved = &isApproved
	}
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filter.DateFrom = &dateFrom
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		filter.DateTo = &dateTo
	}

	expenses, total, err := h.expenseService.ListByBranch(c.Context(), branchID, filter)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"data":  expenses,
		"total": total,
		"page":  filter.Page,
		"size":  filter.PageSize,
	})
}

// Approve approves an expense
// @Summary Approve an expense
// @Tags Expenses
// @Produce json
// @Param id path int true "Expense ID"
// @Success 200 {object} domain.Expense
// @Router /api/v1/expenses/{id}/approve [post]
func (h *ExpenseHandler) Approve(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid expense ID format",
		})
	}

	userID := c.Locals("userID").(int64)

	expense, err := h.expenseService.Approve(c.Context(), id, userID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(expense)
}

// GetTotalByBranchAndDate retrieves total expenses for a branch on a date
// @Summary Get total expenses for a branch on a date
// @Tags Expenses
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/branches/{branch_id}/expenses/total [get]
func (h *ExpenseHandler) GetTotalByBranchAndDate(c *fiber.Ctx) error {
	branchID, err := strconv.ParseInt(c.Params("branch_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid branch ID format",
		})
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format",
		})
	}

	total, err := h.expenseService.GetTotalByBranchAndDate(c.Context(), branchID, date)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(fiber.Map{
		"branch_id": branchID,
		"date":      dateStr,
		"total":     total,
	})
}

// RegisterRoutes registers expense routes
func (h *ExpenseHandler) RegisterRoutes(router fiber.Router, authMiddleware *middleware.AuthMiddleware) {
	// Expense categories
	categories := router.Group("/expenses/categories")
	categories.Use(authMiddleware.Authenticate())
	categories.Post("/", authMiddleware.RequirePermission("expenses:create"), h.CreateCategory)
	categories.Get("/", authMiddleware.RequirePermission("expenses:read"), h.ListCategories)
	categories.Get("/:id", authMiddleware.RequirePermission("expenses:read"), h.GetCategoryByID)
	categories.Put("/:id", authMiddleware.RequirePermission("expenses:update"), h.UpdateCategory)

	// Expenses
	expenses := router.Group("/expenses")
	expenses.Use(authMiddleware.Authenticate())
	expenses.Post("/", authMiddleware.RequirePermission("expenses:create"), h.Create)
	expenses.Get("/", authMiddleware.RequirePermission("expenses:read"), h.List)
	expenses.Get("/:id", authMiddleware.RequirePermission("expenses:read"), h.GetByID)
	expenses.Put("/:id", authMiddleware.RequirePermission("expenses:update"), h.Update)
	expenses.Delete("/:id", authMiddleware.RequirePermission("expenses:delete"), h.Delete)
	expenses.Post("/:id/approve", authMiddleware.RequirePermission("expenses:approve"), h.Approve)

	// Branch-specific expense routes
	branchExpenses := router.Group("/branches/:branch_id/expenses")
	branchExpenses.Use(authMiddleware.Authenticate())
	branchExpenses.Get("/", authMiddleware.RequirePermission("expenses:read"), h.ListByBranch)
	branchExpenses.Get("/total", authMiddleware.RequirePermission("expenses:read"), h.GetTotalByBranchAndDate)
}
