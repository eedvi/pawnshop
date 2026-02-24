import { Activity, Users, FileText, AlertTriangle, TrendingUp } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { AUDIT_ACTIONS } from '@/types'
import type { AuditStats } from '@/types'
import { formatDateTime } from '@/lib/format'

interface AuditDashboardProps {
  stats: AuditStats
  isLoading: boolean
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

export function AuditDashboard({ stats, isLoading }: AuditDashboardProps) {
  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-6">
        {[...Array(4)].map((_, i) => (
          <Card key={i}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <Skeleton className="h-4 w-24" />
              <Skeleton className="h-4 w-4 rounded" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-8 w-16 mb-1" />
              <Skeleton className="h-3 w-32" />
            </CardContent>
          </Card>
        ))}
      </div>
    )
  }

  const topActions = Object.entries(stats.actions_by_type)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 3)

  const topEntities = Object.entries(stats.actions_by_entity)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 3)

  return (
    <div className="space-y-6 mb-6">
      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Acciones</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.total_actions.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">
              Registros en el periodo seleccionado
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Usuarios Activos</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.active_users}</div>
            <p className="text-xs text-muted-foreground">
              Usuarios con actividad registrada
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Acciones Críticas</CardTitle>
            <AlertTriangle className="h-4 w-4 text-destructive" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-destructive">
              {(stats.actions_by_type.delete || 0) + (stats.actions_by_type.reject || 0)}
            </div>
            <p className="text-xs text-muted-foreground">
              Eliminaciones y rechazos
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Entidades Afectadas</CardTitle>
            <FileText className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {Object.keys(stats.actions_by_entity).length}
            </div>
            <p className="text-xs text-muted-foreground">
              Tipos de entidades modificadas
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Details Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {/* Top Actions */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base flex items-center gap-2">
              <TrendingUp className="h-4 w-4" />
              Acciones Más Frecuentes
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {topActions.map(([action, count]) => {
                const actionLabel = AUDIT_ACTIONS.find(a => a.value === action)?.label || action
                const percentage = ((count / stats.total_actions) * 100).toFixed(1)
                return (
                  <div key={action} className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Badge variant={getActionBadgeVariant(action)}>{actionLabel}</Badge>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-medium">{count}</span>
                      <span className="text-xs text-muted-foreground">({percentage}%)</span>
                    </div>
                  </div>
                )
              })}
            </div>
          </CardContent>
        </Card>

        {/* Top Entities */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base flex items-center gap-2">
              <FileText className="h-4 w-4" />
              Entidades Más Modificadas
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {topEntities.map(([entity, count]) => {
                const entityLabel = ENTITY_TYPE_LABELS[entity] || entity
                const percentage = ((count / stats.total_actions) * 100).toFixed(1)
                return (
                  <div key={entity} className="flex items-center justify-between">
                    <span className="text-sm font-medium">{entityLabel}</span>
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-medium">{count}</span>
                      <span className="text-xs text-muted-foreground">({percentage}%)</span>
                    </div>
                  </div>
                )
              })}
            </div>
          </CardContent>
        </Card>

        {/* Top Users */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base flex items-center gap-2">
              <Users className="h-4 w-4" />
              Usuarios Más Activos
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {stats.top_users.slice(0, 5).map((user, index) => {
                const percentage = ((user.count / stats.total_actions) * 100).toFixed(1)
                return (
                  <div key={user.user_id} className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <Badge variant="outline" className="w-6 h-6 p-0 flex items-center justify-center">
                        {index + 1}
                      </Badge>
                      <span className="text-sm font-medium truncate">{user.user_name}</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-medium">{user.count}</span>
                      <span className="text-xs text-muted-foreground">({percentage}%)</span>
                    </div>
                  </div>
                )
              })}
              {stats.top_users.length === 0 && (
                <p className="text-sm text-muted-foreground text-center py-2">
                  No hay datos disponibles
                </p>
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Recent Critical Actions */}
      {stats.recent_critical.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base flex items-center gap-2">
              <AlertTriangle className="h-4 w-4 text-destructive" />
              Acciones Críticas Recientes
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {stats.recent_critical.slice(0, 5).map((log) => {
                const actionLabel = AUDIT_ACTIONS.find(a => a.value === log.action)?.label || log.action
                const entityLabel = ENTITY_TYPE_LABELS[log.entity_type] || log.entity_type
                return (
                  <div
                    key={log.id}
                    className="flex items-center justify-between p-2 rounded-lg border bg-muted/50"
                  >
                    <div className="flex items-center gap-3 flex-1 min-w-0">
                      <Badge variant={getActionBadgeVariant(log.action)}>{actionLabel}</Badge>
                      <div className="flex-1 min-w-0">
                        <p className="text-sm font-medium truncate">
                          {entityLabel}
                          {log.entity_id && <span className="text-muted-foreground"> #{log.entity_id}</span>}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          {log.user_name || 'Sistema'} • {formatDateTime(log.created_at)}
                        </p>
                      </div>
                    </div>
                  </div>
                )
              })}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
