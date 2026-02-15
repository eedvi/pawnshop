import { useAuthStore } from '@/stores'

interface PermissionGuardProps {
  children: React.ReactNode
  permission?: string
  permissions?: string[]
  requireAll?: boolean
  fallback?: React.ReactNode
}

/**
 * Component that conditionally renders children based on user permissions.
 *
 * @param permission - Single permission required
 * @param permissions - Multiple permissions (use with requireAll)
 * @param requireAll - If true, user must have ALL permissions. If false, ANY permission suffices.
 * @param fallback - Optional fallback to render when user lacks permissions.
 */
export function PermissionGuard({
  children,
  permission,
  permissions = [],
  requireAll = false,
  fallback = null,
}: PermissionGuardProps) {
  const { hasAnyPermission, hasAllPermissions } = useAuthStore()

  const allPermissions = permission ? [permission, ...permissions] : permissions

  if (allPermissions.length === 0) {
    return <>{children}</>
  }

  const hasAccess = requireAll
    ? hasAllPermissions(allPermissions)
    : hasAnyPermission(allPermissions)

  if (!hasAccess) {
    return <>{fallback}</>
  }

  return <>{children}</>
}
