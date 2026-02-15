import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Loader2, Building2, Settings2, User, Save } from 'lucide-react'

import { PageHeader } from '@/components/layout/page-header'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useSettings, useSetMultipleSettings } from '@/hooks/use-settings'
import { useBranchStore } from '@/stores/branch-store'
import { useAuthStore } from '@/stores/auth-store'
import { SETTING_KEYS } from '@/types/settings'
import { toast } from 'sonner'

// Schemas
const companySettingsSchema = z.object({
  companyName: z.string().min(1, 'El nombre es requerido'),
  companyAddress: z.string().optional(),
  companyPhone: z.string().optional(),
  companyEmail: z.string().email('Email inválido').optional().or(z.literal('')),
  companyTaxId: z.string().optional(),
})

const loanSettingsSchema = z.object({
  defaultInterestRate: z.coerce.number().min(0).max(100),
  defaultTermDays: z.coerce.number().min(1).max(365),
  defaultGracePeriod: z.coerce.number().min(0).max(30),
  lateFeeRate: z.coerce.number().min(0).max(100),
  minLoanAmount: z.coerce.number().min(0),
  maxLoanAmount: z.coerce.number().min(0),
})

const systemSettingsSchema = z.object({
  timezone: z.string(),
  currency: z.string(),
  dateFormat: z.string(),
  allowNegativeCash: z.boolean(),
})

const notificationSettingsSchema = z.object({
  emailEnabled: z.boolean(),
  smsEnabled: z.boolean(),
  whatsappEnabled: z.boolean(),
  reminderDaysBefore: z.coerce.number().min(1).max(30),
})

type CompanySettings = z.infer<typeof companySettingsSchema>
type LoanSettings = z.infer<typeof loanSettingsSchema>
type SystemSettings = z.infer<typeof systemSettingsSchema>
type NotificationSettings = z.infer<typeof notificationSettingsSchema>

// Helper to get setting value
function getSettingValue(settings: { key: string; value: string }[] | undefined, key: string, defaultValue = ''): string {
  const setting = settings?.find((s) => s.key === key)
  return setting?.value || defaultValue
}

