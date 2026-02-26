# 🚀 Logging Quick Reference

## TL;DR - Reglas de Oro

### ✅ LOGGEAR SIEMPRE:
- Request ID (trazabilidad)
- User ID (auditoría)
- Timestamp (cuando)
- HTTP method, path, status (qué)
- Duration/latency (performance)
- Client IP (quién/donde)

### ❌ NUNCA LOGGEAR:
- **Passwords** - Cualquier tipo
- **Tokens** - Bearer, JWT, API keys, session tokens
- **PII** - Emails, teléfonos, direcciones
- **Datos financieros** - Tarjetas, cuentas bancarias
- **Request/Response bodies completos** - En producción

---

## 📊 Configuración por Ambiente

### Development
```yaml
logging:
  level: "debug"
  format: "console"        # Legible para humanos
  log_request_body: true   # ✅ OK
  sanitize_sensitive: true # ✅ SIEMPRE
```

### Production
```yaml
logging:
  level: "warn"            # Solo warnings y errors
  format: "json"           # Para herramientas
  log_request_body: false  # ❌ NUNCA
  sanitize_sensitive: true # ✅ CRÍTICO
  sample_rate: 0.1         # Solo 10%
```

---

## 🎯 Request Logging - Qué Incluir

```json
{
  // ✅ SIEMPRE
  "request_id": "uuid",
  "user_id": 123,
  "method": "POST",
  "path": "/api/v1/loans",
  "status_code": 201,
  "duration_ms": 53.77,
  "client_ip": "192.168.1.100",

  // ✅ OK: Metadatos
  "request_body_size": 256,
  "request_content_type": "application/json",

  // ❌ NUNCA: Body completo
  "request_body": {...}
}
```

---

## 🔒 Headers - Whitelist Approach

### ✅ SEGURO para loggear:
- Content-Type
- Accept
- X-Request-ID
- X-Correlation-ID
- User-Agent
- Referer

### ❌ NUNCA loggear:
- Authorization
- Cookie
- Set-Cookie
- X-API-Key
- Proxy-Authorization

---

## 🛡️ Sanitización Automática

### Patrones Detectados:
```javascript
// URLs con passwords
postgres://user:password@host → postgres://user:***@host

// Bearer tokens
Bearer eyJhbG... → Bearer ***

// Credit cards
4532-1234-5678-9010 → ****-****-****-****

// API Keys
api_key=abc123 → api_key=***

// JWTs
eyJhbG... → ***JWT***
```

---

## 📋 Checklist Rápido

**Antes de deploy a producción:**

- [ ] `log_request_body: false`
- [ ] `log_response_body: false`
- [ ] `sanitize_sensitive: true`
- [ ] `format: "json"`
- [ ] `level: "warn"` o `"info"`
- [ ] Variables de entorno para secrets
- [ ] Sample rate configurado (<1.0)
- [ ] Headers whitelist implementado
- [ ] IP anonymization (GDPR)
- [ ] Retention policy definido
- [ ] Log rotation configurado

---

## 🚨 Red Flags - Revisar AHORA

Si ves estos patrones en tus logs, **FIX INMEDIATELY**:

```bash
# Buscar passwords
grep -i "password.*:" logs/*.log

# Buscar tokens
grep -E "Bearer [A-Za-z0-9]+" logs/*.log

# Buscar credit cards
grep -E "\b\d{4}[- ]?\d{4}[- ]?\d{4}[- ]?\d{4}\b" logs/*.log

# Buscar emails
grep -E "[a-z0-9]+@[a-z0-9]+" logs/*.log
```

---

## 💡 Examples

### ✅ CORRECTO:
```json
{
  "level": "info",
  "request_id": "550e8400-e29b-41d4",
  "user_id": 123,
  "method": "POST",
  "path": "/api/v1/loans",
  "status_code": 201,
  "duration_ms": 53.77,
  "client_ip": "192.168.1.0",
  "message": "Loan created"
}
```

### ❌ INCORRECTO:
```json
{
  "level": "info",
  "request_body": {
    "username": "admin",
    "password": "secret123",
    "email": "admin@example.com",
    "credit_card": "4532-1234-5678-9010"
  },
  "message": "Processing request"
}
```

---

## 🔗 Recursos

- **Documentación Completa**: `docs/LOGGING_BEST_PRACTICES.md`
- **Configuración**: `config.production-example.yaml`
- **Implementación**: `internal/middleware/request_logging.go`
- **Sanitización**: `pkg/logger/sanitize.go`

---

## 📞 Ayuda Rápida

**¿Tengo que loggear esto?**
- Si es un password/token/secret → **NO**
- Si es PII (email, teléfono) → **NO**
- Si es request/response body → **NO (en prod)**
- Si es para debugging → **Solo en dev**
- Si es para auditoría → **SI (sanitizado)**

**Regla simple:** Cuando dudes, **NO LO LOGGEES**.
