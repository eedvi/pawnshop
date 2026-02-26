import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { cashService } from '@/services/cash-service'
import type {
  CashSessionListParams,
  CreateCashRegisterInput,
  UpdateCashRegisterInput,
  OpenCashSessionInput,
  CloseCashSessionInput,
  CreateCashMovementInput,
} from '@/types'

export const cashKeys = {
  all: ['cash'] as const,
  registers: () => [...cashKeys.all, 'registers'] as const,
  registerList: (params?: { branch_id?: number }) => [...cashKeys.registers(), params] as const,
  register: (id: number) => [...cashKeys.registers(), id] as const,
  sessions: () => [...cashKeys.all, 'sessions'] as const,
  sessionList: (params?: CashSessionListParams) => [...cashKeys.sessions(), 'list', params] as const,
  session: (id: number) => [...cashKeys.sessions(), id] as const,
  sessionSummary: (id: number) => [...cashKeys.sessions(), id, 'summary'] as const,
  currentSession: (registerId: number) => [...cashKeys.sessions(), 'current', registerId] as const,
  movements: (sessionId: number) => [...cashKeys.all, 'movements', sessionId] as const,
}

// Register hooks
export function useCashRegisters(params?: { branch_id?: number }) {
  return useQuery({
    queryKey: cashKeys.registerList(params),
    queryFn: () => cashService.listRegisters(params),
  })
}

export function useCashRegister(id: number) {
  return useQuery({
    queryKey: cashKeys.register(id),
    queryFn: () => cashService.getRegister(id),
    enabled: id > 0,
  })
}

export function useCreateCashRegister() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateCashRegisterInput) => cashService.createRegister(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: cashKeys.registers() })
    },
  })
}

export function useUpdateCashRegister() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: UpdateCashRegisterInput }) =>
      cashService.updateRegister(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: cashKeys.registers() })
      queryClient.invalidateQueries({ queryKey: cashKeys.register(id) })
    },
  })
}

export function useDeleteCashRegister() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => cashService.deleteRegister(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: cashKeys.registers() })
    },
  })
}

// Session hooks
export function useCashSessions(params?: CashSessionListParams) {
  return useQuery({
    queryKey: cashKeys.sessionList(params),
    queryFn: () => cashService.listSessions(params),
  })
}

export function useCashSession(id: number) {
  return useQuery({
    queryKey: cashKeys.session(id),
    queryFn: () => cashService.getSession(id),
    enabled: id > 0,
  })
}

export function useCashSessionSummary(id: number) {
  return useQuery({
    queryKey: cashKeys.sessionSummary(id),
    queryFn: () => cashService.getSessionSummary(id),
    enabled: id > 0,
  })
}

export function useCurrentCashSession(registerId: number) {
  return useQuery({
    queryKey: cashKeys.currentSession(registerId),
    queryFn: () => cashService.getCurrentSession(registerId),
    enabled: registerId > 0,
  })
}

export function useOpenCashSession() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: OpenCashSessionInput) => cashService.openSession(input),
    onSuccess: (_, { cash_register_id }) => {
      queryClient.invalidateQueries({ queryKey: cashKeys.sessions() })
      queryClient.invalidateQueries({ queryKey: cashKeys.currentSession(cash_register_id) })
    },
  })
}

export function useCloseCashSession() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: CloseCashSessionInput }) =>
      cashService.closeSession(id, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: cashKeys.sessions() })
    },
  })
}

// Movement hooks
export function useCashMovements(sessionId: number) {
  return useQuery({
    queryKey: cashKeys.movements(sessionId),
    queryFn: () => cashService.listMovements(sessionId),
    enabled: sessionId > 0,
  })
}

export function useCreateCashMovement() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateCashMovementInput) => cashService.createMovement(input),
    onSuccess: (_, { session_id }) => {
      queryClient.invalidateQueries({ queryKey: cashKeys.movements(session_id) })
      queryClient.invalidateQueries({ queryKey: cashKeys.sessionSummary(session_id) })
    },
  })
}

export function useDeleteCashMovement() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, sessionId }: { id: number; sessionId: number }) =>
      cashService.deleteMovement(id),
    onSuccess: (_, { sessionId }) => {
      queryClient.invalidateQueries({ queryKey: cashKeys.movements(sessionId) })
      queryClient.invalidateQueries({ queryKey: cashKeys.sessionSummary(sessionId) })
    },
  })
}
