import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { DataTable } from '@/components/data-table'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import { ROUTES } from '@/routes/routes'
import { useCustomers, useBlockCustomer, useUnblockCustomer, useDeleteCustomer } from '@/hooks/use-customers'
import { usePagination, useConfirm, useDebounce } from '@/hooks'
import { useBranchStore } from '@/stores/branch-store'
import { Customer } from '@/types'
import { getCustomerColumns } from './columns'
import { BlockCustomerDialog } from './block-customer-dialog'

export default function CustomerListPage() {
  const { pageIndex, pageSize, onPaginationChange } = usePagination()
  const { selectedBranchId } = useBranchStore()
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search, 300)

  const confirmDelete = useConfirm()
  const confirmUnblock = useConfirm()
  const [blockDialogOpen, setBlockDialogOpen] = useState(false)
  const [customerToBlock, setCustomerToBlock] = useState<Customer | null>(null)

  const { data, isLoading } = useCustomers({
    page: pageIndex + 1,
    per_page: pageSize,
    branch_id: selectedBranchId ?? undefined,
    search: debouncedSearch || undefined,
  })

  const blockMutation = useBlockCustomer()
  const unblockMutation = useUnblockCustomer()
  const deleteMutation = useDeleteCustomer()

  const handleBlock = (customer: Customer) => {
    setCustomerToBlock(customer)
    setBlockDialogOpen(true)
  }

  const handleBlockConfirm = (reason: string) => {
    if (customerToBlock) {
      blockMutation.mutate(
        { id: customerToBlock.id, reason },
        {
          onSuccess: () => {
            setBlockDialogOpen(false)
            setCustomerToBlock(null)
          },
        }
      )
    }
  }

  const handleUnblock = async (customer: Customer) => {
    const confirmed = await confirmUnblock.confirm({
      title: 'Desbloquear Cliente',
      description: `¿Estás seguro de desbloquear a "${customer.first_name} ${customer.last_name}"?`,
      confirmLabel: 'Desbloquear',
    })

    if (confirmed) {
      unblockMutation.mutate(customer.id)
    }
  }

  const handleDelete = async (customer: Customer) => {
    const confirmed = await confirmDelete.confirm({
      title: 'Eliminar Cliente',
      description: `¿Estás seguro de eliminar a "${customer.first_name} ${customer.last_name}"? Esta acción no se puede deshacer.`,
      confirmLabel: 'Eliminar',
      variant: 'destructive',
    })

    if (confirmed) {
      deleteMutation.mutate(customer.id)
    }
  }

  const columns = useMemo(
    () =>
      getCustomerColumns({
        onBlock: handleBlock,
        onUnblock: handleUnblock,
        onDelete: handleDelete,
      }),
    []
  )

  const pageCount = data?.meta?.pagination?.total_pages ?? 1

  return (
    <div>
      <PageHeader
        title="Clientes"
        description="Gestión de clientes del sistema"
        actions={
          <Button asChild>
            <Link to={ROUTES.CUSTOMER_CREATE}>
              <Plus className="mr-2 h-4 w-4" />
              Nuevo Cliente
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
        searchPlaceholder="Buscar por nombre o documento..."
        searchValue={search}
        onSearchChange={setSearch}
        emptyMessage="No hay clientes registrados"
      />

      <ConfirmDialog {...confirmDelete} />
      <ConfirmDialog {...confirmUnblock} />
      <BlockCustomerDialog
        open={blockDialogOpen}
        onOpenChange={setBlockDialogOpen}
        customer={customerToBlock}
        onConfirm={handleBlockConfirm}
        isLoading={blockMutation.isPending}
      />
    </div>
  )
}
