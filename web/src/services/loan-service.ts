import { apiGet, apiPost, apiGetPaginated } from '@/lib/api-client'
import { Loan, LoanInstallment, CreateLoanInput, RenewLoanInput, LoanListParams } from '@/types'

export const loanService = {
  // List loans with pagination
  list: async (params: LoanListParams = {}) => {
    return apiGetPaginated<Loan>('/loans', params)
  },

  // Get a loan by ID
  getById: async (id: number): Promise<Loan> => {
    return apiGet<Loan>(`/loans/${id}`)
  },

  // Get loan installments
  getInstallments: async (id: number): Promise<LoanInstallment[]> => {
    return apiGet<LoanInstallment[]>(`/loans/${id}/installments`)
  },

  // Create a new loan
  create: async (input: CreateLoanInput): Promise<Loan> => {
    return apiPost<Loan>('/loans', input)
  },

  // Calculate loan terms (preview before creating)
  calculate: async (input: CreateLoanInput): Promise<{
    loan_amount: number
    interest_rate: number
    interest_amount: number
    total_amount: number
    installment_amount?: number
    installments?: LoanInstallment[]
  }> => {
    return apiPost('/loans/calculate', input)
  },

  // Renew a loan
  renew: async (id: number, input: RenewLoanInput): Promise<Loan> => {
    return apiPost<Loan>(`/loans/${id}/renew`, input)
  },

  // Confiscate (when loan defaults)
  confiscate: async (id: number, notes?: string): Promise<void> => {
    return apiPost(`/loans/${id}/confiscate`, { notes })
  },

  // Get customer's active loans
  getCustomerActiveLoans: async (customerId: number): Promise<Loan[]> => {
    const result = await apiGetPaginated<Loan>('/loans', {
      customer_id: customerId,
      status: 'active',
      per_page: 100
    })
    return result.data || []
  },
}
