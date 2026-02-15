import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { paymentService } from '@/services/payment-service'
import { CreatePaymentInput, ReversePaymentInput, PaymentListParams } from '@/types'
import { loanKeys } from './use-loans'
import { toast } from 'sonner'

// Query keys
export const paymentKeys = {
  all: ['payments'] as const,
  lists: () => [...paymentKeys.all, 'list'] as const,
  list: (params: PaymentListParams) => [...paymentKeys.lists(), params] as const,
  details: () => [...paymentKeys.all, 'detail'] as const,
  detail: (id: number) => [...paymentKeys.details(), id] as const,
  payoff: (loanId: number) => ['loans', loanId, 'payoff'] as const,
  minimum: (loanId: number) => ['loans', loanId, 'minimum'] as const,
}

// Hook to list payments
export function usePayments(params: PaymentListParams = {}) {
  return useQuery({
    queryKey: paymentKeys.list(params),
    queryFn: () => paymentService.list(params),
  })
}

// Hook to get a single payment
export function usePayment(id: number) {
  return useQuery({
    queryKey: paymentKeys.detail(id),
    queryFn: () => paymentService.getById(id),
    enabled: id > 0,
  })
}

// Hook to calculate payoff
export function usePayoffCalculation(loanId: number) {
  return useQuery({
    queryKey: paymentKeys.payoff(loanId),
    queryFn: () => paymentService.calculatePayoff(loanId),
    enabled: loanId > 0,
  })
}

// Hook to calculate minimum payment
export function useMinimumPaymentCalculation(loanId: number) {
  return useQuery({
    queryKey: paymentKeys.minimum(loanId),
    queryFn: () => paymentService.calculateMinimum(loanId),
    enabled: loanId > 0,
  })
}

// Hook to create a payment
export function useCreatePayment() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (input: CreatePaymentInput) => paymentService.create(input),
    onSuccess: (_, { loan_id }) => {
      queryClient.invalidateQueries({ queryKey: paymentKeys.lists() })
      queryClient.invalidateQueries({ queryKey: loanKeys.detail(loan_id) })
      queryClient.invalidateQueries({ queryKey: loanKeys.lists() })
      toast.success('Pago registrado exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al registrar el pago')
    },
  })
}

// Hook to reverse a payment
export function useReversePayment() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, input }: { id: number; input: ReversePaymentInput }) =>
      paymentService.reverse(id, input),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: paymentKeys.lists() })
      queryClient.invalidateQueries({ queryKey: paymentKeys.detail(id) })
      queryClient.invalidateQueries({ queryKey: loanKeys.lists() })
      toast.success('Pago revertido exitosamente')
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Error al revertir el pago')
    },
  })
}
