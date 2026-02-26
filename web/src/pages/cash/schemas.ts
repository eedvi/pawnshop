import { z } from 'zod'

export const cashRegisterFormSchema = z.object({
  name: z.string().min(1, 'El nombre es requerido').max(100, 'Máximo 100 caracteres'),
  code: z.string().min(1, 'El código es requerido').max(20, 'Máximo 20 caracteres'),
  description: z.string().max(500, 'Máximo 500 caracteres').optional(),
})

export type CashRegisterFormValues = z.infer<typeof cashRegisterFormSchema>

export const openSessionSchema = z.object({
  cash_register_id: z.number().min(1, 'La caja es requerida'),
  opening_amount: z.number().min(0, 'El monto debe ser mayor o igual a 0'),
  notes: z.string().max(500, 'Máximo 500 caracteres').optional(),
})

export type OpenSessionFormValues = z.infer<typeof openSessionSchema>

export const closeSessionSchema = z.object({
  closing_amount: z.number().min(0, 'El monto debe ser mayor o igual a 0'),
  notes: z.string().max(500, 'Máximo 500 caracteres').optional(),
})

export type CloseSessionFormValues = z.infer<typeof closeSessionSchema>

export const cashMovementSchema = z.object({
  movement_type: z.enum(['income', 'expense', 'adjustment'], {
    required_error: 'El tipo de movimiento es requerido',
  }),
  amount: z.number().min(0.01, 'El monto debe ser mayor a 0'),
  description: z.string().min(1, 'La descripción es requerida').max(500, 'Máximo 500 caracteres'),
})

export type CashMovementFormValues = z.infer<typeof cashMovementSchema>
