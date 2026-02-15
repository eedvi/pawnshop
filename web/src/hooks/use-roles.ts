import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { roleService } from '@/services/role-service'
import type { CreateRoleInput, UpdateRoleInput } from '@/types'

export const roleKeys = {
  all: ['roles'] as const,
  lists: () => [...roleKeys.all, 'list'] as const,
  details: () => [...roleKeys.all, 'detail'] as const,
  detail: (id: number) => [...roleKeys.details(), id] as const,
  permissions: () => [...roleKeys.all, 'permissions'] as const,
}

export function useRoles() {
  return useQuery({
    queryKey: roleKeys.lists(),
    queryFn: () => roleService.list(),
  })
}

export function useRole(id: number) {
  return useQuery({
    queryKey: roleKeys.detail(id),
    queryFn: () => roleService.getById(id),
    enabled: id > 0,
  })
}

export function usePermissions() {
  return useQuery({
    queryKey: roleKeys.permissions(),
    queryFn: () => roleService.getPermissions(),
  })
}

export function useCreateRole() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateRoleInput) => roleService.create(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: roleKeys.lists() })
    },
  })
}

export function useUpdateRole() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: UpdateRoleInput }) =>
      roleService.update(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: roleKeys.lists() })
      queryClient.invalidateQueries({ queryKey: roleKeys.detail(id) })
    },
  })
}

export function useDeleteRole() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => roleService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: roleKeys.lists() })
    },
  })
}
