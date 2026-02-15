import { useState } from 'react'
import { useParams, useNavigate, useSearchParams } from 'react-router-dom'
import { Pencil, Power, PowerOff, Trash2, Building2, Phone, Mail, MapPin, Clock, Percent, Calendar, AlertTriangle } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { StatusBadge } from '@/components/common/status-badge'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import { ROUTES } from '@/routes/routes'
import { useBranch, useUpdateBranch, useActivateBranch, useDeactivateBranch, useDeleteBranch } from '@/hooks/use-branches'
import { useConfirm } from '@/hooks'
import { formatPercent, formatDate } from '@/lib/format'
import { BranchForm } from './branch-form'
import { BranchFormValues } from './schemas'

export default function BranchDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()
  const isEditing = searchParams.get('edit') === 'true'

  const branchId = parseInt(id ?? '0', 10)
  const { data: branch, isLoading, error } = useBranch(branchId)

  const updateMutation = useUpdateBranch()
  const activateMutation = useActivateBranch()
  const deactivateMutation = useDeactivateBranch()
  const deleteMutation = useDeleteBranch()

  const confirmDelete = useConfirm()
  const confirmDeactivate = useConfirm()

  const handleEdit = () => {
    setSearchParams({ edit: 'true' })
  }

  const handleCancelEdit = () => {
    setSearchParams({})
  }

  const handleSubmit = async (values: BranchFormValues) => {
    try {
      await updateMutation.mutateAsync({
        id: branchId,
        input: {
          name: values.name,
          address: values.address || undefined,
          phone: values.phone || undefined,
          email: values.email || undefined,
          timezone: values.timezone,
          currency: values.currency,
          is_active: values.is_active,
          default_interest_rate: values.default_interest_rate,
          default_loan_term_days: values.default_loan_term_days,
          default_grace_period: values.default_grace_period,
        },
      })
      setSearchParams({})
    } catch {
      // Error handling is done in the mutation
    }
  }

  const handleActivate = () => {
    activateMutation.mutate(branchId)
  }

  const handleDeactivate = async () => {
    const confirmed = await confirmDeactivate.confirm({
      title: 'Desactivar Sucursal',
      description: `¿Estás seguro de desactivar "${branch?.name}"? Las sucursales inactivas no pueden realizar operaciones.`,
      confirmLabel: 'Desactivar',
      variant: 'destructive',
    })

    if (confirmed) {
      deactivateMutation.mutate(branchId)
    }
  }

  const handleDelete = async () => {
    const confirmed = await confirmDelete.confirm({
      title: 'Eliminar Sucursal',
      description: `¿Estás seguro de eliminar "${branch?.name}"? Esta acción no se puede deshacer.`,
      confirmLabel: 'Eliminar',
      variant: 'destructive',
    })

    if (confirmed) {
      await deleteMutation.mutateAsync(branchId)
      navigate(ROUTES.BRANCHES)
    }
  }

  if (isLoading) {
    return (
      <div>
        <PageHeader
          title={<Skeleton className="h-8 w-48" />}
          backUrl={ROUTES.BRANCHES}
        />
        <div className="rounded-lg border bg-card p-6">
          <div className="space-y-4">
            <Skeleton className="h-6 w-32" />
            <Skeleton className="h-6 w-64" />
            <Skeleton className="h-6 w-48" />
          </div>
        </div>
      </div>
    )
  }

  if (error || !branch) {
    return (
      <div>
        <PageHeader title="Error" backUrl={ROUTES.BRANCHES} />
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center gap-2 text-destructive">
            <AlertTriangle className="h-5 w-5" />
            <p>No se pudo cargar la sucursal.</p>
          </div>
        </div>
      </div>
    )
  }

  if (isEditing) {
    return (
      <div>
        <PageHeader
          title={`Editar: ${branch.name}`}
          description="Modificar información de la sucursal"
          backUrl={ROUTES.BRANCHES}
        />
        <div className="rounded-lg border bg-card p-6">
          <BranchForm
            branch={branch}
            onSubmit={handleSubmit}
            onCancel={handleCancelEdit}
            isLoading={updateMutation.isPending}
          />
        </div>
      </div>
    )
  }

  return (
    <div>
      <PageHeader
        title={branch.name}
        description={`Código: ${branch.code}`}
        backUrl={ROUTES.BRANCHES}
        actions={
          <div className="flex gap-2">
            <Button variant="outline" onClick={handleEdit}>
              <Pencil className="mr-2 h-4 w-4" />
              Editar
            </Button>
            {branch.is_active ? (
              <Button variant="outline" onClick={handleDeactivate}>
                <PowerOff className="mr-2 h-4 w-4" />
                Desactivar
              </Button>
            ) : (
              <Button variant="outline" onClick={handleActivate}>
                <Power className="mr-2 h-4 w-4" />
                Activar
              </Button>
            )}
            <Button variant="destructive" onClick={handleDelete}>
              <Trash2 className="mr-2 h-4 w-4" />
              Eliminar
            </Button>
          </div>
        }
      />

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Basic Info */}
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">Información General</h3>
          <dl className="space-y-4">
            <div className="flex items-start gap-3">
              <Building2 className="h-5 w-5 text-muted-foreground mt-0.5" />
              <div>
                <dt className="text-sm text-muted-foreground">Nombre</dt>
                <dd className="font-medium">{branch.name}</dd>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <div className="h-5 w-5 flex items-center justify-center text-muted-foreground">
                <span className="text-xs font-mono">ID</span>
              </div>
              <div>
                <dt className="text-sm text-muted-foreground">Código</dt>
                <dd className="font-mono font-medium">{branch.code}</dd>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <MapPin className="h-5 w-5 text-muted-foreground mt-0.5" />
              <div>
                <dt className="text-sm text-muted-foreground">Dirección</dt>
                <dd className="font-medium">{branch.address || '-'}</dd>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <Phone className="h-5 w-5 text-muted-foreground mt-0.5" />
              <div>
                <dt className="text-sm text-muted-foreground">Teléfono</dt>
                <dd className="font-medium">{branch.phone || '-'}</dd>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <Mail className="h-5 w-5 text-muted-foreground mt-0.5" />
              <div>
                <dt className="text-sm text-muted-foreground">Email</dt>
                <dd className="font-medium">{branch.email || '-'}</dd>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="h-5 w-5" />
              <div>
                <dt className="text-sm text-muted-foreground">Estado</dt>
                <dd>
                  <StatusBadge status={branch.is_active ? 'active' : 'inactive'} />
                </dd>
              </div>
            </div>
          </dl>
        </div>

        {/* Loan Settings */}
        <div className="rounded-lg border bg-card p-6">
          <h3 className="text-lg font-semibold mb-4">Configuración de Préstamos</h3>
          <dl className="space-y-4">
            <div className="flex items-start gap-3">
              <Percent className="h-5 w-5 text-muted-foreground mt-0.5" />
              <div>
                <dt className="text-sm text-muted-foreground">Tasa de Interés</dt>
                <dd className="font-medium">{formatPercent(branch.default_interest_rate)}</dd>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <Calendar className="h-5 w-5 text-muted-foreground mt-0.5" />
              <div>
                <dt className="text-sm text-muted-foreground">Plazo Predeterminado</dt>
                <dd className="font-medium">{branch.default_loan_term_days} días</dd>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <Clock className="h-5 w-5 text-muted-foreground mt-0.5" />
              <div>
                <dt className="text-sm text-muted-foreground">Período de Gracia</dt>
                <dd className="font-medium">{branch.default_grace_period} días</dd>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <div className="h-5 w-5 flex items-center justify-center text-muted-foreground">
                <span className="text-xs">$</span>
              </div>
              <div>
                <dt className="text-sm text-muted-foreground">Moneda</dt>
                <dd className="font-medium">{branch.currency}</dd>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <Clock className="h-5 w-5 text-muted-foreground mt-0.5" />
              <div>
                <dt className="text-sm text-muted-foreground">Zona Horaria</dt>
                <dd className="font-medium">{branch.timezone}</dd>
              </div>
            </div>
          </dl>
        </div>

        {/* Metadata */}
        <div className="rounded-lg border bg-card p-6 lg:col-span-2">
          <h3 className="text-lg font-semibold mb-4">Información del Sistema</h3>
          <dl className="grid gap-4 sm:grid-cols-2">
            <div>
              <dt className="text-sm text-muted-foreground">Creada</dt>
              <dd className="font-medium">{formatDate(branch.created_at)}</dd>
            </div>
            <div>
              <dt className="text-sm text-muted-foreground">Última Actualización</dt>
              <dd className="font-medium">{formatDate(branch.updated_at)}</dd>
            </div>
          </dl>
        </div>
      </div>

      <ConfirmDialog {...confirmDelete} />
      <ConfirmDialog {...confirmDeactivate} />
    </div>
  )
}
