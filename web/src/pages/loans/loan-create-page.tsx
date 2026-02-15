import { PageHeader } from '@/components/layout/page-header'
import { ROUTES } from '@/routes/routes'
import { LoanWizard } from './loan-wizard'

export default function LoanCreatePage() {
  return (
    <div>
      <PageHeader
        title="Nuevo Préstamo"
        description="Crear un nuevo préstamo"
        backUrl={ROUTES.LOANS}
      />

      <LoanWizard />
    </div>
  )
}
