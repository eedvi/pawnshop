import { ColumnDef } from '@tanstack/react-table'
import { Link } from 'react-router-dom'
import { MoreHorizontal, Eye, RefreshCw, AlertTriangle } from 'lucide-react'

import { Loan, LOAN_STATUSES } from '@/types'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { DataTableColumnHeader } from '@/components/data-table/data-table-column-header'
import { Badge } from '@/components/ui/badge'
import { formatCurrency, formatDate } from '@/lib/format'
import { loanRoute, customerRoute, itemRoute } from '@/routes/routes'

interface ColumnActions {
  onRenew: (loan: Loan) => void
  onConfiscate: (loan: Loan) => void
}

export function getLoanColumns(actions: ColumnActions): ColumnDef<Loan>[] {
  return [
    {
      accessorKey: 'loan_number',
      header: ({ column }) => <DataTableColumnHeader column={column} title="No. Préstamo" />,
      cell: ({ row }) => (
        <Link
          to={loanRoute(row.original.id)}
          className="font-mono text-primary hover:underline"
        >
          {row.original.loan_number}
        </Link>
      ),
    },
    {
      accessorKey: 'customer',
      header: 'Cliente',
      cell: ({ row }) => {
        const customer = row.original.customer
        if (!customer) return '-'
        return (
          <Link
            to={customerRoute(customer.id)}
            className="hover:underline"
          >
            {customer.first_name} {customer.last_name}
          </Link>
        )
      },
    },
    {
      accessorKey: 'item',
      header: 'Artículo',
      cell: ({ row }) => {
        const item = row.original.item
        if (!item) return '-'
        return (
          <Link
            to={itemRoute(item.id)}
            className="hover:underline"
          >
            {item.name}
          </Link>
        )
      },
    },
    {
      accessorKey: 'loan_amount',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Monto" />,
      cell: ({ row }) => formatCurrency(row.original.loan_amount),
    },
    {
      accessorKey: 'total_amount',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Total" />,
      cell: ({ row }) => formatCurrency(row.original.total_amount),
    },
    {
      accessorKey: 'amount_paid',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Pagado" />,
      cell: ({ row }) => formatCurrency(row.original.amount_paid),
    },
    {
      accessorKey: 'due_date',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Vencimiento" />,
      cell: ({ row }) => {
        const isOverdue = row.original.days_overdue > 0
        return (
          <span className={isOverdue ? 'text-destructive font-medium' : ''}>
            {formatDate(row.original.due_date)}
            {isOverdue && ` (${row.original.days_overdue}d)`}
          </span>
        )
      },
    },
    {
      accessorKey: 'status',
      header: 'Estado',
      cell: ({ row }) => {
        const status = LOAN_STATUSES.find((s) => s.value === row.original.status)
        const colorMap: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
          green: 'default',
          blue: 'secondary',
          purple: 'secondary',
          red: 'destructive',
          gray: 'outline',
        }
        return (
          <Badge variant={colorMap[status?.color || 'gray'] || 'outline'}>
            {status?.label || row.original.status}
          </Badge>
        )
      },
    },
    {
      id: 'actions',
      cell: ({ row }) => {
        const loan = row.original
        const canRenew = ['active', 'overdue'].includes(loan.status)
        const canConfiscate = ['overdue', 'defaulted'].includes(loan.status)

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
                <Link to={loanRoute(loan.id)}>
                  <Eye className="mr-2 h-4 w-4" />
                  Ver detalles
                </Link>
              </DropdownMenuItem>
              {canRenew && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={() => actions.onRenew(loan)}>
                    <RefreshCw className="mr-2 h-4 w-4" />
                    Renovar
                  </DropdownMenuItem>
                </>
              )}
              {canConfiscate && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    onClick={() => actions.onConfiscate(loan)}
                    className="text-destructive focus:text-destructive"
                  >
                    <AlertTriangle className="mr-2 h-4 w-4" />
                    Confiscar
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
