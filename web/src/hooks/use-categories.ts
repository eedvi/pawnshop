import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { categoryService, CategoryListParams } from '@/services/category-service'
import { CreateCategoryInput, UpdateCategoryInput } from '@/types'
import { toast } from 'sonner'

// Query keys
export const categoryKeys = {
  all: ['categories'] as const,
  lists: () => [...categoryKeys.all, 'list'] as const,
  list: (params: CategoryListParams) => [...categoryKeys.lists(), params] as const,
  tree: () => [...categoryKeys.all, 'tree'] as const,
  details: () => [...categoryKeys.all, 'detail'] as const,
  detail: (id: number) => [...categoryKeys.details(), id] as const,
}

// Hook to list categories (flat)
export function useCategories(params: CategoryListParams = {}) {
  return useQuery({
    queryKey: categoryKeys.list(params),
    queryFn: () => categoryService.list(params),
  })
}

// Hook to list categories as tree
export function useCategoryTree() {
  return useQuery({
    queryKey: categoryKeys.tree(),
    queryFn: () => categoryService.listTree(),
  })
}

// Hook to get a single category
export function useCategory(id: number) {
  return useQuery({
    queryKey: categoryKeys.detail(id),
    queryFn: () => categoryService.getById(id),
    enabled: id > 0,
  })
}

// Hook to create a category
export function useCreateCategory() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateCategoryInput) => categoryService.create(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: categoryKeys.all })
      toast.success('Categoría creada exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al crear la categoría')
    },
  })
}

// Hook to update a category
export function useUpdateCategory() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: UpdateCategoryInput }) =>
      categoryService.update(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: categoryKeys.all })
      queryClient.invalidateQueries({ queryKey: categoryKeys.detail(id) })
      toast.success('Categoría actualizada exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al actualizar la categoría')
    },
  })
}

// Hook to delete a category
export function useDeleteCategory() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => categoryService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: categoryKeys.all })
      toast.success('Categoría eliminada exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al eliminar la categoría')
    },
  })
}
