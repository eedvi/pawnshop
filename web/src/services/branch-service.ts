import { apiGet, apiPost, apiPut, apiDelete, apiGetPaginated } from '@/lib/api-client'
import { Branch, CreateBranchInput, UpdateBranchInput } from '@/types'

export interface BranchListParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  is_active?: boolean
}

export const branchService = {
  // List all branches with pagination
  list: async (params: BranchListParams = {}) => {
    return apiGetPaginated<Branch>('/branches', params)
  },

  // Get a branch by ID
  getById: async (id: number): Promise<Branch> => {
    return apiGet<Branch>(`/branches/${id}`)
  },

  // Get a branch by code
  getByCode: async (code: string): Promise<Branch> => {
    return apiGet<Branch>(`/branches/code/${code}`)
  },

  // Create a new branch
  create: async (input: CreateBranchInput): Promise<Branch> => {
    return apiPost<Branch>('/branches', input)
  },

  // Update a branch
  update: async (id: number, input: UpdateBranchInput): Promise<Branch> => {
    return apiPut<Branch>(`/branches/${id}`, input)
  },

  // Delete a branch
  delete: async (id: number): Promise<void> => {
    return apiDelete(`/branches/${id}`)
  },

  // Activate a branch
  activate: async (id: number): Promise<void> => {
    return apiPost(`/branches/${id}/activate`)
  },

  // Deactivate a branch
  deactivate: async (id: number): Promise<void> => {
    return apiPost(`/branches/${id}/deactivate`)
  },
}
