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

export function useExportLoanReport() {
  return useMutation({
    mutationFn: (filters?: ReportFilters) => reportService.exportLoanReport(filters),
    onSuccess: (blob) => {
      downloadBlob(blob, 'reporte-prestamos.pdf')
    },
  })
}

export function useExportPaymentReport() {
  return useMutation({
    mutationFn: (filters?: ReportFilters) => reportService.exportPaymentReport(filters),
    onSuccess: (blob) => {
      downloadBlob(blob, 'reporte-pagos.pdf')
    },
  })
}

export function useExportSalesReport() {
  return useMutation({
    mutationFn: (filters?: ReportFilters) => reportService.exportSalesReport(filters),
    onSuccess: (blob) => {
      downloadBlob(blob, 'reporte-ventas.pdf')
    },
  })
}

export function useExportOverdueReport() {
  return useMutation({
    mutationFn: (filters?: ReportFilters) => reportService.exportOverdueReport(filters),
    onSuccess: (blob) => {
      downloadBlob(blob, 'reporte-vencidos.pdf')
    },
  })
}

function downloadBlob(blob: Blob, filename: string) {
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  window.URL.revokeObjectURL(url)
}
