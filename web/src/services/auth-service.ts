import { apiPost, apiGet } from '@/lib/api-client'
import {
  LoginInput,
  LoginOutput,
  RefreshInput,
  RefreshOutput,
  TwoFactorVerifyInput,
  UserPublic,
} from '@/types'

/**
 * Authentication service for login, logout, refresh, and 2FA operations
 */
export const authService = {
  /**
   * Login with email and password
   */
  login: async (input: LoginInput): Promise<LoginOutput> => {
    return apiPost<LoginOutput>('/auth/login', input)
  },

  /**
   * Refresh access token using refresh token
   */
  refresh: async (input: RefreshInput): Promise<RefreshOutput> => {
    return apiPost<RefreshOutput>('/auth/refresh', input)
  },

  /**
   * Logout the current user
   */
  logout: async (): Promise<void> => {
    return apiPost('/auth/logout')
  },

  /**
   * Get current authenticated user
   */
  me: async (): Promise<UserPublic> => {
    return apiGet<UserPublic>('/auth/me')
  },

  /**
   * Verify two-factor authentication code
   */
  verifyTwoFactor: async (input: TwoFactorVerifyInput): Promise<LoginOutput> => {
    return apiPost<LoginOutput>('/auth/2fa/verify', input)
  },

  /**
   * Setup two-factor authentication
   * Returns a QR code URL for the authenticator app
   */
  setupTwoFactor: async (): Promise<{ secret: string; qr_code: string }> => {
    return apiPost('/auth/2fa/setup')
  },

  /**
   * Confirm two-factor authentication setup with code
   */
  confirmTwoFactor: async (code: string): Promise<{ backup_codes: string[] }> => {
    return apiPost('/auth/2fa/confirm', { code })
  },

  /**
   * Disable two-factor authentication
   */
  disableTwoFactor: async (password: string): Promise<void> => {
    return apiPost('/auth/2fa/disable', { password })
  },

  /**
   * Request password reset email
   */
  forgotPassword: async (email: string): Promise<void> => {
    return apiPost('/auth/forgot-password', { email })
  },

  /**
   * Reset password with token
   */
  resetPassword: async (token: string, password: string): Promise<void> => {
    return apiPost('/auth/reset-password', { token, password })
  },

  /**
   * Change password (authenticated)
   */
  changePassword: async (currentPassword: string, newPassword: string): Promise<void> => {
    return apiPost('/auth/change-password', {
      current_password: currentPassword,
      new_password: newPassword,
    })
  },
}
