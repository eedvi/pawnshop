import { ColumnDef } from '@tanstack/react-table'
import { Link } from 'react-router-dom'
import { MoreHorizontal, Eye, Pencil, ShoppingCart, Trash2, Tag } from 'lucide-react'

import { Item, ITEM_STATUSES, ITEM_CONDITIONS } from '@/types'
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
import { formatCurrency } from '@/lib/format'
import { itemRoute, itemEditRoute, customerRoute } from '@/routes/routes'

interface ColumnActions {
  onMarkForSale: (item: Item) => void
  onDelete: (item: Item) => void
}

export function getItemColumns(actions: ColumnActions): ColumnDef<Item>[] {
  return [
    {
      accessorKey: 'sku',
      header: ({ column }) => <DataTableColumnHeader column={column} title="SKU" />,
      cell: ({ row }) => (
        <Link
          to={itemRoute(row.original.id)}
          className="font-mono text-primary hover:underline"
        >
          {row.original.sku}
        </Link>
      ),
    },
    {
      accessorKey: 'name',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Nombre" />,
      cell: ({ row }) => (
        <div className="max-w-[200px]">
          <Link
            to={itemRoute(row.original.id)}
            className="font-medium hover:underline block truncate"
          >
            {row.original.name}
          </Link>
          {row.original.brand && (
            <span className="text-xs text-muted-foreground">{row.original.brand}</span>
          )}
        </div>
      ),
    },
    {
      accessorKey: 'category',
      header: 'Categoría',
      cell: ({ row }) => row.original.category?.name || '-',
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
            className="text-primary hover:underline"
          >
            {customer.first_name} {customer.last_name}
          </Link>
        )
      },
    },
    {
      accessorKey: 'condition',
      header: 'Condición',
      cell: ({ row }) => {
        const condition = ITEM_CONDITIONS.find((c) => c.value === row.original.condition)
        return condition?.label || row.original.condition
      },
    },
    {
      accessorKey: 'appraised_value',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Avalúo" />,
      cell: ({ row }) => formatCurrency(row.original.appraised_value),
    },
    {
      accessorKey: 'loan_value',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Préstamo" />,
      cell: ({ row }) => formatCurrency(row.original.loan_value),
    },
    {
      accessorKey: 'status',
      header: 'Estado',
      cell: ({ row }) => {
        const status = ITEM_STATUSES.find((s) => s.value === row.original.status)
        const colorMap: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
          green: 'default',
          blue: 'default',
          purple: 'secondary',
          orange: 'secondary',
          gray: 'outline',
          red: 'destructive',
          cyan: 'secondary',
          yellow: 'secondary',
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
        const item = row.original
        const canMarkForSale = ['available', 'confiscated'].includes(item.status)

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
                <Link to={itemRoute(item.id)}>
                  <Eye className="mr-2 h-4 w-4" />
                  Ver detalles
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link to={itemEditRoute(item.id)}>
                  <Pencil className="mr-2 h-4 w-4" />
                  Editar
                </Link>
              </DropdownMenuItem>
              {canMarkForSale && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={() => actions.onMarkForSale(item)}>
                    <ShoppingCart className="mr-2 h-4 w-4" />
                    Marcar para venta
                  </DropdownMenuItem>
                </>
              )}
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => actions.onDelete(item)}
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
