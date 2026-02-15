import {
  ColumnDef,
  SortingState,
  VisibilityState,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from '@tanstack/react-table'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { DataTablePagination } from './data-table-pagination'
import { DataTableToolbar } from './data-table-toolbar'
import { Skeleton } from '@/components/ui/skeleton'

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
  // Pagination
  pageCount?: number
  pageIndex?: number
  pageSize?: number
  onPaginationChange?: (pageIndex: number, pageSize: number) => void
  // Sorting
  sorting?: SortingState
  onSortingChange?: (sorting: SortingState) => void
  // Search
  searchPlaceholder?: string
  searchValue?: string
  onSearchChange?: (value: string) => void
  // Loading
  isLoading?: boolean
  // Toolbar
  showToolbar?: boolean
  toolbarActions?: React.ReactNode
  // Column visibility
  columnVisibility?: VisibilityState
  onColumnVisibilityChange?: (visibility: VisibilityState) => void
  // Empty state
  emptyMessage?: string
}

export function DataTable<TData, TValue>({
  columns,
  data,
  pageCount = 1,
  pageIndex = 0,
  pageSize = 20,
  onPaginationChange,
  sorting = [],
  onSortingChange,
  searchPlaceholder = 'Buscar...',
  searchValue = '',
  onSearchChange,
  isLoading = false,
  showToolbar = true,
  toolbarActions,
  columnVisibility = {},
  onColumnVisibilityChange,
  emptyMessage = 'No se encontraron resultados.',
}: DataTableProps<TData, TValue>) {
  const table = useReactTable({
    data,
    columns,
    pageCount,
    state: {
      sorting,
      columnVisibility,
      pagination: {
        pageIndex,
        pageSize,
      },
    },
    onSortingChange: (updater) => {
      const newSorting = typeof updater === 'function' ? updater(sorting) : updater
      onSortingChange?.(newSorting)
    },
    onColumnVisibilityChange: (updater) => {
      const newVisibility = typeof updater === 'function' ? updater(columnVisibility) : updater
      onColumnVisibilityChange?.(newVisibility)
    },
    onPaginationChange: (updater) => {
      const current = { pageIndex, pageSize }
      const newPagination = typeof updater === 'function' ? updater(current) : updater
      onPaginationChange?.(newPagination.pageIndex, newPagination.pageSize)
    },
    getCoreRowModel: getCoreRowModel(),
    manualPagination: true,
    manualSorting: true,
  })

  return (
    <div className="space-y-4">
      {showToolbar && (
        <DataTableToolbar
          table={table}
          searchPlaceholder={searchPlaceholder}
          searchValue={searchValue}
          onSearchChange={onSearchChange}
          actions={toolbarActions}
        />
      )}

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id} colSpan={header.colSpan}>
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {isLoading ? (
              // Loading skeleton
              Array.from({ length: pageSize }).map((_, index) => (
                <TableRow key={index}>
                  {columns.map((_, colIndex) => (
                    <TableCell key={colIndex}>
                      <Skeleton className="h-6 w-full" />
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && 'selected'}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  {emptyMessage}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      {onPaginationChange && (
        <DataTablePagination
          table={table}
          totalRows={pageCount * pageSize}
        />
      )}
    </div>
  )
}
