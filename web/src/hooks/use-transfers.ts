import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { transferService } from '@/services/transfer-service'
import type {
  TransferListParams,
  CreateTransferInput,
  ApproveTransferInput,
  ShipTransferInput,
  ReceiveTransferInput,
  CancelTransferInput,
} from '@/types'

export const transferKeys = {
  all: ['transfers'] as const,
  lists: () => [...transferKeys.all, 'list'] as const,
  list: (params?: TransferListParams) => [...transferKeys.lists(), params] as const,
  details: () => [...transferKeys.all, 'detail'] as const,
  detail: (id: number) => [...transferKeys.details(), id] as const,
}

export function useTransfers(params?: TransferListParams) {
  return useQuery({
    queryKey: transferKeys.list(params),
    queryFn: () => transferService.list(params),
  })
}

export function useTransfer(id: number) {
  return useQuery({
    queryKey: transferKeys.detail(id),
    queryFn: () => transferService.getById(id),
    enabled: id > 0,
  })
}

export function useCreateTransfer() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateTransferInput) => transferService.create(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: transferKeys.lists() })
    },
  })
}

export function useApproveTransfer() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input?: ApproveTransferInput }) =>
      transferService.approve(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: transferKeys.lists() })
      queryClient.invalidateQueries({ queryKey: transferKeys.detail(id) })
    },
  })
}

export function useShipTransfer() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input?: ShipTransferInput }) =>
      transferService.ship(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: transferKeys.lists() })
      queryClient.invalidateQueries({ queryKey: transferKeys.detail(id) })
    },
  })
}

export function useReceiveTransfer() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input?: ReceiveTransferInput }) =>
      transferService.receive(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: transferKeys.lists() })
      queryClient.invalidateQueries({ queryKey: transferKeys.detail(id) })
    },
  })
}

export function useCancelTransfer() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: CancelTransferInput }) =>
      transferService.cancel(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: transferKeys.lists() })
      queryClient.invalidateQueries({ queryKey: transferKeys.detail(id) })
    },
  })
}

export function useDeleteTransfer() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => transferService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: transferKeys.lists() })
    },
  })
}
