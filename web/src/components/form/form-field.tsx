import * as React from 'react'
import {
  FieldPath,
  FieldValues,
  useFormContext,
} from 'react-hook-form'
import { cn } from '@/lib/utils'
import { Label } from '@/components/ui/label'

interface FormFieldContextValue {
  name: string
  error?: string
}

const FormFieldContext = React.createContext<FormFieldContextValue>({
  name: '',
})

interface FormFieldProps<
  TFieldValues extends FieldValues = FieldValues,
  TName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>
> {
  name: TName
  children: React.ReactNode
  label?: string
  description?: string
}

export function FormField<
  TFieldValues extends FieldValues = FieldValues,
  TName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>
>({
  name,
  children,
  label,
  description,
}: FormFieldProps<TFieldValues, TName>) {
  const { formState } = useFormContext<TFieldValues>()
  const error = formState.errors[name]?.message as string | undefined

  return (
    <FormFieldContext.Provider value={{ name, error }}>
      <div className="space-y-2">
        {label && (
          <Label htmlFor={name} className={cn(error && 'text-destructive')}>
            {label}
          </Label>
        )}
        {children}
        {description && !error && (
          <p className="text-sm text-muted-foreground">{description}</p>
        )}
        {error && <p className="text-sm text-destructive">{error}</p>}
      </div>
    </FormFieldContext.Provider>
  )
}

export function useFormField() {
  const context = React.useContext(FormFieldContext)
  if (!context) {
    throw new Error('useFormField must be used within a FormField')
  }
  return context
}
