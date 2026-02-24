// Audit types - mirrors internal/domain/audit.go

export type AuditAction = 'create' | 'update' | 'delete' | 'login' | 'logout' | 'view' | 'export' | 'approve' | 'reject' | 'other'

export const AUDIT_ACTIONS: { value: AuditAction; label: string }[] = [
  { value: 'create', label: 'Crear' },
  { value: 'update', label: 'Actualizar' },
  { value: 'delete', label: 'Eliminar' },
  { value: 'login', label: 'Iniciar Sesión' },
  { value: 'logout', label: 'Cerrar Sesión' },
  { value: 'view', label: 'Ver' },
  { value: 'export', label: 'Exportar' },
  { value: 'approve', label: 'Aprobar' },
  { value: 'reject', label: 'Rechazar' },
  { value: 'other', label: 'Otro' },
]

export interface AuditLog {
  id: number
  user_id?: number
  branch_id?: number
  action: AuditAction
  entity_type: string
  entity_id?: number
  old_values?: Record<string, unknown>
  new_values?: Record<string, unknown>
  ip_address?: string
  user_agent?: string
  description?: string
  created_at: string

  // Relations
  user_name?: string
  branch_name?: string
}

export interface AuditLogListParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  user_id?: number
  branch_id?: number
  action?: AuditAction
  entity_type?: string
  entity_id?: number
  date_from?: string
  date_to?: string
}

export const ENTITY_TYPES = [
  'customer',
  'item',
  'loan',
  'payment',
  'sale',
  'user',
  'branch',
  'category',
  'role',
  'cash_register',
  'cash_session',
  'cash_movement',
  'transfer',
  'expense',
  'notification',
  'setting',
]

export interface TopUserStat {
  user_id: number
  user_name: string
  count: number
}

export interface AuditStats {
  total_actions: number
  actions_by_type: Record<string, number>
  actions_by_entity: Record<string, number>
  active_users: number
  top_users: TopUserStat[]
  recent_critical: AuditLog[]
}
