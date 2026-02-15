import { useState, useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2, Package, User } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Form } from '@/components/ui/form'
import { FormInput, FormSelect, FormTextarea } from '@/components/form'
import { ROUTES, itemRoute, customerRoute } from '@/routes/routes'
import { useItem } from '@/hooks/use-items'
import { useCustomer } from '@/hooks/use-customers'
import { useCreateSale } from '@/hooks/use-sales'
import { useBranchStore } from '@/stores/branch-store'
import { PAYMENT_METHODS, SALE_TYPES } from '@/types'
import { formatCurrency } from '@/lib/format'
import { saleFormSchema, SaleFormValues, defaultSaleValues } from './schemas'

export default function SaleCreatePage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const itemIdParam = searchParams.get('item_id')
  const itemId = itemIdParam ? parseInt(itemIdParam, 10) : 0

  const { selectedBranchId } = useBranchStore()

  const [customItemId, setCustomItemId] = useState<number>(itemId)
  const [customCustomerId, setCustomCustomerId] = useState<number>(0)

  const { data: item, isLoading: loadingItem } = useItem(customItemId)
  const { data: customer, isLoading: loadingCustomer } = useCustomer(customCustomerId)

  const createMutation = useCreateSale()

  const form = useForm<SaleFormValues>({
    resolver: zodResolver(saleFormSchema),
    defaultValues: {
      ...defaultSaleValues,
      item_id: itemId || undefined,
    },
  })

  useEffect(() => {
    if (itemId) {
      form.setValue('item_id', itemId)
      setCustomItemId(itemId)
    }
  }, [itemId, form])

  useEffect(() => {
    if (item?.sale_price) {
      form.setValue('sale_price', item.sale_price)
    }
  }, [item, form])

  const handleSubmit = (values: SaleFormValues) => {
    if (!selectedBranchId) {
      return
    }

    createMutation.mutate(
      {
        branch_id: selectedBranchId,
        item_id: values.item_id,
        customer_id: values.customer_id || undefined,
        sale_type: values.sale_type,
        sale_price: values.sale_price,
        discount_amount: values.discount_amount || 0,
        discount_reason: values.discount_reason || undefined,
        payment_method: values.payment_method,
        reference_number: values.reference_number || undefined,
        notes: values.notes || undefined,
      },
      {
        onSuccess: () => {
          navigate(ROUTES.SALES)
        },
      }
    )
  }

  const handleCancel = () => {
    navigate(ROUTES.SALES)
  }

  const salePrice = form.watch('sale_price') || 0
  const discountAmount = form.watch('discount_amount') || 0
  const finalPrice = salePrice - discountAmount

  const paymentMethodOptions = PAYMENT_METHODS.map((m) => ({
    value: m.value,
    label: m.label,
  }))

  const saleTypeOptions = SALE_TYPES.map((t) => ({
    value: t.value,
    label: t.label,
  }))

  return (
    <div>
      <PageHeader
        title="Nueva Venta"
        description="Registrar una nueva venta"
        backUrl={ROUTES.SALES}
      />

      <div className="grid gap-6 md:grid-cols-3">
        <div className="md:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle>Información de la Venta</CardTitle>
            </CardHeader>
            <CardContent>
              <Form {...form}>
                <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-6">
                  {/* Item Selection */}
                  <FormInput
                    control={form.control}
                    name="item_id"
                    label="ID del Artículo"
                    type="number"
                    required
                    description={
                      item
                        ? `Artículo: ${item.name}`
                        : 'Ingrese el ID del artículo a vender'
                    }
                    onChange={(e) => setCustomItemId(Number(e.target.value))}
                  />

                  {loadingItem && customItemId > 0 && (
                    <div className="flex items-center gap-2 text-muted-foreground">
                      <Loader2 className="h-4 w-4 animate-spin" />
                      Cargando artículo...
                    </div>
                  )}

                  {item && (
                    <div className="rounded-lg border p-4 space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Artículo:</span>
                        <span className="font-medium">{item.name}</span>
                      </div>
                      {item.brand && (
                        <div className="flex justify-between">
                          <span className="text-muted-foreground">Marca:</span>
                          <span>{item.brand}</span>
                        </div>
                      )}
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Estado:</span>
                        <span className="capitalize">{item.status}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Precio sugerido:</span>
                        <span className="font-medium">{formatCurrency(item.sale_price || 0)}</span>
                      </div>
                    </div>
                  )}

                  {/* Customer Selection (Optional) */}
                  <FormInput
                    control={form.control}
                    name="customer_id"
                    label="ID del Cliente (Opcional)"
                    type="number"
                    description={
                      customer
                        ? `Cliente: ${customer.first_name} ${customer.last_name}`
                        : 'Dejar vacío para venta sin cliente registrado'
                    }
                    onChange={(e) => setCustomCustomerId(Number(e.target.value))}
                  />

                  {loadingCustomer && customCustomerId > 0 && (
                    <div className="flex items-center gap-2 text-muted-foreground">
                      <Loader2 className="h-4 w-4 animate-spin" />
                      Cargando cliente...
                    </div>
                  )}

                  <FormSelect
                    control={form.control}
                    name="sale_type"
                    label="Tipo de Venta"
                    options={saleTypeOptions}
                    required
                  />

                  {/* Pricing */}
                  <div className="grid gap-4 sm:grid-cols-2">
                    <FormInput
                      control={form.control}
                      name="sale_price"
                      label="Precio de Venta"
                      type="number"
                      required
                    />
                    <FormInput
                      control={form.control}
                      name="discount_amount"
                      label="Descuento"
                      type="number"
                    />
                  </div>

                  {discountAmount > 0 && (
                    <FormInput
                      control={form.control}
                      name="discount_reason"
                      label="Razón del Descuento"
                      placeholder="Descripción del descuento..."
                    />
                  )}

                  <FormSelect
                    control={form.control}
                    name="payment_method"
                    label="Método de Pago"
                    options={paymentMethodOptions}
                    required
                  />

                  <FormInput
                    control={form.control}
                    name="reference_number"
                    label="Número de Referencia"
                    placeholder="Opcional - para transferencias, cheques, etc."
                  />

                  <FormTextarea
                    control={form.control}
                    name="notes"
                    label="Notas"
                    placeholder="Observaciones sobre la venta..."
                    rows={2}
                  />

                  <div className="flex justify-end gap-4">
                    <Button type="button" variant="outline" onClick={handleCancel}>
                      Cancelar
                    </Button>
                    <Button
                      type="submit"
                      disabled={createMutation.isPending || !item || item.status !== 'for_sale'}
                    >
                      {createMutation.isPending && (
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      )}
                      Registrar Venta
                    </Button>
                  </div>
                </form>
              </Form>
            </CardContent>
          </Card>
        </div>

        {/* Summary Panel */}
        <div className="space-y-4">
          {item && (
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm flex items-center gap-2">
                  <Package className="h-4 w-4" />
                  Artículo
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Nombre:</span>
                  <span className="font-medium">{item.name}</span>
                </div>
                {item.category && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Categoría:</span>
                    <span>{item.category.name}</span>
                  </div>
                )}
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Condición:</span>
                  <span className="capitalize">{item.condition}</span>
                </div>
              </CardContent>
            </Card>
          )}

          {customer && (
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm flex items-center gap-2">
                  <User className="h-4 w-4" />
                  Cliente
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Nombre:</span>
                  <span>
                    {customer.first_name} {customer.last_name}
                  </span>
                </div>
                {customer.phone && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Teléfono:</span>
                    <span>{customer.phone}</span>
                  </div>
                )}
              </CardContent>
            </Card>
          )}

          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm">Resumen</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Precio:</span>
                <span>{formatCurrency(salePrice)}</span>
              </div>
              {discountAmount > 0 && (
                <div className="flex justify-between text-destructive">
                  <span>Descuento:</span>
                  <span>-{formatCurrency(discountAmount)}</span>
                </div>
              )}
              <div className="flex justify-between border-t pt-2 font-bold">
                <span>Total:</span>
                <span>{formatCurrency(finalPrice)}</span>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
