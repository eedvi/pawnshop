import { useState, useMemo } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { Plus, Loader2 } from 'lucide-react'

import { ItemTransfer, TransferStatus, TRANSFER_STATUSES } from '@/types'
import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { DataTable } from '@/components/data-table/data-table'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { ROUTES } from '@/routes/routes'
import { useBranchStore } from '@/stores/branch-store'
import {
  useTransfers,
  useApproveTransfer,
  useShipTransfer,
  useReceiveTransfer,
  useCancelTransfer,
} from '@/hooks/use-transfers'
import { getTransferColumns } from './columns'
import { CancelTransferDialog } from './cancel-transfer-dialog'
import { ShipTransferDialog } from './ship-transfer-dialog'

export default function TransferListPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const { selectedBranch } = useBranchStore()

  const [approveTransfer, setApproveTransfer] = useState<ItemTransfer | null>(null)
  const [shipTransfer, setShipTransfer] = useState<ItemTransfer | null>(null)
  const [receiveTransfer, setReceiveTransfer] = useState<ItemTransfer | null>(null)
  const [cancelTransfer, setCancelTransfer] = useState<ItemTransfer | null>(null)

  const page = parseInt(searchParams.get('page') || '1')
  const status = searchParams.get('status') as TransferStatus | null

  const { data: transfersResponse, isLoading } = useTransfers({
    page,
    per_page: 10,
    from_branch_id: selectedBranch?.id,
    status: status || undefined,
    order_by: 'created_at',
    order: 'desc',
  })

  const approveMutation = useApproveTransfer()
  const shipMutation = useShipTransfer()
  const receiveMutation = useReceiveTransfer()
  const cancelMutation = useCancelTransfer()

  const transfers = transfersResponse?.data || []
  const pagination = transfersResponse?.meta?.pagination

  const columns = useMemo(
    () =>
      getTransferColumns({
        onApprove: (transfer) => setApproveTransfer(transfer),
        onShip: (transfer) => setShipTransfer(transfer),
        onReceive: (transfer) => setReceiveTransfer(transfer),
        onCancel: (transfer) => setCancelTransfer(transfer),
      }),
    []
  )

  const handleApproveConfirm = () => {
    if (approveTransfer) {
      approveMutation.mutate(
        { id: approveTransfer.id },
        {
          onSuccess: () => setApproveTransfer(null),
        }
      )
    }
  }

  const handleShipConfirm = (trackingNumber?: string, notes?: string) => {
    if (shipTransfer) {
      shipMutation.mutate(
        {
          id: shipTransfer.id,
          input: {
            tracking_number: trackingNumber || undefined,
            notes: notes || undefined,
          },
        },
        {
          onSuccess: () => setShipTransfer(null),
        }
      )
    }
  }

  const handleReceiveConfirm = () => {
    if (receiveTransfer) {
      receiveMutation.mutate(
        { id: receiveTransfer.id },
        {
          onSuccess: () => setReceiveTransfer(null),
        }
      )
    }
  }

  const handleCancelConfirm = (reason: string) => {
    if (cancelTransfer) {
      cancelMutation.mutate(
        { id: cancelTransfer.id, input: { reason } },
        {
          onSuccess: () => setCancelTransfer(null),
        }
      )
    }
  }

  const handleStatusChange = (value: string) => {
    const params = new URLSearchParams(searchParams)
    if (value === 'all') {
      params.delete('status')
    } else {
      params.set('status', value)
    }
    params.set('page', '1')
    setSearchParams(params)
  }

  const handlePageChange = (newPage: number) => {
    const params = new URLSearchParams(searchParams)
    params.set('page', newPage.toString())
    setSearchParams(params)
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
        title="Transferencias"
        description="Transferencias de artículos entre sucursales"
        actions={
          <Button asChild>
            <Link to={ROUTES.TRANSFER_CREATE}>
              <Plus className="mr-2 h-4 w-4" />
              Nueva Transferencia
            </Link>
          </Button>
        }
      />

      <div className="mb-4 flex gap-4">
        <Select value={status || 'all'} onValueChange={handleStatusChange}>
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Filtrar por estado" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Todos los estados</SelectItem>
            {TRANSFER_STATUSES.map((s) => (
              <SelectItem key={s.value} value={s.value}>
                {s.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <DataTable
        columns={columns}
        data={transfers}
        searchPlaceholder="Buscar transferencias..."
        searchColumn="transfer_number"
        pagination={
          pagination
            ? {
                page: pagination.current_page,
                pageSize: pagination.per_page,
                totalPages: pagination.total_pages,
                totalItems: pagination.total,
                onPageChange: handlePageChange,
              }
            : undefined
        }
      />

      <ConfirmDialog
        open={!!approveTransfer}
        onOpenChange={(open) => !open && setApproveTransfer(null)}
        title="Aprobar Transferencia"
        description={
          approveTransfer
            ? `¿Está seguro de aprobar la transferencia ${approveTransfer.transfer_number}?`
            : ''
        }
        confirmText="Aprobar"
        onConfirm={handleApproveConfirm}
        isLoading={approveMutation.isPending}
      />

      <ShipTransferDialog
        transfer={shipTransfer}
        open={!!shipTransfer}
        onOpenChange={(open) => !open && setShipTransfer(null)}
        onConfirm={handleShipConfirm}
        isLoading={shipMutation.isPending}
      />

      <ConfirmDialog
        open={!!receiveTransfer}
        onOpenChange={(open) => !open && setReceiveTransfer(null)}
        title="Recibir Transferencia"
        description={
          receiveTransfer
            ? `¿Confirma la recepción de la transferencia ${receiveTransfer.transfer_number}?`
            : ''
        }
        confirmText="Confirmar Recepción"
        onConfirm={handleReceiveConfirm}
        isLoading={receiveMutation.isPending}
      />

      <CancelTransferDialog
        transfer={cancelTransfer}
        open={!!cancelTransfer}
        onOpenChange={(open) => !open && setCancelTransfer(null)}
        onConfirm={handleCancelConfirm}
        isLoading={cancelMutation.isPending}
      />
    </div>
  )
}
