import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { branchService, BranchListParams } from '@/services/branch-service'
import { CreateBranchInput, UpdateBranchInput } from '@/types'
import { toast } from 'sonner'

// Query keys
export const branchKeys = {
  all: ['branches'] as const,
  lists: () => [...branchKeys.all, 'list'] as const,
  list: (params: BranchListParams) => [...branchKeys.lists(), params] as const,
  details: () => [...branchKeys.all, 'detail'] as const,
  detail: (id: number) => [...branchKeys.details(), id] as const,
}

// Hook to list branches
export function useBranches(params: BranchListParams = {}) {
  return useQuery({
    queryKey: branchKeys.list(params),
    queryFn: () => branchService.list(params),
  })
}

// Hook to get a single branch
export function useBranch(id: number) {
  return useQuery({
    queryKey: branchKeys.detail(id),
    queryFn: () => branchService.getById(id),
    enabled: id > 0,
  })
}

// Hook to create a branch
export function useCreateBranch() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateBranchInput) => branchService.create(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: branchKeys.lists() })
      toast.success('Sucursal creada exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al crear la sucursal')
    },
  })
}

// Hook to update a branch
export function useUpdateBranch() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: UpdateBranchInput }) =>
      branchService.update(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: branchKeys.lists() })
      queryClient.invalidateQueries({ queryKey: branchKeys.detail(id) })
      toast.success('Sucursal actualizada exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al actualizar la sucursal')
    },
  })
}

// Hook to delete a branch
export function useDeleteBranch() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => branchService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: branchKeys.lists() })
      toast.success('Sucursal eliminada exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al eliminar la sucursal')
    },
  })
}

// Hook to activate a branch
export function useActivateBranch() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => branchService.activate(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: branchKeys.lists() })
      queryClient.invalidateQueries({ queryKey: branchKeys.detail(id) })
      toast.success('Sucursal activada exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al activar la sucursal')
    },
  })
}

// Hook to deactivate a branch
export function useDeactivateBranch() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => branchService.deactivate(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: branchKeys.lists() })
      queryClient.invalidateQueries({ queryKey: branchKeys.detail(id) })
      toast.success('Sucursal desactivada exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al desactivar la sucursal')
    },
  })
}
