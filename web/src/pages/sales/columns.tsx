import { ColumnDef } from '@tanstack/react-table'
import { Link } from 'react-router-dom'
import { MoreHorizontal, Eye, RotateCcw, XCircle } from 'lucide-react'

import { Sale, SALE_STATUSES, SALE_TYPES, PAYMENT_METHODS } from '@/types'
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
import { saleRoute, itemRoute, customerRoute } from '@/routes/routes'

interface ColumnActions {
  onRefund: (sale: Sale) => void
  onCancel: (sale: Sale) => void
}

export function getSaleColumns(actions: ColumnActions): ColumnDef<Sale>[] {
  return [
    {
      accessorKey: 'sale_number',
      header: ({ column }) => <DataTableColumnHeader column={column} title="No. Venta" />,
      cell: ({ row }) => (
        <Link
          to={saleRoute(row.original.id)}
          className="font-mono text-primary hover:underline"
        >
          {row.original.sale_number}
        </Link>
      ),
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
      accessorKey: 'customer',
      header: 'Cliente',
      cell: ({ row }) => {
        const customer = row.original.customer
        if (!customer) return <span className="text-muted-foreground">Sin cliente</span>
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
      accessorKey: 'sale_type',
      header: 'Tipo',
      cell: ({ row }) => {
        const type = SALE_TYPES.find((t) => t.value === row.original.sale_type)
        return type?.label || row.original.sale_type
      },
    },
    {
      accessorKey: 'final_price',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Monto" />,
      cell: ({ row }) => (
        <span className="font-medium">{formatCurrency(row.original.final_price)}</span>
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
      accessorKey: 'sale_date',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Fecha" />,
      cell: ({ row }) => formatDate(row.original.sale_date),
    },
    {
      accessorKey: 'status',
      header: 'Estado',
      cell: ({ row }) => {
        const status = SALE_STATUSES.find((s) => s.value === row.original.status)
        const colorMap: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
          green: 'default',
          yellow: 'secondary',
          red: 'destructive',
          orange: 'secondary',
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
        const sale = row.original
        const canRefund = sale.status === 'completed'
        const canCancel = sale.status === 'pending'

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
                <Link to={saleRoute(sale.id)}>
                  <Eye className="mr-2 h-4 w-4" />
                  Ver detalles
                </Link>
              </DropdownMenuItem>
              {(canRefund || canCancel) && (
                <>
                  <DropdownMenuSeparator />
                  {canRefund && (
                    <DropdownMenuItem
                      onClick={() => actions.onRefund(sale)}
                      className="text-destructive focus:text-destructive"
                    >
                      <RotateCcw className="mr-2 h-4 w-4" />
                      Reembolsar
                    </DropdownMenuItem>
                  )}
                  {canCancel && (
                    <DropdownMenuItem
                      onClick={() => actions.onCancel(sale)}
                      className="text-destructive focus:text-destructive"
                    >
                      <XCircle className="mr-2 h-4 w-4" />
                      Cancelar
                    </DropdownMenuItem>
                  )}
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
    },
  ]
}
