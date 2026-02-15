import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import { UserPublic, JWTClaims } from '@/types'
import { tokenStorage } from '@/lib/api-client'

interface AuthState {
  // State
  user: UserPublic | null
  permissions: string[]
  isAuthenticated: boolean
  isLoading: boolean

  // Actions
  setUser: (user: UserPublic | null) => void
  setPermissions: (permissions: string[]) => void
  setLoading: (loading: boolean) => void
  login: (user: UserPublic, accessToken: string, refreshToken: string) => void
  logout: () => void
  hasPermission: (permission: string) => boolean
  hasAnyPermission: (permissions: string[]) => boolean
  hasAllPermissions: (permissions: string[]) => boolean
}

/**
 * Check if a user permission matches a required permission
 * Supports wildcard matching:
 * - "customers.*" matches "customers.read", "customers.write", etc.
 * - "*" matches everything (super admin)
 */
function matchPermission(userPermission: string, requiredPermission: string): boolean {
  // Exact match
  if (userPermission === requiredPermission) return true

  // Super admin wildcard
  if (userPermission === '*') return true

  // Module wildcard (e.g., "customers.*" matches "customers.read")
  if (userPermission.endsWith('.*')) {
    const module = userPermission.slice(0, -2)
    return requiredPermission.startsWith(module + '.')
  }

  return false
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // Initial state
      user: null,
      permissions: [],
      isAuthenticated: false,
      isLoading: false, // Start as false - will show login if not authenticated

      // Actions
      setUser: (user) =>
        set({
          user,
          isAuthenticated: !!user,
        }),

      setPermissions: (permissions) => set({ permissions }),

      setLoading: (isLoading) => set({ isLoading }),

      login: (user, accessToken, refreshToken) => {
        // Store tokens
        tokenStorage.setAccessToken(accessToken)
        tokenStorage.setRefreshToken(refreshToken)

        // Parse permissions from JWT if available
        let permissions: string[] = []
        try {
          const payload = accessToken.split('.')[1]
          const decoded = JSON.parse(atob(payload)) as JWTClaims
          permissions = decoded.permissions || []
        } catch {
          // If JWT parsing fails, use empty permissions
          console.warn('Failed to parse JWT claims')
        }

        set({
          user,
          permissions,
          isAuthenticated: true,
          isLoading: false,
        })
      },

      logout: () => {
        tokenStorage.clearAll()
        set({
          user: null,
          permissions: [],
          isAuthenticated: false,
          isLoading: false,
        })
      },

      hasPermission: (permission) => {
        const { permissions } = get()
        return permissions.some((p) => matchPermission(p, permission))
      },

      hasAnyPermission: (requiredPermissions) => {
        const { permissions } = get()
        return requiredPermissions.some((required) =>
          permissions.some((p) => matchPermission(p, required))
        )
      },

      hasAllPermissions: (requiredPermissions) => {
        const { permissions } = get()
        return requiredPermissions.every((required) =>
          permissions.some((p) => matchPermission(p, required))
        )
      },
    }),
    {
      name: 'pawnshop_auth',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        user: state.user,
        permissions: state.permissions,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
