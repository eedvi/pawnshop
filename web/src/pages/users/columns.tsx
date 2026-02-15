import { ColumnDef } from '@tanstack/react-table'
import { Link } from 'react-router-dom'
import { MoreHorizontal, Eye, Edit, Key, Lock, Unlock, Power } from 'lucide-react'

import { User } from '@/types'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { DataTableColumnHeader } from '@/components/data-table/data-table-column-header'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { formatDateTime } from '@/lib/format'
import { userRoute, userEditRoute } from '@/routes/routes'

interface ColumnActions {
  onResetPassword: (user: User) => void
  onToggleActive: (user: User) => void
  onUnlock: (user: User) => void
}

export function getUserColumns(actions: ColumnActions): ColumnDef<User>[] {
  return [
    {
      accessorKey: 'avatar',
      header: '',
      cell: ({ row }) => {
        const user = row.original
        const initials = `${user.first_name[0]}${user.last_name[0]}`.toUpperCase()
        return (
          <Avatar className="h-8 w-8">
            <AvatarImage src={user.avatar_url} alt={`${user.first_name} ${user.last_name}`} />
            <AvatarFallback>{initials}</AvatarFallback>
          </Avatar>
        )
      },
    },
    {
      accessorKey: 'name',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Nombre" />,
      cell: ({ row }) => (
        <Link
          to={userRoute(row.original.id)}
          className="font-medium text-primary hover:underline"
        >
          {row.original.first_name} {row.original.last_name}
        </Link>
      ),
    },
    {
      accessorKey: 'email',
      header: ({ column }) => <DataTableColumnHeader column={column} title="Email" />,
      cell: ({ row }) => row.original.email,
    },
    {
      accessorKey: 'role',
      header: 'Rol',
      cell: ({ row }) => row.original.role?.name || '-',
    },
    {
      accessorKey: 'branch',
      header: 'Sucursal',
      cell: ({ row }) => row.original.branch?.name || 'Todas',
    },
    {
      accessorKey: 'is_active',
      header: 'Estado',
      cell: ({ row }) => {
        const isLocked = row.original.locked_until && new Date(row.original.locked_until) > new Date()
        if (isLocked) {
          return (
            <Badge variant="destructive" className="flex items-center gap-1 w-fit">
              <Lock className="h-3 w-3" />
              Bloqueado
            </Badge>
          )
        }
        return (
          <Badge variant={row.original.is_active ? 'default' : 'secondary'}>
            {row.original.is_active ? 'Activo' : 'Inactivo'}
          </Badge>
        )
      },
    },
    {
      accessorKey: 'last_login_at',
      header: 'Último acceso',
      cell: ({ row }) =>
        row.original.last_login_at ? formatDateTime(row.original.last_login_at) : 'Nunca',
    },
    {
      id: 'actions',
      cell: ({ row }) => {
        const user = row.original
        const isLocked = user.locked_until && new Date(user.locked_until) > new Date()

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
                <Link to={userRoute(user.id)}>
                  <Eye className="mr-2 h-4 w-4" />
                  Ver detalles
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link to={userEditRoute(user.id)}>
                  <Edit className="mr-2 h-4 w-4" />
                  Editar
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => actions.onResetPassword(user)}>
                <Key className="mr-2 h-4 w-4" />
                Cambiar contraseña
              </DropdownMenuItem>
              {isLocked && (
                <DropdownMenuItem onClick={() => actions.onUnlock(user)}>
                  <Unlock className="mr-2 h-4 w-4" />
                  Desbloquear
                </DropdownMenuItem>
              )}
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => actions.onToggleActive(user)}
                className={user.is_active ? 'text-destructive' : ''}
              >
                <Power className="mr-2 h-4 w-4" />
                {user.is_active ? 'Desactivar' : 'Activar'}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        )
      },
    },
  ]
}
