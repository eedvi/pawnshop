import { apiGet, apiPost, apiPut, apiDelete } from '@/lib/api-client'
import type {
  Role,
  CreateRoleInput,
  UpdateRoleInput,
} from '@/types'

export const roleService = {
  list: () =>
    apiGet<Role[]>('/roles'),

  getById: (id: number) =>
    apiGet<Role>(`/roles/${id}`),

  create: (input: CreateRoleInput) =>
    apiPost<Role>('/roles', input),

  update: (id: number, input: UpdateRoleInput) =>
    apiPut<Role>(`/roles/${id}`, input),

  delete: (id: number) =>
    apiDelete<void>(`/roles/${id}`),

  getPermissions: () =>
    apiGet<string[]>('/roles/permissions'),
}
