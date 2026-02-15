// Customer types - mirrors internal/domain/customer.go

import type { Branch } from './branch'

export interface Customer {
  id: number
  branch_id: number

  // Personal info
  first_name: string
  last_name: string
  identity_type: IdentityType
  identity_number: string
  birth_date?: string
  gender?: Gender

  // Contact info
  phone: string
  phone_secondary?: string
  email?: string
  address?: string
  city?: string
  state?: string
  postal_code?: string

  // Emergency contact
  emergency_contact_name?: string
  emergency_contact_phone?: string
  emergency_contact_relation?: string

  // Business info
  occupation?: string
  workplace?: string
  monthly_income?: number

  // Credit info
  credit_limit: number
  credit_score: number // 0-100
  total_loans: number
  total_paid: number
  total_defaulted: number

  // Loyalty program
  loyalty_points: number
  loyalty_tier: LoyaltyTier
  loyalty_enrolled_at?: string

  // Status
  is_active: boolean
  is_blocked: boolean
  blocked_reason?: string

  // Notes
  notes?: string
  photo_url?: string

  // Audit
  created_by?: number
  created_at: string
  updated_at: string
  deleted_at?: string

  // Relations
  branch?: Branch
}

export type IdentityType = 'dpi' | 'passport' | 'other'
export type Gender = 'male' | 'female' | 'other'
export type LoyaltyTier = 'standard' | 'silver' | 'gold' | 'platinum'

export const IDENTITY_TYPES: { value: IdentityType; label: string }[] = [
  { value: 'dpi', label: 'DPI' },
  { value: 'passport', label: 'Pasaporte' },
  { value: 'other', label: 'Otro' },
]

export const GENDERS: { value: Gender; label: string }[] = [
  { value: 'male', label: 'Masculino' },
  { value: 'female', label: 'Femenino' },
  { value: 'other', label: 'Otro' },
]

export const LOYALTY_TIERS: { value: LoyaltyTier; label: string; color: string }[] = [
  { value: 'standard', label: 'Estándar', color: 'gray' },
  { value: 'silver', label: 'Plata', color: 'slate' },
  { value: 'gold', label: 'Oro', color: 'yellow' },
  { value: 'platinum', label: 'Platino', color: 'purple' },
]

export interface CreateCustomerInput {
  branch_id: number
  first_name: string
  last_name: string
  identity_type: IdentityType
  identity_number: string
  birth_date?: string
  gender?: Gender
  phone: string
  phone_secondary?: string
  email?: string
  address?: string
  city?: string
  state?: string
  postal_code?: string
  emergency_contact_name?: string
  emergency_contact_phone?: string
  emergency_contact_relation?: string
  occupation?: string
  workplace?: string
  monthly_income?: number
  credit_limit?: number
  notes?: string
}

export interface UpdateCustomerInput {
  first_name?: string
  last_name?: string
  birth_date?: string
  gender?: Gender
  phone?: string
  phone_secondary?: string
  email?: string
  address?: string
  city?: string
  state?: string
  postal_code?: string
  emergency_contact_name?: string
  emergency_contact_phone?: string
  emergency_contact_relation?: string
  occupation?: string
  workplace?: string
  monthly_income?: number
  credit_limit?: number
  notes?: string
  is_active?: boolean
}

export interface CustomerListParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  search?: string
  branch_id?: number
  is_active?: boolean
  is_blocked?: boolean
}

export interface LoyaltyPointsHistory {
  id: number
  customer_id: number
  branch_id?: number
  points_change: number
  points_balance: number
  reference_type?: string
  reference_id?: number
  description?: string
  created_by?: number
  created_at: string
}
