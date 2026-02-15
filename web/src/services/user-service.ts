import { apiGet, apiGetPaginated, apiPost, apiPut, apiDelete } from '@/lib/api-client'
import type {
  User,
  CreateUserInput,
  UpdateUserInput,
  UserListParams,
} from '@/types'

export const userService = {
  list: (params?: UserListParams) =>
    apiGetPaginated<User>('/users', params),

  getById: (id: number) =>
    apiGet<User>(`/users/${id}`),

  create: (input: CreateUserInput) =>
    apiPost<User>('/users', input),

  update: (id: number, input: UpdateUserInput) =>
    apiPut<User>(`/users/${id}`, input),

  delete: (id: number) =>
    apiDelete<void>(`/users/${id}`),

  resetPassword: (id: number, newPassword: string) =>
    apiPost<void>(`/users/${id}/reset-password`, { password: newPassword }),

  toggleActive: (id: number) =>
    apiPost<User>(`/users/${id}/toggle-active`),

  unlock: (id: number) =>
    apiPost<User>(`/users/${id}/unlock`),
}
