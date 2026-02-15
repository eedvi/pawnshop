// Application constants

// Currency formatting
export const CURRENCY_SYMBOL = 'Q'
export const CURRENCY_CODE = 'GTQ'
export const CURRENCY_LOCALE = 'es-GT'

// Date formatting
export const DATE_FORMAT = 'dd/MM/yyyy'
export const DATE_TIME_FORMAT = 'dd/MM/yyyy HH:mm'
export const TIME_FORMAT = 'HH:mm'

// Pagination defaults
export const DEFAULT_PAGE = 1
export const DEFAULT_PER_PAGE = 20
export const PER_PAGE_OPTIONS = [10, 20, 50, 100]

// Debounce delays (ms)
export const SEARCH_DEBOUNCE = 300

// Status colors for badges
export const STATUS_COLORS: Record<string, string> = {
  // Loan statuses
  active: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
  paid: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300',
  overdue: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300',
  defaulted: 'bg-red-200 text-red-900 dark:bg-red-900 dark:text-red-200',
  renewed: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300',
  confiscated: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300',

  // Payment statuses
  completed: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
  pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300',
  reversed: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300',
  failed: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300',

  // Item statuses
  available: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
  pawned: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300',
  collateral: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300',
  for_sale: 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300',
  sold: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300',
  transferred: 'bg-cyan-100 text-cyan-800 dark:bg-cyan-900 dark:text-cyan-300',
  in_transfer: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300',
  damaged: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300',
  lost: 'bg-red-200 text-red-900 dark:bg-red-900 dark:text-red-200',

  // Sale statuses
  cancelled: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300',
  refunded: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300',
  partial_refund: 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300',

  // Transfer statuses
  approved: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300',
  in_transit: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300',
  received: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',

  // Cash session statuses
  open: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
  closed: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300',

  // Expense statuses
  rejected: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300',

  // Boolean statuses
  true: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
  false: 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300',
}

// Status labels (Spanish)
export const STATUS_LABELS: Record<string, string> = {
  // Loan statuses
  active: 'Activo',
  paid: 'Pagado',
  overdue: 'Vencido',
  defaulted: 'En Mora',
  renewed: 'Renovado',
  confiscated: 'Confiscado',

  // Payment statuses
  completed: 'Completado',
  pending: 'Pendiente',
  reversed: 'Revertido',
  failed: 'Fallido',

  // Item statuses
  available: 'Disponible',
  pawned: 'Empeñado',
  collateral: 'En Garantía',
  for_sale: 'En Venta',
  sold: 'Vendido',
  transferred: 'Transferido',
  in_transfer: 'En Tránsito',
  damaged: 'Dañado',
  lost: 'Perdido',

  // Sale statuses
  cancelled: 'Cancelada',
  refunded: 'Reembolsada',
  partial_refund: 'Reembolso Parcial',

  // Transfer statuses
  approved: 'Aprobada',
  in_transit: 'En Tránsito',
  received: 'Recibida',

  // Cash session statuses
  open: 'Abierta',
  closed: 'Cerrada',

  // Expense statuses
  rejected: 'Rechazado',

  // Generic
  true: 'Sí',
  false: 'No',
}

// Sidebar navigation items
export const NAV_ITEMS = [
  { label: 'Dashboard', icon: 'LayoutDashboard', path: '/', permission: null },
  { label: 'Clientes', icon: 'Users', path: '/customers', permission: 'customers.read' },
  { label: 'Artículos', icon: 'Package', path: '/items', permission: 'items.read' },
  { label: 'Préstamos', icon: 'Banknote', path: '/loans', permission: 'loans.read' },
  { label: 'Pagos', icon: 'CreditCard', path: '/payments', permission: 'payments.read' },
  { label: 'Ventas', icon: 'ShoppingCart', path: '/sales', permission: 'sales.read' },
  { label: 'Caja', icon: 'Calculator', path: '/cash', permission: 'cash.read' },
  { label: 'Reportes', icon: 'BarChart3', path: '/reports', permission: 'reports.read' },
  { type: 'separator' as const },
  { label: 'Gastos', icon: 'Receipt', path: '/expenses', permission: 'expenses.read' },
  { label: 'Transferencias', icon: 'ArrowLeftRight', path: '/transfers', permission: 'transfers.read' },
  { type: 'separator' as const },
  { label: 'Usuarios', icon: 'UserCog', path: '/users', permission: 'users.read' },
  { label: 'Sucursales', icon: 'Building2', path: '/branches', permission: 'branches.read' },
  { label: 'Categorías', icon: 'FolderTree', path: '/categories', permission: 'categories.read' },
  { label: 'Roles', icon: 'Shield', path: '/roles', permission: 'roles.read' },
  { label: 'Notificaciones', icon: 'Bell', path: '/notifications', permission: 'notifications.read' },
  { label: 'Auditoría', icon: 'FileSearch', path: '/audit', permission: 'audit.read' },
  { label: 'Configuración', icon: 'Settings', path: '/settings', permission: 'settings.read' },
]
