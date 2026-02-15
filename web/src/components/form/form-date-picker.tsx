import { Control, FieldPath, FieldValues, useController } from 'react-hook-form'
import { format, parseISO } from 'date-fns'
import { es } from 'date-fns/locale'
import { CalendarIcon } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import { Label } from '@/components/ui/label'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { DATE_FORMAT } from '@/lib/constants'

interface FormDatePickerProps<TFieldValues extends FieldValues> {
  control: Control<TFieldValues>
  name: FieldPath<TFieldValues>
  label?: string
  description?: string
  placeholder?: string
  disabled?: boolean
  required?: boolean
  minDate?: Date
  maxDate?: Date
  fromYear?: number
  toYear?: number
  captionLayout?: 'buttons' | 'dropdown' | 'dropdown-buttons'
}

export function FormDatePicker<TFieldValues extends FieldValues = FieldValues>({
  control,
  name,
  label,
  description,
  placeholder = 'Seleccionar fecha',
  disabled,
  required,
  minDate,
  maxDate,
  fromYear,
  toYear,
  captionLayout = 'buttons',
}: FormDatePickerProps<TFieldValues>) {
  const {
    field,
    fieldState: { error },
  } = useController({
    name,
    control,
  })

  const dateValue = field.value
    ? typeof field.value === 'string'
      ? parseISO(field.value)
      : field.value
    : undefined

  return (
    <div className="space-y-2">
      {label && (
        <Label htmlFor={name} className={cn(error && 'text-destructive')}>
          {label}
          {required && <span className="text-destructive ml-1">*</span>}
        </Label>
      )}
      <Popover>
        <PopoverTrigger asChild>
          <Button
            id={name}
            variant="outline"
            className={cn(
              'w-full justify-start text-left font-normal',
              !dateValue && 'text-muted-foreground',
              error && 'border-destructive'
            )}
            disabled={disabled}
          >
            <CalendarIcon className="mr-2 h-4 w-4" />
            {dateValue ? (
              format(dateValue, DATE_FORMAT, { locale: es })
            ) : (
              <span>{placeholder}</span>
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-auto p-0" align="start">
          <Calendar
            mode="single"
            selected={dateValue}
            onSelect={(date) => {
              field.onChange(date?.toISOString().split('T')[0])
            }}
            disabled={(date) => {
              if (minDate && date < minDate) return true
              if (maxDate && date > maxDate) return true
              return false
            }}
            captionLayout={captionLayout}
            fromYear={fromYear}
            toYear={toYear}
            autoFocus
          />
        </PopoverContent>
      </Popover>
      {description && !error && (
        <p className="text-sm text-muted-foreground">{description}</p>
      )}
      {error && <p className="text-sm text-destructive">{error.message}</p>}
    </div>
  )
}
