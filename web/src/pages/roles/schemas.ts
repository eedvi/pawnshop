import { z } from 'zod'

export const roleFormSchema = z.object({
  name: z.string().min(1, 'El nombre es requerido').max(50, 'Máximo 50 caracteres'),
  display_name: z.string().min(1, 'El nombre para mostrar es requerido').max(100, 'Máximo 100 caracteres'),
  description: z.string().max(500, 'Máximo 500 caracteres').optional(),
  permissions: z.array(z.string()),
})

export type RoleFormValues = z.infer<typeof roleFormSchema>

export const defaultRoleValues: Partial<RoleFormValues> = {
  name: '',
  display_name: '',
  description: '',
  permissions: [],
}
