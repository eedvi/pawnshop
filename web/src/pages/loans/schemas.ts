import { z } from 'zod'

export const loanFormSchema = z.object({
  customer_id: z.number().min(1, 'El cliente es requerido'),
  item_id: z.number().min(1, 'El artículo es requerido'),
  loan_amount: z.number().min(1, 'El monto debe ser mayor a 0'),
  interest_rate: z.number().min(0).max(100).optional(),
  payment_plan_type: z.enum(['single', 'minimum_payment', 'installments'], {
    required_error: 'El tipo de pago es requerido',
  }),
  loan_term_days: z.number().min(1).optional(),
  grace_period_days: z.number().min(0).optional(),
  number_of_installments: z.number().min(1).optional(),
  notes: z.string().max(1000, 'Máximo 1000 caracteres').optional(),
})

export type LoanFormValues = z.infer<typeof loanFormSchema>

export const defaultLoanValues: Partial<LoanFormValues> = {
  customer_id: undefined,
  item_id: undefined,
  loan_amount: 0,
  interest_rate: undefined,
  payment_plan_type: 'single',
  loan_term_days: 30,
  grace_period_days: 0,
  number_of_installments: undefined,
  notes: '',
}

export const renewLoanSchema = z.object({
  new_term_days: z.number().min(1, 'El plazo debe ser mayor a 0'),
  new_interest_rate: z.number().min(0).max(100).optional(),
  pay_interest: z.boolean().optional(),
})

export type RenewLoanFormValues = z.infer<typeof renewLoanSchema>

export const confiscateLoanSchema = z.object({
  notes: z.string().max(1000, 'Máximo 1000 caracteres').optional(),
})

export type ConfiscateLoanFormValues = z.infer<typeof confiscateLoanSchema>
