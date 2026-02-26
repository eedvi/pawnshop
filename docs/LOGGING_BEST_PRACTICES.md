# Best Practices de Logging - Requests & Responses

## 📚 Estándares y Normativas

Este documento sigue las siguientes normas y mejores prácticas de la industria:

1. **OWASP** - Application Security Verification Standard (ASVS)
2. **GDPR** - General Data Protection Regulation
3. **PCI DSS** - Payment Card Industry Data Security Standard
4. **NIST** - National Institute of Standards and Technology
5. **OpenTelemetry** - Cloud Native Computing Foundation
6. **12-Factor App** - Methodology for cloud-native applications

## 🎯 Principios Fundamentales

### 1. **Principle of Least Privilege**
> Solo loggear lo necesario para debugging y auditoría

### 2. **Privacy by Design**
> No loggear PII (Personally Identifiable Information)

### 3. **Security First**
> NUNCA loggear credenciales, tokens, o datos sensibles

### 4. **Performance Consideration**
> El logging no debe impactar significativamente el performance

### 5. **Compliance Ready**
> Los logs deben cumplir con regulaciones (GDPR, CCPA, etc.)

## ✅ QUÉ LOGGEAR

### Request Metadata (SIEMPRE)

```json
{
  // Identificación y Trazabilidad
  "timestamp": "2026-02-25T22:56:51Z",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",  // OpenTelemetry
  "span_id": "00f067aa0ba902b7",
  "parent_span_id": "00f067aa0ba902b6",

  // Request Information
  "http": {
    "method": "POST",
    "url": "/api/v1/loans",
    "protocol": "HTTP/1.1",
    "scheme": "https",
    "host": "api.example.com"
  },

  // Client Information
  "client": {
    "ip": "192.168.1.100",           // Considerar GDPR: anonimizar último octeto
    "user_agent": "Mozilla/5.0...",
    "referer": "https://app.example.com/loans"
  },

  // Authentication Context
  "auth": {
    "user_id": 123,
    "tenant_id": "org-456",
    "session_id": "sess-789",       // Hash del session ID
    "auth_method": "jwt",
    "roles": ["admin", "loan_officer"]
  },

  // Business Context
  "business": {
    "resource_type": "loan",
    "resource_id": 23,
    "action": "create",
    "department": "operations",
    "branch_id": 5
  },

  // Performance Metrics
  "metrics": {
    "duration_ms": 53.77,
    "db_queries": 5,
    "db_duration_ms": 45.2,
    "cache_hits": 2,
    "cache_misses": 1
  },

  // Response Information
  "response": {
    "status_code": 201,
    "size_bytes": 1024,
    "content_type": "application/json"
  }
}
```

### Headers Seguros (WHITELIST)

Solo loggear headers específicos y seguros:

```go
var SafeHeaders = []string{
    "Content-Type",
    "Accept",
    "Accept-Language",
    "Accept-Encoding",
    "X-Request-ID",
    "X-Correlation-ID",
    "X-Forwarded-For",
    "X-Real-IP",
    "User-Agent",
    "Referer",
    "Content-Length",
}
```

### Query Parameters (CON CUIDADO)

```json
{
  "query_params": {
    // ✅ OK: Parámetros de paginación y filtrado
    "page": 1,
    "per_page": 20,
    "sort": "created_at",
    "filter": "active",

    // ❌ NO: Parámetros sensibles
    "token": "***REDACTED***",
    "api_key": "***REDACTED***",
    "password": "***REDACTED***"
  }
}
```

## ❌ QUÉ NO LOGGEAR NUNCA

### 1. Credenciales y Autenticación

```json
{
  // ❌ NUNCA
  "password": "secret123",
  "old_password": "old123",
  "new_password": "new123",
  "api_key": "sk_live_123456",
  "api_secret": "secret_key",
  "private_key": "-----BEGIN PRIVATE KEY-----",
  "oauth_token": "ya29.a0AfH6SMB...",
  "bearer_token": "eyJhbGciOiJIUzI1NiIs...",
  "authorization": "Bearer token123",
  "x-api-key": "key-123",
  "session_token": "sess_123",
  "refresh_token": "refresh_abc",
  "jwt": "eyJ...",
  "access_token": "access_xyz"
}
```

### 2. Información Personal (PII - GDPR)

