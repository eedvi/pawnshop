import { apiGet, apiGetPaginated, apiPost, apiDelete } from '@/lib/api-client'
import type {
  Sale,
  CreateSaleInput,
  RefundSaleInput,
  SaleListParams,
  SalesSummary,
} from '@/types'

export const saleService = {
  list: (params?: SaleListParams) =>
    apiGetPaginated<Sale>('/sales', params),

  getById: (id: number) =>
    apiGet<Sale>(`/sales/${id}`),

  create: (input: CreateSaleInput) =>
    apiPost<Sale>('/sales', input),

  refund: (id: number, input: RefundSaleInput) =>
    apiPost<Sale>(`/sales/${id}/refund`, input),

  cancel: (id: number, reason?: string) =>
    apiPost<Sale>(`/sales/${id}/cancel`, { reason }),

  delete: (id: number) =>
    apiDelete<void>(`/sales/${id}`),

  getSummary: (params?: { branch_id?: number; date_from?: string; date_to?: string }) =>
    apiGet<SalesSummary>('/sales/summary', params),
}
