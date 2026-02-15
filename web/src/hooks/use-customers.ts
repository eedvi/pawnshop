import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { customerService, CustomerListParams } from '@/services/customer-service'
import { CreateCustomerInput, UpdateCustomerInput } from '@/types'
import { toast } from 'sonner'

// Query keys
export const customerKeys = {
  all: ['customers'] as const,
  lists: () => [...customerKeys.all, 'list'] as const,
  list: (params: CustomerListParams) => [...customerKeys.lists(), params] as const,
  details: () => [...customerKeys.all, 'detail'] as const,
  detail: (id: number) => [...customerKeys.details(), id] as const,
}

// Hook to list customers
export function useCustomers(params: CustomerListParams = {}) {
  return useQuery({
    queryKey: customerKeys.list(params),
    queryFn: () => customerService.list(params),
  })
}

// Hook to get a single customer
export function useCustomer(id: number) {
  return useQuery({
    queryKey: customerKeys.detail(id),
    queryFn: () => customerService.getById(id),
    enabled: id > 0,
  })
}

// Hook to create a customer
export function useCreateCustomer() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateCustomerInput) => customerService.create(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: customerKeys.lists() })
      toast.success('Cliente creado exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al crear el cliente')
    },
  })
}

// Hook to update a customer
export function useUpdateCustomer() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: UpdateCustomerInput }) =>
      customerService.update(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: customerKeys.lists() })
      queryClient.invalidateQueries({ queryKey: customerKeys.detail(id) })
      toast.success('Cliente actualizado exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al actualizar el cliente')
    },
  })
}

// Hook to delete a customer
export function useDeleteCustomer() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => customerService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: customerKeys.lists() })
      toast.success('Cliente eliminado exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al eliminar el cliente')
    },
  })
}

// Hook to block a customer
export function useBlockCustomer() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, reason }: { id: number; reason: string }) =>
      customerService.block(id, reason),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: customerKeys.lists() })
      queryClient.invalidateQueries({ queryKey: customerKeys.detail(id) })
      toast.success('Cliente bloqueado exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al bloquear el cliente')
    },
  })
}

// Hook to unblock a customer
export function useUnblockCustomer() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => customerService.unblock(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: customerKeys.lists() })
      queryClient.invalidateQueries({ queryKey: customerKeys.detail(id) })
      toast.success('Cliente desbloqueado exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al desbloquear el cliente')
    },
  })
}
