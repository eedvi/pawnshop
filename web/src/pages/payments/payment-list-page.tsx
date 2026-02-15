import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { DataTable } from '@/components/data-table'
import { ROUTES } from '@/routes/routes'
import { usePayments, useReversePayment } from '@/hooks/use-payments'
import { usePagination, useDebounce } from '@/hooks'
import { useBranchStore } from '@/stores/branch-store'
import { Payment } from '@/types'
import { getPaymentColumns } from './columns'
import { ReversePaymentDialog } from './reverse-payment-dialog'

export default function PaymentListPage() {
  const { pageIndex, pageSize, onPaginationChange } = usePagination()
  const { selectedBranchId } = useBranchStore()
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search, 300)

  const [reverseDialogOpen, setReverseDialogOpen] = useState(false)
  const [selectedPayment, setSelectedPayment] = useState<Payment | null>(null)

  const { data, isLoading } = usePayments({
    page: pageIndex + 1,
    per_page: pageSize,
    branch_id: selectedBranchId ?? undefined,
    search: debouncedSearch || undefined,
  })

  const reverseMutation = useReversePayment()

  const handleReverse = (payment: Payment) => {
    setSelectedPayment(payment)
    setReverseDialogOpen(true)
  }

  const handleReverseConfirm = (reason: string) => {
    if (selectedPayment) {
      reverseMutation.mutate(
        { id: selectedPayment.id, input: { reason } },
        {
          onSuccess: () => {
            setReverseDialogOpen(false)
            setSelectedPayment(null)
          },
        }
      )
    }
  }

  const columns = useMemo(
    () =>
      getPaymentColumns({
        onReverse: handleReverse,
      }),
    []
  )

  const pageCount = data?.meta?.pagination?.total_pages ?? 1

  return (
    <div>
      <PageHeader
        title="Pagos"
        description="Registro de pagos recibidos"
        actions={
          <Button asChild>
            <Link to={ROUTES.PAYMENT_CREATE}>
              <Plus className="mr-2 h-4 w-4" />
              Nuevo Pago
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
        searchPlaceholder="Buscar por número de pago o cliente..."
        searchValue={search}
        onSearchChange={setSearch}
        emptyMessage="No hay pagos registrados"
      />

      <ReversePaymentDialog
        open={reverseDialogOpen}
        onOpenChange={setReverseDialogOpen}
        payment={selectedPayment}
        onConfirm={handleReverseConfirm}
        isLoading={reverseMutation.isPending}
      />
    </div>
  )
}
