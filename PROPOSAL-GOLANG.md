# Propuesta Técnica: Sistema de Gestión de Casa de Empeño

## Migración a Golang

---

## 1. Resumen Ejecutivo

### 1.1 Descripción del Proyecto
Sistema integral de gestión para casas de empeño que permite administrar préstamos prendarios, clientes, inventario de artículos, pagos, ventas y generación de documentos legales. El sistema soporta operación multi-sucursal con control de acceso basado en roles.

### 1.2 Justificación de Golang
| Aspecto | Laravel (Actual) | Golang (Propuesto) |
|---------|------------------|-------------------|
| Rendimiento | ~1,000 req/s | ~50,000+ req/s |
| Consumo de memoria | ~100MB por proceso | ~10-20MB por instancia |
| Concurrencia | Basado en procesos | Goroutines nativas |
| Compilación | Interpretado | Binario estático |
| Despliegue | Requiere PHP + extensiones | Un solo binario |
| Tipado | Dinámico (parcial) | Estático fuerte |

### 1.3 Beneficios Esperados
- **50x mejor rendimiento** en operaciones concurrentes
- **Menor costo de infraestructura** (menos recursos de servidor)
- **Despliegue simplificado** (un solo binario, sin dependencias)
- **Mayor seguridad** por tipado estático y compilación
- **Mejor mantenibilidad** a largo plazo

---

## 2. Requerimientos de Negocio

### 2.1 Módulo de Clientes
| ID | Requerimiento | Prioridad |
|----|---------------|-----------|
| CL-01 | Registro de clientes con datos personales completos | Alta |
| CL-02 | Validación de edad mínima (18 años) | Alta |
| CL-03 | Documento de identidad único (DPI/Pasaporte) | Alta |
| CL-04 | Historial crediticio del cliente | Media |
| CL-05 | Límite de crédito configurable | Media |
| CL-06 | Contacto de emergencia | Baja |
| CL-07 | Programa de lealtad/puntos | Baja |

### 2.2 Módulo de Artículos (Inventario)
| ID | Requerimiento | Prioridad |
|----|---------------|-----------|
| IT-01 | Registro de artículos con categorización | Alta |
| IT-02 | Tasación de artículos (valor de mercado vs. préstamo) | Alta |
| IT-03 | Estados: Disponible, En préstamo, Vendido, Confiscado | Alta |
| IT-04 | Número de serie único (cuando aplique) | Media |
| IT-05 | Historial de movimientos del artículo | Media |
| IT-06 | Transferencias entre sucursales | Media |
| IT-07 | Fotografías del artículo | Baja |

### 2.3 Módulo de Préstamos
| ID | Requerimiento | Prioridad |
|----|---------------|-----------|
| LN-01 | Creación de préstamos prendarios | Alta |
| LN-02 | Cálculo automático de intereses (sobre saldo) | Alta |
| LN-03 | Plazo configurable (días) | Alta |
| LN-04 | Estados: Activo, Pagado, Vencido, Confiscado, Renovado | Alta |
| LN-05 | Pago mínimo mensual obligatorio | Alta |
| LN-06 | Período de gracia configurable | Media |
| LN-07 | Plan de cuotas (installments) | Media |
| LN-08 | Renovación de préstamos | Media |
| LN-09 | Recargos por mora | Media |
| LN-10 | Notificaciones de vencimiento | Baja |

### 2.4 Módulo de Pagos
| ID | Requerimiento | Prioridad |
|----|---------------|-----------|
| PY-01 | Registro de pagos parciales y totales | Alta |
| PY-02 | Aplicación: primero a intereses, luego a capital | Alta |
| PY-03 | Métodos de pago: Efectivo, Tarjeta, Transferencia | Alta |
| PY-04 | Generación de recibos | Alta |
| PY-05 | Historial de pagos por préstamo | Media |
| PY-06 | Reversión de pagos (con autorización) | Baja |

### 2.5 Módulo de Ventas
| ID | Requerimiento | Prioridad |
|----|---------------|-----------|
| SL-01 | Venta de artículos confiscados o propios | Alta |
| SL-02 | Descuentos y precio final | Media |
| SL-03 | Generación de comprobante de venta | Alta |
| SL-04 | Vinculación opcional con cliente | Media |

### 2.6 Módulo de Documentos
| ID | Requerimiento | Prioridad |
|----|---------------|-----------|
| DC-01 | Generación de contratos de préstamo (PDF) | Alta |
| DC-02 | Recibos de pago (formato ticket 80mm) | Alta |
| DC-03 | Comprobantes de venta | Alta |
| DC-04 | Plantillas personalizables | Media |
| DC-05 | Branding configurable (logo, colores, términos) | Media |
| DC-06 | Historial de documentos generados | Baja |

### 2.7 Módulo de Caja / POS (Point of Sale)
| ID | Requerimiento | Prioridad |
|----|---------------|-----------|
| POS-01 | Apertura de caja con monto inicial | Alta |
| POS-02 | Cierre de caja con arqueo | Alta |
| POS-03 | Registro de movimientos de efectivo (ingresos/egresos) | Alta |
| POS-04 | Corte parcial (X) y corte final (Z) | Alta |
| POS-05 | Múltiples cajas por sucursal | Media |
| POS-06 | Asignación de caja a usuario | Alta |
| POS-07 | Historial de operaciones por caja | Media |
| POS-08 | Conciliación automática de saldos | Media |
| POS-09 | Soporte para múltiples métodos de pago | Alta |
| POS-10 | Impresión de tickets en impresora térmica | Alta |
| POS-11 | Modo offline con sincronización | Media |
| POS-12 | Lector de códigos de barras/QR | Baja |

### 2.8 Módulo de Contabilidad Básica
| ID | Requerimiento | Prioridad |
|----|---------------|-----------|
| CT-01 | Registro de ingresos por préstamos | Alta |
| CT-02 | Registro de ingresos por intereses | Alta |
| CT-03 | Registro de ingresos por ventas | Alta |
| CT-04 | Registro de gastos operativos | Media |
| CT-05 | Balance diario por sucursal | Alta |
| CT-06 | Reporte de IVA (si aplica) | Media |
| CT-07 | Exportación para sistema contable externo | Baja |

### 2.9 Módulo de Administración
| ID | Requerimiento | Prioridad |
|----|---------------|-----------|
| AD-01 | Gestión de usuarios | Alta |
| AD-02 | Roles y permisos granulares | Alta |
| AD-03 | Multi-sucursal | Alta |
| AD-04 | Configuración global del sistema | Media |
| AD-05 | Auditoría de acciones (logs) | Media |
| AD-06 | Reportes y estadísticas | Media |
| AD-07 | Backup y restauración de datos | Alta |
| AD-08 | Configuración de impresoras por sucursal | Media |

### 2.10 Módulo de Notificaciones
| ID | Requerimiento | Prioridad |
|----|---------------|-----------|
| NT-01 | Notificación de vencimiento de préstamo (SMS/WhatsApp) | Media |
| NT-02 | Recordatorio de pago mínimo | Media |
| NT-03 | Alertas internas para empleados | Baja |
| NT-04 | Notificación de artículo próximo a confiscar | Media |

---

## 3. Requerimientos Técnicos

### 3.1 Requisitos Funcionales

#### 3.1.1 Autenticación y Autorización
- JWT (JSON Web Tokens) para sesiones stateless
- Refresh tokens con rotación
- Autenticación de dos factores (2FA) opcional
- Control de acceso basado en roles (RBAC)
- Permisos granulares por recurso y acción
- Sesiones por dispositivo con revocación

#### 3.1.2 API REST
- Versionado de API (v1, v2, etc.)
- Paginación cursor-based para listas grandes
- Filtrado y ordenamiento flexible
- Rate limiting por usuario/IP
- Respuestas consistentes (JSON:API o similar)
- Documentación OpenAPI/Swagger

#### 3.1.3 Base de Datos
- PostgreSQL como base de datos principal
- Migraciones versionadas
- Soft deletes para datos críticos
- Índices optimizados para consultas frecuentes
- Conexiones pooling
- Transacciones ACID para operaciones críticas

#### 3.1.4 Generación de PDFs
- Motor de renderizado HTML a PDF
- Soporte para plantillas dinámicas
- Formatos: Carta, Ticket (80mm)
- Caché de plantillas compiladas

#### 3.1.5 Tareas Programadas
- Cálculo de intereses diario
- Verificación de vencimientos
- Envío de notificaciones
- Limpieza de datos temporales
- Generación de reportes automáticos

### 3.2 Requisitos No Funcionales

#### 3.2.1 Rendimiento
| Métrica | Objetivo |
|---------|----------|
| Tiempo de respuesta API (p95) | < 100ms |
| Tiempo de respuesta API (p99) | < 500ms |
| Throughput mínimo | 1,000 req/s |
| Generación de PDF | < 2 segundos |
| Consultas de base de datos | < 50ms |

