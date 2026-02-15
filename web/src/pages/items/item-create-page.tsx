import { useNavigate } from 'react-router-dom'
import { Info } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { ROUTES } from '@/routes/routes'
import { useCreateItem } from '@/hooks/use-items'
import { useBranchStore } from '@/stores/branch-store'
import { ItemForm } from './item-form'
import { ItemFormValues } from './schemas'

export default function ItemCreatePage() {
  const navigate = useNavigate()
  const { selectedBranchId } = useBranchStore()
  const createMutation = useCreateItem()

  const handleSubmit = (values: ItemFormValues) => {
    if (!selectedBranchId) return

    createMutation.mutate(
      {
        branch_id: selectedBranchId,
        name: values.name,
        description: values.description || undefined,
        brand: values.brand || undefined,
        model: values.model || undefined,
        serial_number: values.serial_number || undefined,
        color: values.color || undefined,
        condition: values.condition,
        category_id: values.category_id || undefined,
        customer_id: values.customer_id || undefined,
        appraised_value: values.appraised_value,
        loan_value: values.loan_value,
        sale_price: values.sale_price ?? undefined,
        weight: values.weight ?? undefined,
        purity: values.purity || undefined,
        acquisition_type: values.acquisition_type,
        acquisition_price: values.acquisition_price ?? undefined,
        notes: values.notes || undefined,
        tags: values.tags?.length ? values.tags : undefined,
      },
      {
        onSuccess: () => {
          navigate(ROUTES.ITEMS)
        },
      }
    )
  }

  const handleCancel = () => {
    navigate(ROUTES.ITEMS)
  }

  return (
    <div>
      <PageHeader
        title="Nuevo Artículo"
        description="Registrar un nuevo artículo en el inventario"
        backUrl={ROUTES.ITEMS}
      />

      <div className="mb-6 flex items-center gap-3 rounded-lg border bg-muted/50 p-4 text-sm text-muted-foreground">
        <Info className="h-4 w-4 flex-shrink-0" />
        <span>
          Las fotos del artículo se pueden agregar después de crear el registro, desde la página de edición.
        </span>
      </div>

      <div className="rounded-lg border bg-card p-6">
        <ItemForm
          onSubmit={handleSubmit}
          onCancel={handleCancel}
          isLoading={createMutation.isPending}
          showCustomer={true}
        />
      </div>
    </div>
  )
}
