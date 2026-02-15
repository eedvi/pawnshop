// User types - mirrors internal/domain/user.go

import type { Branch } from './branch'
import type { Role } from './role'

export interface User {
  id: number
  branch_id?: number
  role_id: number
  email: string
  first_name: string
  last_name: string
  phone?: string
  avatar_url?: string
  is_active: boolean
  email_verified: boolean
  failed_login_attempts?: number
  locked_until?: string
  password_changed_at?: string
  last_login_at?: string
  last_login_ip?: string
  two_factor_enabled: boolean
  two_factor_confirmed_at?: string
  created_at: string
  updated_at: string
  deleted_at?: string
  branch?: Branch
  role?: Role
}

// Safe version for API responses
export interface UserPublic {
  id: number
  email: string
  first_name: string
  last_name: string
  full_name: string
  phone?: string
  avatar_url?: string
  is_active: boolean
  branch_id?: number
  role_id: number
  branch?: Branch
  role?: Role
  created_at: string
}

export interface CreateUserInput {
  email: string
  password: string
  first_name: string
  last_name: string
  phone?: string
  branch_id?: number
  role_id: number
}

export interface UpdateUserInput {
  email?: string
  first_name?: string
  last_name?: string
  phone?: string
  avatar_url?: string
  is_active?: boolean
  branch_id?: number
  role_id?: number
}

export interface UserListParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  search?: string
  branch_id?: number
  role_id?: number
  is_active?: boolean
}
