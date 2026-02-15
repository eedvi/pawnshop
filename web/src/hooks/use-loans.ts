import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { loanService } from '@/services/loan-service'
import { CreateLoanInput, RenewLoanInput, LoanListParams } from '@/types'
import { toast } from 'sonner'

// Query keys
export const loanKeys = {
  all: ['loans'] as const,
  lists: () => [...loanKeys.all, 'list'] as const,
  list: (params: LoanListParams) => [...loanKeys.lists(), params] as const,
  details: () => [...loanKeys.all, 'detail'] as const,
  detail: (id: number) => [...loanKeys.details(), id] as const,
  installments: (id: number) => [...loanKeys.detail(id), 'installments'] as const,
  customerActive: (customerId: number) => ['customers', customerId, 'loans', 'active'] as const,
}

// Hook to list loans
export function useLoans(params: LoanListParams = {}) {
  return useQuery({
    queryKey: loanKeys.list(params),
    queryFn: () => loanService.list(params),
  })
}

// Hook to get a single loan
export function useLoan(id: number) {
  return useQuery({
    queryKey: loanKeys.detail(id),
    queryFn: () => loanService.getById(id),
    enabled: id > 0,
  })
}

// Hook to get loan installments
export function useLoanInstallments(id: number) {
  return useQuery({
    queryKey: loanKeys.installments(id),
    queryFn: () => loanService.getInstallments(id),
    enabled: id > 0,
  })
}

// Hook to get customer's active loans
export function useCustomerActiveLoans(customerId: number) {
  return useQuery({
    queryKey: loanKeys.customerActive(customerId),
    queryFn: () => loanService.getCustomerActiveLoans(customerId),
    enabled: customerId > 0,
  })
}

// Hook to create a loan
export function useCreateLoan() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreateLoanInput) => loanService.create(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: loanKeys.lists() })
      toast.success('Préstamo creado exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al crear el préstamo')
    },
  })
}

// Hook to calculate loan terms (preview)
export function useCalculateLoan() {
  return useMutation({
    mutationFn: (input: CreateLoanInput) => loanService.calculate(input),
    onError: (error: Error) => {
      toast.error(error.message || 'Error al calcular el préstamo')
    },
  })
}

// Hook to renew a loan
export function useRenewLoan() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: RenewLoanInput }) =>
      loanService.renew(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: loanKeys.lists() })
      queryClient.invalidateQueries({ queryKey: loanKeys.detail(id) })
      toast.success('Préstamo renovado exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al renovar el préstamo')
    },
  })
}

// Hook to confiscate a loan
export function useConfiscateLoan() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, notes }: { id: number; notes?: string }) =>
      loanService.confiscate(id, notes),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: loanKeys.lists() })
      queryClient.invalidateQueries({ queryKey: loanKeys.detail(id) })
      toast.success('Préstamo confiscado exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al confiscar el préstamo')
    },
  })
}
