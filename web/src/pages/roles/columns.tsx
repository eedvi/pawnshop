import { ColumnDef } from '@tanstack/react-table'
import { Link } from 'react-router-dom'
import { MoreHorizontal, Eye, Edit, Trash2, Shield } from 'lucide-react'

import { Role } from '@/types'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Badge } from '@/components/ui/badge'
import { roleRoute, roleEditRoute } from '@/routes/routes'

interface ColumnActions {
  onDelete: (role: Role) => void
}

export function getRoleColumns(actions: ColumnActions): ColumnDef<Role>[] {
  return [
    {
      accessorKey: 'display_name',
      header: 'Nombre',
      cell: ({ row }) => (
        <Link
          to={roleRoute(row.original.id)}
          className="font-medium text-primary hover:underline flex items-center gap-2"
        >
          <Shield className="h-4 w-4" />
          {row.original.display_name}
        </Link>
      ),
    },
    {
      accessorKey: 'name',
      header: 'Código',
      cell: ({ row }) => (
        <span className="font-mono text-muted-foreground">{row.original.name}</span>
      ),
    },
    {
      accessorKey: 'description',
      header: 'Descripción',
      cell: ({ row }) => row.original.description || '-',
    },
    {
      accessorKey: 'permissions',
      header: 'Permisos',
      cell: ({ row }) => (
        <Badge variant="secondary">{row.original.permissions.length} permisos</Badge>
      ),
    },
    {
      accessorKey: 'is_system',
      header: 'Tipo',
      cell: ({ row }) => (
        <Badge variant={row.original.is_system ? 'outline' : 'default'}>
          {row.original.is_system ? 'Sistema' : 'Personalizado'}
        </Badge>
      ),
    },
    {
      id: 'actions',
      cell: ({ row }) => {
        const role = row.original

        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Abrir menú</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem asChild>
                <Link to={roleRoute(role.id)}>
                  <Eye className="mr-2 h-4 w-4" />
                  Ver detalles
                </Link>
              </DropdownMenuItem>
              {!role.is_system && (
                <>
                  <DropdownMenuItem asChild>
                    <Link to={roleEditRoute(role.id)}>
                      <Edit className="mr-2 h-4 w-4" />
                      Editar
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    onClick={() => actions.onDelete(role)}
                    className="text-destructive focus:text-destructive"
                  >
                    <Trash2 className="mr-2 h-4 w-4" />
                    Eliminar
                  </DropdownMenuItem>
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
    },
  ]
}
