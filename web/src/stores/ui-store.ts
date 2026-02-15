import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'

type Theme = 'light' | 'dark' | 'system'

interface UIState {
  // State
  sidebarCollapsed: boolean
  theme: Theme
  isMobileMenuOpen: boolean

  // Actions
  toggleSidebar: () => void
  setSidebarCollapsed: (collapsed: boolean) => void
  setTheme: (theme: Theme) => void
  toggleMobileMenu: () => void
  closeMobileMenu: () => void
}

// Apply theme to document
function applyTheme(theme: Theme) {
  const root = document.documentElement
  const systemDark = window.matchMedia('(prefers-color-scheme: dark)').matches

  if (theme === 'dark' || (theme === 'system' && systemDark)) {
    root.classList.add('dark')
  } else {
    root.classList.remove('dark')
  }
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      // Initial state
      sidebarCollapsed: false,
      theme: 'system',
      isMobileMenuOpen: false,

      // Actions
      toggleSidebar: () =>
        set((state) => ({
          sidebarCollapsed: !state.sidebarCollapsed,
        })),

      setSidebarCollapsed: (collapsed) =>
        set({ sidebarCollapsed: collapsed }),

      setTheme: (theme) => {
        applyTheme(theme)
        set({ theme })
      },

      toggleMobileMenu: () =>
        set((state) => ({
          isMobileMenuOpen: !state.isMobileMenuOpen,
        })),

      closeMobileMenu: () => set({ isMobileMenuOpen: false }),
    }),
    {
      name: 'pawnshop_ui',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        sidebarCollapsed: state.sidebarCollapsed,
        theme: state.theme,
      }),
      onRehydrateStorage: () => (state) => {
        // Apply theme after rehydration
        if (state?.theme) {
          applyTheme(state.theme)
        }
      },
    }
  )
)

// Listen for system theme changes
if (typeof window !== 'undefined') {
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
    const theme = useUIStore.getState().theme
    if (theme === 'system') {
      applyTheme('system')
    }
  })
}
