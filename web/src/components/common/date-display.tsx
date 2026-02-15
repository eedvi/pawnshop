import { formatDate, formatDateTime, formatRelativeTime } from '@/lib/format'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'

interface DateDisplayProps {
  date: string | undefined | null
  showTime?: boolean
  relative?: boolean
  className?: string
}

export function DateDisplay({
  date,
  showTime = false,
  relative = false,
  className,
}: DateDisplayProps) {
  if (!date) {
    return <span className={className}>-</span>
  }

  if (relative) {
    return (
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger className={className}>
            {formatRelativeTime(date)}
          </TooltipTrigger>
          <TooltipContent>
            {showTime ? formatDateTime(date) : formatDate(date)}
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    )
  }

  return (
    <span className={className}>
      {showTime ? formatDateTime(date) : formatDate(date)}
    </span>
  )
}
