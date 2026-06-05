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

### Eventos (público, sin token)
| Método | Ruta          | Descripción                          |
|--------|---------------|--------------------------------------|
| GET    | /events       | Catálogo (filtro opcional `?categoria=`) |
| GET    | /events/:id   | Detalle de un evento                 |

El endpoint `GET /events` acepta un filtro opcional por categoría:
```
GET /events?categoria=musica
```
Si no se indica categoría, devuelve todos los eventos no cancelados. El filtro se aplica directamente en la consulta SQL, no en memoria.

### Tickets (protegido)
| Método | Ruta                     | Descripción                          |
|--------|--------------------------|--------------------------------------|
| POST   | /tickets                 | Comprar entrada (valida cupo)        |
| GET    | /tickets/mine            | Mis entradas                         |
| DELETE | /tickets/:id             | Cancelar entrada (libera/reasigna cupo) |
| PUT    | /tickets/:id/transfer    | Transferir entrada a otro usuario    |

Al cancelar una entrada, el sistema resta 1 a `EntradasVendidas` en el evento (liberando el cupo) y cambia el estado del ticket a `cancelado`. El ticket no se elimina — queda en la base de datos para mantener el historial.

### Lista de espera (protegido)
| Método | Ruta                     | Descripción                          |
|--------|--------------------------|--------------------------------------|
| POST   | /events/:id/waitlist     | Anotarse en la lista de espera (solo si el evento está agotado) |
| GET    | /waitlist/mine           | Mis anotaciones y asignaciones       |

### Administrador (protegido, requiere rol admin)
| Método | Ruta                     | Descripción                          |
|--------|--------------------------|--------------------------------------|
| POST   | /events                  | Crear nuevo evento                   |
| PUT    | /events/:id              | Actualizar datos de un evento        |
| PATCH  | /events/:id/cancel       | Cancelar un evento                   |
| GET    | /events/:id/report       | Reporte de ocupación y compradores   |

> Estos endpoints están implementados en el backend pero **no tienen vista en el frontend**. Para usarlos se requiere una herramienta como Postman con el token de un usuario admin en el header `Authorization: Bearer <token>`. Las vistas de administrador corresponden al hito 2 (entrega final).

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

## Crear un usuario administrador

La app solo permite registrar usuarios con rol `cliente`. Para crear un admin:

1. Registrá el usuario normalmente desde la app (esto genera el hash de contraseña correctamente).
2. Abrí MySQL Workbench, buscá la tabla `users` y cambiá el campo `rol` de `cliente` a `admin` para ese usuario.

## Decisiones de diseño

### Seguridad de contraseñas: SHA-256 + Salt

Las contraseñas nunca se almacenan en texto plano. Al registrarse, el sistema:

1. Genera un **salt** aleatorio de 16 bytes único por usuario.
2. Concatena la contraseña con el salt y aplica **SHA-256**.
3. Guarda en la base de datos solo el hash resultante y el salt — nunca la contraseña original.

Al hacer login, repite el proceso con la contraseña ingresada y compara el hash con el almacenado.

El salt garantiza que dos usuarios con la misma contraseña tengan hashes distintos en la base de datos, protegiéndose contra ataques de rainbow table.

La contraseña nunca aparece en las respuestas de la API — el DTO `UserResponse` solo expone `id`, `nombre`, `email` y `rol`.

### Borrado de registros: Soft Delete

Las tres entidades principales (`User`, `Event`, `Ticket`) usan **soft delete** en lugar de borrado físico.

Esto significa que al "eliminar" un registro, GORM no ejecuta un `DELETE` en la base de datos. En cambio, escribe la fecha y hora actual en la columna `deleted_at`. Las consultas normales filtran automáticamente los registros con `deleted_at IS NOT NULL`, haciéndolos invisibles para la aplicación, pero el dato sigue existiendo en la tabla.

Este comportamiento viene del `gorm.Model` embebido en cada struct, que agrega cuatro columnas automáticamente:

```
id           → clave primaria autoincremental
created_at   → fecha de creación
updated_at   → fecha de última modificación
deleted_at   → NULL mientras existe; fecha de borrado si fue eliminado (soft delete)
```

**Por qué lo elegimos:**
- **Integridad referencial:** un `Ticket` siempre puede mostrar los datos del `Event` al que pertenece, aunque ese evento se haya cancelado. Si borráramos el evento físicamente, el ticket quedaría huérfano.
- **Auditoría:** se conserva el historial completo de eventos y compras, útil para reportes.
- **Cancelación de eventos:** marcar `Cancelado = true` en el evento comunica el estado de negocio; el soft delete protege el registro histórico.

**Implicancia práctica:** si se necesita consultar registros borrados (p. ej. para un reporte de eventos cancelados), se puede usar `db.Unscoped()` en GORM para ignorar el filtro de `deleted_at`.

## Estructura

```
backend/   → API REST en Go (Gin + GORM)
frontend/  → SPA en React + Vite
docs/      → Diagrama de BD y documentación
```

### Lista de espera (bonus)

Cuando un evento está agotado, un usuario puede anotarse en la lista de espera. Al cancelarse una entrada, si hay alguien esperando, el cupo liberado no vuelve al stock: se le crea automáticamente un ticket activo al primero de la lista (FIFO por fecha de anotación) y su anotación pasa a `asignado`. Solo si la lista está vacía se decrementa `EntradasVendidas`.
