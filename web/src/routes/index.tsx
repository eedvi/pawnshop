import { lazy, Suspense } from 'react'
import { createBrowserRouter, Outlet } from 'react-router-dom'
import { ProtectedRoute } from './protected-route'
import { PermissionRoute, ForbiddenPage } from './permission-route'
import { ROUTES } from './routes'

// Loading fallback
function PageLoader() {
  return (
    <div className="flex h-[50vh] items-center justify-center">
      <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
    </div>
  )
}

// Lazy load pages
const LoginPage = lazy(() => import('@/pages/auth/login-page'))
const TwoFactorPage = lazy(() => import('@/pages/auth/two-factor-page'))

const DashboardPage = lazy(() => import('@/pages/dashboard/dashboard-page'))

const CustomerListPage = lazy(() => import('@/pages/customers/customer-list-page'))
const CustomerCreatePage = lazy(() => import('@/pages/customers/customer-create-page'))
const CustomerDetailPage = lazy(() => import('@/pages/customers/customer-detail-page'))
const CustomerEditPage = lazy(() => import('@/pages/customers/customer-edit-page'))

const ItemListPage = lazy(() => import('@/pages/items/item-list-page'))
const ItemCreatePage = lazy(() => import('@/pages/items/item-create-page'))
const ItemDetailPage = lazy(() => import('@/pages/items/item-detail-page'))
const ItemEditPage = lazy(() => import('@/pages/items/item-edit-page'))

const LoanListPage = lazy(() => import('@/pages/loans/loan-list-page'))
const LoanCreatePage = lazy(() => import('@/pages/loans/loan-create-page'))
const LoanDetailPage = lazy(() => import('@/pages/loans/loan-detail-page'))

const PaymentListPage = lazy(() => import('@/pages/payments/payment-list-page'))
const PaymentCreatePage = lazy(() => import('@/pages/payments/payment-create-page'))
const PaymentDetailPage = lazy(() => import('@/pages/payments/payment-detail-page'))

const SaleListPage = lazy(() => import('@/pages/sales/sale-list-page'))
const SaleCreatePage = lazy(() => import('@/pages/sales/sale-create-page'))
const SaleDetailPage = lazy(() => import('@/pages/sales/sale-detail-page'))

const CashPage = lazy(() => import('@/pages/cash/cash-page'))

const ReportsPage = lazy(() => import('@/pages/reports/reports-page'))

const ExpenseListPage = lazy(() => import('@/pages/expenses/expense-list-page'))
const ExpenseCreatePage = lazy(() => import('@/pages/expenses/expense-create-page'))
const ExpenseDetailPage = lazy(() => import('@/pages/expenses/expense-detail-page'))

const TransferListPage = lazy(() => import('@/pages/transfers/transfer-list-page'))
const TransferCreatePage = lazy(() => import('@/pages/transfers/transfer-create-page'))
const TransferDetailPage = lazy(() => import('@/pages/transfers/transfer-detail-page'))

const UserListPage = lazy(() => import('@/pages/users/user-list-page'))
const UserCreatePage = lazy(() => import('@/pages/users/user-create-page'))
const UserDetailPage = lazy(() => import('@/pages/users/user-detail-page'))
const UserEditPage = lazy(() => import('@/pages/users/user-edit-page'))

const BranchListPage = lazy(() => import('@/pages/branches/branch-list-page'))
const BranchCreatePage = lazy(() => import('@/pages/branches/branch-create-page'))
const BranchDetailPage = lazy(() => import('@/pages/branches/branch-detail-page'))

const CategoryPage = lazy(() => import('@/pages/categories/category-page'))

const RoleListPage = lazy(() => import('@/pages/roles/role-list-page'))
const RoleCreatePage = lazy(() => import('@/pages/roles/role-create-page'))
const RoleDetailPage = lazy(() => import('@/pages/roles/role-detail-page'))
const RoleEditPage = lazy(() => import('@/pages/roles/role-edit-page'))

const NotificationPage = lazy(() => import('@/pages/notifications/notification-page'))

const AuditPage = lazy(() => import('@/pages/audit/audit-page'))

const SettingsPage = lazy(() => import('@/pages/settings/settings-page'))

// Layout will be imported once created
const AppLayout = lazy(() => import('@/components/layout/app-layout'))

// Helper to wrap pages with Suspense
function withSuspense(Component: React.ComponentType) {
  return (
    <Suspense fallback={<PageLoader />}>
      <Component />
    </Suspense>
  )
}

// Helper to wrap with permission check
function withPermission(Component: React.ComponentType, permission: string) {
  return (
    <PermissionRoute permission={permission} fallback={<ForbiddenPage />}>
      <Suspense fallback={<PageLoader />}>
        <Component />
      </Suspense>
    </PermissionRoute>
  )
}

