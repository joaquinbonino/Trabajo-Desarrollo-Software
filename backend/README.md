# Backend — Sistema de Gestión de Eventos y Entradas

API REST escrita en Go con Gin. Se conecta a una base de datos MySQL 8 mediante GORM y expone los endpoints que consume el frontend React.

---

## Estructura de carpetas

```
backend/
├── clients/
│   └── db.go
├── controllers/
│   ├── auth_controller.go
│   ├── auth_controller_test.go
│   └── health_controller.go
├── dao/
│   ├── interfaces.go
│   ├── user_dao.go
│   ├── event_dao.go
│   └── ticket_dao.go
├── domain/
│   ├── models.go
│   └── dto.go
├── services/
│   ├── interfaces.go
│   ├── user_service.go
│   ├── user_service_test.go
│   ├── event_service.go
│   └── ticket_service.go
├── utils/
│   ├── hash.go
│   ├── jwt.go
│   ├── middleware.go
│   └── response.go
├── main.go
├── go.mod
├── go.sum
├── .env.example
└── Dockerfile
```

---

## Flujo de una request

Toda request sigue este camino, sin excepciones:

```
controller  →  service  →  dao  →  GORM / MySQL
```

Cada capa tiene una responsabilidad única y no puede saltear a la siguiente.

---

## Descripción de cada carpeta y archivo

### `main.go`

Punto de entrada del servidor. Se encarga de:

1. Cargar las variables de entorno desde `.env` (via `godotenv`).
2. Abrir la conexión a MySQL y ejecutar `AutoMigrate`.
3. Instanciar todos los DAOs, services y controllers en orden (inyección de dependencias por constructor).
4. Registrar las rutas en el router de Gin.
5. Iniciar el servidor HTTP en el puerto configurado (por defecto `8080`).

---

### `clients/`

Contiene los clientes para servicios externos que el backend necesita.

| Archivo | Qué hace |
|---------|----------|
| `db.go` | Arma el DSN de MySQL a partir de las variables de entorno y abre la conexión con GORM. También expone `AutoMigrate`, que sincroniza los structs `User`, `Event` y `Ticket` con las tablas de la base de datos. |

---

### `domain/`

Define las entidades del negocio y los DTOs de la API. **No depende de ninguna otra capa del proyecto.**

| Archivo | Qué contiene |
|---------|--------------|
| `models.go` | Structs mapeados a tablas por GORM: `User`, `Event` y `Ticket`. Cada uno embebe `gorm.Model` (que agrega `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`). Las relaciones entre `Ticket` → `Event` y `Ticket` → `User` se declaran aquí con las tags `foreignKey` / `references`. |
| `dto.go` | DTOs de request y response separados de los modelos. Evita exponer campos internos (como `PasswordHash`) en las respuestas JSON. Incluye los DTOs de auth (`RegisterRequest`, `LoginRequest`, `AuthResponse`), de eventos (`CreateEventRequest`, `EventResponse`) y de tickets (`BuyTicketRequest`, `TransferTicketRequest`, `TicketResponse`). |

**Entidades principales:**

- **`User`** — representa a un usuario del sistema. Campos clave: `Email` (único), `PasswordHash`, `PasswordSalt`, `Rol` (`cliente` | `admin`).
- **`Event`** — representa un evento. Campos clave: `CapacidadTotal`, `EntradasVendidas`, `Cancelado` (borrado lógico).
- **`Ticket`** — representa una entrada comprada. Tiene FK a `Event` y a `User`. El campo `Estado` puede ser `activo`, `cancelado` o `transferido`.

---

### `dao/`

Data Access Objects: la única capa que habla con GORM. Ningún otro paquete importa la conexión a la base de datos directamente.

| Archivo | Qué hace |
|---------|----------|
| `interfaces.go` | Declara las interfaces `IUserDAO`, `IEventDAO` e `ITicketDAO`. Permiten mockear el dao en los tests de service sin tocar MySQL. |
| `user_dao.go` | Implementación de `IUserDAO`: `Create`, `FindByID`, `FindByEmail`, `Update`. |
| `event_dao.go` | Implementación de `IEventDAO`: `Create`, `FindByID`, `FindAll` (con filtro opcional por categoría), `Update`. |
| `ticket_dao.go` | Implementación de `ITicketDAO`: `Create`, `FindByID`, `FindByUserID`, `Update`. |

---

### `services/`

