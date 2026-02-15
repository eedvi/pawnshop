import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import {
  Loader2,
  Edit,
  Key,
  Power,
  Unlock,
  User as UserIcon,
  Mail,
  Phone,
  Building2,
  Shield,
  Calendar,
  AlertTriangle,
  Lock,
} from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Separator } from '@/components/ui/separator'
import { ROUTES, userEditRoute } from '@/routes/routes'
import {
  useUser,
  useResetPassword,
  useToggleUserActive,
  useUnlockUser,
} from '@/hooks/use-users'
import { formatDateTime } from '@/lib/format'
import { ResetPasswordDialog } from './reset-password-dialog'

export default function UserDetailPage() {
  const { id } = useParams()
  const userId = id ? parseInt(id, 10) : 0

  const [resetDialogOpen, setResetDialogOpen] = useState(false)

  const { data: user, isLoading, error } = useUser(userId)
  const resetPasswordMutation = useResetPassword()
  const toggleActiveMutation = useToggleUserActive()
  const unlockMutation = useUnlockUser()

  const handleResetConfirm = (password: string) => {
    resetPasswordMutation.mutate(
      { id: userId, password },
      {
        onSuccess: () => setResetDialogOpen(false),
      }
    )
  }

  const handleToggleActive = () => {
    if (user && confirm(`¿Está seguro de ${user.is_active ? 'desactivar' : 'activar'} este usuario?`)) {
      toggleActiveMutation.mutate(userId)
    }
  }

  const handleUnlock = () => {
    if (confirm('¿Está seguro de desbloquear este usuario?')) {
      unlockMutation.mutate(userId)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (error || !user) {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-4">
        <AlertTriangle className="h-12 w-12 text-destructive" />
        <p className="text-muted-foreground">Error al cargar el usuario</p>
        <Button asChild variant="outline">
          <Link to={ROUTES.USERS}>Volver a usuarios</Link>
        </Button>
      </div>
    )
  }

  const initials = `${user.first_name[0]}${user.last_name[0]}`.toUpperCase()
  const isLocked = user.locked_until && new Date(user.locked_until) > new Date()

  return (
    <div>
      <PageHeader
        title={`${user.first_name} ${user.last_name}`}
        description="Detalles del usuario"
        backUrl={ROUTES.USERS}
        actions={
          <div className="flex gap-2">
            <Button variant="outline" onClick={() => setResetDialogOpen(true)}>
              <Key className="mr-2 h-4 w-4" />
              Cambiar Contraseña
            </Button>
            <Button asChild variant="outline">
              <Link to={userEditRoute(userId)}>
                <Edit className="mr-2 h-4 w-4" />
                Editar
              </Link>
            </Button>
          </div>
        }
      />

      <div className="grid gap-6 md:grid-cols-3">
        {/* Main Info */}
        <div className="md:col-span-2 space-y-6">
          {/* Profile Card */}
          <Card>
            <CardContent className="pt-6">
              <div className="flex items-start gap-6">
                <Avatar className="h-24 w-24">
                  <AvatarImage src={user.avatar_url} alt={`${user.first_name} ${user.last_name}`} />
                  <AvatarFallback className="text-2xl">{initials}</AvatarFallback>
                </Avatar>
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <h2 className="text-2xl font-bold">
                      {user.first_name} {user.last_name}
                    </h2>
                    {isLocked ? (
                      <Badge variant="destructive" className="flex items-center gap-1">
                        <Lock className="h-3 w-3" />
                        Bloqueado
                      </Badge>
                    ) : (
                      <Badge variant={user.is_active ? 'default' : 'secondary'}>
                        {user.is_active ? 'Activo' : 'Inactivo'}
                      </Badge>
                    )}
                  </div>
                  <div className="space-y-1 text-muted-foreground">
                    <div className="flex items-center gap-2">
                      <Mail className="h-4 w-4" />
                      <span>{user.email}</span>
                      {user.email_verified && (
                        <Badge variant="outline" className="text-xs">Verificado</Badge>
                      )}
                    </div>
                    {user.phone && (
                      <div className="flex items-center gap-2">
                        <Phone className="h-4 w-4" />
                        <span>{user.phone}</span>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Details */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <UserIcon className="h-5 w-5" />
                Información del Usuario
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <p className="text-sm text-muted-foreground">Rol</p>
                  <div className="flex items-center gap-2 mt-1">
                    <Shield className="h-4 w-4" />
                    <span className="font-medium">{user.role?.display_name || user.role?.name}</span>
                  </div>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Sucursal</p>
                  <div className="flex items-center gap-2 mt-1">
                    <Building2 className="h-4 w-4" />
                    <span>{user.branch?.name || 'Todas las sucursales'}</span>
                  </div>
                </div>
              </div>

              <Separator />

              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <p className="text-sm text-muted-foreground">Último acceso</p>
                  <p className="mt-1">
                    {user.last_login_at ? formatDateTime(user.last_login_at) : 'Nunca'}
                  </p>
                </div>
                {user.last_login_ip && (
                  <div>
                    <p className="text-sm text-muted-foreground">Última IP</p>
                    <p className="mt-1 font-mono">{user.last_login_ip}</p>
                  </div>
                )}
              </div>

              {user.two_factor_enabled && (
                <>
                  <Separator />
                  <div className="flex items-center gap-2">
                    <Shield className="h-4 w-4 text-green-600" />
                    <span className="text-sm">Autenticación de dos factores habilitada</span>
                  </div>
                </>
              )}
            </CardContent>
          </Card>

          {/* Locked Warning */}
          {isLocked && (
            <Card className="border-destructive">
              <CardHeader className="pb-2">
                <CardTitle className="text-base text-destructive flex items-center gap-2">
                  <Lock className="h-4 w-4" />
                  Usuario Bloqueado
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <p className="text-sm">
                  Este usuario está bloqueado por intentos fallidos de inicio de sesión.
                </p>
                <div>
                  <p className="text-sm text-muted-foreground">Bloqueado hasta:</p>
                  <p className="font-medium">{formatDateTime(user.locked_until!)}</p>
                </div>
                <Button variant="outline" onClick={handleUnlock}>
                  <Unlock className="mr-2 h-4 w-4" />
                  Desbloquear Usuario
                </Button>
              </CardContent>
            </Card>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-4">
          {/* Actions */}
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm">Acciones</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <Button
                variant="outline"
                className="w-full justify-start"
                onClick={() => setResetDialogOpen(true)}
              >
                <Key className="mr-2 h-4 w-4" />
                Cambiar Contraseña
              </Button>
              {isLocked && (
                <Button
                  variant="outline"
                  className="w-full justify-start"
                  onClick={handleUnlock}
                >
                  <Unlock className="mr-2 h-4 w-4" />
                  Desbloquear
                </Button>
              )}
              <Button
                variant={user.is_active ? 'destructive' : 'default'}
                className="w-full justify-start"
                onClick={handleToggleActive}
              >
                <Power className="mr-2 h-4 w-4" />
                {user.is_active ? 'Desactivar Usuario' : 'Activar Usuario'}
              </Button>
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
                <span>{formatDateTime(user.created_at)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Actualizado:</span>
                <span>{formatDateTime(user.updated_at)}</span>
              </div>
              {user.password_changed_at && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Contraseña:</span>
                  <span>{formatDateTime(user.password_changed_at)}</span>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>

      <ResetPasswordDialog
        open={resetDialogOpen}
        onOpenChange={setResetDialogOpen}
        user={user}
        onConfirm={handleResetConfirm}
        isLoading={resetPasswordMutation.isPending}
      />
    </div>
  )
}