#### 3.2.2 Disponibilidad
- Uptime objetivo: 99.9%
- Recuperación ante fallos: < 5 minutos
- Backup automático diario
- Retención de backups: 30 días

#### 3.2.3 Seguridad
- HTTPS obligatorio (TLS 1.3)
- Encriptación de datos sensibles en reposo
- Hashing de contraseñas (Argon2id)
- Protección contra OWASP Top 10
- Headers de seguridad (CSP, HSTS, etc.)
- Sanitización de inputs
- Logs de auditoría inmutables

#### 3.2.4 Escalabilidad
- Diseño stateless para escalar horizontalmente
- Caché distribuida (Redis)
- Base de datos con réplicas de lectura
- Colas de mensajes para tareas asíncronas

---

## 4. Arquitectura Propuesta

### 4.1 Arquitectura General

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLIENTES                                  │
├─────────────┬─────────────┬─────────────┬─────────────┬─────────┤
│  Web App    │  Mobile App │   POS       │  Admin      │   API   │
│  (React)    │  (Flutter)  │  (Electron) │  (React)    │  Ext.   │
└──────┬──────┴──────┬──────┴──────┬──────┴──────┬──────┴────┬────┘
       │             │             │             │           │
       └─────────────┴─────────────┴─────────────┴───────────┘
                                   │
                    ┌──────────────▼──────────────┐
                    │      API Gateway / LB       │
                    │    (Traefik / Nginx)        │
                    └──────────────┬──────────────┘
                                   │
       ┌───────────────────────────┼───────────────────────────┐
       │                           │                           │
┌──────▼──────┐            ┌───────▼───────┐           ┌───────▼───────┐
│  Auth       │            │    Core       │           │   Documents   │
│  Service    │            │   Service     │           │   Service     │
│             │            │               │           │               │
│ - Login     │            │ - Customers   │           │ - PDF Gen     │
│ - JWT       │            │ - Items       │           │ - Templates   │
│ - 2FA       │            │ - Loans       │           │ - Storage     │
│ - RBAC      │            │ - Payments    │           │               │
└──────┬──────┘            │ - Sales       │           └───────┬───────┘
       │                   └───────┬───────┘                   │
       │                           │                           │
       └───────────────────────────┼───────────────────────────┘
                                   │
                    ┌──────────────▼──────────────┐
                    │         Message Queue       │
                    │       (Redis / NATS)        │
                    └──────────────┬──────────────┘
                                   │
       ┌───────────────────────────┼───────────────────────────┐
       │                           │                           │
┌──────▼──────┐            ┌───────▼───────┐           ┌───────▼───────┐
│  PostgreSQL │            │     Redis     │           │  Object Store │
│  (Primary)  │            │    (Cache)    │           │   (MinIO/S3)  │
├─────────────┤            └───────────────┘           └───────────────┘
│  PostgreSQL │
│  (Replica)  │
└─────────────┘
```

### 4.2 Arquitectura Simplificada (Monolito Modular)

Para una implementación inicial más práctica:

```
┌─────────────────────────────────────────────────────────────────┐
│                      Go Application                              │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   HTTP      │  │   Handlers  │  │  Middleware │             │
│  │   Router    │──│   (API)     │──│  (Auth,Log) │             │
│  └─────────────┘  └──────┬──────┘  └─────────────┘             │
│                          │                                      │
│  ┌───────────────────────┼───────────────────────┐             │
│  │                  Services Layer               │             │
│  ├─────────────┬─────────────┬─────────────┬─────┤             │
│  │  Customer   │    Loan     │   Payment   │ ... │             │
│  │  Service    │   Service   │   Service   │     │             │
│  └──────┬──────┴──────┬──────┴──────┬──────┴─────┘             │
│         │             │             │                           │
│  ┌──────▼─────────────▼─────────────▼──────┐                   │
│  │            Repository Layer              │                   │
│  └──────────────────┬───────────────────────┘                   │
│                     │                                           │
└─────────────────────┼───────────────────────────────────────────┘
                      │
        ┌─────────────┼─────────────┐
        │             │             │
   ┌────▼────┐   ┌────▼────┐   ┌────▼────┐
   │PostgreSQL│   │  Redis  │   │  MinIO  │
   └──────────┘   └─────────┘   └─────────┘
```

### 4.3 Estructura del Proyecto

```
pawnshop/
├── cmd/
│   ├── api/                    # Punto de entrada API
│   │   └── main.go
│   ├── worker/                 # Tareas en background
│   │   └── main.go
│   └── migrate/                # Herramienta de migraciones
│       └── main.go
├── internal/
│   ├── config/                 # Configuración
│   │   └── config.go
│   ├── domain/                 # Entidades de dominio
│   │   ├── customer.go
│   │   ├── item.go
│   │   ├── loan.go
│   │   ├── payment.go
│   │   ├── sale.go
│   │   └── user.go
│   ├── repository/             # Capa de datos
│   │   ├── postgres/
│   │   │   ├── customer_repo.go
│   │   │   ├── loan_repo.go
│   │   │   └── ...
│   │   └── interfaces.go
│   ├── service/                # Lógica de negocio
│   │   ├── customer_service.go
│   │   ├── loan_service.go
│   │   ├── payment_service.go
│   │   └── ...
│   ├── handler/                # Handlers HTTP
│   │   ├── customer_handler.go
│   │   ├── loan_handler.go
│   │   └── ...
│   ├── middleware/             # Middlewares
│   │   ├── auth.go
│   │   ├── logging.go
│   │   ├── ratelimit.go
│   │   └── cors.go
│   ├── pdf/                    # Generación de PDFs
│   │   ├── generator.go
│   │   └── templates/
│   └── scheduler/              # Tareas programadas
│       ├── interest_calculator.go
│       └── notification_sender.go
├── pkg/                        # Paquetes reutilizables
│   ├── auth/
│   │   ├── jwt.go
│   │   └── password.go
│   ├── validator/
│   │   └── validator.go
│   └── response/
│       └── json.go
├── migrations/                 # Migraciones SQL
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   └── ...
├── api/                        # Documentación OpenAPI
│   └── openapi.yaml
├── web/                        # Frontend (si aplica)
│   └── ...
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
├── scripts/
│   ├── build.sh
│   └── deploy.sh
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 5. Stack Tecnológico

### 5.1 Backend

| Componente | Tecnología | Justificación |
|------------|------------|---------------|
| Lenguaje | Go 1.22+ | Rendimiento, concurrencia, tipado |
| Framework HTTP | Fiber v2 / Echo v4 | Alto rendimiento, middleware ecosystem |
| ORM/Query Builder | GORM / sqlx | Productividad vs control |
| Migraciones | golang-migrate | Estándar de la industria |
| Validación | go-playground/validator | Completo, bien mantenido |
| JWT | golang-jwt | Estándar para auth |
| PDF | wkhtmltopdf + go wrapper | Calidad de renderizado HTML |
| Cache | go-redis | Cliente oficial Redis |
| Scheduler | robfig/cron | Cron jobs en Go |
| Logging | zerolog / zap | Estructurado, performante |
| Config | viper | Flexible, múltiples fuentes |

### 5.2 Base de Datos

| Componente | Tecnología | Justificación |
|------------|------------|---------------|
| Principal | PostgreSQL 16+ | ACID, JSON, extensible |
| Cache | Redis 7+ | Sesiones, cache, colas |
| Object Storage | MinIO / S3 | Documentos, imágenes |

### 5.3 Infraestructura

| Componente | Tecnología | Justificación |
|------------|------------|---------------|
| Contenedores | Docker | Portabilidad |
| Orquestación | Docker Compose / K8s | Desarrollo / Producción |
| Reverse Proxy | Traefik / Nginx | SSL, routing, LB |
| CI/CD | GitHub Actions | Integrado con repo |
| Monitoreo | Prometheus + Grafana | Métricas y alertas |
| Logs | Loki | Agregación de logs |

### 5.4 Frontend - Panel Administrativo

| Componente | Tecnología | Justificación |
|------------|------------|---------------|
| Framework | React 18 + Vite | Ecosistema maduro, rendimiento |
| UI Library | Tailwind CSS + shadcn/ui | Moderno, customizable, accesible |
| State Management | Zustand + TanStack Query | Simple, cache inteligente |
| Forms | React Hook Form + Zod | Validación tipada, rendimiento |
| Tables | TanStack Table | Paginación, filtros, ordenamiento |
| Charts | Recharts / Chart.js | Dashboards y reportes |
| Router | React Router v6 | Estándar de la industria |
| Auth | JWT + Refresh tokens | Sesiones seguras |

### 5.5 Frontend - Sistema POS (Caja) - **Wails**

