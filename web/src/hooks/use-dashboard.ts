import { useQuery } from '@tanstack/react-query'
import { dashboardService } from '@/services/dashboard-service'
import { useBranchStore } from '@/stores/branch-store'

// Query keys
export const dashboardKeys = {
  all: ['dashboard'] as const,
  stats: (branchId?: number) => [...dashboardKeys.all, 'stats', branchId] as const,
  recentLoans: (branchId?: number) => [...dashboardKeys.all, 'recentLoans', branchId] as const,
  recentPayments: (branchId?: number) => [...dashboardKeys.all, 'recentPayments', branchId] as const,
  overdueCount: (branchId?: number) => [...dashboardKeys.all, 'overdueCount', branchId] as const,
}

// Hook to get dashboard statistics
export function useDashboardStats() {
  const branchId = useBranchStore((state) => state.selectedBranchId)

  return useQuery({
    queryKey: dashboardKeys.stats(branchId ?? undefined),
    queryFn: () => dashboardService.getStats(branchId ?? undefined),
    staleTime: 1000 * 60 * 5, // 5 minutes
    refetchInterval: 1000 * 60 * 5, // Auto-refresh every 5 minutes
  })
}

// Hook to get recent loans
export function useRecentLoans(limit = 5) {
  const branchId = useBranchStore((state) => state.selectedBranchId)

  return useQuery({
    queryKey: dashboardKeys.recentLoans(branchId ?? undefined),
    queryFn: () => dashboardService.getRecentLoans(branchId ?? undefined, limit),
    staleTime: 1000 * 60 * 2, // 2 minutes
  })
}

// Hook to get recent payments
export function useRecentPayments(limit = 5) {
  const branchId = useBranchStore((state) => state.selectedBranchId)

  return useQuery({
    queryKey: dashboardKeys.recentPayments(branchId ?? undefined),
    queryFn: () => dashboardService.getRecentPayments(branchId ?? undefined, limit),
    staleTime: 1000 * 60 * 2, // 2 minutes
  })
}

// Hook to get overdue loans count (for alerts)
export function useOverdueCount() {
  const branchId = useBranchStore((state) => state.selectedBranchId)

  return useQuery({
    queryKey: dashboardKeys.overdueCount(branchId ?? undefined),
    queryFn: () => dashboardService.getOverdueCount(branchId ?? undefined),
    staleTime: 1000 * 60 * 5, // 5 minutes
    refetchInterval: 1000 * 60 * 10, // Auto-refresh every 10 minutes
  })
}
