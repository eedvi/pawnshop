import { useState } from 'react'
import { format, subDays, startOfMonth, endOfMonth } from 'date-fns'
import {
  FileText,
  CreditCard,
  ShoppingBag,
  AlertTriangle,
  Download,
  Loader2,
  Calendar,
} from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  useLoanReport,
  usePaymentReport,
  useSalesReport,
  useOverdueReport,
  useExportLoanReport,
  useExportPaymentReport,
  useExportSalesReport,
  useExportOverdueReport,
} from '@/hooks/use-reports'
import { useBranchStore } from '@/stores/branch-store'
import { formatCurrency, formatDate } from '@/lib/format'
import type { ReportFilters } from '@/types'

export default function ReportsPage() {
  const { selectedBranchId } = useBranchStore()
  const today = new Date()

  const [dateFrom, setDateFrom] = useState(format(startOfMonth(today), 'yyyy-MM-dd'))
  const [dateTo, setDateTo] = useState(format(endOfMonth(today), 'yyyy-MM-dd'))

  const filters: ReportFilters = {
    branch_id: selectedBranchId ?? undefined,
    date_from: dateFrom,
    date_to: dateTo,
  }

  // Queries
  const { data: loanReport, isLoading: loadingLoans } = useLoanReport(filters)
  const { data: paymentReport, isLoading: loadingPayments } = usePaymentReport(filters)
  const { data: salesReport, isLoading: loadingSales } = useSalesReport(filters)
  const { data: overdueReport, isLoading: loadingOverdue } = useOverdueReport(filters)

  // Export mutations
  const exportLoansMutation = useExportLoanReport()
  const exportPaymentsMutation = useExportPaymentReport()
  const exportSalesMutation = useExportSalesReport()
  const exportOverdueMutation = useExportOverdueReport()

  const setThisMonth = () => {
    setDateFrom(format(startOfMonth(today), 'yyyy-MM-dd'))
    setDateTo(format(endOfMonth(today), 'yyyy-MM-dd'))
  }

  const setLast7Days = () => {
    setDateFrom(format(subDays(today, 7), 'yyyy-MM-dd'))
    setDateTo(format(today, 'yyyy-MM-dd'))
  }

  const setLast30Days = () => {
    setDateFrom(format(subDays(today, 30), 'yyyy-MM-dd'))
    setDateTo(format(today, 'yyyy-MM-dd'))
  }

  return (
    <div>
      <PageHeader
        title="Reportes"
        description="Reportes y estadísticas del sistema"
      />

      {/* Date Range Filter */}
      <Card className="mb-6">
        <CardHeader className="pb-3">
          <CardTitle className="text-base flex items-center gap-2">
            <Calendar className="h-4 w-4" />
            Rango de Fechas
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap items-end gap-4">
            <div className="space-y-1">
              <Label htmlFor="date-from">Desde</Label>
              <Input
                id="date-from"
                type="date"
                value={dateFrom}
                onChange={(e) => setDateFrom(e.target.value)}
                className="w-40"
              />
            </div>
            <div className="space-y-1">
              <Label htmlFor="date-to">Hasta</Label>
              <Input
                id="date-to"
                type="date"
                value={dateTo}
                onChange={(e) => setDateTo(e.target.value)}
                className="w-40"
              />
            </div>
            <div className="flex gap-2">
              <Button variant="outline" size="sm" onClick={setLast7Days}>
                Últimos 7 días
              </Button>
              <Button variant="outline" size="sm" onClick={setLast30Days}>
                Últimos 30 días
              </Button>
              <Button variant="outline" size="sm" onClick={setThisMonth}>
                Este mes
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <Tabs defaultValue="loans" className="space-y-4">
        <TabsList>
          <TabsTrigger value="loans" className="flex items-center gap-2">
            <FileText className="h-4 w-4" />
            Préstamos
          </TabsTrigger>
          <TabsTrigger value="payments" className="flex items-center gap-2">
            <CreditCard className="h-4 w-4" />
            Pagos
          </TabsTrigger>
          <TabsTrigger value="sales" className="flex items-center gap-2">
            <ShoppingBag className="h-4 w-4" />
            Ventas
          </TabsTrigger>
          <TabsTrigger value="overdue" className="flex items-center gap-2">
            <AlertTriangle className="h-4 w-4" />
            Vencidos
          </TabsTrigger>
        </TabsList>

        {/* Loans Report */}
        <TabsContent value="loans">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle>Reporte de Préstamos</CardTitle>
                <CardDescription>
                  Préstamos del período seleccionado
                </CardDescription>
              </div>
              <Button
                variant="outline"
                onClick={() => exportLoansMutation.mutate(filters)}
                disabled={exportLoansMutation.isPending}
              >
                {exportLoansMutation.isPending ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <Download className="mr-2 h-4 w-4" />
                )}
                Exportar PDF
              </Button>
            </CardHeader>
            <CardContent>
              {/* Summary */}
              {loanReport?.summary && (
                <div className="grid gap-4 md:grid-cols-4 mb-6">
                  <Card>
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold">{loanReport.summary.total_loans}</div>
                      <p className="text-xs text-muted-foreground">Total Préstamos</p>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold">
                        {formatCurrency(loanReport.summary.total_principal)}
                      </div>
                      <p className="text-xs text-muted-foreground">Capital Total</p>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold">
                        {formatCurrency(loanReport.summary.total_interest)}
                      </div>
                      <p className="text-xs text-muted-foreground">Interés Total</p>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold">
                        {formatCurrency(loanReport.summary.total_remaining)}
                      </div>
                      <p className="text-xs text-muted-foreground">Saldo Pendiente</p>
                    </CardContent>
                  </Card>
                </div>
              )}

              {loadingLoans ? (
                <div className="flex justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                </div>
              ) : !loanReport?.items?.length ? (
                <div className="text-center py-8 text-muted-foreground">
                  No hay datos para el período seleccionado
                </div>
              ) : (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>No. Préstamo</TableHead>
                      <TableHead>Cliente</TableHead>
                      <TableHead>Artículo</TableHead>
                      <TableHead>Monto</TableHead>
                      <TableHead>Saldo</TableHead>
                      <TableHead>Estado</TableHead>
                      <TableHead>Vencimiento</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {loanReport.items.map((item) => (
                      <TableRow key={item.loan_id}>
                        <TableCell className="font-mono">{item.loan_number}</TableCell>
                        <TableCell>{item.customer_name}</TableCell>
                        <TableCell>{item.item_name}</TableCell>
                        <TableCell>{formatCurrency(item.loan_amount)}</TableCell>
                        <TableCell>{formatCurrency(item.total_remaining)}</TableCell>
                        <TableCell>
                          <Badge variant={item.status === 'active' ? 'default' : 'secondary'}>
                            {item.status}
                          </Badge>
                        </TableCell>
                        <TableCell>{formatDate(item.due_date)}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Payments Report */}
        <TabsContent value="payments">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle>Reporte de Pagos</CardTitle>
                <CardDescription>Pagos del período seleccionado</CardDescription>
              </div>
              <Button
                variant="outline"
                onClick={() => exportPaymentsMutation.mutate(filters)}
                disabled={exportPaymentsMutation.isPending}
              >
                {exportPaymentsMutation.isPending ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <Download className="mr-2 h-4 w-4" />
                )}
                Exportar PDF
              </Button>
            </CardHeader>
            <CardContent>
              {/* Summary */}
              {paymentReport?.summary && (
                <div className="grid gap-4 md:grid-cols-4 mb-6">
                  <Card>
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold">{paymentReport.summary.total_payments}</div>
                      <p className="text-xs text-muted-foreground">Total Pagos</p>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold">
                        {formatCurrency(paymentReport.summary.total_amount)}
                      </div>
                      <p className="text-xs text-muted-foreground">Monto Total</p>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold">
                        {formatCurrency(paymentReport.summary.total_principal)}
                      </div>
                      <p className="text-xs text-muted-foreground">Capital Cobrado</p>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold">
                        {formatCurrency(paymentReport.summary.total_interest)}
                      </div>
                      <p className="text-xs text-muted-foreground">Interés Cobrado</p>
                    </CardContent>
                  </Card>
                </div>
              )}

              {loadingPayments ? (
                <div className="flex justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                </div>
              ) : !paymentReport?.items?.length ? (
                <div className="text-center py-8 text-muted-foreground">
                  No hay datos para el período seleccionado
                </div>
              ) : (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>No. Pago</TableHead>
                      <TableHead>Préstamo</TableHead>
                      <TableHead>Cliente</TableHead>
                      <TableHead>Monto</TableHead>
                      <TableHead>Método</TableHead>
                      <TableHead>Fecha</TableHead>
                      <TableHead>Estado</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {paymentReport.items.map((item) => (
                      <TableRow key={item.payment_id}>
                        <TableCell className="font-mono">{item.payment_number}</TableCell>
                        <TableCell className="font-mono">{item.loan_number}</TableCell>
                        <TableCell>{item.customer_name}</TableCell>
                        <TableCell>{formatCurrency(item.amount)}</TableCell>
                        <TableCell>{item.payment_method}</TableCell>
                        <TableCell>{formatDate(item.payment_date)}</TableCell>
                        <TableCell>
                          <Badge variant={item.status === 'completed' ? 'default' : 'secondary'}>
                            {item.status}
                          </Badge>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Sales Report */}
        <TabsContent value="sales">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle>Reporte de Ventas</CardTitle>
                <CardDescription>Ventas del período seleccionado</CardDescription>
              </div>
              <Button
                variant="outline"
                onClick={() => exportSalesMutation.mutate(filters)}
                disabled={exportSalesMutation.isPending}
              >
                {exportSalesMutation.isPending ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <Download className="mr-2 h-4 w-4" />
                )}
                Exportar PDF
              </Button>
            </CardHeader>
            <CardContent>
              {/* Summary */}
              {salesReport?.summary && (
                <div className="grid gap-4 md:grid-cols-4 mb-6">
                  <Card>
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold">{salesReport.summary.total_sales}</div>
                      <p className="text-xs text-muted-foreground">Total Ventas</p>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold">
                        {formatCurrency(salesReport.summary.gross_amount)}
                      </div>
                      <p className="text-xs text-muted-foreground">Monto Bruto</p>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold text-destructive">
                        -{formatCurrency(salesReport.summary.total_discounts)}
                      </div>
                      <p className="text-xs text-muted-foreground">Descuentos</p>
                    </CardContent>
                  </Card>
                  <Card>
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold text-green-600">
                        {formatCurrency(salesReport.summary.net_amount)}
                      </div>
                      <p className="text-xs text-muted-foreground">Monto Neto</p>
                    </CardContent>
                  </Card>
                </div>
              )}

              {loadingSales ? (
                <div className="flex justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                </div>
              ) : !salesReport?.items?.length ? (
                <div className="text-center py-8 text-muted-foreground">
                  No hay datos para el período seleccionado
                </div>
              ) : (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>No. Venta</TableHead>
                      <TableHead>Artículo</TableHead>
                      <TableHead>Cliente</TableHead>
                      <TableHead>Precio</TableHead>
                      <TableHead>Descuento</TableHead>
                      <TableHead>Total</TableHead>
                      <TableHead>Fecha</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {salesReport.items.map((item) => (
                      <TableRow key={item.sale_id}>
                        <TableCell className="font-mono">{item.sale_number}</TableCell>
                        <TableCell>{item.item_name}</TableCell>
                        <TableCell>{item.customer_name || '-'}</TableCell>
                        <TableCell>{formatCurrency(item.sale_price)}</TableCell>
                        <TableCell className="text-destructive">
                          {item.discount_amount > 0
                            ? `-${formatCurrency(item.discount_amount)}`
                            : '-'}
                        </TableCell>
                        <TableCell className="font-medium">
                          {formatCurrency(item.final_price)}
                        </TableCell>
                        <TableCell>{formatDate(item.sale_date)}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Overdue Report */}
        <TabsContent value="overdue">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle className="text-destructive flex items-center gap-2">
                  <AlertTriangle className="h-5 w-5" />
                  Préstamos Vencidos
                </CardTitle>
                <CardDescription>Préstamos con pagos pendientes</CardDescription>
              </div>
              <Button
                variant="outline"
                onClick={() => exportOverdueMutation.mutate(filters)}
                disabled={exportOverdueMutation.isPending}
              >
                {exportOverdueMutation.isPending ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <Download className="mr-2 h-4 w-4" />
                )}
                Exportar PDF
              </Button>
            </CardHeader>
            <CardContent>
              {/* Summary */}
              {overdueReport?.summary && (
                <div className="grid gap-4 md:grid-cols-3 mb-6">
                  <Card className="border-destructive">
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold text-destructive">
                        {overdueReport.summary.total_overdue_loans}
                      </div>
                      <p className="text-xs text-muted-foreground">Préstamos Vencidos</p>
                    </CardContent>
                  </Card>
                  <Card className="border-destructive">
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold text-destructive">
                        {formatCurrency(overdueReport.summary.total_overdue_amount)}
                      </div>
                      <p className="text-xs text-muted-foreground">Monto Vencido</p>
                    </CardContent>
                  </Card>
                  <Card className="border-destructive">
                    <CardContent className="pt-4">
                      <div className="text-2xl font-bold text-destructive">
                        {formatCurrency(overdueReport.summary.total_late_fees)}
                      </div>
                      <p className="text-xs text-muted-foreground">Mora Acumulada</p>
                    </CardContent>
                  </Card>
                </div>
              )}

              {loadingOverdue ? (
                <div className="flex justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                </div>
              ) : !overdueReport?.items?.length ? (
                <div className="text-center py-8 text-muted-foreground">
                  No hay préstamos vencidos
                </div>
              ) : (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>No. Préstamo</TableHead>
                      <TableHead>Cliente</TableHead>
                      <TableHead>Teléfono</TableHead>
                      <TableHead>Artículo</TableHead>
                      <TableHead>Saldo</TableHead>
                      <TableHead>Días Vencido</TableHead>
                      <TableHead>Mora</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {overdueReport.items.map((item) => (
                      <TableRow key={item.loan_id} className="bg-destructive/5">
                        <TableCell className="font-mono">{item.loan_number}</TableCell>
                        <TableCell>{item.customer_name}</TableCell>
                        <TableCell>{item.customer_phone}</TableCell>
                        <TableCell>{item.item_name}</TableCell>
                        <TableCell>{formatCurrency(item.total_remaining)}</TableCell>
                        <TableCell>
                          <Badge variant="destructive">{item.days_overdue} días</Badge>
                        </TableCell>
                        <TableCell className="text-destructive font-medium">
                          {formatCurrency(item.late_fee_amount)}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
