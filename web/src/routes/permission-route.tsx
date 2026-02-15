import { Navigate } from 'react-router-dom'
import { useAuthStore } from '@/stores'
import { ROUTES } from './routes'

interface PermissionRouteProps {
  children: React.ReactNode
  permission?: string
  permissions?: string[]
  requireAll?: boolean
  fallback?: React.ReactNode
}

/**
 * Route guard that requires specific permissions.
 * Must be used inside a ProtectedRoute (assumes user is authenticated).
 *
 * @param permission - Single permission required
 * @param permissions - Multiple permissions (use with requireAll)
 * @param requireAll - If true, user must have ALL permissions. If false, ANY permission suffices.
 * @param fallback - Optional custom fallback. Defaults to redirect to dashboard.
 */
export function PermissionRoute({
  children,
  permission,
  permissions = [],
  requireAll = false,
  fallback,
}: PermissionRouteProps) {
  const { hasAnyPermission, hasAllPermissions } = useAuthStore()

  // Check permissions
  const allPermissions = permission ? [permission, ...permissions] : permissions

  if (allPermissions.length === 0) {
    // No permissions required
    return <>{children}</>
  }

  const hasAccess = requireAll
    ? hasAllPermissions(allPermissions)
    : hasAnyPermission(allPermissions)

  if (!hasAccess) {
    if (fallback) {
      return <>{fallback}</>
    }
    // Redirect to dashboard if no access
    return <Navigate to={ROUTES.DASHBOARD} replace />
  }

  return <>{children}</>
}

/**
 * Forbidden page component shown when user lacks permissions
 */
export function ForbiddenPage() {
  return (
    <div className="flex h-[50vh] flex-col items-center justify-center gap-4">
      <div className="text-6xl">🚫</div>
      <h1 className="text-2xl font-bold">Acceso Denegado</h1>
      <p className="text-muted-foreground">
        No tienes permisos para acceder a esta página.
      </p>
    </div>
  )
}
