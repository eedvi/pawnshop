import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { Role } from '@/types'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Separator } from '@/components/ui/separator'
import { roleFormSchema, RoleFormValues, defaultRoleValues } from './schemas'
import { PermissionEditor } from './permission-editor'

interface RoleFormProps {
  role?: Role
  onSubmit: (values: RoleFormValues) => void
  onCancel: () => void
  isLoading?: boolean
}

export function RoleForm({ role, onSubmit, onCancel, isLoading }: RoleFormProps) {
  const form = useForm<RoleFormValues>({
    resolver: zodResolver(roleFormSchema),
    defaultValues: role
      ? {
          name: role.name,
          display_name: role.display_name,
          description: role.description || '',
          permissions: role.permissions,
        }
      : defaultRoleValues,
  })

  const isEditing = !!role

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <div className="grid gap-4 sm:grid-cols-2">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Código</FormLabel>
                <FormControl>
                  <Input
                    placeholder="admin_ventas"
                    {...field}
                    disabled={isEditing}
                    className="font-mono"
                  />
                </FormControl>
                <FormDescription>
                  Identificador único (sin espacios)
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="display_name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Nombre para mostrar</FormLabel>
                <FormControl>
                  <Input placeholder="Administrador de Ventas" {...field} />
                </FormControl>
                <FormDescription>
                  Nombre visible en la interfaz
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>

        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Descripción</FormLabel>
              <FormControl>
                <Textarea
                  placeholder="Describe las responsabilidades de este rol..."
                  {...field}
                  rows={3}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Separator />

        <FormField
          control={form.control}
          name="permissions"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="text-base">Permisos</FormLabel>
              <FormDescription>
                Selecciona los permisos que tendrá este rol
              </FormDescription>
              <FormControl>
                <PermissionEditor
                  value={field.value}
                  onChange={field.onChange}
                  disabled={isLoading}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="flex gap-2 justify-end">
          <Button type="button" variant="outline" onClick={onCancel} disabled={isLoading}>
            Cancelar
          </Button>
          <Button type="submit" disabled={isLoading}>
            {isLoading ? 'Guardando...' : isEditing ? 'Guardar Cambios' : 'Crear Rol'}
          </Button>
        </div>
      </form>
    </Form>
  )
}
