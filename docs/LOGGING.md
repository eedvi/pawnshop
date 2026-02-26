# Sistema de Logging Mejorado

Este documento describe las mejoras implementadas en el sistema de logging siguiendo las mejores prácticas de la industria.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Configuración](#configuración)
- [Sanitización de Datos](#sanitización-de-datos)
- [Slow Query Logging](#slow-query-logging)
- [Business Events](#business-events)
- [Ejemplos de Uso](#ejemplos-de-uso)

## ✨ Características

### 1. **Sanitización Automática de Datos Sensibles**

El sistema ahora sanitiza automáticamente datos sensibles en los logs:

- **Passwords**: `password`, `secret`, `pin`, `otp`
- **Tokens**: `token`, `api_key`, `access_token`, `refresh_token`
- **Datos financieros**: Números de tarjetas de crédito
- **Credenciales**: Passwords en URLs y headers Bearer

#### Campos Sensibles Detectados

El sistema reconoce automáticamente campos con estos nombres:
```go
password, token, secret, api_key, apikey, authorization,
credit_card, ssn, pin, otp, refresh_token, access_token
```

### 2. **Slow Query Logging**

Detecta y registra automáticamente queries de base de datos lentas:

```yaml
logging:
  slow_query_threshold: "1s"  # Queries > 1s se loggean como WARNING
  log_all_queries: false      # Debug: loggea todas las queries
```

**Información registrada:**
- Query SQL (sanitizada)
- Duración en milisegundos
- Request ID (trazabilidad)
- User ID
- Tipo de query (query, exec, query_row)

**Ejemplo de log:**
```json
{
  "level": "warn",
  "query": "SELECT * FROM loans WHERE ... (truncated)",
  "duration_ms": 1234,
  "threshold": "1s",
  "query_type": "query",
  "request_id": "abc-123",
  "user_id": 1,
  "message": "Slow database query detected"
}
```

### 3. **Business Event Logging**

Logging estructurado de eventos de negocio críticos:

#### Eventos de Préstamos
- `loan_created`: Préstamo creado
- `loan_approved`: Préstamo aprobado
- `loan_renewed`: Préstamo renovado
- `loan_completed`: Préstamo pagado completamente
- `loan_defaulted`: Préstamo en default

#### Eventos de Pagos
- `payment_received`: Pago recibido
- `payment_failed`: Pago fallido

#### Eventos de Ventas
- `sale_completed`: Venta completada
- `sale_refunded`: Venta reembolsada

#### Eventos de Caja
- `cash_session_opened`: Sesión abierta
- `cash_session_closed`: Sesión cerrada
- `cash_discrepancy`: Discrepancia de efectivo

#### Eventos de Artículos
- `item_received`: Artículo recibido
- `item_redeemed`: Artículo redimido
- `item_sold`: Artículo vendido

#### Eventos de Autenticación
- `user_login`: Login exitoso
- `login_failed`: Login fallido
- `user_logout`: Logout

### 4. **Structured Logging**

Todos los logs usan formato estructurado (JSON en producción):

```json
{
  "level": "info",
  "time": "2026-02-25T18:41:33-06:00",
  "request_id": "abc-123",
  "user_id": 1,
  "event_type": "loan_created",
  "loan_id": 42,
  "customer_id": 15,
  "amount": 1000.00,
  "interest_rate": 5.5,
  "service": "loan",
  "message": "Loan created successfully"
}
```

## ⚙️ Configuración

### Archivo de Configuración (`config.yaml`)

```yaml
logging:
  # Nivel de log: debug, info, warn, error
  level: "info"

  # Formato: json (producción) o console (desarrollo)
  format: "console"

  # Umbral para slow queries (1 segundo recomendado)
  slow_query_threshold: "1s"

  # Loggear todas las queries (solo en debug)
  log_all_queries: false
```

### Variables de Entorno

Puedes sobrescribir la configuración con variables de entorno:

```bash
export PAWN_LOGGING_LEVEL=debug
export PAWN_LOGGING_FORMAT=json
export PAWN_LOGGING_SLOW_QUERY_THRESHOLD=500ms
export PAWN_LOGGING_LOG_ALL_QUERIES=true
```

### Niveles de Log

| Nivel | Descripción | Cuándo usar |
|-------|-------------|-------------|
| `debug` | Información detallada | Desarrollo y debugging |
| `info` | Eventos normales | Producción (recomendado) |
| `warn` | Situaciones anómalas no críticas | Siempre habilitado |
| `error` | Errores que requieren atención | Siempre habilitado |

## 🔒 Sanitización de Datos

### Uso en Código

#### Sanitizar Strings
```go
import "pawnshop/pkg/logger"

// Sanitiza un string individual
sanitized := logger.Sanitize("postgres://user:password@localhost/db")
// Resultado: "postgres://$user:***@localhost/db"
```

#### Sanitizar Maps
```go
data := map[string]interface{}{
    "username": "admin",
    "password": "secret123",
    "email": "admin@example.com",
}

sanitized := logger.SanitizeMap(data)
// password estará como "***REDACTED***"
```

#### Sanitizar SQL Queries
```go
query := "UPDATE users SET password = 'secret123' WHERE id = 1"
sanitized := logger.SanitizeSQL(query)
// Resultado: "UPDATE users SET password = '***' WHERE id = 1"
```

#### Tipo SanitizedString
```go
import "github.com/rs/zerolog/log"

password := "secret123"
log.Info().
    Str("username", "admin").
    Str("password", string(logger.SanitizedString(password))).
    Msg("User login attempt")
// password se muestra como "***REDACTED***"
```

## 🐢 Slow Query Logging

### Cómo Funciona

El sistema automáticamente mide el tiempo de cada query de base de datos:

1. **Queries normales** (<1s): Solo se loggean en modo debug
2. **Queries lentas** (>1s): Se loggean como WARNING
3. **Errores**: Siempre se loggean como ERROR

### Ejemplo de Log de Query Lenta

```
2026-02-25T18:41:33-06:00 WRN Slow database query detected
  query="SELECT * FROM loans WHERE ..."
  duration_ms=1234
  threshold=1s
  query_type=query
  request_id=abc-123
  user_id=1
```

### Optimización

Cuando veas una query lenta:

1. Copia el query del log
2. Ejecuta `EXPLAIN ANALYZE` en PostgreSQL
3. Agrega índices si es necesario
4. Optimiza el query

## 📊 Business Events

### Uso en Servicios

```go
type LoanService struct {
    // ...
    businessLogger *logger.BusinessLogger
}

func NewLoanService(...) *LoanService {
    serviceLogger := log.With().Str("service", "loan").Logger()
    return &LoanService{
        // ...
        businessLogger: logger.NewBusinessLogger(serviceLogger),
    }
}

func (s *LoanService) Create(ctx context.Context, input CreateLoanInput) (*domain.Loan, error) {
    // ... crear préstamo ...

    // Log business event
    s.businessLogger.LoanCreated(ctx, loan.ID, input.CustomerID, input.LoanAmount, input.InterestRate)

    return loan, nil
}
```

### Helper Function para Eventos Personalizados

```go
import "pawnshop/pkg/logger"

logger.LogBusinessEvent(ctx, "custom_event", "Something important happened", map[string]interface{}{
    "entity_id": 123,
    "action": "approved",
    "amount": 1000.00,
})
```

## 📝 Ejemplos de Uso

### 1. Log de Error con Sanitización

```go
func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var input LoginInput
    if err := c.BodyParser(&input); err != nil {
        // El password en input se sanitiza automáticamente
        log.Error().
            Err(err).
            Interface("input", logger.SanitizeMap(map[string]interface{}{
                "email": input.Email,
                "password": input.Password,
            })).
            Msg("Failed to parse login request")
        return response.BadRequest(c, "Invalid request")
    }
    // ...
}
```

### 2. Log de Evento de Negocio

```go
func (s *PaymentService) ProcessPayment(ctx context.Context, input PaymentInput) (*domain.Payment, error) {
    // ... procesar pago ...

    // Log business event
    s.businessLogger.PaymentReceived(ctx, payment.ID, payment.LoanID, payment.Amount, string(payment.PaymentMethod))

    return payment, nil
}
```

### 3. Query con Context Tracking

```go
func (r *LoanRepository) GetByID(ctx context.Context, id int64) (*domain.Loan, error) {
    // El context incluye request_id y user_id automáticamente
    row := r.db.QueryRowContext(ctx,
        "SELECT * FROM loans WHERE id = $1", id)

    // Si la query es lenta (>1s), se loggea automáticamente con:
    // - request_id del contexto
    // - user_id del contexto
    // - duración de la query

    // ...
}
```

### 4. Logging en Diferentes Entornos

#### Desarrollo (console format)
```
2026-02-25T18:41:33-06:00 INF Loan created successfully
  loan_id=42
  customer_id=15
  amount=1000.00
  event_type=loan_created
```

#### Producción (JSON format)
```json
{
  "level":"info",
  "time":"2026-02-25T18:41:33-06:00",
  "loan_id":42,
  "customer_id":15,
  "amount":1000.00,
  "event_type":"loan_created",
  "request_id":"abc-123",
  "message":"Loan created successfully"
}
```

## 🎯 Best Practices

### 1. **Usa el Context Apropiado**
Siempre pasa el context con request_id y user_id:
```go
func (s *Service) DoSomething(ctx context.Context) error {
    // El logger automáticamente incluye request_id y user_id del context
    logger := logger.FromContext(ctx, s.logger)
    logger.Info().Msg("Doing something")
}
```

### 2. **Log en el Nivel Correcto**
- **Debug**: Información detallada de debugging
- **Info**: Eventos normales del negocio
- **Warn**: Situaciones anómalas pero no críticas
- **Error**: Errores que requieren atención

### 3. **No Loggees Datos Sensibles**
```go
// ❌ MAL
log.Info().Str("password", user.Password).Msg("User created")

// ✅ BIEN
log.Info().
    Str("password", string(logger.SanitizedString(user.Password))).
    Msg("User created")
```

### 4. **Usa Business Events**
Para eventos importantes del negocio, usa BusinessLogger:
```go
// ✅ BIEN
s.businessLogger.LoanCreated(ctx, loanID, customerID, amount, rate)

// En lugar de:
// log.Info().Int64("loan_id", loanID).Msg("Loan created") ❌
```

### 5. **Log Errores con Contexto**
```go
// ✅ BIEN
log.Error().
    Err(err).
    Int64("loan_id", loanID).
    Str("operation", "create_payment").
    Msg("Failed to create payment")
```

## 🔍 Análisis de Logs

### Buscar Queries Lentas
```bash
# En logs JSON
cat logs/app.log | jq 'select(.event_type == "slow_query")'

# Queries > 2 segundos
cat logs/app.log | jq 'select(.duration_ms > 2000)'
```

### Buscar Eventos de Negocio
```bash
# Todos los préstamos creados hoy
cat logs/app.log | jq 'select(.event_type == "loan_created")'

# Pagos fallidos
cat logs/app.log | jq 'select(.event_type == "payment_failed")'
```

### Trazabilidad de Request
```bash
# Seguir todos los logs de un request específico
cat logs/app.log | jq 'select(.request_id == "abc-123")'
```

## 📈 Métricas y Monitoring

Los logs estructurados pueden ser importados a herramientas como:

- **ELK Stack** (Elasticsearch, Logstash, Kibana)
- **Grafana Loki**
- **Datadog**
- **New Relic**
- **CloudWatch**

Ejemplo de dashboard en Grafana:
- Queries lentas por minuto
- Eventos de negocio por tipo
- Errores por servicio
- Discrepancias de caja

## 🚀 Próximos Pasos

Mejoras futuras planificadas:

1. **OpenTelemetry**: Integración con distributed tracing
2. **Log Sampling**: Reducir volumen en producción
3. **Log Rotation**: Rotación automática de archivos
4. **Dynamic Log Levels**: Cambiar niveles sin reiniciar
5. **Alertas**: Notificaciones automáticas en errores críticos

---

Para más información, ver:
- [Configuración](../config.example.yaml)
- [API Documentation](./API.md)
- [Database Schema](./DATABASE.md)
