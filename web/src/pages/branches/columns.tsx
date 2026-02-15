import { ColumnDef } from '@tanstack/react-table'
import { Link } from 'react-router-dom'
import { MoreHorizontal, Eye, Pencil, Power, PowerOff, Trash2 } from 'lucide-react'

import { Branch } from '@/types'
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
import { formatPercent } from '@/lib/format'

interface ColumnActions {
  onActivate: (branch: Branch) => void
  onDeactivate: (branch: Branch) => void
  onDelete: (branch: Branch) => void
}

export function getBranchColumns(actions: ColumnActions): ColumnDef<Branch>[] {
  return [
    {
      accessorKey: 'code',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Código" />,
      cell: ({ row }) => (
        <Link
          to={`/branches/${row.original.id}`}
          className="font-mono font-medium text-primary hover:underline"
        >
          {row.original.code}
        </Link>
      ),
    },
    {
      accessorKey: 'name',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Nombre" />,
      cell: ({ row }) => (
        <Link
          to={`/branches/${row.original.id}`}
          className="font-medium hover:underline"
        >
          {row.original.name}
        </Link>
      ),
    },
    {
      accessorKey: 'phone',
      header: 'Teléfono',
      cell: ({ row }) => row.original.phone || '-',
    },
    {
      accessorKey: 'default_interest_rate',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Tasa Interés" />,
      cell: ({ row }) => formatPercent(row.original.default_interest_rate),
    },
    {
      accessorKey: 'default_loan_term_days',
      header: 'Plazo',
      cell: ({ row }) => `${row.original.default_loan_term_days} días`,
    },
    {
      accessorKey: 'is_active',
      header: 'Estado',
      cell: ({ row }) => (
        <StatusBadge status={row.original.is_active ? 'active' : 'inactive'} />
      ),
    },
    {
      id: 'actions',
      cell: ({ row }) => {
        const branch = row.original

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
                <Link to={`/branches/${branch.id}`}>
                  <Eye className="mr-2 h-4 w-4" />
                  Ver detalles
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link to={`/branches/${branch.id}?edit=true`}>
                  <Pencil className="mr-2 h-4 w-4" />
                  Editar
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              {branch.is_active ? (
                <DropdownMenuItem onClick={() => actions.onDeactivate(branch)}>
                  <PowerOff className="mr-2 h-4 w-4" />
                  Desactivar
                </DropdownMenuItem>
              ) : (
                <DropdownMenuItem onClick={() => actions.onActivate(branch)}>
                  <Power className="mr-2 h-4 w-4" />
                  Activar
                </DropdownMenuItem>
              )}
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => actions.onDelete(branch)}
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
