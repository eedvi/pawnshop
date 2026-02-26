// Cash types - mirrors internal/domain/cash*.go

import type { Branch } from './branch'
import type { User } from './user'

export type CashSessionStatus = 'open' | 'closed'
export type CashMovementType = 'income' | 'expense'

export const CASH_SESSION_STATUSES: { value: CashSessionStatus; label: string; color: string }[] = [
  { value: 'open', label: 'Abierta', color: 'green' },
  { value: 'closed', label: 'Cerrada', color: 'gray' },
]

export const CASH_MOVEMENT_TYPES: { value: CashMovementType; label: string; color: string }[] = [
  { value: 'income', label: 'Ingreso', color: 'green' },
  { value: 'expense', label: 'Egreso', color: 'red' },
]

export interface CashRegister {
  id: number
  branch_id: number
  name: string
  code: string
  description?: string
  is_active: boolean
  created_at: string
  updated_at: string
  branch?: Branch
}

export interface CashSession {
  id: number
  cash_register_id: number
  branch_id: number
  user_id: number
  opening_amount: number
  closing_amount?: number
  expected_amount?: number
  difference?: number
  status: CashSessionStatus
  opened_at: string
  closed_at?: string
  opening_notes?: string
  closing_notes?: string
  closed_by?: number
  created_at: string
  updated_at: string
  register?: CashRegister
  branch?: Branch
  user?: User
  movements?: CashMovement[]
}

export interface CashMovement {
  id: number
  branch_id: number
  session_id: number
  movement_type: CashMovementType
  amount: number
  payment_method: string
  reference_type?: string
  reference_id?: number
  description: string
  balance_after: number
  created_by?: number
  created_at: string
  cash_session?: CashSession
  branch?: Branch
}

export interface CreateCashRegisterInput {
  branch_id: number
  name: string
  description?: string
}

export interface UpdateCashRegisterInput {
  name?: string
  description?: string
  is_active?: boolean
}

export interface OpenCashSessionInput {
  cash_register_id: number
  opening_amount: number
  opening_notes?: string
}

export interface CloseCashSessionInput {
  closing_amount: number
  closing_notes?: string
}

export interface CreateCashMovementInput {
  branch_id: number
  session_id: number
  movement_type: CashMovementType
  amount: number
  payment_method: string
  reference_type?: string
  reference_id?: number
  description: string
}

export interface CashSessionSummary {
  session_id: number
  opening_amount: number
  total_income: number
  total_expense: number
  expected_balance: number
  closing_amount?: number
  difference?: number
  movements_count: number
  movements_by_type: {
    type: CashMovementType
    count: number
    total: number
  }[]
}

export interface CashSessionListParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  register_id?: number
  branch_id?: number
  user_id?: number
  status?: CashSessionStatus
  date_from?: string
  date_to?: string
}
