import { Link } from 'react-router-dom'
import {
  AreaChart,
  Area,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
import {
  Banknote,
  CreditCard,
  Users,
  Package,
  AlertTriangle,
  TrendingUp,
  ArrowRight,
  DollarSign,
  Clock,
} from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { StatCard } from '@/components/common/stat-card'
import { StatusBadge } from '@/components/common/status-badge'
import { LoadingSpinner } from '@/components/common/loading-spinner'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import { useDashboardStats, useRecentLoans, useRecentPayments } from '@/hooks/use-dashboard'
import { useBranchStore } from '@/stores/branch-store'
import { formatCurrency, formatDate, formatRelativeTime } from '@/lib/format'
import { cn } from '@/lib/utils'

// Sample chart data - in production, this would come from an API
const loanChartData = [
  { month: 'Ene', prestamos: 45, monto: 125000 },
  { month: 'Feb', prestamos: 52, monto: 145000 },
  { month: 'Mar', prestamos: 48, monto: 132000 },
  { month: 'Abr', prestamos: 61, monto: 168000 },
  { month: 'May', prestamos: 55, monto: 152000 },
  { month: 'Jun', prestamos: 67, monto: 185000 },
]

const paymentChartData = [
  { month: 'Ene', pagos: 78, monto: 98000 },
  { month: 'Feb', pagos: 82, monto: 105000 },
  { month: 'Mar', pagos: 91, monto: 118000 },
  { month: 'Abr', pagos: 85, monto: 108000 },
  { month: 'May', pagos: 96, monto: 125000 },
  { month: 'Jun', pagos: 102, monto: 135000 },
]

export default function DashboardPage() {
  const { data: stats, isLoading: statsLoading, error: statsError } = useDashboardStats()
  const { data: recentLoans, isLoading: loansLoading } = useRecentLoans(5)
  const { data: recentPayments, isLoading: paymentsLoading } = useRecentPayments(5)
  const selectedBranch = useBranchStore((state) => state.selectedBranch)

  // Error state
  if (statsError) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <div className="text-center">
          <AlertTriangle className="mx-auto h-12 w-12 text-destructive" />
          <h3 className="mt-4 text-lg font-semibold">Error al cargar el dashboard</h3>
          <p className="mt-2 text-muted-foreground">No se pudieron obtener las estadísticas.</p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Dashboard"
        description={
          selectedBranch
            ? `Resumen de ${selectedBranch.name}`
            : 'Resumen general del sistema'
        }
      />

      {/* Alert for overdue loans */}
      {stats && stats.overdue_loans > 0 && (
        <div className="flex items-center gap-3 rounded-lg border border-destructive/50 bg-destructive/10 p-4">
          <AlertTriangle className="h-5 w-5 text-destructive" />
          <div className="flex-1">
            <p className="font-medium text-destructive">
              {stats.overdue_loans} préstamo{stats.overdue_loans > 1 ? 's' : ''} vencido{stats.overdue_loans > 1 ? 's' : ''}
            </p>
            <p className="text-sm text-muted-foreground">
              Monto total: {formatCurrency(stats.overdue_loans_amount)}
            </p>
          </div>
          <Button variant="outline" size="sm" asChild>
            <Link to="/loans?status=overdue">
              Ver préstamos
              <ArrowRight className="ml-2 h-4 w-4" />
            </Link>
          </Button>
        </div>
      )}

      {/* Primary Stats */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {statsLoading ? (
          <>
            {[...Array(4)].map((_, i) => (
              <div key={i} className="rounded-lg border bg-card p-6">
                <Skeleton className="h-4 w-24" />
                <Skeleton className="mt-2 h-8 w-32" />
                <Skeleton className="mt-1 h-4 w-20" />
              </div>
            ))}
          </>
        ) : (
          <>
            <StatCard
              title="Préstamos Activos"
              value={stats?.active_loans ?? 0}
              description={formatCurrency(stats?.active_loans_amount)}
              icon={Banknote}
              trend={
                stats?.loans_trend !== undefined
                  ? { value: stats.loans_trend, isPositive: stats.loans_trend >= 0 }
                  : undefined
              }
            />
            <StatCard
              title="Ingresos del Mes"
              value={formatCurrency(stats?.payments_this_month_amount)}
              description={`${stats?.payments_this_month ?? 0} pagos recibidos`}
              icon={DollarSign}
              trend={
                stats?.payments_trend !== undefined
                  ? { value: stats.payments_trend, isPositive: stats.payments_trend >= 0 }
                  : undefined
              }
            />
            <StatCard
              title="Clientes"
              value={stats?.total_customers ?? 0}
              description={`${stats?.new_customers_this_month ?? 0} nuevos este mes`}
              icon={Users}
              trend={
                stats?.customers_trend !== undefined
                  ? { value: stats.customers_trend, isPositive: stats.customers_trend >= 0 }
                  : undefined
              }
            />
            <StatCard
              title="Artículos"
              value={(stats?.items_available ?? 0) + (stats?.items_for_sale ?? 0) + (stats?.items_pawned ?? 0)}
              description={`${stats?.items_for_sale ?? 0} en venta`}
              icon={Package}
            />
          </>
        )}
      </div>

      {/* Secondary Stats */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {statsLoading ? (
          <>
            {[...Array(4)].map((_, i) => (
              <div key={i} className="rounded-lg border bg-card p-4">
                <Skeleton className="h-4 w-20" />
                <Skeleton className="mt-2 h-6 w-16" />
              </div>
            ))}
          </>
        ) : (
          <>
            <div className="rounded-lg border bg-card p-4">
              <p className="text-sm text-muted-foreground">Pagos Hoy</p>
              <p className="mt-1 text-xl font-semibold">{formatCurrency(stats?.payments_today_amount)}</p>
              <p className="text-xs text-muted-foreground">{stats?.payments_today ?? 0} transacciones</p>
            </div>
            <div className="rounded-lg border bg-card p-4">
              <p className="text-sm text-muted-foreground">Ventas Hoy</p>
              <p className="mt-1 text-xl font-semibold">{formatCurrency(stats?.sales_today_amount)}</p>
              <p className="text-xs text-muted-foreground">{stats?.sales_today ?? 0} ventas</p>
            </div>
            <div className="rounded-lg border bg-card p-4">
              <p className="text-sm text-muted-foreground">Vencen Hoy</p>
              <p className={cn(
                'mt-1 text-xl font-semibold',
                (stats?.loans_due_today ?? 0) > 0 && 'text-orange-600 dark:text-orange-400'
              )}>
                {stats?.loans_due_today ?? 0}
              </p>
              <p className="text-xs text-muted-foreground">préstamos</p>
            </div>
            <div className="rounded-lg border bg-card p-4">
              <p className="text-sm text-muted-foreground">Vencen Esta Semana</p>
              <p className={cn(
                'mt-1 text-xl font-semibold',
                (stats?.loans_due_this_week ?? 0) > 5 && 'text-orange-600 dark:text-orange-400'
              )}>
                {stats?.loans_due_this_week ?? 0}
              </p>
              <p className="text-xs text-muted-foreground">préstamos</p>
            </div>
          </>
        )}
      </div>

      {/* Charts */}
      <div className="grid gap-6 lg:grid-cols-2">
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold">Préstamos por Mes</h3>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </div>
          <div className="mt-4 h-64">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={loanChartData}>
                <defs>
                  <linearGradient id="colorPrestamos" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="hsl(var(--primary))" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="hsl(var(--primary))" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
                <XAxis dataKey="month" className="text-xs" tick={{ fill: 'hsl(var(--muted-foreground))' }} />
                <YAxis className="text-xs" tick={{ fill: 'hsl(var(--muted-foreground))' }} />
                <Tooltip
                  contentStyle={{
                    backgroundColor: 'hsl(var(--card))',
                    borderColor: 'hsl(var(--border))',
                    borderRadius: '8px',
                  }}
                  labelStyle={{ color: 'hsl(var(--foreground))' }}
                />
                <Area
                  type="monotone"
                  dataKey="prestamos"
                  stroke="hsl(var(--primary))"
                  fillOpacity={1}
                  fill="url(#colorPrestamos)"
                  name="Préstamos"
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold">Pagos Recibidos</h3>
            <CreditCard className="h-4 w-4 text-muted-foreground" />
          </div>
          <div className="mt-4 h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={paymentChartData}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
                <XAxis dataKey="month" className="text-xs" tick={{ fill: 'hsl(var(--muted-foreground))' }} />
                <YAxis className="text-xs" tick={{ fill: 'hsl(var(--muted-foreground))' }} />
                <Tooltip
                  contentStyle={{
                    backgroundColor: 'hsl(var(--card))',
                    borderColor: 'hsl(var(--border))',
                    borderRadius: '8px',
                  }}
                  labelStyle={{ color: 'hsl(var(--foreground))' }}
                  formatter={(value: number) => [formatCurrency(value), 'Monto']}
                />
                <Bar dataKey="monto" fill="hsl(var(--primary))" radius={[4, 4, 0, 0]} name="Monto" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      {/* Recent Activity */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Recent Loans */}
        <div className="rounded-lg border bg-card">
          <div className="flex items-center justify-between p-6 pb-4">
            <h3 className="text-lg font-semibold">Préstamos Recientes</h3>
            <Button variant="ghost" size="sm" asChild>
              <Link to="/loans">
                Ver todos
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
          <div className="px-6 pb-6">
            {loansLoading ? (
              <div className="space-y-4">
                {[...Array(5)].map((_, i) => (
                  <div key={i} className="flex items-center gap-4">
                    <Skeleton className="h-10 w-10 rounded-full" />
                    <div className="flex-1">
                      <Skeleton className="h-4 w-32" />
                      <Skeleton className="mt-1 h-3 w-24" />
                    </div>
                    <Skeleton className="h-5 w-20" />
                  </div>
                ))}
              </div>
            ) : recentLoans && recentLoans.length > 0 ? (
              <div className="space-y-4">
                {recentLoans.map((loan) => (
                  <Link
                    key={loan.loan_id}
                    to={`/loans/${loan.loan_id}`}
                    className="flex items-center gap-4 rounded-lg p-2 transition-colors hover:bg-muted/50"
                  >
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10">
                      <Banknote className="h-5 w-5 text-primary" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="font-medium truncate">{loan.customer_name}</p>
                      <p className="text-sm text-muted-foreground truncate">
                        {loan.loan_number} - {loan.item_name}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-medium">{formatCurrency(loan.loan_amount)}</p>
                      <StatusBadge status={loan.status} size="sm" />
                    </div>
                  </Link>
                ))}
              </div>
            ) : (
              <div className="flex h-32 items-center justify-center text-muted-foreground">
                No hay préstamos recientes
              </div>
            )}
          </div>
        </div>

        {/* Recent Payments */}
        <div className="rounded-lg border bg-card">
          <div className="flex items-center justify-between p-6 pb-4">
            <h3 className="text-lg font-semibold">Pagos Recientes</h3>
            <Button variant="ghost" size="sm" asChild>
              <Link to="/payments">
                Ver todos
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
          <div className="px-6 pb-6">
            {paymentsLoading ? (
              <div className="space-y-4">
                {[...Array(5)].map((_, i) => (
                  <div key={i} className="flex items-center gap-4">
                    <Skeleton className="h-10 w-10 rounded-full" />
                    <div className="flex-1">
                      <Skeleton className="h-4 w-32" />
                      <Skeleton className="mt-1 h-3 w-24" />
                    </div>
                    <Skeleton className="h-5 w-20" />
                  </div>
                ))}
              </div>
            ) : recentPayments && recentPayments.length > 0 ? (
              <div className="space-y-4">
                {recentPayments.map((payment) => (
                  <Link
                    key={payment.payment_id}
                    to={`/payments/${payment.payment_id}`}
                    className="flex items-center gap-4 rounded-lg p-2 transition-colors hover:bg-muted/50"
                  >
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-green-100 dark:bg-green-900/30">
                      <CreditCard className="h-5 w-5 text-green-600 dark:text-green-400" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="font-medium truncate">{payment.customer_name}</p>
                      <p className="text-sm text-muted-foreground truncate">
                        {payment.payment_number} - {payment.loan_number}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="font-medium text-green-600 dark:text-green-400">
                        +{formatCurrency(payment.amount)}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        {formatRelativeTime(payment.payment_date)}
                      </p>
                    </div>
                  </Link>
                ))}
              </div>
            ) : (
              <div className="flex h-32 items-center justify-center text-muted-foreground">
                No hay pagos recientes
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="rounded-lg border bg-card p-6">
        <h3 className="text-lg font-semibold">Acciones Rápidas</h3>
        <div className="mt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <Button asChild variant="outline" className="h-auto py-4">
            <Link to="/loans/new" className="flex flex-col items-center gap-2">
              <Banknote className="h-6 w-6" />
              <span>Nuevo Préstamo</span>
            </Link>
          </Button>
          <Button asChild variant="outline" className="h-auto py-4">
            <Link to="/payments/new" className="flex flex-col items-center gap-2">
              <CreditCard className="h-6 w-6" />
              <span>Registrar Pago</span>
            </Link>
          </Button>
          <Button asChild variant="outline" className="h-auto py-4">
            <Link to="/customers/new" className="flex flex-col items-center gap-2">
              <Users className="h-6 w-6" />
              <span>Nuevo Cliente</span>
            </Link>
          </Button>
          <Button asChild variant="outline" className="h-auto py-4">
            <Link to="/items/new" className="flex flex-col items-center gap-2">
              <Package className="h-6 w-6" />
              <span>Nuevo Artículo</span>
            </Link>
          </Button>
        </div>
      </div>
    </div>
  )
}
