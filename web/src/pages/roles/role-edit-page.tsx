import { useNavigate, useParams, Link } from 'react-router-dom'
import { Loader2, AlertTriangle } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { ROUTES, roleRoute } from '@/routes/routes'
import { useRole, useUpdateRole } from '@/hooks/use-roles'
import { RoleForm } from './role-form'
import { RoleFormValues } from './schemas'

export default function RoleEditPage() {
  const { id } = useParams()
  const roleId = id ? parseInt(id, 10) : 0
  const navigate = useNavigate()

  const { data: role, isLoading, error } = useRole(roleId)
  const updateMutation = useUpdateRole()

  const handleSubmit = (values: RoleFormValues) => {
    updateMutation.mutate(
      {
        id: roleId,
        input: {
          name: values.name,
          display_name: values.display_name,
          description: values.description || undefined,
          permissions: values.permissions,
        },
      },
      {
        onSuccess: () => {
          navigate(roleRoute(roleId))
        },
      }
    )
  }

  const handleCancel = () => {
    navigate(roleRoute(roleId))
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

  if (role.is_system) {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-4">
        <AlertTriangle className="h-12 w-12 text-warning" />
        <p className="text-muted-foreground">
          Los roles de sistema no se pueden editar
        </p>
        <Button asChild variant="outline">
          <Link to={roleRoute(roleId)}>Volver al rol</Link>
        </Button>
      </div>
    )
  }

  return (
    <div>
      <PageHeader
        title={`Editar ${role.display_name}`}
        description="Modificar información y permisos del rol"
        backUrl={roleRoute(roleId)}
      />

      <Card>
        <CardHeader>
          <CardTitle>Información del Rol</CardTitle>
        </CardHeader>
        <CardContent>
          <RoleForm
            role={role}
            onSubmit={handleSubmit}
            onCancel={handleCancel}
            isLoading={updateMutation.isPending}
          />
        </CardContent>
      </Card>
    </div>
  )
}
