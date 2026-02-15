import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2 } from 'lucide-react'

import { Customer, IDENTITY_TYPES, GENDERS } from '@/types'
import { Button } from '@/components/ui/button'
import { Form } from '@/components/ui/form'
import { FormInput, FormSelect, FormTextarea, FormDatePicker, FormSwitch, FormCurrencyInput } from '@/components/form'
import { customerFormSchema, CustomerFormValues, defaultCustomerValues } from './schemas'

interface CustomerFormProps {
  customer?: Customer
  onSubmit: (values: CustomerFormValues) => void
  onCancel: () => void
  isLoading?: boolean
}

export function CustomerForm({ customer, onSubmit, onCancel, isLoading }: CustomerFormProps) {
  const form = useForm<CustomerFormValues>({
    resolver: zodResolver(customerFormSchema),
    defaultValues: customer
      ? {
          first_name: customer.first_name,
          last_name: customer.last_name,
          identity_type: customer.identity_type,
          identity_number: customer.identity_number,
          birth_date: customer.birth_date || '',
          gender: customer.gender || '__none__',
          phone: customer.phone,
          phone_secondary: customer.phone_secondary || '',
          email: customer.email || '',
          address: customer.address || '',
          city: customer.city || '',
          state: customer.state || '',
          postal_code: customer.postal_code || '',
          emergency_contact_name: customer.emergency_contact_name || '',
          emergency_contact_phone: customer.emergency_contact_phone || '',
          emergency_contact_relation: customer.emergency_contact_relation || '',
          occupation: customer.occupation || '',
          workplace: customer.workplace || '',
          monthly_income: customer.monthly_income ?? null,
          credit_limit: customer.credit_limit,
          notes: customer.notes || '',
          is_active: customer.is_active,
        }
      : defaultCustomerValues,
  })

  const isEditing = !!customer

  const identityTypeOptions = IDENTITY_TYPES.map((t) => ({
    value: t.value,
    label: t.label,
  }))

  const genderOptions = [
    { value: '__none__', label: 'Sin especificar' },
    ...GENDERS.map((g) => ({
      value: g.value,
      label: g.label,
    })),
  ]

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        {/* Personal Information */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Información Personal</h3>
          <div className="grid gap-4 sm:grid-cols-2">
            <FormInput
              control={form.control}
              name="first_name"
              label="Nombre"
              placeholder="Juan"
              required
            />
            <FormInput
              control={form.control}
              name="last_name"
              label="Apellido"
              placeholder="Pérez"
              required
            />
          </div>
          <div className="grid gap-4 sm:grid-cols-3">
            <FormSelect
              control={form.control}
              name="identity_type"
              label="Tipo Documento"
              options={identityTypeOptions}
              required
            />
            <FormInput
              control={form.control}
              name="identity_number"
              label="Número Documento"
              placeholder="1234 56789 0101"
              required
              disabled={isEditing}
              description={isEditing ? 'El documento no puede ser modificado' : undefined}
            />
            <FormSelect
              control={form.control}
              name="gender"
              label="Género"
              options={genderOptions}
            />
          </div>
          <div className="grid gap-4 sm:grid-cols-2">
            <FormDatePicker
              control={form.control}
              name="birth_date"
              label="Fecha de Nacimiento"
              maxDate={new Date()}
              captionLayout="dropdown"
              fromYear={1920}
              toYear={new Date().getFullYear()}
            />
          </div>
        </div>

        {/* Contact Information */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Información de Contacto</h3>
          <div className="grid gap-4 sm:grid-cols-2">
            <FormInput
              control={form.control}
              name="phone"
              label="Teléfono"
              placeholder="5555-1234"
              required
            />
            <FormInput
              control={form.control}
              name="phone_secondary"
              label="Teléfono Secundario"
              placeholder="5555-5678"
            />
          </div>
          <FormInput
            control={form.control}
            name="email"
            label="Email"
            type="email"
            placeholder="juan@ejemplo.com"
          />
          <FormInput
            control={form.control}
            name="address"
            label="Dirección"
            placeholder="12 Calle 1-25 Zona 1"
          />
          <div className="grid gap-4 sm:grid-cols-3">
            <FormInput
              control={form.control}
              name="city"
              label="Ciudad"
              placeholder="Guatemala"
            />
            <FormInput
              control={form.control}
              name="state"
              label="Departamento"
              placeholder="Guatemala"
            />
            <FormInput
              control={form.control}
              name="postal_code"
              label="Código Postal"
              placeholder="01001"
            />
          </div>
        </div>

        {/* Emergency Contact */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Contacto de Emergencia</h3>
          <div className="grid gap-4 sm:grid-cols-3">
            <FormInput
              control={form.control}
              name="emergency_contact_name"
              label="Nombre"
              placeholder="María García"
            />
            <FormInput
              control={form.control}
              name="emergency_contact_phone"
              label="Teléfono"
              placeholder="5555-9999"
            />
            <FormInput
              control={form.control}
              name="emergency_contact_relation"
              label="Relación"
              placeholder="Hermano/a"
            />
          </div>
        </div>

        {/* Work Information */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Información Laboral</h3>
          <div className="grid gap-4 sm:grid-cols-2">
            <FormInput
              control={form.control}
              name="occupation"
              label="Ocupación"
              placeholder="Comerciante"
            />
            <FormInput
              control={form.control}
              name="workplace"
              label="Lugar de Trabajo"
              placeholder="Mercado Central"
            />
          </div>
          <div className="grid gap-4 sm:grid-cols-2">
            <FormCurrencyInput
              control={form.control}
              name="monthly_income"
              label="Ingreso Mensual"
            />
            <FormCurrencyInput
              control={form.control}
              name="credit_limit"
              label="Límite de Crédito"
            />
          </div>
        </div>

        {/* Notes */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Notas</h3>
          <FormTextarea
            control={form.control}
            name="notes"
            placeholder="Observaciones sobre el cliente..."
            rows={4}
          />
        </div>

        {/* Status */}
        {isEditing && (
          <FormSwitch
            control={form.control}
            name="is_active"
            label="Cliente activo"
            description="Los clientes inactivos no pueden realizar operaciones"
          />
        )}

        {/* Actions */}
        <div className="flex justify-end gap-4">
          <Button type="button" variant="outline" onClick={onCancel}>
            Cancelar
          </Button>
          <Button type="submit" disabled={isLoading}>
            {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {isEditing ? 'Guardar Cambios' : 'Crear Cliente'}
          </Button>
        </div>
      </form>
    </Form>
  )
}
