import { z } from 'zod'

export const userFormSchema = z.object({
  email: z.string().email('Email inválido'),
  password: z.string().min(8, 'La contraseña debe tener al menos 8 caracteres').optional(),
  first_name: z.string().min(1, 'El nombre es requerido').max(100, 'Máximo 100 caracteres'),
  last_name: z.string().min(1, 'El apellido es requerido').max(100, 'Máximo 100 caracteres'),
  phone: z.string().max(20, 'Máximo 20 caracteres').optional(),
  branch_id: z.union([z.number(), z.string()]).optional(),
  role_id: z.number().min(1, 'El rol es requerido'),
})

export type UserFormValues = z.infer<typeof userFormSchema>

export const createUserSchema = userFormSchema.extend({
  password: z.string().min(8, 'La contraseña debe tener al menos 8 caracteres'),
})

export type CreateUserFormValues = z.infer<typeof createUserSchema>

export const resetPasswordSchema = z.object({
  password: z.string().min(8, 'La contraseña debe tener al menos 8 caracteres'),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: 'Las contraseñas no coinciden',
  path: ['confirmPassword'],
})

export type ResetPasswordFormValues = z.infer<typeof resetPasswordSchema>

export const defaultUserValues: Partial<UserFormValues> = {
  email: '',
  first_name: '',
  last_name: '',
  phone: '',
  branch_id: 'all',
  role_id: undefined,
}
