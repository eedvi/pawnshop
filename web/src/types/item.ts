// Item types - mirrors internal/domain/item.go

import type { Branch } from './branch'
import type { Category } from './category'
import type { Customer } from './customer'

export type ItemStatus =
  | 'available'
  | 'pawned'
  | 'collateral'
  | 'for_sale'
  | 'sold'
  | 'confiscated'
  | 'transferred'
  | 'in_transfer'
  | 'damaged'
  | 'lost'

export type ItemCondition = 'new' | 'excellent' | 'good' | 'fair' | 'poor'
export type AcquisitionType = 'pawn' | 'purchase' | 'consignment'

export const ITEM_STATUSES: { value: ItemStatus; label: string; color: string }[] = [
  { value: 'available', label: 'Disponible', color: 'green' },
  { value: 'pawned', label: 'Empeñado', color: 'blue' },
  { value: 'collateral', label: 'En Garantía', color: 'purple' },
  { value: 'for_sale', label: 'En Venta', color: 'orange' },
  { value: 'sold', label: 'Vendido', color: 'gray' },
  { value: 'confiscated', label: 'Confiscado', color: 'red' },
  { value: 'transferred', label: 'Transferido', color: 'cyan' },
  { value: 'in_transfer', label: 'En Tránsito', color: 'yellow' },
  { value: 'damaged', label: 'Dañado', color: 'red' },
  { value: 'lost', label: 'Perdido', color: 'red' },
]

export const ITEM_CONDITIONS: { value: ItemCondition; label: string }[] = [
  { value: 'new', label: 'Nuevo' },
  { value: 'excellent', label: 'Excelente' },
  { value: 'good', label: 'Bueno' },
  { value: 'fair', label: 'Regular' },
  { value: 'poor', label: 'Malo' },
]

export interface Item {
  id: number
  branch_id: number
  category_id?: number
  customer_id?: number

  // Identification
  sku: string
  name: string
  description?: string
  brand?: string
  model?: string
  serial_number?: string
  color?: string
  condition: ItemCondition

  // Valuation
  appraised_value: number
  loan_value: number
  sale_price?: number

  // Status
  status: ItemStatus

  // Physical details
  weight?: number
  purity?: string

  // Additional info
  notes?: string
  tags?: string[]

  // Acquisition
  acquisition_type: AcquisitionType
  acquisition_date: string
  acquisition_price?: number

  // Media
  photos?: string[]

  // Delivery tracking
  delivered_at?: string

  // Audit
  created_by?: number
  updated_by?: number
  created_at: string
  updated_at: string
  deleted_at?: string

  // Relations
  branch?: Branch
  category?: Category
  customer?: Customer
}

export interface CreateItemInput {
  branch_id: number
  category_id?: number
  customer_id?: number
  name: string
  description?: string
  brand?: string
  model?: string
  serial_number?: string
  color?: string
  condition: ItemCondition
  appraised_value: number
  loan_value: number
  sale_price?: number
  weight?: number
  purity?: string
  notes?: string
  tags?: string[]
  acquisition_type: AcquisitionType
  acquisition_price?: number
}

export interface UpdateItemInput {
  category_id?: number
  name?: string
  description?: string
  brand?: string
  model?: string
  serial_number?: string
  color?: string
  condition?: ItemCondition
  appraised_value?: number
  loan_value?: number
  sale_price?: number
  weight?: number
  purity?: string
  notes?: string
  tags?: string[]
}

export interface ItemListParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  search?: string
  branch_id?: number
  category_id?: number
  customer_id?: number
  status?: ItemStatus
  condition?: ItemCondition
}

export interface ItemHistory {
  id: number
  item_id: number
  action: string
  old_status?: string
  new_status?: string
  old_branch_id?: number
  new_branch_id?: number
  reference_type?: string
  reference_id?: number
  notes?: string
  created_by?: number
  created_at: string
}
