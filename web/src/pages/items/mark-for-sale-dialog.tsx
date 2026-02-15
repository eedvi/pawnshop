import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Form } from '@/components/ui/form'
import { FormCurrencyInput } from '@/components/form'
import { Item } from '@/types'
import { markForSaleSchema, MarkForSaleFormValues } from './schemas'
import { formatCurrency } from '@/lib/format'

interface MarkForSaleDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  item: Item | null
  onConfirm: (salePrice: number) => void
  isLoading?: boolean
}

export function MarkForSaleDialog({
  open,
  onOpenChange,
  item,
  onConfirm,
  isLoading = false,
}: MarkForSaleDialogProps) {
  const form = useForm<MarkForSaleFormValues>({
    resolver: zodResolver(markForSaleSchema),
    defaultValues: {
      sale_price: item?.appraised_value || 0,
    },
  })

  const handleSubmit = (values: MarkForSaleFormValues) => {
    onConfirm(values.sale_price)
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      form.reset()
    }
    onOpenChange(newOpen)
  }

  // Reset form when item changes
  if (item && form.getValues('sale_price') === 0) {
    form.setValue('sale_price', item.appraised_value)
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Marcar para Venta</DialogTitle>
          <DialogDescription>
            {item && (
              <>
                Marcar <strong>{item.name}</strong> para venta.
                Avalúo actual: {formatCurrency(item.appraised_value)}
              </>
            )}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormCurrencyInput
              control={form.control}
              name="sale_price"
              label="Precio de Venta"
              required
            />

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => handleOpenChange(false)}
                disabled={isLoading}
              >
                Cancelar
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? 'Procesando...' : 'Confirmar'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
