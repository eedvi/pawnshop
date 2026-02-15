import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { DataTable } from '@/components/data-table'
import { ROUTES } from '@/routes/routes'
import { useSales, useRefundSale, useCancelSale } from '@/hooks/use-sales'
import { usePagination, useDebounce } from '@/hooks'
import { useBranchStore } from '@/stores/branch-store'
import { Sale } from '@/types'
import { getSaleColumns } from './columns'
import { RefundSaleDialog } from './refund-sale-dialog'
import { CancelSaleDialog } from './cancel-sale-dialog'

export default function SaleListPage() {
  const { pageIndex, pageSize, onPaginationChange } = usePagination()
  const { selectedBranchId } = useBranchStore()
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search, 300)

  const [refundDialogOpen, setRefundDialogOpen] = useState(false)
  const [cancelDialogOpen, setCancelDialogOpen] = useState(false)
  const [selectedSale, setSelectedSale] = useState<Sale | null>(null)

  const { data, isLoading } = useSales({
    page: pageIndex + 1,
    per_page: pageSize,
    branch_id: selectedBranchId ?? undefined,
    search: debouncedSearch || undefined,
  })

  const refundMutation = useRefundSale()
  const cancelMutation = useCancelSale()

  const handleRefund = (sale: Sale) => {
    setSelectedSale(sale)
    setRefundDialogOpen(true)
  }

  const handleCancel = (sale: Sale) => {
    setSelectedSale(sale)
    setCancelDialogOpen(true)
  }

  const handleRefundConfirm = (amount: number | undefined, reason: string) => {
    if (selectedSale) {
      refundMutation.mutate(
        { id: selectedSale.id, input: { amount, reason } },
        {
          onSuccess: () => {
            setRefundDialogOpen(false)
            setSelectedSale(null)
          },
        }
      )
    }
  }

  const handleCancelConfirm = (reason: string) => {
    if (selectedSale) {
      cancelMutation.mutate(
        { id: selectedSale.id, reason: reason || undefined },
        {
          onSuccess: () => {
            setCancelDialogOpen(false)
            setSelectedSale(null)
          },
        }
      )
    }
  }

  const columns = useMemo(
    () =>
      getSaleColumns({
        onRefund: handleRefund,
        onCancel: handleCancel,
      }),
    []
  )

  const pageCount = data?.meta?.pagination?.total_pages ?? 1

  return (
    <div>
      <PageHeader
        title="Ventas"
        description="Registro de ventas realizadas"
        actions={
          <Button asChild>
            <Link to={ROUTES.SALE_CREATE}>
              <Plus className="mr-2 h-4 w-4" />
              Nueva Venta
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
        searchPlaceholder="Buscar por número de venta o artículo..."
        searchValue={search}
        onSearchChange={setSearch}
        emptyMessage="No hay ventas registradas"
      />

      <RefundSaleDialog
        open={refundDialogOpen}
        onOpenChange={setRefundDialogOpen}
        sale={selectedSale}
        onConfirm={handleRefundConfirm}
        isLoading={refundMutation.isPending}
      />

      <CancelSaleDialog
        open={cancelDialogOpen}
        onOpenChange={setCancelDialogOpen}
        sale={selectedSale}
        onConfirm={handleCancelConfirm}
        isLoading={cancelMutation.isPending}
      />
    </div>
  )
}