export default function SettingsPage() {
  const [activeTab, setActiveTab] = useState('company')
  const { selectedBranch } = useBranchStore()
  const { user } = useAuthStore()

  const { data: settings, isLoading } = useSettings()
  const setMultipleMutation = useSetMultipleSettings()

  // Company settings form
  const companyForm = useForm<CompanySettings>({
    resolver: zodResolver(companySettingsSchema),
    values: {
      companyName: getSettingValue(settings, SETTING_KEYS.COMPANY_NAME, 'Mi Empeño'),
      companyAddress: getSettingValue(settings, SETTING_KEYS.COMPANY_ADDRESS),
      companyPhone: getSettingValue(settings, SETTING_KEYS.COMPANY_PHONE),
      companyEmail: getSettingValue(settings, SETTING_KEYS.COMPANY_EMAIL),
      companyTaxId: getSettingValue(settings, SETTING_KEYS.COMPANY_TAX_ID),
    },
  })

  // Loan settings form
  const loanForm = useForm<LoanSettings>({
    resolver: zodResolver(loanSettingsSchema),
    values: {
      defaultInterestRate: parseFloat(getSettingValue(settings, SETTING_KEYS.DEFAULT_INTEREST_RATE, '5')),
      defaultTermDays: parseInt(getSettingValue(settings, SETTING_KEYS.DEFAULT_LOAN_TERM_DAYS, '30')),
      defaultGracePeriod: parseInt(getSettingValue(settings, SETTING_KEYS.DEFAULT_GRACE_PERIOD, '3')),
      lateFeeRate: parseFloat(getSettingValue(settings, SETTING_KEYS.LATE_FEE_RATE, '2')),
      minLoanAmount: parseFloat(getSettingValue(settings, SETTING_KEYS.MIN_LOAN_AMOUNT, '100')),
      maxLoanAmount: parseFloat(getSettingValue(settings, SETTING_KEYS.MAX_LOAN_AMOUNT, '100000')),
    },
  })

  // System settings form
  const systemForm = useForm<SystemSettings>({
    resolver: zodResolver(systemSettingsSchema),
    values: {
      timezone: getSettingValue(settings, SETTING_KEYS.TIMEZONE, 'America/Mexico_City'),
      currency: getSettingValue(settings, SETTING_KEYS.CURRENCY, 'MXN'),
      dateFormat: getSettingValue(settings, SETTING_KEYS.DATE_FORMAT, 'dd/MM/yyyy'),
      allowNegativeCash: getSettingValue(settings, SETTING_KEYS.ALLOW_NEGATIVE_CASH, 'false') === 'true',
    },
  })

  // Notification settings form
  const notificationForm = useForm<NotificationSettings>({
    resolver: zodResolver(notificationSettingsSchema),
    values: {
      emailEnabled: getSettingValue(settings, SETTING_KEYS.EMAIL_ENABLED, 'true') === 'true',
      smsEnabled: getSettingValue(settings, SETTING_KEYS.SMS_ENABLED, 'false') === 'true',
      whatsappEnabled: getSettingValue(settings, SETTING_KEYS.WHATSAPP_ENABLED, 'false') === 'true',
      reminderDaysBefore: parseInt(getSettingValue(settings, SETTING_KEYS.REMINDER_DAYS_BEFORE, '3')),
    },
  })

  const handleCompanySave = (values: CompanySettings) => {
    setMultipleMutation.mutate(
      {
        settings: [
          { key: SETTING_KEYS.COMPANY_NAME, value: values.companyName, data_type: 'string' },
          { key: SETTING_KEYS.COMPANY_ADDRESS, value: values.companyAddress || '', data_type: 'string' },
          { key: SETTING_KEYS.COMPANY_PHONE, value: values.companyPhone || '', data_type: 'string' },
          { key: SETTING_KEYS.COMPANY_EMAIL, value: values.companyEmail || '', data_type: 'string' },
          { key: SETTING_KEYS.COMPANY_TAX_ID, value: values.companyTaxId || '', data_type: 'string' },
        ],
      },
      {
        onSuccess: () => toast.success('Configuración de empresa guardada'),
      }
    )
  }

  const handleLoanSave = (values: LoanSettings) => {
    setMultipleMutation.mutate(
      {
        settings: [
          { key: SETTING_KEYS.DEFAULT_INTEREST_RATE, value: values.defaultInterestRate.toString(), data_type: 'number' },
          { key: SETTING_KEYS.DEFAULT_LOAN_TERM_DAYS, value: values.defaultTermDays.toString(), data_type: 'number' },
          { key: SETTING_KEYS.DEFAULT_GRACE_PERIOD, value: values.defaultGracePeriod.toString(), data_type: 'number' },
          { key: SETTING_KEYS.LATE_FEE_RATE, value: values.lateFeeRate.toString(), data_type: 'number' },
          { key: SETTING_KEYS.MIN_LOAN_AMOUNT, value: values.minLoanAmount.toString(), data_type: 'number' },
          { key: SETTING_KEYS.MAX_LOAN_AMOUNT, value: values.maxLoanAmount.toString(), data_type: 'number' },
        ],
      },
      {
        onSuccess: () => toast.success('Configuración de préstamos guardada'),
      }
    )
  }

  const handleSystemSave = (values: SystemSettings) => {
    setMultipleMutation.mutate(
      {
        settings: [
          { key: SETTING_KEYS.TIMEZONE, value: values.timezone, data_type: 'string' },
          { key: SETTING_KEYS.CURRENCY, value: values.currency, data_type: 'string' },
          { key: SETTING_KEYS.DATE_FORMAT, value: values.dateFormat, data_type: 'string' },
          { key: SETTING_KEYS.ALLOW_NEGATIVE_CASH, value: values.allowNegativeCash.toString(), data_type: 'boolean' },
        ],
      },
      {
        onSuccess: () => toast.success('Configuración del sistema guardada'),
      }
    )
  }

  const handleNotificationSave = (values: NotificationSettings) => {
    setMultipleMutation.mutate(
      {
        settings: [
          { key: SETTING_KEYS.EMAIL_ENABLED, value: values.emailEnabled.toString(), data_type: 'boolean' },
          { key: SETTING_KEYS.SMS_ENABLED, value: values.smsEnabled.toString(), data_type: 'boolean' },
          { key: SETTING_KEYS.WHATSAPP_ENABLED, value: values.whatsappEnabled.toString(), data_type: 'boolean' },
          { key: SETTING_KEYS.REMINDER_DAYS_BEFORE, value: values.reminderDaysBefore.toString(), data_type: 'number' },
        ],
      },
      {
        onSuccess: () => toast.success('Configuración de notificaciones guardada'),
      }
    )
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  return (
    <div>
      <PageHeader
        title="Configuración"
        description="Configuración general del sistema"
      />

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="mb-6">
          <TabsTrigger value="company" className="gap-2">
            <Building2 className="h-4 w-4" />
            Empresa
          </TabsTrigger>
          <TabsTrigger value="loans" className="gap-2">
            <Settings2 className="h-4 w-4" />
            Préstamos
          </TabsTrigger>
          <TabsTrigger value="system" className="gap-2">
            <Settings2 className="h-4 w-4" />
            Sistema
          </TabsTrigger>
          <TabsTrigger value="notifications" className="gap-2">
            <Settings2 className="h-4 w-4" />
            Notificaciones
          </TabsTrigger>
          <TabsTrigger value="account" className="gap-2">
            <User className="h-4 w-4" />
            Mi Cuenta
          </TabsTrigger>
        </TabsList>

        {/* Company Settings */}
        <TabsContent value="company">
          <Card>
            <CardHeader>
              <CardTitle>Información de la Empresa</CardTitle>
              <CardDescription>
                Datos generales de la empresa que aparecerán en tickets y documentos
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Form {...companyForm}>
                <form onSubmit={companyForm.handleSubmit(handleCompanySave)} className="space-y-4">
                  <FormField
                    control={companyForm.control}
                    name="companyName"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Nombre de la Empresa</FormLabel>
                        <FormControl>
                          <Input placeholder="Mi Empeño S.A." {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={companyForm.control}
                    name="companyAddress"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Dirección</FormLabel>
                        <FormControl>
                          <Input placeholder="Calle Principal #123, Colonia Centro" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className="grid gap-4 sm:grid-cols-2">
                    <FormField
                      control={companyForm.control}
                      name="companyPhone"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Teléfono</FormLabel>
                          <FormControl>
                            <Input placeholder="+52 55 1234 5678" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={companyForm.control}
                      name="companyEmail"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Correo Electrónico</FormLabel>
                          <FormControl>
                            <Input type="email" placeholder="contacto@miempeno.com" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  <FormField
                    control={companyForm.control}
                    name="companyTaxId"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>RFC / ID Fiscal</FormLabel>
                        <FormControl>
                          <Input placeholder="XAXX010101000" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className="flex justify-end">
                    <Button type="submit" disabled={setMultipleMutation.isPending}>
                      <Save className="mr-2 h-4 w-4" />
                      Guardar Cambios
                    </Button>
                  </div>
                </form>
              </Form>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Loan Settings */}
        <TabsContent value="loans">
          <Card>
            <CardHeader>
              <CardTitle>Configuración de Préstamos</CardTitle>
              <CardDescription>
                Valores predeterminados para nuevos préstamos
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Form {...loanForm}>
                <form onSubmit={loanForm.handleSubmit(handleLoanSave)} className="space-y-4">
                  <div className="grid gap-4 sm:grid-cols-2">
                    <FormField
                      control={loanForm.control}
                      name="defaultInterestRate"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Tasa de Interés Mensual (%)</FormLabel>
                          <FormControl>
                            <Input type="number" step="0.1" {...field} />
                          </FormControl>
                          <FormDescription>
                            Porcentaje mensual aplicado sobre el capital
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={loanForm.control}
                      name="defaultTermDays"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Plazo Predeterminado (días)</FormLabel>
                          <FormControl>
                            <Input type="number" {...field} />
                          </FormControl>
                          <FormDescription>
                            Duración del préstamo en días
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  <div className="grid gap-4 sm:grid-cols-2">
                    <FormField
                      control={loanForm.control}
                      name="defaultGracePeriod"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Período de Gracia (días)</FormLabel>
                          <FormControl>
                            <Input type="number" {...field} />
                          </FormControl>
                          <FormDescription>
                            Días después del vencimiento antes de aplicar mora
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={loanForm.control}
                      name="lateFeeRate"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Tasa de Mora (%)</FormLabel>
                          <FormControl>
                            <Input type="number" step="0.1" {...field} />
                          </FormControl>
                          <FormDescription>
                            Porcentaje adicional por pagos atrasados
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  <div className="grid gap-4 sm:grid-cols-2">
                    <FormField
                      control={loanForm.control}
                      name="minLoanAmount"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Monto Mínimo de Préstamo</FormLabel>
                          <FormControl>
                            <Input type="number" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={loanForm.control}
                      name="maxLoanAmount"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Monto Máximo de Préstamo</FormLabel>
                          <FormControl>
                            <Input type="number" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  <div className="flex justify-end">
                    <Button type="submit" disabled={setMultipleMutation.isPending}>
                      <Save className="mr-2 h-4 w-4" />
                      Guardar Cambios
                    </Button>
                  </div>
                </form>
              </Form>
            </CardContent>
          </Card>
        </TabsContent>

        {/* System Settings */}
        <TabsContent value="system">
          <Card>
            <CardHeader>
              <CardTitle>Configuración del Sistema</CardTitle>
              <CardDescription>
                Opciones generales del sistema
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Form {...systemForm}>
                <form onSubmit={systemForm.handleSubmit(handleSystemSave)} className="space-y-4">
                  <div className="grid gap-4 sm:grid-cols-2">
                    <FormField
                      control={systemForm.control}
                      name="timezone"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Zona Horaria</FormLabel>
                          <Select onValueChange={field.onChange} value={field.value}>
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue placeholder="Seleccionar zona horaria" />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              <SelectItem value="America/Mexico_City">Ciudad de México (GMT-6)</SelectItem>
                              <SelectItem value="America/Tijuana">Tijuana (GMT-8)</SelectItem>
                              <SelectItem value="America/Monterrey">Monterrey (GMT-6)</SelectItem>
                              <SelectItem value="America/Cancun">Cancún (GMT-5)</SelectItem>
                            </SelectContent>
                          </Select>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={systemForm.control}
                      name="currency"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Moneda</FormLabel>
                          <Select onValueChange={field.onChange} value={field.value}>
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue placeholder="Seleccionar moneda" />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              <SelectItem value="MXN">Peso Mexicano (MXN)</SelectItem>
                              <SelectItem value="USD">Dólar Americano (USD)</SelectItem>
                              <SelectItem value="EUR">Euro (EUR)</SelectItem>
                            </SelectContent>
                          </Select>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  <FormField
                    control={systemForm.control}
                    name="dateFormat"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Formato de Fecha</FormLabel>
                        <Select onValueChange={field.onChange} value={field.value}>
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Seleccionar formato" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            <SelectItem value="dd/MM/yyyy">DD/MM/AAAA (31/12/2024)</SelectItem>
                            <SelectItem value="MM/dd/yyyy">MM/DD/AAAA (12/31/2024)</SelectItem>
                            <SelectItem value="yyyy-MM-dd">AAAA-MM-DD (2024-12-31)</SelectItem>
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={systemForm.control}
                    name="allowNegativeCash"
                    render={({ field }) => (
                      <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                        <div className="space-y-0.5">
                          <FormLabel className="text-base">Permitir Caja Negativa</FormLabel>
                          <FormDescription>
                            Permite realizar operaciones aunque la caja quede en negativo
                          </FormDescription>
                        </div>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />

                  <div className="flex justify-end">
                    <Button type="submit" disabled={setMultipleMutation.isPending}>
                      <Save className="mr-2 h-4 w-4" />
                      Guardar Cambios
                    </Button>
                  </div>
                </form>
              </Form>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Notification Settings */}
        <TabsContent value="notifications">
          <Card>
            <CardHeader>
              <CardTitle>Configuración de Notificaciones</CardTitle>
              <CardDescription>
                Canales y opciones de notificación a clientes
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Form {...notificationForm}>
                <form onSubmit={notificationForm.handleSubmit(handleNotificationSave)} className="space-y-4">
                  <FormField
                    control={notificationForm.control}
                    name="emailEnabled"
                    render={({ field }) => (
                      <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                        <div className="space-y-0.5">
                          <FormLabel className="text-base">Notificaciones por Email</FormLabel>
                          <FormDescription>
                            Enviar recordatorios y notificaciones por correo electrónico
                          </FormDescription>
                        </div>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={notificationForm.control}
                    name="smsEnabled"
                    render={({ field }) => (
                      <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                        <div className="space-y-0.5">
                          <FormLabel className="text-base">Notificaciones por SMS</FormLabel>
                          <FormDescription>
                            Enviar mensajes de texto a clientes
                          </FormDescription>
                        </div>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={notificationForm.control}
                    name="whatsappEnabled"
                    render={({ field }) => (
                      <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                        <div className="space-y-0.5">
                          <FormLabel className="text-base">Notificaciones por WhatsApp</FormLabel>
                          <FormDescription>
                            Enviar mensajes de WhatsApp a clientes
                          </FormDescription>
                        </div>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={notificationForm.control}
                    name="reminderDaysBefore"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Días de Anticipación para Recordatorios</FormLabel>
                        <FormControl>
                          <Input type="number" min="1" max="30" {...field} />
                        </FormControl>
                        <FormDescription>
                          Enviar recordatorio de pago X días antes del vencimiento
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className="flex justify-end">
                    <Button type="submit" disabled={setMultipleMutation.isPending}>
                      <Save className="mr-2 h-4 w-4" />
                      Guardar Cambios
                    </Button>
                  </div>
                </form>
              </Form>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Account Settings */}
        <TabsContent value="account">
          <Card>
            <CardHeader>
              <CardTitle>Mi Cuenta</CardTitle>
              <CardDescription>
                Información de tu cuenta de usuario
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <p className="text-sm text-muted-foreground">Nombre</p>
                  <p className="font-medium">{user?.first_name} {user?.last_name}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Correo Electrónico</p>
                  <p className="font-medium">{user?.email}</p>
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <p className="text-sm text-muted-foreground">Rol</p>
                  <p className="font-medium">{user?.role?.display_name || user?.role?.name}</p>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Sucursal</p>
                  <p className="font-medium">{selectedBranch?.name || 'Todas las sucursales'}</p>
                </div>
              </div>

              <div className="pt-4">
                <p className="text-sm text-muted-foreground mb-2">
                  Para cambiar tu contraseña o información personal, contacta a un administrador.
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
