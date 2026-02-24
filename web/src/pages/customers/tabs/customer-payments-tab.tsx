import { Link } from 'react-router-dom'
import { ExternalLink } from 'lucide-react'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { paymentRoute, loanRoute } from '@/routes/routes'
import { formatCurrency, formatDate } from '@/lib/format'
import { usePayments } from '@/hooks/use-payments'
import { PAYMENT_METHODS } from '@/types/payment'

interface CustomerPaymentsTabProps {
  customerId: number
}

const PAYMENT_STATUS_LABELS: Record<string, { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline' }> = {
  completed: { label: 'Completado', variant: 'default' },
  pending: { label: 'Pendiente', variant: 'secondary' },
  cancelled: { label: 'Cancelado', variant: 'outline' },
  refunded: { label: 'Reembolsado', variant: 'destructive' },
}

export function CustomerPaymentsTab({ customerId }: CustomerPaymentsTabProps) {
  const { data: paymentsResponse, isLoading } = usePayments({ customer_id: customerId, per_page: 100 })
  const payments = paymentsResponse?.data || []

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Pagos</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {Array.from({ length: 3 }).map((_, i) => (
              <Skeleton key={i} className="h-12 w-full" />
            ))}
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Pagos del Cliente</CardTitle>
      </CardHeader>
      <CardContent>
        {payments.length === 0 ? (
          <p className="text-center text-sm text-muted-foreground py-8">
            Este cliente no tiene pagos registrados
          </p>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>No. Pago</TableHead>
                <TableHead>Préstamo</TableHead>
                <TableHead>Método</TableHead>
                <TableHead className="text-right">Monto</TableHead>
                <TableHead>Estado</TableHead>
                <TableHead>Fecha</TableHead>
                <TableHead></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {payments.map((payment) => {
                const statusConfig = PAYMENT_STATUS_LABELS[payment.status] || { label: payment.status, variant: 'outline' as const }
                return (
                  <TableRow key={payment.id}>
                    <TableCell className="font-mono">{payment.payment_number}</TableCell>
                    <TableCell>
                      <Link
                        to={loanRoute(payment.loan_id)}
                        className="text-primary hover:underline"
                      >
                        {payment.loan?.loan_number || `Préstamo #${payment.loan_id}`}
                      </Link>
                    </TableCell>
                    <TableCell>
                      {PAYMENT_METHODS.find(m => m.value === payment.payment_method)?.label || payment.payment_method}
                    </TableCell>
                    <TableCell className="text-right font-medium">{formatCurrency(payment.amount)}</TableCell>
                    <TableCell>
                      <Badge variant={statusConfig.variant}>{statusConfig.label}</Badge>
                    </TableCell>
                    <TableCell>{formatDate(payment.created_at)}</TableCell>
                    <TableCell>
                      <Link
                        to={paymentRoute(payment.id)}
                        className="inline-flex items-center text-primary hover:underline"
                      >
                        <ExternalLink className="h-4 w-4" />
                      </Link>
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>
        )}
      </CardContent>
    </Card>
  )
}
