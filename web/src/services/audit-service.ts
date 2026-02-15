import { apiGet, apiGetPaginated } from '@/lib/api-client'
import type { AuditLog, AuditLogListParams } from '@/types'

export const auditService = {
  list: (params?: AuditLogListParams) =>
    apiGetPaginated<AuditLog>('/audit-logs', params),

  getById: (id: number) =>
    apiGet<AuditLog>(`/audit-logs/${id}`),

  // Get stats/summary for audit dashboard
  getStats: (params?: { date_from?: string; date_to?: string; branch_id?: number }) =>
    apiGet<{
      total_actions: number
      actions_by_type: Record<string, number>
      actions_by_entity: Record<string, number>
      active_users: number
    }>('/audit-logs/stats', params),
}
