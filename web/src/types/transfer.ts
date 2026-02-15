// Transfer types - mirrors internal/domain/transfer.go

import type { Branch } from './branch'
import type { Item } from './item'
import type { User } from './user'

export type TransferStatus =
  | 'pending'
  | 'in_transit'
  | 'completed'
  | 'cancelled'

export const TRANSFER_STATUSES: { value: TransferStatus; label: string; color: string }[] = [
  { value: 'pending', label: 'Pendiente', color: 'yellow' },
  { value: 'in_transit', label: 'En Tránsito', color: 'purple' },
  { value: 'completed', label: 'Completada', color: 'green' },
  { value: 'cancelled', label: 'Cancelada', color: 'gray' },
]

export interface ItemTransfer {
  id: number
  transfer_number: string
  item_id: number
  from_branch_id: number
  to_branch_id: number
  status: TransferStatus

  // Users involved
  requested_by: number
  approved_by?: number
  received_by?: number

  // Dates
  requested_at: string
  approved_at?: string
  shipped_at?: string
  received_at?: string
  cancelled_at?: string

  // Notes
  request_notes?: string
  approval_notes?: string
  receipt_notes?: string
  cancellation_reason?: string

  // Timestamps
  created_at: string
  updated_at: string

  // Relations
  item?: Item
  from_branch?: Branch
  to_branch?: Branch
  requester?: User
  approver?: User
  receiver?: User
}

export interface CreateTransferInput {
  item_id: number
  to_branch_id: number
  request_notes?: string
}

export interface ApproveTransferInput {
  approval_notes?: string
}

export interface ShipTransferInput {
  // No specific fields - just action
}

export interface ReceiveTransferInput {
  receipt_notes?: string
}

export interface CancelTransferInput {
  cancellation_reason: string
}

export interface TransferListParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  search?: string
  item_id?: number
  from_branch_id?: number
  to_branch_id?: number
  status?: TransferStatus
  date_from?: string
  date_to?: string
}
