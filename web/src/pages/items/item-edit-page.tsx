import { useNavigate, useParams } from 'react-router-dom'
import { Loader2, Plus, Trash2, Image as ImageIcon } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ConfirmDialog } from '@/components/common/confirm-dialog'
import { itemRoute } from '@/routes/routes'
import { useItem, useUpdateItem, useUploadItemPhotos, useDeleteItemPhoto } from '@/hooks/use-items'
import { useConfirm } from '@/hooks'
import { ItemForm } from './item-form'
import { ItemFormValues } from './schemas'

export default function ItemEditPage() {
  const { id } = useParams()
  const itemId = parseInt(id!, 10)
  const navigate = useNavigate()

  const { data: item, isLoading: isLoadingItem } = useItem(itemId)
  const updateMutation = useUpdateItem()
  const uploadPhotosMutation = useUploadItemPhotos()
  const deletePhotoMutation = useDeleteItemPhoto()
  const confirmDeletePhoto = useConfirm()

  const handleUploadPhotos = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files
    if (files && files.length > 0) {
      uploadPhotosMutation.mutate({ id: itemId, files: Array.from(files) })
    }
    // Reset input to allow re-uploading same file
    event.target.value = ''
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

  const handleSubmit = (values: ItemFormValues) => {
    updateMutation.mutate(
      {
        id: itemId,
        input: {
          name: values.name,
          description: values.description || undefined,
          brand: values.brand || undefined,
          model: values.model || undefined,
          serial_number: values.serial_number || undefined,
          color: values.color || undefined,
          condition: values.condition,
          category_id: values.category_id || undefined,
          appraised_value: values.appraised_value,
          loan_value: values.loan_value,
          sale_price: values.sale_price ?? undefined,
          weight: values.weight ?? undefined,
          purity: values.purity || undefined,
          notes: values.notes || undefined,
          tags: values.tags?.length ? values.tags : undefined,
        },
      },
      {
        onSuccess: () => {
          navigate(itemRoute(itemId))
        },
      }
    )
  }

  const handleCancel = () => {
    navigate(itemRoute(itemId))
  }

  if (isLoadingItem) {
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

  return (
    <div>
      <PageHeader
        title={`Editar: ${item.name}`}
        description="Modificar información del artículo"
        backUrl={itemRoute(itemId)}
      />

      <div className="rounded-lg border bg-card p-6">
        <ItemForm
          item={item}
          onSubmit={handleSubmit}
          onCancel={handleCancel}
          isLoading={updateMutation.isPending}
        />
      </div>

      {/* Photo Gallery */}
      <Card className="mt-6">
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
                    type="button"
                    onClick={() => handleDeletePhoto(photo)}
                    disabled={deletePhotoMutation.isPending}
                    className="absolute top-2 right-2 p-1 bg-destructive text-destructive-foreground rounded-full opacity-0 group-hover:opacity-100 transition-opacity"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-center text-sm text-muted-foreground py-8">
              No hay fotos del artículo. Haz clic en "Agregar Fotos" para subir imágenes.
            </p>
          )}
        </CardContent>
      </Card>

      <ConfirmDialog {...confirmDeletePhoto} />
    </div>
  )
}
