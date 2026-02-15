import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2 } from 'lucide-react'

import { Branch } from '@/types'
import { Button } from '@/components/ui/button'
import { Form } from '@/components/ui/form'
import { FormInput, FormSelect, FormSwitch } from '@/components/form'
import { branchFormSchema, BranchFormValues, defaultBranchValues } from './schemas'

interface BranchFormProps {
  branch?: Branch
  onSubmit: (values: BranchFormValues) => void
  onCancel: () => void
  isLoading?: boolean
}

const timezoneOptions = [
  { value: 'America/Guatemala', label: 'Guatemala (UTC-6)' },
  { value: 'America/Mexico_City', label: 'Ciudad de México (UTC-6)' },
  { value: 'America/El_Salvador', label: 'El Salvador (UTC-6)' },
  { value: 'America/Tegucigalpa', label: 'Honduras (UTC-6)' },
  { value: 'America/Managua', label: 'Nicaragua (UTC-6)' },
  { value: 'America/Costa_Rica', label: 'Costa Rica (UTC-6)' },
  { value: 'America/Panama', label: 'Panamá (UTC-5)' },
]

const currencyOptions = [
  { value: 'GTQ', label: 'Quetzal (GTQ)' },
  { value: 'USD', label: 'Dólar (USD)' },
  { value: 'MXN', label: 'Peso Mexicano (MXN)' },
]

export function BranchForm({ branch, onSubmit, onCancel, isLoading }: BranchFormProps) {
  const form = useForm<BranchFormValues>({
    resolver: zodResolver(branchFormSchema),
    defaultValues: branch
      ? {
          name: branch.name,
          code: branch.code,
          address: branch.address || '',
          phone: branch.phone || '',
          email: branch.email || '',
          timezone: branch.timezone,
          currency: branch.currency,
          default_interest_rate: branch.default_interest_rate,
          default_loan_term_days: branch.default_loan_term_days,
          default_grace_period: branch.default_grace_period,
          is_active: branch.is_active,
        }
      : defaultBranchValues,
  })

  const isEditing = !!branch

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        {/* Basic Info */}
        <div className="grid gap-4 sm:grid-cols-2">
          <FormInput
            control={form.control}
            name="name"
            label="Nombre"
            placeholder="Sucursal Central"
            required
          />
          <FormInput
            control={form.control}
            name="code"
            label="Código"
            placeholder="SC01"
            required
            disabled={isEditing}
            description={isEditing ? 'El código no puede ser modificado' : undefined}
          />
        </div>

        <FormInput
          control={form.control}
          name="address"
          label="Dirección"
          placeholder="12 Calle 1-25 Zona 1, Guatemala"
        />

        <div className="grid gap-4 sm:grid-cols-2">
          <FormInput
            control={form.control}
            name="phone"
            label="Teléfono"
            placeholder="2222-3333"
          />
          <FormInput
            control={form.control}
            name="email"
            label="Email"
            type="email"
            placeholder="sucursal@empresa.com"
          />
        </div>

        {/* Regional Settings */}
        <div className="grid gap-4 sm:grid-cols-2">
          <FormSelect
            control={form.control}
            name="timezone"
            label="Zona Horaria"
            options={timezoneOptions}
          />
          <FormSelect
            control={form.control}
            name="currency"
            label="Moneda"
            options={currencyOptions}
          />
        </div>

        {/* Loan Defaults */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Configuración de Préstamos</h3>
          <div className="grid gap-4 sm:grid-cols-3">
            <FormInput
              control={form.control}
              name="default_interest_rate"
              label="Tasa de Interés (%)"
              type="number"
              min={0}
              max={100}
              step={0.1}
            />
            <FormInput
              control={form.control}
              name="default_loan_term_days"
              label="Plazo (días)"
              type="number"
              min={1}
              max={365}
            />
            <FormInput
              control={form.control}
              name="default_grace_period"
              label="Período de Gracia (días)"
              type="number"
              min={0}
              max={30}
            />
          </div>
        </div>

        {/* Status */}
        {isEditing && (
          <FormSwitch
            control={form.control}
            name="is_active"
            label="Sucursal activa"
            description="Las sucursales inactivas no pueden realizar operaciones"
          />
        )}

        {/* Actions */}
        <div className="flex justify-end gap-4">
          <Button type="button" variant="outline" onClick={onCancel}>
            Cancelar
          </Button>
          <Button type="submit" disabled={isLoading}>
            {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {isEditing ? 'Guardar Cambios' : 'Crear Sucursal'}
          </Button>
        </div>
      </form>
    </Form>
  )
}