export const router = createBrowserRouter([
  // Public routes
  {
    path: ROUTES.LOGIN,
    element: withSuspense(LoginPage),
  },
  {
    path: ROUTES.TWO_FACTOR,
    element: withSuspense(TwoFactorPage),
  },

  // Protected routes
  {
    path: '/',
    element: (
      <ProtectedRoute>
        <Suspense fallback={<PageLoader />}>
          <AppLayout />
        </Suspense>
      </ProtectedRoute>
    ),
    children: [
      // Dashboard
      {
        index: true,
        element: withSuspense(DashboardPage),
      },

      // Customers
      {
        path: 'customers',
        element: <Outlet />,
        children: [
          { index: true, element: withPermission(CustomerListPage, 'customers.read') },
          { path: 'new', element: withPermission(CustomerCreatePage, 'customers.create') },
          { path: ':id', element: withPermission(CustomerDetailPage, 'customers.read') },
          { path: ':id/edit', element: withPermission(CustomerEditPage, 'customers.update') },
        ],
      },

      // Items
      {
        path: 'items',
        element: <Outlet />,
        children: [
          { index: true, element: withPermission(ItemListPage, 'items.read') },
          { path: 'new', element: withPermission(ItemCreatePage, 'items.create') },
          { path: ':id', element: withPermission(ItemDetailPage, 'items.read') },
          { path: ':id/edit', element: withPermission(ItemEditPage, 'items.update') },
        ],
      },

      // Loans
      {
        path: 'loans',
        element: <Outlet />,
        children: [
          { index: true, element: withPermission(LoanListPage, 'loans.read') },
          { path: 'new', element: withPermission(LoanCreatePage, 'loans.create') },
          { path: ':id', element: withPermission(LoanDetailPage, 'loans.read') },
        ],
      },

      // Payments
      {
        path: 'payments',
        element: <Outlet />,
        children: [
          { index: true, element: withPermission(PaymentListPage, 'payments.read') },
          { path: 'new', element: withPermission(PaymentCreatePage, 'payments.create') },
          { path: ':id', element: withPermission(PaymentDetailPage, 'payments.read') },
        ],
      },

      // Sales
      {
        path: 'sales',
        element: <Outlet />,
        children: [
          { index: true, element: withPermission(SaleListPage, 'sales.read') },
          { path: 'new', element: withPermission(SaleCreatePage, 'sales.create') },
          { path: ':id', element: withPermission(SaleDetailPage, 'sales.read') },
        ],
      },

      // Cash
      {
        path: 'cash',
        element: withPermission(CashPage, 'cash.read'),
      },

      // Reports
      {
        path: 'reports',
        element: withPermission(ReportsPage, 'reports.read'),
      },

      // Expenses
      {
        path: 'expenses',
        element: <Outlet />,
        children: [
          { index: true, element: withPermission(ExpenseListPage, 'expenses.read') },
          { path: 'new', element: withPermission(ExpenseCreatePage, 'expenses.create') },
          { path: ':id', element: withPermission(ExpenseDetailPage, 'expenses.read') },
        ],
      },

      // Transfers
      {
        path: 'transfers',
        element: <Outlet />,
        children: [
          { index: true, element: withPermission(TransferListPage, 'transfers.read') },
          { path: 'new', element: withPermission(TransferCreatePage, 'transfers.create') },
          { path: ':id', element: withPermission(TransferDetailPage, 'transfers.read') },
        ],
      },

      // Users
      {
        path: 'users',
        element: <Outlet />,
        children: [
          { index: true, element: withPermission(UserListPage, 'users.read') },
          { path: 'new', element: withPermission(UserCreatePage, 'users.create') },
          { path: ':id', element: withPermission(UserDetailPage, 'users.read') },
          { path: ':id/edit', element: withPermission(UserEditPage, 'users.update') },
        ],
      },

      // Branches
      {
        path: 'branches',
        element: <Outlet />,
        children: [
          { index: true, element: withPermission(BranchListPage, 'branches.read') },
          { path: 'new', element: withPermission(BranchCreatePage, 'branches.create') },
          { path: ':id', element: withPermission(BranchDetailPage, 'branches.read') },
        ],
      },

      // Categories
      {
        path: 'categories',
        element: withPermission(CategoryPage, 'categories.read'),
      },

      // Roles
      {
        path: 'roles',
        element: <Outlet />,
        children: [
          { index: true, element: withPermission(RoleListPage, 'roles.read') },
          { path: 'new', element: withPermission(RoleCreatePage, 'roles.create') },
          { path: ':id', element: withPermission(RoleDetailPage, 'roles.read') },
          { path: ':id/edit', element: withPermission(RoleEditPage, 'roles.update') },
        ],
      },

      // Notifications
      {
        path: 'notifications',
        element: withPermission(NotificationPage, 'notifications.read'),
      },

      // Audit
      {
        path: 'audit',
        element: withPermission(AuditPage, 'audit.read'),
      },

      // Settings
      {
        path: 'settings',
        element: withPermission(SettingsPage, 'settings.read'),
      },
    ],
  },
])

export { ROUTES } from './routes'
