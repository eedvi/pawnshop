import { ColumnDef } from '@tanstack/react-table'
import { Eye, ExternalLink } from 'lucide-react'
import { Link } from 'react-router-dom'

import { AuditLog, AUDIT_ACTIONS } from '@/types'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { formatDateTime } from '@/lib/format'
import {
  customerRoute,
  itemRoute,
  loanRoute,
  paymentRoute,
  saleRoute,
  userRoute,
  branchRoute,
  expenseRoute,
  transferRoute,
  roleRoute,
} from '@/routes/routes'

interface AuditColumnOptions {
  onViewDetails: (log: AuditLog) => void
}

function getActionBadgeVariant(action: string): 'default' | 'secondary' | 'destructive' | 'outline' {
  switch (action) {
    case 'create':
      return 'default'
    case 'update':
      return 'secondary'
    case 'delete':
      return 'destructive'
    case 'login':
    case 'logout':
      return 'outline'
    case 'approve':
      return 'default'
    case 'reject':
      return 'destructive'
    default:
      return 'secondary'
  }
}

const ENTITY_TYPE_LABELS: Record<string, string> = {
  customer: 'Cliente',
  item: 'Artículo',
  loan: 'Préstamo',
  payment: 'Pago',
  sale: 'Venta',
  user: 'Usuario',
  branch: 'Sucursal',
  category: 'Categoría',
  role: 'Rol',
  cash_register: 'Caja',
  cash_session: 'Sesión de Caja',
  cash_movement: 'Movimiento de Caja',
  transfer: 'Transferencia',
  expense: 'Gasto',
  notification: 'Notificación',
  setting: 'Configuración',
}

// Helper function to get route for entity
function getEntityRoute(entityType: string, entityId: number): string | null {
  const routes: Record<string, (id: number) => string> = {
    customer: customerRoute,
    item: itemRoute,
    loan: loanRoute,
    payment: paymentRoute,
    sale: saleRoute,
    user: userRoute,
    branch: branchRoute,
    expense: expenseRoute,
    transfer: transferRoute,
    role: roleRoute,
  }

  const routeFn = routes[entityType]
  return routeFn ? routeFn(entityId) : null
}

export function getAuditColumns(options: AuditColumnOptions): ColumnDef<AuditLog>[] {
  const { onViewDetails } = options

  return [
    {
      accessorKey: 'created_at',
      header: 'Fecha/Hora',
      cell: ({ row }) => (
        <span className="text-sm">{formatDateTime(row.original.created_at)}</span>
      ),
    },
    {
      accessorKey: 'user_name',
      header: 'Usuario',
      cell: ({ row }) => (
        <span className="font-medium">{row.original.user_name || 'Sistema'}</span>
      ),
    },
    {
      accessorKey: 'action',
      header: 'Acción',
      cell: ({ row }) => {
        const action = AUDIT_ACTIONS.find((a) => a.value === row.original.action)
        return (
          <Badge variant={getActionBadgeVariant(row.original.action)}>
            {action?.label || row.original.action}
          </Badge>
        )
      },
    },
    {
      accessorKey: 'entity_type',
      header: 'Entidad',
      cell: ({ row }) => {
        const label = ENTITY_TYPE_LABELS[row.original.entity_type] || row.original.entity_type
        const route = row.original.entity_id
          ? getEntityRoute(row.original.entity_type, row.original.entity_id)
          : null

        if (route) {
          return (
            <Link
              to={route}
              className="group flex items-start gap-1 hover:text-primary transition-colors"
            >
              <div>
                <div className="flex items-center gap-1">
                  <p className="font-medium">{label}</p>
                  <ExternalLink className="h-3 w-3 opacity-0 group-hover:opacity-100 transition-opacity" />
                </div>
                <p className="text-xs text-muted-foreground">ID: {row.original.entity_id}</p>
              </div>
            </Link>
          )
        }

        return (
          <div>
            <p>{label}</p>
            {row.original.entity_id && (
              <p className="text-xs text-muted-foreground">ID: {row.original.entity_id}</p>
            )}
          </div>
        )
      },
    },
    {
      accessorKey: 'description',
      header: 'Descripción',
      cell: ({ row }) => (
        <span className="text-sm text-muted-foreground line-clamp-2">
          {row.original.description || '-'}
        </span>
      ),
    },
    {
      accessorKey: 'branch_name',
      header: 'Sucursal',
      cell: ({ row }) => (
        <span className="text-sm">{row.original.branch_name || '-'}</span>
      ),
    },
    {
      accessorKey: 'ip_address',
      header: 'IP',
      cell: ({ row }) => (
        <span className="text-xs font-mono text-muted-foreground">
          {row.original.ip_address || '-'}
        </span>
      ),
    },
    {
      id: 'actions',
      cell: ({ row }) => {
        const log = row.original
        const hasDetails = log.old_values || log.new_values

        return hasDetails ? (
          <Button variant="ghost" size="sm" onClick={() => onViewDetails(log)}>
            <Eye className="h-4 w-4" />
          </Button>
        ) : null
      },
    },
  ]
}
