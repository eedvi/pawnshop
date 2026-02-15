import { Control, FieldPath, FieldValues, useController } from 'react-hook-form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { cn } from '@/lib/utils'

interface FormPhoneInputProps<TFieldValues extends FieldValues> {
  control: Control<TFieldValues>
  name: FieldPath<TFieldValues>
  label?: string
  description?: string
  placeholder?: string
  disabled?: boolean
  required?: boolean
}

export function FormPhoneInput<TFieldValues extends FieldValues = FieldValues>({
  control,
  name,
  label,
  description,
  placeholder = '0000-0000',
  disabled,
  required,
}: FormPhoneInputProps<TFieldValues>) {
  const {
    field,
    fieldState: { error },
  } = useController({
    name,
    control,
  })

  const formatPhone = (value: string): string => {
    // Remove non-numeric characters
    const cleaned = value.replace(/\D/g, '')

    // Format as XXXX-XXXX for Guatemala
    if (cleaned.length <= 4) {
      return cleaned
    }
    return `${cleaned.slice(0, 4)}-${cleaned.slice(4, 8)}`
  }

  const parsePhone = (value: string): string => {
    return value.replace(/\D/g, '')
  }

  return (
    <div className="space-y-2">
      {label && (
        <Label htmlFor={name} className={cn(error && 'text-destructive')}>
          {label}
          {required && <span className="text-destructive ml-1">*</span>}
        </Label>
      )}
      <Input
        id={name}
        type="tel"
        inputMode="tel"
        placeholder={placeholder}
        disabled={disabled}
        className={cn(error && 'border-destructive')}
        value={formatPhone(field.value || '')}
        onChange={(e) => {
          const raw = parsePhone(e.target.value)
          if (raw.length <= 8) {
            field.onChange(raw)
          }
        }}
        maxLength={9} // XXXX-XXXX format
      />
      {description && !error && (
        <p className="text-sm text-muted-foreground">{description}</p>
      )}
      {error && <p className="text-sm text-destructive">{error.message}</p>}
    </div>
  )
}
