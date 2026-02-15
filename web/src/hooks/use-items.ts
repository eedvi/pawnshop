import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { itemService } from '@/services/item-service'
import { CreateItemInput, UpdateItemInput, ItemListParams, ItemStatus } from '@/types'
import { toast } from 'sonner'

// Query keys
export const itemKeys = {
  all: ['items'] as const,
  lists: () => [...itemKeys.all, 'list'] as const,
  list: (params: ItemListParams) => [...itemKeys.lists(), params] as const,
  details: () => [...itemKeys.all, 'detail'] as const,
  detail: (id: number) => [...itemKeys.details(), id] as const,
}

// Hook to list items
export function useItems(params: ItemListParams = {}) {
  return useQuery({
    queryKey: itemKeys.list(params),
    queryFn: () => itemService.list(params),
  })
}

// Hook to get a single item
export function useItem(id: number) {
  return useQuery({
    queryKey: itemKeys.detail(id),
    queryFn: () => itemService.getById(id),
    enabled: id > 0,
  })
}

// Hook to create an item
export function useCreateItem() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateItemInput) => itemService.create(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: itemKeys.lists() })
      toast.success('Artículo creado exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al crear el artículo')
    },
  })
}

// Hook to update an item
export function useUpdateItem() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: UpdateItemInput }) =>
      itemService.update(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: itemKeys.lists() })
      queryClient.invalidateQueries({ queryKey: itemKeys.detail(id) })
      toast.success('Artículo actualizado exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al actualizar el artículo')
    },
  })
}

// Hook to delete an item
export function useDeleteItem() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => itemService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: itemKeys.lists() })
      toast.success('Artículo eliminado exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al eliminar el artículo')
    },
  })
}

// Hook to update item status
export function useUpdateItemStatus() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, status, notes }: { id: number; status: ItemStatus; notes?: string }) =>
      itemService.updateStatus(id, status, notes),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: itemKeys.lists() })
      queryClient.invalidateQueries({ queryKey: itemKeys.detail(id) })
      toast.success('Estado actualizado exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al actualizar el estado')
    },
  })
}

// Hook to mark item for sale
export function useMarkItemForSale() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, salePrice }: { id: number; salePrice: number }) =>
      itemService.markForSale(id, salePrice),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: itemKeys.lists() })
      queryClient.invalidateQueries({ queryKey: itemKeys.detail(id) })
      toast.success('Artículo marcado para venta')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al marcar para venta')
    },
  })
}

// Hook to upload photos
export function useUploadItemPhotos() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, files }: { id: number; files: File[] }) =>
      itemService.uploadPhotos(id, files),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: itemKeys.detail(id) })
      toast.success('Fotos subidas exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al subir las fotos')
    },
  })
}

// Hook to delete photo
export function useDeleteItemPhoto() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, photoUrl }: { id: number; photoUrl: string }) =>
      itemService.deletePhoto(id, photoUrl),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: itemKeys.detail(id) })
      toast.success('Foto eliminada exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al eliminar la foto')
    },
  })
}
