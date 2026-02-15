import { apiGet, apiGetPaginated, apiPost, apiPut, apiDelete } from '@/lib/api-client'
import type {
  Expense,
  ExpenseCategory,
  ExpenseListParams,
  ExpenseSummary,
  CreateExpenseInput,
  UpdateExpenseInput,
  ApproveExpenseInput,
  CreateExpenseCategoryInput,
  UpdateExpenseCategoryInput,
} from '@/types'

export const expenseService = {
  // Expenses
  list: (params?: ExpenseListParams) => {
    // Convert frontend 'approved' param to backend 'is_approved'
    const backendParams = params ? {
      ...params,
      is_approved: params.approved,
      approved: undefined, // Remove frontend param
    } : undefined
    return apiGetPaginated<Expense>('/expenses', backendParams)
  },

  getById: (id: number) =>
    apiGet<Expense>(`/expenses/${id}`),

  create: (input: CreateExpenseInput) =>
    apiPost<Expense>('/expenses', input),

  update: (id: number, input: UpdateExpenseInput) =>
    apiPut<Expense>(`/expenses/${id}`, input),

  delete: (id: number) =>
    apiDelete<void>(`/expenses/${id}`),

  approve: (id: number, _input?: ApproveExpenseInput) =>
    apiPost<Expense>(`/expenses/${id}/approve`, {}),

  // Note: Backend doesn't support reject - removed

  getSummary: (branchId?: number, dateFrom?: string, dateTo?: string) =>
    apiGet<ExpenseSummary>('/expenses/summary', { branch_id: branchId, date_from: dateFrom, date_to: dateTo }),

  // Expense Categories
  listCategories: () =>
    apiGet<ExpenseCategory[]>('/expenses/categories'),

  getCategoryById: (id: number) =>
    apiGet<ExpenseCategory>(`/expenses/categories/${id}`),

  createCategory: (input: CreateExpenseCategoryInput) =>
    apiPost<ExpenseCategory>('/expenses/categories', input),

  updateCategory: (id: number, input: UpdateExpenseCategoryInput) =>
    apiPut<ExpenseCategory>(`/expenses/categories/${id}`, input),

  deleteCategory: (id: number) =>
    apiDelete<void>(`/expenses/categories/${id}`),
}
