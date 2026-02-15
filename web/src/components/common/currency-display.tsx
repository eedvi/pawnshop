import { formatCurrency } from '@/lib/format'
import { cn } from '@/lib/utils'

interface CurrencyDisplayProps {
  amount: number | undefined | null
  className?: string
  showSign?: boolean
}

export function CurrencyDisplay({
  amount,
  className,
  showSign = false,
}: CurrencyDisplayProps) {
  const isNegative = amount !== undefined && amount !== null && amount < 0
  const isPositive = amount !== undefined && amount !== null && amount > 0

  return (
    <span
      className={cn(
        showSign && isNegative && 'text-destructive',
        showSign && isPositive && 'text-green-600 dark:text-green-400',
        className
      )}
    >
      {showSign && isPositive && '+'}
      {formatCurrency(amount)}
    </span>
  )
}
