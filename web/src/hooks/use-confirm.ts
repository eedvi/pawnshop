import { useState, useCallback } from 'react'

interface ConfirmState {
  isOpen: boolean
  title: string
  description?: string
  confirmLabel?: string
  cancelLabel?: string
  variant?: 'default' | 'destructive'
  onConfirm: () => void
}

interface UseConfirmReturn {
  isOpen: boolean
  open: boolean // Alias for ConfirmDialog compatibility
  title: string
  description?: string
  confirmLabel?: string
  cancelLabel?: string
  variant?: 'default' | 'destructive'
  onConfirm: () => void
  onOpenChange: (open: boolean) => void
  confirm: (options: ConfirmOptions) => Promise<boolean>
}

interface ConfirmOptions {
  title: string
  description?: string
  confirmLabel?: string
  cancelLabel?: string
  variant?: 'default' | 'destructive'
}

/**
 * Hook for managing confirmation dialogs imperatively.
 *
 * Usage:
 * ```tsx
 * const confirm = useConfirm()
 *
 * const handleDelete = async () => {
 *   const confirmed = await confirm.confirm({
 *     title: 'Eliminar cliente',
 *     description: '¿Estás seguro de eliminar este cliente?',
 *     variant: 'destructive',
 *   })
 *
 *   if (confirmed) {
 *     // Perform delete
 *   }
 * }
 *
 * return (
 *   <>
 *     <Button onClick={handleDelete}>Eliminar</Button>
 *     <ConfirmDialog {...confirm} />
 *   </>
 * )
 * ```
 */
export function useConfirm(): UseConfirmReturn {
  const [state, setState] = useState<ConfirmState>({
    isOpen: false,
    title: '',
    onConfirm: () => {},
  })

  const [resolveRef, setResolveRef] = useState<{
    resolve: (value: boolean) => void
  } | null>(null)

  const confirm = useCallback((options: ConfirmOptions): Promise<boolean> => {
    return new Promise((resolve) => {
      setResolveRef({ resolve })
      setState({
        isOpen: true,
        ...options,
        onConfirm: () => {
          resolve(true)
          setState((prev) => ({ ...prev, isOpen: false }))
        },
      })
    })
  }, [])

  const onOpenChange = useCallback(
    (open: boolean) => {
      if (!open && resolveRef) {
        resolveRef.resolve(false)
        setResolveRef(null)
      }
      setState((prev) => ({ ...prev, isOpen: open }))
    },
    [resolveRef]
  )

  return {
    isOpen: state.isOpen,
    open: state.isOpen, // Alias for ConfirmDialog compatibility
    title: state.title,
    description: state.description,
    confirmLabel: state.confirmLabel,
    cancelLabel: state.cancelLabel,
    variant: state.variant,
    onConfirm: state.onConfirm,
    onOpenChange,
    confirm,
  }
}
