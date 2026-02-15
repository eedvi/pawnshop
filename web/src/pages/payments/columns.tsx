import { ColumnDef } from '@tanstack/react-table'
import { Link } from 'react-router-dom'
import { MoreHorizontal, Eye, RotateCcw } from 'lucide-react'

import { Payment, PAYMENT_STATUSES, PAYMENT_METHODS } from '@/types'
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
import { paymentRoute, loanRoute, customerRoute } from '@/routes/routes'

interface ColumnActions {
  onReverse: (payment: Payment) => void
}

export function getPaymentColumns(actions: ColumnActions): ColumnDef<Payment>[] {
  return [
    {
      accessorKey: 'payment_number',
      header: ({ column }) => <DataTableColumnHeader column={column} title="No. Pago" />,
      cell: ({ row }) => (
        <Link
          to={paymentRoute(row.original.id)}
          className="font-mono text-primary hover:underline"
        >
          {row.original.payment_number}
        </Link>
      ),
    },
    {
      accessorKey: 'loan',
      header: 'Préstamo',
      cell: ({ row }) => {
        const loan = row.original.loan
        if (!loan) return '-'
        return (
          <Link
            to={loanRoute(loan.id)}
            className="font-mono hover:underline"
          >
            {loan.loan_number}
          </Link>
        )
      },
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
      accessorKey: 'amount',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Monto" />,
      cell: ({ row }) => (
        <span className="font-medium">{formatCurrency(row.original.amount)}</span>
      ),
    },
    {
      accessorKey: 'payment_method',
      header: 'Método',
      cell: ({ row }) => {
        const method = PAYMENT_METHODS.find((m) => m.value === row.original.payment_method)
        return method?.label || row.original.payment_method
      },
    },
    {
      accessorKey: 'payment_date',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Fecha" />,
      cell: ({ row }) => formatDate(row.original.payment_date),
    },
    {
      accessorKey: 'status',
      header: 'Estado',
      cell: ({ row }) => {
        const status = PAYMENT_STATUSES.find((s) => s.value === row.original.status)
        const colorMap: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
          green: 'default',
          yellow: 'secondary',
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
        const payment = row.original
        const canReverse = payment.status === 'completed'

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
                <Link to={paymentRoute(payment.id)}>
                  <Eye className="mr-2 h-4 w-4" />
                  Ver detalles
                </Link>
              </DropdownMenuItem>
              {canReverse && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    onClick={() => actions.onReverse(payment)}
                    className="text-destructive focus:text-destructive"
                  >
                    <RotateCcw className="mr-2 h-4 w-4" />
                    Revertir
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
