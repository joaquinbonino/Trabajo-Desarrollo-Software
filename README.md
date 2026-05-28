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
