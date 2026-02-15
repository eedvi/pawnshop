import { apiGet, apiGetPaginated, apiPost, apiDelete } from '@/lib/api-client'
import type {
  ItemTransfer,
  TransferListParams,
  CreateTransferInput,
  ApproveTransferInput,
  ShipTransferInput,
  ReceiveTransferInput,
  CancelTransferInput,
} from '@/types'

export const transferService = {
  list: (params?: TransferListParams) =>
    apiGetPaginated<ItemTransfer>('/transfers', params),

  getById: (id: number) =>
    apiGet<ItemTransfer>(`/transfers/${id}`),

  create: (input: CreateTransferInput) =>
    apiPost<ItemTransfer>('/transfers', input),

  approve: (id: number, input?: ApproveTransferInput) =>
    apiPost<ItemTransfer>(`/transfers/${id}/approve`, input || {}),

  ship: (id: number, input?: ShipTransferInput) =>
    apiPost<ItemTransfer>(`/transfers/${id}/ship`, input || {}),

  receive: (id: number, input?: ReceiveTransferInput) =>
    apiPost<ItemTransfer>(`/transfers/${id}/receive`, input || {}),

  cancel: (id: number, input: CancelTransferInput) =>
    apiPost<ItemTransfer>(`/transfers/${id}/cancel`, input),

  delete: (id: number) =>
    apiDelete<void>(`/transfers/${id}`),
}