| Componente | Tecnología | Justificación |
|------------|------------|---------------|
| Framework | **Wails v2 + React** | App nativa con backend Go compartido |
| Backend Desktop | **Go (mismo código que API)** | Reutilización 100% de lógica de negocio |
| UI | React + Tailwind + shadcn/ui | Mismos componentes que panel web |
| Impresión | **go-escpos** (nativo Go) | Acceso directo USB/Serial a térmicas |
| Offline Storage | **SQLite + Sync** | BD local con sincronización |
| Barcode Scanner | Keyboard wedge / Serial | Lectores estándar |
| Cash Drawer | **Serial/USB directo desde Go** | Sin dependencias JS |

#### ¿Por qué Wails en lugar de Electron?

| Aspecto | Electron | **Wails** |
|---------|----------|-----------|
| Tamaño instalador | 150-200 MB | **8-15 MB** |
| Consumo RAM | 200-400 MB | **30-80 MB** |
| Arranque | 3-5 segundos | **<1 segundo** |
| Backend | Node.js (duplicar código) | **Go (código compartido)** |
| Acceso hardware | Via Node addons | **Nativo en Go** |

#### Arquitectura Wails POS

```
┌─────────────────────────────────────────────────────────────┐
│                    WAILS DESKTOP APP                         │
├─────────────────────────────────────────────────────────────┤
│  ┌───────────────────────────────────────────────────────┐  │
│  │                 FRONTEND (React)                       │  │
│  │  • Componentes UI (shadcn/ui)                         │  │
│  │  • Estado local (Zustand)                             │  │
│  │  • Mismo código que panel web                         │  │
│  └─────────────────────────┬─────────────────────────────┘  │
│                            │ Wails Bindings                  │
│                            ▼                                 │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                 BACKEND (Go)                           │  │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐     │  │
│  │  │   Loan      │ │  Payment    │ │    Cash     │     │  │
│  │  │  Service    │ │  Service    │ │   Service   │     │  │
│  │  └──────┬──────┘ └──────┬──────┘ └──────┬──────┘     │  │
│  │         │               │               │             │  │
│  │         └───────────────┼───────────────┘             │  │
│  │                         ▼                             │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │            SERVICIOS COMPARTIDOS                │ │  │
│  │  │  (Mismo código que API del servidor)            │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                         │                             │  │
│  │         ┌───────────────┼───────────────┐            │  │
│  │         ▼               ▼               ▼            │  │
│  │  ┌───────────┐   ┌───────────┐   ┌───────────┐      │  │
│  │  │  SQLite   │   │  Printer  │   │   Sync    │      │  │
│  │  │  Local    │   │  Driver   │   │  Service  │      │  │
│  │  └───────────┘   └───────────┘   └───────────┘      │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                              │
│  Tamaño: ~12MB  │  RAM: ~50MB  │  Startup: <1s             │
└─────────────────────────────────────────────────────────────┘
```

---

## 5.6 Arquitectura de Interfaces de Usuario

### 5.6.1 Panel Administrativo (Web)

```
┌─────────────────────────────────────────────────────────────────────┐
│                    PANEL ADMINISTRATIVO                              │
├─────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐                                                   │
│  │   SIDEBAR    │   ┌─────────────────────────────────────────┐    │
│  │              │   │              HEADER                      │    │
│  │ • Dashboard  │   │  [Sucursal ▼]  [🔔]  [👤 Usuario ▼]     │    │
│  │ • Clientes   │   └─────────────────────────────────────────┘    │
│  │ • Artículos  │                                                   │
│  │ • Préstamos  │   ┌─────────────────────────────────────────┐    │
│  │ • Pagos      │   │                                         │    │
│  │ • Ventas     │   │            CONTENT AREA                 │    │
│  │ • Caja       │   │                                         │    │
│  │ ─────────    │   │   • Tablas con filtros y búsqueda      │    │
│  │ • Reportes   │   │   • Formularios modales                │    │
│  │ • Config     │   │   • Gráficos y estadísticas            │    │
│  │ • Usuarios   │   │                                         │    │
│  │              │   │                                         │    │
│  └──────────────┘   └─────────────────────────────────────────┘    │
│                                                                      │
│  Características:                                                    │
│  • Responsive (desktop, tablet)                                     │
│  • Tema claro/oscuro                                                │
│  • Atajos de teclado                                                │
│  • Exportación a Excel/PDF                                          │
└─────────────────────────────────────────────────────────────────────┘
```

### 5.6.2 Sistema POS / Caja (Desktop App)

```
┌─────────────────────────────────────────────────────────────────────┐
│  SISTEMA POS - CASA DE EMPEÑO              [Caja #1] [Juan Pérez]   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────────────────────┐  ┌────────────────────────────┐   │
│  │      OPERACIÓN ACTUAL       │  │      ACCESO RÁPIDO         │   │
│  │                             │  │                            │   │
│  │  Cliente: _______________   │  │  ┌────────┐  ┌────────┐   │   │
│  │  [🔍 Buscar]                │  │  │  NUEVO │  │ BUSCAR │   │   │
│  │                             │  │  │PRÉSTAMO│  │PRÉSTAMO│   │   │
│  │  ┌─────────────────────┐   │  │  └────────┘  └────────┘   │   │
│  │  │ Artículo: Laptop HP │   │  │                            │   │
│  │  │ Valor: Q 3,500.00   │   │  │  ┌────────┐  ┌────────┐   │   │
│  │  │ Préstamo: Q 2,500.00│   │  │  │ COBRAR │  │ NUEVA  │   │   │
│  │  └─────────────────────┘   │  │  │  PAGO  │  │ VENTA  │   │   │
│  │                             │  │  └────────┘  └────────┘   │   │
│  │  Interés: Q 250.00 (10%)   │  │                            │   │
│  │  ─────────────────────────  │  │  ┌────────┐  ┌────────┐   │   │
│  │  TOTAL: Q 2,750.00         │  │  │MOVIM.  │  │ CERRAR │   │   │
│  │                             │  │  │ CAJA   │  │  CAJA  │   │   │
│  │  [CONFIRMAR]  [CANCELAR]   │  │  └────────┘  └────────┘   │   │
│  │                             │  │                            │   │
│  └─────────────────────────────┘  └────────────────────────────┘   │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │  RESUMEN DE CAJA                                             │   │
│  │  ────────────────────────────────────────────────────────── │   │
│  │  Apertura: Q 500.00  │  Ingresos: Q 5,230.00  │  Egresos: Q 200.00  │
│  │  ────────────────────────────────────────────────────────── │   │
│  │  SALDO ACTUAL: Q 5,530.00                    [📄 IMPRIMIR]  │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                      │
│  [F1-Ayuda] [F2-Nuevo] [F3-Buscar] [F5-Pago] [F8-Imprimir] [F12-Salir] │
└─────────────────────────────────────────────────────────────────────┘

Características POS:
• Interfaz optimizada para velocidad (teclas de función)
• Botones grandes para pantalla táctil
• Impresión automática de tickets
• Funciona sin conexión a internet
• Sincronización automática cuando hay conexión
• Integración con impresora térmica 80mm
• Apertura automática de cajón de dinero
```

### 5.6.3 Flujo de Operaciones POS

```
┌──────────────────────────────────────────────────────────────────┐
│                    FLUJO DIARIO DE CAJA                          │
└──────────────────────────────────────────────────────────────────┘

    ┌─────────────┐
    │   INICIO    │
    │   DEL DÍA   │
    └──────┬──────┘
           │
           ▼
    ┌─────────────┐     ┌─────────────────────────────────┐
    │  APERTURA   │     │  • Verificar monto inicial      │
    │  DE CAJA    │────▶│  • Contar efectivo físico       │
    │             │     │  • Registrar diferencias        │
    └──────┬──────┘     └─────────────────────────────────┘
           │
           ▼
    ┌─────────────┐
    │ OPERACIONES │◄────────────────────────────────────┐
    │   DEL DÍA   │                                     │
    └──────┬──────┘                                     │
           │                                             │
     ┌─────┴─────┬─────────────┬─────────────┐         │
     ▼           ▼             ▼             ▼         │
┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │
│ NUEVO   │ │ COBRAR  │ │  NUEVA  │ │MOVIM.   │       │
│PRÉSTAMO │ │  PAGO   │ │  VENTA  │ │DE CAJA  │       │
└────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘       │
     │           │           │           │             │
     ▼           ▼           ▼           ▼             │
┌─────────────────────────────────────────────┐       │
│           IMPRESIÓN DE TICKET               │───────┘
│    (automático después de cada operación)   │
└─────────────────────────────────────────────┘
           │
           ▼
    ┌─────────────┐     ┌─────────────────────────────────┐
    │   CORTE X   │     │  • Resumen parcial sin cerrar   │
    │  (PARCIAL)  │────▶│  • Verificación de saldo        │
    │             │     │  • No afecta operaciones        │
    └──────┬──────┘     └─────────────────────────────────┘
           │
           ▼
    ┌─────────────┐     ┌─────────────────────────────────┐
    │  CIERRE DE  │     │  • Arqueo final obligatorio     │
    │  CAJA (Z)   │────▶│  • Registrar diferencias        │
    │             │     │  • Generar reporte del día      │
    └──────┬──────┘     │  • Bloquear más operaciones     │
           │            └─────────────────────────────────┘
           ▼
    ┌─────────────┐
    │   FIN DEL   │
    │     DÍA     │
    └─────────────┘
```

