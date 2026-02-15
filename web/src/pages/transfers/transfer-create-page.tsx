import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Search, Package } from 'lucide-react'

import { Item } from '@/types'
import { PageHeader } from '@/components/layout/page-header'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { ROUTES, transferRoute } from '@/routes/routes'
import { useBranches } from '@/hooks/use-branches'
import { useBranchStore } from '@/stores/branch-store'
import { useItems } from '@/hooks/use-items'
import { useCreateTransfer } from '@/hooks/use-transfers'
import { transferFormSchema, TransferFormValues, defaultTransferValues } from './schemas'
import { formatCurrency } from '@/lib/format'

export default function TransferCreatePage() {
  const navigate = useNavigate()
  const { selectedBranch } = useBranchStore()
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedItem, setSelectedItem] = useState<Item | null>(null)

  const { data: branches } = useBranches()
  const { data: itemsResponse } = useItems({
    page: 1,
    per_page: 10,
    branch_id: selectedBranch?.id,
    status: 'in_stock',
    search: searchQuery,
  })
  const createMutation = useCreateTransfer()

  const items = itemsResponse?.data || []

  const form = useForm<TransferFormValues>({
    resolver: zodResolver(transferFormSchema),
    defaultValues: defaultTransferValues,
  })

  const handleItemSelect = (item: Item) => {
    setSelectedItem(item)
    form.setValue('item_id', item.id)
    setSearchQuery('')
  }

  const handleSubmit = (values: TransferFormValues) => {
    createMutation.mutate(
      {
        item_id: values.item_id,
        to_branch_id: values.to_branch_id,
        reason: values.reason || undefined,
        notes: values.notes || undefined,
      },
      {
        onSuccess: (response) => {
          if (response.data) {
            navigate(transferRoute(response.data.id))
          } else {
            navigate(ROUTES.TRANSFERS)
          }
        },
      }
    )
  }

  const handleCancel = () => {
    navigate(ROUTES.TRANSFERS)
  }

  // Filter branches to exclude current branch
  const availableBranches = branches?.data?.filter((b) => b.id !== selectedBranch?.id) || []

  return (
    <div>
      <PageHeader
        title="Nueva Transferencia"
        description="Transferir un artículo a otra sucursal"
        backUrl={ROUTES.TRANSFERS}
      />

      <div className="grid gap-6 md:grid-cols-3">
        <div className="md:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle>Información de la Transferencia</CardTitle>
            </CardHeader>
            <CardContent>
              <Form {...form}>
                <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-6">
                  {/* Item Search */}
                  <div className="space-y-4">
                    <FormLabel>Artículo a Transferir</FormLabel>
                    {!selectedItem ? (
                      <div className="space-y-2">
                        <div className="relative">
                          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                          <Input
                            placeholder="Buscar artículo por nombre o SKU..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            className="pl-10"
                          />
                        </div>
                        {searchQuery && items.length > 0 && (
                          <div className="border rounded-lg divide-y max-h-60 overflow-y-auto">
                            {items.map((item) => (
                              <button
                                key={item.id}
                                type="button"
                                onClick={() => handleItemSelect(item)}
                                className="w-full p-3 text-left hover:bg-muted/50 transition-colors"
                              >
                                <div className="flex items-center gap-3">
                                  <Package className="h-8 w-8 text-muted-foreground" />
                                  <div className="flex-1">
                                    <p className="font-medium">{item.name}</p>
                                    <p className="text-sm text-muted-foreground">
                                      SKU: {item.sku} | Valor: {formatCurrency(item.appraisal_value)}
                                    </p>
                                  </div>
                                </div>
                              </button>
                            ))}
                          </div>
                        )}
                        {searchQuery && items.length === 0 && (
                          <p className="text-sm text-muted-foreground text-center py-4">
                            No se encontraron artículos disponibles
                          </p>
                        )}
                      </div>
                    ) : (
                      <div className="border rounded-lg p-4 bg-muted/20">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-3">
                            <Package className="h-10 w-10 text-muted-foreground" />
                            <div>
                              <p className="font-medium">{selectedItem.name}</p>
                              <p className="text-sm text-muted-foreground">
                                SKU: {selectedItem.sku}
                              </p>
                            </div>
                          </div>
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={() => {
                              setSelectedItem(null)
                              form.setValue('item_id', 0)
                            }}
                          >
                            Cambiar
                          </Button>
                        </div>
                      </div>
                    )}
                    <FormField
                      control={form.control}
                      name="item_id"
                      render={() => <FormMessage />}
                    />
                  </div>

                  <FormField
                    control={form.control}
                    name="to_branch_id"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Sucursal Destino</FormLabel>
                        <Select
                          onValueChange={field.onChange}
                          value={field.value?.toString()}
                        >
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Seleccionar sucursal destino" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {availableBranches.map((branch) => (
                              <SelectItem key={branch.id} value={branch.id.toString()}>
                                {branch.name}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormDescription>
                          Sucursal a la que se enviará el artículo
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="reason"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Motivo de la Transferencia (opcional)</FormLabel>
                        <FormControl>
                          <Input
                            placeholder="Ej: Solicitud de cliente, Reabastecimiento..."
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="notes"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Notas (opcional)</FormLabel>
                        <FormControl>
                          <Textarea
                            placeholder="Notas adicionales..."
                            {...field}
                            rows={3}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className="flex gap-2 justify-end">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={handleCancel}
                      disabled={createMutation.isPending}
                    >
                      Cancelar
                    </Button>
                    <Button type="submit" disabled={createMutation.isPending}>
                      {createMutation.isPending ? 'Creando...' : 'Crear Transferencia'}
                    </Button>
                  </div>
                </form>
              </Form>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-4">
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm">Información</CardTitle>
            </CardHeader>
            <CardContent className="text-sm text-muted-foreground space-y-2">
              <p>La transferencia quedará en estado <Badge variant="secondary">Pendiente</Badge> hasta que sea aprobada.</p>
              <p>Una vez aprobada, podrá ser enviada y posteriormente recibida en la sucursal destino.</p>
            </CardContent>
          </Card>

          {selectedItem && (
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm">Artículo Seleccionado</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Categoría:</span>
                  <span>{selectedItem.category?.name || '-'}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Valor:</span>
                  <span className="font-medium">{formatCurrency(selectedItem.appraisal_value)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Sucursal actual:</span>
                  <span>{selectedBranch?.name || '-'}</span>
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  )
}
