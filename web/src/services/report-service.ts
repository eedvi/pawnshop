import { apiGet } from '@/lib/api-client'
import type {
  DashboardStats,
  LoanReport,
  PaymentReport,
  SalesReport,
  OverdueReport,
  ReportFilters,
} from '@/types'

export const reportService = {
  getDashboardStats: (branchId?: number) =>
    apiGet<DashboardStats>('/reports/dashboard', { branch_id: branchId }),

  getLoanReport: (filters?: ReportFilters) =>
    apiGet<LoanReport>('/reports/loans', filters),

  getPaymentReport: (filters?: ReportFilters) =>
    apiGet<PaymentReport>('/reports/payments', filters),

  getSalesReport: (filters?: ReportFilters) =>
    apiGet<SalesReport>('/reports/sales', filters),

  getOverdueReport: (filters?: ReportFilters) =>
    apiGet<OverdueReport>('/reports/overdue', filters),

  exportLoanReport: (filters?: ReportFilters) =>
    apiGet<Blob>('/reports/loans/export', { ...filters, responseType: 'blob' }),

  exportPaymentReport: (filters?: ReportFilters) =>
    apiGet<Blob>('/reports/payments/export', { ...filters, responseType: 'blob' }),

  exportSalesReport: (filters?: ReportFilters) =>
    apiGet<Blob>('/reports/sales/export', { ...filters, responseType: 'blob' }),

  exportOverdueReport: (filters?: ReportFilters) =>
    apiGet<Blob>('/reports/overdue/export', { ...filters, responseType: 'blob' }),
}
