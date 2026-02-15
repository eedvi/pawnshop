import { useState, useMemo } from 'react'
import { Link } from 'react-router-dom'
import { Plus, Loader2 } from 'lucide-react'

import { Role } from '@/types'
import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { DataTable } from '@/components/data-table/data-table'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import { ROUTES } from '@/routes/routes'
import { useRoles, useDeleteRole } from '@/hooks/use-roles'
import { getRoleColumns } from './columns'

export default function RoleListPage() {
  const [deleteRole, setDeleteRole] = useState<Role | null>(null)

  const { data: rolesResponse, isLoading } = useRoles()
  const deleteMutation = useDeleteRole()

  const roles = Array.isArray(rolesResponse) ? rolesResponse : rolesResponse?.data || []

  const columns = useMemo(
    () =>
      getRoleColumns({
        onDelete: (role) => setDeleteRole(role),
      }),
    []
  )

  const handleDeleteConfirm = () => {
    if (deleteRole) {
      deleteMutation.mutate(deleteRole.id, {
        onSuccess: () => setDeleteRole(null),
      })
    }
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
        title="Roles"
        description="Gestión de roles y permisos del sistema"
        actions={
          <Button asChild>
            <Link to={ROUTES.ROLE_CREATE}>
              <Plus className="mr-2 h-4 w-4" />
              Nuevo Rol
            </Link>
          </Button>
        }
      />

      <DataTable
        columns={columns}
        data={roles}
        searchPlaceholder="Buscar roles..."
        searchColumn="display_name"
      />

      <ConfirmDialog
        open={!!deleteRole}
        onOpenChange={(open) => !open && setDeleteRole(null)}
        title="Eliminar Rol"
        description={
          deleteRole
            ? `¿Está seguro de eliminar el rol "${deleteRole.display_name}"? Esta acción no se puede deshacer.`
            : ''
        }
        confirmText="Eliminar"
        variant="destructive"
        onConfirm={handleDeleteConfirm}
        isLoading={deleteMutation.isPending}
      />
    </div>
  )
}
