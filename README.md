# Pawnshop - Sistema de Gestión para Casas de Empeño

Sistema integral de gestión para casas de empeño desarrollado en Go (Golang).

## Características

- **Multi-sucursal**: Gestión de múltiples sucursales con configuración independiente
- **Gestión de Clientes**: Registro completo con validación de identidad y límites de crédito
- **Gestión de Artículos**: Inventario con tasación, estados y transferencias entre sucursales
- **Préstamos Prendarios**: Creación, renovación, cálculo de intereses y confiscación
- **Sistema de Pagos**: Pagos parciales y totales con aplicación a intereses/capital
- **Ventas**: Venta de artículos confiscados o propios
- **Sistema POS/Caja**: Apertura/cierre de caja, movimientos, cortes X y Z
- **Generación de Documentos**: Contratos, recibos y comprobantes en PDF
- **Autenticación JWT**: Tokens de acceso y refresh con roles y permisos
- **API REST**: Documentada con OpenAPI/Swagger

## Requisitos

- Go 1.22+
- PostgreSQL 16+
- Redis 7+ (opcional, para caché)
- Docker & Docker Compose (opcional)

## Inicio Rápido

### Con Docker Compose

```bash
# Iniciar todos los servicios
cd docker
docker-compose up -d

# Ver logs
docker-compose logs -f api
```

### Desarrollo Local

```bash
# Instalar dependencias
go mod download

# Configurar variables de entorno (copiar y editar)
cp config.yaml config.local.yaml

# Ejecutar migraciones
make migrate-up

# Ejecutar aplicación
make run

# O con hot reload
make dev
```

## Estructura del Proyecto

```
pawnshop/
├── cmd/
│   └── api/             # Punto de entrada de la API
├── internal/
│   ├── config/          # Configuración
│   ├── domain/          # Entidades de dominio
│   ├── repository/      # Capa de datos
│   ├── service/         # Lógica de negocio
│   ├── handler/         # Handlers HTTP
│   ├── middleware/      # Middlewares
│   └── pdf/             # Generación de PDFs
├── pkg/
│   ├── auth/            # JWT y passwords
│   ├── validator/       # Validación
│   └── response/        # Respuestas JSON
├── migrations/          # Migraciones SQL
├── docker/              # Docker configs
├── api/                 # Documentación OpenAPI
└── web/                 # Frontend Admin Panel (React)
    ├── src/             # Código fuente
    └── e2e/             # Tests E2E con Playwright
```

## Frontend (Admin Panel)

El panel de administración está desarrollado con:
- React 18 + Vite + TypeScript
- Tailwind CSS + shadcn/ui
- TanStack Query + Zustand
- React Hook Form + Zod

### Iniciar el Frontend

```bash
cd web
npm install
npm run dev
```

El frontend estará disponible en `http://localhost:5173`

### Build de Producción

```bash
cd web
npm run build
npm run preview
```

## Tests E2E (Playwright)

El proyecto incluye una suite completa de tests end-to-end con Playwright que cubren todos los flujos del sistema.

### Prerrequisitos

1. Backend corriendo en `http://localhost:8080`
2. Frontend corriendo en `http://localhost:5173`
3. Base de datos con datos de prueba

### Comandos de Testing

```bash
cd web

# Ejecutar todos los tests
npm test

# Ejecutar con interfaz visual (recomendado para desarrollo)
npm run test:ui

# Ejecutar mostrando el navegador
npm run test:headed

# Modo debug (paso a paso)
npm run test:debug

# Ejecutar un archivo específico
npx playwright test e2e/21-complete-pawn-flow.spec.ts

# Ejecutar tests que coincidan con un patrón
npx playwright test -g "customer"

# Ver reporte HTML del último test
npm run test:report
```

### Estructura de Tests

