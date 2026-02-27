import { useQuery, useMutation } from '@tanstack/react-query'
import { reportService } from '@/services/report-service'
import type { ReportFilters } from '@/types'

export const reportKeys = {
  all: ['reports'] as const,
  dashboard: (branchId?: number) => [...reportKeys.all, 'dashboard', branchId] as const,
  loans: (filters?: ReportFilters) => [...reportKeys.all, 'loans', filters] as const,
  payments: (filters?: ReportFilters) => [...reportKeys.all, 'payments', filters] as const,
  sales: (filters?: ReportFilters) => [...reportKeys.all, 'sales', filters] as const,
  overdue: (filters?: ReportFilters) => [...reportKeys.all, 'overdue', filters] as const,
}

export function useDashboardStats(branchId?: number) {
  return useQuery({
    queryKey: reportKeys.dashboard(branchId),
    queryFn: () => reportService.getDashboardStats(branchId),
  })
}

export function useLoanReport(filters?: ReportFilters) {
  return useQuery({
    queryKey: reportKeys.loans(filters),
    queryFn: () => reportService.getLoanReport(filters),
  })
}

export function usePaymentReport(filters?: ReportFilters) {
  return useQuery({
    queryKey: reportKeys.payments(filters),
    queryFn: () => reportService.getPaymentReport(filters),
  })
}

export function useSalesReport(filters?: ReportFilters) {
  return useQuery({
    queryKey: reportKeys.sales(filters),
    queryFn: () => reportService.getSalesReport(filters),
  })
}

export function useOverdueReport(filters?: ReportFilters) {
  return useQuery({
    queryKey: reportKeys.overdue(filters),
    queryFn: () => reportService.getOverdueReport(filters),
  })
}

// Report export mutations - apiDownload handles the download automatically
export function useExportLoanReport() {
  return useMutation({
    mutationFn: (filters?: ReportFilters) => reportService.exportLoanReport(filters),
  })
}

export function useExportPaymentReport() {
  return useMutation({
    mutationFn: (filters?: ReportFilters) => reportService.exportPaymentReport(filters),
  })
}

export function useExportSalesReport() {
  return useMutation({
    mutationFn: (filters?: ReportFilters) => reportService.exportSalesReport(filters),
  })
}

export function useExportOverdueReport() {
  return useMutation({
    mutationFn: (filters?: ReportFilters) => reportService.exportOverdueReport(filters),
  })
}

// Individual document export mutations
export function useExportLoanContract() {
  return useMutation({
    mutationFn: (loanId: number) => reportService.exportLoanContract(loanId),
  })
}

export function useExportPaymentReceipt() {
  return useMutation({
    mutationFn: (paymentId: number) => reportService.exportPaymentReceipt(paymentId),
  })
}

export function useExportSaleReceipt() {
  return useMutation({
    mutationFn: (saleId: number) => reportService.exportSaleReceipt(saleId),
  })
}

export function useExportDailyReport() {
  return useMutation({
    mutationFn: (params?: { date?: string; branchId?: number }) =>
      reportService.exportDailyReport(params?.date, params?.branchId),
  })
}
