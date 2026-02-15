import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Check, ChevronRight, User, Package, Calculator, FileText } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { Textarea } from '@/components/ui/textarea'
import { Skeleton } from '@/components/ui/skeleton'
import { ROUTES, customerRoute, itemRoute } from '@/routes/routes'
import { useCustomers } from '@/hooks/use-customers'
import { useItems } from '@/hooks/use-items'
import { useCreateLoan, useCalculateLoan } from '@/hooks/use-loans'
import { useBranchStore } from '@/stores/branch-store'
import { useDebounce } from '@/hooks'
import { Customer, Item, PAYMENT_PLAN_TYPES, PaymentPlanType } from '@/types'
import { formatCurrency } from '@/lib/format'

const STEPS = [
  { id: 'customer', title: 'Cliente', icon: User },
  { id: 'item', title: 'Artículo', icon: Package },
  { id: 'terms', title: 'Condiciones', icon: Calculator },
  { id: 'review', title: 'Confirmar', icon: FileText },
]

export function LoanWizard() {
  const navigate = useNavigate()
  const { selectedBranchId } = useBranchStore()
  const [currentStep, setCurrentStep] = useState(0)

  // Form state
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(null)
  const [selectedItem, setSelectedItem] = useState<Item | null>(null)
  const [loanAmount, setLoanAmount] = useState<number>(0)
  const [interestRate, setInterestRate] = useState<number>(15)
  const [paymentPlanType, setPaymentPlanType] = useState<PaymentPlanType>('single')
  const [loanTermDays, setLoanTermDays] = useState<number>(30)
  const [gracePeriodDays, setGracePeriodDays] = useState<number>(0)
  const [numberOfInstallments, setNumberOfInstallments] = useState<number>(3)
  const [notes, setNotes] = useState<string>('')

  // Search state
  const [customerSearch, setCustomerSearch] = useState('')
  const [itemSearch, setItemSearch] = useState('')
  const debouncedCustomerSearch = useDebounce(customerSearch, 300)
  const debouncedItemSearch = useDebounce(itemSearch, 300)

  // Queries
  const { data: customersData, isLoading: loadingCustomers } = useCustomers({
    search: debouncedCustomerSearch || undefined,
    branch_id: selectedBranchId ?? undefined,
    is_active: true,
    is_blocked: false,
    per_page: 10,
  })

  const { data: itemsData, isLoading: loadingItems } = useItems({
    search: debouncedItemSearch || undefined,
    branch_id: selectedBranchId ?? undefined,
    customer_id: selectedCustomer?.id,
    status: 'available',
    per_page: 10,
  })

  // Mutations
  const createMutation = useCreateLoan()
  const calculateMutation = useCalculateLoan()

  // Calculate on terms change
  const handleCalculate = () => {
    if (!selectedBranchId || !selectedCustomer || !selectedItem || loanAmount <= 0) return

    calculateMutation.mutate({
      branch_id: selectedBranchId,
      customer_id: selectedCustomer.id,
      item_id: selectedItem.id,
      loan_amount: loanAmount,
      interest_rate: interestRate,
      payment_plan_type: paymentPlanType,
      loan_term_days: loanTermDays,
      grace_period_days: gracePeriodDays,
      number_of_installments: paymentPlanType === 'installments' ? numberOfInstallments : undefined,
    })
  }

  const handleSubmit = () => {
    if (!selectedBranchId || !selectedCustomer || !selectedItem) return

    createMutation.mutate(
      {
        branch_id: selectedBranchId,
        customer_id: selectedCustomer.id,
        item_id: selectedItem.id,
        loan_amount: loanAmount,
        interest_rate: interestRate,
        payment_plan_type: paymentPlanType,
        loan_term_days: loanTermDays,
        grace_period_days: gracePeriodDays,
        number_of_installments: paymentPlanType === 'installments' ? numberOfInstallments : undefined,
        notes: notes || undefined,
      },
      {
        onSuccess: () => {
          navigate(ROUTES.LOANS)
        },
      }
    )
  }

  const canProceed = () => {
    switch (currentStep) {
      case 0:
        return selectedCustomer !== null
      case 1:
        return selectedItem !== null
      case 2:
        return loanAmount > 0
      case 3:
        return true
      default:
        return false
    }
  }

  const handleNext = () => {
    if (currentStep === 2) {
      handleCalculate()
    }
    if (currentStep < STEPS.length - 1) {
      setCurrentStep(currentStep + 1)
    }
  }

  const handleBack = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1)
    }
  }

  const selectCustomer = (customer: Customer) => {
    setSelectedCustomer(customer)
    setSelectedItem(null) // Reset item when customer changes
  }

  const selectItem = (item: Item) => {
    setSelectedItem(item)
    setLoanAmount(item.loan_value)
  }

  return (
    <div className="space-y-6">
      {/* Steps indicator */}
      <div className="flex items-center justify-center">
        {STEPS.map((step, index) => (
          <div key={step.id} className="flex items-center">
            <div
              className={`flex items-center justify-center w-10 h-10 rounded-full border-2 ${
                index < currentStep
                  ? 'bg-primary border-primary text-primary-foreground'
                  : index === currentStep
                  ? 'border-primary text-primary'
                  : 'border-muted text-muted-foreground'
              }`}
            >
              {index < currentStep ? (
                <Check className="w-5 h-5" />
              ) : (
                <step.icon className="w-5 h-5" />
              )}
            </div>
            <span
              className={`ml-2 text-sm font-medium ${
                index === currentStep ? 'text-primary' : 'text-muted-foreground'
              }`}
            >
              {step.title}
            </span>
            {index < STEPS.length - 1 && (
              <ChevronRight className="w-5 h-5 mx-4 text-muted-foreground" />
            )}
          </div>
        ))}
      </div>

      {/* Step content */}
      <Card>
        <CardHeader>
          <CardTitle>{STEPS[currentStep].title}</CardTitle>
          <CardDescription>
            {currentStep === 0 && 'Selecciona el cliente para el préstamo'}
            {currentStep === 1 && 'Selecciona el artículo a empeñar'}
            {currentStep === 2 && 'Define las condiciones del préstamo'}
            {currentStep === 3 && 'Revisa y confirma el préstamo'}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {/* Step 1: Customer */}
          {currentStep === 0 && (
            <div className="space-y-4">
              <Input
                placeholder="Buscar por nombre o documento..."
                value={customerSearch}
                onChange={(e) => setCustomerSearch(e.target.value)}
              />
              <div className="space-y-2 max-h-80 overflow-y-auto">
                {loadingCustomers ? (
                  Array.from({ length: 3 }).map((_, i) => (
                    <Skeleton key={i} className="h-16 w-full" />
                  ))
                ) : customersData?.data?.length === 0 ? (
                  <p className="text-center text-muted-foreground py-4">
                    No se encontraron clientes
                  </p>
                ) : (
                  customersData?.data?.map((customer) => (
                    <div
                      key={customer.id}
                      onClick={() => selectCustomer(customer)}
                      className={`p-4 rounded-lg border cursor-pointer transition-colors ${
                        selectedCustomer?.id === customer.id
                          ? 'border-primary bg-primary/5'
                          : 'hover:bg-muted'
                      }`}
                    >
                      <p className="font-medium">
                        {customer.first_name} {customer.last_name}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {customer.identity_number} • Límite: {formatCurrency(customer.credit_limit)}
                      </p>
                    </div>
                  ))
                )}
              </div>
            </div>
          )}

          {/* Step 2: Item */}
          {currentStep === 1 && (
            <div className="space-y-4">
              <Input
                placeholder="Buscar artículo..."
                value={itemSearch}
                onChange={(e) => setItemSearch(e.target.value)}
              />
              <div className="space-y-2 max-h-80 overflow-y-auto">
                {loadingItems ? (
                  Array.from({ length: 3 }).map((_, i) => (
                    <Skeleton key={i} className="h-16 w-full" />
                  ))
                ) : itemsData?.data?.length === 0 ? (
                  <p className="text-center text-muted-foreground py-4">
                    No hay artículos disponibles
                  </p>
                ) : (
                  itemsData?.data?.map((item) => (
                    <div
                      key={item.id}
                      onClick={() => selectItem(item)}
                      className={`p-4 rounded-lg border cursor-pointer transition-colors ${
                        selectedItem?.id === item.id
                          ? 'border-primary bg-primary/5'
                          : 'hover:bg-muted'
                      }`}
                    >
                      <p className="font-medium">{item.name}</p>
                      <p className="text-sm text-muted-foreground">
                        {item.sku} • Avalúo: {formatCurrency(item.appraised_value)} •
                        Préstamo: {formatCurrency(item.loan_value)}
                      </p>
                    </div>
                  ))
                )}
              </div>
            </div>
          )}

          {/* Step 3: Terms */}
          {currentStep === 2 && (
            <div className="space-y-6">
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <Label>Monto del Préstamo</Label>
                  <Input
                    type="number"
                    value={loanAmount}
                    onChange={(e) => setLoanAmount(Number(e.target.value))}
                  />
                  <p className="text-xs text-muted-foreground">
                    Máximo sugerido: {formatCurrency(selectedItem?.loan_value || 0)}
                  </p>
                </div>
                <div className="space-y-2">
                  <Label>Tasa de Interés (%)</Label>
                  <Input
                    type="number"
                    value={interestRate}
                    onChange={(e) => setInterestRate(Number(e.target.value))}
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label>Tipo de Pago</Label>
                <RadioGroup
                  value={paymentPlanType}
                  onValueChange={(v) => setPaymentPlanType(v as PaymentPlanType)}
                >
                  {PAYMENT_PLAN_TYPES.map((type) => (
                    <div key={type.value} className="flex items-start space-x-3">
                      <RadioGroupItem value={type.value} id={type.value} />
                      <div>
                        <Label htmlFor={type.value} className="cursor-pointer">
                          {type.label}
                        </Label>
                        <p className="text-sm text-muted-foreground">{type.description}</p>
                      </div>
                    </div>
                  ))}
                </RadioGroup>
              </div>

              <div className="grid gap-4 sm:grid-cols-3">
                <div className="space-y-2">
                  <Label>Plazo (días)</Label>
                  <Input
                    type="number"
                    value={loanTermDays}
                    onChange={(e) => setLoanTermDays(Number(e.target.value))}
                    disabled={paymentPlanType === 'installments'}
                  />
                  {paymentPlanType === 'installments' && (
                    <p className="text-xs text-muted-foreground">
                      Se calcula automáticamente basado en las cuotas mensuales
                    </p>
                  )}
                </div>
                <div className="space-y-2">
                  <Label>Período de Gracia (días)</Label>
                  <Input
                    type="number"
                    value={gracePeriodDays}
                    onChange={(e) => setGracePeriodDays(Number(e.target.value))}
                  />
                </div>
                {paymentPlanType === 'installments' && (
                  <div className="space-y-2">
                    <Label>Número de Cuotas</Label>
                    <Input
                      type="number"
                      value={numberOfInstallments}
                      onChange={(e) => setNumberOfInstallments(Number(e.target.value))}
                    />
                  </div>
                )}
              </div>

              <div className="space-y-2">
                <Label>Notas</Label>
                <Textarea
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                  placeholder="Observaciones sobre el préstamo..."
                  rows={2}
                />
              </div>
            </div>
          )}

          {/* Step 4: Review */}
          {currentStep === 3 && (
            <div className="space-y-6">
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="p-4 rounded-lg border">
                  <p className="text-sm text-muted-foreground">Cliente</p>
                  <p className="font-medium">
                    {selectedCustomer?.first_name} {selectedCustomer?.last_name}
                  </p>
                  <p className="text-sm">{selectedCustomer?.identity_number}</p>
                </div>
                <div className="p-4 rounded-lg border">
                  <p className="text-sm text-muted-foreground">Artículo</p>
                  <p className="font-medium">{selectedItem?.name}</p>
                  <p className="text-sm">{selectedItem?.sku}</p>
                </div>
              </div>

              <div className="p-4 rounded-lg border space-y-2">
                <p className="font-medium">Resumen del Préstamo</p>
                <div className="grid grid-cols-2 gap-2 text-sm">
                  <span className="text-muted-foreground">Monto:</span>
                  <span>{formatCurrency(loanAmount)}</span>
                  <span className="text-muted-foreground">Tasa de Interés:</span>
                  <span>{interestRate}%</span>
                  <span className="text-muted-foreground">Tipo de Pago:</span>
                  <span>
                    {PAYMENT_PLAN_TYPES.find((t) => t.value === paymentPlanType)?.label}
                  </span>
                  <span className="text-muted-foreground">Plazo:</span>
                  <span>
                    {paymentPlanType === 'installments'
                      ? `${numberOfInstallments} meses (~${numberOfInstallments * 30} días)`
                      : `${loanTermDays} días`}
                  </span>
                  {calculateMutation.data && (
                    <>
                      <span className="text-muted-foreground">Interés:</span>
                      <span>{formatCurrency(calculateMutation.data.interest_amount)}</span>
                      <span className="text-muted-foreground font-medium">Total a Pagar:</span>
                      <span className="font-bold">
                        {formatCurrency(calculateMutation.data.total_amount)}
                      </span>
                      {paymentPlanType === 'installments' && numberOfInstallments > 0 && (
                        <>
                          <span className="text-muted-foreground">Número de Cuotas:</span>
                          <span>{numberOfInstallments}</span>
                          <span className="text-muted-foreground">Monto por Cuota:</span>
                          <span className="font-medium">
                            {formatCurrency(calculateMutation.data.total_amount / numberOfInstallments)}
                          </span>
                        </>
                      )}
                    </>
                  )}
                </div>
              </div>

              {notes && (
                <div className="p-4 rounded-lg border">
                  <p className="text-sm text-muted-foreground">Notas</p>
                  <p>{notes}</p>
                </div>
              )}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Navigation */}
      <div className="flex justify-between">
        <Button
          variant="outline"
          onClick={handleBack}
          disabled={currentStep === 0}
        >
          Anterior
        </Button>
        <div className="flex gap-2">
          <Button
            variant="outline"
            onClick={() => navigate(ROUTES.LOANS)}
          >
            Cancelar
          </Button>
          {currentStep < STEPS.length - 1 ? (
            <Button onClick={handleNext} disabled={!canProceed()}>
              Siguiente
            </Button>
          ) : (
            <Button
              onClick={handleSubmit}
              disabled={createMutation.isPending}
            >
              {createMutation.isPending ? 'Creando...' : 'Crear Préstamo'}
            </Button>
          )}
        </div>
      </div>
    </div>
  )
}
