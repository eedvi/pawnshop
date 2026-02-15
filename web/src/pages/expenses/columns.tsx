import { ColumnDef } from '@tanstack/react-table'
import { Link } from 'react-router-dom'
import { MoreHorizontal, Eye, Check, Trash2 } from 'lucide-react'

import { Expense, EXPENSE_STATUSES, getExpenseStatus } from '@/types'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Badge } from '@/components/ui/badge'
import { expenseRoute } from '@/routes/routes'
import { formatCurrency, formatDate } from '@/lib/format'

interface ColumnActions {
  onApprove: (expense: Expense) => void
  onReject: (expense: Expense) => void  // Kept for interface compat, not used
  onDelete: (expense: Expense) => void
}

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

export function getExpenseColumns(actions: ColumnActions): ColumnDef<Expense>[] {
  return [
    {
      accessorKey: 'id',
      header: 'ID',
      cell: ({ row }) => (
        <Link
          to={expenseRoute(row.original.id)}
          className="font-medium text-primary hover:underline"
        >
          #{row.original.id}
        </Link>
      ),
    },
    {
      accessorKey: 'expense_date',
      header: 'Fecha',
      cell: ({ row }) => formatDate(row.original.expense_date),
    },
    {
      accessorKey: 'category',
      header: 'Categoría',
      cell: ({ row }) => row.original.category?.name || '-',
    },
    {
      accessorKey: 'description',
      header: 'Descripción',
      cell: ({ row }) => (
        <span className="max-w-[200px] truncate block">
          {row.original.description}
        </span>
      ),
    },
    {
      accessorKey: 'amount',
      header: 'Monto',
      cell: ({ row }) => (
        <span className="font-medium">{formatCurrency(row.original.amount)}</span>
      ),
    },
    {
      accessorKey: 'branch',
      header: 'Sucursal',
      cell: ({ row }) => row.original.branch?.name || '-',
    },
    {
      id: 'status',
      header: 'Estado',
      cell: ({ row }) => getStatusBadge(getExpenseStatus(row.original)),
    },
    {
      id: 'actions',
      cell: ({ row }) => {
        const expense = row.original
        const isPending = getExpenseStatus(expense) === 'pending'

        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Abrir menú</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem asChild>
                <Link to={expenseRoute(expense.id)}>
                  <Eye className="mr-2 h-4 w-4" />
                  Ver detalles
                </Link>
              </DropdownMenuItem>
              {isPending && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={() => actions.onApprove(expense)}>
                    <Check className="mr-2 h-4 w-4" />
                    Aprobar
                  </DropdownMenuItem>
                </>
              )}
              {isPending && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    onClick={() => actions.onDelete(expense)}
                    className="text-destructive focus:text-destructive"
                  >
                    <Trash2 className="mr-2 h-4 w-4" />
                    Eliminar
                  </DropdownMenuItem>
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
    },
  ]
}
