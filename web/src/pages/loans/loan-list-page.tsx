import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { DataTable } from '@/components/data-table'
import { ROUTES } from '@/routes/routes'
import { useLoans, useRenewLoan, useConfiscateLoan } from '@/hooks/use-loans'
import { usePagination, useDebounce } from '@/hooks'
import { useBranchStore } from '@/stores/branch-store'
import { Loan, RenewLoanInput } from '@/types'
import { getLoanColumns } from './columns'
import { RenewLoanDialog } from './renew-loan-dialog'
import { ConfiscateLoanDialog } from './confiscate-loan-dialog'

export default function LoanListPage() {
  const { pageIndex, pageSize, onPaginationChange } = usePagination()
  const { selectedBranchId } = useBranchStore()
  const [search, setSearch] = useState('')
  const debouncedSearch = useDebounce(search, 300)

  const [renewDialogOpen, setRenewDialogOpen] = useState(false)
  const [confiscateDialogOpen, setConfiscateDialogOpen] = useState(false)
  const [selectedLoan, setSelectedLoan] = useState<Loan | null>(null)

  const { data, isLoading } = useLoans({
    page: pageIndex + 1,
    per_page: pageSize,
    branch_id: selectedBranchId ?? undefined,
    search: debouncedSearch || undefined,
  })

  const renewMutation = useRenewLoan()
  const confiscateMutation = useConfiscateLoan()

  const handleRenew = (loan: Loan) => {
    setSelectedLoan(loan)
    setRenewDialogOpen(true)
  }

  const handleRenewConfirm = (values: RenewLoanInput) => {
    if (selectedLoan) {
      renewMutation.mutate(
        { id: selectedLoan.id, input: values },
        {
          onSuccess: () => {
            setRenewDialogOpen(false)
            setSelectedLoan(null)
          },
        }
      )
    }
  }

  const handleConfiscate = (loan: Loan) => {
    setSelectedLoan(loan)
    setConfiscateDialogOpen(true)
  }

  const handleConfiscateConfirm = (notes?: string) => {
    if (selectedLoan) {
      confiscateMutation.mutate(
        { id: selectedLoan.id, notes },
        {
          onSuccess: () => {
            setConfiscateDialogOpen(false)
            setSelectedLoan(null)
          },
        }
      )
    }
  }

  const columns = useMemo(
    () =>
      getLoanColumns({
        onRenew: handleRenew,
        onConfiscate: handleConfiscate,
      }),
    []
  )

  const pageCount = data?.meta?.pagination?.total_pages ?? 1

  return (
    <div>
      <PageHeader
        title="Préstamos"
        description="Gestión de préstamos"
        actions={
          <Button asChild>
            <Link to={ROUTES.LOAN_CREATE}>
              <Plus className="mr-2 h-4 w-4" />
              Nuevo Préstamo
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
        searchPlaceholder="Buscar por número de préstamo o cliente..."
        searchValue={search}
        onSearchChange={setSearch}
        emptyMessage="No hay préstamos registrados"
      />

      <RenewLoanDialog
        open={renewDialogOpen}
        onOpenChange={setRenewDialogOpen}
        loan={selectedLoan}
        onConfirm={handleRenewConfirm}
        isLoading={renewMutation.isPending}
      />

      <ConfiscateLoanDialog
        open={confiscateDialogOpen}
        onOpenChange={setConfiscateDialogOpen}
        loan={selectedLoan}
        onConfirm={handleConfiscateConfirm}
        isLoading={confiscateMutation.isPending}
      />
    </div>
  )
}
