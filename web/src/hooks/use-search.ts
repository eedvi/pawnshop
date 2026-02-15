import { useSearchParams } from 'react-router-dom'
import { useCallback, useMemo } from 'react'
import { useDebounce } from './use-debounce'
import { SEARCH_DEBOUNCE } from '@/lib/constants'

interface UseSearchOptions {
  paramName?: string
  debounce?: number
}

interface UseSearchReturn {
  search: string
  debouncedSearch: string
  setSearch: (value: string) => void
  clearSearch: () => void
}

/**
 * Hook that syncs search state with URL search params.
 * Includes debounced value for API calls.
 */
export function useSearch({
  paramName = 'search',
  debounce = SEARCH_DEBOUNCE,
}: UseSearchOptions = {}): UseSearchReturn {
  const [searchParams, setSearchParams] = useSearchParams()

  const search = useMemo(() => {
    return searchParams.get(paramName) || ''
  }, [searchParams, paramName])

  const debouncedSearch = useDebounce(search, debounce)

  const setSearch = useCallback(
    (value: string) => {
      setSearchParams((prev) => {
        const params = new URLSearchParams(prev)
        if (value) {
          params.set(paramName, value)
        } else {
          params.delete(paramName)
        }
        // Reset to first page when searching
        params.delete('page')
        return params
      })
    },
    [setSearchParams, paramName]
  )

  const clearSearch = useCallback(() => {
    setSearch('')
  }, [setSearch])

  return {
    search,
    debouncedSearch,
    setSearch,
    clearSearch,
  }
}
