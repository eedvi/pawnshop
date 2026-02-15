import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2, Search } from 'lucide-react'

import { Item, ITEM_CONDITIONS, Customer } from '@/types'
import { Button } from '@/components/ui/button'
import { Form } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { FormInput, FormSelect, FormTextarea, FormCurrencyInput } from '@/components/form'
import { itemFormSchema, ItemFormValues, defaultItemValues } from './schemas'
import { useCategories } from '@/hooks/use-categories'
import { useCustomers } from '@/hooks/use-customers'
import { useDebounce } from '@/hooks'
import { useBranchStore } from '@/stores/branch-store'

interface ItemFormProps {
  item?: Item
  onSubmit: (values: ItemFormValues) => void
  onCancel: () => void
  isLoading?: boolean
  showCustomer?: boolean
}

const ACQUISITION_TYPES = [
  { value: 'pawn', label: 'Empeño' },
  { value: 'purchase', label: 'Compra' },
  { value: 'consignment', label: 'Consignación' },
]

export function ItemForm({ item, onSubmit, onCancel, isLoading, showCustomer = false }: ItemFormProps) {
  const { selectedBranchId } = useBranchStore()
  const [customerSearch, setCustomerSearch] = useState('')
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(item?.customer || null)
  const debouncedCustomerSearch = useDebounce(customerSearch, 300)

  const form = useForm<ItemFormValues>({
    resolver: zodResolver(itemFormSchema),
    defaultValues: item
      ? {
          name: item.name,
          description: item.description || '',
          brand: item.brand || '',
          model: item.model || '',
          serial_number: item.serial_number || '',
          color: item.color || '',
          condition: item.condition,
          category_id: item.category_id ?? 'none',
          customer_id: item.customer_id ?? null,
          appraised_value: item.appraised_value,
          loan_value: item.loan_value,
          sale_price: item.sale_price ?? null,
          weight: item.weight ?? null,
          purity: item.purity || '',
          acquisition_type: item.acquisition_type,
          acquisition_price: item.acquisition_price ?? null,
          notes: item.notes || '',
          tags: item.tags || [],
        }
      : defaultItemValues,
  })

  const isEditing = !!item

  // Fetch categories for select
  const { data: categories } = useCategories()

  // Fetch customers for selector
  const { data: customersData, isLoading: loadingCustomers } = useCustomers({
    search: debouncedCustomerSearch || undefined,
    branch_id: selectedBranchId ?? undefined,
    is_active: true,
    per_page: 10,
  })
  const categoryOptions = [
    { value: 'none', label: 'Sin categoría' },
    ...(categories?.map((c) => ({
      value: c.id.toString(),
      label: c.name,
    })) || []),
  ]

  const conditionOptions = ITEM_CONDITIONS.map((c) => ({
    value: c.value,
    label: c.label,
  }))

  const acquisitionTypeOptions = ACQUISITION_TYPES.map((t) => ({
    value: t.value,
    label: t.label,
  }))

  const handleFormSubmit = (values: ItemFormValues) => {
    // Convert category_id from string to number
    let categoryId: number | null = null
    if (values.category_id && values.category_id !== 'none') {
      categoryId = typeof values.category_id === 'string'
        ? parseInt(values.category_id, 10)
        : values.category_id
    }

    onSubmit({
      ...values,
      category_id: categoryId,
      customer_id: selectedCustomer?.id ?? null,
    })
  }

  const handleSelectCustomer = (customer: Customer) => {
    setSelectedCustomer(customer)
    setCustomerSearch('')
  }

  const handleClearCustomer = () => {
    setSelectedCustomer(null)
    setCustomerSearch('')
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleFormSubmit)} className="space-y-8">
        {/* Basic Information */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Información Básica</h3>
          <div className="grid gap-4 sm:grid-cols-2">
            <FormInput
              control={form.control}
              name="name"
              label="Nombre del Artículo"
              placeholder="Anillo de oro 14k"
              required
            />
            <FormSelect
              control={form.control}
              name="category_id"
              label="Categoría"
              options={categoryOptions}
            />
          </div>
          <FormTextarea
            control={form.control}
            name="description"
            label="Descripción"
            placeholder="Descripción detallada del artículo..."
            rows={3}
          />
          <div className="grid gap-4 sm:grid-cols-4">
            <FormInput
              control={form.control}
              name="brand"
              label="Marca"
              placeholder="Samsung"
            />
            <FormInput
              control={form.control}
              name="model"
              label="Modelo"
              placeholder="Galaxy S23"
            />
            <FormInput
              control={form.control}
              name="serial_number"
              label="No. Serie"
              placeholder="ABC123456"
            />
            <FormInput
              control={form.control}
              name="color"
              label="Color"
              placeholder="Dorado"
            />
          </div>
        </div>

        {/* Customer Selection */}
        {showCustomer && (
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Cliente</h3>
            {selectedCustomer ? (
              <div className="flex items-center justify-between p-4 rounded-lg border bg-primary/5 border-primary">
                <div>
                  <p className="font-medium">
                    {selectedCustomer.first_name} {selectedCustomer.last_name}
                  </p>
                  <p className="text-sm text-muted-foreground">
                    {selectedCustomer.identity_number} • {selectedCustomer.phone}
                  </p>
                </div>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={handleClearCustomer}
                >
                  Cambiar
                </Button>
              </div>
            ) : (
              <div className="space-y-3">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                  <Input
                    placeholder="Buscar cliente por nombre o documento..."
                    value={customerSearch}
                    onChange={(e) => setCustomerSearch(e.target.value)}
                    className="pl-10"
                  />
                </div>
                {customerSearch && (
                  <div className="max-h-48 overflow-y-auto rounded-lg border">
                    {loadingCustomers ? (
                      <div className="p-4 text-center text-muted-foreground">
                        <Loader2 className="h-4 w-4 animate-spin mx-auto" />
                      </div>
                    ) : customersData?.data?.length === 0 ? (
                      <p className="p-4 text-center text-muted-foreground text-sm">
                        No se encontraron clientes
                      </p>
                    ) : (
                      customersData?.data?.map((customer) => (
                        <div
                          key={customer.id}
                          onClick={() => handleSelectCustomer(customer)}
                          className="p-3 cursor-pointer hover:bg-muted border-b last:border-b-0"
                        >
                          <p className="font-medium text-sm">
                            {customer.first_name} {customer.last_name}
                          </p>
                          <p className="text-xs text-muted-foreground">
                            {customer.identity_number}
                          </p>
                        </div>
                      ))
                    )}
                  </div>
                )}
                <p className="text-xs text-muted-foreground">
                  Opcional: Asigna este artículo a un cliente para poder usarlo en préstamos.
                </p>
              </div>
            )}
          </div>
        )}

        {/* Condition & Physical Details */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Condición y Detalles Físicos</h3>
          <div className="grid gap-4 sm:grid-cols-3">
            <FormSelect
              control={form.control}
              name="condition"
              label="Condición"
              options={conditionOptions}
              required
            />
            <FormInput
              control={form.control}
              name="weight"
              label="Peso (gramos)"
              type="number"
              placeholder="10.5"
            />
            <FormInput
              control={form.control}
              name="purity"
              label="Pureza/Quilates"
              placeholder="14k / 585"
            />
          </div>
        </div>

        {/* Valuation */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Valoración</h3>
          <div className="grid gap-4 sm:grid-cols-3">
            <FormCurrencyInput
              control={form.control}
              name="appraised_value"
              label="Valor de Avalúo"
              required
            />
            <FormCurrencyInput
              control={form.control}
              name="loan_value"
              label="Valor de Préstamo"
              required
            />
            <FormCurrencyInput
              control={form.control}
              name="sale_price"
              label="Precio de Venta"
            />
          </div>
        </div>

        {/* Acquisition */}
        {!isEditing && (
          <div className="space-y-4">
            <h3 className="text-lg font-medium">Adquisición</h3>
            <div className="grid gap-4 sm:grid-cols-2">
              <FormSelect
                control={form.control}
                name="acquisition_type"
                label="Tipo de Adquisición"
                options={acquisitionTypeOptions}
                required
              />
              <FormCurrencyInput
                control={form.control}
                name="acquisition_price"
                label="Precio de Adquisición"
              />
            </div>
          </div>
        )}

        {/* Notes */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Notas</h3>
          <FormTextarea
            control={form.control}
            name="notes"
            placeholder="Observaciones sobre el artículo..."
            rows={4}
          />
        </div>

        {/* Actions */}
        <div className="flex justify-end gap-4">
          <Button type="button" variant="outline" onClick={onCancel}>
            Cancelar
          </Button>
          <Button type="submit" disabled={isLoading}>
            {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {isEditing ? 'Guardar Cambios' : 'Crear Artículo'}
          </Button>
        </div>
      </form>
    </Form>
  )
}
