import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { authService } from '@/services/auth-service'
import { useAuthStore } from '@/stores'
import { queryKeys } from '@/lib/query-client'
import { ROUTES } from '@/routes/routes'
import type { LoginInput, TwoFactorVerifyInput } from '@/types'

/**
 * Hook for fetching current user
 */
export function useCurrentUser() {
  const { setUser, setLoading, isAuthenticated } = useAuthStore()

  return useQuery({
    queryKey: queryKeys.auth.me,
    queryFn: async () => {
      const user = await authService.me()
      setUser(user)
      setLoading(false)
      return user
    },
    enabled: isAuthenticated,
    staleTime: 5 * 60 * 1000, // 5 minutes
    retry: false,
  })
}

/**
 * Hook for login mutation
 */
export function useLogin() {
  const navigate = useNavigate()
  const { login } = useAuthStore()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: LoginInput) => authService.login(input),
    onSuccess: (data) => {
      if (data.two_factor_required) {
        // Navigate to 2FA page
        navigate(ROUTES.TWO_FACTOR, {
          state: { challengeToken: data.two_factor_challenge_token },
        })
      } else if (data.user && data.access_token && data.refresh_token) {
        // Complete login
        login(data.user, data.access_token, data.refresh_token)
        queryClient.invalidateQueries({ queryKey: queryKeys.auth.me })
        navigate(ROUTES.DASHBOARD)
      }
    },
  })
}

/**
 * Hook for 2FA verification
 */
export function useTwoFactorVerify() {
  const navigate = useNavigate()
  const { login } = useAuthStore()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: TwoFactorVerifyInput) => authService.verifyTwoFactor(input),
    onSuccess: (data) => {
      if (data.user && data.access_token && data.refresh_token) {
        login(data.user, data.access_token, data.refresh_token)
        queryClient.invalidateQueries({ queryKey: queryKeys.auth.me })
        navigate(ROUTES.DASHBOARD)
      }
    },
  })
}

/**
 * Hook for logout
 */
export function useLogout() {
  const navigate = useNavigate()
  const { logout } = useAuthStore()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => authService.logout(),
    onSettled: () => {
      // Always logout locally, even if API call fails
      logout()
      queryClient.clear()
      navigate(ROUTES.LOGIN)
    },
  })
}

/**
 * Hook for password change
 */
export function useChangePassword() {
  return useMutation({
    mutationFn: ({ currentPassword, newPassword }: { currentPassword: string; newPassword: string }) =>
      authService.changePassword(currentPassword, newPassword),
  })
}

/**
 * Hook for forgot password
 */
export function useForgotPassword() {
  return useMutation({
    mutationFn: (email: string) => authService.forgotPassword(email),
  })
}

/**
 * Hook for reset password
 */
export function useResetPassword() {
  const navigate = useNavigate()

  return useMutation({
    mutationFn: ({ token, password }: { token: string; password: string }) =>
      authService.resetPassword(token, password),
    onSuccess: () => {
      navigate(ROUTES.LOGIN)
    },
  })
}

/**
 * Hook for 2FA setup
 */
export function useTwoFactorSetup() {
  return useMutation({
    mutationFn: () => authService.setupTwoFactor(),
  })
}

/**
 * Hook for 2FA confirmation
 */
export function useTwoFactorConfirm() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (code: string) => authService.confirmTwoFactor(code),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.auth.me })
    },
  })
}

/**
 * Hook for 2FA disable
 */
export function useTwoFactorDisable() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (password: string) => authService.disableTwoFactor(password),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.auth.me })
    },
  })
}
