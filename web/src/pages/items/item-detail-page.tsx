import { useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import {
  Loader2,
  Pencil,
  ShoppingCart,
  Tag,
  Package,
  Barcode,
  Calendar,
  DollarSign,
  Scale,
  User,
  Folder,
  Image as ImageIcon,
  Trash2,
  Plus,
} from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import { ROUTES, itemEditRoute, customerRoute } from '@/routes/routes'
import { useItem, useMarkItemForSale, useDeleteItemPhoto, useUploadItemPhotos } from '@/hooks/use-items'
import { useConfirm } from '@/hooks'
import { formatCurrency, formatDate } from '@/lib/format'
import { ITEM_STATUSES, ITEM_CONDITIONS } from '@/types'
import { MarkForSaleDialog } from './mark-for-sale-dialog'

export default function ItemDetailPage() {
  const { id } = useParams()
  const itemId = parseInt(id!, 10)

  const { data: item, isLoading } = useItem(itemId)
  const markForSaleMutation = useMarkItemForSale()
  const deletePhotoMutation = useDeleteItemPhoto()
  const uploadPhotosMutation = useUploadItemPhotos()

  const confirmDeletePhoto = useConfirm()
  const [markForSaleOpen, setMarkForSaleOpen] = useState(false)

  const handleMarkForSale = () => {
    setMarkForSaleOpen(true)
  }

  const handleMarkForSaleConfirm = (salePrice: number) => {
    markForSaleMutation.mutate(
      { id: itemId, salePrice },
      {
        onSuccess: () => {
          setMarkForSaleOpen(false)
        },
      }
    )
  }

  const handleDeletePhoto = async (photoUrl: string) => {
    const confirmed = await confirmDeletePhoto.confirm({
      title: 'Eliminar Foto',
      description: '¿Estás seguro de eliminar esta foto?',
      confirmLabel: 'Eliminar',
      variant: 'destructive',
    })

    if (confirmed) {
      deletePhotoMutation.mutate({ id: itemId, photoUrl })
    }
  }

  const handleUploadPhotos = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files
    if (files && files.length > 0) {
      uploadPhotosMutation.mutate({ id: itemId, files: Array.from(files) })
    }
  }

  if (isLoading) {
    return (
      <div className="flex h-96 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (!item) {
    return (
      <div className="flex h-96 items-center justify-center">
        <p className="text-muted-foreground">Artículo no encontrado</p>
      </div>
    )
  }

  const status = ITEM_STATUSES.find((s) => s.value === item.status)
  const condition = ITEM_CONDITIONS.find((c) => c.value === item.condition)
  const canMarkForSale = ['available', 'confiscated'].includes(item.status)

  const colorMap: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
    green: 'default',
    blue: 'default',
    purple: 'secondary',
    orange: 'secondary',
    gray: 'outline',
    red: 'destructive',
    cyan: 'secondary',
    yellow: 'secondary',
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title={item.name}
        description={`SKU: ${item.sku}`}
        backUrl={ROUTES.ITEMS}
        actions={
          <div className="flex gap-2">
            {canMarkForSale && (
              <Button variant="outline" onClick={handleMarkForSale}>
                <ShoppingCart className="mr-2 h-4 w-4" />
                Marcar para Venta
              </Button>
            )}
            <Button asChild>
              <Link to={itemEditRoute(itemId)}>
                <Pencil className="mr-2 h-4 w-4" />
                Editar
              </Link>
            </Button>
          </div>
        }
      />

      {/* Status and Valuation Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Estado</CardTitle>
            <Tag className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <Badge variant={colorMap[status?.color || 'gray'] || 'outline'}>
              {status?.label || item.status}
            </Badge>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Avalúo</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(item.appraised_value)}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Préstamo</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(item.loan_value)}</div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Precio Venta</CardTitle>
            <ShoppingCart className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {item.sale_price ? formatCurrency(item.sale_price) : '-'}
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        {/* Item Details */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Package className="h-5 w-5" />
              Detalles del Artículo
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm text-muted-foreground">Condición</p>
                <p className="font-medium">{condition?.label || item.condition}</p>
              </div>
              {item.brand && (
                <div>
                  <p className="text-sm text-muted-foreground">Marca</p>
                  <p className="font-medium">{item.brand}</p>
                </div>
              )}
              {item.model && (
                <div>
                  <p className="text-sm text-muted-foreground">Modelo</p>
                  <p className="font-medium">{item.model}</p>
                </div>
              )}
              {item.serial_number && (
                <div>
                  <p className="text-sm text-muted-foreground">No. Serie</p>
                  <p className="font-mono">{item.serial_number}</p>
                </div>
              )}
              {item.color && (
                <div>
                  <p className="text-sm text-muted-foreground">Color</p>
                  <p className="font-medium">{item.color}</p>
                </div>
              )}
              {item.weight && (
                <div>
                  <p className="text-sm text-muted-foreground">Peso</p>
                  <p className="font-medium">{item.weight} g</p>
                </div>
              )}
              {item.purity && (
                <div>
                  <p className="text-sm text-muted-foreground">Pureza/Quilates</p>
                  <p className="font-medium">{item.purity}</p>
                </div>
              )}
            </div>

            {item.description && (
              <div>
                <p className="text-sm text-muted-foreground">Descripción</p>
                <p className="whitespace-pre-wrap">{item.description}</p>
              </div>
            )}

            {item.notes && (
              <div>
                <p className="text-sm text-muted-foreground">Notas</p>
                <p className="whitespace-pre-wrap">{item.notes}</p>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Related Info */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Barcode className="h-5 w-5" />
              Información Relacionada
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center gap-3">
              <Folder className="h-4 w-4 text-muted-foreground" />
              <div>
                <p className="text-sm text-muted-foreground">Categoría</p>
                <p className="font-medium">{item.category?.name || 'Sin categoría'}</p>
              </div>
            </div>

            {item.customer && (
              <div className="flex items-center gap-3">
                <User className="h-4 w-4 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Cliente</p>
                  <Link
                    to={customerRoute(item.customer.id)}
                    className="font-medium text-primary hover:underline"
                  >
                    {item.customer.first_name} {item.customer.last_name}
                  </Link>
                </div>
              </div>
            )}

            <div className="flex items-center gap-3">
              <Calendar className="h-4 w-4 text-muted-foreground" />
              <div>
                <p className="text-sm text-muted-foreground">Fecha de Adquisición</p>
                <p className="font-medium">{formatDate(item.acquisition_date)}</p>
              </div>
            </div>

            <div className="flex items-center gap-3">
              <Scale className="h-4 w-4 text-muted-foreground" />
              <div>
                <p className="text-sm text-muted-foreground">Tipo de Adquisición</p>
                <p className="font-medium">
                  {item.acquisition_type === 'pawn' && 'Empeño'}
                  {item.acquisition_type === 'purchase' && 'Compra'}
                  {item.acquisition_type === 'confiscation' && 'Confiscación'}
                </p>
              </div>
            </div>

            {item.acquisition_price !== undefined && item.acquisition_price !== null && (
              <div className="flex items-center gap-3">
                <DollarSign className="h-4 w-4 text-muted-foreground" />
                <div>
                  <p className="text-sm text-muted-foreground">Precio de Adquisición</p>
                  <p className="font-medium">{formatCurrency(item.acquisition_price)}</p>
                </div>
              </div>
            )}

            {item.branch && (
              <div className="mt-4 pt-4 border-t">
                <p className="text-sm text-muted-foreground">Sucursal</p>
                <p className="font-medium">{item.branch.name}</p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Photo Gallery */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle className="flex items-center gap-2">
            <ImageIcon className="h-5 w-5" />
            Galería de Fotos
          </CardTitle>
          <label className="cursor-pointer">
            <input
              type="file"
              multiple
              accept="image/*"
              className="hidden"
              onChange={handleUploadPhotos}
              disabled={uploadPhotosMutation.isPending}
            />
            <Button variant="outline" size="sm" asChild>
              <span>
                {uploadPhotosMutation.isPending ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <Plus className="mr-2 h-4 w-4" />
                )}
                Agregar Fotos
              </span>
            </Button>
          </label>
        </CardHeader>
        <CardContent>
          {item.photos && item.photos.length > 0 ? (
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
              {item.photos.map((photo, index) => (
                <div key={index} className="relative group aspect-square">
                  <img
                    src={photo}
                    alt={`${item.name} - Foto ${index + 1}`}
                    className="w-full h-full object-cover rounded-lg"
                  />
                  <button
                    onClick={() => handleDeletePhoto(photo)}
                    className="absolute top-2 right-2 p-1 bg-destructive text-destructive-foreground rounded-full opacity-0 group-hover:opacity-100 transition-opacity"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-center text-sm text-muted-foreground py-8">
              No hay fotos del artículo
            </p>
          )}
        </CardContent>
      </Card>

      {/* Audit Info */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Información del Registro</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex gap-8 text-sm text-muted-foreground">
            <div>
              <span className="font-medium">Creado:</span> {formatDate(item.created_at)}
            </div>
            <div>
              <span className="font-medium">Actualizado:</span> {formatDate(item.updated_at)}
            </div>
          </div>
        </CardContent>
      </Card>

      <ConfirmDialog {...confirmDeletePhoto} />
      <MarkForSaleDialog
        open={markForSaleOpen}
        onOpenChange={setMarkForSaleOpen}
        item={item}
        onConfirm={handleMarkForSaleConfirm}
        isLoading={markForSaleMutation.isPending}
      />
    </div>
  )
}
