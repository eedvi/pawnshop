import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useLogin } from '@/hooks/use-auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { ApiErrorException } from '@/types/api'

const loginSchema = z.object({
  email: z.string().trim().email('Email inválido'),
  password: z.string().min(1, 'La contraseña es requerida'),
})

type LoginFormData = z.infer<typeof loginSchema>

export default function LoginPage() {
  const loginMutation = useLogin()

  const form = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  })

  const onSubmit = (data: LoginFormData) => {
    loginMutation.mutate(data)
  }

  const getErrorMessage = () => {
    if (!loginMutation.error) return null
    if (loginMutation.error instanceof ApiErrorException) {
      // Mensajes amigables según el código de error
      const code = loginMutation.error.code
      if (code === 'UNAUTHORIZED' || code === 'INVALID_CREDENTIALS') {
        return 'Email o contraseña incorrectos'
      }
      if (code === 'USER_INACTIVE') {
        return 'Tu cuenta está desactivada. Contacta al administrador.'
      }
      if (code === 'TOO_MANY_REQUESTS') {
        return 'Demasiados intentos. Espera unos minutos.'
      }
      return loginMutation.error.message
    }
    // Error de red u otro tipo
    const errMsg = (loginMutation.error as Error)?.message || ''
    if (errMsg.includes('Network') || errMsg.includes('network')) {
      return 'Error de conexión. Verifica que el servidor esté corriendo.'
    }
    return 'Error inesperado. Intenta de nuevo.'
  }

  const errorMessage = getErrorMessage()

  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="w-full max-w-md space-y-8 rounded-lg border bg-card p-8 shadow-lg">
        <div className="text-center">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-lg bg-primary text-primary-foreground text-xl font-bold">
            P
          </div>
          <h2 className="mt-4 text-2xl font-bold">PawnShop Admin</h2>
          <p className="mt-2 text-sm text-muted-foreground">
            Ingresa tus credenciales para continuar
          </p>
        </div>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            {errorMessage && (
              <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
                {errorMessage}
              </div>
            )}

            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input
                      type="email"
                      placeholder="tu@email.com"
                      autoComplete="email"
                      disabled={loginMutation.isPending}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Contraseña</FormLabel>
                  <FormControl>
                    <Input
                      type="password"
                      placeholder="••••••••"
                      autoComplete="current-password"
                      disabled={loginMutation.isPending}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" className="w-full" disabled={loginMutation.isPending}>
              {loginMutation.isPending ? 'Ingresando...' : 'Ingresar'}
            </Button>
          </form>
        </Form>
      </div>
    </div>
  )
}
