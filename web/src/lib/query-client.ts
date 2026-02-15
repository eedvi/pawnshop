import { QueryClient } from '@tanstack/react-query'
import { ApiErrorException } from '@/types/api'

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // Data is fresh for 30 seconds
      staleTime: 30_000,
      // Cache data for 5 minutes
      gcTime: 5 * 60 * 1000,
      // Retry once on failure
      retry: 1,
      // Don't refetch on window focus (can be noisy)
      refetchOnWindowFocus: false,
      // Refetch on reconnect
      refetchOnReconnect: true,
    },
    mutations: {
      // Don't retry mutations by default
      retry: false,
      // Handle errors globally if needed
      onError: (error) => {
        if (error instanceof ApiErrorException) {
          console.error(`API Error: ${error.code} - ${error.message}`)
        } else {
          console.error('Mutation error:', error)
        }
      },
    },
  },
})

// Query keys factory for consistent key management
export const queryKeys = {
  // Auth
  auth: {
    me: ['auth', 'me'] as const,
  },

  // Customers
  customers: {
    all: ['customers'] as const,
    lists: () => [...queryKeys.customers.all, 'list'] as const,
    list: (params: Record<string, unknown>) => [...queryKeys.customers.lists(), params] as const,
    details: () => [...queryKeys.customers.all, 'detail'] as const,
    detail: (id: number) => [...queryKeys.customers.details(), id] as const,
  },

  // Items
  items: {
    all: ['items'] as const,
    lists: () => [...queryKeys.items.all, 'list'] as const,
    list: (params: Record<string, unknown>) => [...queryKeys.items.lists(), params] as const,
    details: () => [...queryKeys.items.all, 'detail'] as const,
    detail: (id: number) => [...queryKeys.items.details(), id] as const,
    forSale: () => [...queryKeys.items.all, 'for-sale'] as const,
  },

  // Loans
  loans: {
    all: ['loans'] as const,
    lists: () => [...queryKeys.loans.all, 'list'] as const,
    list: (params: Record<string, unknown>) => [...queryKeys.loans.lists(), params] as const,
    details: () => [...queryKeys.loans.all, 'detail'] as const,
    detail: (id: number) => [...queryKeys.loans.details(), id] as const,
    payments: (id: number) => [...queryKeys.loans.detail(id), 'payments'] as const,
    installments: (id: number) => [...queryKeys.loans.detail(id), 'installments'] as const,
    overdue: () => [...queryKeys.loans.all, 'overdue'] as const,
  },

  // Payments
  payments: {
    all: ['payments'] as const,
    lists: () => [...queryKeys.payments.all, 'list'] as const,
    list: (params: Record<string, unknown>) => [...queryKeys.payments.lists(), params] as const,
    details: () => [...queryKeys.payments.all, 'detail'] as const,
    detail: (id: number) => [...queryKeys.payments.details(), id] as const,
    calculatePayoff: (loanId: number) => [...queryKeys.payments.all, 'payoff', loanId] as const,
    calculateMinimum: (loanId: number) => [...queryKeys.payments.all, 'minimum', loanId] as const,
  },

  // Sales
  sales: {
    all: ['sales'] as const,
    lists: () => [...queryKeys.sales.all, 'list'] as const,
    list: (params: Record<string, unknown>) => [...queryKeys.sales.lists(), params] as const,
    details: () => [...queryKeys.sales.all, 'detail'] as const,
    detail: (id: number) => [...queryKeys.sales.details(), id] as const,
    summary: (params: Record<string, unknown>) => [...queryKeys.sales.all, 'summary', params] as const,
  },

  // Branches
  branches: {
    all: ['branches'] as const,
    lists: () => [...queryKeys.branches.all, 'list'] as const,
    list: (params: Record<string, unknown>) => [...queryKeys.branches.lists(), params] as const,
    details: () => [...queryKeys.branches.all, 'detail'] as const,
    detail: (id: number) => [...queryKeys.branches.details(), id] as const,
  },

  // Categories
  categories: {
    all: ['categories'] as const,
    lists: () => [...queryKeys.categories.all, 'list'] as const,
    list: (params: Record<string, unknown>) => [...queryKeys.categories.lists(), params] as const,
    tree: () => [...queryKeys.categories.all, 'tree'] as const,
    details: () => [...queryKeys.categories.all, 'detail'] as const,
    detail: (id: number) => [...queryKeys.categories.details(), id] as const,
  },

  // Users
  users: {
    all: ['users'] as const,
    lists: () => [...queryKeys.users.all, 'list'] as const,
    list: (params: Record<string, unknown>) => [...queryKeys.users.lists(), params] as const,
    details: () => [...queryKeys.users.all, 'detail'] as const,
    detail: (id: number) => [...queryKeys.users.details(), id] as const,
  },

  // Roles
  roles: {
    all: ['roles'] as const,
    lists: () => [...queryKeys.roles.all, 'list'] as const,
    list: () => [...queryKeys.roles.lists()] as const,
    details: () => [...queryKeys.roles.all, 'detail'] as const,
    detail: (id: number) => [...queryKeys.roles.details(), id] as const,
    permissions: () => [...queryKeys.roles.all, 'permissions'] as const,
  },

  // Cash
  cash: {
    registers: {
      all: ['cash', 'registers'] as const,
      list: (params: Record<string, unknown>) => [...queryKeys.cash.registers.all, 'list', params] as const,
      detail: (id: number) => [...queryKeys.cash.registers.all, 'detail', id] as const,
    },
    sessions: {
      all: ['cash', 'sessions'] as const,
      list: (params: Record<string, unknown>) => [...queryKeys.cash.sessions.all, 'list', params] as const,
      detail: (id: number) => [...queryKeys.cash.sessions.all, 'detail', id] as const,
      current: () => [...queryKeys.cash.sessions.all, 'current'] as const,
      summary: (id: number) => [...queryKeys.cash.sessions.all, 'summary', id] as const,
      movements: (id: number) => [...queryKeys.cash.sessions.all, 'movements', id] as const,
    },
    movements: {
      all: ['cash', 'movements'] as const,
      list: (params: Record<string, unknown>) => [...queryKeys.cash.movements.all, 'list', params] as const,
    },
  },

  // Reports
  reports: {
    dashboard: (branchId?: number) => ['reports', 'dashboard', branchId] as const,
    loans: (params: Record<string, unknown>) => ['reports', 'loans', params] as const,
    payments: (params: Record<string, unknown>) => ['reports', 'payments', params] as const,
    sales: (params: Record<string, unknown>) => ['reports', 'sales', params] as const,
    overdue: (params: Record<string, unknown>) => ['reports', 'overdue', params] as const,
  },

  // Settings
  settings: {
    all: ['settings'] as const,
    list: () => [...queryKeys.settings.all, 'list'] as const,
    merged: (branchId?: number) => [...queryKeys.settings.all, 'merged', branchId] as const,
    byKey: (key: string) => [...queryKeys.settings.all, 'key', key] as const,
  },

  // Notifications
  notifications: {
    internal: {
      all: ['notifications', 'internal'] as const,
      me: () => [...queryKeys.notifications.internal.all, 'me'] as const,
      unread: () => [...queryKeys.notifications.internal.all, 'unread'] as const,
      unreadCount: () => [...queryKeys.notifications.internal.all, 'unread-count'] as const,
    },
    templates: {
      all: ['notifications', 'templates'] as const,
      list: () => [...queryKeys.notifications.templates.all, 'list'] as const,
      detail: (id: number) => [...queryKeys.notifications.templates.all, 'detail', id] as const,
    },
  },

  // Audit
  audit: {
    all: ['audit'] as const,
    list: (params: Record<string, unknown>) => [...queryKeys.audit.all, 'list', params] as const,
  },

  // Transfers
  transfers: {
    all: ['transfers'] as const,
    lists: () => [...queryKeys.transfers.all, 'list'] as const,
    list: (params: Record<string, unknown>) => [...queryKeys.transfers.lists(), params] as const,
    details: () => [...queryKeys.transfers.all, 'detail'] as const,
    detail: (id: number) => [...queryKeys.transfers.details(), id] as const,
  },

  // Expenses
  expenses: {
    all: ['expenses'] as const,
    lists: () => [...queryKeys.expenses.all, 'list'] as const,
    list: (params: Record<string, unknown>) => [...queryKeys.expenses.lists(), params] as const,
    details: () => [...queryKeys.expenses.all, 'detail'] as const,
    detail: (id: number) => [...queryKeys.expenses.details(), id] as const,
    categories: () => [...queryKeys.expenses.all, 'categories'] as const,
  },
}
