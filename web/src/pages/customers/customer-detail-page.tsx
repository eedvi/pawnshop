import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import {
  Loader2,
  Pencil,
  Ban,
  CheckCircle,
  Phone,
  Mail,
  MapPin,
  User,
  Briefcase,
  AlertCircle,
  CreditCard,
  Award,
  Calendar,
} from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import { StatusBadge } from '@/components/common/status-badge'
import { ROUTES, customerEditRoute } from '@/routes/routes'
import { useCustomer, useBlockCustomer, useUnblockCustomer } from '@/hooks/use-customers'
import { useConfirm } from '@/hooks'
import { formatCurrency, formatPhone, formatDate } from '@/lib/format'
import { IDENTITY_TYPES, GENDERS, LOYALTY_TIERS } from '@/types'
import { BlockCustomerDialog } from './block-customer-dialog'
import { CustomerLoansTab } from './tabs/customer-loans-tab'
import { CustomerPaymentsTab } from './tabs/customer-payments-tab'
import { CustomerItemsTab } from './tabs/customer-items-tab'

export default function CustomerDetailPage() {
  const { id } = useParams()
  const customerId = parseInt(id!, 10)

  const { data: customer, isLoading } = useCustomer(customerId)
  const blockMutation = useBlockCustomer()
  const unblockMutation = useUnblockCustomer()

  const confirmUnblock = useConfirm()
  const [blockDialogOpen, setBlockDialogOpen] = useState(false)

  const handleBlock = () => {
    setBlockDialogOpen(true)
  }

  const handleBlockConfirm = (reason: string) => {
    blockMutation.mutate(
      { id: customerId, reason },
      {
        onSuccess: () => {
          setBlockDialogOpen(false)
        },
      }
    )
  }

  const handleUnblock = async () => {
    const confirmed = await confirmUnblock.confirm({
      title: 'Desbloquear Cliente',
      description: `¿Estás seguro de desbloquear a "${customer?.first_name} ${customer?.last_name}"?`,
      confirmLabel: 'Desbloquear',
    })

    if (confirmed) {
      unblockMutation.mutate(customerId)
    }
  }

  if (isLoading) {
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

  const identityTypeLabel = IDENTITY_TYPES.find((t) => t.value === customer.identity_type)?.label
  const genderLabel = customer.gender
    ? GENDERS.find((g) => g.value === customer.gender)?.label
    : null
  const loyaltyTier = LOYALTY_TIERS.find((t) => t.value === customer.loyalty_tier)

  return (
    <div className="space-y-6">
      <PageHeader
        title={`${customer.first_name} ${customer.last_name}`}
        description={`${identityTypeLabel}: ${customer.identity_number}`}
        backUrl={ROUTES.CUSTOMERS}
        actions={
          <div className="flex gap-2">
            {customer.is_blocked ? (
              <Button variant="outline" onClick={handleUnblock}>
                <CheckCircle className="mr-2 h-4 w-4" />
                Desbloquear
              </Button>
            ) : (
              <Button variant="outline" onClick={handleBlock}>
                <Ban className="mr-2 h-4 w-4" />
                Bloquear
              </Button>
            )}
            <Button asChild>
              <Link to={customerEditRoute(customerId)}>
                <Pencil className="mr-2 h-4 w-4" />
                Editar
              </Link>
            </Button>
          </div>
        }
      />

      {/* Status Alert */}
      {customer.is_blocked && (
        <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-4">
          <div className="flex items-center gap-2 text-destructive">
            <AlertCircle className="h-5 w-5" />
            <span className="font-medium">Cliente bloqueado</span>
          </div>
          {customer.blocked_reason && (
            <p className="mt-1 text-sm text-muted-foreground">
              Razón: {customer.blocked_reason}
            </p>
          )}
        </div>
      )}

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Estado</CardTitle>
            <User className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <StatusBadge
              status={customer.is_blocked ? 'blocked' : customer.is_active ? 'active' : 'inactive'}
            />
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Límite de Crédito</CardTitle>
            <CreditCard className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(customer.credit_limit)}</div>
            <p className="text-xs text-muted-foreground">
              Score: {customer.credit_score}/100
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Préstamos</CardTitle>
            <Briefcase className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{customer.total_loans}</div>
            <p className="text-xs text-muted-foreground">
              Pagados: {formatCurrency(customer.total_paid)}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Fidelidad</CardTitle>
            <Award className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-2">
              <Badge variant="outline">{loyaltyTier?.label}</Badge>
              <span className="text-sm text-muted-foreground">
                {customer.loyalty_points} pts
              </span>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Tabs */}
      <Tabs defaultValue="info" className="space-y-4">
        <TabsList>
          <TabsTrigger value="info">Información</TabsTrigger>
          <TabsTrigger value="loans">Préstamos</TabsTrigger>
          <TabsTrigger value="payments">Pagos</TabsTrigger>
          <TabsTrigger value="items">Artículos</TabsTrigger>
        </TabsList>

        <TabsContent value="info" className="space-y-4">
          <div className="grid gap-4 md:grid-cols-2">
            {/* Contact Info */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Información de Contacto</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center gap-3">
                  <Phone className="h-4 w-4 text-muted-foreground" />
                  <div>
                    <p className="text-sm font-medium">{formatPhone(customer.phone)}</p>
                    {customer.phone_secondary && (
                      <p className="text-xs text-muted-foreground">
                        Secundario: {formatPhone(customer.phone_secondary)}
                      </p>
                    )}
                  </div>
                </div>

                {customer.email && (
                  <div className="flex items-center gap-3">
                    <Mail className="h-4 w-4 text-muted-foreground" />
                    <p className="text-sm">{customer.email}</p>
                  </div>
                )}

                {customer.address && (
                  <div className="flex items-start gap-3">
                    <MapPin className="mt-0.5 h-4 w-4 text-muted-foreground" />
                    <div>
                      <p className="text-sm">{customer.address}</p>
                      {(customer.city || customer.state || customer.postal_code) && (
                        <p className="text-xs text-muted-foreground">
                          {[customer.city, customer.state, customer.postal_code]
                            .filter(Boolean)
                            .join(', ')}
                        </p>
                      )}
                    </div>
                  </div>
                )}

                {customer.birth_date && (
                  <div className="flex items-center gap-3">
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                    <div>
                      <p className="text-sm">{formatDate(customer.birth_date)}</p>
                      {genderLabel && (
                        <p className="text-xs text-muted-foreground">{genderLabel}</p>
                      )}
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Emergency Contact */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Contacto de Emergencia</CardTitle>
              </CardHeader>
              <CardContent>
                {customer.emergency_contact_name ? (
                  <div className="space-y-2">
                    <p className="font-medium">{customer.emergency_contact_name}</p>
                    {customer.emergency_contact_relation && (
                      <p className="text-sm text-muted-foreground">
                        {customer.emergency_contact_relation}
                      </p>
                    )}
                    {customer.emergency_contact_phone && (
                      <p className="text-sm">{formatPhone(customer.emergency_contact_phone)}</p>
                    )}
                  </div>
                ) : (
                  <p className="text-sm text-muted-foreground">Sin contacto de emergencia</p>
                )}
              </CardContent>
            </Card>

            {/* Work Info */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Información Laboral</CardTitle>
              </CardHeader>
              <CardContent>
                {customer.occupation || customer.workplace || customer.monthly_income ? (
                  <div className="space-y-2">
                    {customer.occupation && (
                      <p className="font-medium">{customer.occupation}</p>
                    )}
                    {customer.workplace && (
                      <p className="text-sm text-muted-foreground">{customer.workplace}</p>
                    )}
                    {customer.monthly_income !== undefined && customer.monthly_income !== null && (
                      <p className="text-sm">
                        Ingreso mensual: {formatCurrency(customer.monthly_income)}
                      </p>
                    )}
                  </div>
                ) : (
                  <p className="text-sm text-muted-foreground">Sin información laboral</p>
                )}
              </CardContent>
            </Card>

            {/* Notes */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Notas</CardTitle>
              </CardHeader>
              <CardContent>
                {customer.notes ? (
                  <p className="whitespace-pre-wrap text-sm">{customer.notes}</p>
                ) : (
                  <p className="text-sm text-muted-foreground">Sin notas</p>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Audit Info */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Información del Registro</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex gap-8 text-sm text-muted-foreground">
                <div>
                  <span className="font-medium">Creado:</span> {formatDate(customer.created_at)}
                </div>
                <div>
                  <span className="font-medium">Actualizado:</span> {formatDate(customer.updated_at)}
                </div>
                {customer.branch && (
                  <div>
                    <span className="font-medium">Sucursal:</span> {customer.branch.name}
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="loans">
          <CustomerLoansTab customerId={customerId} />
        </TabsContent>

        <TabsContent value="payments">
          <CustomerPaymentsTab customerId={customerId} />
        </TabsContent>

        <TabsContent value="items">
          <CustomerItemsTab customerId={customerId} />
        </TabsContent>
      </Tabs>

      <ConfirmDialog {...confirmUnblock} />
      <BlockCustomerDialog
        open={blockDialogOpen}
        onOpenChange={setBlockDialogOpen}
        customer={customer}
        onConfirm={handleBlockConfirm}
        isLoading={blockMutation.isPending}
      />
    </div>
  )
}
