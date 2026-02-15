import { useState } from 'react'
import { formatDistanceToNow } from 'date-fns'
import { es } from 'date-fns/locale'
import {
  Bell,
  BellOff,
  Check,
  CheckCheck,
  Loader2,
  AlertTriangle,
  Info,
  DollarSign,
  Clock,
  Package,
} from 'lucide-react'

import { InternalNotification } from '@/types'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'
import {
  useInternalNotifications,
  useMarkAsRead,
  useMarkAllAsRead,
  useUnreadCount,
} from '@/hooks/use-notifications'

function getNotificationIcon(type: string) {
  switch (type) {
    case 'payment':
      return DollarSign
    case 'overdue':
      return AlertTriangle
    case 'loan':
      return Clock
    case 'transfer':
      return Package
    default:
      return Info
  }
}

function NotificationItem({
  notification,
  onMarkAsRead,
  isMarkingRead,
}: {
  notification: InternalNotification
  onMarkAsRead: (id: number) => void
  isMarkingRead: boolean
}) {
  const Icon = getNotificationIcon(notification.notification_type)

  return (
    <div
      className={cn(
        'flex gap-4 p-4 border-b last:border-b-0 transition-colors',
        !notification.is_read && 'bg-muted/30'
      )}
    >
      <div
        className={cn(
          'flex-shrink-0 h-10 w-10 rounded-full flex items-center justify-center',
          notification.is_read ? 'bg-muted' : 'bg-primary/10'
        )}
      >
        <Icon
          className={cn(
            'h-5 w-5',
            notification.is_read ? 'text-muted-foreground' : 'text-primary'
          )}
        />
      </div>

      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-2">
          <div>
            <p
              className={cn(
                'font-medium',
                notification.is_read && 'text-muted-foreground'
              )}
            >
              {notification.title}
            </p>
            <p className="text-sm text-muted-foreground mt-1">
              {notification.message}
            </p>
          </div>
          {!notification.is_read && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onMarkAsRead(notification.id)}
              disabled={isMarkingRead}
            >
              <Check className="h-4 w-4" />
            </Button>
          )}
        </div>
        <p className="text-xs text-muted-foreground mt-2">
          {formatDistanceToNow(new Date(notification.created_at), {
            addSuffix: true,
            locale: es,
          })}
        </p>
      </div>
    </div>
  )
}

export function InternalNotificationList() {
  const [showUnreadOnly, setShowUnreadOnly] = useState(false)
  const [page, setPage] = useState(1)

  const { data: notificationsResponse, isLoading } = useInternalNotifications({
    page,
    per_page: 20,
    is_read: showUnreadOnly ? false : undefined,
  })

  const { data: unreadCount } = useUnreadCount()
  const markAsReadMutation = useMarkAsRead()
  const markAllAsReadMutation = useMarkAllAsRead()

  const notifications = notificationsResponse?.data || []
  const pagination = notificationsResponse?.meta?.pagination

  const handleMarkAsRead = (id: number) => {
    markAsReadMutation.mutate(id)
  }

  const handleMarkAllAsRead = () => {
    markAllAsReadMutation.mutate()
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
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Button
            variant={showUnreadOnly ? 'default' : 'outline'}
            size="sm"
            onClick={() => {
              setShowUnreadOnly(!showUnreadOnly)
              setPage(1)
            }}
          >
            {showUnreadOnly ? (
              <>
                <BellOff className="mr-2 h-4 w-4" />
                Solo no leídas
              </>
            ) : (
              <>
                <Bell className="mr-2 h-4 w-4" />
                Todas
              </>
            )}
          </Button>
          {unreadCount !== undefined && unreadCount > 0 && (
            <Badge variant="secondary">{unreadCount} sin leer</Badge>
          )}
        </div>

        {unreadCount !== undefined && unreadCount > 0 && (
          <Button
            variant="outline"
            size="sm"
            onClick={handleMarkAllAsRead}
            disabled={markAllAsReadMutation.isPending}
          >
            <CheckCheck className="mr-2 h-4 w-4" />
            Marcar todas como leídas
          </Button>
        )}
      </div>

      {/* Notifications List */}
      <Card>
        <CardContent className="p-0">
          {notifications.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 text-muted-foreground">
              <Bell className="h-12 w-12 mb-4 opacity-50" />
              <p>No hay notificaciones</p>
            </div>
          ) : (
            <>
              {notifications.map((notification) => (
                <NotificationItem
                  key={notification.id}
                  notification={notification}
                  onMarkAsRead={handleMarkAsRead}
                  isMarkingRead={markAsReadMutation.isPending}
                />
              ))}
            </>
          )}
        </CardContent>
      </Card>

      {/* Pagination */}
      {pagination && pagination.total_pages > 1 && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-muted-foreground">
            Mostrando {notifications.length} de {pagination.total}
          </p>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1}
            >
              Anterior
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setPage((p) => Math.min(pagination.total_pages, p + 1))}
              disabled={page >= pagination.total_pages}
            >
              Siguiente
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}
