import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { DataTable } from '@/components/data-table'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import { ROUTES } from '@/routes/routes'
import { useItems, useMarkItemForSale, useDeleteItem } from '@/hooks/use-items'
import { usePagination, useConfirm, useDebounce } from '@/hooks'
import { useBranchStore } from '@/stores/branch-store'
import { Item } from '@/types'
import { getItemColumns } from './columns'
import { MarkForSaleDialog } from './mark-for-sale-dialog'

export default function ItemListPage() {
  const { pageIndex, pageSize, onPaginationChange } = usePagination()
  const { selectedBranchId } = useBranchStore()
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search, 300)

  const confirmDelete = useConfirm()
  const [markForSaleOpen, setMarkForSaleOpen] = useState(false)
  const [selectedItem, setSelectedItem] = useState<Item | null>(null)

  const { data, isLoading } = useItems({
    page: pageIndex + 1,
    per_page: pageSize,
    branch_id: selectedBranchId ?? undefined,
    search: debouncedSearch || undefined,
  })

  const markForSaleMutation = useMarkItemForSale()
  const deleteMutation = useDeleteItem()

  const handleMarkForSale = (item: Item) => {
    setSelectedItem(item)
    setMarkForSaleOpen(true)
  }

  const handleMarkForSaleConfirm = (salePrice: number) => {
    if (selectedItem) {
      markForSaleMutation.mutate(
        { id: selectedItem.id, salePrice },
        {
          onSuccess: () => {
            setMarkForSaleOpen(false)
            setSelectedItem(null)
          },
        }
      )
    }
  }

  const handleDelete = async (item: Item) => {
    const confirmed = await confirmDelete.confirm({
      title: 'Eliminar Artículo',
      description: `¿Estás seguro de eliminar "${item.name}"? Esta acción no se puede deshacer.`,
      confirmLabel: 'Eliminar',
      variant: 'destructive',
    })

    if (confirmed) {
      deleteMutation.mutate(item.id)
    }
  }

  const columns = useMemo(
    () =>
      getItemColumns({
        onMarkForSale: handleMarkForSale,
        onDelete: handleDelete,
      }),
    []
  )

  const pageCount = data?.meta?.pagination?.total_pages ?? 1

  return (
    <div>
      <PageHeader
        title="Artículos"
        description="Inventario de artículos"
        actions={
          <Button asChild>
            <Link to={ROUTES.ITEM_CREATE}>
              <Plus className="mr-2 h-4 w-4" />
              Nuevo Artículo
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
        searchPlaceholder="Buscar por SKU, nombre o marca..."
        searchValue={search}
        onSearchChange={setSearch}
        emptyMessage="No hay artículos registrados"
      />

      <ConfirmDialog {...confirmDelete} />
      <MarkForSaleDialog
        open={markForSaleOpen}
        onOpenChange={setMarkForSaleOpen}
        item={selectedItem}
        onConfirm={handleMarkForSaleConfirm}
        isLoading={markForSaleMutation.isPending}
      />
    </div>
  )
}
