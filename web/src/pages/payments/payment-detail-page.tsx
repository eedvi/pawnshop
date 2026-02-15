import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { Loader2, RotateCcw, CreditCard, Calendar, User, FileText, AlertTriangle } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { ROUTES, loanRoute, customerRoute } from '@/routes/routes'
import { usePayment, useReversePayment } from '@/hooks/use-payments'
import { PAYMENT_STATUSES, PAYMENT_METHODS } from '@/types'
import { formatCurrency, formatDate, formatDateTime } from '@/lib/format'
import { ReversePaymentDialog } from './reverse-payment-dialog'

export default function PaymentDetailPage() {
  const { id } = useParams()
  const paymentId = id ? parseInt(id, 10) : 0

  const [reverseDialogOpen, setReverseDialogOpen] = useState(false)

  const { data: payment, isLoading, error } = usePayment(paymentId)
  const reverseMutation = useReversePayment()

  const handleReverseConfirm = (reason: string) => {
    reverseMutation.mutate(
      { id: paymentId, input: { reason } },
      {
        onSuccess: () => {
          setReverseDialogOpen(false)
        },
      }
    )
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (error || !payment) {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-4">
        <AlertTriangle className="h-12 w-12 text-destructive" />
        <p className="text-muted-foreground">Error al cargar el pago</p>
        <Button asChild variant="outline">
          <Link to={ROUTES.PAYMENTS}>Volver a pagos</Link>
        </Button>
      </div>
    )
  }

  const status = PAYMENT_STATUSES.find((s) => s.value === payment.status)
  const method = PAYMENT_METHODS.find((m) => m.value === payment.payment_method)
  const canReverse = payment.status === 'completed'

  const statusColorMap: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
    green: 'default',
    yellow: 'secondary',
    red: 'destructive',
    gray: 'outline',
  }

  return (
    <div>
      <PageHeader
        title={`Pago ${payment.payment_number}`}
        description="Detalles del pago"
        backUrl={ROUTES.PAYMENTS}
        actions={
          canReverse && (
            <Button
              variant="destructive"
              onClick={() => setReverseDialogOpen(true)}
            >
              <RotateCcw className="mr-2 h-4 w-4" />
              Revertir Pago
            </Button>
          )
        }
      />

      <div className="grid gap-6 md:grid-cols-3">
        {/* Main Info */}
        <div className="md:col-span-2 space-y-6">
          {/* Payment Summary */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <CreditCard className="h-5 w-5" />
                Información del Pago
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <p className="text-sm text-muted-foreground">Número de Pago</p>
                  <p className="font-mono text-lg">{payment.payment_number}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Estado</p>
                  <Badge variant={statusColorMap[status?.color || 'gray'] || 'outline'}>
                    {status?.label || payment.status}
                  </Badge>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Monto</p>
                  <p className="text-2xl font-bold">{formatCurrency(payment.amount)}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Método de Pago</p>
                  <p className="font-medium">{method?.label || payment.payment_method}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Fecha de Pago</p>
                  <p>{formatDate(payment.payment_date)}</p>
                </div>
                {payment.reference_number && (
                  <div>
                    <p className="text-sm text-muted-foreground">Número de Referencia</p>
                    <p className="font-mono">{payment.reference_number}</p>
                  </div>
                )}
              </div>

              {payment.notes && (
                <>
                  <Separator />
                  <div>
                    <p className="text-sm text-muted-foreground mb-1">Notas</p>
                    <p className="text-sm">{payment.notes}</p>
                  </div>
                </>
              )}
            </CardContent>
          </Card>

          {/* Payment Breakdown */}
          {(payment.principal_applied || payment.interest_applied || payment.late_fee_applied) && (
            <Card>
              <CardHeader>
                <CardTitle className="text-base">Distribución del Pago</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Capital aplicado:</span>
                    <span className="font-medium">{formatCurrency(payment.principal_applied || 0)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Interés aplicado:</span>
                    <span className="font-medium">{formatCurrency(payment.interest_applied || 0)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Mora aplicada:</span>
                    <span className="font-medium">{formatCurrency(payment.late_fee_applied || 0)}</span>
                  </div>
                  <Separator />
                  <div className="flex justify-between font-bold">
                    <span>Total:</span>
                    <span>{formatCurrency(payment.amount)}</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Reversal Info */}
          {payment.status === 'reversed' && payment.reversal_reason && (
            <Card className="border-destructive">
              <CardHeader className="pb-2">
                <CardTitle className="text-base text-destructive flex items-center gap-2">
                  <AlertTriangle className="h-4 w-4" />
                  Pago Revertido
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                <div>
                  <p className="text-sm text-muted-foreground">Razón:</p>
                  <p className="text-sm">{payment.reversal_reason}</p>
                </div>
                {payment.reversed_at && (
                  <div>
                    <p className="text-sm text-muted-foreground">Fecha de reversión:</p>
                    <p className="text-sm">{formatDateTime(payment.reversed_at)}</p>
                  </div>
                )}
                {payment.reversed_by && (
                  <div>
                    <p className="text-sm text-muted-foreground">Revertido por:</p>
                    <p className="text-sm">{payment.reversed_by.first_name} {payment.reversed_by.last_name}</p>
                  </div>
                )}
              </CardContent>
            </Card>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-4">
          {/* Loan Info */}
          {payment.loan && (
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm flex items-center gap-2">
                  <FileText className="h-4 w-4" />
                  Préstamo Asociado
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Número:</span>
                  <Link
                    to={loanRoute(payment.loan.id)}
                    className="font-mono text-primary hover:underline"
                  >
                    {payment.loan.loan_number}
                  </Link>
                </div>
                {payment.loan.principal_amount && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Monto original:</span>
                    <span>{formatCurrency(payment.loan.principal_amount)}</span>
                  </div>
                )}
                {payment.loan.status && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Estado:</span>
                    <span className="capitalize">{payment.loan.status}</span>
                  </div>
                )}
              </CardContent>
            </Card>
          )}

          {/* Customer Info */}
          {payment.customer && (
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm flex items-center gap-2">
                  <User className="h-4 w-4" />
                  Cliente
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Nombre:</span>
                  <Link
                    to={customerRoute(payment.customer.id)}
                    className="text-primary hover:underline"
                  >
                    {payment.customer.first_name} {payment.customer.last_name}
                  </Link>
                </div>
                {payment.customer.phone && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Teléfono:</span>
                    <span>{payment.customer.phone}</span>
                  </div>
                )}
              </CardContent>
            </Card>
          )}

          {/* Audit Info */}
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm flex items-center gap-2">
                <Calendar className="h-4 w-4" />
                Información de Auditoría
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2 text-sm">
              {payment.processed_by && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Procesado por:</span>
                  <span>{payment.processed_by.first_name} {payment.processed_by.last_name}</span>
                </div>
              )}
              {payment.created_at && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Creado:</span>
                  <span>{formatDateTime(payment.created_at)}</span>
                </div>
              )}
              {payment.updated_at && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Actualizado:</span>
                  <span>{formatDateTime(payment.updated_at)}</span>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>

      <ReversePaymentDialog
        open={reverseDialogOpen}
        onOpenChange={setReverseDialogOpen}
        payment={payment}
        onConfirm={handleReverseConfirm}
        isLoading={reverseMutation.isPending}
      />
    </div>
  )
}
