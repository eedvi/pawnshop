import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { notificationService } from '@/services/notification-service'
import type {
  NotificationListParams,
  CreateNotificationInput,
  CreateFromTemplateInput,
  CreateNotificationTemplateInput,
  UpdateNotificationTemplateInput,
} from '@/types'

export const notificationKeys = {
  all: ['notifications'] as const,
  lists: () => [...notificationKeys.all, 'list'] as const,
  list: (params?: NotificationListParams) => [...notificationKeys.lists(), params] as const,
  details: () => [...notificationKeys.all, 'detail'] as const,
  detail: (id: number) => [...notificationKeys.details(), id] as const,
  internal: () => [...notificationKeys.all, 'internal'] as const,
  internalList: (params?: { page?: number; per_page?: number; is_read?: boolean }) =>
    [...notificationKeys.internal(), 'list', params] as const,
  internalDetail: (id: number) => [...notificationKeys.internal(), 'detail', id] as const,
  unreadCount: () => [...notificationKeys.internal(), 'unread'] as const,
  templates: () => [...notificationKeys.all, 'templates'] as const,
  templateDetail: (id: number) => [...notificationKeys.templates(), id] as const,
}

// Notifications
export function useNotifications(params?: NotificationListParams) {
  return useQuery({
    queryKey: notificationKeys.list(params),
    queryFn: () => notificationService.list(params),
  })
}

export function useNotification(id: number) {
  return useQuery({
    queryKey: notificationKeys.detail(id),
    queryFn: () => notificationService.getById(id),
    enabled: id > 0,
  })
}

export function useCreateNotification() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateNotificationInput) => notificationService.create(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: notificationKeys.lists() })
    },
  })
}

export function useCreateFromTemplate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateFromTemplateInput) => notificationService.createFromTemplate(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: notificationKeys.lists() })
    },
  })
}

export function useCancelNotification() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => notificationService.cancel(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: notificationKeys.lists() })
      queryClient.invalidateQueries({ queryKey: notificationKeys.detail(id) })
    },
  })
}

export function useRetryNotification() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => notificationService.retry(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: notificationKeys.lists() })
      queryClient.invalidateQueries({ queryKey: notificationKeys.detail(id) })
    },
  })
}

// Internal Notifications
export function useInternalNotifications(params?: { page?: number; per_page?: number; is_read?: boolean }) {
  return useQuery({
    queryKey: notificationKeys.internalList(params),
    queryFn: () => notificationService.listInternal(params),
  })
}

export function useInternalNotification(id: number) {
  return useQuery({
    queryKey: notificationKeys.internalDetail(id),
    queryFn: () => notificationService.getInternalById(id),
    enabled: id > 0,
  })
}

export function useUnreadCount() {
  return useQuery({
    queryKey: notificationKeys.unreadCount(),
    queryFn: async () => {
      const result = await notificationService.getUnreadCount()
      return (result as { count?: number })?.count || 0
    },
    refetchInterval: 60000, // Refetch every minute
  })
}

export function useMarkAsRead() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => notificationService.markAsRead(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: notificationKeys.internal() })
    },
  })
}

export function useMarkAllAsRead() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => notificationService.markAllAsRead(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: notificationKeys.internal() })
    },
  })
}

// Templates
export function useNotificationTemplates() {
  return useQuery({
    queryKey: notificationKeys.templates(),
    queryFn: () => notificationService.listTemplates(),
  })
}

export function useNotificationTemplate(id: number) {
  return useQuery({
    queryKey: notificationKeys.templateDetail(id),
    queryFn: () => notificationService.getTemplateById(id),
    enabled: id > 0,
  })
}

export function useCreateTemplate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateNotificationTemplateInput) => notificationService.createTemplate(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: notificationKeys.templates() })
    },
  })
}

export function useUpdateTemplate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: UpdateNotificationTemplateInput }) =>
      notificationService.updateTemplate(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: notificationKeys.templates() })
      queryClient.invalidateQueries({ queryKey: notificationKeys.templateDetail(id) })
    },
  })
}

export function useDeleteTemplate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => notificationService.deleteTemplate(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: notificationKeys.templates() })
    },
  })
}

export function useToggleTemplate() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => notificationService.toggleTemplate(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: notificationKeys.templates() })
      queryClient.invalidateQueries({ queryKey: notificationKeys.templateDetail(id) })
    },
  })
}
