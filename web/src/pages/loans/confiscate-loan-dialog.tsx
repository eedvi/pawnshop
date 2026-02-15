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
import { Loan } from '@/types'
import { formatCurrency, formatDate } from '@/lib/format'

interface ConfiscateLoanDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  loan: Loan | null
  onConfirm: (notes?: string) => void
  isLoading?: boolean
}

export function ConfiscateLoanDialog({
  open,
  onOpenChange,
  loan,
  onConfirm,
  isLoading = false,
}: ConfiscateLoanDialogProps) {
  const [notes, setNotes] = useState('')

  const handleConfirm = () => {
    onConfirm(notes || undefined)
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      setNotes('')
    }
    onOpenChange(newOpen)
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-destructive">
            <AlertTriangle className="h-5 w-5" />
            Confiscar Préstamo
          </DialogTitle>
          <DialogDescription>
            Esta acción marcará el préstamo como confiscado y el artículo pasará a propiedad
            de la casa de empeño.
          </DialogDescription>
        </DialogHeader>

        {loan && (
          <div className="rounded-lg border p-4 space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Préstamo:</span>
              <span className="font-mono">{loan.loan_number}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Cliente:</span>
              <span>{loan.customer?.first_name} {loan.customer?.last_name}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Artículo:</span>
              <span>{loan.item?.name}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Saldo pendiente:</span>
              <span className="font-medium">
                {formatCurrency(loan.principal_remaining + loan.interest_remaining + loan.late_fee_amount)}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Días vencido:</span>
              <span className="text-destructive font-medium">{loan.days_overdue} días</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Fecha vencimiento:</span>
              <span>{formatDate(loan.due_date)}</span>
            </div>
          </div>
        )}

        <div className="space-y-2">
          <Label htmlFor="notes">Notas (opcional)</Label>
          <Textarea
            id="notes"
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
            placeholder="Razón de la confiscación..."
            rows={2}
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
            disabled={isLoading}
          >
            {isLoading ? 'Procesando...' : 'Confiscar'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
