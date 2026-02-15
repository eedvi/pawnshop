import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

import {
  NotificationTemplate,
  NotificationChannel,
  NotificationType,
  NOTIFICATION_CHANNELS,
} from '@/types'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

const NOTIFICATION_TYPES: { value: NotificationType; label: string }[] = [
  { value: 'payment_reminder', label: 'Recordatorio de pago' },
  { value: 'overdue_notice', label: 'Aviso de mora' },
  { value: 'loan_expiry', label: 'Vencimiento de préstamo' },
  { value: 'promotional', label: 'Promocional' },
  { value: 'system', label: 'Sistema' },
]

const createTemplateSchema = z.object({
  name: z.string().min(1, 'El nombre es requerido'),
  code: z
    .string()
    .min(1, 'El código es requerido')
    .regex(/^[a-z0-9_]+$/, 'Solo minúsculas, números y guiones bajos'),
  channel: z.enum(['email', 'sms', 'whatsapp', 'internal'] as const),
  notification_type: z.enum([
    'payment_reminder',
    'overdue_notice',
    'loan_expiry',
    'promotional',
    'system',
  ] as const),
  subject: z.string().optional(),
  content: z.string().min(1, 'El contenido es requerido'),
})

const editTemplateSchema = z.object({
  name: z.string().min(1, 'El nombre es requerido'),
  subject: z.string().optional(),
  content: z.string().min(1, 'El contenido es requerido'),
})

type CreateFormValues = z.infer<typeof createTemplateSchema>
type EditFormValues = z.infer<typeof editTemplateSchema>

interface TemplateFormDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  template?: NotificationTemplate | null
  onSubmit: (data: CreateFormValues | EditFormValues) => void
  isLoading?: boolean
}

export function TemplateFormDialog({
  open,
  onOpenChange,
  template,
  onSubmit,
  isLoading,
}: TemplateFormDialogProps) {
  const isEditing = !!template

  const createForm = useForm<CreateFormValues>({
    resolver: zodResolver(createTemplateSchema),
    defaultValues: {
      name: '',
      code: '',
      channel: 'email',
      notification_type: 'system',
      subject: '',
      content: '',
    },
  })

  const editForm = useForm<EditFormValues>({
    resolver: zodResolver(editTemplateSchema),
    defaultValues: {
      name: '',
      subject: '',
      content: '',
    },
  })

  const form = isEditing ? editForm : createForm

  useEffect(() => {
    if (template) {
      editForm.reset({
        name: template.name,
        subject: template.subject || '',
        content: template.content,
      })
    } else {
      createForm.reset({
        name: '',
        code: '',
        channel: 'email',
        notification_type: 'system',
        subject: '',
        content: '',
      })
    }
  }, [template, createForm, editForm])

  const handleSubmit = (data: CreateFormValues | EditFormValues) => {
    onSubmit(data)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>
            {isEditing ? 'Editar Plantilla' : 'Nueva Plantilla'}
          </DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Nombre</FormLabel>
                  <FormControl>
                    <Input placeholder="Ej: Recordatorio de pago próximo" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {!isEditing && (
              <>
                <FormField
                  control={createForm.control}
                  name="code"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Código</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="Ej: payment_reminder_3_days"
                          {...field}
                          className="font-mono"
                        />
                      </FormControl>
                      <FormDescription>
                        Identificador único. Solo minúsculas, números y guiones bajos.
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="grid gap-4 sm:grid-cols-2">
                  <FormField
                    control={createForm.control}
                    name="channel"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Canal</FormLabel>
                        <Select onValueChange={field.onChange} value={field.value}>
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Seleccionar canal" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {NOTIFICATION_CHANNELS.map((channel) => (
                              <SelectItem key={channel.value} value={channel.value}>
                                {channel.label}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="notification_type"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Tipo</FormLabel>
                        <Select onValueChange={field.onChange} value={field.value}>
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Seleccionar tipo" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {NOTIFICATION_TYPES.map((type) => (
                              <SelectItem key={type.value} value={type.value}>
                                {type.label}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
              </>
            )}

            <FormField
              control={form.control}
              name="subject"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Asunto (opcional)</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="Ej: Recordatorio de pago - {{loan_number}}"
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    Para email. Puede usar variables como {'{{customer_name}}'}.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="content"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Contenido</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Escriba el contenido de la notificación..."
                      rows={6}
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    Variables disponibles: {'{{customer_name}}'}, {'{{loan_number}}'},{' '}
                    {'{{amount}}'}, {'{{due_date}}'}, {'{{branch_name}}'}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={isLoading}
              >
                Cancelar
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? 'Guardando...' : isEditing ? 'Guardar Cambios' : 'Crear Plantilla'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
