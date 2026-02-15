import { ColumnDef } from '@tanstack/react-table'
import { Link } from 'react-router-dom'
import { MoreHorizontal, Eye, Pencil, Ban, CheckCircle, Trash2 } from 'lucide-react'

import { Customer } from '@/types'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { DataTableColumnHeader } from '@/components/data-table/data-table-column-header'
import { StatusBadge } from '@/components/common/status-badge'
import { formatPhone, formatCurrency } from '@/lib/format'

interface ColumnActions {
  onBlock: (customer: Customer) => void
  onUnblock: (customer: Customer) => void
  onDelete: (customer: Customer) => void
}

export function getCustomerColumns(actions: ColumnActions): ColumnDef<Customer>[] {
  return [
    {
      accessorKey: 'identity_number',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Documento" />,
      cell: ({ row }) => (
        <Link
          to={`/customers/${row.original.id}`}
          className="font-mono text-primary hover:underline"
        >
          {row.original.identity_number}
        </Link>
      ),
    },
    {
      id: 'name',
      accessorFn: (row) => `${row.first_name} ${row.last_name}`,
      header: ({ column }) => <DataTableColumnHeader column={column} title="Nombre" />,
      cell: ({ row }) => (
        <Link
          to={`/customers/${row.original.id}`}
          className="font-medium hover:underline"
        >
          {row.original.first_name} {row.original.last_name}
        </Link>
      ),
    },
    {
      accessorKey: 'phone',
      header: 'Teléfono',
      cell: ({ row }) => formatPhone(row.original.phone),
    },
    {
      accessorKey: 'total_loans',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Préstamos" />,
      cell: ({ row }) => row.original.total_loans,
    },
    {
      accessorKey: 'credit_limit',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Límite Crédito" />,
      cell: ({ row }) => formatCurrency(row.original.credit_limit),
    },
    {
      accessorKey: 'is_blocked',
      header: 'Estado',
      cell: ({ row }) => {
        if (row.original.is_blocked) {
          return <StatusBadge status="blocked" />
        }
        return <StatusBadge status={row.original.is_active ? 'active' : 'inactive'} />
      },
    },
    {
      id: 'actions',
      cell: ({ row }) => {
        const customer = row.original

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
                <Link to={`/customers/${customer.id}`}>
                  <Eye className="mr-2 h-4 w-4" />
                  Ver detalles
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link to={`/customers/${customer.id}/edit`}>
                  <Pencil className="mr-2 h-4 w-4" />
                  Editar
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              {customer.is_blocked ? (
                <DropdownMenuItem onClick={() => actions.onUnblock(customer)}>
                  <CheckCircle className="mr-2 h-4 w-4" />
                  Desbloquear
                </DropdownMenuItem>
              ) : (
                <DropdownMenuItem onClick={() => actions.onBlock(customer)}>
                  <Ban className="mr-2 h-4 w-4" />
                  Bloquear
                </DropdownMenuItem>
              )}
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => actions.onDelete(customer)}
                className="text-destructive focus:text-destructive"
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Eliminar
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
    },
  ]
}
