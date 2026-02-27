import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import {
  Loader2,
  RefreshCw,
  AlertTriangle,
  DollarSign,
  Calendar,
  Percent,
  User,
  Package,
  CreditCard,
  CheckCircle,
  Clock,
  FileText,
} from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Progress } from '@/components/ui/progress'
import { ROUTES, customerRoute, itemRoute } from '@/routes/routes'
import { useLoan, useLoanInstallments, useRenewLoan, useConfiscateLoan } from '@/hooks/use-loans'
import { useExportLoanContract } from '@/hooks/use-reports'
import { formatCurrency, formatDate } from '@/lib/format'
import { LOAN_STATUSES, PAYMENT_PLAN_TYPES, RenewLoanInput } from '@/types'
import { RenewLoanDialog } from './renew-loan-dialog'
import { ConfiscateLoanDialog } from './confiscate-loan-dialog'

export default function LoanDetailPage() {
  const { id } = useParams()
  const loanId = parseInt(id!, 10)

  const { data: loan, isLoading } = useLoan(loanId)
  const { data: installments } = useLoanInstallments(loanId)
  const renewMutation = useRenewLoan()
  const confiscateMutation = useConfiscateLoan()
  const exportContractMutation = useExportLoanContract()

  const [renewDialogOpen, setRenewDialogOpen] = useState(false)
  const [confiscateDialogOpen, setConfiscateDialogOpen] = useState(false)

  const handleRenewConfirm = (values: RenewLoanInput) => {
    renewMutation.mutate(
      { id: loanId, input: values },
      {
        onSuccess: () => {
          setRenewDialogOpen(false)
        },
      }
    )
  }

  const handleConfiscateConfirm = (notes?: string) => {
    confiscateMutation.mutate(
      { id: loanId, notes },
      {
        onSuccess: () => {
          setConfiscateDialogOpen(false)
        },
      }
    )
  }

  if (isLoading) {
    return (
      <div className="flex h-96 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (!loan) {
    return (
      <div className="flex h-96 items-center justify-center">
        <p className="text-muted-foreground">Préstamo no encontrado</p>
      </div>
    )
  }

  const status = LOAN_STATUSES.find((s) => s.value === loan.status)
  const paymentPlan = PAYMENT_PLAN_TYPES.find((p) => p.value === loan.payment_plan_type)
  // Use late_fee_remaining (pending) instead of late_fee_amount (historical total)
  const lateFeeRemaining = loan.late_fee_remaining ?? loan.late_fee_amount ?? 0
  const balance = loan.principal_remaining + loan.interest_remaining + lateFeeRemaining
  const progressPercent = loan.total_amount > 0 ? (loan.amount_paid / loan.total_amount) * 100 : 0
  const canRenew = ['active', 'overdue'].includes(loan.status)
  const canAcceptPayments = !['paid', 'confiscated'].includes(loan.status)

  // Calculate grace period end and days until automatic confiscation
  const calculateGracePeriodInfo = () => {
    if (loan.status !== 'overdue') return { isOutsideGracePeriod: false, daysUntilConfiscation: 0 }

    const dueDate = new Date(loan.due_date)
    const gracePeriodEnd = new Date(dueDate)
    gracePeriodEnd.setDate(gracePeriodEnd.getDate() + loan.grace_period_days)
    // Set to end of day
    gracePeriodEnd.setHours(23, 59, 59, 999)

    const now = new Date()
    const isOutsideGracePeriod = now > gracePeriodEnd
    const diffTime = gracePeriodEnd.getTime() - now.getTime()
    const daysUntilConfiscation = Math.max(0, Math.ceil(diffTime / (1000 * 60 * 60 * 24)))

    return { isOutsideGracePeriod, daysUntilConfiscation }
  }

  const gracePeriodInfo = calculateGracePeriodInfo()
  // Only allow manual confiscation as backup if grace period has passed
  const canConfiscate = loan.status === 'overdue' && gracePeriodInfo.isOutsideGracePeriod

  const colorMap: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
    green: 'default',
    blue: 'secondary',
    purple: 'secondary',
    red: 'destructive',
    gray: 'outline',
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title={`Préstamo ${loan.loan_number}`}
        description={`Creado el ${formatDate(loan.created_at)}`}
        backUrl={ROUTES.LOANS}
        actions={
          <div className="flex gap-2">
            <Button
              variant="outline"
              onClick={() => exportContractMutation.mutate(loanId)}
              disabled={exportContractMutation.isPending}
            >
              {exportContractMutation.isPending ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <FileText className="mr-2 h-4 w-4" />
              )}
              Contrato
            </Button>
            {canRenew && (
              <Button variant="outline" onClick={() => setRenewDialogOpen(true)}>
                <RefreshCw className="mr-2 h-4 w-4" />
                Renovar
              </Button>
            )}
            {canConfiscate && (
              <Button variant="destructive" onClick={() => setConfiscateDialogOpen(true)}>
                <AlertTriangle className="mr-2 h-4 w-4" />
                Confiscar
              </Button>
            )}
            {canAcceptPayments && (
              <Button asChild>
                <Link to={`/payments/new?loan_id=${loanId}`}>
                  <CreditCard className="mr-2 h-4 w-4" />
                  Registrar Pago
                </Link>
              </Button>
            )}
          </div>
        }
      />

      {/* Status Alert */}
      {loan.status === 'overdue' && (
        <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-4">
          <div className="flex items-center gap-2 text-destructive">
            <AlertTriangle className="h-5 w-5" />
            <span className="font-medium">
              Préstamo vencido
            </span>
          </div>
          {gracePeriodInfo.daysUntilConfiscation > 0 ? (
            <p className="mt-2 text-sm text-amber-600 dark:text-amber-500">
              <Clock className="inline h-4 w-4 mr-1" />
              Confiscación automática en {gracePeriodInfo.daysUntilConfiscation} día(s)
            </p>
          ) : gracePeriodInfo.isOutsideGracePeriod && loan.status === 'overdue' ? (
            <p className="mt-2 text-sm text-destructive">
              <AlertTriangle className="inline h-4 w-4 mr-1" />
              Periodo de gracia terminado - Elegible para confiscación
            </p>
          ) : null}
        </div>
      )}

      {loan.status === 'confiscated' && (
        <div className="rounded-lg border border-muted bg-muted/50 p-4">
          <div className="flex items-center gap-2">
            <AlertTriangle className="h-5 w-5" />
            <span className="font-medium">Préstamo confiscado</span>
          </div>
          <p className="mt-1 text-sm text-muted-foreground">
            El artículo ha sido confiscado y está disponible para venta
          </p>
        </div>
      )}

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Estado</CardTitle>
          </CardHeader>
          <CardContent>
            <Badge variant={colorMap[status?.color || 'gray'] || 'outline'} className="text-base">
              {status?.label || loan.status}
            </Badge>
            {loan.renewal_count > 0 && (
              <p className="text-xs text-muted-foreground mt-1">
                {loan.renewal_count} renovación(es)
              </p>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Monto Préstamo</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(loan.loan_amount)}</div>
            <p className="text-xs text-muted-foreground">
              Total: {formatCurrency(loan.total_amount)}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Saldo Pendiente</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(balance)}</div>
            <div className="mt-2">
              <Progress value={progressPercent} className="h-2" />
              <p className="text-xs text-muted-foreground mt-1">
                {progressPercent.toFixed(0)}% pagado
              </p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Vencimiento</CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className={`text-2xl font-bold ${loan.days_overdue > 0 ? 'text-destructive' : ''}`}>
              {formatDate(loan.due_date)}
            </div>
            {loan.next_payment_due_date && (
              <p className="text-xs text-muted-foreground">
                Próximo pago: {formatDate(loan.next_payment_due_date)}
              </p>
            )}
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Loan Details */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Percent className="h-5 w-5" />
              Detalles del Préstamo
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <p className="text-muted-foreground">Tipo de Pago</p>
                <p className="font-medium">{paymentPlan?.label}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Tasa de Interés</p>
                <p className="font-medium">{loan.interest_rate}%</p>
              </div>
              <div>
                <p className="text-muted-foreground">Plazo</p>
                <p className="font-medium">{loan.loan_term_days} días</p>
              </div>
              <div>
                <p className="text-muted-foreground">Período de Gracia</p>
                <p className="font-medium">{loan.grace_period_days} días</p>
              </div>
              <div>
                <p className="text-muted-foreground">Capital</p>
                <p className="font-medium">{formatCurrency(loan.loan_amount)}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Interés</p>
                <p className="font-medium">{formatCurrency(loan.interest_amount)}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Mora</p>
                <p className="font-medium">
                  {formatCurrency(loan.late_fee_amount)}
                  {loan.status === 'paid' && loan.late_fee_amount > 0 && (
                    <span className="text-xs text-muted-foreground ml-1">(pagada)</span>
                  )}
                  {loan.status === 'overdue' && loan.late_fee_remaining !== undefined &&
                   loan.late_fee_remaining < loan.late_fee_amount && (
                    <span className="text-xs text-muted-foreground ml-1">
                      (pendiente: {formatCurrency(loan.late_fee_remaining)})
                    </span>
                  )}
                </p>
              </div>
              <div>
                <p className="text-muted-foreground">Total Pagado</p>
                <p className="font-medium">{formatCurrency(loan.amount_paid)}</p>
              </div>
            </div>

            {loan.notes && (
              <div className="pt-4 border-t">
                <p className="text-muted-foreground text-sm">Notas</p>
                <p className="whitespace-pre-wrap">{loan.notes}</p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Related Info */}
        <Card>
          <CardHeader>
            <CardTitle>Información Relacionada</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {loan.customer && (
              <div className="flex items-center gap-3 p-3 rounded-lg bg-muted/50">
                <User className="h-8 w-8 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Cliente</p>
                  <Link
                    to={customerRoute(loan.customer.id)}
                    className="font-medium text-primary hover:underline"
                  >
                    {loan.customer.first_name} {loan.customer.last_name}
                  </Link>
                </div>
              </div>
            )}

            {loan.item && (
              <div className="flex items-center gap-3 p-3 rounded-lg bg-muted/50">
                <Package className="h-8 w-8 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Artículo</p>
                  <Link
                    to={itemRoute(loan.item.id)}
                    className="font-medium text-primary hover:underline"
                  >
                    {loan.item.name}
                  </Link>
                  <p className="text-xs text-muted-foreground">{loan.item.sku}</p>
                </div>
              </div>
            )}

            <div className="pt-4 border-t text-sm text-muted-foreground">
              <div className="flex justify-between">
                <span>Fecha inicio:</span>
                <span>{formatDate(loan.start_date)}</span>
              </div>
              {loan.paid_date && (
                <div className="flex justify-between">
                  <span>Fecha de pago:</span>
                  <span>{formatDate(loan.paid_date)}</span>
                </div>
              )}
              {loan.branch && (
                <div className="flex justify-between">
                  <span>Sucursal:</span>
                  <span>{loan.branch.name}</span>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Installments Table */}
      {installments && installments.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Tabla de Cuotas</CardTitle>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Cuota</TableHead>
                  <TableHead>Fecha Vencimiento</TableHead>
                  <TableHead className="text-right">Capital</TableHead>
                  <TableHead className="text-right">Interés</TableHead>
                  <TableHead className="text-right">Total</TableHead>
                  <TableHead className="text-right">Pagado</TableHead>
                  <TableHead>Estado</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {installments.map((inst) => (
                  <TableRow key={inst.id}>
                    <TableCell className="font-medium">#{inst.installment_number}</TableCell>
                    <TableCell>{formatDate(inst.due_date)}</TableCell>
                    <TableCell className="text-right">
                      {formatCurrency(inst.principal_amount)}
                    </TableCell>
                    <TableCell className="text-right">
                      {formatCurrency(inst.interest_amount)}
                    </TableCell>
                    <TableCell className="text-right font-medium">
                      {formatCurrency(inst.total_amount)}
                    </TableCell>
                    <TableCell className="text-right">
                      {formatCurrency(inst.amount_paid)}
                    </TableCell>
                    <TableCell>
                      {inst.is_paid ? (
                        <Badge variant="secondary" className="gap-1">
                          <CheckCircle className="h-3 w-3" />
                          Pagada
                        </Badge>
                      ) : (
                        <Badge variant="outline" className="gap-1">
                          <Clock className="h-3 w-3" />
                          Pendiente
                        </Badge>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      )}

      <RenewLoanDialog
        open={renewDialogOpen}
        onOpenChange={setRenewDialogOpen}
        loan={loan}
        onConfirm={handleRenewConfirm}
        isLoading={renewMutation.isPending}
      />

      <ConfiscateLoanDialog
        open={confiscateDialogOpen}
        onOpenChange={setConfiscateDialogOpen}
        loan={loan}
        onConfirm={handleConfiscateConfirm}
        isLoading={confiscateMutation.isPending}
      />
    </div>
  )
}
