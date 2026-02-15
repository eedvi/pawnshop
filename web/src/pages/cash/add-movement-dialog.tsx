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
import { FormInput, FormSelect, FormTextarea } from '@/components/form'
import { Loader2 } from 'lucide-react'
import { CASH_MOVEMENT_TYPES } from '@/types'
import { cashMovementSchema, CashMovementFormValues } from './schemas'

interface AddMovementDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  sessionId: number
  onConfirm: (data: CashMovementFormValues) => void
  isLoading?: boolean
}

export function AddMovementDialog({
  open,
  onOpenChange,
  sessionId,
  onConfirm,
  isLoading = false,
}: AddMovementDialogProps) {
  const form = useForm<CashMovementFormValues>({
    resolver: zodResolver(cashMovementSchema),
    defaultValues: {
      movement_type: 'income',
      amount: 0,
      description: '',
    },
  })

  const handleSubmit = (values: CashMovementFormValues) => {
    onConfirm(values)
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      form.reset()
    }
    onOpenChange(newOpen)
  }

  const movementTypeOptions = CASH_MOVEMENT_TYPES.map((t) => ({
    value: t.value,
    label: t.label,
  }))

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Agregar Movimiento</DialogTitle>
          <DialogDescription>
            Registre un ingreso, egreso o ajuste en la caja
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormSelect
              control={form.control}
              name="movement_type"
              label="Tipo de Movimiento"
              options={movementTypeOptions}
              required
            />

            <FormInput
              control={form.control}
              name="amount"
              label="Monto"
              type="number"
              required
            />

            <FormTextarea
              control={form.control}
              name="description"
              label="Descripción"
              placeholder="Descripción del movimiento..."
              required
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
                Agregar
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
