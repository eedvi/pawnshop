import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'

import { PageHeader } from '@/components/layout/page-header'
import { ROUTES } from '@/routes/routes'
import { useCreateCustomer } from '@/hooks/use-customers'
import { useBranchStore } from '@/stores/branch-store'
import { CustomerForm } from './customer-form'
import { CustomerFormValues } from './schemas'

export default function CustomerCreatePage() {
  const navigate = useNavigate()
  const { selectedBranchId } = useBranchStore()
  const createMutation = useCreateCustomer()

  const handleSubmit = (values: CustomerFormValues) => {
    if (!selectedBranchId) {
      toast.error('Debe seleccionar una sucursal primero')
      return
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const { is_active, ...createData } = values

    createMutation.mutate(
      {
        ...createData,
        branch_id: selectedBranchId,
        birth_date: values.birth_date || undefined,
        gender: values.gender || undefined,
        phone_secondary: values.phone_secondary || undefined,
        email: values.email || undefined,
        address: values.address || undefined,
        city: values.city || undefined,
        state: values.state || undefined,
        postal_code: values.postal_code || undefined,
        emergency_contact_name: values.emergency_contact_name || undefined,
        emergency_contact_phone: values.emergency_contact_phone || undefined,
        emergency_contact_relation: values.emergency_contact_relation || undefined,
        occupation: values.occupation || undefined,
        workplace: values.workplace || undefined,
        monthly_income: values.monthly_income ?? undefined,
        notes: values.notes || undefined,
      },
      {
        onSuccess: () => {
          toast.success('Cliente creado exitosamente')
          navigate(ROUTES.CUSTOMERS)
        },
        onError: (error) => {
          toast.error(error.message || 'Error al crear el cliente')
        },
      }
    )
  }

  const handleCancel = () => {
    navigate(ROUTES.CUSTOMERS)
  }

  return (
    <div>
      <PageHeader
        title="Nuevo Cliente"
        description="Registrar un nuevo cliente en el sistema"
        backUrl={ROUTES.CUSTOMERS}
      />

      <div className="rounded-lg border bg-card p-6">
        <CustomerForm
          onSubmit={handleSubmit}
          onCancel={handleCancel}
          isLoading={createMutation.isPending}
        />
      </div>
    </div>
  )
}
