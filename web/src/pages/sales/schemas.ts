import { z } from 'zod'

export const saleFormSchema = z.object({
  item_id: z.number().min(1, 'El artículo es requerido'),
  customer_id: z.number().optional(),
  sale_type: z.enum(['direct', 'layaway'], {
    required_error: 'El tipo de venta es requerido',
  }),
  sale_price: z.number().min(0.01, 'El precio debe ser mayor a 0'),
  discount_amount: z.number().min(0, 'El descuento no puede ser negativo').default(0),
  discount_reason: z.string().max(200, 'Máximo 200 caracteres').optional(),
  payment_method: z.enum(['cash', 'card', 'transfer', 'check', 'other'], {
    required_error: 'El método de pago es requerido',
  }),
  reference_number: z.string().max(100, 'Máximo 100 caracteres').optional(),
  notes: z.string().max(1000, 'Máximo 1000 caracteres').optional(),
})

export type SaleFormValues = z.infer<typeof saleFormSchema>

export const defaultSaleValues: Partial<SaleFormValues> = {
  item_id: undefined,
  customer_id: undefined,
  sale_type: 'direct',
  sale_price: 0,
  discount_amount: 0,
  discount_reason: '',
  payment_method: 'cash',
  reference_number: '',
  notes: '',
}

export const refundSaleSchema = z.object({
  amount: z.number().min(0.01, 'El monto debe ser mayor a 0').optional(),
  reason: z.string().min(1, 'La razón es requerida').max(500, 'Máximo 500 caracteres'),
  full_refund: z.boolean().default(true),
})

export type RefundSaleFormValues = z.infer<typeof refundSaleSchema>