Contiene la lógica y las reglas de negocio. Recibe llamadas desde los controllers y usa el dao para persistir. **No conoce a Gin ni al protocolo HTTP.**

| Archivo | Qué hace |
|---------|----------|
| `interfaces.go` | Declara `IUserService`, `IEventService` e `ITicketService`. Facilitan el mockeo en los tests de controller. |
| `user_service.go` | Lógica de registro y login: hashea la contraseña con salt antes de guardar, genera el JWT al autenticar. |
| `user_service_test.go` | Tests unitarios de `UserService` con dao mockeado. |
| `event_service.go` | Lógica de listado, detalle y creación/cancelación de eventos. |
| `ticket_service.go` | Lógica de compra (verifica cupo disponible), cancelación (libera el cupo), transferencia (cambia el titular en una transacción atómica) y consulta de "Mis Entradas". |

---

### `controllers/`

Handlers HTTP de Gin. Parsean y validan el input, llaman al service correspondiente y arman la respuesta HTTP con el status code correcto. **No contienen lógica de negocio.**

| Archivo | Qué hace |
|---------|----------|
| `health_controller.go` | Endpoint `GET /health`: ejecuta `SELECT 1` para verificar que la base de datos responde. Devuelve `200` si todo está bien o `503` si la DB no está disponible. |
| `auth_controller.go` | Endpoints `POST /auth/register` y `POST /auth/login`. Delega en `IUserService` y devuelve el JWT junto con los datos del usuario. |
| `auth_controller_test.go` | Tests de integración de los endpoints de auth usando `httptest` y un service mockeado. |

---

### `utils/`

Helpers reutilizables que no pertenecen a ninguna capa de negocio en particular.

| Archivo | Qué expone |
|---------|------------|
| `hash.go` | `HashPassword(password, salt)` — SHA-256 del concatenado `password+salt`. `GenerateSalt()` — genera 16 bytes aleatorios en hex. Las contraseñas nunca se guardan en texto plano. |
| `jwt.go` | `GenerateToken(userID, rol)` — firma un JWT con claims `user_id`, `rol` y `exp` (24 h) usando `JWT_SECRET`. `ValidateToken(tokenStr)` — parsea y valida el token, devuelve los claims. |
| `middleware.go` | `AuthMiddleware()` — middleware de Gin que extrae el Bearer token del header `Authorization`, lo valida con `ValidateToken` y propaga `user_id` y `rol` al contexto de la request. Aborta con `401` si el token es inválido o está ausente. |
| `response.go` | `Success(c, code, data)` y `Error(c, code, message)` — envuelven `c.JSON` con un formato de respuesta uniforme: `{"data": ...}` para éxito y `{"error": "..."}` para errores. |

---

### Archivos raíz

| Archivo | Descripción |
|---------|-------------|
| `go.mod` | Módulo Go (`backend`). Declara las dependencias directas: Gin, GORM + driver MySQL, golang-jwt, godotenv, testify. |
| `go.sum` | Hashes de integridad de todas las dependencias (directas e indirectas). No editar a mano. |
| `.env.example` | Plantilla de variables de entorno. Copiar a `.env` (que está en `.gitignore`) y completar con los valores reales antes de levantar el servidor. |
| `Dockerfile` | Build multi-stage: la primera etapa compila el binario con `golang:1.26-alpine`; la segunda copia solo el binario a `alpine:3.19` para una imagen final liviana. Expone el puerto `8080`. |

---

## Variables de entorno

Definidas en `.env` (ver `.env.example`):

| Variable | Descripción |
|----------|-------------|
| `DB_HOST` | Host de MySQL (e.g. `localhost` o nombre del servicio Docker) |
| `DB_PORT` | Puerto de MySQL (por defecto `3306`) |
| `DB_USER` | Usuario de la base de datos |
| `DB_PASSWORD` | Contraseña de la base de datos |
| `DB_NAME` | Nombre de la base de datos (e.g. `eventos_db`) |
| `JWT_SECRET` | Clave secreta para firmar los JWT. Usar una cadena larga y aleatoria en producción. |
| `PORT` | Puerto donde escucha el servidor (por defecto `8080`) |

---

## Comandos útiles

```bash
# Instalar / actualizar dependencias
go mod tidy

# Levantar el servidor
go run .

# Correr todos los tests
go test ./...

# Tests con cobertura por paquete
go test ./... -cover

# Reporte de cobertura total
go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out

# Ver cobertura en el navegador
go tool cover -html=coverage.out

# Lint básico
go vet ./... && gofmt -l .
```
