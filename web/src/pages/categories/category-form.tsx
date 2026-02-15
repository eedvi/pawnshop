import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2 } from 'lucide-react'

import { Category } from '@/types'
import { Button } from '@/components/ui/button'
import { Form } from '@/components/ui/form'
import { FormInput, FormTextarea, FormSelect, FormSwitch } from '@/components/form'
import { categoryFormSchema, CategoryFormValues, defaultCategoryValues } from './schemas'

interface CategoryFormProps {
  category?: Category
  categories: Category[]
  onSubmit: (values: CategoryFormValues) => void
  onCancel: () => void
  isLoading?: boolean
}

export function CategoryForm({
  category,
  categories,
  onSubmit,
  onCancel,
  isLoading,
}: CategoryFormProps) {
  const form = useForm<CategoryFormValues>({
    resolver: zodResolver(categoryFormSchema),
    defaultValues: category
      ? {
          name: category.name,
          description: category.description || '',
          parent_id: category.parent_id ?? 'none',
          icon: category.icon || '',
          default_interest_rate: category.default_interest_rate,
          min_loan_amount: category.min_loan_amount ?? null,
          max_loan_amount: category.max_loan_amount ?? null,
          loan_to_value_ratio: category.loan_to_value_ratio,
          sort_order: category.sort_order,
          is_active: category.is_active,
        }
      : defaultCategoryValues,
  })

  const isEditing = !!category

  // Filter out the current category and its children from parent options
  const getParentOptions = () => {
    const flattenCategories = (cats: Category[]): Category[] => {
      const result: Category[] = []
      for (const cat of cats) {
        result.push(cat)
        if (cat.children) {
          result.push(...flattenCategories(cat.children))
        }
      }
      return result
    }

    const allCategories = flattenCategories(categories)

    // When editing, exclude current category and its descendants
    const excludeIds = new Set<number>()
    if (category) {
      const addDescendants = (cats: Category[]) => {
        for (const cat of cats) {
          if (cat.id === category.id) {
            excludeIds.add(cat.id)
            if (cat.children) {
              cat.children.forEach((c) => excludeIds.add(c.id))
            }
          }
          if (cat.children) {
            addDescendants(cat.children)
          }
        }
      }
      addDescendants(categories)
      excludeIds.add(category.id)
    }

    return [
      { value: 'none', label: 'Ninguna (categoría raíz)' },
      ...allCategories
        .filter((cat) => !excludeIds.has(cat.id))
        .map((cat) => ({
          value: cat.id.toString(),
          label: cat.parent_id ? `  └ ${cat.name}` : cat.name,
        })),
    ]
  }

  const handleSubmit = (values: CategoryFormValues) => {
    onSubmit({
      ...values,
      parent_id: values.parent_id === 'none' ? undefined : values.parent_id,
    })
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
        <FormInput
          control={form.control}
          name="name"
          label="Nombre"
          placeholder="Electrónicos"
          required
        />

        <FormTextarea
          control={form.control}
          name="description"
          label="Descripción"
          placeholder="Categoría para artículos electrónicos..."
          rows={3}
        />

        <FormSelect
          control={form.control}
          name="parent_id"
          label="Categoría Padre"
          options={getParentOptions()}
          placeholder="Seleccionar categoría padre"
        />

        <FormInput
          control={form.control}
          name="icon"
          label="Icono"
          placeholder="laptop, phone, watch..."
          description="Nombre del icono de Lucide"
        />

        <div className="grid gap-4 sm:grid-cols-2">
          <FormInput
            control={form.control}
            name="default_interest_rate"
            label="Tasa de Interés (%)"
            type="number"
            min={0}
            max={100}
            step={0.1}
          />
          <FormInput
            control={form.control}
            name="loan_to_value_ratio"
            label="Ratio Préstamo/Valor (%)"
            type="number"
            min={0}
            max={100}
            step={1}
          />
        </div>

        <div className="grid gap-4 sm:grid-cols-2">
          <FormInput
            control={form.control}
            name="min_loan_amount"
            label="Monto Mínimo (Q)"
            type="number"
            min={0}
            step={0.01}
          />
          <FormInput
            control={form.control}
            name="max_loan_amount"
            label="Monto Máximo (Q)"
            type="number"
            min={0}
            step={0.01}
          />
        </div>

        <FormInput
          control={form.control}
          name="sort_order"
          label="Orden"
          type="number"
          min={0}
        />

        {isEditing && (
          <FormSwitch
            control={form.control}
            name="is_active"
            label="Categoría activa"
            description="Las categorías inactivas no aparecen en los formularios"
          />
        )}

        <div className="flex justify-end gap-4 pt-4">
          <Button type="button" variant="outline" onClick={onCancel}>
            Cancelar
          </Button>
          <Button type="submit" disabled={isLoading}>
            {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {isEditing ? 'Guardar Cambios' : 'Crear Categoría'}
          </Button>
        </div>
      </form>
    </Form>
  )
}
