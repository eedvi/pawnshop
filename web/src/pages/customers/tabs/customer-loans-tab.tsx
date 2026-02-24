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
import { loanRoute } from '@/routes/routes'
import { formatCurrency, formatDate } from '@/lib/format'
import { useLoans } from '@/hooks/use-loans'
import { calculateRemainingBalance } from '@/types/loan'

interface CustomerLoansTabProps {
  customerId: number
}

const LOAN_STATUS_LABELS: Record<string, { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline' }> = {
  active: { label: 'Activo', variant: 'default' },
  paid: { label: 'Pagado', variant: 'secondary' },
  overdue: { label: 'Vencido', variant: 'destructive' },
  defaulted: { label: 'Incobrable', variant: 'destructive' },
  cancelled: { label: 'Cancelado', variant: 'outline' },
}

export function CustomerLoansTab({ customerId }: CustomerLoansTabProps) {
  const { data: loansResponse, isLoading } = useLoans({ customer_id: customerId, per_page: 100 })
  const loans = loansResponse?.data || []

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Préstamos</CardTitle>
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
        <CardTitle>Préstamos del Cliente</CardTitle>
      </CardHeader>
      <CardContent>
        {loans.length === 0 ? (
          <p className="text-center text-sm text-muted-foreground py-8">
            Este cliente no tiene préstamos registrados
          </p>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>No. Préstamo</TableHead>
                <TableHead>Estado</TableHead>
                <TableHead className="text-right">Monto</TableHead>
                <TableHead className="text-right">Pagado</TableHead>
                <TableHead className="text-right">Saldo</TableHead>
                <TableHead>Vencimiento</TableHead>
                <TableHead>Fecha</TableHead>
                <TableHead></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loans.map((loan) => {
                const statusConfig = LOAN_STATUS_LABELS[loan.status] || { label: loan.status, variant: 'outline' as const }
                return (
                  <TableRow key={loan.id}>
                    <TableCell className="font-mono">{loan.loan_number}</TableCell>
                    <TableCell>
                      <Badge variant={statusConfig.variant}>{statusConfig.label}</Badge>
                    </TableCell>
                    <TableCell className="text-right">{formatCurrency(loan.loan_amount)}</TableCell>
                    <TableCell className="text-right">{formatCurrency(loan.amount_paid)}</TableCell>
                    <TableCell className="text-right font-medium">{formatCurrency(calculateRemainingBalance(loan))}</TableCell>
                    <TableCell>{formatDate(loan.due_date)}</TableCell>
                    <TableCell>{formatDate(loan.created_at)}</TableCell>
                    <TableCell>
                      <Link
                        to={loanRoute(loan.id)}
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
