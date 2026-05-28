# CLAUDE.md

Guía de contexto para el agente. Lee esto **antes de tocar código**. Si una instrucción
del usuario contradice este archivo, avisa antes de continuar.

---

## 1. Qué estamos construyendo

**Sistema de Gestión de Eventos y Entradas (tipo Ticketek)** — trabajo integrador de la
materia *Desarrollo de Software 2026* (UCC, Facultad de Ingeniería).

Dos roles, una sola API:
- **Cliente**: explora el catálogo de eventos, ve el detalle, compra entradas, consulta
  "Mis Entradas", cancela una compra y transfiere una entrada a otro usuario.
- **Administrador**: crea, edita y cancela eventos, y consulta reportes de ocupación/ventas.

El backend (Go) expone una API REST. El frontend (React) la consume. Están **desacoplados**.

> Es un trabajo evaluado. La nota depende de respetar la arquitectura, los tests y la
> documentación tanto como de que "funcione". No tomes atajos que rompan las convenciones
> de abajo.

---

## 2. Stack y decisiones fijas

| Capa        | Tecnología                                                       |
|-------------|------------------------------------------------------------------|
| Backend     | Go (>= 1.26), router **Gin** (`github.com/gin-gonic/gin`)        |
| ORM         | **GORM** (`gorm.io/gorm` + `gorm.io/driver/mysql`) — obligatorio |
| Auth        | JWT (`github.com/golang-jwt/jwt/v5`)                             |
| Base datos  | **MySQL 8**                                                      |
| Frontend    | React + **Vite**, React Router, Axios                            |
| Tests BE    | `testing` + `testify`, `net/http/httptest` para controllers     |
| Contenedores| Docker + Docker Compose (un servicio por componente)            |

**Decisiones que NO se negocian** (vienen del enunciado):
- Las tablas se crean y mapean **estrictamente con GORM** (`AutoMigrate` + structs con tags).
  Prohibido crear tablas con SQL crudo.
- Las contraseñas **nunca** se guardan en texto plano: hashing con **SHA-256 + salt**
  (el enunciado admite MD5/SHA256; usamos SHA-256). El hash vive en `utils`.
- El JWT debe distinguir **dos roles** (`cliente`, `admin`) y **expirar** (usar `exp`).
- Las operaciones sensibles (escritura, compra, etc.) van **detrás del middleware de auth**.

Si vas a cambiar alguna de estas decisiones, **pregunta primero**.

---

## 3. Estructura del repositorio

Monorepo con frontend y backend separados:

```
/
├── backend/
│   ├── clients/        # conexión a MySQL y cualquier servicio externo
│   ├── controllers/    # handlers HTTP (Gin). Punto de entrada de cada feature.
│   ├── dao/            # Data Access Object: queries GORM, mapeo modelo<->tabla
│   ├── domain/         # entidades del negocio (structs GORM) + DTOs + reglas
│   ├── services/       # lógica de negocio. Llamada desde controllers, usa dao.
│   ├── utils/          # helpers: JWT, hashing, respuestas HTTP, errores
│   ├── main.go         # arranque: server HTTP, routers, inyección de servicios
│   ├── go.mod
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── api/        # cliente axios + funciones por endpoint
│   │   ├── components/ # componentes reutilizables
│   │   ├── pages/      # vistas (Catálogo, Detalle, MisEntradas, Admin, ...)
│   │   ├── context/    # auth context / estado de sesión
│   │   └── App.jsx
│   ├── package.json
│   └── Dockerfile
├── docs/               # diagrama de BD (fuente + imagen) y otros recursos
├── docker-compose.yml
└── README.md
```

### Regla de oro de la arquitectura (backend)

El flujo de una request es **siempre**:

```
controller  ->  service  ->  dao  ->  (GORM / MySQL)
```

- El **controller** no contiene lógica de negocio: parsea/valida input, llama al service,
  arma la respuesta HTTP y el status code.
- El **service** tiene la lógica y las reglas de negocio (validaciones de cupo, propiedad
  del ticket, etc.). No conoce a Gin ni a `http`.
- El **dao** es lo único que habla con GORM. Ningún otro paquete importa la conexión.
- El **domain** define las entidades y los DTOs. No depende de las otras capas.

Nunca saltees capas (p. ej. un controller que llama directo al dao).

---

## 4. Modelo de dominio (mínimo)

Tres entidades obligatorias. Agregar intermedias solo si el diseño lo necesita.

- **User**: `ID`, `Nombre`, `Email` (único), `PasswordHash`, `Rol` (`cliente`|`admin`),
  timestamps.
- **Event**: `ID`, `Titulo`, `Descripcion`, `Categoria`, `FechaHora`, `Duracion`,
  `CapacidadTotal`, `EntradasVendidas` (o calculado), `Foto`, timestamps.
- **Ticket**: `ID`, `EventID` (FK), `UserID` (FK, titular actual), `Estado`
  (`activo`|`cancelado`|`transferido`), `FechaCompra`, timestamps.

Relaciones con **claves foráneas reales** configuradas en los structs (tags GORM:
`foreignKey`, `references`). Decidir y documentar el tipo de borrado de eventos
(soft delete vs. cancelación lógica) — es una de las "decisiones de diseño" que pide el README.

