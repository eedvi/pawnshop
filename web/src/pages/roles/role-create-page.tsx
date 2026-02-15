import { useNavigate } from 'react-router-dom'

import { PageHeader } from '@/components/layout/page-header'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ROUTES, roleRoute } from '@/routes/routes'
import { useCreateRole } from '@/hooks/use-roles'
import { RoleForm } from './role-form'
import { RoleFormValues } from './schemas'

export default function RoleCreatePage() {
  const navigate = useNavigate()
  const createMutation = useCreateRole()

  const handleSubmit = (values: RoleFormValues) => {
    createMutation.mutate(
      {
        name: values.name,
        display_name: values.display_name,
        description: values.description || undefined,
        permissions: values.permissions,
      },
      {
        onSuccess: (response) => {
          if (response.data) {
            navigate(roleRoute(response.data.id))
          } else {
            navigate(ROUTES.ROLES)
          }
        },
      }
    )
  }

  const handleCancel = () => {
    navigate(ROUTES.ROLES)
  }

  return (
    <div>
      <PageHeader
        title="Nuevo Rol"
        description="Crear un nuevo rol con permisos personalizados"
        backUrl={ROUTES.ROLES}
      />

      <Card>
        <CardHeader>
          <CardTitle>Información del Rol</CardTitle>
        </CardHeader>
        <CardContent>
          <RoleForm
            onSubmit={handleSubmit}
            onCancel={handleCancel}
            isLoading={createMutation.isPending}
          />
        </CardContent>
      </Card>
    </div>
  )
}
