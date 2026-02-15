import { useState } from 'react'
import {
  Plus,
  Play,
  Square,
  Loader2,
  Banknote,
  ArrowUpCircle,
  ArrowDownCircle,
  Settings,
  Trash2,
} from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  useCashRegisters,
  useCashSessions,
  useCashMovements,
  useCashSessionSummary,
  useCurrentCashSession,
  useCreateCashRegister,
  useUpdateCashRegister,
  useDeleteCashRegister,
  useOpenCashSession,
  useCloseCashSession,
  useCreateCashMovement,
} from '@/hooks/use-cash'
import { useBranchStore } from '@/stores/branch-store'
import {
  CashRegister,
  CashSession,
  CASH_SESSION_STATUSES,
  CASH_MOVEMENT_TYPES,
} from '@/types'
import { formatCurrency, formatDateTime } from '@/lib/format'
import { RegisterFormDialog } from './register-form-dialog'
import { OpenSessionDialog } from './open-session-dialog'
import { CloseSessionDialog } from './close-session-dialog'
import { AddMovementDialog } from './add-movement-dialog'
import { CashRegisterFormValues, OpenSessionFormValues, CloseSessionFormValues, CashMovementFormValues } from './schemas'

export default function CashPage() {
  const { selectedBranchId } = useBranchStore()

  // State for dialogs
  const [registerDialogOpen, setRegisterDialogOpen] = useState(false)
  const [selectedRegister, setSelectedRegister] = useState<CashRegister | null>(null)
  const [openSessionDialogOpen, setOpenSessionDialogOpen] = useState(false)
  const [closeSessionDialogOpen, setCloseSessionDialogOpen] = useState(false)
  const [selectedSession, setSelectedSession] = useState<CashSession | null>(null)
  const [movementDialogOpen, setMovementDialogOpen] = useState(false)
  const [activeSessionId, setActiveSessionId] = useState<number>(0)

  // Queries
  const { data: registers, isLoading: loadingRegisters } = useCashRegisters({
    branch_id: selectedBranchId ?? undefined,
  })
  const { data: sessionsData, isLoading: loadingSessions } = useCashSessions({
    branch_id: selectedBranchId ?? undefined,
    per_page: 20,
  })
  const { data: movements, isLoading: loadingMovements } = useCashMovements(activeSessionId)
  const { data: sessionSummary } = useCashSessionSummary(selectedSession?.id || 0)

  // Mutations
  const createRegisterMutation = useCreateCashRegister()
  const updateRegisterMutation = useUpdateCashRegister()
  const deleteRegisterMutation = useDeleteCashRegister()
  const openSessionMutation = useOpenCashSession()
  const closeSessionMutation = useCloseCashSession()
  const createMovementMutation = useCreateCashMovement()

  // Handlers
  const handleCreateRegister = () => {
    setSelectedRegister(null)
    setRegisterDialogOpen(true)
  }

  const handleEditRegister = (register: CashRegister) => {
    setSelectedRegister(register)
    setRegisterDialogOpen(true)
  }

  const handleRegisterSubmit = (data: CashRegisterFormValues) => {
    if (selectedRegister) {
      updateRegisterMutation.mutate(
        { id: selectedRegister.id, input: data },
        {
          onSuccess: () => setRegisterDialogOpen(false),
        }
      )
    } else if (selectedBranchId) {
      createRegisterMutation.mutate(
        { ...data, branch_id: selectedBranchId },
        {
          onSuccess: () => setRegisterDialogOpen(false),
        }
      )
    }
  }

  const handleDeleteRegister = (id: number) => {
    if (confirm('¿Está seguro de eliminar esta caja?')) {
      deleteRegisterMutation.mutate(id)
    }
  }

  const handleOpenSession = (register: CashRegister) => {
    setSelectedRegister(register)
    setOpenSessionDialogOpen(true)
  }

  const handleOpenSessionSubmit = (data: OpenSessionFormValues) => {
    openSessionMutation.mutate(data, {
      onSuccess: () => setOpenSessionDialogOpen(false),
    })
  }

  const handleCloseSession = (session: CashSession) => {
    setSelectedSession(session)
    setCloseSessionDialogOpen(true)
  }

  const handleCloseSessionSubmit = (data: CloseSessionFormValues) => {
    if (selectedSession) {
      closeSessionMutation.mutate(
        { id: selectedSession.id, input: data },
        {
          onSuccess: () => {
            setCloseSessionDialogOpen(false)
            setSelectedSession(null)
          },
        }
      )
    }
  }

  const handleViewMovements = (session: CashSession) => {
    setActiveSessionId(session.id)
    setSelectedSession(session)
  }

  const handleAddMovement = () => {
    setMovementDialogOpen(true)
  }

  const handleMovementSubmit = (data: CashMovementFormValues) => {
    if (activeSessionId) {
      createMovementMutation.mutate(
        { ...data, session_id: activeSessionId },
        {
          onSuccess: () => setMovementDialogOpen(false),
        }
      )
    }
  }

  const sessions = sessionsData?.data ?? []

  return (
    <div>
      <PageHeader
        title="Caja"
        description="Gestión de cajas registradoras y sesiones"
      />

      <Tabs defaultValue="registers" className="space-y-4">
        <TabsList>
          <TabsTrigger value="registers">Cajas</TabsTrigger>
          <TabsTrigger value="sessions">Sesiones</TabsTrigger>
          <TabsTrigger value="movements">Movimientos</TabsTrigger>
        </TabsList>

        {/* Registers Tab */}
        <TabsContent value="registers">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle>Cajas Registradoras</CardTitle>
                <CardDescription>Administre las cajas de la sucursal</CardDescription>
              </div>
              <Button onClick={handleCreateRegister}>
                <Plus className="mr-2 h-4 w-4" />
                Nueva Caja
              </Button>
            </CardHeader>
            <CardContent>
              {loadingRegisters ? (
                <div className="flex justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                </div>
              ) : !registers?.length ? (
                <div className="text-center py-8 text-muted-foreground">
                  No hay cajas registradas
                </div>
              ) : (
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                  {registers.map((register) => (
                    <Card key={register.id}>
                      <CardHeader className="pb-2">
                        <div className="flex items-center justify-between">
                          <CardTitle className="text-base">{register.name}</CardTitle>
                          <Badge variant={register.is_active ? 'default' : 'secondary'}>
                            {register.is_active ? 'Activa' : 'Inactiva'}
                          </Badge>
                        </div>
                        <CardDescription className="font-mono">{register.code}</CardDescription>
                      </CardHeader>
                      <CardContent>
                        {register.description && (
                          <p className="text-sm text-muted-foreground mb-4">
                            {register.description}
                          </p>
                        )}
                        <div className="flex gap-2">
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => handleOpenSession(register)}
                          >
                            <Play className="mr-1 h-3 w-3" />
                            Abrir
                          </Button>
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => handleEditRegister(register)}
                          >
                            <Settings className="h-4 w-4" />
                          </Button>
                          <Button
                            size="sm"
                            variant="ghost"
                            className="text-destructive"
                            onClick={() => handleDeleteRegister(register.id)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Sessions Tab */}
        <TabsContent value="sessions">
          <Card>
            <CardHeader>
              <CardTitle>Sesiones de Caja</CardTitle>
              <CardDescription>Historial de sesiones abiertas y cerradas</CardDescription>
            </CardHeader>
            <CardContent>
              {loadingSessions ? (
                <div className="flex justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                </div>
              ) : !sessions.length ? (
                <div className="text-center py-8 text-muted-foreground">
                  No hay sesiones registradas
                </div>
              ) : (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Caja</TableHead>
                      <TableHead>Usuario</TableHead>
                      <TableHead>Apertura</TableHead>
                      <TableHead>Cierre</TableHead>
                      <TableHead>Monto Inicial</TableHead>
                      <TableHead>Monto Cierre</TableHead>
                      <TableHead>Estado</TableHead>
                      <TableHead>Acciones</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {sessions.map((session) => {
                      const status = CASH_SESSION_STATUSES.find(
                        (s) => s.value === session.status
                      )
                      return (
                        <TableRow key={session.id}>
                          <TableCell className="font-medium">
                            {session.register?.name || '-'}
                          </TableCell>
                          <TableCell>
                            {session.user
                              ? `${session.user.first_name} ${session.user.last_name}`
                              : '-'}
                          </TableCell>
                          <TableCell>{formatDateTime(session.opened_at)}</TableCell>
                          <TableCell>
                            {session.closed_at ? formatDateTime(session.closed_at) : '-'}
                          </TableCell>
                          <TableCell>{formatCurrency(session.opening_amount)}</TableCell>
                          <TableCell>
                            {session.closing_amount !== undefined
                              ? formatCurrency(session.closing_amount)
                              : '-'}
                          </TableCell>
                          <TableCell>
                            <Badge
                              variant={
                                status?.color === 'green' ? 'default' : 'secondary'
                              }
                            >
                              {status?.label || session.status}
                            </Badge>
                          </TableCell>
                          <TableCell>
                            <div className="flex gap-1">
                              <Button
                                size="sm"
                                variant="ghost"
                                onClick={() => handleViewMovements(session)}
                              >
                                <Banknote className="h-4 w-4" />
                              </Button>
                              {session.status === 'open' && (
                                <Button
                                  size="sm"
                                  variant="ghost"
                                  onClick={() => handleCloseSession(session)}
                                >
                                  <Square className="h-4 w-4" />
                                </Button>
                              )}
                            </div>
                          </TableCell>
                        </TableRow>
                      )
                    })}
                  </TableBody>
                </Table>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Movements Tab */}
        <TabsContent value="movements">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle>Movimientos de Caja</CardTitle>
                <CardDescription>
                  {selectedSession
                    ? `Sesión: ${selectedSession.register?.name || 'Caja'} - ${formatDateTime(selectedSession.opened_at)}`
                    : 'Seleccione una sesión para ver sus movimientos'}
                </CardDescription>
              </div>
              {selectedSession?.status === 'open' && (
                <Button onClick={handleAddMovement}>
                  <Plus className="mr-2 h-4 w-4" />
                  Agregar Movimiento
                </Button>
              )}
            </CardHeader>
            <CardContent>
              {!activeSessionId ? (
                <div className="text-center py-8 text-muted-foreground">
                  Seleccione una sesión en la pestaña "Sesiones" para ver sus movimientos
                </div>
              ) : loadingMovements ? (
                <div className="flex justify-center py-8">
                  <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                </div>
              ) : !movements?.length ? (
                <div className="text-center py-8 text-muted-foreground">
                  No hay movimientos en esta sesión
                </div>
              ) : (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Tipo</TableHead>
                      <TableHead>Descripción</TableHead>
                      <TableHead>Monto</TableHead>
                      <TableHead>Fecha</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {movements.map((movement) => {
                      const type = CASH_MOVEMENT_TYPES.find(
                        (t) => t.value === movement.movement_type
                      )
                      const isIncome = movement.movement_type === 'income'
                      const isExpense = movement.movement_type === 'expense'
                      return (
                        <TableRow key={movement.id}>
                          <TableCell>
                            <div className="flex items-center gap-2">
                              {isIncome && (
                                <ArrowUpCircle className="h-4 w-4 text-green-500" />
                              )}
                              {isExpense && (
                                <ArrowDownCircle className="h-4 w-4 text-red-500" />
                              )}
                              {!isIncome && !isExpense && (
                                <Banknote className="h-4 w-4 text-yellow-500" />
                              )}
                              <Badge
                                variant={
                                  isIncome
                                    ? 'default'
                                    : isExpense
                                    ? 'destructive'
                                    : 'secondary'
                                }
                              >
                                {type?.label || movement.movement_type}
                              </Badge>
                            </div>
                          </TableCell>
                          <TableCell>{movement.description}</TableCell>
                          <TableCell
                            className={
                              isIncome
                                ? 'text-green-600 font-medium'
                                : isExpense
                                ? 'text-red-600 font-medium'
                                : 'font-medium'
                            }
                          >
                            {isIncome ? '+' : isExpense ? '-' : ''}
                            {formatCurrency(movement.amount)}
                          </TableCell>
                          <TableCell>{formatDateTime(movement.created_at)}</TableCell>
                        </TableRow>
                      )
                    })}
                  </TableBody>
                </Table>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      {/* Dialogs */}
      <RegisterFormDialog
        open={registerDialogOpen}
        onOpenChange={setRegisterDialogOpen}
        register={selectedRegister}
        onConfirm={handleRegisterSubmit}
        isLoading={createRegisterMutation.isPending || updateRegisterMutation.isPending}
      />

      <OpenSessionDialog
        open={openSessionDialogOpen}
        onOpenChange={setOpenSessionDialogOpen}
        register={selectedRegister}
        onConfirm={handleOpenSessionSubmit}
        isLoading={openSessionMutation.isPending}
      />

      <CloseSessionDialog
        open={closeSessionDialogOpen}
        onOpenChange={setCloseSessionDialogOpen}
        session={selectedSession}
        summary={sessionSummary || null}
        onConfirm={handleCloseSessionSubmit}
        isLoading={closeSessionMutation.isPending}
      />

      <AddMovementDialog
        open={movementDialogOpen}
        onOpenChange={setMovementDialogOpen}
        sessionId={activeSessionId}
        onConfirm={handleMovementSubmit}
        isLoading={createMovementMutation.isPending}
      />
    </div>
  )
}
