import { apiGet, apiGetPaginated, apiPost, apiPut, apiDelete } from '@/lib/api-client'
import type {
  CashRegister,
  CashSession,
  CashMovement,
  CashSessionSummary,
  CreateCashRegisterInput,
  UpdateCashRegisterInput,
  OpenCashSessionInput,
  CloseCashSessionInput,
  CreateCashMovementInput,
  CashSessionListParams,
  ApiResponse,
} from '@/types'

export const cashService = {
  // Cash Registers
  listRegisters: (params?: { branch_id?: number }) =>
    apiGet<CashRegister[]>('/cash/registers', params),

  getRegister: (id: number) =>
    apiGet<CashRegister>(`/cash/registers/${id}`),

  createRegister: (input: CreateCashRegisterInput) =>
    apiPost<CashRegister>('/cash/registers', input),

  updateRegister: (id: number, input: UpdateCashRegisterInput) =>
    apiPut<CashRegister>(`/cash/registers/${id}`, input),

  deleteRegister: (id: number) =>
    apiDelete<void>(`/cash/registers/${id}`),

  // Cash Sessions
  listSessions: (params?: CashSessionListParams) =>
    apiGetPaginated<CashSession>('/cash/sessions', params),

  getSession: (id: number) =>
    apiGet<CashSession>(`/cash/sessions/${id}`),

  getSessionSummary: (id: number) =>
    apiGet<CashSessionSummary>(`/cash/sessions/${id}/summary`),

  openSession: (input: OpenCashSessionInput) =>
    apiPost<CashSession>('/cash/sessions/open', input),

  closeSession: (id: number, input: CloseCashSessionInput) =>
    apiPost<CashSession>(`/cash/sessions/${id}/close`, input),

  getCurrentSession: (registerId: number) =>
    apiGet<CashSession>(`/cash/registers/${registerId}/current-session`),

  // Cash Movements
  listMovements: (sessionId: number) =>
    apiGet<CashMovement[]>(`/cash/sessions/${sessionId}/movements`),

  createMovement: (input: CreateCashMovementInput) =>
    apiPost<CashMovement>('/cash/movements', input),

  deleteMovement: (id: number) =>
    apiDelete<void>(`/cash/movements/${id}`),
}