### 5.6.4 Pantallas del Sistema POS

| Pantalla | Descripción | Acceso Rápido |
|----------|-------------|---------------|
| **Apertura de Caja** | Ingreso de monto inicial, conteo de billetes | Al iniciar sesión |
| **Dashboard POS** | Vista principal con accesos rápidos | F1 |
| **Nuevo Préstamo** | Wizard: Cliente → Artículo → Condiciones → Confirmar | F2 |
| **Buscar Préstamo** | Por número, cliente, artículo | F3 |
| **Registrar Pago** | Seleccionar préstamo, ingresar monto | F5 |
| **Nueva Venta** | Seleccionar artículo, cliente opcional | F6 |
| **Movimientos de Caja** | Ingresos/egresos manuales | F7 |
| **Imprimir Último Ticket** | Reimpresión del último documento | F8 |
| **Corte X (Parcial)** | Resumen sin cerrar caja | F9 |
| **Corte Z (Cierre)** | Cierre definitivo del día | F10 |
| **Configuración** | Impresora, sonidos, atajos | F11 |
| **Cerrar Sesión** | Salir del sistema | F12 |

---

## 6. Diseño de Base de Datos

### 6.1 Diagrama ER Simplificado

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   branches   │     │    users     │     │    roles     │
├──────────────┤     ├──────────────┤     ├──────────────┤
│ id           │◄────│ branch_id    │     │ id           │
│ name         │     │ id           │────►│ name         │
│ address      │     │ name         │     │ permissions  │
│ phone        │     │ email        │     └──────────────┘
│ is_active    │     │ password     │
└──────────────┘     │ role_id      │
       │             └──────────────┘
       │
       ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  customers   │     │    items     │     │  categories  │
├──────────────┤     ├──────────────┤     ├──────────────┤
│ id           │◄────│ customer_id  │     │ id           │
│ branch_id    │     │ id           │────►│ name         │
│ first_name   │     │ category_id  │     │ slug         │
│ last_name    │     │ branch_id    │     └──────────────┘
│ identity_num │     │ name         │
│ phone        │     │ brand        │
│ email        │     │ serial_num   │
│ credit_limit │     │ appraised_val│
│ is_active    │     │ status       │
└──────────────┘     └──────────────┘
       │                    │
       │                    │
       ▼                    ▼
┌──────────────────────────────────────┐
│               loans                   │
├──────────────────────────────────────┤
│ id                                   │
│ loan_number (unique)                 │
│ customer_id ─────────────────────────│
│ item_id ─────────────────────────────│
│ branch_id                            │
│ loan_amount                          │
│ interest_rate                        │
│ interest_amount                      │
│ principal_remaining                  │
│ total_amount                         │
│ start_date                           │
│ due_date                             │
│ status                               │
│ requires_minimum_payment             │
│ minimum_monthly_payment              │
└──────────────────────────────────────┘
       │
       │
       ▼
┌──────────────┐     ┌──────────────┐
│   payments   │     │    sales     │
├──────────────┤     ├──────────────┤
│ id           │     │ id           │
│ payment_num  │     │ sale_number  │
│ loan_id      │     │ item_id      │
│ branch_id    │     │ customer_id  │
│ amount       │     │ branch_id    │
│ principal_amt│     │ sale_price   │
│ interest_amt │     │ discount     │
│ payment_date │     │ final_price  │
│ method       │     │ sale_date    │
│ status       │     │ status       │
└──────────────┘     └──────────────┘
```

### 6.2 Diagrama de Caja/POS

```
┌──────────────────┐
│   cash_registers │
├──────────────────┤
│ id               │
│ branch_id        │───────┐
│ name             │       │
│ is_active        │       │
└────────┬─────────┘       │
         │                 │
         ▼                 │
┌──────────────────┐       │
│  cash_sessions   │       │
├──────────────────┤       │
│ id               │       │
│ cash_register_id │◄──────┘
│ user_id          │
│ opening_amount   │
│ closing_amount   │
│ expected_amount  │
│ difference       │
│ status           │ (open, closed)
│ opened_at        │
│ closed_at        │
│ notes            │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐     ┌──────────────────┐
│ cash_movements   │     │   Relacionado    │
├──────────────────┤     ├──────────────────┤
│ id               │     │ • payments       │
│ cash_session_id  │     │ • sales          │
│ type             │ ◄───│ • loans          │
│ amount           │     │   (desembolsos)  │
│ payment_method   │     └──────────────────┘
│ reference_type   │
│ reference_id     │
│ description      │
│ created_by       │
│ created_at       │
└──────────────────┘

Tipos de movimiento (type):
• income_loan_disbursement  - Desembolso de préstamo (egreso)
• income_payment           - Cobro de pago (ingreso)
• income_sale              - Venta (ingreso)
• income_other             - Otros ingresos
• expense_return           - Devolución
• expense_supplier         - Pago a proveedor
• expense_other            - Otros egresos
• adjustment_positive      - Ajuste positivo
• adjustment_negative      - Ajuste negativo
```

### 6.3 Ejemplo de Migración

```sql
-- migrations/000003_create_loans.up.sql

CREATE TYPE loan_status AS ENUM (
    'active',
    'paid',
    'overdue',
    'defaulted',
    'renewed',
    'confiscated'
);

CREATE TYPE payment_plan_type AS ENUM (
    'single',
    'minimum_payment',
    'installments'
);

CREATE TABLE loans (
    id              BIGSERIAL PRIMARY KEY,
    loan_number     VARCHAR(50) NOT NULL UNIQUE,
    customer_id     BIGINT NOT NULL REFERENCES customers(id),
    item_id         BIGINT NOT NULL REFERENCES items(id),
    branch_id       BIGINT NOT NULL REFERENCES branches(id),
    created_by      BIGINT REFERENCES users(id),

    -- Amounts
    loan_amount         DECIMAL(12,2) NOT NULL CHECK (loan_amount > 0),
    interest_rate       DECIMAL(5,2) NOT NULL CHECK (interest_rate >= 0),
    interest_amount     DECIMAL(12,2) NOT NULL DEFAULT 0,
    principal_remaining DECIMAL(12,2) NOT NULL,
    total_amount        DECIMAL(12,2) NOT NULL,
    amount_paid         DECIMAL(12,2) NOT NULL DEFAULT 0,

    -- Dates
    start_date  DATE NOT NULL DEFAULT CURRENT_DATE,
    due_date    DATE,

    -- Payment plan
    payment_plan_type           payment_plan_type NOT NULL DEFAULT 'minimum_payment',
    loan_term_days              INTEGER,
    requires_minimum_payment    BOOLEAN NOT NULL DEFAULT false,
    minimum_monthly_payment     DECIMAL(12,2),
    next_minimum_payment_date   DATE,
    grace_period_days           INTEGER DEFAULT 5,

    -- Status
    status      loan_status NOT NULL DEFAULT 'active',
    notes       TEXT,

    -- Timestamps
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,

    -- Indexes
    CONSTRAINT loans_item_unique CHECK (
        status != 'active' OR item_id NOT IN (
            SELECT item_id FROM loans WHERE status = 'active' AND id != loans.id
        )
    )
);

CREATE INDEX idx_loans_customer ON loans(customer_id);
CREATE INDEX idx_loans_branch ON loans(branch_id);
CREATE INDEX idx_loans_status ON loans(status);
CREATE INDEX idx_loans_due_date ON loans(due_date) WHERE status = 'active';
CREATE INDEX idx_loans_number ON loans(loan_number);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER loans_updated_at
    BEFORE UPDATE ON loans
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();
```

---

## 7. Diseño de API

### 7.1 Endpoints Principales

```yaml
# Autenticación
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
POST   /api/v1/auth/refresh
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password

# Usuarios
GET    /api/v1/users
POST   /api/v1/users
GET    /api/v1/users/{id}
PUT    /api/v1/users/{id}
DELETE /api/v1/users/{id}
GET    /api/v1/users/me

# Clientes
GET    /api/v1/customers
POST   /api/v1/customers
GET    /api/v1/customers/{id}
PUT    /api/v1/customers/{id}
DELETE /api/v1/customers/{id}
GET    /api/v1/customers/{id}/loans
GET    /api/v1/customers/{id}/payments

