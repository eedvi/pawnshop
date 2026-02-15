import { Navigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/stores'
import { ROUTES } from './routes'

interface ProtectedRouteProps {
  children: React.ReactNode
}

/**
 * Route guard that requires authentication.
 * Redirects to login page if user is not authenticated.
 * Preserves the intended destination in location state.
 */
export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading } = useAuthStore()
  const location = useLocation()

  // Show nothing while checking auth state
  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    )
  }

  if (!isAuthenticated) {
    // Redirect to login, preserving the intended destination
    return <Navigate to={ROUTES.LOGIN} state={{ from: location }} replace />
  }

  return <>{children}</>
}
