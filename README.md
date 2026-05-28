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

| Variable     | Descripción                        |
|--------------|------------------------------------|
| DB_HOST      | Host de MySQL (default: localhost) |
| DB_PORT      | Puerto MySQL (default: 3306)       |
| DB_USER      | Usuario MySQL                      |
| DB_PASSWORD  | Contraseña MySQL                   |
| DB_NAME      | Nombre de la base de datos         |
| JWT_SECRET   | Secreto para firmar los tokens JWT |
| PORT         | Puerto del backend (default: 8080) |

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
