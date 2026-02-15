import { Control, FieldPath, FieldValues, useController } from 'react-hook-form'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { cn } from '@/lib/utils'

interface FormSwitchProps<TFieldValues extends FieldValues> {
  control: Control<TFieldValues>
  name: FieldPath<TFieldValues>
  label?: string
  description?: string
  disabled?: boolean
  className?: string
}

export function FormSwitch<TFieldValues extends FieldValues = FieldValues>({
  control,
  name,
  label,
  description,
  disabled,
  className,
}: FormSwitchProps<TFieldValues>) {
  const {
    field,
    fieldState: { error },
  } = useController({
    name,
    control,
  })

  return (
    <div className={cn('flex items-center justify-between rounded-lg border p-4', className)}>
      <div className="space-y-0.5">
        {label && (
          <Label htmlFor={name} className={cn('text-base', error && 'text-destructive')}>
            {label}
          </Label>
        )}
        {description && (
          <p className="text-sm text-muted-foreground">{description}</p>
        )}
        {error && <p className="text-sm text-destructive">{error.message}</p>}
      </div>
      <Switch
        id={name}
        checked={field.value}
        onCheckedChange={field.onChange}
        disabled={disabled}
      />
    </div>
  )
}
