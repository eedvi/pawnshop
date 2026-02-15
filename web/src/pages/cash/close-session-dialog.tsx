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
import { Loader2, AlertTriangle } from 'lucide-react'
import { CashSession, CashSessionSummary } from '@/types'
import { formatCurrency } from '@/lib/format'
import { closeSessionSchema, CloseSessionFormValues } from './schemas'

interface CloseSessionDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  session: CashSession | null
  summary: CashSessionSummary | null
  onConfirm: (data: CloseSessionFormValues) => void
  isLoading?: boolean
}

export function CloseSessionDialog({
  open,
  onOpenChange,
  session,
  summary,
  onConfirm,
  isLoading = false,
}: CloseSessionDialogProps) {
  const form = useForm<CloseSessionFormValues>({
    resolver: zodResolver(closeSessionSchema),
    defaultValues: {
      closing_amount: summary?.expected_balance || 0,
      notes: '',
    },
  })

  const closingAmount = form.watch('closing_amount')
  const expectedAmount = summary?.expected_balance || 0
  const difference = closingAmount - expectedAmount

  const handleSubmit = (values: CloseSessionFormValues) => {
    onConfirm(values)
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      form.reset()
    }
    onOpenChange(newOpen)
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Cerrar Sesión de Caja</DialogTitle>
          <DialogDescription>
            Verifique el monto en efectivo y cierre la sesión
          </DialogDescription>
        </DialogHeader>

        {summary && (
          <div className="rounded-lg border p-4 space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Monto inicial:</span>
              <span>{formatCurrency(summary.opening_amount)}</span>
            </div>
            <div className="flex justify-between text-green-600">
              <span>Total ingresos:</span>
              <span>+{formatCurrency(summary.total_income)}</span>
            </div>
            <div className="flex justify-between text-red-600">
              <span>Total egresos:</span>
              <span>-{formatCurrency(summary.total_expense)}</span>
            </div>
            <div className="flex justify-between border-t pt-2 font-medium">
              <span>Saldo esperado:</span>
              <span>{formatCurrency(expectedAmount)}</span>
            </div>
          </div>
        )}

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormInput
              control={form.control}
              name="closing_amount"
              label="Monto en Caja"
              type="number"
              required
              description="Monto en efectivo contado al cerrar"
            />

            {difference !== 0 && (
              <div
                className={`flex items-center gap-2 p-3 rounded-lg ${
                  difference > 0
                    ? 'bg-green-50 text-green-700 dark:bg-green-950 dark:text-green-300'
                    : 'bg-red-50 text-red-700 dark:bg-red-950 dark:text-red-300'
                }`}
              >
                <AlertTriangle className="h-4 w-4" />
                <span className="text-sm">
                  Diferencia: {difference > 0 ? '+' : ''}
                  {formatCurrency(difference)}
                  {difference > 0 ? ' (sobrante)' : ' (faltante)'}
                </span>
              </div>
            )}

            <FormTextarea
              control={form.control}
              name="notes"
              label="Notas"
              placeholder="Observaciones al cerrar la caja..."
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
                Cerrar Sesión
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
