import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import {
  Loader2,
  Check,
  Truck,
  PackageCheck,
  X,
  ArrowRightLeft,
  Building2,
  Package,
  Calendar,
  AlertTriangle,
} from 'lucide-react'

import { TRANSFER_STATUSES } from '@/types'
import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import { ROUTES, itemRoute } from '@/routes/routes'
import {
  useTransfer,
  useApproveTransfer,
  useShipTransfer,
  useReceiveTransfer,
  useCancelTransfer,
} from '@/hooks/use-transfers'
import { formatCurrency, formatDateTime } from '@/lib/format'
import { CancelTransferDialog } from './cancel-transfer-dialog'
import { ShipTransferDialog } from './ship-transfer-dialog'

function getStatusBadge(status: string) {
  const statusConfig = TRANSFER_STATUSES.find((s) => s.value === status)
  if (!statusConfig) return <Badge variant="secondary">{status}</Badge>

  const variants: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
    pending: 'secondary',
    approved: 'outline',
    in_transit: 'default',
    received: 'default',
    cancelled: 'destructive',
  }

  return <Badge variant={variants[status] || 'secondary'}>{statusConfig.label}</Badge>
}

export default function TransferDetailPage() {
  const { id } = useParams()
  const transferId = id ? parseInt(id, 10) : 0
  const navigate = useNavigate()

  const [approveDialogOpen, setApproveDialogOpen] = useState(false)
  const [shipDialogOpen, setShipDialogOpen] = useState(false)
  const [receiveDialogOpen, setReceiveDialogOpen] = useState(false)
  const [cancelDialogOpen, setCancelDialogOpen] = useState(false)

  const { data: transfer, isLoading, error } = useTransfer(transferId)
  const approveMutation = useApproveTransfer()
  const shipMutation = useShipTransfer()
  const receiveMutation = useReceiveTransfer()
  const cancelMutation = useCancelTransfer()

  const handleApprove = () => {
    approveMutation.mutate(
      { id: transferId },
      {
        onSuccess: () => setApproveDialogOpen(false),
      }
    )
  }

  const handleShip = (trackingNumber?: string, notes?: string) => {
    shipMutation.mutate(
      {
        id: transferId,
        input: {
          tracking_number: trackingNumber || undefined,
          notes: notes || undefined,
        },
      },
      {
        onSuccess: () => setShipDialogOpen(false),
      }
    )
  }

  const handleReceive = () => {
    receiveMutation.mutate(
      { id: transferId },
      {
        onSuccess: () => setReceiveDialogOpen(false),
      }
    )
  }

  const handleCancel = (reason: string) => {
    cancelMutation.mutate(
      { id: transferId, input: { reason } },
      {
        onSuccess: () => setCancelDialogOpen(false),
      }
    )
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (error || !transfer) {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-4">
        <AlertTriangle className="h-12 w-12 text-destructive" />
        <p className="text-muted-foreground">Error al cargar la transferencia</p>
        <Button asChild variant="outline">
          <Link to={ROUTES.TRANSFERS}>Volver a transferencias</Link>
        </Button>
      </div>
    )
  }

  const { status } = transfer
  const canApprove = status === 'pending'
  const canShip = status === 'approved'
  const canReceive = status === 'in_transit'
  const canCancel = status === 'pending' || status === 'approved'

  return (
    <div>
      <PageHeader
        title={transfer.transfer_number}
        description="Detalles de la transferencia"
        backUrl={ROUTES.TRANSFERS}
        actions={
          <div className="flex gap-2">
            {canApprove && (
              <Button onClick={() => setApproveDialogOpen(true)}>
                <Check className="mr-2 h-4 w-4" />
                Aprobar
              </Button>
            )}
            {canShip && (
              <Button onClick={() => setShipDialogOpen(true)}>
                <Truck className="mr-2 h-4 w-4" />
                Enviar
              </Button>
            )}
            {canReceive && (
              <Button onClick={() => setReceiveDialogOpen(true)}>
                <PackageCheck className="mr-2 h-4 w-4" />
                Recibir
              </Button>
            )}
            {canCancel && (
              <Button variant="destructive" onClick={() => setCancelDialogOpen(true)}>
                <X className="mr-2 h-4 w-4" />
                Cancelar
              </Button>
            )}
          </div>
        }
      />

      <div className="grid gap-6 md:grid-cols-3">
        {/* Main Info */}
        <div className="md:col-span-2 space-y-6">
          {/* Transfer Info */}
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle className="flex items-center gap-2">
                  <ArrowRightLeft className="h-5 w-5" />
                  Información de la Transferencia
                </CardTitle>
                {getStatusBadge(transfer.status)}
              </div>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <p className="text-sm text-muted-foreground">Sucursal Origen</p>
                  <p className="font-medium">{transfer.from_branch?.name || '-'}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Sucursal Destino</p>
                  <p className="font-medium">{transfer.to_branch?.name || '-'}</p>
                </div>
              </div>

              {transfer.reason && (
                <>
                  <Separator />
                  <div>
                    <p className="text-sm text-muted-foreground">Motivo</p>
                    <p className="mt-1">{transfer.reason}</p>
                  </div>
                </>
              )}

              {transfer.notes && (
                <div>
                  <p className="text-sm text-muted-foreground">Notas</p>
                  <p className="mt-1 whitespace-pre-wrap">{transfer.notes}</p>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Item Info */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Package className="h-5 w-5" />
                Artículo
              </CardTitle>
            </CardHeader>
            <CardContent>
              {transfer.item ? (
                <div className="flex items-center gap-4">
                  <Package className="h-12 w-12 text-muted-foreground" />
                  <div className="flex-1">
                    <Link
                      to={itemRoute(transfer.item.id)}
                      className="font-medium text-primary hover:underline"
                    >
                      {transfer.item.name}
                    </Link>
                    <p className="text-sm text-muted-foreground font-mono">
                      SKU: {transfer.item.sku}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      Categoría: {transfer.item.category?.name || '-'}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="text-sm text-muted-foreground">Valor</p>
                    <p className="font-medium">{formatCurrency(transfer.item.appraisal_value)}</p>
                  </div>
                </div>
              ) : (
                <p className="text-muted-foreground">Artículo no disponible</p>
              )}
            </CardContent>
          </Card>

          {/* Timeline / Status History */}
          <Card>
            <CardHeader>
              <CardTitle className="text-base">Historial</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {/* Created */}
                <div className="flex gap-4">
                  <div className="flex flex-col items-center">
                    <div className="h-3 w-3 rounded-full bg-primary" />
                    <div className="w-0.5 flex-1 bg-border" />
                  </div>
                  <div className="pb-4">
                    <p className="font-medium">Creada</p>
                    <p className="text-sm text-muted-foreground">
                      {formatDateTime(transfer.created_at)}
                    </p>
                  </div>
                </div>

                {/* Approved */}
                {transfer.approved_at && (
                  <div className="flex gap-4">
                    <div className="flex flex-col items-center">
                      <div className="h-3 w-3 rounded-full bg-primary" />
                      <div className="w-0.5 flex-1 bg-border" />
                    </div>
                    <div className="pb-4">
                      <p className="font-medium">Aprobada</p>
                      <p className="text-sm text-muted-foreground">
                        {formatDateTime(transfer.approved_at)} por{' '}
                        {transfer.approved_by_user?.first_name} {transfer.approved_by_user?.last_name}
                      </p>
                    </div>
                  </div>
                )}

                {/* Shipped */}
                {transfer.shipped_at && (
                  <div className="flex gap-4">
                    <div className="flex flex-col items-center">
                      <div className="h-3 w-3 rounded-full bg-primary" />
                      <div className="w-0.5 flex-1 bg-border" />
                    </div>
                    <div className="pb-4">
                      <p className="font-medium">Enviada</p>
                      <p className="text-sm text-muted-foreground">
                        {formatDateTime(transfer.shipped_at)} por{' '}
                        {transfer.shipped_by_user?.first_name} {transfer.shipped_by_user?.last_name}
                      </p>
                      {transfer.tracking_number && (
                        <p className="text-sm text-muted-foreground">
                          Seguimiento: {transfer.tracking_number}
                        </p>
                      )}
                    </div>
                  </div>
                )}

                {/* Received */}
                {transfer.received_at && (
                  <div className="flex gap-4">
                    <div className="flex flex-col items-center">
                      <div className="h-3 w-3 rounded-full bg-green-600" />
                    </div>
                    <div>
                      <p className="font-medium text-green-600">Recibida</p>
                      <p className="text-sm text-muted-foreground">
                        {formatDateTime(transfer.received_at)} por{' '}
                        {transfer.received_by_user?.first_name} {transfer.received_by_user?.last_name}
                      </p>
                    </div>
                  </div>
                )}

                {/* Cancelled */}
                {transfer.cancelled_at && (
                  <div className="flex gap-4">
                    <div className="flex flex-col items-center">
                      <div className="h-3 w-3 rounded-full bg-destructive" />
                    </div>
                    <div>
                      <p className="font-medium text-destructive">Cancelada</p>
                      <p className="text-sm text-muted-foreground">
                        {formatDateTime(transfer.cancelled_at)}
                      </p>
                      {transfer.cancellation_reason && (
                        <p className="text-sm mt-1">{transfer.cancellation_reason}</p>
                      )}
                    </div>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-4">
          {/* Branches */}
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm flex items-center gap-2">
                <Building2 className="h-4 w-4" />
                Sucursales
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div>
                <p className="text-xs text-muted-foreground uppercase">Origen</p>
                <p className="font-medium">{transfer.from_branch?.name}</p>
                {transfer.from_branch?.address && (
                  <p className="text-sm text-muted-foreground">{transfer.from_branch.address}</p>
                )}
              </div>
              <Separator />
              <div>
                <p className="text-xs text-muted-foreground uppercase">Destino</p>
                <p className="font-medium">{transfer.to_branch?.name}</p>
                {transfer.to_branch?.address && (
                  <p className="text-sm text-muted-foreground">{transfer.to_branch.address}</p>
                )}
              </div>
            </CardContent>
          </Card>

          {/* Audit Info */}
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm flex items-center gap-2">
                <Calendar className="h-4 w-4" />
                Información de Auditoría
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Creado:</span>
                <span>{formatDateTime(transfer.created_at)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Actualizado:</span>
                <span>{formatDateTime(transfer.updated_at)}</span>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      <ConfirmDialog
        open={approveDialogOpen}
        onOpenChange={setApproveDialogOpen}
        title="Aprobar Transferencia"
        description={`¿Está seguro de aprobar la transferencia ${transfer.transfer_number}?`}
        confirmText="Aprobar"
        onConfirm={handleApprove}
        isLoading={approveMutation.isPending}
      />

      <ShipTransferDialog
        transfer={transfer}
        open={shipDialogOpen}
        onOpenChange={setShipDialogOpen}
        onConfirm={handleShip}
        isLoading={shipMutation.isPending}
      />

      <ConfirmDialog
        open={receiveDialogOpen}
        onOpenChange={setReceiveDialogOpen}
        title="Recibir Transferencia"
        description={`¿Confirma la recepción de la transferencia ${transfer.transfer_number}?`}
        confirmText="Confirmar Recepción"
        onConfirm={handleReceive}
        isLoading={receiveMutation.isPending}
      />

      <CancelTransferDialog
        transfer={transfer}
        open={cancelDialogOpen}
        onOpenChange={setCancelDialogOpen}
        onConfirm={handleCancel}
        isLoading={cancelMutation.isPending}
      />
    </div>
  )
}
