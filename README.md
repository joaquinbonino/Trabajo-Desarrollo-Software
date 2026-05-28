# Sistema de Gestión de Eventos y Entradas

Trabajo integrador — Desarrollo de Software 2026 (UCC, Facultad de Ingeniería).

Backend REST en Go + Frontend en React. Permite explorar eventos, comprar entradas, cancelar y transferirlas.

## Prerrequisitos

- Go >= 1.22
- Node >= 20
- Docker + Docker Compose (para levantar todo junto)
- MySQL 8 (si corrés sin Docker)

## Levantar con Docker (recomendado)

```bash
cp .env.example .env        # completar JWT_SECRET y DB_PASSWORD
docker compose up --build
```

- Frontend: http://localhost
- Backend:  http://localhost:8080
- Health:   http://localhost:8080/health

> La primera vez que levanta, GORM crea las tablas automáticamente vía `AutoMigrate`. No hace falta crear nada en SQL.

## Levantar sin Docker

### Base de datos
Levantá MySQL 8 localmente y creá la base de datos:
```sql
CREATE DATABASE eventos_db;
```

### Backend
```bash
cd backend
cp .env.example .env   # completar con tus datos locales
go mod tidy
go run .
```

### Frontend
```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

## Variables de entorno

| Variable     | Descripción                        | Ejemplo         |
|--------------|------------------------------------|-----------------|
| DB_HOST      | Host de MySQL                      | localhost       |
| DB_PORT      | Puerto MySQL                       | 3306            |
| DB_USER      | Usuario MySQL                      | root            |
| DB_PASSWORD  | Contraseña MySQL                   | secret          |
| DB_NAME      | Nombre de la base de datos         | eventos_db      |
| JWT_SECRET   | Secreto para firmar los tokens JWT | clave_secreta   |
| PORT         | Puerto del backend                 | 8080            |

## Endpoints disponibles

### Auth (público)
| Método | Ruta              | Descripción                        |
|--------|-------------------|------------------------------------|
| POST   | /auth/register    | Registrar usuario nuevo            |
| POST   | /auth/login       | Login, devuelve JWT                |

### Eventos (público)
| Método | Ruta          | Descripción                          |
|--------|---------------|--------------------------------------|
| GET    | /events       | Catálogo (filtro opcional `?categoria=`) |
| GET    | /events/:id   | Detalle de un evento                 |

### Tickets (protegido)
| Método | Ruta                     | Descripción                          |
|--------|--------------------------|--------------------------------------|
| POST   | /tickets                 | Comprar entrada (valida cupo)        |
| GET    | /tickets/mine            | Mis entradas (saca el user del JWT)  |
| DELETE | /tickets/:id             | Cancelar entrada (libera/reasigna cupo) |
| PUT    | /tickets/:id/transfer    | Transferir entrada a otro usuario    |

### Lista de espera (protegido)
| Método | Ruta                     | Descripción                          |
|--------|--------------------------|--------------------------------------|
| POST   | /events/:id/waitlist     | Anotarse en la lista de espera (solo si el evento está agotado) |
| GET    | /waitlist/mine           | Mis anotaciones y asignaciones       |

### Health
| Método | Ruta      | Descripción                        |
|--------|-----------|-------------------------------------|
| GET    | /health   | Estado del servidor y conexión DB  |

### Usar el token JWT
Todos los endpoints protegidos requieren el header:
```
Authorization: Bearer <token>
```

## Comandos útiles (backend)

```bash
go test ./...                    # correr todos los tests
go test ./... -cover             # cobertura por paquete
go vet ./... && gofmt -l .       # lint
```

## Estructura

```
backend/   → API REST en Go (Gin + GORM)
frontend/  → SPA en React + Vite
docs/      → Diagrama de BD y documentación
```

## Decisiones de diseño

- **Lista de espera (bonus):** cuando un evento está agotado, un usuario puede anotarse
  (`WaitlistEntry`, estados `pendiente`/`asignado`/`cancelado`). Al cancelarse una entrada,
  si hay alguien esperando, el cupo liberado **no vuelve al stock**: se le crea
  automáticamente un ticket activo al primero de la lista (FIFO por fecha de anotación) y su
  anotación pasa a `asignado` con `FechaAsignacion`. Esa asignación queda como **registro**
  visible en `GET /waitlist/mine` y la entrada nueva aparece en "Mis Entradas" (no hay envío
  de email porque no hay infraestructura de correo). Solo si la lista está vacía se decrementa
  `EntradasVendidas`.
- **Atomicidad:** la cancelación con reasignación encadena varias operaciones de DAO sin una
  transacción explícita, manteniendo el estilo del resto de los servicios. Queda anotado como
  posible mejora envolver la operación en una transacción GORM.
