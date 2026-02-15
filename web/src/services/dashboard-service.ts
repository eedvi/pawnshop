import { apiGet } from '@/lib/api-client'
import { DashboardStats, LoanReportItem, PaymentReportItem } from '@/types'

export interface RecentActivityItem {
  type: 'loan' | 'payment' | 'sale'
  id: number
  description: string
  amount: number
  customer_name: string
  created_at: string
}

export interface DashboardData {
  stats: DashboardStats
  recent_loans: LoanReportItem[]
  recent_payments: PaymentReportItem[]
  recent_activity: RecentActivityItem[]
}

// Dashboard endpoints
export const dashboardService = {
  // Get dashboard statistics
  getStats: async (branchId?: number): Promise<DashboardStats> => {
    const params: Record<string, unknown> = {}
    if (branchId) {
      params.branch_id = branchId
    }
    return apiGet<DashboardStats>('/reports/dashboard', params)
  },

  // Get recent loans (for chart/list)
  getRecentLoans: async (branchId?: number, limit = 5): Promise<LoanReportItem[]> => {
    const params: Record<string, unknown> = { limit }
    if (branchId) {
      params.branch_id = branchId
    }
    // Use loan report endpoint with small limit
    const report = await apiGet<{ items: LoanReportItem[] }>('/reports/loans', params)
    return report.items?.slice(0, limit) ?? []
  },

  // Get recent payments (for chart/list)
  getRecentPayments: async (branchId?: number, limit = 5): Promise<PaymentReportItem[]> => {
    const params: Record<string, unknown> = { limit }
    if (branchId) {
      params.branch_id = branchId
    }
    const report = await apiGet<{ items: PaymentReportItem[] }>('/reports/payments', params)
    return report.items?.slice(0, limit) ?? []
  },

  // Get overdue loans count (for alerts)
  getOverdueCount: async (branchId?: number): Promise<number> => {
    const params: Record<string, unknown> = {}
    if (branchId) {
      params.branch_id = branchId
    }
    const report = await apiGet<{ summary: { total_overdue_loans: number } }>('/reports/overdue', params)
    return report.summary?.total_overdue_loans ?? 0
  },
}