# Artículos
GET    /api/v1/items
POST   /api/v1/items
GET    /api/v1/items/{id}
PUT    /api/v1/items/{id}
DELETE /api/v1/items/{id}
POST   /api/v1/items/{id}/transfer

# Préstamos
GET    /api/v1/loans
POST   /api/v1/loans
GET    /api/v1/loans/{id}
PUT    /api/v1/loans/{id}
POST   /api/v1/loans/{id}/renew
POST   /api/v1/loans/{id}/confiscate
GET    /api/v1/loans/{id}/payments
GET    /api/v1/loans/{id}/installments

# Pagos
GET    /api/v1/payments
POST   /api/v1/payments
GET    /api/v1/payments/{id}
POST   /api/v1/payments/{id}/reverse

# Ventas
GET    /api/v1/sales
POST   /api/v1/sales
GET    /api/v1/sales/{id}

# Documentos
GET    /api/v1/documents/loan-contract/{loanId}
GET    /api/v1/documents/loan-receipt/{loanId}
GET    /api/v1/documents/payment-receipt/{paymentId}
GET    /api/v1/documents/sale-receipt/{saleId}

# Caja / POS
GET    /api/v1/cash-registers
POST   /api/v1/cash-registers
GET    /api/v1/cash-registers/{id}
PUT    /api/v1/cash-registers/{id}

# Sesiones de Caja
GET    /api/v1/cash-sessions
POST   /api/v1/cash-sessions/open              # Apertura de caja
POST   /api/v1/cash-sessions/{id}/close        # Cierre de caja (Corte Z)
GET    /api/v1/cash-sessions/{id}
GET    /api/v1/cash-sessions/{id}/summary      # Corte X (parcial)
GET    /api/v1/cash-sessions/{id}/movements
GET    /api/v1/cash-sessions/current           # Sesión activa del usuario

# Movimientos de Caja
GET    /api/v1/cash-movements
POST   /api/v1/cash-movements                  # Ingreso/egreso manual
GET    /api/v1/cash-movements/{id}

# Reportes
GET    /api/v1/reports/daily-summary
GET    /api/v1/reports/loans-by-status
GET    /api/v1/reports/overdue-loans
GET    /api/v1/reports/revenue
GET    /api/v1/reports/cash-flow              # Flujo de caja
GET    /api/v1/reports/cash-by-branch         # Resumen por sucursal
GET    /api/v1/reports/sales-by-period
GET    /api/v1/reports/interest-earned

# Impresión
POST   /api/v1/print/ticket                    # Enviar a impresora térmica
GET    /api/v1/print/preview/{type}/{id}       # Vista previa

# Configuración
GET    /api/v1/settings
PUT    /api/v1/settings
GET    /api/v1/branches
POST   /api/v1/branches
GET    /api/v1/printers
POST   /api/v1/printers
```

### 7.2 Ejemplo de Request/Response

```json
// POST /api/v1/loans
// Request
{
  "customer_id": 123,
  "item_id": 456,
  "loan_amount": 2500.00,
  "interest_rate": 10.00,
  "loan_term_days": 30,
  "payment_plan_type": "minimum_payment",
  "requires_minimum_payment": true,
  "minimum_monthly_payment": 275.00,
  "grace_period_days": 5
}

// Response 201 Created
{
  "data": {
    "id": 789,
    "loan_number": "LN-2024-0789",
    "customer": {
      "id": 123,
      "full_name": "Juan Carlos Pérez"
    },
    "item": {
      "id": 456,
      "name": "Laptop HP Pavilion",
      "appraised_value": 3500.00
    },
    "loan_amount": 2500.00,
    "interest_rate": 10.00,
    "interest_amount": 250.00,
    "total_amount": 2750.00,
    "principal_remaining": 2500.00,
    "start_date": "2024-01-15",
    "due_date": "2024-02-14",
    "status": "active",
    "payment_plan_type": "minimum_payment",
    "requires_minimum_payment": true,
    "minimum_monthly_payment": 275.00,
    "next_minimum_payment_date": "2024-02-14",
    "created_at": "2024-01-15T10:30:00Z"
  },
  "meta": {
    "documents": {
      "contract": "/api/v1/documents/loan-contract/789",
      "receipt": "/api/v1/documents/loan-receipt/789"
    }
  }
}
```

### 7.3 Códigos de Error

```json
// 400 Bad Request - Validación
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Los datos proporcionados no son válidos",
    "details": [
      {
        "field": "loan_amount",
        "message": "El monto debe ser mayor a 0"
      },
      {
        "field": "customer_id",
        "message": "El cliente no existe"
      }
    ]
  }
}

// 401 Unauthorized
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Token de autenticación inválido o expirado"
  }
}

// 403 Forbidden
{
  "error": {
    "code": "FORBIDDEN",
    "message": "No tiene permisos para realizar esta acción"
  }
}

// 404 Not Found
{
  "error": {
    "code": "NOT_FOUND",
    "message": "El recurso solicitado no existe"
  }
}

// 409 Conflict
{
  "error": {
    "code": "CONFLICT",
    "message": "El artículo ya está en un préstamo activo"
  }
}
```

---

## 8. Estándares y Mejores Prácticas

### 8.1 Estándares de Código

#### 8.1.1 Go - Convenciones
| Estándar | Descripción |
|----------|-------------|
| **gofmt** | Formateo automático obligatorio |
| **golint** | Linting estricto |
| **go vet** | Análisis estático |
| **errcheck** | Verificar errores no manejados |
| **staticcheck** | Análisis avanzado |
| **Effective Go** | Guía oficial de estilo |

#### 8.1.2 Nomenclatura
```go
// Packages: minúsculas, singular
package customer  // ✓
package customers // ✗

// Interfaces: sufijo -er cuando sea posible
type Reader interface { ... }
type CustomerRepository interface { ... }

// Structs: PascalCase
type Customer struct { ... }

// Métodos públicos: PascalCase
func (c *Customer) FullName() string { ... }

// Métodos privados: camelCase
func (c *Customer) calculateScore() int { ... }

// Constantes: PascalCase o SCREAMING_SNAKE_CASE
const MaxRetries = 3
const DEFAULT_TIMEOUT = 30
```

#### 8.1.3 Estructura de Archivos
```
// Un archivo por entidad principal
customer.go         // Struct + métodos
customer_test.go    // Tests
customer_mock.go    // Mocks para tests

// Archivos auxiliares con sufijo descriptivo
customer_validation.go
customer_repository.go
```

### 8.2 Estándares de API (REST)

#### 8.2.1 Convenciones URL
```
# Recursos en plural, kebab-case
GET  /api/v1/customers
GET  /api/v1/cash-registers
POST /api/v1/loan-renewals

# Anidamiento máximo 2 niveles
GET /api/v1/customers/{id}/loans     ✓
GET /api/v1/customers/{id}/loans/{id}/payments  ✗ (usar /payments?loan_id=X)

# Acciones como sub-recursos
POST /api/v1/loans/{id}/renew
POST /api/v1/cash-sessions/{id}/close
```

#### 8.2.2 Formato de Respuesta (JSON:API inspirado)
```json
// Respuesta exitosa
{
  "data": { ... },
  "meta": {
    "timestamp": "2024-01-15T10:30:00Z",
    "request_id": "abc-123"
  }
}

// Respuesta paginada
{
  "data": [ ... ],
  "meta": {
    "current_page": 1,
    "per_page": 20,
    "total_items": 150,
    "total_pages": 8
  },
  "links": {
    "first": "/api/v1/customers?page=1",
    "prev": null,
    "next": "/api/v1/customers?page=2",
    "last": "/api/v1/customers?page=8"
  }
}

// Respuesta de error
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Datos inválidos",
    "details": [ ... ]
  },
  "meta": {
    "timestamp": "...",
    "request_id": "..."
  }
}
```

#### 8.2.3 Códigos HTTP
| Código | Uso |
|--------|-----|
| 200 | GET exitoso, PUT/PATCH exitoso |
| 201 | POST exitoso (recurso creado) |
| 204 | DELETE exitoso (sin contenido) |
| 400 | Error de validación |
| 401 | No autenticado |
| 403 | Sin permisos |
| 404 | Recurso no encontrado |
| 409 | Conflicto (duplicado, estado inválido) |
| 422 | Entidad no procesable |
| 429 | Rate limit excedido |
| 500 | Error interno del servidor |

### 8.3 Estándares de Base de Datos

#### 8.3.1 Nomenclatura SQL
```sql
-- Tablas: plural, snake_case
CREATE TABLE customers ( ... );
CREATE TABLE cash_registers ( ... );

-- Columnas: snake_case
customer_id, created_at, is_active

-- Foreign keys: tabla_singular_id
customer_id, loan_id, branch_id

