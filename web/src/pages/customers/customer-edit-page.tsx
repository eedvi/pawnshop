import { useNavigate, useParams } from 'react-router-dom'
import { Loader2 } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { customerRoute } from '@/routes/routes'
import { useCustomer, useUpdateCustomer } from '@/hooks/use-customers'
import { CustomerForm } from './customer-form'
import { CustomerFormValues } from './schemas'

export default function CustomerEditPage() {
  const { id } = useParams()
  const customerId = parseInt(id!, 10)
  const navigate = useNavigate()

  const { data: customer, isLoading: isLoadingCustomer } = useCustomer(customerId)
  const updateMutation = useUpdateCustomer()

  const handleSubmit = (values: CustomerFormValues) => {
    updateMutation.mutate(
      {
        id: customerId,
        input: {
          first_name: values.first_name,
          last_name: values.last_name,
          birth_date: values.birth_date || undefined,
          gender: values.gender || undefined,
          phone: values.phone,
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
          credit_limit: values.credit_limit,
          notes: values.notes || undefined,
          is_active: values.is_active,
        },
      },
      {
        onSuccess: () => {
          navigate(customerRoute(customerId))
        },
      }
    )
  }

  const handleCancel = () => {
    navigate(customerRoute(customerId))
  }

  if (isLoadingCustomer) {
    return (
      <div className="flex h-96 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (!customer) {
    return (
      <div className="flex h-96 items-center justify-center">
        <p className="text-muted-foreground">Cliente no encontrado</p>
      </div>
    )
  }

  return (
    <div>
      <PageHeader
        title={`Editar: ${customer.first_name} ${customer.last_name}`}
        description="Modificar información del cliente"
        backUrl={customerRoute(customerId)}
      />

      <div className="rounded-lg border bg-card p-6">
        <CustomerForm
          customer={customer}
          onSubmit={handleSubmit}
          onCancel={handleCancel}
          isLoading={updateMutation.isPending}
        />
      </div>
    </div>
  )
}
