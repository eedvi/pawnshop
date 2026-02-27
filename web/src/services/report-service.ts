import { apiGet, apiDownload } from '@/lib/api-client'
import type {
  DashboardStats,
  LoanReport,
  PaymentReport,
  SalesReport,
  OverdueReport,
  ReportFilters,
} from '@/types'

// Build query string from filters
function buildQueryString(filters?: ReportFilters): string {
  if (!filters) return ''
  const params = new URLSearchParams()
  if (filters.branch_id) params.append('branch_id', filters.branch_id.toString())
  if (filters.date_from) params.append('date_from', filters.date_from)
  if (filters.date_to) params.append('date_to', filters.date_to)
  const query = params.toString()
  return query ? `?${query}` : ''
}

export const reportService = {
  // Dashboard
  getDashboardStats: (branchId?: number) =>
    apiGet<DashboardStats>('/reports/dashboard', { branch_id: branchId }),

  // Reports (JSON data)
  getLoanReport: (filters?: ReportFilters) =>
    apiGet<LoanReport>('/reports/loans', filters),

  getPaymentReport: (filters?: ReportFilters) =>
    apiGet<PaymentReport>('/reports/payments', filters),

  getSalesReport: (filters?: ReportFilters) =>
    apiGet<SalesReport>('/reports/sales', filters),

  getOverdueReport: (filters?: ReportFilters) =>
    apiGet<OverdueReport>('/reports/overdue', filters),

  // Report exports (PDF)
  exportLoanReport: (filters?: ReportFilters) => {
    const dateFrom = filters?.date_from || new Date().toISOString().slice(0, 10)
    const dateTo = filters?.date_to || new Date().toISOString().slice(0, 10)
    return apiDownload(
      `/reports/loans/export${buildQueryString(filters)}`,
      `reporte_prestamos_${dateFrom}_${dateTo}.pdf`
    )
  },

  exportPaymentReport: (filters?: ReportFilters) => {
    const dateFrom = filters?.date_from || new Date().toISOString().slice(0, 10)
    const dateTo = filters?.date_to || new Date().toISOString().slice(0, 10)
    return apiDownload(
      `/reports/payments/export${buildQueryString(filters)}`,
      `reporte_pagos_${dateFrom}_${dateTo}.pdf`
    )
  },

  exportSalesReport: (filters?: ReportFilters) => {
    const dateFrom = filters?.date_from || new Date().toISOString().slice(0, 10)
    const dateTo = filters?.date_to || new Date().toISOString().slice(0, 10)
    return apiDownload(
      `/reports/sales/export${buildQueryString(filters)}`,
      `reporte_ventas_${dateFrom}_${dateTo}.pdf`
    )
  },

  exportOverdueReport: (filters?: ReportFilters) => {
    const today = new Date().toISOString().slice(0, 10)
    return apiDownload(
      `/reports/overdue/export${buildQueryString(filters)}`,
      `reporte_vencidos_${today}.pdf`
    )
  },

  // Individual document exports (PDF)
  exportLoanContract: (loanId: number) =>
    apiDownload(`/reports/export/loan/${loanId}/contract`, `contrato_prestamo_${loanId}.pdf`),

  exportPaymentReceipt: (paymentId: number) =>
    apiDownload(`/reports/export/payment/${paymentId}/receipt`, `recibo_pago_${paymentId}.pdf`),

  exportSaleReceipt: (saleId: number) =>
    apiDownload(`/reports/export/sale/${saleId}/receipt`, `recibo_venta_${saleId}.pdf`),

  exportDailyReport: (date?: string, branchId?: number) => {
    const reportDate = date || new Date().toISOString().slice(0, 10)
    const params = new URLSearchParams({ date: reportDate })
    if (branchId) params.append('branch_id', branchId.toString())
    return apiDownload(`/reports/export/daily?${params}`, `reporte_diario_${reportDate}.pdf`)
  },
}