-- Índices: idx_tabla_columna(s)
CREATE INDEX idx_loans_customer_id ON loans(customer_id);
CREATE INDEX idx_loans_status_due ON loans(status, due_date);

-- Constraints: tabla_columna_tipo
CONSTRAINT customers_email_unique UNIQUE(email)
CONSTRAINT loans_amount_positive CHECK(loan_amount > 0)
```

#### 8.3.2 Campos Obligatorios
```sql
-- Toda tabla debe tener:
id              BIGSERIAL PRIMARY KEY,
created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

-- Tablas con soft delete:
deleted_at      TIMESTAMPTZ
```

### 8.4 Estándares de Testing

#### 8.4.1 Cobertura Mínima
| Tipo | Cobertura | Descripción |
|------|-----------|-------------|
| Unit Tests | 80%+ | Lógica de negocio |
| Integration Tests | 60%+ | Repositorios, handlers |
| E2E Tests | Flujos críticos | Préstamo completo, pago, cierre caja |

#### 8.4.2 Nomenclatura de Tests
```go
// Formato: Test<Función>_<Escenario>_<Resultado>
func TestCreateLoan_ValidData_ReturnsLoan(t *testing.T) { ... }
func TestCreateLoan_InvalidCustomer_ReturnsError(t *testing.T) { ... }
func TestCreateLoan_ItemNotAvailable_ReturnsConflict(t *testing.T) { ... }
```

### 8.5 Estándares de Documentación

#### 8.5.1 Código
```go
// Package customer maneja la lógica de negocio de clientes.
// Incluye creación, actualización, validación de crédito
// y gestión del historial crediticio.
package customer

// Customer representa un cliente del sistema de empeño.
// Contiene información personal, de contacto y crediticia.
type Customer struct {
    // ID es el identificador único del cliente.
    ID int64 `json:"id"`

    // CreditLimit es el monto máximo que puede solicitar.
    // Se calcula basado en el historial crediticio.
    CreditLimit float64 `json:"credit_limit"`
}

// CalculateCreditScore evalúa el historial del cliente
// y retorna un score entre 0 y 100.
//
// El cálculo considera:
//   - Préstamos pagados a tiempo
//   - Monto total de operaciones
//   - Antigüedad como cliente
func (c *Customer) CalculateCreditScore() int {
    // ...
}
```

#### 8.5.2 API (OpenAPI/Swagger)
```yaml
/api/v1/loans:
  post:
    summary: Crear nuevo préstamo
    description: |
      Crea un nuevo préstamo prendario. El artículo debe estar
      disponible y el cliente activo con crédito suficiente.
    tags:
      - Préstamos
    security:
      - bearerAuth: []
    requestBody:
      required: true
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/CreateLoanRequest'
    responses:
      '201':
        description: Préstamo creado exitosamente
      '400':
        description: Datos de entrada inválidos
      '409':
        description: Artículo no disponible
```

### 8.6 Estándares de Git

#### 8.6.1 Commits (Conventional Commits)
```
<tipo>(<ámbito>): <descripción>

[cuerpo opcional]

[footer opcional]
```

| Tipo | Uso |
|------|-----|
| feat | Nueva funcionalidad |
| fix | Corrección de bug |
| docs | Documentación |
| style | Formateo (no afecta lógica) |
| refactor | Refactorización |
| test | Tests |
| chore | Tareas de mantenimiento |

```bash
# Ejemplos
feat(loans): agregar cálculo de interés compuesto
fix(payments): corregir aplicación de pagos parciales
docs(api): documentar endpoints de caja
refactor(customer): extraer validación a servicio
```

#### 8.6.2 Branches
```
main              # Producción
develop           # Desarrollo
feature/POS-123   # Nueva funcionalidad
bugfix/POS-456    # Corrección
hotfix/POS-789    # Corrección urgente producción
release/v1.2.0    # Preparación de release
```

### 8.7 Estándares de Seguridad (OWASP)

| Vulnerabilidad | Mitigación |
|----------------|------------|
| SQL Injection | Prepared statements, ORM |
| XSS | Escape de output, CSP headers |
| CSRF | Tokens en formularios |
| Broken Auth | JWT seguro, rate limiting |
| Sensitive Data | Encriptación, HTTPS |
| Security Misconfiguration | Headers seguros, sin defaults |
| Insufficient Logging | Audit logs completos |

---

## 9. Seguridad

### 9.1 Autenticación

```go
// pkg/auth/jwt.go
type JWTClaims struct {
    UserID    int64    `json:"user_id"`
    Email     string   `json:"email"`
    RoleID    int64    `json:"role_id"`
    BranchID  int64    `json:"branch_id"`
    Permissions []string `json:"permissions"`
    jwt.RegisteredClaims
}

// Access token: 15 minutos
// Refresh token: 7 días (rotación obligatoria)
```

### 9.2 Contraseñas

```go
// pkg/auth/password.go
// Argon2id con parámetros seguros
func HashPassword(password string) (string, error) {
    return argon2id.CreateHash(password, &argon2id.Params{
        Memory:      64 * 1024,
        Iterations:  3,
        Parallelism: 2,
        SaltLength:  16,
        KeyLength:   32,
    })
}
```

### 9.3 Rate Limiting

```go
// Por IP: 100 requests/minuto (general)
// Por Usuario: 1000 requests/minuto
// Login: 5 intentos/15 minutos
// Generación PDF: 10/minuto
```

### 9.4 Headers de Seguridad

```go
// middleware/security.go
func SecurityHeaders() fiber.Handler {
    return func(c *fiber.Ctx) error {
        c.Set("X-Content-Type-Options", "nosniff")
        c.Set("X-Frame-Options", "DENY")
        c.Set("X-XSS-Protection", "1; mode=block")
        c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Set("Content-Security-Policy", "default-src 'self'")
        c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
        return c.Next()
    }
}
```

---

## 10. Implementación por Fases

### Fase 1: Fundamentos (4-5 semanas)
| Semana | Entregables |
|--------|-------------|
| 1 | Estructura proyecto, config, conexión BD |
| 2 | Migraciones, modelos base, repositorios |
| 3 | Autenticación JWT, middleware auth |
| 4 | CRUD usuarios, roles, permisos |
| 5 | Middleware (logging, cors, rate limit), manejo errores |

**Criterios de aceptación:**
- [ ] Login/logout funcional
- [ ] CRUD usuarios completo
- [ ] Tests unitarios >80%
- [ ] API documentada (Swagger)

### Fase 2: Core del Negocio (6-7 semanas)
| Semana | Entregables |
|--------|-------------|
| 6 | CRUD sucursales, categorías |
| 7 | CRUD clientes con validaciones |
| 8 | CRUD artículos, estados, transferencias |
| 9-10 | Gestión de préstamos completa |
| 11 | Sistema de pagos, aplicación a intereses/capital |
| 12 | Ventas, estados, transiciones |

**Criterios de aceptación:**
- [ ] Flujo completo: Cliente → Artículo → Préstamo → Pago
- [ ] Cálculo de intereses correcto
- [ ] Tests de integración
- [ ] Validaciones de negocio

### Fase 3: Sistema POS / Caja (4-5 semanas)
| Semana | Entregables |
|--------|-------------|
| 13 | Backend: Cajas, sesiones, movimientos |
| 14 | Backend: Corte X, Corte Z, reportes caja |
| 15 | Frontend POS: Estructura Electron/React |
| 16 | Frontend POS: Operaciones principales |
| 17 | Integración impresora térmica, offline mode |

**Criterios de aceptación:**
- [ ] Apertura/cierre de caja funcional
- [ ] Impresión de tickets en térmica 80mm
- [ ] Funciona sin conexión a internet
- [ ] Sincronización automática

### Fase 4: Documentos y Reportes (3-4 semanas)
| Semana | Entregables |
|--------|-------------|
| 18 | Generación de PDFs (contratos, recibos) |
| 19 | Plantillas personalizables, branding |
| 20 | Reportes: préstamos, pagos, ventas |
| 21 | Dashboard con gráficos, KPIs |

**Criterios de aceptación:**
- [ ] Contratos PDF generados correctamente
- [ ] Tickets térmicos formateados
- [ ] Al menos 10 reportes funcionales
- [ ] Dashboard con datos en tiempo real

### Fase 5: Características Avanzadas (3-4 semanas)
| Semana | Entregables |
|--------|-------------|
| 22 | Tareas programadas (intereses, vencimientos) |
| 23 | Notificaciones (email, SMS básico) |
| 24 | Auditoría completa, logs inmutables |
| 25 | Renovaciones, confiscaciones, flujos especiales |

**Criterios de aceptación:**
- [ ] Cálculo automático de intereses diario
- [ ] Alertas de vencimiento funcionando
- [ ] Historial de auditoría completo

### Fase 6: Panel Administrativo Web (4-5 semanas)
| Semana | Entregables |
|--------|-------------|
| 26 | Estructura React, autenticación |
| 27 | Módulos: Clientes, Artículos |
| 28 | Módulos: Préstamos, Pagos, Ventas |
| 29 | Módulos: Caja, Configuración |
| 30 | Dashboard, reportes, gráficos |

**Criterios de aceptación:**
- [ ] CRUD completo de todas las entidades
- [ ] Responsive (desktop, tablet)
- [ ] Filtros, búsqueda, exportación

### Fase 7: QA, Optimización y Deploy (3-4 semanas)
| Semana | Entregables |
|--------|-------------|
| 31 | Tests E2E, corrección de bugs |
| 32 | Optimización, caché Redis, índices |
| 33 | Docker, CI/CD, documentación deploy |
| 34 | Monitoreo (Prometheus/Grafana), alertas |

**Criterios de aceptación:**
- [ ] Tests E2E pasando
- [ ] Tiempo respuesta <100ms (p95)
- [ ] Deploy automatizado funcionando
- [ ] Alertas configuradas

---

### Resumen de Timeline

```
┌─────────────────────────────────────────────────────────────────┐
│                    CRONOGRAMA DEL PROYECTO                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Mes 1    Mes 2    Mes 3    Mes 4    Mes 5    Mes 6    Mes 7   │
│  ────────────────────────────────────────────────────────────   │
│  ▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   │
│  Fase 1: Fundamentos                                             │
│                                                                  │
│  ░░░░░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░   │
│           Fase 2: Core del Negocio                               │
│                                                                  │
│  ░░░░░░░░░░░░░░░░░░░░░░▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░░░░░░░░░   │
│                         Fase 3: Sistema POS                      │
│                                                                  │
│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓▓▓▓▓▓▓▓░░░░░░░░░░░░░░░░░░   │
│                                   Fase 4: Documentos             │
│                                                                  │
│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓▓▓▓▓▓▓▓░░░░░░░░░░   │
│                                           Fase 5: Avanzado       │
│                                                                  │
│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓▓▓▓▓▓▓▓▓▓   │
│                                                   Fase 6: Admin  │
│                                                                  │
│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░▓▓▓▓▓   │
│                                                         Fase 7   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘

