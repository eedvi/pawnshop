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
import { FormInput, FormTextarea } from '@/components/form'
import { Loader2 } from 'lucide-react'
import { CashRegister } from '@/types'
import { openSessionSchema, OpenSessionFormValues } from './schemas'

interface OpenSessionDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  register: CashRegister | null
  onConfirm: (data: OpenSessionFormValues) => void
  isLoading?: boolean
}

export function OpenSessionDialog({
  open,
  onOpenChange,
  register,
  onConfirm,
  isLoading = false,
}: OpenSessionDialogProps) {
  const form = useForm<OpenSessionFormValues>({
    resolver: zodResolver(openSessionSchema),
    defaultValues: {
      register_id: register?.id || 0,
      opening_amount: 0,
      notes: '',
    },
  })

  const handleSubmit = (values: OpenSessionFormValues) => {
    onConfirm({
      ...values,
      register_id: register?.id || 0,
    })
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      form.reset()
    }
    onOpenChange(newOpen)
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Abrir Sesión de Caja</DialogTitle>
          <DialogDescription>
            Ingrese el monto inicial para abrir la sesión en {register?.name}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormInput
              control={form.control}
              name="opening_amount"
              label="Monto Inicial"
              type="number"
              required
              description="Monto en efectivo al iniciar la sesión"
            />

            <FormTextarea
              control={form.control}
              name="notes"
              label="Notas"
              placeholder="Observaciones al abrir la caja..."
              rows={2}
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
                {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Abrir Sesión
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
