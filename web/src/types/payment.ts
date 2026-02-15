// Payment types - mirrors internal/domain/payment.go

import type { Branch } from './branch'
import type { Customer } from './customer'
import type { Loan } from './loan'

export type PaymentMethod = 'cash' | 'card' | 'transfer' | 'check' | 'other'
export type PaymentStatus = 'completed' | 'pending' | 'reversed' | 'failed'

export const PAYMENT_METHODS: { value: PaymentMethod; label: string }[] = [
  { value: 'cash', label: 'Efectivo' },
  { value: 'card', label: 'Tarjeta' },
  { value: 'transfer', label: 'Transferencia' },
  { value: 'check', label: 'Cheque' },
  { value: 'other', label: 'Otro' },
]

export const PAYMENT_STATUSES: { value: PaymentStatus; label: string; color: string }[] = [
  { value: 'completed', label: 'Completado', color: 'green' },
  { value: 'pending', label: 'Pendiente', color: 'yellow' },
  { value: 'reversed', label: 'Revertido', color: 'red' },
  { value: 'failed', label: 'Fallido', color: 'gray' },
]

export interface Payment {
  id: number
  payment_number: string
  branch_id: number
  loan_id: number
  customer_id: number

  // Payment details
  amount: number
  principal_amount: number
  interest_amount: number
  late_fee_amount: number

  // Method
  payment_method: PaymentMethod
  reference_number?: string

  // Status
  status: PaymentStatus
  payment_date: string

  // Balances after payment
  loan_balance_after: number
  interest_balance_after: number

  // Reversal info
  reversed_at?: string
  reversed_by?: number
  reversal_reason?: string

  // Notes
  notes?: string

  // Cash session reference
  cash_session_id?: number

  // Audit
  created_by?: number
  created_at: string
  updated_at: string

  // Relations
  branch?: Branch
  loan?: Loan
  customer?: Customer
}

export interface CreatePaymentInput {
  loan_id: number
  amount: number
  payment_method: PaymentMethod
  reference_number?: string
  notes?: string
}

export interface ReversePaymentInput {
  reason: string
}

export interface PaymentListParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  search?: string
  branch_id?: number
  loan_id?: number
  customer_id?: number
  status?: PaymentStatus
  payment_method?: PaymentMethod
  date_from?: string
  date_to?: string
}

export interface PayoffCalculation {
  loan_id: number
  principal_remaining: number
  interest_remaining: number
  late_fee_amount: number
  total_payoff: number
}

export interface MinimumPaymentCalculation {
  loan_id: number
  minimum_amount: number
  due_date: string
  installment_number?: number
}
