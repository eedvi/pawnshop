import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { expenseService } from '@/services/expense-service'
import type {
  ExpenseListParams,
  CreateExpenseInput,
  UpdateExpenseInput,
  ApproveExpenseInput,
  CreateExpenseCategoryInput,
  UpdateExpenseCategoryInput,
} from '@/types'

export const expenseKeys = {
  all: ['expenses'] as const,
  lists: () => [...expenseKeys.all, 'list'] as const,
  list: (params?: ExpenseListParams) => [...expenseKeys.lists(), params] as const,
  details: () => [...expenseKeys.all, 'detail'] as const,
  detail: (id: number) => [...expenseKeys.details(), id] as const,
  summary: (branchId?: number, dateFrom?: string, dateTo?: string) =>
    [...expenseKeys.all, 'summary', branchId, dateFrom, dateTo] as const,
  categories: () => [...expenseKeys.all, 'categories'] as const,
  categoryDetail: (id: number) => [...expenseKeys.categories(), id] as const,
}

export function useExpenses(params?: ExpenseListParams) {
  return useQuery({
    queryKey: expenseKeys.list(params),
    queryFn: () => expenseService.list(params),
  })
}

export function useExpense(id: number) {
  return useQuery({
    queryKey: expenseKeys.detail(id),
    queryFn: () => expenseService.getById(id),
    enabled: id > 0,
  })
}

export function useExpenseSummary(branchId?: number, dateFrom?: string, dateTo?: string) {
  return useQuery({
    queryKey: expenseKeys.summary(branchId, dateFrom, dateTo),
    queryFn: () => expenseService.getSummary(branchId, dateFrom, dateTo),
  })
}

export function useCreateExpense() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateExpenseInput) => expenseService.create(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: expenseKeys.lists() })
    },
  })
}

export function useUpdateExpense() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: UpdateExpenseInput }) =>
      expenseService.update(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: expenseKeys.lists() })
      queryClient.invalidateQueries({ queryKey: expenseKeys.detail(id) })
    },
  })
}

export function useDeleteExpense() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => expenseService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: expenseKeys.lists() })
    },
  })
}

export function useApproveExpense() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input?: ApproveExpenseInput }) =>
      expenseService.approve(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: expenseKeys.lists() })
      queryClient.invalidateQueries({ queryKey: expenseKeys.detail(id) })
      queryClient.invalidateQueries({ queryKey: expenseKeys.summary() })
    },
  })
}

// Note: useRejectExpense removed - backend doesn't support reject

// Expense Categories
export function useExpenseCategories() {
  return useQuery({
    queryKey: expenseKeys.categories(),
    queryFn: () => expenseService.listCategories(),
  })
}

export function useExpenseCategory(id: number) {
  return useQuery({
    queryKey: expenseKeys.categoryDetail(id),
    queryFn: () => expenseService.getCategoryById(id),
    enabled: id > 0,
  })
}

export function useCreateExpenseCategory() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateExpenseCategoryInput) => expenseService.createCategory(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: expenseKeys.categories() })
    },
  })
}

export function useUpdateExpenseCategory() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: UpdateExpenseCategoryInput }) =>
      expenseService.updateCategory(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: expenseKeys.categories() })
      queryClient.invalidateQueries({ queryKey: expenseKeys.categoryDetail(id) })
    },
  })
}

export function useDeleteExpenseCategory() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => expenseService.deleteCategory(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: expenseKeys.categories() })
    },
  })
}
