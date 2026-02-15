import { useSearchParams } from 'react-router-dom'
import { useCallback, useMemo } from 'react'
import { DEFAULT_PAGE, DEFAULT_PER_PAGE } from '@/lib/constants'

interface UsePaginationOptions {
  defaultPage?: number
  defaultPerPage?: number
}

interface UsePaginationReturn {
  page: number
  perPage: number
  setPage: (page: number) => void
  setPerPage: (perPage: number) => void
  offset: number
  // For TanStack Table
  pageIndex: number
  pageSize: number
  onPaginationChange: (pageIndex: number, pageSize: number) => void
}

/**
 * Hook that syncs pagination state with URL search params.
 * Makes pages bookmarkable and shareable.
 */
export function usePagination({
  defaultPage = DEFAULT_PAGE,
  defaultPerPage = DEFAULT_PER_PAGE,
}: UsePaginationOptions = {}): UsePaginationReturn {
  const [searchParams, setSearchParams] = useSearchParams()

  const page = useMemo(() => {
    const pageParam = searchParams.get('page')
    const parsed = pageParam ? parseInt(pageParam, 10) : defaultPage
    return isNaN(parsed) || parsed < 1 ? defaultPage : parsed
  }, [searchParams, defaultPage])

  const perPage = useMemo(() => {
    const perPageParam = searchParams.get('per_page')
    const parsed = perPageParam ? parseInt(perPageParam, 10) : defaultPerPage
    return isNaN(parsed) || parsed < 1 ? defaultPerPage : parsed
  }, [searchParams, defaultPerPage])

  const setPage = useCallback(
    (newPage: number) => {
      setSearchParams((prev) => {
        const params = new URLSearchParams(prev)
        if (newPage === defaultPage) {
          params.delete('page')
        } else {
          params.set('page', newPage.toString())
        }
        return params
      })
    },
    [setSearchParams, defaultPage]
  )

  const setPerPage = useCallback(
    (newPerPage: number) => {
      setSearchParams((prev) => {
        const params = new URLSearchParams(prev)
        if (newPerPage === defaultPerPage) {
          params.delete('per_page')
        } else {
          params.set('per_page', newPerPage.toString())
        }
        // Reset to first page when changing per_page
        params.delete('page')
        return params
      })
    },
    [setSearchParams, defaultPerPage]
  )

  // For compatibility with TanStack Table (0-indexed)
  const pageIndex = page - 1
  const pageSize = perPage

  const onPaginationChange = useCallback(
    (newPageIndex: number, newPageSize: number) => {
      setSearchParams((prev) => {
        const params = new URLSearchParams(prev)

        const newPage = newPageIndex + 1
        if (newPage === defaultPage) {
          params.delete('page')
        } else {
          params.set('page', newPage.toString())
        }

        if (newPageSize === defaultPerPage) {
          params.delete('per_page')
        } else {
          params.set('per_page', newPageSize.toString())
        }

        return params
      })
    },
    [setSearchParams, defaultPage, defaultPerPage]
  )

  return {
    page,
    perPage,
    setPage,
    setPerPage,
    offset: (page - 1) * perPage,
    // TanStack Table compatibility
    pageIndex,
    pageSize,
    onPaginationChange,
  }
}
