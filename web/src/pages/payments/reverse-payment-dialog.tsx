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
import { AlertTriangle } from 'lucide-react'
import { Payment } from '@/types'
import { formatCurrency, formatDate } from '@/lib/format'

interface ReversePaymentDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  payment: Payment | null
  onConfirm: (reason: string) => void
  isLoading?: boolean
}

export function ReversePaymentDialog({
  open,
  onOpenChange,
  payment,
  onConfirm,
  isLoading = false,
}: ReversePaymentDialogProps) {
  const [reason, setReason] = useState('')

  const handleConfirm = () => {
    if (reason.trim()) {
      onConfirm(reason.trim())
    }
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
            <AlertTriangle className="h-5 w-5" />
            Revertir Pago
          </DialogTitle>
          <DialogDescription>
            Esta acción revertirá el pago y actualizará el saldo del préstamo.
          </DialogDescription>
        </DialogHeader>

        {payment && (
          <div className="rounded-lg border p-4 space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Pago:</span>
              <span className="font-mono">{payment.payment_number}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Monto:</span>
              <span className="font-medium">{formatCurrency(payment.amount)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Fecha:</span>
              <span>{formatDate(payment.payment_date)}</span>
            </div>
            {payment.loan && (
              <div className="flex justify-between">
                <span className="text-muted-foreground">Préstamo:</span>
                <span className="font-mono">{payment.loan.loan_number}</span>
              </div>
            )}
          </div>
        )}

        <div className="space-y-2">
          <Label htmlFor="reason">Razón de la reversión *</Label>
          <Textarea
            id="reason"
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            placeholder="Ingrese la razón de la reversión..."
            rows={3}
          />
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
            disabled={isLoading || !reason.trim()}
          >
            {isLoading ? 'Procesando...' : 'Revertir'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
