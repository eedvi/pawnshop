import { cn } from '@/lib/utils'
import { STATUS_COLORS, STATUS_LABELS } from '@/lib/constants'

interface StatusBadgeProps {
  status: string
  className?: string
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const colorClass = STATUS_COLORS[status] || STATUS_COLORS['false']
  const label = STATUS_LABELS[status] || status

  return (
    <span
      className={cn(
        'inline-flex items-center rounded-md px-2 py-1 text-xs font-medium',
        colorClass,
        className
      )}
    >
      {label}
    </span>
  )
}