```json
{
  // ❌ NUNCA (o solo con consentimiento explícito)
  "email": "user@example.com",
  "phone": "+1234567890",
  "address": "123 Main Street",
  "city": "New York",
  "postal_code": "10001",
  "birth_date": "1990-01-01",
  "ssn": "123-45-6789",
  "passport_number": "AB123456",
  "drivers_license": "D1234567",
  "national_id": "12345678A",
  "full_name": "John Doe",

  // ✅ OK: IDs internos (no PII)
  "user_id": 123,
  "customer_id": 456
}
```

### 3. Datos Financieros (PCI DSS)

```json
{
  // ❌ NUNCA
  "credit_card_number": "4532-1234-5678-9010",
  "cvv": "123",
  "card_expiry": "12/25",
  "bank_account": "1234567890",
  "routing_number": "123456789",
  "iban": "GB82 WEST 1234 5698 7654 32",
  "pin": "1234",

  // ✅ OK: Solo últimos 4 dígitos
  "card_last4": "9010",
  "card_brand": "visa"
}
```

### 4. Request/Response Bodies Completos

```json
{
  // ❌ NUNCA en producción
  "request_body": {
    "username": "admin",
    "password": "secret123",
    "credit_card": "4532123456789010"
  },

  // ✅ OK: Solo metadatos
  "request": {
    "body_size": 256,
    "content_type": "application/json",
    "has_file_upload": false,
    "field_count": 5
  }
}
```

### 5. Datos de Sesión

```json
{
  // ❌ NUNCA
  "session_data": {
    "cart": [...],
    "preferences": {...}
  },
  "cookie": "sessionid=abc123; token=xyz789",

  // ✅ OK: Solo hash/ID
  "session_id_hash": "sha256:abc123..."
}
```

## 🔒 Sanitización Obligatoria

### Patrones a Sanitizar Automáticamente

```go
var SensitivePatterns = []Pattern{
    // Passwords en URLs
    {Regex: `://[^:]+:([^@]+)@`, Replace: "://$user:***@"},

    // Bearer tokens
    {Regex: `Bearer\s+[A-Za-z0-9\-_\.]+`, Replace: "Bearer ***"},

    // API Keys
    {Regex: `api[_-]?key["']?\s*[:=]\s*["']?([^"'\s]+)`, Replace: "api_key=***"},

    // Credit cards (básico)
    {Regex: `\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`, Replace: "****-****-****-****"},

    // JWTs
    {Regex: `eyJ[A-Za-z0-9\-_]+\.[A-Za-z0-9\-_]+\.[A-Za-z0-9\-_]+`, Replace: "***JWT***"},

    // Emails (opcional según contexto)
    {Regex: `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`, Replace: "***@***.***"},
}
```

## 📊 Niveles de Logging por Ambiente

### Development
```yaml
logging:
  level: "debug"
  format: "console"
  request_logging:
    log_request_body: true          # ✅ OK
    log_response_body: true         # ✅ OK
    sanitize_sensitive: true        # ✅ SIEMPRE
    max_body_size: 10240            # 10KB
    log_headers: true
    log_query_params: true
```

### Staging
```yaml
logging:
  level: "info"
  format: "json"
  request_logging:
    log_request_body: false         # ❌ Metadatos solo
    log_response_body: false        # ❌ Metadatos solo
    sanitize_sensitive: true        # ✅ SIEMPRE
    max_body_size: 1024             # 1KB
    log_headers: true               # Solo whitelist
    log_query_params: true          # Sanitizado
```

### Production
```yaml
logging:
  level: "warn"                     # Solo warnings y errors
  format: "json"
  request_logging:
    log_request_body: false         # ❌ NUNCA
    log_response_body: false        # ❌ NUNCA
    sanitize_sensitive: true        # ✅ CRÍTICO
    max_body_size: 0                # No loggear bodies
    log_headers: false              # Solo en errores
    log_query_params: false         # Solo en errores
    sample_rate: 0.1                # Solo 10% de requests exitosos
    always_log_errors: true         # Siempre loggear errores
```

## 🎯 Casos de Uso Específicos

### 1. Debugging de Producción

**NUNCA** habilitar body logging en producción. En su lugar:

```bash
# Usar request_id para correlacionar con otros sistemas
# Ejemplo: buscar en logs de aplicación
grep "request_id=abc-123" /var/log/app.log

