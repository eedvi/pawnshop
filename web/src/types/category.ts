// Category types - mirrors internal/domain/category.go

export interface Category {
  id: number
  name: string
  slug: string
  description?: string
  parent_id?: number
  icon?: string
  default_interest_rate: number
  min_loan_amount?: number
  max_loan_amount?: number
  loan_to_value_ratio: number
  sort_order: number
  is_active: boolean
  created_at: string
  updated_at: string
  parent?: Category
  children?: Category[]
}

export interface CreateCategoryInput {
  name: string
  description?: string
  parent_id?: number
  icon?: string
  default_interest_rate?: number
  min_loan_amount?: number
  max_loan_amount?: number
  loan_to_value_ratio?: number
  sort_order?: number
}

export interface UpdateCategoryInput {
  name?: string
  description?: string
  parent_id?: number
  icon?: string
  default_interest_rate?: number
  min_loan_amount?: number
  max_loan_amount?: number
  loan_to_value_ratio?: number
  sort_order?: number
  is_active?: boolean
}

export interface CategoryListParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  search?: string
  parent_id?: number
  is_active?: boolean
}
