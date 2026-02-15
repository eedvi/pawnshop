import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import { Branch } from '@/types'

interface BranchState {
  // State
  branches: Branch[]
  selectedBranchId: number | null
  isLoading: boolean

  // Computed
  selectedBranch: Branch | null

  // Actions
  setBranches: (branches: Branch[]) => void
  setSelectedBranchId: (id: number | null) => void
  setLoading: (loading: boolean) => void
  reset: () => void
}

export const useBranchStore = create<BranchState>()(
  persist(
    (set, get) => ({
      // Initial state
      branches: [],
      selectedBranchId: null,
      isLoading: false,

      // Computed getter
      get selectedBranch() {
        const { branches, selectedBranchId } = get()
        return branches.find((b) => b.id === selectedBranchId) ?? null
      },

      // Actions
      setBranches: (branches) => {
        const { selectedBranchId } = get()

        // If current selection is not in new branches list, select first one
        const currentExists = branches.some((b) => b.id === selectedBranchId)

        set({
          branches,
          selectedBranchId: currentExists
            ? selectedBranchId
            : branches.length > 0
            ? branches[0].id
            : null,
        })
      },

      setSelectedBranchId: (id) => set({ selectedBranchId: id }),

      setLoading: (isLoading) => set({ isLoading }),

      reset: () =>
        set({
          branches: [],
          selectedBranchId: null,
          isLoading: false,
        }),
    }),
    {
      name: 'pawnshop_branch',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        selectedBranchId: state.selectedBranchId,
      }),
    }
  )
)
