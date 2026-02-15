import { z } from 'zod'

export const categoryFormSchema = z.object({
  name: z.string().min(1, 'El nombre es requerido').max(100, 'Máximo 100 caracteres'),
  description: z.string().max(500, 'Máximo 500 caracteres').optional(),
  parent_id: z.union([z.number(), z.string()]).nullable().optional(),
  icon: z.string().max(50, 'Máximo 50 caracteres').optional(),
  default_interest_rate: z
    .number()
    .min(0, 'Debe ser mayor o igual a 0')
    .max(100, 'Máximo 100%')
    .default(10),
  min_loan_amount: z.number().min(0).optional().nullable(),
  max_loan_amount: z.number().min(0).optional().nullable(),
  loan_to_value_ratio: z
    .number()
    .min(0, 'Debe ser mayor o igual a 0')
    .max(100, 'Máximo 100%')
    .default(70),
  sort_order: z.number().int().min(0).default(0),
  is_active: z.boolean().default(true),
})

export type CategoryFormValues = z.infer<typeof categoryFormSchema>

export const defaultCategoryValues: CategoryFormValues = {
  name: '',
  description: '',
  parent_id: 'none',
  icon: '',
  default_interest_rate: 10,
  min_loan_amount: null,
  max_loan_amount: null,
  loan_to_value_ratio: 70,
  sort_order: 0,
  is_active: true,
}
