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
import { Textarea } from '@/components/ui/textarea'
import { cancelTransferSchema, CancelTransferFormValues } from './schemas'

interface CancelTransferDialogProps {
  transfer: ItemTransfer | null
  open: boolean
  onOpenChange: (open: boolean) => void
  onConfirm: (reason: string) => void
  isLoading?: boolean
}

export function CancelTransferDialog({
  transfer,
  open,
  onOpenChange,
  onConfirm,
  isLoading,
}: CancelTransferDialogProps) {
  const form = useForm<CancelTransferFormValues>({
    resolver: zodResolver(cancelTransferSchema),
    defaultValues: {
      reason: '',
    },
  })

  const handleSubmit = (values: CancelTransferFormValues) => {
    onConfirm(values.reason)
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
          <DialogTitle>Cancelar Transferencia</DialogTitle>
          <DialogDescription>
            Cancelar la transferencia {transfer.transfer_number}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="reason"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Motivo de Cancelación</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Explica el motivo de la cancelación..."
                      {...field}
                      rows={4}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button type="button" variant="outline" onClick={handleClose} disabled={isLoading}>
                Cerrar
              </Button>
              <Button type="submit" variant="destructive" disabled={isLoading}>
                {isLoading ? 'Cancelando...' : 'Cancelar Transferencia'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
