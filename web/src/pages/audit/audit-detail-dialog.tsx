import { AuditLog, AUDIT_ACTIONS } from '@/types'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { formatDateTime } from '@/lib/format'

interface AuditDetailDialogProps {
  log: AuditLog | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

function JsonViewer({ data, title }: { data: Record<string, unknown> | undefined; title: string }) {
  if (!data || Object.keys(data).length === 0) {
    return (
      <div>
        <p className="text-sm font-medium mb-2">{title}</p>
        <p className="text-sm text-muted-foreground italic">Sin datos</p>
      </div>
    )
  }

  return (
    <div>
      <p className="text-sm font-medium mb-2">{title}</p>
      <div className="bg-muted rounded-lg p-3 overflow-x-auto">
        <pre className="text-xs font-mono">
          {JSON.stringify(data, null, 2)}
        </pre>
      </div>
    </div>
  )
}

function ChangesViewer({
  oldValues,
  newValues,
}: {
  oldValues?: Record<string, unknown>
  newValues?: Record<string, unknown>
}) {
  // Get all unique keys from both objects
  const allKeys = new Set([
    ...Object.keys(oldValues || {}),
    ...Object.keys(newValues || {}),
  ])

  if (allKeys.size === 0) {
    return <p className="text-sm text-muted-foreground">Sin cambios registrados</p>
  }

  const changes = Array.from(allKeys).map((key) => {
    const oldVal = oldValues?.[key]
    const newVal = newValues?.[key]
    const hasChanged = JSON.stringify(oldVal) !== JSON.stringify(newVal)

    return { key, oldVal, newVal, hasChanged }
  })

  return (
    <div className="space-y-3">
      {changes.map(({ key, oldVal, newVal, hasChanged }) => (
        <div key={key} className="grid grid-cols-3 gap-2 text-sm">
          <div className="font-medium">{key}</div>
          <div className={`font-mono text-xs ${hasChanged ? 'text-red-500 line-through' : 'text-muted-foreground'}`}>
            {oldVal !== undefined ? JSON.stringify(oldVal) : '-'}
          </div>
          <div className={`font-mono text-xs ${hasChanged ? 'text-green-600' : 'text-muted-foreground'}`}>
            {newVal !== undefined ? JSON.stringify(newVal) : '-'}
          </div>
        </div>
      ))}
    </div>
  )
}

export function AuditDetailDialog({ log, open, onOpenChange }: AuditDetailDialogProps) {
  if (!log) return null

  const action = AUDIT_ACTIONS.find((a) => a.value === log.action)
  const showChangesView = log.action === 'update' && log.old_values && log.new_values

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Detalles del Registro de Auditoría</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          {/* Basic Info */}
          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <p className="text-sm text-muted-foreground">Fecha/Hora</p>
              <p className="font-medium">{formatDateTime(log.created_at)}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Usuario</p>
              <p className="font-medium">{log.user_name || 'Sistema'}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Acción</p>
              <Badge variant="outline">{action?.label || log.action}</Badge>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Entidad</p>
              <p className="font-medium">
                {log.entity_type}
                {log.entity_id && <span className="text-muted-foreground"> #{log.entity_id}</span>}
              </p>
            </div>
          </div>

          {log.description && (
            <div>
              <p className="text-sm text-muted-foreground">Descripción</p>
              <p className="mt-1">{log.description}</p>
            </div>
          )}

          <Separator />

          {/* Technical Details */}
          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <p className="text-sm text-muted-foreground">Sucursal</p>
              <p>{log.branch_name || '-'}</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Dirección IP</p>
              <p className="font-mono text-sm">{log.ip_address || '-'}</p>
            </div>
          </div>

          {log.user_agent && (
            <div>
              <p className="text-sm text-muted-foreground">User Agent</p>
              <p className="text-xs font-mono text-muted-foreground break-all">
                {log.user_agent}
              </p>
            </div>
          )}

          <Separator />

          {/* Values */}
          {showChangesView ? (
            <div>
              <p className="text-sm font-medium mb-3">Cambios Realizados</p>
              <div className="grid grid-cols-3 gap-2 text-xs text-muted-foreground mb-2 font-medium">
                <div>Campo</div>
                <div>Valor Anterior</div>
                <div>Valor Nuevo</div>
              </div>
              <ChangesViewer oldValues={log.old_values} newValues={log.new_values} />
            </div>
          ) : (
            <div className="grid gap-4 sm:grid-cols-2">
              <JsonViewer data={log.old_values} title="Valores Anteriores" />
              <JsonViewer data={log.new_values} title="Valores Nuevos" />
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
