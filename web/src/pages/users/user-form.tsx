import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2 } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Form } from '@/components/ui/form'
import { FormInput, FormSelect } from '@/components/form'
import { useBranches } from '@/hooks/use-branches'
import { useRoles } from '@/hooks/use-roles'
import { User } from '@/types'
import {
  userFormSchema,
  createUserSchema,
  UserFormValues,
  CreateUserFormValues,
  defaultUserValues,
} from './schemas'

interface UserFormProps {
  user?: User
  onSubmit: (values: UserFormValues | CreateUserFormValues) => void
  onCancel: () => void
  isLoading?: boolean
}

export function UserForm({ user, onSubmit, onCancel, isLoading = false }: UserFormProps) {
  const isEditing = !!user

  const { data: branches } = useBranches()
  const { data: roles } = useRoles()

  const form = useForm<UserFormValues | CreateUserFormValues>({
    resolver: zodResolver(isEditing ? userFormSchema : createUserSchema),
    defaultValues: user
      ? {
          email: user.email,
          first_name: user.first_name,
          last_name: user.last_name,
          phone: user.phone || '',
          branch_id: user.branch_id ?? 'all',
          role_id: user.role_id,
        }
      : defaultUserValues,
  })

  const branchOptions = [
    { value: 'all', label: 'Todas las sucursales' },
    ...(branches?.data?.map((b) => ({
      value: b.id.toString(),
      label: b.name,
    })) || []),
  ]

  const roleOptions =
    roles?.data?.map((r) => ({
      value: r.id.toString(),
      label: r.name,
    })) || []

  const handleSubmit = (values: UserFormValues | CreateUserFormValues) => {
    onSubmit({
      ...values,
      branch_id: values.branch_id === 'all' ? undefined : values.branch_id,
    })
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-6">
        <div className="grid gap-4 sm:grid-cols-2">
          <FormInput
            control={form.control}
            name="first_name"
            label="Nombre"
            required
          />
          <FormInput
            control={form.control}
            name="last_name"
            label="Apellido"
            required
          />
        </div>

        <FormInput
          control={form.control}
          name="email"
          label="Email"
          type="email"
          required
        />

        {!isEditing && (
          <FormInput
            control={form.control}
            name="password"
            label="Contraseña"
            type="password"
            required
            description="Mínimo 8 caracteres"
          />
        )}

        <FormInput
          control={form.control}
          name="phone"
          label="Teléfono"
          placeholder="Opcional"
        />

        <div className="grid gap-4 sm:grid-cols-2">
          <FormSelect
            control={form.control}
            name="role_id"
            label="Rol"
            options={roleOptions}
            required
          />
          <FormSelect
            control={form.control}
            name="branch_id"
            label="Sucursal"
            options={branchOptions}
            description="Dejar vacío para acceso a todas"
          />
        </div>

        <div className="flex justify-end gap-4">
          <Button type="button" variant="outline" onClick={onCancel}>
            Cancelar
          </Button>
          <Button type="submit" disabled={isLoading}>
            {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {isEditing ? 'Guardar Cambios' : 'Crear Usuario'}
          </Button>
        </div>
      </form>
    </Form>
  )
}
