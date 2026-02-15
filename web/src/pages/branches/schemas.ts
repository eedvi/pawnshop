import { z } from 'zod'

export const branchFormSchema = z.object({
  name: z.string().min(1, 'El nombre es requerido').max(100, 'Máximo 100 caracteres'),
  code: z
    .string()
    .min(1, 'El código es requerido')
    .max(10, 'Máximo 10 caracteres')
    .regex(/^[A-Z0-9]+$/, 'Solo letras mayúsculas y números'),
  address: z.string().max(255, 'Máximo 255 caracteres').optional(),
  phone: z
    .string()
    .regex(/^\d{4}-?\d{4}$/, 'Formato: XXXX-XXXX')
    .optional()
    .or(z.literal('')),
  email: z.string().email('Email inválido').optional().or(z.literal('')),
  timezone: z.string().default('America/Guatemala'),
  currency: z.string().default('GTQ'),
  default_interest_rate: z
    .number()
    .min(0, 'Debe ser mayor o igual a 0')
    .max(100, 'Máximo 100%')
    .default(10),
  default_loan_term_days: z
    .number()
    .int()
    .min(1, 'Mínimo 1 día')
    .max(365, 'Máximo 365 días')
    .default(30),
  default_grace_period: z
    .number()
    .int()
    .min(0, 'Debe ser mayor o igual a 0')
    .max(30, 'Máximo 30 días')
    .default(3),
  is_active: z.boolean().default(true),
})

export type BranchFormValues = z.infer<typeof branchFormSchema>

export const defaultBranchValues: BranchFormValues = {
  name: '',
  code: '',
  address: '',
  phone: '',
  email: '',
  timezone: 'America/Guatemala',
  currency: 'GTQ',
  default_interest_rate: 10,
  default_loan_term_days: 30,
  default_grace_period: 3,
  is_active: true,
}
