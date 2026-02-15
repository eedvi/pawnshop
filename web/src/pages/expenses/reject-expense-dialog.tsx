import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

import { Expense } from '@/types'
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
import { formatCurrency } from '@/lib/format'
import { rejectExpenseSchema, RejectExpenseFormValues } from './schemas'

interface RejectExpenseDialogProps {
  expense: Expense | null
  open: boolean
  onOpenChange: (open: boolean) => void
  onConfirm: (reason: string) => void
  isLoading?: boolean
}

export function RejectExpenseDialog({
  expense,
  open,
  onOpenChange,
  onConfirm,
  isLoading,
}: RejectExpenseDialogProps) {
  const form = useForm<RejectExpenseFormValues>({
    resolver: zodResolver(rejectExpenseSchema),
    defaultValues: {
      reason: '',
    },
  })

  const handleSubmit = (values: RejectExpenseFormValues) => {
    onConfirm(values.reason)
  }

  const handleClose = () => {
    form.reset()
    onOpenChange(false)
  }

  if (!expense) return null

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Rechazar Gasto</DialogTitle>
          <DialogDescription>
            Rechazar el gasto #{expense.id} por {formatCurrency(expense.amount)}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="reason"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Motivo del Rechazo</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Explica el motivo del rechazo..."
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
                Cancelar
              </Button>
              <Button type="submit" variant="destructive" disabled={isLoading}>
                {isLoading ? 'Rechazando...' : 'Rechazar Gasto'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
