import { useQuery } from '@tanstack/react-query'
import { auditService } from '@/services/audit-service'
import type { AuditLogListParams } from '@/types'

export const auditKeys = {
  all: ['audit'] as const,
  lists: () => [...auditKeys.all, 'list'] as const,
  list: (params?: AuditLogListParams) => [...auditKeys.lists(), params] as const,
  details: () => [...auditKeys.all, 'detail'] as const,
  detail: (id: number) => [...auditKeys.details(), id] as const,
  stats: (params?: { date_from?: string; date_to?: string; branch_id?: number }) =>
    [...auditKeys.all, 'stats', params] as const,
}

export function useAuditLogs(params?: AuditLogListParams) {
  return useQuery({
    queryKey: auditKeys.list(params),
    queryFn: () => auditService.list(params),
  })
}

export function useAuditLog(id: number) {
  return useQuery({
    queryKey: auditKeys.detail(id),
    queryFn: () => auditService.getById(id),
    enabled: id > 0,
  })
}

export function useAuditStats(params?: { date_from?: string; date_to?: string; branch_id?: number }) {
  return useQuery({
    queryKey: auditKeys.stats(params),
    queryFn: () => auditService.getStats(params),
  })
}
