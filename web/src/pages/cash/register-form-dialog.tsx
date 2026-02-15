import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Form } from '@/components/ui/form'
import { FormInput, FormTextarea } from '@/components/form'
import { Loader2 } from 'lucide-react'
import { CashRegister } from '@/types'
import { cashRegisterFormSchema, CashRegisterFormValues } from './schemas'

interface RegisterFormDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  register?: CashRegister | null
  onConfirm: (data: CashRegisterFormValues) => void
  isLoading?: boolean
}

export function RegisterFormDialog({
  open,
  onOpenChange,
  register,
  onConfirm,
  isLoading = false,
}: RegisterFormDialogProps) {
  const isEditing = !!register

  const form = useForm<CashRegisterFormValues>({
    resolver: zodResolver(cashRegisterFormSchema),
    defaultValues: {
      name: '',
      code: '',
      description: '',
    },
  })

  useEffect(() => {
    if (register) {
      form.reset({
        name: register.name,
        code: register.code,
        description: register.description || '',
      })
    } else {
      form.reset({
        name: '',
        code: '',
        description: '',
      })
    }
  }, [register, form])

  const handleSubmit = (values: CashRegisterFormValues) => {
    onConfirm(values)
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!newOpen) {
      form.reset()
    }
    onOpenChange(newOpen)
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEditing ? 'Editar Caja' : 'Nueva Caja'}</DialogTitle>
          <DialogDescription>
            {isEditing
              ? 'Modifique los datos de la caja registradora'
              : 'Ingrese los datos para crear una nueva caja registradora'}
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormInput
              control={form.control}
              name="name"
              label="Nombre"
              placeholder="Caja Principal"
              required
            />

            <FormInput
              control={form.control}
              name="code"
              label="Código"
              placeholder="CAJA-01"
              required
              disabled={isEditing}
            />

            <FormTextarea
              control={form.control}
              name="description"
              label="Descripción"
              placeholder="Descripción de la caja..."
              rows={2}
            />

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => handleOpenChange(false)}
                disabled={isLoading}
              >
                Cancelar
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {isEditing ? 'Guardar' : 'Crear'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
