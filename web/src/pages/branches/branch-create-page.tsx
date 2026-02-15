import { useNavigate } from 'react-router-dom'

import { PageHeader } from '@/components/layout/page-header'
import { ROUTES } from '@/routes/routes'
import { useCreateBranch } from '@/hooks/use-branches'
import { BranchForm } from './branch-form'
import { BranchFormValues } from './schemas'

export default function BranchCreatePage() {
  const navigate = useNavigate()
  const createMutation = useCreateBranch()

  const handleSubmit = async (values: BranchFormValues) => {
    try {
      await createMutation.mutateAsync({
        name: values.name,
        code: values.code,
        address: values.address || undefined,
        phone: values.phone || undefined,
        email: values.email || undefined,
        timezone: values.timezone,
        currency: values.currency,
        default_interest_rate: values.default_interest_rate,
        default_loan_term_days: values.default_loan_term_days,
        default_grace_period: values.default_grace_period,
      })
      navigate(ROUTES.BRANCHES)
    } catch {
      // Error handling is done in the mutation
    }
  }

  const handleCancel = () => {
    navigate(ROUTES.BRANCHES)
  }

  return (
    <div>
      <PageHeader
        title="Nueva Sucursal"
        description="Registrar una nueva sucursal en el sistema"
        backUrl={ROUTES.BRANCHES}
      />

      <div className="rounded-lg border bg-card p-6">
        <BranchForm
          onSubmit={handleSubmit}
          onCancel={handleCancel}
          isLoading={createMutation.isPending}
        />
      </div>
    </div>
  )
}
