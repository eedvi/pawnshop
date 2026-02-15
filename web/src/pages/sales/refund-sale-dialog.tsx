import { useState, useEffect } from 'react'
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
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { AlertTriangle } from 'lucide-react'
import { Sale } from '@/types'
import { formatCurrency, formatDate } from '@/lib/format'

interface RefundSaleDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  sale: Sale | null
  onConfirm: (amount: number | undefined, reason: string) => void
  isLoading?: boolean
}

export function RefundSaleDialog({
  open,
  onOpenChange,
  sale,
  onConfirm,
  isLoading = false,
}: RefundSaleDialogProps) {
  const [reason, setReason] = useState('')
  const [fullRefund, setFullRefund] = useState(true)
  const [amount, setAmount] = useState<number>(0)

  useEffect(() => {
    if (sale) {
      setAmount(sale.final_price)
    }
  }, [sale])

  const handleConfirm = () => {
    if (reason.trim()) {
      onConfirm(fullRefund ? undefined : amount, reason.trim())
    }
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      setReason('')
      setFullRefund(true)
      setAmount(0)
    }
    onOpenChange(newOpen)
  }

  const isValidAmount = fullRefund || (amount > 0 && amount <= (sale?.final_price || 0))

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-destructive">
            <AlertTriangle className="h-5 w-5" />
            Reembolsar Venta
          </DialogTitle>
          <DialogDescription>
            Esta acción reembolsará la venta y el artículo volverá a estar disponible.
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

        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <Label htmlFor="full-refund">Reembolso completo</Label>
            <Switch
              id="full-refund"
              checked={fullRefund}
              onCheckedChange={setFullRefund}
            />
          </div>

          {!fullRefund && (
            <div className="space-y-2">
              <Label htmlFor="amount">Monto a reembolsar *</Label>
              <Input
                id="amount"
                type="number"
                value={amount}
                onChange={(e) => setAmount(Number(e.target.value))}
                placeholder="0.00"
                step="0.01"
                min="0.01"
                max={sale?.final_price}
              />
              {sale && amount > sale.final_price && (
                <p className="text-sm text-destructive">
                  El monto no puede ser mayor al precio de venta
                </p>
              )}
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="reason">Razón del reembolso *</Label>
            <Textarea
              id="reason"
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              placeholder="Ingrese la razón del reembolso..."
              rows={3}
            />
          </div>
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => handleOpenChange(false)}
            disabled={isLoading}
          >
            Cancelar
          </Button>
          <Button
            variant="destructive"
            onClick={handleConfirm}
            disabled={isLoading || !reason.trim() || !isValidAmount}
          >
            {isLoading ? 'Procesando...' : 'Reembolsar'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
