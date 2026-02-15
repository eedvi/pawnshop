import { useEffect } from 'react'
import { Outlet } from 'react-router-dom'
import { useUIStore, useBranchStore } from '@/stores'
import { branchService } from '@/services/branch-service'
import { Sidebar } from './sidebar'
import { Header } from './header'
import { cn } from '@/lib/utils'

export default function AppLayout() {
  const { sidebarCollapsed } = useUIStore()
  const { branches, setBranches } = useBranchStore()

  // Load branches into store on mount
  useEffect(() => {
    if (branches.length === 0) {
      branchService.list({ is_active: true }).then((response) => {
        const branchList = response?.items || response?.data || response
        if (Array.isArray(branchList) && branchList.length > 0) {
          setBranches(branchList)
        }
      }).catch(console.error)
    }
  }, [branches.length, setBranches])

  return (
    <div className="min-h-screen bg-background">
      <Sidebar />
      <Header />
      <main
        className={cn(
          'pt-16 transition-all duration-300',
          sidebarCollapsed ? 'pl-16' : 'pl-64'
        )}
      >
        <div className="container mx-auto p-6">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