Total: 30-34 semanas (7-8 meses)
```

### MVP (Producto Mínimo Viable) - 16 semanas

Si se requiere un MVP más rápido, se puede entregar en **16 semanas** con:
- ✅ Fases 1-3 completas (Backend + POS básico)
- ✅ PDFs básicos (sin plantillas personalizables)
- ⏸️ Panel Admin simplificado
- ⏸️ Reportes básicos
- ❌ Notificaciones (posterior)
- ❌ Dashboard avanzado (posterior)

---

## 10. Recursos Necesarios

### 10.1 Equipo de Desarrollo

| Rol | Cantidad | Dedicación |
|-----|----------|------------|
| Tech Lead / Arquitecto | 1 | 100% |
| Backend Developer Sr | 1-2 | 100% |
| Backend Developer Jr/Mid | 1 | 100% |
| Frontend Developer (si aplica) | 1 | 50-100% |
| DevOps / SRE | 1 | 25-50% |
| QA Engineer | 1 | 50% |

### 10.2 Infraestructura (Producción)

| Componente | Especificación | Costo Mensual Est. |
|------------|----------------|-------------------|
| Servidor API | 2 vCPU, 4GB RAM | $40-80 |
| PostgreSQL | 2 vCPU, 4GB RAM, 50GB SSD | $50-100 |
| Redis | 1GB RAM | $15-30 |
| Object Storage | 50GB | $5-10 |
| Backup Storage | 100GB | $5-10 |
| **Total** | | **$115-230/mes** |

### 10.3 Herramientas y Licencias

| Herramienta | Propósito | Costo |
|-------------|-----------|-------|
| GitHub | Repositorio | Gratis (público) / $4/user (privado) |
| GitHub Actions | CI/CD | Gratis (2000 min/mes) |
| Sentry | Error tracking | Gratis (5k eventos/mes) |
| Grafana Cloud | Monitoreo | Gratis (tier básico) |

---

## 11. Riesgos y Mitigación

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|--------------|---------|------------|
| Curva de aprendizaje Go | Media | Medio | Capacitación previa, pair programming |
| Migración de datos | Alta | Alto | Scripts de migración, validación, rollback plan |
| Compatibilidad PDF | Media | Medio | Pruebas exhaustivas, fallback a librería alternativa |
| Tiempo de desarrollo | Media | Alto | Fases bien definidas, MVPs incrementales |
| Rendimiento real vs esperado | Baja | Medio | Benchmarks tempranos, optimización continua |

---

## 12. Conclusiones

### Ventajas de la Migración
1. **Rendimiento**: 50x mejor throughput
2. **Costos**: Menor consumo de recursos = menor costo de infra
3. **Mantenibilidad**: Tipado estático, menos errores en runtime
4. **Despliegue**: Binario único, sin dependencias
5. **Escalabilidad**: Diseño preparado para crecer

### Consideraciones
1. Requiere inversión inicial en tiempo y recursos
2. Equipo debe conocer o aprender Go
3. Migración de datos debe ser cuidadosa
4. Período de transición donde ambos sistemas coexisten

### Recomendación
Proceder con la migración en fases, comenzando con un MVP que cubra las funcionalidades críticas (préstamos, pagos, clientes), para validar el enfoque antes de migrar características secundarias.

---

## Apéndice A: Ejemplo de Código

### Handler de Préstamos

```go
// internal/handler/loan_handler.go
package handler

import (
    "github.com/gofiber/fiber/v2"
    "pawnshop/internal/domain"
    "pawnshop/internal/service"
    "pawnshop/pkg/response"
)

type LoanHandler struct {
    loanService service.LoanService
}

func NewLoanHandler(ls service.LoanService) *LoanHandler {
    return &LoanHandler{loanService: ls}
}

type CreateLoanRequest struct {
    CustomerID              int64   `json:"customer_id" validate:"required"`
    ItemID                  int64   `json:"item_id" validate:"required"`
    LoanAmount              float64 `json:"loan_amount" validate:"required,gt=0"`
    InterestRate            float64 `json:"interest_rate" validate:"required,gte=0,lte=100"`
    LoanTermDays            int     `json:"loan_term_days" validate:"required,gt=0"`
    PaymentPlanType         string  `json:"payment_plan_type" validate:"required,oneof=single minimum_payment installments"`
    RequiresMinimumPayment  bool    `json:"requires_minimum_payment"`
    MinimumMonthlyPayment   float64 `json:"minimum_monthly_payment" validate:"gte=0"`
    GracePeriodDays         int     `json:"grace_period_days" validate:"gte=0,lte=30"`
}

func (h *LoanHandler) Create(c *fiber.Ctx) error {
    var req CreateLoanRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body")
    }

    if err := validate.Struct(req); err != nil {
        return response.ValidationError(c, err)
    }

    // Get user from context (set by auth middleware)
    user := c.Locals("user").(*domain.User)

    loan, err := h.loanService.Create(c.Context(), service.CreateLoanInput{
        CustomerID:             req.CustomerID,
        ItemID:                 req.ItemID,
        BranchID:               user.BranchID,
        CreatedBy:              user.ID,
        LoanAmount:             req.LoanAmount,
        InterestRate:           req.InterestRate,
        LoanTermDays:           req.LoanTermDays,
        PaymentPlanType:        req.PaymentPlanType,
        RequiresMinimumPayment: req.RequiresMinimumPayment,
        MinimumMonthlyPayment:  req.MinimumMonthlyPayment,
        GracePeriodDays:        req.GracePeriodDays,
    })

    if err != nil {
        return response.HandleServiceError(c, err)
    }

    return response.Created(c, loan)
}

func (h *LoanHandler) GetByID(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil {
        return response.BadRequest(c, "Invalid loan ID")
    }

    loan, err := h.loanService.GetByID(c.Context(), int64(id))
    if err != nil {
        return response.HandleServiceError(c, err)
    }

    return response.OK(c, loan)
}

func (h *LoanHandler) List(c *fiber.Ctx) error {
    user := c.Locals("user").(*domain.User)

    params := service.ListLoansParams{
        BranchID: user.BranchID,
        Status:   c.Query("status"),
        Page:     c.QueryInt("page", 1),
        PerPage:  c.QueryInt("per_page", 20),
    }

    result, err := h.loanService.List(c.Context(), params)
    if err != nil {
        return response.HandleServiceError(c, err)
    }

    return response.Paginated(c, result.Loans, result.Meta)
}
```

### Servicio de Préstamos

```go
// internal/service/loan_service.go
package service

import (
    "context"
    "errors"
    "fmt"
    "time"

    "pawnshop/internal/domain"
    "pawnshop/internal/repository"
)

