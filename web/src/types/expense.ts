// Expense types - mirrors internal/domain/expense.go

import type { Branch } from './branch'
import type { User } from './user'

// Note: Backend doesn't have explicit status - approval is determined by approved_by != null
// We derive status in frontend for display purposes
export type ExpenseStatus = 'pending' | 'approved'

export const EXPENSE_STATUSES: { value: ExpenseStatus; label: string; color: string }[] = [
  { value: 'pending', label: 'Pendiente', color: 'yellow' },
  { value: 'approved', label: 'Aprobado', color: 'green' },
]

// Helper to derive status from expense
export function getExpenseStatus(expense: { approved_by?: number }): ExpenseStatus {
  return expense.approved_by ? 'approved' : 'pending'
}

export interface ExpenseCategory {
  id: number
  code: string
  name: string
  description?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Expense {
  id: number
  expense_number: string
  branch_id: number
  category_id?: number

  // Details
  description: string
  amount: number
  expense_date: string

  // Payment info
  payment_method: string
  receipt_number?: string
  vendor?: string

  // Approval
  approved_by?: number
  approved_at?: string

  // Audit
  created_by?: number
  created_at: string
  updated_at: string

  // Relations
  branch?: Branch
  category?: ExpenseCategory
}

export interface CreateExpenseInput {
  branch_id: number
  category_id?: number
  description: string
  amount: number
  expense_date?: string
  payment_method: string
  receipt_number?: string
  vendor?: string
}

export interface UpdateExpenseInput {
  category_id?: number
  description?: string
  amount?: number
  expense_date?: string
  payment_method?: string
  receipt_number?: string
  vendor?: string
}

export interface ApproveExpenseInput {
  // Approval doesn't need extra fields
}

export interface CreateExpenseCategoryInput {
  code: string
  name: string
  description?: string
}

export interface UpdateExpenseCategoryInput {
  code?: string
  name?: string
  description?: string
  is_active?: boolean
}

export interface ExpenseListParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  search?: string
  branch_id?: number
  category_id?: number
  approved?: boolean
  date_from?: string
  date_to?: string
}

export interface ExpenseSummary {
  total_amount: number
  approved_amount: number
  unapproved_amount: number
  by_category: {
    category_id: number
    category_name: string
    count: number
    total: number
  }[]
}