# O habilitar temporalmente para un usuario específico
curl -H "X-Debug-User: test@example.com" https://api.example.com
```

### 2. Auditoría de Seguridad

Loggear **quién** hizo **qué** y **cuándo**:

```json
{
  "type": "security_audit",
  "timestamp": "2026-02-25T22:56:51Z",
  "user_id": 123,
  "action": "loan.create",
  "resource_id": 23,
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "result": "success",
  "details": {
    "amount": 5000,
    "customer_id": 17
  }
}
```

### 3. Compliance (GDPR, CCPA)

```json
{
  // ✅ OK: Datos anonimizados
  "user_id_hash": "sha256:abc123...",
  "ip_anonymized": "192.168.1.0",   // Último octeto = 0
  "country": "US",
  "region": "CA",

  // ❌ NO: PII identificable
  "email": "user@example.com",
  "full_ip": "192.168.1.100"
}
```

### 4. Performance Monitoring

```json
{
  "type": "performance",
  "endpoint": "/api/v1/loans",
  "method": "POST",
  "duration_ms": 53.77,
  "db_queries": 5,
  "db_duration_ms": 45.2,
  "cache_hit_rate": 0.67,
  "memory_mb": 128,
  "cpu_percent": 15
}
```

## 🚨 Detección de Violaciones

### Automated Scanning

Implementar escaneo automático para detectar:

```bash
# Buscar passwords en logs
grep -i "password.*:" /var/log/app.log

# Buscar tokens
grep -E "(Bearer|token).*[A-Za-z0-9]{20,}" /var/log/app.log

# Buscar credit cards
grep -E "\b\d{4}[- ]?\d{4}[- ]?\d{4}[- ]?\d{4}\b" /var/log/app.log

# Buscar emails
grep -E "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}" /var/log/app.log
```

### Pre-commit Hooks

```bash
#!/bin/bash
# .git/hooks/pre-commit

# Verificar que no se commiteen logs con datos sensibles
if git diff --cached | grep -iE "(password|api_key|secret).*:"; then
    echo "ERROR: Possible sensitive data in logs"
    exit 1
fi
```

## 📋 Checklist de Implementación

- [ ] **Sanitización automática** de datos sensibles habilitada
- [ ] **PII removido** de todos los logs
- [ ] **Headers sensibles** filtrados (whitelist approach)
- [ ] **Body logging deshabilitado** en producción
- [ ] **Sample rate** configurado para requests normales
- [ ] **Errores siempre loggeados** con contexto completo
- [ ] **Request IDs** propagados en todos los logs
- [ ] **Trace context** (OpenTelemetry) implementado
- [ ] **Log rotation** configurado
- [ ] **Retention policy** definido (30-90 días)
- [ ] **Access controls** en archivos de log
- [ ] **Encryption at rest** para logs sensibles
- [ ] **SIEM integration** para alertas de seguridad
- [ ] **Compliance review** completado (GDPR, PCI DSS)
- [ ] **Documentation** actualizada
- [ ] **Team training** completado

## 🔗 Referencias

1. **OWASP Logging Cheat Sheet**
   https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html

2. **GDPR Article 32 - Security of Processing**
   https://gdpr-info.eu/art-32-gdpr/

3. **PCI DSS Requirement 10**
   https://www.pcisecuritystandards.org/

4. **NIST Special Publication 800-92**
   Guide to Computer Security Log Management

5. **OpenTelemetry Specification**
   https://opentelemetry.io/docs/specs/otel/

6. **12-Factor App - Logs**
   https://12factor.net/logs

7. **RFC 5424 - Syslog Protocol**
   https://tools.ietf.org/html/rfc5424

8. **Cloud Native Computing Foundation (CNCF)**
   https://www.cncf.io/

## 💡 Ejemplos de Empresas

### Google Cloud Logging
- No loggean request/response bodies
- Solo metadatos y métricas
- Sampling rate configurable

### AWS CloudWatch
- Headers sanitizados automáticamente
- PII detection automática
- Retention policies enforced

### DataDog
- Sensitive Data Scanner
- Automatic PII redaction
- Compliance templates (GDPR, HIPAA, PCI)

### Elastic (ELK Stack)
- Field-level security
- Data masking
- Audit logging

---

**Regla de Oro:** Cuando tengas duda sobre si loggear algo, **NO LO LOGGEES**.
Es mejor tener menos información que violar privacidad o seguridad.
