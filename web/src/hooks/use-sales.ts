import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { saleService } from '@/services/sale-service'
import type { SaleListParams, CreateSaleInput, RefundSaleInput } from '@/types'

export const saleKeys = {
  all: ['sales'] as const,
  lists: () => [...saleKeys.all, 'list'] as const,
  list: (params?: SaleListParams) => [...saleKeys.lists(), params] as const,
  details: () => [...saleKeys.all, 'detail'] as const,
  detail: (id: number) => [...saleKeys.details(), id] as const,
  summary: (params?: { branch_id?: number; date_from?: string; date_to?: string }) =>
    [...saleKeys.all, 'summary', params] as const,
}

export function useSales(params?: SaleListParams) {
  return useQuery({
    queryKey: saleKeys.list(params),
    queryFn: () => saleService.list(params),
  })
}

export function useSale(id: number) {
  return useQuery({
    queryKey: saleKeys.detail(id),
    queryFn: () => saleService.getById(id),
    enabled: id > 0,
  })
}

export function useSalesSummary(params?: {
  branch_id?: number
  date_from?: string
  date_to?: string
}) {
  return useQuery({
    queryKey: saleKeys.summary(params),
    queryFn: () => saleService.getSummary(params),
  })
}

export function useCreateSale() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateSaleInput) => saleService.create(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: saleKeys.lists() })
    },
  })
}

export function useRefundSale() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: RefundSaleInput }) =>
      saleService.refund(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: saleKeys.lists() })
      queryClient.invalidateQueries({ queryKey: saleKeys.detail(id) })
    },
  })
}

export function useCancelSale() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, reason }: { id: number; reason?: string }) =>
      saleService.cancel(id, reason),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: saleKeys.lists() })
      queryClient.invalidateQueries({ queryKey: saleKeys.detail(id) })
    },
  })
}

export function useDeleteSale() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => saleService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: saleKeys.lists() })
    },
  })
}
