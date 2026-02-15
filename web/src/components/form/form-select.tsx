import { Control, FieldPath, FieldValues, useController } from 'react-hook-form'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Label } from '@/components/ui/label'
import { cn } from '@/lib/utils'

interface SelectOption {
  value: string
  label: string
}

interface FormSelectProps<TFieldValues extends FieldValues> {
  control: Control<TFieldValues>
  name: FieldPath<TFieldValues>
  label?: string
  description?: string
  placeholder?: string
  options: SelectOption[]
  required?: boolean
  disabled?: boolean
  className?: string
}

export function FormSelect<TFieldValues extends FieldValues = FieldValues>({
  control,
  name,
  label,
  description,
  placeholder = 'Seleccionar...',
  options,
  required,
  disabled,
  className,
}: FormSelectProps<TFieldValues>) {
  const {
    field,
    fieldState: { error },
  } = useController({
    name,
    control,
  })

  return (
    <div className="space-y-2">
      {label && (
        <Label htmlFor={name} className={cn(error && 'text-destructive')}>
          {label}
          {required && <span className="text-destructive ml-1">*</span>}
        </Label>
      )}
      <Select
        value={field.value != null ? field.value.toString() : undefined}
        onValueChange={field.onChange}
        disabled={disabled}
      >
        <SelectTrigger
          id={name}
          className={cn(error && 'border-destructive', className)}
        >
          <SelectValue placeholder={placeholder} />
        </SelectTrigger>
        <SelectContent>
          {options
            .filter((option) => option.value !== '')
            .map((option) => (
              <SelectItem key={option.value} value={option.value}>
                {option.label}
              </SelectItem>
            ))}
        </SelectContent>
      </Select>
      {description && !error && (
        <p className="text-sm text-muted-foreground">{description}</p>
      )}
      {error && <p className="text-sm text-destructive">{error.message}</p>}
    </div>
  )
}
