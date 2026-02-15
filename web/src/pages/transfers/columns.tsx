import { ColumnDef } from '@tanstack/react-table'
import { Link } from 'react-router-dom'
import { MoreHorizontal, Eye, Check, Truck, PackageCheck, X } from 'lucide-react'

import { ItemTransfer, TRANSFER_STATUSES } from '@/types'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Badge } from '@/components/ui/badge'
import { transferRoute } from '@/routes/routes'
import { formatDate } from '@/lib/format'

interface ColumnActions {
  onApprove: (transfer: ItemTransfer) => void
  onShip: (transfer: ItemTransfer) => void
  onReceive: (transfer: ItemTransfer) => void
  onCancel: (transfer: ItemTransfer) => void
}

function getStatusBadge(status: string) {
  const statusConfig = TRANSFER_STATUSES.find((s) => s.value === status)
  if (!statusConfig) return <Badge variant="secondary">{status}</Badge>

  const variants: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
    pending: 'secondary',
    approved: 'outline',
    in_transit: 'default',
    received: 'default',
    cancelled: 'destructive',
  }

  return <Badge variant={variants[status] || 'secondary'}>{statusConfig.label}</Badge>
}

export function getTransferColumns(actions: ColumnActions): ColumnDef<ItemTransfer>[] {
  return [
    {
      accessorKey: 'transfer_number',
      header: 'Número',
      cell: ({ row }) => (
        <Link
          to={transferRoute(row.original.id)}
          className="font-medium text-primary hover:underline font-mono"
        >
          {row.original.transfer_number}
        </Link>
      ),
    },
    {
      accessorKey: 'item',
      header: 'Artículo',
      cell: ({ row }) => (
        <div>
          <p className="font-medium">{row.original.item?.name || '-'}</p>
          <p className="text-sm text-muted-foreground font-mono">
            {row.original.item?.sku}
          </p>
        </div>
      ),
    },
    {
      accessorKey: 'from_branch',
      header: 'Origen',
      cell: ({ row }) => row.original.from_branch?.name || '-',
    },
    {
      accessorKey: 'to_branch',
      header: 'Destino',
      cell: ({ row }) => row.original.to_branch?.name || '-',
    },
    {
      accessorKey: 'created_at',
      header: 'Fecha',
      cell: ({ row }) => formatDate(row.original.created_at),
    },
    {
      accessorKey: 'status',
      header: 'Estado',
      cell: ({ row }) => getStatusBadge(row.original.status),
    },
    {
      id: 'actions',
      cell: ({ row }) => {
        const transfer = row.original
        const { status } = transfer

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
                <Link to={transferRoute(transfer.id)}>
                  <Eye className="mr-2 h-4 w-4" />
                  Ver detalles
                </Link>
              </DropdownMenuItem>

              {status === 'pending' && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={() => actions.onApprove(transfer)}>
                    <Check className="mr-2 h-4 w-4" />
                    Aprobar
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onClick={() => actions.onCancel(transfer)}
                    className="text-destructive focus:text-destructive"
                  >
                    <X className="mr-2 h-4 w-4" />
                    Cancelar
                  </DropdownMenuItem>
                </>
              )}

              {status === 'approved' && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={() => actions.onShip(transfer)}>
                    <Truck className="mr-2 h-4 w-4" />
                    Enviar
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onClick={() => actions.onCancel(transfer)}
                    className="text-destructive focus:text-destructive"
                  >
                    <X className="mr-2 h-4 w-4" />
                    Cancelar
                  </DropdownMenuItem>
                </>
              )}

              {status === 'in_transit' && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={() => actions.onReceive(transfer)}>
                    <PackageCheck className="mr-2 h-4 w-4" />
                    Recibir
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
