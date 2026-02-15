import axios, { AxiosError, AxiosInstance, AxiosRequestConfig, InternalAxiosRequestConfig } from 'axios'
import { ApiResponse, ApiErrorException } from '@/types/api'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

// Storage keys
const ACCESS_TOKEN_KEY = 'pawnshop_access_token'
const REFRESH_TOKEN_KEY = 'pawnshop_refresh_token'

// Token management
export const tokenStorage = {
  getAccessToken: () => localStorage.getItem(ACCESS_TOKEN_KEY),
  setAccessToken: (token: string) => localStorage.setItem(ACCESS_TOKEN_KEY, token),
  removeAccessToken: () => localStorage.removeItem(ACCESS_TOKEN_KEY),

  getRefreshToken: () => localStorage.getItem(REFRESH_TOKEN_KEY),
  setRefreshToken: (token: string) => localStorage.setItem(REFRESH_TOKEN_KEY, token),
  removeRefreshToken: () => localStorage.removeItem(REFRESH_TOKEN_KEY),

  clearAll: () => {
    localStorage.removeItem(ACCESS_TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
    localStorage.removeItem('pawnshop_user')
  }
}

// Create axios instance
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000, // 30 seconds
})

// Flag to prevent multiple refresh attempts
let isRefreshing = false
let refreshSubscribers: ((token: string) => void)[] = []

// Subscribe to token refresh
function subscribeTokenRefresh(callback: (token: string) => void) {
  refreshSubscribers.push(callback)
}

// Notify all subscribers with new token
function onTokenRefreshed(token: string) {
  refreshSubscribers.forEach((callback) => callback(token))
  refreshSubscribers = []
}

// Request interceptor - attach auth token
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = tokenStorage.getAccessToken()
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor - handle errors and token refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ApiResponse<unknown>>) => {
    const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean }

    // Handle 401 Unauthorized
    if (error.response?.status === 401 && !originalRequest._retry) {
      // Don't retry refresh or login endpoints
      if (originalRequest.url?.includes('/auth/refresh') || originalRequest.url?.includes('/auth/login')) {
        tokenStorage.clearAll()
        window.location.href = '/login'
        return Promise.reject(error)
      }

      originalRequest._retry = true

      if (!isRefreshing) {
        isRefreshing = true
        const refreshToken = tokenStorage.getRefreshToken()

        if (!refreshToken) {
          tokenStorage.clearAll()
          window.location.href = '/login'
          return Promise.reject(error)
        }

        try {
          const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refresh_token: refreshToken,
          })

          const { access_token, refresh_token } = response.data.data
          tokenStorage.setAccessToken(access_token)
          tokenStorage.setRefreshToken(refresh_token)

          isRefreshing = false
          onTokenRefreshed(access_token)

          // Retry original request with new token
          if (originalRequest.headers) {
            originalRequest.headers.Authorization = `Bearer ${access_token}`
          }
          return apiClient(originalRequest)
        } catch (refreshError) {
          isRefreshing = false
          tokenStorage.clearAll()
          window.location.href = '/login'
          return Promise.reject(refreshError)
        }
      }

      // Wait for token refresh and retry
      return new Promise((resolve) => {
        subscribeTokenRefresh((token: string) => {
          if (originalRequest.headers) {
            originalRequest.headers.Authorization = `Bearer ${token}`
          }
          resolve(apiClient(originalRequest))
        })
      })
    }

    // Transform error response
    if (error.response?.data?.error) {
      const apiError = error.response.data.error
      throw new ApiErrorException(apiError)
    }

    // Network error or other issues
    throw new ApiErrorException({
      code: 'NETWORK_ERROR',
      message: error.message || 'Error de conexión',
    })
  }
)

// Generic API helpers that unwrap the response envelope

export async function apiGet<T>(url: string, params?: Record<string, unknown>): Promise<T> {
  const response = await apiClient.get<ApiResponse<T>>(url, { params })
  if (!response.data.success) {
    throw new ApiErrorException(response.data.error || { code: 'UNKNOWN', message: 'Error desconocido' })
  }
  return response.data.data as T
}

export async function apiPost<T>(url: string, data?: unknown): Promise<T> {
  const response = await apiClient.post<ApiResponse<T>>(url, data)
  if (!response.data.success) {
    throw new ApiErrorException(response.data.error || { code: 'UNKNOWN', message: 'Error desconocido' })
  }
  return response.data.data as T
}

export async function apiPut<T>(url: string, data?: unknown): Promise<T> {
  const response = await apiClient.put<ApiResponse<T>>(url, data)
  if (!response.data.success) {
    throw new ApiErrorException(response.data.error || { code: 'UNKNOWN', message: 'Error desconocido' })
  }
  return response.data.data as T
}

export async function apiDelete<T>(url: string): Promise<T> {
  const response = await apiClient.delete<ApiResponse<T>>(url)
  // Handle 204 No Content responses
  if (response.status === 204 || !response.data) {
    return undefined as T
  }
  if (!response.data.success) {
    throw new ApiErrorException(response.data.error || { code: 'UNKNOWN', message: 'Error desconocido' })
  }
  return response.data.data as T
}

// For paginated endpoints - returns the full response with meta
export async function apiGetPaginated<T>(url: string, params?: Record<string, unknown>) {
  const response = await apiClient.get<ApiResponse<T[]>>(url, { params })
  if (!response.data.success) {
    throw new ApiErrorException(response.data.error || { code: 'UNKNOWN', message: 'Error desconocido' })
  }
  return {
    data: response.data.data as T[],
    meta: response.data.meta,
  }
}

// For file uploads
export async function apiUpload<T>(url: string, formData: FormData): Promise<T> {
  const response = await apiClient.post<ApiResponse<T>>(url, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
  if (!response.data.success) {
    throw new ApiErrorException(response.data.error || { code: 'UNKNOWN', message: 'Error desconocido' })
  }
  return response.data.data as T
}

// For downloading files (PDFs, etc.)
export async function apiDownload(url: string, filename: string): Promise<void> {
  const response = await apiClient.get(url, {
    responseType: 'blob',
  })

  const blob = new Blob([response.data])
  const downloadUrl = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = downloadUrl
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  window.URL.revokeObjectURL(downloadUrl)
}

export default apiClient
