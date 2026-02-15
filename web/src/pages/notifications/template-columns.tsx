import { ColumnDef } from '@tanstack/react-table'
import { MoreHorizontal, Pencil, Trash, Power, PowerOff } from 'lucide-react'

import { NotificationTemplate, NOTIFICATION_CHANNELS } from '@/types'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

interface TemplateColumnOptions {
  onEdit: (template: NotificationTemplate) => void
  onToggle: (template: NotificationTemplate) => void
  onDelete: (template: NotificationTemplate) => void
}

export function getTemplateColumns(options: TemplateColumnOptions): ColumnDef<NotificationTemplate>[] {
  const { onEdit, onToggle, onDelete } = options

  return [
    {
      accessorKey: 'name',
      header: 'Nombre',
      cell: ({ row }) => (
        <div>
          <p className="font-medium">{row.original.name}</p>
          <p className="text-sm text-muted-foreground font-mono">{row.original.code}</p>
        </div>
      ),
    },
    {
      accessorKey: 'channel',
      header: 'Canal',
      cell: ({ row }) => {
        const channel = NOTIFICATION_CHANNELS.find((c) => c.value === row.original.channel)
        return <Badge variant="outline">{channel?.label || row.original.channel}</Badge>
      },
    },
    {
      accessorKey: 'notification_type',
      header: 'Tipo',
      cell: ({ row }) => {
        const typeLabels: Record<string, string> = {
          payment_reminder: 'Recordatorio de pago',
          overdue_notice: 'Aviso de mora',
          loan_expiry: 'Vencimiento de préstamo',
          promotional: 'Promocional',
          system: 'Sistema',
        }
        return <span>{typeLabels[row.original.notification_type] || row.original.notification_type}</span>
      },
    },
    {
      accessorKey: 'subject',
      header: 'Asunto',
      cell: ({ row }) => (
        <span className="text-muted-foreground">{row.original.subject || '-'}</span>
      ),
    },
    {
      accessorKey: 'is_active',
      header: 'Estado',
      cell: ({ row }) =>
        row.original.is_active ? (
          <Badge variant="default">Activa</Badge>
        ) : (
          <Badge variant="secondary">Inactiva</Badge>
        ),
    },
    {
      id: 'actions',
      cell: ({ row }) => {
        const template = row.original

        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Abrir menú</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => onEdit(template)}>
                <Pencil className="mr-2 h-4 w-4" />
                Editar
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => onToggle(template)}>
                {template.is_active ? (
                  <>
                    <PowerOff className="mr-2 h-4 w-4" />
                    Desactivar
                  </>
                ) : (
                  <>
                    <Power className="mr-2 h-4 w-4" />
                    Activar
                  </>
                )}
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => onDelete(template)}
                className="text-destructive"
              >
                <Trash className="mr-2 h-4 w-4" />
                Eliminar
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
    },
  ]
}
