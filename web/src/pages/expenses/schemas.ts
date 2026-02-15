import { z } from 'zod'

export const expenseFormSchema = z.object({
  branch_id: z.coerce.number().min(1, 'La sucursal es requerida'),
  category_id: z.coerce.number().optional(),
  amount: z.coerce.number().positive('El monto debe ser mayor a 0'),
  description: z.string().min(1, 'La descripción es requerida').max(500, 'Máximo 500 caracteres'),
  payment_method: z.string().min(1, 'El método de pago es requerido'),
  receipt_number: z.string().max(100, 'Máximo 100 caracteres').optional(),
  vendor: z.string().max(200, 'Máximo 200 caracteres').optional(),
  expense_date: z.string().optional(),
})

export type ExpenseFormValues = z.infer<typeof expenseFormSchema>

export const defaultExpenseValues: Partial<ExpenseFormValues> = {
  branch_id: 0,
  category_id: undefined,
  amount: 0,
  description: '',
  payment_method: 'cash',
  receipt_number: '',
  vendor: '',
  expense_date: new Date().toISOString().split('T')[0],
}

export const rejectExpenseSchema = z.object({
  reason: z.string().min(1, 'El motivo de rechazo es requerido').max(500, 'Máximo 500 caracteres'),
})

export type RejectExpenseFormValues = z.infer<typeof rejectExpenseSchema>

export const expenseCategoryFormSchema = z.object({
  code: z.string().min(1, 'El código es requerido').max(20, 'Máximo 20 caracteres'),
  name: z.string().min(1, 'El nombre es requerido').max(100, 'Máximo 100 caracteres'),
  description: z.string().max(500, 'Máximo 500 caracteres').optional(),
})

export type ExpenseCategoryFormValues = z.infer<typeof expenseCategoryFormSchema>
