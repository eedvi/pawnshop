// Sale types - mirrors internal/domain/sale.go

import type { Branch } from './branch'
import type { Customer } from './customer'
import type { Item } from './item'
import type { PaymentMethod } from './payment'

export type SaleStatus = 'completed' | 'pending' | 'cancelled' | 'refunded' | 'partial_refund'
export type SaleType = 'direct' | 'layaway'

export const SALE_STATUSES: { value: SaleStatus; label: string; color: string }[] = [
  { value: 'completed', label: 'Completada', color: 'green' },
  { value: 'pending', label: 'Pendiente', color: 'yellow' },
  { value: 'cancelled', label: 'Cancelada', color: 'gray' },
  { value: 'refunded', label: 'Reembolsada', color: 'red' },
  { value: 'partial_refund', label: 'Reembolso Parcial', color: 'orange' },
]

export const SALE_TYPES: { value: SaleType; label: string }[] = [
  { value: 'direct', label: 'Venta Directa' },
  { value: 'layaway', label: 'Apartado' },
]

export interface Sale {
  id: number
  branch_id: number
  item_id: number
  customer_id?: number
  sale_number: string
  sale_type: SaleType

  // Pricing
  sale_price: number
  discount_amount: number
  discount_reason?: string
  final_price: number

  // Payment
  payment_method: PaymentMethod
  reference_number?: string

  // Status
  status: SaleStatus
  sale_date: string

  // Refund info
  refund_amount?: number
  refund_reason?: string
  refunded_at?: string
  refunded_by?: number

  // Notes
  notes?: string

  // Cash session reference
  cash_session_id?: number

  // Audit
  created_by?: number
  updated_by?: number
  created_at: string
  updated_at: string
  deleted_at?: string

  // Relations
  branch?: Branch
  item?: Item
  customer?: Customer
}

export interface CreateSaleInput {
  branch_id: number
  item_id: number
  customer_id?: number
  sale_type?: SaleType
  sale_price?: number
  discount_amount?: number
  discount_reason?: string
  payment_method: PaymentMethod
  reference_number?: string
  notes?: string
}

export interface RefundSaleInput {
  amount?: number
  reason: string
}

export interface SaleListParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  search?: string
  branch_id?: number
  item_id?: number
  customer_id?: number
  status?: SaleStatus
  sale_type?: SaleType
  payment_method?: PaymentMethod
  date_from?: string
  date_to?: string
}

export interface SalesSummary {
  total_sales: number
  total_amount: number
  total_discount: number
  net_amount: number
  total_refunds: number
  refund_amount: number
  sales_by_method: {
    method: PaymentMethod
    count: number
    amount: number
  }[]
}
