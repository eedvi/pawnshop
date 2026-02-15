import { apiGet, apiGetPaginated, apiPost, apiPut, apiDelete } from '@/lib/api-client'
import type {
  Notification,
  InternalNotification,
  NotificationTemplate,
  NotificationListParams,
  CreateNotificationInput,
  CreateFromTemplateInput,
  CreateNotificationTemplateInput,
  UpdateNotificationTemplateInput,
} from '@/types'

export const notificationService = {
  // Notifications
  list: (params?: NotificationListParams) =>
    apiGetPaginated<Notification>('/notifications', params),

  getById: (id: number) =>
    apiGet<Notification>(`/notifications/${id}`),

  create: (input: CreateNotificationInput) =>
    apiPost<Notification>('/notifications', input),

  createFromTemplate: (input: CreateFromTemplateInput) =>
    apiPost<Notification>('/notifications/from-template', input),

  cancel: (id: number) =>
    apiPost<Notification>(`/notifications/${id}/cancel`, {}),

  retry: (id: number) =>
    apiPost<Notification>(`/notifications/${id}/retry`, {}),

  // Internal Notifications
  listInternal: (params?: { page?: number; per_page?: number; is_read?: boolean }) =>
    apiGetPaginated<InternalNotification>('/notifications/internal', params),

  getInternalById: (id: number) =>
    apiGet<InternalNotification>(`/notifications/internal/${id}`),

  markAsRead: (id: number) =>
    apiPost<InternalNotification>(`/notifications/internal/${id}/read`, {}),

  markAllAsRead: () =>
    apiPost<void>('/notifications/internal/read-all', {}),

  getUnreadCount: () =>
    apiGet<{ count: number }>('/notifications/internal/unread-count'),

  // Templates
  listTemplates: () =>
    apiGet<NotificationTemplate[]>('/notifications/templates'),

  getTemplateById: (id: number) =>
    apiGet<NotificationTemplate>(`/notifications/templates/${id}`),

  createTemplate: (input: CreateNotificationTemplateInput) =>
    apiPost<NotificationTemplate>('/notifications/templates', input),

  updateTemplate: (id: number, input: UpdateNotificationTemplateInput) =>
    apiPut<NotificationTemplate>(`/notifications/templates/${id}`, input),

  deleteTemplate: (id: number) =>
    apiDelete<void>(`/notifications/templates/${id}`),

  toggleTemplate: (id: number) =>
    apiPost<NotificationTemplate>(`/notifications/templates/${id}/toggle`, {}),
}
