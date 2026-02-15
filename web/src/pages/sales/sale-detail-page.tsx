import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import {
  Loader2,
  RotateCcw,
  XCircle,
  ShoppingBag,
  Calendar,
  User,
  Package,
  AlertTriangle,
} from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { ROUTES, itemRoute, customerRoute } from '@/routes/routes'
import { useSale, useRefundSale, useCancelSale } from '@/hooks/use-sales'
import { SALE_STATUSES, SALE_TYPES, PAYMENT_METHODS } from '@/types'
import { formatCurrency, formatDate, formatDateTime } from '@/lib/format'
import { RefundSaleDialog } from './refund-sale-dialog'
import { CancelSaleDialog } from './cancel-sale-dialog'

export default function SaleDetailPage() {
  const { id } = useParams()
  const saleId = id ? parseInt(id, 10) : 0

  const [refundDialogOpen, setRefundDialogOpen] = useState(false)
  const [cancelDialogOpen, setCancelDialogOpen] = useState(false)

  const { data: sale, isLoading, error } = useSale(saleId)
  const refundMutation = useRefundSale()
  const cancelMutation = useCancelSale()

  const handleRefundConfirm = (amount: number | undefined, reason: string) => {
    refundMutation.mutate(
      { id: saleId, input: { amount, reason } },
      {
        onSuccess: () => {
          setRefundDialogOpen(false)
        },
      }
    )
  }

  const handleCancelConfirm = (reason: string) => {
    cancelMutation.mutate(
      { id: saleId, reason: reason || undefined },
      {
        onSuccess: () => {
          setCancelDialogOpen(false)
        },
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

  if (error || !sale) {
    return (
      <div className="flex flex-col items-center justify-center h-64 gap-4">
        <AlertTriangle className="h-12 w-12 text-destructive" />
        <p className="text-muted-foreground">Error al cargar la venta</p>
        <Button asChild variant="outline">
          <Link to={ROUTES.SALES}>Volver a ventas</Link>
        </Button>
      </div>
    )
  }

  const status = SALE_STATUSES.find((s) => s.value === sale.status)
  const saleType = SALE_TYPES.find((t) => t.value === sale.sale_type)
  const method = PAYMENT_METHODS.find((m) => m.value === sale.payment_method)
  const canRefund = sale.status === 'completed'
  const canCancel = sale.status === 'pending'

  const statusColorMap: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
    green: 'default',
    yellow: 'secondary',
    red: 'destructive',
    orange: 'secondary',
    gray: 'outline',
  }

  return (
    <div>
      <PageHeader
        title={`Venta ${sale.sale_number}`}
        description="Detalles de la venta"
        backUrl={ROUTES.SALES}
        actions={
          <div className="flex gap-2">
            {canCancel && (
              <Button variant="outline" onClick={() => setCancelDialogOpen(true)}>
                <XCircle className="mr-2 h-4 w-4" />
                Cancelar
              </Button>
            )}
            {canRefund && (
              <Button variant="destructive" onClick={() => setRefundDialogOpen(true)}>
                <RotateCcw className="mr-2 h-4 w-4" />
                Reembolsar
              </Button>
            )}
          </div>
        }
      />

      <div className="grid gap-6 md:grid-cols-3">
        {/* Main Info */}
        <div className="md:col-span-2 space-y-6">
          {/* Sale Summary */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <ShoppingBag className="h-5 w-5" />
                Información de la Venta
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <p className="text-sm text-muted-foreground">Número de Venta</p>
                  <p className="font-mono text-lg">{sale.sale_number}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Estado</p>
                  <Badge variant={statusColorMap[status?.color || 'gray'] || 'outline'}>
                    {status?.label || sale.status}
                  </Badge>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Tipo de Venta</p>
                  <p className="font-medium">{saleType?.label || sale.sale_type}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Fecha de Venta</p>
                  <p>{formatDate(sale.sale_date)}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Método de Pago</p>
                  <p className="font-medium">{method?.label || sale.payment_method}</p>
                </div>
                {sale.reference_number && (
                  <div>
                    <p className="text-sm text-muted-foreground">Número de Referencia</p>
                    <p className="font-mono">{sale.reference_number}</p>
                  </div>
                )}
              </div>

              {sale.notes && (
                <>
                  <Separator />
                  <div>
                    <p className="text-sm text-muted-foreground mb-1">Notas</p>
                    <p className="text-sm">{sale.notes}</p>
                  </div>
                </>
              )}
            </CardContent>
          </Card>

          {/* Pricing Breakdown */}
          <Card>
            <CardHeader>
              <CardTitle className="text-base">Desglose de Precio</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Precio de venta:</span>
                  <span className="font-medium">{formatCurrency(sale.sale_price)}</span>
                </div>
                {sale.discount_amount > 0 && (
                  <>
                    <div className="flex justify-between text-destructive">
                      <span>Descuento:</span>
                      <span>-{formatCurrency(sale.discount_amount)}</span>
                    </div>
                    {sale.discount_reason && (
                      <div className="text-sm text-muted-foreground pl-4">
                        Razón: {sale.discount_reason}
                      </div>
                    )}
                  </>
                )}
                <Separator />
                <div className="flex justify-between font-bold text-lg">
                  <span>Total:</span>
                  <span>{formatCurrency(sale.final_price)}</span>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Refund Info */}
          {(sale.status === 'refunded' || sale.status === 'partial_refund') && (
            <Card className="border-destructive">
              <CardHeader className="pb-2">
                <CardTitle className="text-base text-destructive flex items-center gap-2">
                  <AlertTriangle className="h-4 w-4" />
                  {sale.status === 'refunded' ? 'Venta Reembolsada' : 'Reembolso Parcial'}
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                {sale.refund_amount && (
                  <div>
                    <p className="text-sm text-muted-foreground">Monto reembolsado:</p>
                    <p className="font-medium">{formatCurrency(sale.refund_amount)}</p>
                  </div>
                )}
                {sale.refund_reason && (
                  <div>
                    <p className="text-sm text-muted-foreground">Razón:</p>
                    <p className="text-sm">{sale.refund_reason}</p>
                  </div>
                )}
                {sale.refunded_at && (
                  <div>
                    <p className="text-sm text-muted-foreground">Fecha de reembolso:</p>
                    <p className="text-sm">{formatDateTime(sale.refunded_at)}</p>
                  </div>
                )}
              </CardContent>
            </Card>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-4">
          {/* Item Info */}
          {sale.item && (
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm flex items-center gap-2">
                  <Package className="h-4 w-4" />
                  Artículo Vendido
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Nombre:</span>
                  <Link
                    to={itemRoute(sale.item.id)}
                    className="text-primary hover:underline"
                  >
                    {sale.item.name}
                  </Link>
                </div>
                {sale.item.brand && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Marca:</span>
                    <span>{sale.item.brand}</span>
                  </div>
                )}
                {sale.item.model && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Modelo:</span>
                    <span>{sale.item.model}</span>
                  </div>
                )}
                {sale.item.serial_number && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Serie:</span>
                    <span className="font-mono">{sale.item.serial_number}</span>
                  </div>
                )}
              </CardContent>
            </Card>
          )}

          {/* Customer Info */}
          {sale.customer && (
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm flex items-center gap-2">
                  <User className="h-4 w-4" />
                  Cliente
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Nombre:</span>
                  <Link
                    to={customerRoute(sale.customer.id)}
                    className="text-primary hover:underline"
                  >
                    {sale.customer.first_name} {sale.customer.last_name}
                  </Link>
                </div>
                {sale.customer.phone && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Teléfono:</span>
                    <span>{sale.customer.phone}</span>
                  </div>
                )}
                {sale.customer.email && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Email:</span>
                    <span>{sale.customer.email}</span>
                  </div>
                )}
              </CardContent>
            </Card>
          )}

          {!sale.customer && (
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm flex items-center gap-2">
                  <User className="h-4 w-4" />
                  Cliente
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground">
                  Venta sin cliente registrado
                </p>
              </CardContent>
            </Card>
          )}

          {/* Audit Info */}
          <Card>
            <CardHeader className="pb-2">
              <CardTitle className="text-sm flex items-center gap-2">
                <Calendar className="h-4 w-4" />
                Información de Auditoría
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2 text-sm">
              {sale.created_at && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Creado:</span>
                  <span>{formatDateTime(sale.created_at)}</span>
                </div>
              )}
              {sale.updated_at && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Actualizado:</span>
                  <span>{formatDateTime(sale.updated_at)}</span>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>

      <RefundSaleDialog
        open={refundDialogOpen}
        onOpenChange={setRefundDialogOpen}
        sale={sale}
        onConfirm={handleRefundConfirm}
        isLoading={refundMutation.isPending}
      />

      <CancelSaleDialog
        open={cancelDialogOpen}
        onOpenChange={setCancelDialogOpen}
        sale={sale}
        onConfirm={handleCancelConfirm}
        isLoading={cancelMutation.isPending}
      />
    </div>
  )
}
