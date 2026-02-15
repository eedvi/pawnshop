import { useState, useMemo } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { Plus, Loader2 } from 'lucide-react'

import { Expense, EXPENSE_STATUSES } from '@/types'
import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { DataTable } from '@/components/data-table/data-table'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { ROUTES } from '@/routes/routes'
import { useBranchStore } from '@/stores/branch-store'
import {
  useExpenses,
  useApproveExpense,
  useDeleteExpense,
} from '@/hooks/use-expenses'
import { getExpenseColumns } from './columns'

export default function ExpenseListPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const { selectedBranch } = useBranchStore()

  const [approveExpense, setApproveExpense] = useState<Expense | null>(null)
  const [deleteExpense, setDeleteExpense] = useState<Expense | null>(null)

  const page = parseInt(searchParams.get('page') || '1')
  const statusFilter = searchParams.get('status')

  // Convert status filter to is_approved boolean for backend
  const isApprovedFilter = statusFilter === 'approved' ? true : statusFilter === 'pending' ? false : undefined

  const { data: expensesResponse, isLoading } = useExpenses({
    page,
    per_page: 10,
    branch_id: selectedBranch?.id,
    approved: isApprovedFilter,
    order_by: 'created_at',
    order: 'desc',
  })

  const approveMutation = useApproveExpense()
  const deleteMutation = useDeleteExpense()

  const expenses = expensesResponse?.data || []
  const pagination = expensesResponse?.meta?.pagination

  const columns = useMemo(
    () =>
      getExpenseColumns({
        onApprove: (expense) => setApproveExpense(expense),
        onReject: () => {}, // Backend doesn't support reject
        onDelete: (expense) => setDeleteExpense(expense),
      }),
    []
  )

  const handleApproveConfirm = () => {
    if (approveExpense) {
      approveMutation.mutate(
        { id: approveExpense.id },
        {
          onSuccess: () => setApproveExpense(null),
        }
      )
    }
  }

  const handleDeleteConfirm = () => {
    if (deleteExpense) {
      deleteMutation.mutate(deleteExpense.id, {
        onSuccess: () => setDeleteExpense(null),
      })
    }
  }

  const handleStatusChange = (value: string) => {
    const params = new URLSearchParams(searchParams)
    if (value === 'all') {
      params.delete('status')
    } else {
      params.set('status', value)
    }
    params.set('page', '1')
    setSearchParams(params)
  }

  const handlePageChange = (newPage: number) => {
    const params = new URLSearchParams(searchParams)
    params.set('page', newPage.toString())
    setSearchParams(params)
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  return (
    <div>
      <PageHeader
        title="Gastos"
        description="Registro y aprobación de gastos operativos"
        actions={
          <Button asChild>
            <Link to={ROUTES.EXPENSE_CREATE}>
              <Plus className="mr-2 h-4 w-4" />
              Nuevo Gasto
            </Link>
          </Button>
        }
      />

      <div className="mb-4 flex gap-4">
        <Select value={statusFilter || 'all'} onValueChange={handleStatusChange}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Filtrar por estado" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Todos los estados</SelectItem>
            {EXPENSE_STATUSES.map((s) => (
              <SelectItem key={s.value} value={s.value}>
                {s.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <DataTable
        columns={columns}
        data={expenses}
        searchPlaceholder="Buscar gastos..."
        searchColumn="description"
        pagination={
          pagination
            ? {
                page: pagination.current_page,
                pageSize: pagination.per_page,
                totalPages: pagination.total_pages,
                totalItems: pagination.total,
                onPageChange: handlePageChange,
              }
            : undefined
        }
      />

      <ConfirmDialog
        open={!!approveExpense}
        onOpenChange={(open) => !open && setApproveExpense(null)}
        title="Aprobar Gasto"
        description={
          approveExpense
            ? `¿Está seguro de aprobar el gasto #${approveExpense.id} por ${new Intl.NumberFormat('es-MX', { style: 'currency', currency: 'MXN' }).format(approveExpense.amount)}?`
            : ''
        }
        confirmText="Aprobar"
        onConfirm={handleApproveConfirm}
        isLoading={approveMutation.isPending}
      />

      <ConfirmDialog
        open={!!deleteExpense}
        onOpenChange={(open) => !open && setDeleteExpense(null)}
        title="Eliminar Gasto"
        description={
          deleteExpense
            ? `¿Está seguro de eliminar el gasto #${deleteExpense.id}? Esta acción no se puede deshacer.`
            : ''
        }
        confirmText="Eliminar"
        variant="destructive"
        onConfirm={handleDeleteConfirm}
        isLoading={deleteMutation.isPending}
      />
    </div>
  )
}
