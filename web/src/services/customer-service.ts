import { apiGet, apiPost, apiPut, apiDelete, apiGetPaginated } from '@/lib/api-client'
import { Customer, CreateCustomerInput, UpdateCustomerInput } from '@/types'

export interface CustomerListParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  search?: string
  branch_id?: number
  is_active?: boolean
  is_blocked?: boolean
}

export const customerService = {
  // List customers with pagination
  list: async (params: CustomerListParams = {}) => {
    return apiGetPaginated<Customer>('/customers', params)
  },

  // Get a customer by ID
  getById: async (id: number): Promise<Customer> => {
    return apiGet<Customer>(`/customers/${id}`)
  },

  // Create a new customer
  create: async (input: CreateCustomerInput): Promise<Customer> => {
    return apiPost<Customer>('/customers', input)
  },

  // Update a customer
  update: async (id: number, input: UpdateCustomerInput): Promise<Customer> => {
    return apiPut<Customer>(`/customers/${id}`, input)
  },

  // Delete a customer
  delete: async (id: number): Promise<void> => {
    return apiDelete(`/customers/${id}`)
  },

  // Block a customer
  block: async (id: number, reason: string): Promise<void> => {
    return apiPost(`/customers/${id}/block`, { reason })
  },

  // Unblock a customer
  unblock: async (id: number): Promise<void> => {
    return apiPost(`/customers/${id}/unblock`)
  },
}