Reglas de negocio que el agente debe respetar:
- No vender entradas si `EntradasVendidas >= CapacidadTotal`.
- Cancelar una compra **libera el cupo** del evento.
- Transferir un ticket cambia el `UserID` titular de forma íntegra (transacción atómica).
- "Mis Entradas" devuelve **solo** los tickets del usuario autenticado (sacar el ID del JWT,
  nunca de un parámetro que mande el cliente).

---

## 5. Comandos

### Backend
```bash
cd backend
go mod tidy
go run .                          # levanta el server
go test ./...                     # corre todos los tests
go test ./... -cover              # cobertura por paquete
go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out  # total
go tool cover -html=coverage.out  # ver cobertura en el navegador
go vet ./... && gofmt -l .        # lint básico antes de commitear
```

### Frontend
```bash
cd frontend
npm install
npm run dev      # desarrollo
npm run build    # build de producción
npm run lint
```

### Todo junto (Docker)
```bash
docker compose up --build   # levanta frontend + backend + MySQL
```

---

## 6. Convenciones de código

**Go**
- `gofmt` siempre. Nombres exportados en `PascalCase`, internos en `camelCase`.
- Errores: devolver `error` hacia arriba; el controller los traduce a status codes.
  Centralizar respuestas de error en `utils` (no repetir `c.JSON(...)` a mano por todos lados).
- DTOs de request/response separados de las entidades de dominio (no exponer el `PasswordHash`).
- Inyección de dependencias por constructor (`NewXService(dao)`), no globales. Facilita testear.

**React**
- Componentes funcionales + hooks. Nada de `localStorage` para estado de UI; sí se puede usar
  para el token JWT (es lo habitual), centralizado en el auth context.
- Llamadas a la API solo desde `src/api/`, nunca `fetch` suelto dentro de un componente.
- Validar formularios del lado cliente **antes** de pegarle al backend.

---

## 7. Seguridad

- Hash de contraseña en `utils` (SHA-256 + salt). El password plano no sale nunca de la capa
  de auth ni se loguea.
- JWT firmado con secreto desde variable de entorno (`JWT_SECRET`). Claims: `user_id`, `rol`,
  `exp`. Middleware de Gin que valida el token y, para endpoints de admin, valida el rol.
- Endpoints públicos: catálogo y detalle de evento. Todo lo demás, protegido.
- No commitear secretos. Usar `.env` (gitignored) y un `.env.example` versionado.

---

## 8. Testing (se evalúa la cobertura)

- Unitarios sobre **services** (lógica) y de integración sobre **controllers** con `httptest`.
- Mockear el dao en los tests de service (interfaz + mock con testify) para no depender de MySQL.
- Probar casos de **éxito y de error**: comprar sin cupo, acceder sin token, token vencido,
  cancelar ticket ajeno, etc. Verificar el status code esperado en cada uno.
- Objetivo de cobertura: **40% para regularidad**, **80% para el final** (services + controllers).
  Al cerrar una feature, corré `-cover` y reportá el número.

---

## 9. Git / GitFlow

- Ramas: `main` (estable), `develop` (integración), `feature/<nombre>` por funcionalidad.
- Trabajar siempre en una `feature/*` que sale de `develop` y vuelve a `develop` por PR.
- Commits chicos y descriptivos. El trabajo es grupal y **se evalúa que todos contribuyan**:
  no juntes el trabajo de varias personas en un solo commit.
- Mensajes en imperativo y claros: `feat: endpoint de compra de entrada`, `test: cobertura service eventos`.

---

## 10. Alcance: qué priorizar

El trabajo tiene dos hitos. **No mezclar el orden.**

### Entrega para REGULARIDAD (12–19 jun) — construir primero
- Estructura base BE + FE y conexión a MySQL.
- Registro y login con hash de contraseña.
- JWT básico (solo **autenticación**).
- **Todo** el flujo del rol **Cliente** (endpoints + vistas): catálogo, detalle, compra,
  Mis Entradas, cancelación, transferencia.
- README con pasos para levantar el proyecto.
- Cobertura mínima **40%**.
- *Fuera de este hito:* roles/permisos, todo lo del Admin, Docker, Bonus.

### Entrega para FINAL — construir después
- JWT con **autorización** (validación estricta de roles en endpoints protegidos).
- Rol **Admin** completo (endpoints + vistas: ABM de eventos + reportes).
- **Dockerización** completa (compose con los 3 servicios).
- **Bonus Track** (funcionalidad extra acordada con el docente).
- Cobertura **80%**.
- Documentación definitiva (diagrama de BD en `/docs` e incrustado en el README,
  sección de decisiones de diseño, capturas).

Cuando arranques una tarea, ubicala en uno de los dos hitos y, si no está claro, preguntá.

---

## 11. Reglas para el agente

- Antes de implementar una feature, decí en qué capa va cada pieza y respetá el flujo
  controller → service → dao.
- No introduzcas dependencias nuevas sin avisar.
- No borres ni reescribas tests existentes para "que pasen"; arreglá el código.
- Si un requisito del enunciado es ambiguo (atributos de una entidad, criterio de filtro,
  método de borrado), proponé una opción y pedí confirmación en vez de asumir en silencio.
- Mantené el README y el `.env.example` actualizados cuando agregues config.