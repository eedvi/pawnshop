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
import { Customer } from '@/types'

interface BlockCustomerDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  customer: Customer | null
  onConfirm: (reason: string) => void
  isLoading?: boolean
}

export function BlockCustomerDialog({
  open,
  onOpenChange,
  customer,
  onConfirm,
  isLoading = false,
}: BlockCustomerDialogProps) {
  const [reason, setReason] = useState('')

  const handleConfirm = () => {
    if (reason.trim()) {
      onConfirm(reason.trim())
      setReason('')
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
          <DialogTitle>Bloquear Cliente</DialogTitle>
          <DialogDescription>
            {customer && (
              <>
                Bloquear a <strong>{customer.first_name} {customer.last_name}</strong>.
                Los clientes bloqueados no pueden realizar operaciones.
              </>
            )}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-2">
          <Label htmlFor="reason">Razón del bloqueo *</Label>
          <Textarea
            id="reason"
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            placeholder="Ingrese la razón del bloqueo..."
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
            {isLoading ? 'Bloqueando...' : 'Bloquear'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
