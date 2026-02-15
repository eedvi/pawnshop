import { Control, FieldPath, FieldValues, useController } from 'react-hook-form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { cn } from '@/lib/utils'
import { CURRENCY_SYMBOL } from '@/lib/constants'

interface FormCurrencyInputProps<TFieldValues extends FieldValues> {
  control: Control<TFieldValues>
  name: FieldPath<TFieldValues>
  label?: string
  description?: string
  placeholder?: string
  disabled?: boolean
  required?: boolean
  min?: number
  max?: number
}

export function FormCurrencyInput<TFieldValues extends FieldValues = FieldValues>({
  control,
  name,
  label,
  description,
  placeholder = '0.00',
  disabled,
  required,
  min,
  max,
}: FormCurrencyInputProps<TFieldValues>) {
  const {
    field,
    fieldState: { error },
  } = useController({
    name,
    control,
  })

  const formatCurrency = (value: number | string | null | undefined): string => {
    if (value === null || value === undefined) return ''
    const num = typeof value === 'string' ? parseFloat(value) : value
    if (isNaN(num)) return ''
    return num.toFixed(2)
  }

  const parseCurrency = (value: string): number => {
    const cleaned = value.replace(/[^\d.]/g, '')
    const num = parseFloat(cleaned)
    return isNaN(num) ? 0 : num
  }

  return (
    <div className="space-y-2">
      {label && (
        <Label htmlFor={name} className={cn(error && 'text-destructive')}>
          {label}
          {required && <span className="text-destructive ml-1">*</span>}
        </Label>
      )}
      <div className="relative">
        <span className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
          {CURRENCY_SYMBOL}
        </span>
        <Input
          id={name}
          type="text"
          inputMode="decimal"
          placeholder={placeholder}
          disabled={disabled}
          className={cn('pl-7', error && 'border-destructive')}
          value={field.value != null ? formatCurrency(field.value) : ''}
          onChange={(e) => {
            const value = parseCurrency(e.target.value)
            if (min !== undefined && value < min) return
            if (max !== undefined && value > max) return
            field.onChange(value)
          }}
          onBlur={(e) => {
            const value = parseCurrency(e.target.value)
            field.onChange(value)
            field.onBlur()
          }}
        />
      </div>
      {description && !error && (
        <p className="text-sm text-muted-foreground">{description}</p>
      )}
      {error && <p className="text-sm text-destructive">{error.message}</p>}
    </div>
  )
}
