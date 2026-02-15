// Loan types - mirrors internal/domain/loan.go

import type { Branch } from './branch'
import type { Customer } from './customer'
import type { Item } from './item'

export type LoanStatus =
  | 'active'
  | 'paid'
  | 'overdue'
  | 'defaulted'
  | 'renewed'
  | 'confiscated'

export type PaymentPlanType = 'single' | 'minimum_payment' | 'installments'

export const LOAN_STATUSES: { value: LoanStatus; label: string; color: string }[] = [
  { value: 'active', label: 'Activo', color: 'green' },
  { value: 'paid', label: 'Pagado', color: 'blue' },
  { value: 'overdue', label: 'Vencido', color: 'red' },
  { value: 'defaulted', label: 'En Mora', color: 'red' },
  { value: 'renewed', label: 'Renovado', color: 'purple' },
  { value: 'confiscated', label: 'Confiscado', color: 'gray' },
]

export const PAYMENT_PLAN_TYPES: { value: PaymentPlanType; label: string; description: string }[] = [
  { value: 'single', label: 'Pago Único', description: 'Un solo pago al vencimiento' },
  { value: 'minimum_payment', label: 'Pago Mínimo', description: 'Pagos mínimos mensuales + pago final' },
  { value: 'installments', label: 'Cuotas', description: 'Pagos iguales divididos en cuotas' },
]

export interface Loan {
  id: number
  loan_number: string
  branch_id: number
  customer_id: number
  item_id: number

  // Amounts
  loan_amount: number
  interest_rate: number
  interest_amount: number
  principal_remaining: number
  interest_remaining: number
  total_amount: number
  amount_paid: number

  // Late fees
  late_fee_rate: number
  late_fee_amount: number

  // Dates
  start_date: string
  due_date: string
  paid_date?: string
  confiscated_date?: string

  // Payment plan
  payment_plan_type: PaymentPlanType
  loan_term_days: number
  requires_minimum_payment: boolean
  minimum_payment_amount?: number
  next_payment_due_date?: string
  grace_period_days: number

  // Installments
  number_of_installments?: number
  installment_amount?: number

  // Status
  status: LoanStatus
  days_overdue: number

  // Renewal
  renewed_from_id?: number
  renewal_count: number

  // Notes
  notes?: string

  // Audit
  created_by?: number
  updated_by?: number
  created_at: string
  updated_at: string
  deleted_at?: string

  // Relations
  branch?: Branch
  customer?: Customer
  item?: Item
}

export interface LoanInstallment {
  id: number
  loan_id: number
  installment_number: number
  due_date: string
  principal_amount: number
  interest_amount: number
  total_amount: number
  amount_paid: number
  is_paid: boolean
  paid_date?: string
  created_at: string
  updated_at: string
}

export interface CreateLoanInput {
  branch_id: number
  customer_id: number
  item_id: number
  loan_amount: number
  interest_rate?: number
  payment_plan_type: PaymentPlanType
  loan_term_days?: number
  grace_period_days?: number
  number_of_installments?: number
  notes?: string
}

export interface RenewLoanInput {
  new_term_days: number
  new_interest_rate?: number
  pay_interest?: boolean
}

export interface LoanListParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  search?: string
  branch_id?: number
  customer_id?: number
  item_id?: number
  status?: LoanStatus
  payment_plan_type?: PaymentPlanType
  is_overdue?: boolean
  date_from?: string
  date_to?: string
}

// Calculate remaining balance
export function calculateRemainingBalance(loan: Loan): number {
  return loan.principal_remaining + loan.interest_remaining + loan.late_fee_amount
}
