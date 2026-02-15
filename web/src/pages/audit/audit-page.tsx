import { useState, useMemo } from 'react'
import { useSearchParams } from 'react-router-dom'
import { Loader2, Filter, X, Download } from 'lucide-react'
import { format, subDays } from 'date-fns'

import {
  AuditLog,
  AuditAction,
  AUDIT_ACTIONS,
  ENTITY_TYPES,
} from '@/types'
import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { DataTable } from '@/components/data-table/data-table'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { useBranchStore } from '@/stores/branch-store'
import { useBranches } from '@/hooks/use-branches'
import { useUsers } from '@/hooks/use-users'
import { useAuditLogs } from '@/hooks/use-audit'
import { getAuditColumns } from './columns'
import { AuditDetailDialog } from './audit-detail-dialog'

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

const QUICK_DATE_FILTERS = [
  { label: 'Hoy', days: 0 },
  { label: 'Últimos 7 días', days: 7 },
  { label: 'Últimos 30 días', days: 30 },
  { label: 'Últimos 90 días', days: 90 },
]

export default function AuditPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const { selectedBranch } = useBranchStore()

  const [selectedLog, setSelectedLog] = useState<AuditLog | null>(null)
  const [showFilters, setShowFilters] = useState(false)

  // Parse filters from URL
  const page = parseInt(searchParams.get('page') || '1')
  const userId = searchParams.get('user_id') ? parseInt(searchParams.get('user_id')!) : undefined
  const branchId = searchParams.get('branch_id')
    ? parseInt(searchParams.get('branch_id')!)
    : selectedBranch?.id
  const action = searchParams.get('action') as AuditAction | undefined
  const entityType = searchParams.get('entity_type') || undefined
  const entityId = searchParams.get('entity_id') ? parseInt(searchParams.get('entity_id')!) : undefined
  const dateFrom = searchParams.get('date_from') || undefined
  const dateTo = searchParams.get('date_to') || undefined

  const { data: branches } = useBranches()
  const { data: usersResponse } = useUsers({ per_page: 100 })
  const users = usersResponse?.data || []

  const { data: logsResponse, isLoading } = useAuditLogs({
    page,
    per_page: 20,
    user_id: userId,
    branch_id: branchId,
    action,
    entity_type: entityType,
    entity_id: entityId,
    date_from: dateFrom,
    date_to: dateTo,
    order_by: 'created_at',
    order: 'desc',
  })

  const logs = logsResponse?.data || []
  const pagination = logsResponse?.meta?.pagination

  const columns = useMemo(
    () =>
      getAuditColumns({
        onViewDetails: (log) => setSelectedLog(log),
      }),
    []
  )

  // Count active filters
  const activeFilters = [
    userId,
    branchId !== selectedBranch?.id ? branchId : undefined,
    action,
    entityType,
    entityId,
    dateFrom,
    dateTo,
  ].filter(Boolean).length

  const updateFilter = (key: string, value: string | undefined) => {
    const params = new URLSearchParams(searchParams)
    if (value) {
      params.set(key, value)
    } else {
      params.delete(key)
    }
    params.set('page', '1')
    setSearchParams(params)
  }

  const clearFilters = () => {
    setSearchParams({ page: '1' })
  }

  const handleQuickDateFilter = (days: number) => {
    const params = new URLSearchParams(searchParams)
    if (days === 0) {
      params.set('date_from', format(new Date(), 'yyyy-MM-dd'))
      params.set('date_to', format(new Date(), 'yyyy-MM-dd'))
    } else {
      params.set('date_from', format(subDays(new Date(), days), 'yyyy-MM-dd'))
      params.set('date_to', format(new Date(), 'yyyy-MM-dd'))
    }
    params.set('page', '1')
    setSearchParams(params)
  }

  const handlePageChange = (newPage: number) => {
    const params = new URLSearchParams(searchParams)
    params.set('page', newPage.toString())
    setSearchParams(params)
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  return (
    <div>
      <PageHeader
        title="Auditoría"
        description="Registro de actividades del sistema"
        actions={
          <div className="flex gap-2">
            <Button
              variant="outline"
              onClick={() => setShowFilters(!showFilters)}
            >
              <Filter className="mr-2 h-4 w-4" />
              Filtros
              {activeFilters > 0 && (
                <Badge variant="secondary" className="ml-2">
                  {activeFilters}
                </Badge>
              )}
            </Button>
          </div>
        }
      />

      {/* Filters Panel */}
      {showFilters && (
        <Card className="mb-4">
          <CardContent className="pt-4">
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
              {/* Quick Date Filters */}
              <div className="sm:col-span-2 lg:col-span-4">
                <Label className="mb-2 block">Filtros Rápidos</Label>
                <div className="flex flex-wrap gap-2">
                  {QUICK_DATE_FILTERS.map((filter) => (
                    <Button
                      key={filter.days}
                      variant="outline"
                      size="sm"
                      onClick={() => handleQuickDateFilter(filter.days)}
                    >
                      {filter.label}
                    </Button>
                  ))}
                </div>
              </div>

              {/* Date Range */}
              <div>
                <Label htmlFor="date_from">Fecha Desde</Label>
                <Input
                  id="date_from"
                  type="date"
                  value={dateFrom || ''}
                  onChange={(e) => updateFilter('date_from', e.target.value || undefined)}
                />
              </div>
              <div>
                <Label htmlFor="date_to">Fecha Hasta</Label>
                <Input
                  id="date_to"
                  type="date"
                  value={dateTo || ''}
                  onChange={(e) => updateFilter('date_to', e.target.value || undefined)}
                />
              </div>

              {/* User Filter */}
              <div>
                <Label htmlFor="user">Usuario</Label>
                <Select
                  value={userId?.toString() || 'all'}
                  onValueChange={(value) => updateFilter('user_id', value === 'all' ? undefined : value)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Todos los usuarios" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">Todos los usuarios</SelectItem>
                    {users.map((user) => (
                      <SelectItem key={user.id} value={user.id.toString()}>
                        {user.first_name} {user.last_name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              {/* Branch Filter */}
              <div>
                <Label htmlFor="branch">Sucursal</Label>
                <Select
                  value={branchId?.toString() || 'all'}
                  onValueChange={(value) => updateFilter('branch_id', value === 'all' ? undefined : value)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Todas las sucursales" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">Todas las sucursales</SelectItem>
                    {branches?.data?.map((branch) => (
                      <SelectItem key={branch.id} value={branch.id.toString()}>
                        {branch.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              {/* Action Filter */}
              <div>
                <Label htmlFor="action">Acción</Label>
                <Select
                  value={action || 'all'}
                  onValueChange={(value) => updateFilter('action', value === 'all' ? undefined : value)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Todas las acciones" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">Todas las acciones</SelectItem>
                    {AUDIT_ACTIONS.map((a) => (
                      <SelectItem key={a.value} value={a.value}>
                        {a.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              {/* Entity Type Filter */}
              <div>
                <Label htmlFor="entity_type">Tipo de Entidad</Label>
                <Select
                  value={entityType || 'all'}
                  onValueChange={(value) => updateFilter('entity_type', value === 'all' ? undefined : value)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Todas las entidades" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">Todas las entidades</SelectItem>
                    {ENTITY_TYPES.map((type) => (
                      <SelectItem key={type} value={type}>
                        {ENTITY_TYPE_LABELS[type] || type}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              {/* Entity ID Filter */}
              <div>
                <Label htmlFor="entity_id">ID de Entidad</Label>
                <Input
                  id="entity_id"
                  type="number"
                  placeholder="Ej: 123"
                  value={entityId || ''}
                  onChange={(e) => updateFilter('entity_id', e.target.value || undefined)}
                />
              </div>

              {/* Clear Button */}
              <div className="flex items-end">
                <Button
                  variant="outline"
                  onClick={clearFilters}
                  disabled={activeFilters === 0}
                >
                  <X className="mr-2 h-4 w-4" />
                  Limpiar Filtros
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Active Filters Summary */}
      {activeFilters > 0 && !showFilters && (
        <div className="mb-4 flex flex-wrap items-center gap-2">
          <span className="text-sm text-muted-foreground">Filtros activos:</span>
          {dateFrom && (
            <Badge variant="secondary">
              Desde: {dateFrom}
              <button
                className="ml-1 hover:text-destructive"
                onClick={() => updateFilter('date_from', undefined)}
              >
                ×
              </button>
            </Badge>
          )}
          {dateTo && (
            <Badge variant="secondary">
              Hasta: {dateTo}
              <button
                className="ml-1 hover:text-destructive"
                onClick={() => updateFilter('date_to', undefined)}
              >
                ×
              </button>
            </Badge>
          )}
          {action && (
            <Badge variant="secondary">
              Acción: {AUDIT_ACTIONS.find((a) => a.value === action)?.label}
              <button
                className="ml-1 hover:text-destructive"
                onClick={() => updateFilter('action', undefined)}
              >
                ×
              </button>
            </Badge>
          )}
          {entityType && (
            <Badge variant="secondary">
              Entidad: {ENTITY_TYPE_LABELS[entityType] || entityType}
              <button
                className="ml-1 hover:text-destructive"
                onClick={() => updateFilter('entity_type', undefined)}
              >
                ×
              </button>
            </Badge>
          )}
          <Button variant="ghost" size="sm" onClick={clearFilters}>
            Limpiar todos
          </Button>
        </div>
      )}

      {/* Data Table */}
      <DataTable
        columns={columns}
        data={logs}
        pagination={
          pagination
            ? {
                page: pagination.current_page,
                pageSize: pagination.per_page,
                totalPages: pagination.total_pages,
                totalItems: pagination.total,
                onPageChange: handlePageChange,
              }
            : undefined
        }
      />

      {/* Detail Dialog */}
      <AuditDetailDialog
        log={selectedLog}
        open={!!selectedLog}
        onOpenChange={(open) => !open && setSelectedLog(null)}
      />
    </div>
  )
}