| Archivo | Descripción | Tests |
|---------|-------------|-------|
| `01-navigation.spec.ts` | Navegación básica y sidebar | 15 |
| `02-customers.spec.ts` | CRUD de clientes | 20 |
| `03-items.spec.ts` | Gestión de artículos | 18 |
| `04-loans.spec.ts` | Préstamos y wizard | 22 |
| `05-payments.spec.ts` | Sistema de pagos | 16 |
| `06-sales.spec.ts` | Ventas de artículos | 14 |
| `07-cash.spec.ts` | Apertura/cierre de caja | 19 |
| `08-users.spec.ts` | Gestión de usuarios | 17 |
| `09-settings.spec.ts` | Configuración | 12 |
| `10-dashboard.spec.ts` | KPIs del dashboard | 15 |
| `11-crud-customers.spec.ts` | CRUD completo clientes | 18 |
| `12-crud-items.spec.ts` | CRUD completo artículos | 16 |
| `13-loan-lifecycle.spec.ts` | Ciclo de vida préstamos | 20 |
| `14-form-validation.spec.ts` | Validación de formularios | 25 |
| `15-search-filters.spec.ts` | Búsqueda y filtros | 22 |
| `16-cashier-workflow.spec.ts` | Flujo de cajero | 27 |
| `17-admin-workflow.spec.ts` | Flujo de administrador | 36 |
| `18-manager-workflow.spec.ts` | Flujo de gerente | 34 |
| `19-permissions.spec.ts` | Permisos por rol | 39 |
| `20-edge-cases.spec.ts` | Casos borde y errores | 31 |
| `21-complete-pawn-flow.spec.ts` | **Flujo completo de empeño** | 18 |

### Test Destacado: Flujo Completo de Empeño

El archivo `21-complete-pawn-flow.spec.ts` simula el proceso real de una casa de empeño:

1. **Crear cliente** - Registro de nuevo cliente
2. **Crear artículo** - Registro de artículo a empeñar
3. **Wizard de préstamo** - Proceso de 4 pasos:
   - Selección de cliente
   - Selección de artículo
   - Definición de términos (monto, plazo, interés)
   - Confirmación y creación
4. **Registrar pago** - Pago de cuota
5. **Verificación** - Validar historial del cliente

```bash
# Ejecutar solo el flujo completo
npx playwright test e2e/21-complete-pawn-flow.spec.ts --headed
```

### Autenticación en Tests

Los tests usan autenticación global configurada en `e2e/fixtures.ts`. El login se realiza una vez y se reutiliza en todos los tests mediante `storageState`.

## API Endpoints

### Autenticación
- `POST /api/v1/auth/login` - Iniciar sesión
- `POST /api/v1/auth/refresh` - Refrescar token
- `POST /api/v1/auth/logout` - Cerrar sesión
- `GET /api/v1/auth/me` - Usuario actual
- `POST /api/v1/auth/change-password` - Cambiar contraseña

### Usuarios
- `GET /api/v1/users` - Listar usuarios
- `POST /api/v1/users` - Crear usuario
- `GET /api/v1/users/:id` - Obtener usuario
- `PUT /api/v1/users/:id` - Actualizar usuario
- `DELETE /api/v1/users/:id` - Eliminar usuario

### (Más endpoints en desarrollo)

## Comandos Make

```bash
make build          # Compilar aplicación
make run            # Ejecutar aplicación
make dev            # Ejecutar con hot reload
make test           # Ejecutar tests
make lint           # Ejecutar linter
make migrate-up     # Ejecutar migraciones
make migrate-down   # Revertir migraciones
make docker-build   # Construir imagen Docker
make swagger        # Generar documentación Swagger
make help           # Ver todos los comandos
```

## Configuración

La aplicación se configura mediante:

1. Archivo `config.yaml`
2. Variables de entorno con prefijo `PAWN_`

Ejemplo de variables de entorno:
```bash
PAWN_DATABASE_HOST=localhost
PAWN_DATABASE_PORT=5432
PAWN_DATABASE_USER=postgres
PAWN_DATABASE_PASSWORD=postgres
PAWN_DATABASE_DBNAME=pawnshop
PAWN_JWT_SECRET=your-secret-key
```

## Usuario por Defecto

- **Email**: admin@pawnshop.com
- **Contraseña**: admin123 (cambiar en producción)

## Licencia

Propietario - Todos los derechos reservados
