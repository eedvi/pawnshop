import { Link, useLocation } from 'react-router-dom'
import { cn } from '@/lib/utils'
import { useAuthStore, useUIStore } from '@/stores'
import { NAV_ITEMS } from '@/lib/constants'
import {
  LayoutDashboard,
  Users,
  Package,
  Banknote,
  CreditCard,
  ShoppingCart,
  Calculator,
  BarChart3,
  Receipt,
  ArrowLeftRight,
  UserCog,
  Building2,
  FolderTree,
  Shield,
  Bell,
  FileSearch,
  Settings,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'

// Icon mapping
const iconMap: Record<string, React.ComponentType<{ className?: string }>> = {
  LayoutDashboard,
  Users,
  Package,
  Banknote,
  CreditCard,
  ShoppingCart,
  Calculator,
  BarChart3,
  Receipt,
  ArrowLeftRight,
  UserCog,
  Building2,
  FolderTree,
  Shield,
  Bell,
  FileSearch,
  Settings,
}

export function Sidebar() {
  const location = useLocation()
  const { hasPermission } = useAuthStore()
  const { sidebarCollapsed, toggleSidebar } = useUIStore()

  // Filter nav items by permissions
  const visibleItems = NAV_ITEMS.filter((item) => {
    if ('type' in item && item.type === 'separator') return true
    if (!('permission' in item)) return true
    if (!item.permission) return true
    return hasPermission(item.permission)
  })

  // Remove consecutive separators and leading/trailing separators
  const cleanedItems = visibleItems.filter((item, index, arr) => {
    if (!('type' in item) || item.type !== 'separator') return true
    // Remove if first or last
    if (index === 0 || index === arr.length - 1) return false
    // Remove if previous item is also separator
    const prev = arr[index - 1]
    if ('type' in prev && prev.type === 'separator') return false
    return true
  })

  return (
    <aside
      className={cn(
        'fixed left-0 top-0 z-40 flex h-screen flex-col border-r bg-card transition-all duration-300',
        sidebarCollapsed ? 'w-16' : 'w-64'
      )}
    >
      {/* Logo */}
      <div className="flex h-16 items-center justify-between border-b px-4">
        {!sidebarCollapsed && (
          <Link to="/" className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground font-bold">
              P
            </div>
            <span className="font-semibold">PawnShop</span>
          </Link>
        )}
        <Button
          variant="ghost"
          size="icon"
          onClick={toggleSidebar}
          className={cn('h-8 w-8', sidebarCollapsed && 'mx-auto')}
        >
          {sidebarCollapsed ? (
            <ChevronRight className="h-4 w-4" />
          ) : (
            <ChevronLeft className="h-4 w-4" />
          )}
        </Button>
      </div>

      {/* Navigation */}
      <nav className="flex-1 overflow-y-auto p-2">
        <TooltipProvider delayDuration={0}>
          <ul className="space-y-1">
            {cleanedItems.map((item, index) => {
              if ('type' in item && item.type === 'separator') {
                return (
                  <li key={`sep-${index}`} className="my-2">
                    <div className="h-px bg-border" />
                  </li>
                )
              }

              const navItem = item as {
                label: string
                icon: string
                path: string
                permission: string | null
              }
              const Icon = iconMap[navItem.icon]
              const isActive =
                navItem.path === '/'
                  ? location.pathname === '/'
                  : location.pathname.startsWith(navItem.path)

              const linkContent = (
                <Link
                  to={navItem.path}
                  className={cn(
                    'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                    isActive
                      ? 'bg-primary text-primary-foreground'
                      : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
                    sidebarCollapsed && 'justify-center px-2'
                  )}
                >
                  {Icon && <Icon className="h-5 w-5 flex-shrink-0" />}
                  {!sidebarCollapsed && <span>{navItem.label}</span>}
                </Link>
              )

              if (sidebarCollapsed) {
                return (
                  <li key={navItem.path}>
                    <Tooltip>
                      <TooltipTrigger asChild>{linkContent}</TooltipTrigger>
                      <TooltipContent side="right">
                        {navItem.label}
                      </TooltipContent>
                    </Tooltip>
                  </li>
                )
              }

              return <li key={navItem.path}>{linkContent}</li>
            })}
          </ul>
        </TooltipProvider>
      </nav>
    </aside>
  )
}
