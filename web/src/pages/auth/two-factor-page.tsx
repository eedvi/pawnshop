import { useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { useTwoFactorVerify } from '@/hooks/use-auth'
import { Button } from '@/components/ui/button'
import { ROUTES } from '@/routes/routes'
import { ApiErrorException } from '@/types/api'

export default function TwoFactorPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const [code, setCode] = useState('')
  const verifyMutation = useTwoFactorVerify()

  const challengeToken = (location.state as { challengeToken?: string })?.challengeToken

  // Redirect if no challenge token
  if (!challengeToken) {
    navigate(ROUTES.LOGIN, { replace: true })
    return null
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    verifyMutation.mutate({
      challenge_token: challengeToken,
      code,
    })
  }

  const errorMessage = verifyMutation.error instanceof ApiErrorException
    ? verifyMutation.error.message
    : verifyMutation.error
    ? 'Error de verificación'
    : null

  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="w-full max-w-md space-y-8 rounded-lg border bg-card p-8 shadow-lg">
        <div className="text-center">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-lg bg-primary text-primary-foreground text-xl font-bold">
            P
          </div>
          <h2 className="mt-4 text-2xl font-bold">Verificación de Dos Pasos</h2>
          <p className="mt-2 text-sm text-muted-foreground">
            Ingresa el código de 6 dígitos de tu aplicación autenticadora
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {errorMessage && (
            <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
              {errorMessage}
            </div>
          )}

          <div className="space-y-2">
            <label htmlFor="code" className="text-sm font-medium">
              Código
            </label>
            <input
              id="code"
              type="text"
              inputMode="numeric"
              pattern="[0-9]*"
              maxLength={6}
              className="flex h-12 w-full rounded-md border border-input bg-transparent px-3 py-1 text-center text-2xl tracking-widest shadow-sm transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
              placeholder="000000"
              value={code}
              onChange={(e) => setCode(e.target.value.replace(/\D/g, ''))}
              disabled={verifyMutation.isPending}
            />
          </div>

          <Button
            type="submit"
            className="w-full"
            disabled={verifyMutation.isPending || code.length !== 6}
          >
            {verifyMutation.isPending ? 'Verificando...' : 'Verificar'}
          </Button>

          <button
            type="button"
            onClick={() => navigate(ROUTES.LOGIN)}
            className="w-full text-sm text-muted-foreground hover:text-foreground"
          >
            Volver al inicio de sesión
          </button>
        </form>
      </div>
    </div>
  )
}
