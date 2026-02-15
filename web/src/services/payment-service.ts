import { apiGet, apiPost, apiGetPaginated } from '@/lib/api-client'
import {
  Payment,
  CreatePaymentInput,
  ReversePaymentInput,
  PaymentListParams,
  PayoffCalculation,
  MinimumPaymentCalculation,
} from '@/types'

export const paymentService = {
  // List payments with pagination
  list: async (params: PaymentListParams = {}) => {
    return apiGetPaginated<Payment>('/payments', params)
  },

  // Get a payment by ID
  getById: async (id: number): Promise<Payment> => {
    return apiGet<Payment>(`/payments/${id}`)
  },

  // Create a new payment
  create: async (input: CreatePaymentInput): Promise<Payment> => {
    return apiPost<Payment>('/payments', input)
  },

  // Reverse a payment
  reverse: async (id: number, input: ReversePaymentInput): Promise<void> => {
    return apiPost(`/payments/${id}/reverse`, input)
  },

  // Calculate payoff amount for a loan
  calculatePayoff: async (loanId: number): Promise<PayoffCalculation> => {
    return apiGet<PayoffCalculation>('/payments/calculate-payoff', { loan_id: loanId })
  },

  // Calculate minimum payment for a loan
  calculateMinimum: async (loanId: number): Promise<MinimumPaymentCalculation> => {
    return apiGet<MinimumPaymentCalculation>('/payments/calculate-minimum', { loan_id: loanId })
  },
}
