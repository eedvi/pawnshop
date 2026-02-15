import { z } from 'zod'

export const transferFormSchema = z.object({
  item_id: z.coerce.number().min(1, 'El artículo es requerido'),
  to_branch_id: z.coerce.number().min(1, 'La sucursal destino es requerida'),
  reason: z.string().max(500, 'Máximo 500 caracteres').optional(),
  notes: z.string().max(1000, 'Máximo 1000 caracteres').optional(),
})

export type TransferFormValues = z.infer<typeof transferFormSchema>

export const defaultTransferValues: Partial<TransferFormValues> = {
  item_id: 0,
  to_branch_id: 0,
  reason: '',
  notes: '',
}

export const shipTransferSchema = z.object({
  tracking_number: z.string().max(100, 'Máximo 100 caracteres').optional(),
  notes: z.string().max(500, 'Máximo 500 caracteres').optional(),
})

export type ShipTransferFormValues = z.infer<typeof shipTransferSchema>

export const cancelTransferSchema = z.object({
  reason: z.string().min(1, 'El motivo de cancelación es requerido').max(500, 'Máximo 500 caracteres'),
})

export type CancelTransferFormValues = z.infer<typeof cancelTransferSchema>
