// Role types - mirrors internal/domain/role.go

export interface Role {
  id: number
  name: string
  display_name: string
  description?: string
  permissions: string[]
  is_system: boolean
  created_at: string
  updated_at: string
}

export interface CreateRoleInput {
  name: string
  display_name: string
  description?: string
  permissions: string[]
}

export interface UpdateRoleInput {
  name?: string
  display_name?: string
  description?: string
  permissions?: string[]
}

// Predefined role names
export const ROLE_SUPER_ADMIN = 'super_admin'
export const ROLE_ADMIN = 'admin'
export const ROLE_MANAGER = 'manager'
export const ROLE_CASHIER = 'cashier'
export const ROLE_SELLER = 'seller'

// Permission groups
export const PERMISSION_GROUPS = {
  customers: ['customers.read', 'customers.create', 'customers.update', 'customers.delete'],
  items: ['items.read', 'items.create', 'items.update', 'items.delete'],
  loans: ['loans.read', 'loans.create', 'loans.update', 'loans.delete'],
  payments: ['payments.read', 'payments.create', 'payments.update', 'payments.delete'],
  sales: ['sales.read', 'sales.create', 'sales.update', 'sales.delete'],
  cash: ['cash.read', 'cash.create', 'cash.update', 'cash.delete'],
  reports: ['reports.read', 'reports.export'],
  users: ['users.read', 'users.create', 'users.update', 'users.delete'],
  branches: ['branches.read', 'branches.create', 'branches.update', 'branches.delete'],
  categories: ['categories.read', 'categories.create', 'categories.update', 'categories.delete'],
  roles: ['roles.read', 'roles.create', 'roles.update', 'roles.delete'],
  settings: ['settings.read', 'settings.update'],
  audit: ['audit.read'],
  notifications: ['notifications.read', 'notifications.create', 'notifications.manage'],
  expenses: ['expenses.read', 'expenses.create', 'expenses.update', 'expenses.delete', 'expenses.approve'],
  transfers: ['transfers.read', 'transfers.create', 'transfers.approve', 'transfers.ship', 'transfers.receive', 'transfers.cancel'],
}

// Check if user has permission (with wildcard support)
export function hasPermission(userPermissions: string[], permission: string): boolean {
  for (const p of userPermissions) {
    if (p === '*' || p === permission) {
      return true
    }
    // Check wildcard permissions (e.g., "customers.*" matches "customers.read")
    if (p.endsWith('.*')) {
      const prefix = p.slice(0, -2)
      if (permission.startsWith(prefix)) {
        return true
      }
    }
  }
  return false
}