type LoanService interface {
    Create(ctx context.Context, input CreateLoanInput) (*domain.Loan, error)
    GetByID(ctx context.Context, id int64) (*domain.Loan, error)
    List(ctx context.Context, params ListLoansParams) (*ListLoansResult, error)
    ApplyPayment(ctx context.Context, loanID int64, amount float64) (*domain.Payment, error)
}

type loanService struct {
    loanRepo     repository.LoanRepository
    itemRepo     repository.ItemRepository
    customerRepo repository.CustomerRepository
    paymentRepo  repository.PaymentRepository
}

func NewLoanService(
    lr repository.LoanRepository,
    ir repository.ItemRepository,
    cr repository.CustomerRepository,
    pr repository.PaymentRepository,
) LoanService {
    return &loanService{
        loanRepo:     lr,
        itemRepo:     ir,
        customerRepo: cr,
        paymentRepo:  pr,
    }
}

func (s *loanService) Create(ctx context.Context, input CreateLoanInput) (*domain.Loan, error) {
    // Validate customer exists and is active
    customer, err := s.customerRepo.GetByID(ctx, input.CustomerID)
    if err != nil {
        return nil, fmt.Errorf("customer not found: %w", err)
    }
    if !customer.IsActive {
        return nil, errors.New("customer is not active")
    }

    // Validate item exists and is available
    item, err := s.itemRepo.GetByID(ctx, input.ItemID)
    if err != nil {
        return nil, fmt.Errorf("item not found: %w", err)
    }
    if item.Status != domain.ItemStatusAvailable {
        return nil, errors.New("item is not available for loan")
    }

    // Validate loan amount against item value
    if input.LoanAmount > item.AppraisedValue {
        return nil, errors.New("loan amount cannot exceed appraised value")
    }

    // Calculate interest
    interestAmount := input.LoanAmount * (input.InterestRate / 100)
    totalAmount := input.LoanAmount + interestAmount

    // Generate loan number
    loanNumber, err := s.loanRepo.GenerateNumber(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to generate loan number: %w", err)
    }

    // Create loan
    loan := &domain.Loan{
        LoanNumber:             loanNumber,
        CustomerID:             input.CustomerID,
        ItemID:                 input.ItemID,
        BranchID:               input.BranchID,
        CreatedBy:              input.CreatedBy,
        LoanAmount:             input.LoanAmount,
        InterestRate:           input.InterestRate,
        InterestAmount:         interestAmount,
        PrincipalRemaining:     input.LoanAmount,
        TotalAmount:            totalAmount,
        StartDate:              time.Now(),
        DueDate:                time.Now().AddDate(0, 0, input.LoanTermDays),
        Status:                 domain.LoanStatusActive,
        PaymentPlanType:        input.PaymentPlanType,
        LoanTermDays:           input.LoanTermDays,
        RequiresMinimumPayment: input.RequiresMinimumPayment,
        MinimumMonthlyPayment:  input.MinimumMonthlyPayment,
        GracePeriodDays:        input.GracePeriodDays,
    }

    if loan.RequiresMinimumPayment {
        nextPaymentDate := time.Now().AddDate(0, 0, 30)
        loan.NextMinimumPaymentDate = &nextPaymentDate
    }

    // Start transaction
    tx, err := s.loanRepo.BeginTx(ctx)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    // Save loan
    if err := s.loanRepo.CreateTx(ctx, tx, loan); err != nil {
        return nil, fmt.Errorf("failed to create loan: %w", err)
    }

    // Update item status to collateral
    if err := s.itemRepo.UpdateStatusTx(ctx, tx, item.ID, domain.ItemStatusCollateral); err != nil {
        return nil, fmt.Errorf("failed to update item status: %w", err)
    }

    if err := tx.Commit(); err != nil {
        return nil, err
    }

    // Load relations for response
    loan.Customer = customer
    loan.Item = item

    return loan, nil
}
```

---

## Apéndice B: Resumen de Requerimientos

### Tabla Resumen de Módulos

| Módulo | Requerimientos | Prioridad Alta | Prioridad Media | Prioridad Baja |
|--------|----------------|----------------|-----------------|----------------|
| Clientes | 7 | 4 | 2 | 1 |
| Artículos | 7 | 4 | 2 | 1 |
| Préstamos | 10 | 5 | 4 | 1 |
| Pagos | 6 | 4 | 1 | 1 |
| Ventas | 4 | 2 | 2 | 0 |
| Documentos | 6 | 3 | 2 | 1 |
| **Caja/POS** | **12** | **6** | **4** | **2** |
| Contabilidad | 7 | 3 | 2 | 2 |
| Administración | 8 | 4 | 3 | 1 |
| Notificaciones | 4 | 0 | 3 | 1 |
| **TOTAL** | **71** | **35** | **25** | **11** |

### Matriz de Funcionalidades por Rol

| Funcionalidad | Super Admin | Admin | Gerente | Cajero | Vendedor |
|---------------|:-----------:|:-----:|:-------:|:------:|:--------:|
| Gestión Usuarios | ✓ | ✓ | - | - | - |
| Configuración Sistema | ✓ | ✓ | - | - | - |
| Ver Todas Sucursales | ✓ | ✓ | - | - | - |
| Crear Préstamos | ✓ | ✓ | ✓ | ✓ | - |
| Aprobar Préstamos | ✓ | ✓ | ✓ | - | - |
| Recibir Pagos | ✓ | ✓ | ✓ | ✓ | - |
| Realizar Ventas | ✓ | ✓ | ✓ | ✓ | ✓ |
| Abrir/Cerrar Caja | ✓ | ✓ | ✓ | ✓ | - |
| Movimientos Caja | ✓ | ✓ | ✓ | ✓ | - |
| Anular Operaciones | ✓ | ✓ | ✓ | - | - |
| Confiscar Artículos | ✓ | ✓ | ✓ | - | - |
| Ver Reportes | ✓ | ✓ | ✓ | - | - |
| Exportar Datos | ✓ | ✓ | ✓ | - | - |

### Checklist de Entregables por Fase

#### Fase 1: Fundamentos
- [ ] Repositorio configurado con CI/CD
- [ ] Base de datos con migraciones
- [ ] API de autenticación (login, logout, refresh)
- [ ] CRUD de usuarios con roles
- [ ] Documentación Swagger/OpenAPI
- [ ] Tests unitarios (>80% cobertura)

#### Fase 2: Core del Negocio
- [ ] CRUD completo: Clientes, Artículos, Categorías
- [ ] Flujo de préstamos completo
- [ ] Cálculo de intereses (simple y sobre saldo)
- [ ] Sistema de pagos con aplicación correcta
- [ ] Ventas de artículos
- [ ] Tests de integración

#### Fase 3: Sistema POS
- [ ] Aplicación de escritorio instalable
- [ ] Apertura y cierre de caja
- [ ] Registro de operaciones
- [ ] Impresión en térmica 80mm
- [ ] Modo offline con sync
- [ ] Cortes X y Z

#### Fase 4: Documentos
- [ ] Generación de contratos PDF
- [ ] Tickets de pago
- [ ] Comprobantes de venta
- [ ] Plantillas personalizables
- [ ] Logo y branding configurables

#### Fase 5: Avanzado
- [ ] Scheduler de tareas automáticas
- [ ] Cálculo diario de intereses
- [ ] Alertas de vencimiento
- [ ] Logs de auditoría
- [ ] Renovaciones de préstamos

#### Fase 6: Panel Admin
- [ ] Dashboard con KPIs
- [ ] Todos los módulos CRUD
- [ ] Reportes con gráficos
- [ ] Exportación Excel/PDF
- [ ] Configuración del sistema

#### Fase 7: Deploy
- [ ] Docker images optimizadas
- [ ] CI/CD completo
- [ ] Monitoreo con alertas
- [ ] Documentación de operaciones
- [ ] Manual de usuario

---

## Apéndice C: Glosario

| Término | Definición |
|---------|------------|
| **Préstamo Prendario** | Préstamo garantizado con un artículo como garantía |
| **Capital** | Monto principal del préstamo sin intereses |
| **Interés** | Porcentaje cobrado sobre el capital o saldo |
| **Pago Mínimo** | Cantidad mínima a pagar mensualmente para mantener el préstamo activo |
| **Confiscación** | Proceso de tomar posesión del artículo cuando el préstamo no se paga |
| **Renovación** | Extensión del plazo del préstamo pagando intereses |
| **Corte X** | Resumen parcial de caja sin cerrar la sesión |
| **Corte Z** | Cierre definitivo de caja al final del día |
| **Arqueo** | Conteo físico del dinero en caja |
| **Tasación** | Evaluación del valor de un artículo |

---

*Documento preparado para: Sistema de Casa de Empeño*
*Versión: 1.1*
*Fecha: Febrero 2026*
*Autor: Equipo de Desarrollo*
