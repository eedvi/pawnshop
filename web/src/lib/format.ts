import { format, formatDistanceToNow, parseISO, isValid } from 'date-fns'
import { es } from 'date-fns/locale'
import { CURRENCY_SYMBOL, CURRENCY_LOCALE, DATE_FORMAT, DATE_TIME_FORMAT } from './constants'

/**
 * Format a number as currency
 */
export function formatCurrency(amount: number | undefined | null): string {
  if (amount === undefined || amount === null) return `${CURRENCY_SYMBOL} 0.00`

  return `${CURRENCY_SYMBOL} ${amount.toLocaleString(CURRENCY_LOCALE, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })}`
}

/**
 * Format a number with thousands separator
 */
export function formatNumber(value: number | undefined | null, decimals = 0): string {
  if (value === undefined || value === null) return '0'

  return value.toLocaleString(CURRENCY_LOCALE, {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  })
}

/**
 * Format a percentage
 */
export function formatPercent(value: number | undefined | null, decimals = 1): string {
  if (value === undefined || value === null) return '0%'

  return `${value.toFixed(decimals)}%`
}

/**
 * Format a date string
 */
export function formatDate(dateString: string | undefined | null): string {
  if (!dateString) return '-'

  try {
    const date = parseISO(dateString)
    if (!isValid(date)) return '-'
    return format(date, DATE_FORMAT, { locale: es })
  } catch {
    return '-'
  }
}

/**
 * Format a date and time string
 */
export function formatDateTime(dateString: string | undefined | null): string {
  if (!dateString) return '-'

  try {
    const date = parseISO(dateString)
    if (!isValid(date)) return '-'
    return format(date, DATE_TIME_FORMAT, { locale: es })
  } catch {
    return '-'
  }
}

/**
 * Format a date as relative time (e.g., "hace 2 días")
 */
export function formatRelativeTime(dateString: string | undefined | null): string {
  if (!dateString) return '-'

  try {
    const date = parseISO(dateString)
    if (!isValid(date)) return '-'
    return formatDistanceToNow(date, { addSuffix: true, locale: es })
  } catch {
    return '-'
  }
}

/**
 * Format a phone number
 */
export function formatPhone(phone: string | undefined | null): string {
  if (!phone) return '-'

  // Remove non-numeric characters
  const cleaned = phone.replace(/\D/g, '')

  // Format as Guatemala phone: XXXX-XXXX
  if (cleaned.length === 8) {
    return `${cleaned.slice(0, 4)}-${cleaned.slice(4)}`
  }

  return phone
}

/**
 * Get initials from a name
 */
export function getInitials(name: string | undefined | null): string {
  if (!name) return '?'

  const parts = name.trim().split(' ')
  if (parts.length === 1) {
    return parts[0].charAt(0).toUpperCase()
  }
  return (parts[0].charAt(0) + parts[parts.length - 1].charAt(0)).toUpperCase()
}

/**
 * Truncate text with ellipsis
 */
export function truncate(text: string | undefined | null, maxLength: number): string {
  if (!text) return ''
  if (text.length <= maxLength) return text
  return text.slice(0, maxLength - 3) + '...'
}

/**
 * Format file size
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'

  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`
}

/**
 * Format days overdue
 */
export function formatDaysOverdue(days: number): string {
  if (days === 0) return 'Al día'
  if (days === 1) return '1 día vencido'
  return `${days} días vencido`
}

/**
 * Format loan term
 */
export function formatLoanTerm(days: number): string {
  if (days < 30) return `${days} días`

  const months = Math.floor(days / 30)
  const remainingDays = days % 30

  if (remainingDays === 0) {
    return months === 1 ? '1 mes' : `${months} meses`
  }

  const monthStr = months === 1 ? '1 mes' : `${months} meses`
  const dayStr = remainingDays === 1 ? '1 día' : `${remainingDays} días`

  return `${monthStr} y ${dayStr}`
}
