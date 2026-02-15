import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import {
  Loader2,
  Edit,
  Trash2,
  Shield,
  Calendar,
  AlertTriangle,
  CheckCircle2,
  XCircle,
} from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import { ROUTES, roleEditRoute } from '@/routes/routes'
import { useRole, useDeleteRole } from '@/hooks/use-roles'
import { formatDateTime } from '@/lib/format'
import { PERMISSION_GROUPS } from '@/types/role'

const GROUP_LABELS: Record<string, string> = {
  customers: 'Clientes',
  items: 'Artículos',
  loans: 'Préstamos',
  payments: 'Pagos',
  sales: 'Ventas',
  cash: 'Caja',
  reports: 'Reportes',
  users: 'Usuarios',
  branches: 'Sucursales',
  categories: 'Categorías',
  roles: 'Roles',
  settings: 'Configuración',
  audit: 'Auditoría',
  notifications: 'Notificaciones',
  expenses: 'Gastos',
  transfers: 'Transferencias',
}

const PERMISSION_LABELS: Record<string, string> = {
  read: 'Ver',
  create: 'Crear',
  update: 'Editar',
  delete: 'Eliminar',
  export: 'Exportar',
  approve: 'Aprobar',
  ship: 'Enviar',
  receive: 'Recibir',
  cancel: 'Cancelar',
  manage: 'Gestionar',
}

function getPermissionLabel(permission: string): string {
  const action = permission.split('.').pop() || permission
  return PERMISSION_LABELS[action] || action
}

export default function RoleDetailPage() {
  const { id } = useParams()
  const roleId = id ? parseInt(id, 10) : 0
  const navigate = useNavigate()

  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  const { data: role, isLoading, error } = useRole(roleId)
  const deleteMutation = useDeleteRole()

  const handleDelete = () => {
    deleteMutation.mutate(roleId, {
      onSuccess: () => navigate(ROUTES.ROLES),
    })
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (error || !role) {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-4">
        <AlertTriangle className="h-12 w-12 text-destructive" />
        <p className="text-muted-foreground">Error al cargar el rol</p>
        <Button asChild variant="outline">
          <Link to={ROUTES.ROLES}>Volver a roles</Link>
        </Button>
      </div>
    )
  }

  return (
    <div>
      <PageHeader
        title={role.display_name}
        description="Detalles del rol"
        backUrl={ROUTES.ROLES}
        actions={
          !role.is_system && (
            <div className="flex gap-2">
              <Button asChild variant="outline">
                <Link to={roleEditRoute(roleId)}>
                  <Edit className="mr-2 h-4 w-4" />
                  Editar
                </Link>
              </Button>
              <Button variant="destructive" onClick={() => setDeleteDialogOpen(true)}>
                <Trash2 className="mr-2 h-4 w-4" />
                Eliminar
              </Button>
            </div>
          )
        }
      />

      <div className="grid gap-6 md:grid-cols-3">
        {/* Main Info */}
        <div className="md:col-span-2 space-y-6">
          {/* Role Info */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Shield className="h-5 w-5" />
                Información del Rol
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <p className="text-sm text-muted-foreground">Nombre</p>
                  <p className="font-medium">{role.display_name}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Código</p>
                  <p className="font-mono">{role.name}</p>
                </div>
              </div>

              {role.description && (
                <>
                  <Separator />
                  <div>
                    <p className="text-sm text-muted-foreground">Descripción</p>
                    <p className="mt-1">{role.description}</p>
                  </div>
                </>
              )}

              <Separator />

              <div className="flex items-center gap-2">
                <p className="text-sm text-muted-foreground">Tipo:</p>
                <Badge variant={role.is_system ? 'outline' : 'default'}>
                  {role.is_system ? 'Sistema' : 'Personalizado'}
                </Badge>
                {role.is_system && (
                  <span className="text-xs text-muted-foreground">
                    (No se puede modificar ni eliminar)
                  </span>
                )}
              </div>
            </CardContent>
          </Card>

          {/* Permissions */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Permisos</CardTitle>
                <Badge variant="secondary">{role.permissions.length} permisos</Badge>
              </div>
            </CardHeader>
            <CardContent>
              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {Object.entries(PERMISSION_GROUPS).map(([groupKey, permissions]) => {
                  const hasAny = permissions.some((p) => role.permissions.includes(p))
                  if (!hasAny) return null

                  return (
                    <div key={groupKey} className="space-y-2">
                      <h4 className="font-medium text-sm">
                        {GROUP_LABELS[groupKey] || groupKey}
                      </h4>
                      <div className="space-y-1">
                        {permissions.map((permission) => {
                          const hasPermission = role.permissions.includes(permission)
                          return (
                            <div
                              key={permission}
                              className="flex items-center gap-2 text-sm"
                            >
                              {hasPermission ? (
                                <CheckCircle2 className="h-4 w-4 text-green-600" />
                              ) : (
                                <XCircle className="h-4 w-4 text-muted-foreground/30" />
                              )}
                              <span
                                className={
                                  hasPermission
                                    ? ''
                                    : 'text-muted-foreground/50 line-through'
                                }
                              >
                                {getPermissionLabel(permission)}
                              </span>
                            </div>
                          )
                        })}
                      </div>
                    </div>
                  )
                })}
              </div>

              {role.permissions.length === 0 && (
                <p className="text-center text-muted-foreground py-4">
                  Este rol no tiene permisos asignados
                </p>
              )}
            </CardContent>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-4">
          {/* Quick Stats */}
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm">Resumen</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Total permisos:</span>
                <span className="font-medium">{role.permissions.length}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Módulos:</span>
                <span className="font-medium">
                  {
                    Object.entries(PERMISSION_GROUPS).filter(([, perms]) =>
                      perms.some((p) => role.permissions.includes(p))
                    ).length
                  }
                </span>
              </div>
            </CardContent>
          </Card>

          {/* Audit Info */}
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm flex items-center gap-2">
                <Calendar className="h-4 w-4" />
                Información de Auditoría
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Creado:</span>
                <span>{formatDateTime(role.created_at)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Actualizado:</span>
                <span>{formatDateTime(role.updated_at)}</span>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      <ConfirmDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        title="Eliminar Rol"
        description={`¿Está seguro de eliminar el rol "${role.display_name}"? Los usuarios con este rol perderán sus permisos.`}
        confirmText="Eliminar"
        variant="destructive"
        onConfirm={handleDelete}
        isLoading={deleteMutation.isPending}
      />
    </div>
  )
}
