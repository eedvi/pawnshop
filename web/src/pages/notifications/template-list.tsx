import { useState, useMemo } from 'react'
import { Plus, Loader2 } from 'lucide-react'

import { NotificationTemplate } from '@/types'
import { Button } from '@/components/ui/button'
import { DataTable } from '@/components/data-table/data-table'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import {
  useNotificationTemplates,
  useCreateTemplate,
  useUpdateTemplate,
  useDeleteTemplate,
  useToggleTemplate,
} from '@/hooks/use-notifications'
import { getTemplateColumns } from './template-columns'
import { TemplateFormDialog } from './template-form-dialog'

export function TemplateList() {
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [editTemplate, setEditTemplate] = useState<NotificationTemplate | null>(null)
  const [deleteTemplate, setDeleteTemplate] = useState<NotificationTemplate | null>(null)
  const [toggleTemplate, setToggleTemplate] = useState<NotificationTemplate | null>(null)

  const { data: templates, isLoading } = useNotificationTemplates()
  const createMutation = useCreateTemplate()
  const updateMutation = useUpdateTemplate()
  const deleteMutation = useDeleteTemplate()
  const toggleMutation = useToggleTemplate()

  const columns = useMemo(
    () =>
      getTemplateColumns({
        onEdit: (template) => setEditTemplate(template),
        onToggle: (template) => setToggleTemplate(template),
        onDelete: (template) => setDeleteTemplate(template),
      }),
    []
  )

  const handleCreate = (data: {
    name: string
    code: string
    channel: 'email' | 'sms' | 'whatsapp' | 'internal'
    notification_type: 'payment_reminder' | 'overdue_notice' | 'loan_expiry' | 'promotional' | 'system'
    subject?: string
    content: string
  }) => {
    createMutation.mutate(
      {
        name: data.name,
        code: data.code,
        channel: data.channel,
        notification_type: data.notification_type,
        subject: data.subject || undefined,
        content: data.content,
      },
      {
        onSuccess: () => setCreateDialogOpen(false),
      }
    )
  }

  const handleUpdate = (data: { name: string; subject?: string; content: string }) => {
    if (editTemplate) {
      updateMutation.mutate(
        {
          id: editTemplate.id,
          input: {
            name: data.name,
            subject: data.subject || undefined,
            content: data.content,
          },
        },
        {
          onSuccess: () => setEditTemplate(null),
        }
      )
    }
  }

  const handleDelete = () => {
    if (deleteTemplate) {
      deleteMutation.mutate(deleteTemplate.id, {
        onSuccess: () => setDeleteTemplate(null),
      })
    }
  }

  const handleToggle = () => {
    if (toggleTemplate) {
      toggleMutation.mutate(toggleTemplate.id, {
        onSuccess: () => setToggleTemplate(null),
      })
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-end">
        <Button onClick={() => setCreateDialogOpen(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Nueva Plantilla
        </Button>
      </div>

      <DataTable
        columns={columns}
        data={templates || []}
        searchPlaceholder="Buscar plantillas..."
        searchColumn="name"
      />

      {/* Create Dialog */}
      <TemplateFormDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
        onSubmit={handleCreate}
        isLoading={createMutation.isPending}
      />

      {/* Edit Dialog */}
      <TemplateFormDialog
        open={!!editTemplate}
        onOpenChange={(open) => !open && setEditTemplate(null)}
        template={editTemplate}
        onSubmit={handleUpdate}
        isLoading={updateMutation.isPending}
      />

      {/* Delete Confirmation */}
      <ConfirmDialog
        open={!!deleteTemplate}
        onOpenChange={(open) => !open && setDeleteTemplate(null)}
        title="Eliminar Plantilla"
        description={
          deleteTemplate
            ? `¿Está seguro de eliminar la plantilla "${deleteTemplate.name}"? Esta acción no se puede deshacer.`
            : ''
        }
        confirmText="Eliminar"
        variant="destructive"
        onConfirm={handleDelete}
        isLoading={deleteMutation.isPending}
      />

      {/* Toggle Confirmation */}
      <ConfirmDialog
        open={!!toggleTemplate}
        onOpenChange={(open) => !open && setToggleTemplate(null)}
        title={toggleTemplate?.is_active ? 'Desactivar Plantilla' : 'Activar Plantilla'}
        description={
          toggleTemplate
            ? toggleTemplate.is_active
              ? `¿Está seguro de desactivar la plantilla "${toggleTemplate.name}"?`
              : `¿Está seguro de activar la plantilla "${toggleTemplate.name}"?`
            : ''
        }
        confirmText={toggleTemplate?.is_active ? 'Desactivar' : 'Activar'}
        onConfirm={handleToggle}
        isLoading={toggleMutation.isPending}
      />
    </div>
  )
}
