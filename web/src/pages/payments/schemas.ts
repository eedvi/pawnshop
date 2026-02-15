import { z } from 'zod'

export const paymentFormSchema = z.object({
  loan_id: z.number().min(1, 'El préstamo es requerido'),
  amount: z.number().min(1, 'El monto debe ser mayor a 0'),
  payment_method: z.enum(['cash', 'card', 'transfer', 'check', 'other'], {
    required_error: 'El método de pago es requerido',
  }),
  reference_number: z.string().max(100, 'Máximo 100 caracteres').optional(),
  notes: z.string().max(1000, 'Máximo 1000 caracteres').optional(),
})

export type PaymentFormValues = z.infer<typeof paymentFormSchema>

export const defaultPaymentValues: Partial<PaymentFormValues> = {
  loan_id: undefined,
  amount: 0,
  payment_method: 'cash',
  reference_number: '',
  notes: '',
}

export const reversePaymentSchema = z.object({
  reason: z.string().min(1, 'La razón es requerida').max(500, 'Máximo 500 caracteres'),
})

export type ReversePaymentFormValues = z.infer<typeof reversePaymentSchema>
