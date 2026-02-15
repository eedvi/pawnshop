// API Response types - mirrors pkg/response/json.go

export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: ApiError
  meta?: ResponseMeta
}

export interface PaginatedResponse<T> {
  success: boolean
  data: T[]
  meta: ResponseMeta
}

export interface ApiError {
  code: string
  message: string
  details?: FieldError[]
}

export interface FieldError {
  field: string
  message: string
}

export interface ResponseMeta {
  request_id?: string
  timestamp: string
  pagination?: Pagination
}

export interface Pagination {
  current_page: number
  per_page: number
  total_items: number
  total_pages: number
}

export interface PaginationParams {
  page?: number
  per_page?: number
  order_by?: string
  order?: 'asc' | 'desc'
  search?: string
}

// Error class for API errors
export class ApiErrorException extends Error {
  code: string
  details?: FieldError[]

  constructor(error: ApiError) {
    super(error.message)
    this.name = 'ApiErrorException'
    this.code = error.code
    this.details = error.details
  }
}
