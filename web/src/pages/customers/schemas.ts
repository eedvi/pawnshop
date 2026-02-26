import { z } from 'zod'

export const customerFormSchema = z.object({
  // Personal info
  first_name: z.string().min(1, 'El nombre es requerido').max(100, 'Máximo 100 caracteres'),
  last_name: z.string().min(1, 'El apellido es requerido').max(100, 'Máximo 100 caracteres'),
  identity_type: z.enum(['dpi', 'passport', 'other'], {
    required_error: 'El tipo de documento es requerido',
  }),
  identity_number: z.string().min(1, 'El número de documento es requerido').max(50, 'Máximo 50 caracteres'),
  birth_date: z.string().optional().refine((val) => {
    if (!val) return true // optional field
    const birthDate = new Date(val)
    const today = new Date()
    const age = today.getFullYear() - birthDate.getFullYear()
    const monthDiff = today.getMonth() - birthDate.getMonth()
    const dayDiff = today.getDate() - birthDate.getDate()

    // Adjust age if birthday hasn't occurred yet this year
    const adjustedAge = monthDiff < 0 || (monthDiff === 0 && dayDiff < 0) ? age - 1 : age

    return adjustedAge >= 18
  }, {
    message: 'El cliente debe tener al menos 18 años de edad'
  }),
  gender: z.enum(['male', 'female', 'other', '__none__']).optional().transform(val => val === '__none__' ? undefined : val),

  // Contact info
  phone: z
    .string()
    .min(1, 'El teléfono es requerido')
    .min(8, 'Mínimo 8 dígitos'),
  phone_secondary: z
    .string()
    .min(8, 'Mínimo 8 dígitos')
    .optional()
    .or(z.literal('')),
  email: z.string().email('Email inválido').optional().or(z.literal('')),
  address: z.string().max(255, 'Máximo 255 caracteres').optional(),
  city: z.string().max(100, 'Máximo 100 caracteres').optional(),
  state: z.string().max(100, 'Máximo 100 caracteres').optional(),
  postal_code: z.string().max(20, 'Máximo 20 caracteres').optional(),

  // Emergency contact
  emergency_contact_name: z.string().max(200, 'Máximo 200 caracteres').optional(),
  emergency_contact_phone: z
    .string()
    .optional()
    .or(z.literal('')),
  emergency_contact_relation: z.string().max(50, 'Máximo 50 caracteres').optional(),

  // Business info
  occupation: z.string().max(100, 'Máximo 100 caracteres').optional(),
  workplace: z.string().max(200, 'Máximo 200 caracteres').optional(),
  monthly_income: z.number().min(0).optional().nullable(),

  // Credit info
  credit_limit: z.number().min(0, 'Debe ser mayor o igual a 0').default(5000),

  // Notes
  notes: z.string().max(1000, 'Máximo 1000 caracteres').optional(),

  // Status (for editing)
  is_active: z.boolean().default(true),
})

export type CustomerFormValues = z.infer<typeof customerFormSchema>

export const defaultCustomerValues: CustomerFormValues = {
  first_name: '',
  last_name: '',
  identity_type: 'dpi',
  identity_number: '',
  birth_date: '',
  gender: '__none__',
  phone: '',
  phone_secondary: '',
  email: '',
  address: '',
  city: '',
  state: '',
  postal_code: '',
  emergency_contact_name: '',
  emergency_contact_phone: '',
  emergency_contact_relation: '',
  occupation: '',
  workplace: '',
  monthly_income: null,
  credit_limit: 5000,
  notes: '',
  is_active: true,
}
