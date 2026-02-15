import { useNavigate, useParams } from 'react-router-dom'
import { Loader2, AlertTriangle } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { ROUTES, userRoute } from '@/routes/routes'
import { useUser, useUpdateUser } from '@/hooks/use-users'
import { UserForm } from './user-form'
import { UserFormValues } from './schemas'

export default function UserEditPage() {
  const { id } = useParams()
  const userId = id ? parseInt(id, 10) : 0
  const navigate = useNavigate()

  const { data: user, isLoading, error } = useUser(userId)
  const updateMutation = useUpdateUser()

  const handleSubmit = (values: UserFormValues) => {
    updateMutation.mutate(
      {
        id: userId,
        input: {
          email: values.email,
          first_name: values.first_name,
          last_name: values.last_name,
          phone: values.phone || undefined,
          branch_id: values.branch_id,
          role_id: values.role_id,
        },
      },
      {
        onSuccess: () => {
          navigate(userRoute(userId))
        },
      }
    )
  }

  const handleCancel = () => {
    navigate(userRoute(userId))
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
          <a href={ROUTES.USERS}>Volver a usuarios</a>
        </Button>
      </div>
    )
  }

  return (
    <div>
      <PageHeader
        title={`Editar ${user.first_name} ${user.last_name}`}
        description="Modificar información del usuario"
        backUrl={userRoute(userId)}
      />

      <Card>
        <CardHeader>
          <CardTitle>Información del Usuario</CardTitle>
        </CardHeader>
        <CardContent>
          <UserForm
            user={user}
            onSubmit={handleSubmit}
            onCancel={handleCancel}
            isLoading={updateMutation.isPending}
          />
        </CardContent>
      </Card>
    </div>
  )
}
