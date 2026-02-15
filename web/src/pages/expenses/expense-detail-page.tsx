import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import {
  Loader2,
  Check,
  Trash2,
  Receipt,
  Building2,
  FolderTree,
  Calendar,
  AlertTriangle,
} from 'lucide-react'

import { EXPENSE_STATUSES, getExpenseStatus } from '@/types'
import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import { ROUTES } from '@/routes/routes'
import {
  useExpense,
  useApproveExpense,
  useDeleteExpense,
} from '@/hooks/use-expenses'
import { formatCurrency, formatDate, formatDateTime } from '@/lib/format'

function getStatusBadge(status: string) {
  const statusConfig = EXPENSE_STATUSES.find((s) => s.value === status)
  if (!statusConfig) return <Badge variant="secondary">{status}</Badge>

  const variants: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
    pending: 'secondary',
    approved: 'default',
    rejected: 'destructive',
  }

  return <Badge variant={variants[status] || 'secondary'}>{statusConfig.label}</Badge>
}

export default function ExpenseDetailPage() {
  const { id } = useParams()
  const expenseId = id ? parseInt(id, 10) : 0
  const navigate = useNavigate()

  const [approveDialogOpen, setApproveDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  const { data: expense, isLoading, error } = useExpense(expenseId)
  const approveMutation = useApproveExpense()
  const deleteMutation = useDeleteExpense()

  const handleApprove = () => {
    approveMutation.mutate(
      { id: expenseId },
      {
        onSuccess: () => setApproveDialogOpen(false),
      }
    )
  }

  const handleDelete = () => {
    deleteMutation.mutate(expenseId, {
      onSuccess: () => navigate(ROUTES.EXPENSES),
    })
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (error || !expense) {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-4">
        <AlertTriangle className="h-12 w-12 text-destructive" />
        <p className="text-muted-foreground">Error al cargar el gasto</p>
        <Button asChild variant="outline">
          <Link to={ROUTES.EXPENSES}>Volver a gastos</Link>
        </Button>
      </div>
    )
  }

  const status = getExpenseStatus(expense)
  const isPending = status === 'pending'

  return (
    <div>
      <PageHeader
        title={`Gasto #${expense.id}`}
        description="Detalles del gasto"
        backUrl={ROUTES.EXPENSES}
        actions={
          isPending && (
            <div className="flex gap-2">
              <Button onClick={() => setApproveDialogOpen(true)}>
                <Check className="mr-2 h-4 w-4" />
                Aprobar
              </Button>
              <Button variant="destructive" onClick={() => setDeleteDialogOpen(true)}>
                <Trash2 className="mr-2 h-4 w-4" />
                Eliminar
              </Button>
            </div>
          )
        }
      />

      <div className="grid gap-6 md:grid-cols-3">
        {/* Main Info */}
        <div className="md:col-span-2 space-y-6">
          {/* Expense Info */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle className="flex items-center gap-2">
                  <Receipt className="h-5 w-5" />
                  Información del Gasto
                </CardTitle>
                {getStatusBadge(status)}
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <p className="text-sm text-muted-foreground">Monto</p>
                  <p className="text-2xl font-bold">{formatCurrency(expense.amount)}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Fecha del Gasto</p>
                  <p className="font-medium">{formatDate(expense.expense_date)}</p>
                </div>
              </div>

              <Separator />

              <div>
                <p className="text-sm text-muted-foreground">Descripción</p>
                <p className="mt-1">{expense.description}</p>
              </div>

              {expense.payment_method && (
                <div>
                  <p className="text-sm text-muted-foreground">Método de Pago</p>
                  <p className="mt-1">{expense.payment_method === 'cash' ? 'Efectivo' : expense.payment_method === 'card' ? 'Tarjeta' : expense.payment_method === 'transfer' ? 'Transferencia' : expense.payment_method}</p>
                </div>
              )}

              {expense.vendor && (
                <div>
                  <p className="text-sm text-muted-foreground">Proveedor</p>
                  <p className="mt-1">{expense.vendor}</p>
                </div>
              )}

              {expense.receipt_number && (
                <div>
                  <p className="text-sm text-muted-foreground">Número de Recibo</p>
                  <p className="font-mono mt-1">{expense.receipt_number}</p>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Approval Info */}
          {status === 'approved' && expense.approved_at && (
            <Card className="border-green-500">
              <CardHeader className="pb-2">
                <CardTitle className="text-base text-green-600 flex items-center gap-2">
                  <Check className="h-4 w-4" />
                  Aprobación
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                <div>
                  <p className="text-sm text-muted-foreground">Fecha de aprobación</p>
                  <p>{formatDateTime(expense.approved_at)}</p>
                </div>
              </CardContent>
            </Card>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-4">
          {/* Category */}
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm flex items-center gap-2">
                <FolderTree className="h-4 w-4" />
                Categoría
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="font-medium">{expense.category?.name || 'Sin categoría'}</p>
              {expense.category?.code && (
                <p className="text-sm text-muted-foreground font-mono">
                  {expense.category.code}
                </p>
              )}
            </CardContent>
          </Card>

          {/* Branch */}
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm flex items-center gap-2">
                <Building2 className="h-4 w-4" />
                Sucursal
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="font-medium">{expense.branch?.name || '-'}</p>
              {expense.branch?.address && (
                <p className="text-sm text-muted-foreground">{expense.branch.address}</p>
              )}
            </CardContent>
          </Card>

          {/* Audit Info */}
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm flex items-center gap-2">
                <Calendar className="h-4 w-4" />
                Información de Auditoría
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Creado:</span>
                <span>{formatDateTime(expense.created_at)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Actualizado:</span>
                <span>{formatDateTime(expense.updated_at)}</span>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      <ConfirmDialog
        open={approveDialogOpen}
        onOpenChange={setApproveDialogOpen}
        title="Aprobar Gasto"
        description={`¿Está seguro de aprobar el gasto #${expense.id} por ${formatCurrency(expense.amount)}?`}
        confirmText="Aprobar"
        onConfirm={handleApprove}
        isLoading={approveMutation.isPending}
      />

      <ConfirmDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        title="Eliminar Gasto"
        description={`¿Está seguro de eliminar el gasto #${expense.id}? Esta acción no se puede deshacer.`}
        confirmText="Eliminar"
        variant="destructive"
        onConfirm={handleDelete}
        isLoading={deleteMutation.isPending}
      />
    </div>
  )
}
