// Branch types - mirrors internal/domain/branch.go

export interface Branch {
  id: number
  name: string
  code: string
  address?: string
  phone?: string
  email?: string
  is_active: boolean
  timezone: string
  currency: string
  default_interest_rate: number
  default_loan_term_days: number
  default_grace_period: number
  created_at: string
  updated_at: string
  deleted_at?: string
}

export interface CreateBranchInput {
  name: string
  code: string
  address?: string
  phone?: string
  email?: string
  timezone?: string
  currency?: string
  default_interest_rate?: number
  default_loan_term_days?: number
  default_grace_period?: number
}

export interface UpdateBranchInput {
  name?: string
  address?: string
  phone?: string
  email?: string
  timezone?: string
  currency?: string
  is_active?: boolean
  default_interest_rate?: number
  default_loan_term_days?: number
  default_grace_period?: number
}
