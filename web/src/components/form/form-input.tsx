import { Control, FieldPath, FieldValues, useController } from 'react-hook-form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { cn } from '@/lib/utils'

interface FormInputProps<TFieldValues extends FieldValues>
  extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'name'> {
  control: Control<TFieldValues>
  name: FieldPath<TFieldValues>
  label?: string
  description?: string
  required?: boolean
}

export function FormInput<TFieldValues extends FieldValues = FieldValues>({
  control,
  name,
  label,
  description,
  className,
  required,
  type = 'text',
  ...props
}: FormInputProps<TFieldValues>) {
  const {
    field,
    fieldState: { error },
  } = useController({
    name,
    control,
  })

  // Handle number type conversion
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (type === 'number') {
      const value = e.target.value
      field.onChange(value === '' ? undefined : parseFloat(value))
    } else {
      field.onChange(e.target.value)
    }
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
        type={type}
        className={cn(error && 'border-destructive', className)}
        {...field}
        value={field.value ?? ''}
        onChange={handleChange}
        {...props}
      />
      {description && !error && (
        <p className="text-sm text-muted-foreground">{description}</p>
      )}
      {error && <p className="text-sm text-destructive">{error.message}</p>}
    </div>
  )
}
