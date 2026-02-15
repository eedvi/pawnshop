// Auth types - mirrors auth handler inputs/outputs

import type { UserPublic } from './user'

export interface LoginInput {
  email: string
  password: string
}

export interface LoginOutput {
  user: UserPublic
  access_token: string
  refresh_token: string
  expires_at: string
  token_type: string
  two_factor_required?: boolean
  two_factor_challenge_token?: string
}

export interface RefreshInput {
  refresh_token: string
}

export interface RefreshOutput {
  access_token: string
  refresh_token: string
  expires_at: string
  token_type: string
}

export interface ChangePasswordInput {
  current_password: string
  new_password: string
}

export interface TwoFactorSetupOutput {
  secret: string
  qr_code: string
  backup_codes: string[]
}

export interface TwoFactorVerifyInput {
  code: string
  challenge_token?: string
}

export interface TwoFactorChallenge {
  challenge_token: string
  expires_at: string
}

// JWT claims structure (decoded from token)
export interface JWTClaims {
  user_id: number
  email: string
  role_id: number
  branch_id?: number
  permissions: string[]
  token_type: string
  exp: number
  iat: number
}
