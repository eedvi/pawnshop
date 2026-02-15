import { z } from 'zod'

export const itemFormSchema = z.object({
  // Identification
  name: z.string().min(1, 'El nombre es requerido').max(200, 'Máximo 200 caracteres'),
  description: z.string().max(1000, 'Máximo 1000 caracteres').optional(),
  brand: z.string().max(100, 'Máximo 100 caracteres').optional(),
  model: z.string().max(100, 'Máximo 100 caracteres').optional(),
  serial_number: z.string().max(100, 'Máximo 100 caracteres').optional(),
  color: z.string().max(50, 'Máximo 50 caracteres').optional(),
  condition: z.enum(['new', 'excellent', 'good', 'fair', 'poor'], {
    required_error: 'La condición es requerida',
  }),

  // Category
  category_id: z.union([z.number(), z.string()]).optional().nullable(),
  customer_id: z.number().optional().nullable(),

  // Valuation
  appraised_value: z.number().min(0, 'Debe ser mayor o igual a 0'),
  loan_value: z.number().min(0, 'Debe ser mayor o igual a 0'),
  sale_price: z.number().min(0).optional().nullable(),

  // Physical details
  weight: z.number().min(0).optional().nullable(),
  purity: z.string().max(50, 'Máximo 50 caracteres').optional(),

  // Acquisition
  acquisition_type: z.enum(['pawn', 'purchase', 'consignment'], {
    required_error: 'El tipo de adquisición es requerido',
  }),
  acquisition_price: z.number().min(0).optional().nullable(),

  // Notes
  notes: z.string().max(1000, 'Máximo 1000 caracteres').optional(),
  tags: z.array(z.string()).optional(),
})

export type ItemFormValues = z.infer<typeof itemFormSchema>

export const defaultItemValues: ItemFormValues = {
  name: '',
  description: '',
  brand: '',
  model: '',
  serial_number: '',
  color: '',
  condition: 'good',
  category_id: 'none',
  customer_id: null,
  appraised_value: 0,
  loan_value: 0,
  sale_price: null,
  weight: null,
  purity: '',
  acquisition_type: 'pawn',
  acquisition_price: null,
  notes: '',
  tags: [],
}

export const markForSaleSchema = z.object({
  sale_price: z.number().min(1, 'El precio debe ser mayor a 0'),
})

export type MarkForSaleFormValues = z.infer<typeof markForSaleSchema>
