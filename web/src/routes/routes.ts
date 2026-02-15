// Route path constants
export const ROUTES = {
  // Public
  LOGIN: '/login',
  TWO_FACTOR: '/two-factor',
  FORGOT_PASSWORD: '/forgot-password',
  RESET_PASSWORD: '/reset-password',

  // Dashboard
  DASHBOARD: '/',

  // Customers
  CUSTOMERS: '/customers',
  CUSTOMER_CREATE: '/customers/new',
  CUSTOMER_DETAIL: '/customers/:id',
  CUSTOMER_EDIT: '/customers/:id/edit',

  // Items
  ITEMS: '/items',
  ITEM_CREATE: '/items/new',
  ITEM_DETAIL: '/items/:id',
  ITEM_EDIT: '/items/:id/edit',

  // Loans
  LOANS: '/loans',
  LOAN_CREATE: '/loans/new',
  LOAN_DETAIL: '/loans/:id',

  // Payments
  PAYMENTS: '/payments',
  PAYMENT_CREATE: '/payments/new',
  PAYMENT_DETAIL: '/payments/:id',

  // Sales
  SALES: '/sales',
  SALE_CREATE: '/sales/new',
  SALE_DETAIL: '/sales/:id',

  // Cash
  CASH: '/cash',

  // Reports
  REPORTS: '/reports',

  // Expenses
  EXPENSES: '/expenses',
  EXPENSE_CREATE: '/expenses/new',
  EXPENSE_DETAIL: '/expenses/:id',

  // Transfers
  TRANSFERS: '/transfers',
  TRANSFER_CREATE: '/transfers/new',
  TRANSFER_DETAIL: '/transfers/:id',

  // Admin
  USERS: '/users',
  USER_CREATE: '/users/new',
  USER_DETAIL: '/users/:id',

  BRANCHES: '/branches',
  BRANCH_CREATE: '/branches/new',
  BRANCH_DETAIL: '/branches/:id',

  CATEGORIES: '/categories',

  ROLES: '/roles',
  ROLE_CREATE: '/roles/new',
  ROLE_DETAIL: '/roles/:id',

  NOTIFICATIONS: '/notifications',

  AUDIT: '/audit',

  SETTINGS: '/settings',
} as const

// Helper to generate dynamic routes
export function customerRoute(id: number | string) {
  return `/customers/${id}`
}

export function customerEditRoute(id: number | string) {
  return `/customers/${id}/edit`
}

export function itemRoute(id: number | string) {
  return `/items/${id}`
}

export function itemEditRoute(id: number | string) {
  return `/items/${id}/edit`
}

export function loanRoute(id: number | string) {
  return `/loans/${id}`
}

export function paymentRoute(id: number | string) {
  return `/payments/${id}`
}

export function saleRoute(id: number | string) {
  return `/sales/${id}`
}

export function expenseRoute(id: number | string) {
  return `/expenses/${id}`
}

export function transferRoute(id: number | string) {
  return `/transfers/${id}`
}

export function userRoute(id: number | string) {
  return `/users/${id}`
}

export function userEditRoute(id: number | string) {
  return `/users/${id}/edit`
}

export function branchRoute(id: number | string) {
  return `/branches/${id}`
}

export function roleRoute(id: number | string) {
  return `/roles/${id}`
}

export function roleEditRoute(id: number | string) {
  return `/roles/${id}/edit`
}
