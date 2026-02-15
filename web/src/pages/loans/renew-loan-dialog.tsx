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
import { FormInput, FormSwitch } from '@/components/form'
import { Loan } from '@/types'
import { renewLoanSchema, RenewLoanFormValues } from './schemas'
import { formatCurrency } from '@/lib/format'

interface RenewLoanDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  loan: Loan | null
  onConfirm: (values: RenewLoanFormValues) => void
  isLoading?: boolean
}

export function RenewLoanDialog({
  open,
  onOpenChange,
  loan,
  onConfirm,
  isLoading = false,
}: RenewLoanDialogProps) {
  const form = useForm<RenewLoanFormValues>({
    resolver: zodResolver(renewLoanSchema),
    defaultValues: {
      new_term_days: loan?.loan_term_days ?? 30,
      new_interest_rate: loan?.interest_rate,
      pay_interest: false,
    },
  })

  const handleSubmit = (values: RenewLoanFormValues) => {
    onConfirm(values)
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      form.reset()
    }
    onOpenChange(newOpen)
  }

  const balance = loan
    ? loan.principal_remaining + loan.interest_remaining + loan.late_fee_amount
    : 0

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Renovar Préstamo</DialogTitle>
          <DialogDescription>
            {loan && (
              <>
                Renovar préstamo <strong>{loan.loan_number}</strong>.
                Saldo pendiente: {formatCurrency(balance)}
              </>
            )}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormInput
              control={form.control}
              name="new_term_days"
              label="Nuevo Plazo (días)"
              type="number"
              placeholder="30"
            />

            <FormInput
              control={form.control}
              name="new_interest_rate"
              label="Nueva Tasa de Interés (%)"
              type="number"
              placeholder="15"
            />

            <FormSwitch
              control={form.control}
              name="pay_interest"
              label="Pagar Interés"
              description="Marcar si se pagará el interés acumulado antes de renovar"
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
                {isLoading ? 'Procesando...' : 'Renovar'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
