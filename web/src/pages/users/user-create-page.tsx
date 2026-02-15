import { useNavigate } from 'react-router-dom'

import { PageHeader } from '@/components/layout/page-header'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ROUTES } from '@/routes/routes'
import { useCreateUser } from '@/hooks/use-users'
import { UserForm } from './user-form'
import { CreateUserFormValues } from './schemas'

export default function UserCreatePage() {
  const navigate = useNavigate()
  const createMutation = useCreateUser()

  const handleSubmit = (values: CreateUserFormValues) => {
    createMutation.mutate(
      {
        email: values.email,
        password: values.password,
        first_name: values.first_name,
        last_name: values.last_name,
        phone: values.phone || undefined,
        branch_id: values.branch_id,
        role_id: values.role_id,
      },
      {
        onSuccess: () => {
          navigate(ROUTES.USERS)
        },
      }
    )
  }

  const handleCancel = () => {
    navigate(ROUTES.USERS)
  }

  return (
    <div>
      <PageHeader
        title="Nuevo Usuario"
        description="Registrar un nuevo usuario en el sistema"
        backUrl={ROUTES.USERS}
      />

      <Card>
        <CardHeader>
          <CardTitle>Información del Usuario</CardTitle>
        </CardHeader>
        <CardContent>
          <UserForm
            onSubmit={handleSubmit}
            onCancel={handleCancel}
            isLoading={createMutation.isPending}
          />
        </CardContent>
      </Card>
    </div>
  )
}
