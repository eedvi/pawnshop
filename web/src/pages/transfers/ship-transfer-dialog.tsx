import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { ItemTransfer } from '@/types'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { shipTransferSchema, ShipTransferFormValues } from './schemas'

interface ShipTransferDialogProps {
  transfer: ItemTransfer | null
  open: boolean
  onOpenChange: (open: boolean) => void
  onConfirm: (trackingNumber?: string, notes?: string) => void
  isLoading?: boolean
}

export function ShipTransferDialog({
  transfer,
  open,
  onOpenChange,
  onConfirm,
  isLoading,
}: ShipTransferDialogProps) {
  const form = useForm<ShipTransferFormValues>({
    resolver: zodResolver(shipTransferSchema),
    defaultValues: {
      tracking_number: '',
      notes: '',
    },
  })

  const handleSubmit = (values: ShipTransferFormValues) => {
    onConfirm(values.tracking_number, values.notes)
  }

  const handleClose = () => {
    form.reset()
    onOpenChange(false)
  }

  if (!transfer) return null

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Enviar Transferencia</DialogTitle>
          <DialogDescription>
            Registrar el envío de la transferencia {transfer.transfer_number}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="tracking_number"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Número de Seguimiento (opcional)</FormLabel>
                  <FormControl>
                    <Input placeholder="Ej: TRACK-12345" {...field} />
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
                      placeholder="Notas adicionales sobre el envío..."
                      {...field}
                      rows={3}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button type="button" variant="outline" onClick={handleClose} disabled={isLoading}>
                Cancelar
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? 'Enviando...' : 'Confirmar Envío'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
