import { useState, useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2, Calculator, DollarSign } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Form } from '@/components/ui/form'
import { FormInput, FormSelect, FormTextarea } from '@/components/form'
import { ROUTES, loanRoute } from '@/routes/routes'
import { useLoan } from '@/hooks/use-loans'
import { useCreatePayment, usePayoffCalculation, useMinimumPaymentCalculation } from '@/hooks/use-payments'
import { PAYMENT_METHODS } from '@/types'
import { formatCurrency } from '@/lib/format'
import { paymentFormSchema, PaymentFormValues, defaultPaymentValues } from './schemas'

export default function PaymentCreatePage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const loanIdParam = searchParams.get('loan_id')
  const loanId = loanIdParam ? parseInt(loanIdParam, 10) : 0

  const [customLoanId, setCustomLoanId] = useState<number>(loanId)

  const { data: loan, isLoading: loadingLoan } = useLoan(customLoanId)
  const { data: payoff } = usePayoffCalculation(customLoanId)
  const { data: minimum } = useMinimumPaymentCalculation(customLoanId)

  const createMutation = useCreatePayment()

  const form = useForm<PaymentFormValues>({
    resolver: zodResolver(paymentFormSchema),
    defaultValues: {
      ...defaultPaymentValues,
      loan_id: loanId || undefined,
    },
  })

  useEffect(() => {
    if (loanId) {
      form.setValue('loan_id', loanId)
      setCustomLoanId(loanId)
    }
  }, [loanId, form])

  const handleSubmit = (values: PaymentFormValues) => {
    createMutation.mutate(
      {
        loan_id: values.loan_id,
        amount: values.amount,
        payment_method: values.payment_method,
        reference_number: values.reference_number || undefined,
        notes: values.notes || undefined,
      },
      {
        onSuccess: () => {
          navigate(ROUTES.PAYMENTS)
        },
      }
    )
  }

  const handleCancel = () => {
    navigate(ROUTES.PAYMENTS)
  }

  const setPayoffAmount = () => {
    if (payoff) {
      form.setValue('amount', payoff.total_payoff)
    }
  }

  const setMinimumAmount = () => {
    if (minimum) {
      form.setValue('amount', minimum.minimum_amount)
    }
  }

  const paymentMethodOptions = PAYMENT_METHODS.map((m) => ({
    value: m.value,
    label: m.label,
  }))

  const balance = loan
    ? loan.principal_remaining + loan.interest_remaining + loan.late_fee_amount
    : 0

  return (
    <div>
      <PageHeader
        title="Nuevo Pago"
        description="Registrar un nuevo pago"
        backUrl={ROUTES.PAYMENTS}
      />

      <div className="grid gap-6 md:grid-cols-3">
        <div className="md:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle>Información del Pago</CardTitle>
            </CardHeader>
            <CardContent>
              <Form {...form}>
                <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-6">
                  <FormInput
                    control={form.control}
                    name="loan_id"
                    label="ID del Préstamo"
                    type="number"
                    required
                    description={loan ? `Préstamo: ${loan.loan_number}` : 'Ingrese el ID del préstamo'}
                    onChange={(e) => setCustomLoanId(Number(e.target.value))}
                  />

                  {loadingLoan && customLoanId > 0 && (
                    <div className="flex items-center gap-2 text-muted-foreground">
                      <Loader2 className="h-4 w-4 animate-spin" />
                      Cargando préstamo...
                    </div>
                  )}

                  {loan && (
                    <div className="rounded-lg border p-4 space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Cliente:</span>
                        <span>
                          {loan.customer?.first_name} {loan.customer?.last_name}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-muted-foreground">Saldo pendiente:</span>
                        <span className="font-medium">{formatCurrency(balance)}</span>
                      </div>
                    </div>
                  )}

                  <div className="flex items-end gap-4">
                    <div className="flex-1">
                      <FormInput
                        control={form.control}
                        name="amount"
                        label="Monto del Pago"
                        type="number"
                        required
                      />
                    </div>
                    {payoff && (
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={setPayoffAmount}
                      >
                        Liquidar ({formatCurrency(payoff.total_payoff)})
                      </Button>
                    )}
                    {minimum && (
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={setMinimumAmount}
                      >
                        Mínimo ({formatCurrency(minimum.minimum_amount)})
                      </Button>
                    )}
                  </div>

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
                    placeholder="Observaciones sobre el pago..."
                    rows={2}
                  />

                  <div className="flex justify-end gap-4">
                    <Button type="button" variant="outline" onClick={handleCancel}>
                      Cancelar
                    </Button>
                    <Button type="submit" disabled={createMutation.isPending || !loan}>
                      {createMutation.isPending && (
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      )}
                      Registrar Pago
                    </Button>
                  </div>
                </form>
              </Form>
            </CardContent>
          </Card>
        </div>

        {/* Calculator Panel */}
        <div className="space-y-4">
          {payoff && (
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm flex items-center gap-2">
                  <Calculator className="h-4 w-4" />
                  Liquidación Total
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Capital:</span>
                  <span>{formatCurrency(payoff.principal_remaining)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Interés:</span>
                  <span>{formatCurrency(payoff.interest_remaining)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Mora:</span>
                  <span>{formatCurrency(payoff.late_fee_amount)}</span>
                </div>
                <div className="flex justify-between border-t pt-2 font-medium">
                  <span>Total:</span>
                  <span>{formatCurrency(payoff.total_payoff)}</span>
                </div>
              </CardContent>
            </Card>
          )}

          {minimum && (
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm flex items-center gap-2">
                  <DollarSign className="h-4 w-4" />
                  Pago Mínimo
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Monto mínimo:</span>
                  <span className="font-medium">{formatCurrency(minimum.minimum_amount)}</span>
                </div>
                {minimum.installment_number && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Cuota #:</span>
                    <span>{minimum.installment_number}</span>
                  </div>
                )}
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  )
}
