import { useMemo } from 'react'
import { Link } from 'react-router-dom'
import { Plus } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { DataTable } from '@/components/data-table'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import { ROUTES } from '@/routes/routes'
import { useBranches, useActivateBranch, useDeactivateBranch, useDeleteBranch } from '@/hooks/use-branches'
import { usePagination, useConfirm } from '@/hooks'
import { Branch } from '@/types'
import { getBranchColumns } from './columns'

export default function BranchListPage() {
  const { pageIndex, pageSize, onPaginationChange } = usePagination()
  const confirmDelete = useConfirm()
  const confirmDeactivate = useConfirm()

  const { data, isLoading } = useBranches({
    page: pageIndex + 1,
    per_page: pageSize,
  })

  const activateMutation = useActivateBranch()
  const deactivateMutation = useDeactivateBranch()
  const deleteMutation = useDeleteBranch()

  const handleActivate = (branch: Branch) => {
    activateMutation.mutate(branch.id)
  }

  const handleDeactivate = async (branch: Branch) => {
    const confirmed = await confirmDeactivate.confirm({
      title: 'Desactivar Sucursal',
      description: `¿Estás seguro de desactivar "${branch.name}"? Las sucursales inactivas no pueden realizar operaciones.`,
      confirmLabel: 'Desactivar',
      variant: 'destructive',
    })

    if (confirmed) {
      deactivateMutation.mutate(branch.id)
    }
  }

  const handleDelete = async (branch: Branch) => {
    const confirmed = await confirmDelete.confirm({
      title: 'Eliminar Sucursal',
      description: `¿Estás seguro de eliminar "${branch.name}"? Esta acción no se puede deshacer.`,
      confirmLabel: 'Eliminar',
      variant: 'destructive',
    })

    if (confirmed) {
      deleteMutation.mutate(branch.id)
    }
  }

  const columns = useMemo(
    () =>
      getBranchColumns({
        onActivate: handleActivate,
        onDeactivate: handleDeactivate,
        onDelete: handleDelete,
      }),
    []
  )

  const pageCount = data?.meta?.pagination?.total_pages ?? 1

  return (
    <div>
      <PageHeader
        title="Sucursales"
        description="Gestión de sucursales del sistema"
        actions={
          <Button asChild>
            <Link to={ROUTES.BRANCH_CREATE}>
              <Plus className="mr-2 h-4 w-4" />
              Nueva Sucursal
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
        showToolbar={false}
        emptyMessage="No hay sucursales registradas"
      />

      <ConfirmDialog {...confirmDelete} />
      <ConfirmDialog {...confirmDeactivate} />
    </div>
  )
}
