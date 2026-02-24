import { Link } from 'react-router-dom'
import { ExternalLink } from 'lucide-react'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { itemRoute } from '@/routes/routes'
import { formatCurrency, formatDate } from '@/lib/format'
import { useItems } from '@/hooks/use-items'

interface CustomerItemsTabProps {
  customerId: number
}

const ITEM_STATUS_LABELS: Record<string, { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline' }> = {
  in_custody: { label: 'En custodia', variant: 'default' },
  available_for_sale: { label: 'En venta', variant: 'secondary' },
  sold: { label: 'Vendido', variant: 'outline' },
  returned: { label: 'Devuelto', variant: 'outline' },
  transferred: { label: 'Transferido', variant: 'outline' },
  lost: { label: 'Perdido', variant: 'destructive' },
  damaged: { label: 'Dañado', variant: 'destructive' },
}

export function CustomerItemsTab({ customerId }: CustomerItemsTabProps) {
  const { data: itemsResponse, isLoading } = useItems({ customer_id: customerId, per_page: 100 })
  const items = itemsResponse?.data || []

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Artículos</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {Array.from({ length: 3 }).map((_, i) => (
              <Skeleton key={i} className="h-12 w-full" />
            ))}
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Artículos del Cliente</CardTitle>
      </CardHeader>
      <CardContent>
        {items.length === 0 ? (
          <p className="text-center text-sm text-muted-foreground py-8">
            Este cliente no tiene artículos registrados
          </p>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Artículo</TableHead>
                <TableHead>Categoría</TableHead>
                <TableHead>Estado</TableHead>
                <TableHead className="text-right">Avalúo</TableHead>
                <TableHead className="text-right">Préstamo</TableHead>
                <TableHead>Fecha</TableHead>
                <TableHead></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {items.map((item) => {
                const statusConfig = ITEM_STATUS_LABELS[item.status] || { label: item.status, variant: 'outline' as const }
                return (
                  <TableRow key={item.id}>
                    <TableCell className="font-medium">{item.name}</TableCell>
                    <TableCell>{item.category?.name || '-'}</TableCell>
                    <TableCell>
                      <Badge variant={statusConfig.variant}>{statusConfig.label}</Badge>
                    </TableCell>
                    <TableCell className="text-right">{formatCurrency(item.appraised_value)}</TableCell>
                    <TableCell className="text-right">{formatCurrency(item.loan_value)}</TableCell>
                    <TableCell>{formatDate(item.created_at)}</TableCell>
                    <TableCell>
                      <Link
                        to={itemRoute(item.id)}
                        className="inline-flex items-center text-primary hover:underline"
                      >
                        <ExternalLink className="h-4 w-4" />
                      </Link>
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>
        )}
      </CardContent>
    </Card>
  )
}
