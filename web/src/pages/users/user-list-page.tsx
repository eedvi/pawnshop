import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { DataTable } from '@/components/data-table'
import { ROUTES } from '@/routes/routes'
import {
  useUsers,
  useResetPassword,
  useToggleUserActive,
  useUnlockUser,
} from '@/hooks/use-users'
import { usePagination, useDebounce } from '@/hooks'
import { useBranchStore } from '@/stores/branch-store'
import { User } from '@/types'
import { getUserColumns } from './columns'
import { ResetPasswordDialog } from './reset-password-dialog'

export default function UserListPage() {
  const { pageIndex, pageSize, onPaginationChange } = usePagination()
  const { selectedBranchId } = useBranchStore()
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search, 300)

  const [resetDialogOpen, setResetDialogOpen] = useState(false)
  const [selectedUser, setSelectedUser] = useState<User | null>(null)

  const { data, isLoading } = useUsers({
    page: pageIndex + 1,
    per_page: pageSize,
    branch_id: selectedBranchId ?? undefined,
    search: debouncedSearch || undefined,
  })

  const resetPasswordMutation = useResetPassword()
  const toggleActiveMutation = useToggleUserActive()
  const unlockMutation = useUnlockUser()

  const handleResetPassword = (user: User) => {
    setSelectedUser(user)
    setResetDialogOpen(true)
  }

  const handleResetConfirm = (password: string) => {
    if (selectedUser) {
      resetPasswordMutation.mutate(
        { id: selectedUser.id, password },
        {
          onSuccess: () => {
            setResetDialogOpen(false)
            setSelectedUser(null)
          },
        }
      )
    }
  }

  const handleToggleActive = (user: User) => {
    if (confirm(`¿Está seguro de ${user.is_active ? 'desactivar' : 'activar'} este usuario?`)) {
      toggleActiveMutation.mutate(user.id)
    }
  }

  const handleUnlock = (user: User) => {
    if (confirm('¿Está seguro de desbloquear este usuario?')) {
      unlockMutation.mutate(user.id)
    }
  }

  const columns = useMemo(
    () =>
      getUserColumns({
        onResetPassword: handleResetPassword,
        onToggleActive: handleToggleActive,
        onUnlock: handleUnlock,
      }),
    []
  )

  const pageCount = data?.meta?.pagination?.total_pages ?? 1

  return (
    <div>
      <PageHeader
        title="Usuarios"
        description="Gestión de usuarios del sistema"
        actions={
          <Button asChild>
            <Link to={ROUTES.USER_CREATE}>
              <Plus className="mr-2 h-4 w-4" />
              Nuevo Usuario
            </Link>
          </Button>
        }
      />

      <DataTable
        columns={columns}
        data={data?.data ?? []}
        pageCount={pageCount}
        pageIndex={pageIndex}
        pageSize={pageSize}
        onPaginationChange={onPaginationChange}
        isLoading={isLoading}
        searchPlaceholder="Buscar por nombre o email..."
        searchValue={search}
        onSearchChange={setSearch}
        emptyMessage="No hay usuarios registrados"
      />

      <ResetPasswordDialog
        open={resetDialogOpen}
        onOpenChange={setResetDialogOpen}
        user={selectedUser}
        onConfirm={handleResetConfirm}
        isLoading={resetPasswordMutation.isPending}
      />
    </div>
  )
}
