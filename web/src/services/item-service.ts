import { apiGet, apiPost, apiPut, apiDelete, apiGetPaginated, apiUpload } from '@/lib/api-client'
import { Item, CreateItemInput, UpdateItemInput, ItemListParams, ItemStatus } from '@/types'

export const itemService = {
  // List items with pagination
  list: async (params: ItemListParams = {}) => {
    return apiGetPaginated<Item>('/items', params)
  },

  // Get an item by ID
  getById: async (id: number): Promise<Item> => {
    return apiGet<Item>(`/items/${id}`)
  },

  // Create a new item
  create: async (input: CreateItemInput): Promise<Item> => {
    return apiPost<Item>('/items', input)
  },

  // Update an item
  update: async (id: number, input: UpdateItemInput): Promise<Item> => {
    return apiPut<Item>(`/items/${id}`, input)
  },

  // Delete an item
  delete: async (id: number): Promise<void> => {
    return apiDelete(`/items/${id}`)
  },

  // Update item status
  updateStatus: async (id: number, status: ItemStatus, notes?: string): Promise<void> => {
    return apiPost(`/items/${id}/status`, { status, notes })
  },

  // Mark item for sale
  markForSale: async (id: number, salePrice: number): Promise<void> => {
    return apiPost(`/items/${id}/mark-for-sale`, { sale_price: salePrice })
  },

  // Upload photos (uploads one at a time to match backend)
  uploadPhotos: async (id: number, files: File[]): Promise<string[]> => {
    const results: string[] = []
    for (const file of files) {
      const formData = new FormData()
      formData.append('image', file)
      const result = await apiUpload<{ url: string }>(`/items/${id}/images`, formData)
      if (result?.url) {
        results.push(result.url)
      }
    }
    return results
  },

  // Delete photo
  deletePhoto: async (id: number, photoUrl: string): Promise<void> => {
    return apiDelete(`/items/${id}/images?url=${encodeURIComponent(photoUrl)}`)
  },

  // Mark item as delivered to customer
  markAsDelivered: async (id: number, notes?: string): Promise<Item> => {
    const response = await apiPost<{ message: string; item: Item }>(
      `/items/${id}/mark-as-delivered`,
      { notes }
    )
    return response.item
  },
}
