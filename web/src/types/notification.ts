// Notification types - mirrors internal/domain/notification.go

export type NotificationChannel = 'email' | 'sms' | 'whatsapp' | 'internal'
export type NotificationStatus = 'pending' | 'sent' | 'delivered' | 'failed' | 'cancelled'
export type NotificationType = 'payment_reminder' | 'overdue_notice' | 'loan_expiry' | 'promotional' | 'system'

export const NOTIFICATION_CHANNELS: { value: NotificationChannel; label: string }[] = [
  { value: 'email', label: 'Email' },
  { value: 'sms', label: 'SMS' },
  { value: 'whatsapp', label: 'WhatsApp' },
  { value: 'internal', label: 'Interna' },
]

export const NOTIFICATION_STATUSES: { value: NotificationStatus; label: string; color: string }[] = [
  { value: 'pending', label: 'Pendiente', color: 'yellow' },
  { value: 'sent', label: 'Enviada', color: 'blue' },
  { value: 'delivered', label: 'Entregada', color: 'green' },
  { value: 'failed', label: 'Fallida', color: 'red' },
  { value: 'cancelled', label: 'Cancelada', color: 'gray' },
]

export interface NotificationTemplate {
  id: number
  name: string
  code: string
  channel: NotificationChannel
  notification_type: NotificationType
  subject?: string
  content: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Notification {
  id: number
  customer_id: number
  branch_id?: number
  template_id?: number
  channel: NotificationChannel
  notification_type: NotificationType
  subject?: string
  content: string
  status: NotificationStatus
  scheduled_at?: string
  sent_at?: string
  delivered_at?: string
  failed_at?: string
  error_message?: string
  reference_type?: string
  reference_id?: number
  created_by?: number
  created_at: string
  updated_at: string
}

export interface InternalNotification {
  id: number
  user_id: number
  branch_id?: number
  title: string
  message: string
  notification_type: string
  reference_type?: string
  reference_id?: number
  is_read: boolean
  read_at?: string
  created_by?: number
  created_at: string
}

export interface CustomerNotificationPreference {
  id: number
  customer_id: number
  email_enabled: boolean
  sms_enabled: boolean
  whatsapp_enabled: boolean
  payment_reminders: boolean
  overdue_notices: boolean
  promotional: boolean
  created_at: string
  updated_at: string
}

export interface CreateNotificationTemplateInput {
  name: string
  code: string
  channel: NotificationChannel
  notification_type: NotificationType
  subject?: string
  content: string
}

export interface UpdateNotificationTemplateInput {
  name?: string
  subject?: string
  content?: string
  is_active?: boolean
}

export interface CreateNotificationInput {
  customer_id: number
  channel: NotificationChannel
  notification_type: NotificationType
  subject?: string
  content: string
  scheduled_at?: string
  reference_type?: string
  reference_id?: number
}

export interface CreateFromTemplateInput {
  customer_id: number
  template_code: string
  variables?: Record<string, string>
  scheduled_at?: string
  reference_type?: string
  reference_id?: number
}

export interface CreateInternalNotificationInput {
  user_id: number
  title: string
  message: string
  notification_type?: string
  reference_type?: string
  reference_id?: number
}

export interface UpdateNotificationPreferencesInput {
  email_enabled?: boolean
  sms_enabled?: boolean
  whatsapp_enabled?: boolean
  payment_reminders?: boolean
  overdue_notices?: boolean
  promotional?: boolean
}

export interface NotificationListParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  customer_id?: number
  branch_id?: number
  channel?: NotificationChannel
  status?: NotificationStatus
  notification_type?: NotificationType
  date_from?: string
  date_to?: string
}
