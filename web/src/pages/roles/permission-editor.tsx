import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { PERMISSION_GROUPS } from '@/types/role'
import {
  Users,
  Package,
  HandCoins,
  CreditCard,
  ShoppingCart,
  Wallet,
  BarChart3,
  UserCog,
  Building2,
  FolderTree,
  Shield,
  Settings,
  FileSearch,
  Bell,
  Receipt,
  ArrowRightLeft,
} from 'lucide-react'

interface PermissionEditorProps {
  value: string[]
  onChange: (permissions: string[]) => void
  disabled?: boolean
}

const GROUP_CONFIG: Record<string, { label: string; icon: React.ElementType }> = {
  customers: { label: 'Clientes', icon: Users },
  items: { label: 'Artículos', icon: Package },
  loans: { label: 'Préstamos', icon: HandCoins },
  payments: { label: 'Pagos', icon: CreditCard },
  sales: { label: 'Ventas', icon: ShoppingCart },
  cash: { label: 'Caja', icon: Wallet },
  reports: { label: 'Reportes', icon: BarChart3 },
  users: { label: 'Usuarios', icon: UserCog },
  branches: { label: 'Sucursales', icon: Building2 },
  categories: { label: 'Categorías', icon: FolderTree },
  roles: { label: 'Roles', icon: Shield },
  settings: { label: 'Configuración', icon: Settings },
  audit: { label: 'Auditoría', icon: FileSearch },
  notifications: { label: 'Notificaciones', icon: Bell },
  expenses: { label: 'Gastos', icon: Receipt },
  transfers: { label: 'Transferencias', icon: ArrowRightLeft },
}

const PERMISSION_LABELS: Record<string, string> = {
  read: 'Ver',
  create: 'Crear',
  update: 'Editar',
  delete: 'Eliminar',
  export: 'Exportar',
  approve: 'Aprobar',
  ship: 'Enviar',
  receive: 'Recibir',
  cancel: 'Cancelar',
  manage: 'Gestionar',
}

function getPermissionLabel(permission: string): string {
  const action = permission.split('.').pop() || permission
  return PERMISSION_LABELS[action] || action
}

export function PermissionEditor({ value, onChange, disabled }: PermissionEditorProps) {
  const handlePermissionChange = (permission: string, checked: boolean) => {
    if (checked) {
      onChange([...value, permission])
    } else {
      onChange(value.filter((p) => p !== permission))
    }
  }

  const handleGroupSelectAll = (groupKey: string, permissions: string[]) => {
    const allSelected = permissions.every((p) => value.includes(p))
    if (allSelected) {
      onChange(value.filter((p) => !permissions.includes(p)))
    } else {
      const newPermissions = [...value]
      permissions.forEach((p) => {
        if (!newPermissions.includes(p)) {
          newPermissions.push(p)
        }
      })
      onChange(newPermissions)
    }
  }

  const handleSelectAll = () => {
    const allPermissions = Object.values(PERMISSION_GROUPS).flat()
    const allSelected = allPermissions.every((p) => value.includes(p))
    if (allSelected) {
      onChange([])
    } else {
      onChange(allPermissions)
    }
  }

  const allPermissions = Object.values(PERMISSION_GROUPS).flat()
  const allSelected = allPermissions.every((p) => value.includes(p))
  const someSelected = value.length > 0 && !allSelected

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Checkbox
            id="select-all"
            checked={allSelected}
            onCheckedChange={handleSelectAll}
            disabled={disabled}
            className={someSelected ? 'data-[state=checked]:bg-primary/50' : ''}
          />
          <Label htmlFor="select-all" className="font-medium">
            Seleccionar todos
          </Label>
        </div>
        <Badge variant="secondary">
          {value.length} / {allPermissions.length} permisos
        </Badge>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {Object.entries(PERMISSION_GROUPS).map(([groupKey, permissions]) => {
          const config = GROUP_CONFIG[groupKey] || { label: groupKey, icon: Shield }
          const Icon = config.icon
          const groupAllSelected = permissions.every((p) => value.includes(p))
          const groupSomeSelected = permissions.some((p) => value.includes(p)) && !groupAllSelected

          return (
            <Card key={groupKey} className="overflow-hidden">
              <CardHeader className="py-3 px-4 bg-muted/50">
                <div className="flex items-center justify-between">
                  <CardTitle className="text-sm flex items-center gap-2">
                    <Icon className="h-4 w-4" />
                    {config.label}
                  </CardTitle>
                  <Checkbox
                    checked={groupAllSelected}
                    onCheckedChange={() => handleGroupSelectAll(groupKey, permissions)}
                    disabled={disabled}
                    className={groupSomeSelected ? 'data-[state=checked]:bg-primary/50' : ''}
                  />
                </div>
              </CardHeader>
              <CardContent className="py-3 px-4">
                <div className="space-y-2">
                  {permissions.map((permission) => (
                    <div key={permission} className="flex items-center gap-2">
                      <Checkbox
                        id={permission}
                        checked={value.includes(permission)}
                        onCheckedChange={(checked) =>
                          handlePermissionChange(permission, checked as boolean)
                        }
                        disabled={disabled}
                      />
                      <Label
                        htmlFor={permission}
                        className="text-sm font-normal cursor-pointer"
                      >
                        {getPermissionLabel(permission)}
                      </Label>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )
        })}
      </div>
    </div>
  )
}
