import { apiGet, apiPost, apiPut, apiDelete } from '@/lib/api-client'
import { Category, CreateCategoryInput, UpdateCategoryInput } from '@/types'

export interface CategoryListParams {
  parent_id?: number
  is_active?: boolean
}

export const categoryService = {
  // List all categories (flat)
  list: async (params: CategoryListParams = {}): Promise<Category[]> => {
    return apiGet<Category[]>('/categories', params)
  },

  // List categories as tree (with children)
  listTree: async (): Promise<Category[]> => {
    return apiGet<Category[]>('/categories/tree')
  },

  // Get a category by ID
  getById: async (id: number): Promise<Category> => {
    return apiGet<Category>(`/categories/${id}`)
  },

  // Get a category by slug
  getBySlug: async (slug: string): Promise<Category> => {
    return apiGet<Category>(`/categories/slug/${slug}`)
  },

  // Create a new category
  create: async (input: CreateCategoryInput): Promise<Category> => {
    return apiPost<Category>('/categories', input)
  },

  // Update a category
  update: async (id: number, input: UpdateCategoryInput): Promise<Category> => {
    return apiPut<Category>(`/categories/${id}`, input)
  },

  // Delete a category
  delete: async (id: number): Promise<void> => {
    return apiDelete(`/categories/${id}`)
  },
}
