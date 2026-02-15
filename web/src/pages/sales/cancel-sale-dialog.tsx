import { useState } from 'react'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { XCircle } from 'lucide-react'
import { Sale } from '@/types'
import { formatCurrency, formatDate } from '@/lib/format'

interface CancelSaleDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  sale: Sale | null
  onConfirm: (reason: string) => void
  isLoading?: boolean
}

export function CancelSaleDialog({
  open,
  onOpenChange,
  sale,
  onConfirm,
  isLoading = false,
}: CancelSaleDialogProps) {
  const [reason, setReason] = useState('')

  const handleConfirm = () => {
    onConfirm(reason.trim())
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      setReason('')
    }
    onOpenChange(newOpen)
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-destructive">
            <XCircle className="h-5 w-5" />
            Cancelar Venta
          </DialogTitle>
          <DialogDescription>
            Esta acción cancelará la venta pendiente. El artículo volverá a estar disponible.
          </DialogDescription>
        </DialogHeader>

        {sale && (
          <div className="rounded-lg border p-4 space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Venta:</span>
              <span className="font-mono">{sale.sale_number}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Artículo:</span>
              <span>{sale.item?.name || '-'}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Monto:</span>
              <span className="font-medium">{formatCurrency(sale.final_price)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Fecha:</span>
              <span>{formatDate(sale.sale_date)}</span>
            </div>
          </div>
        )}

        <div className="space-y-2">
          <Label htmlFor="reason">Razón de la cancelación (opcional)</Label>
          <Textarea
            id="reason"
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            placeholder="Ingrese la razón de la cancelación..."
            rows={3}
          />
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => handleOpenChange(false)}
            disabled={isLoading}
          >
            Volver
          </Button>
          <Button
            variant="destructive"
            onClick={handleConfirm}
            disabled={isLoading}
          >
            {isLoading ? 'Procesando...' : 'Cancelar Venta'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
