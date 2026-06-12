# Testing — Resultados

Estado de las pruebas del backend del Sistema de Gestión de Eventos y Entradas.

Comando de medición:

```bash
cd backend
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out        # detalle por función + total
```

---

## Resultados de cobertura

| Paquete | Cobertura | Nota |
|---------|-----------|------|
| `backend/services` | **98.9%** | 100% del código alcanzable |
| `backend/controllers` | **99.3%** | 100% del código alcanzable |
| `backend/utils` | **98.0%** | 100% del código alcanzable |
| **services + controllers (foco del enunciado)** | **99.1%** | Supera ampliamente el 80% exigido |
| `backend/dao` | 0.0% | Ver "Qué faltó y por qué" |
| `backend/clients` | 0.0% | Ver "Qué faltó y por qué" |
| `backend/main.go` | 0.0% | Ver "Qué faltó y por qué" |
| `backend/domain` | — | Solo structs/DTOs, sin lógica que testear |
| **TOTAL backend** | **76.5%** | El total lo bajan dao/clients/main |
| **Frontend** | 0% | Sin tests (bonus opcional del enunciado) |

> El enunciado (§5) pide 80% *"principalmente en las capas de servicios (lógica) y controladores"*.
> Ese subconjunto está en **99.1%**. El total global (76.5%) es menor porque incluye paquetes que no
> se pueden testear con las herramientas autorizadas (ver abajo).

**Todos los tests pasan** (`go test ./...` en verde) y **`go vet ./...` está limpio.**

---

## Qué se testeó

**Pruebas unitarias (services + utils), con el DAO mockeado (testify):**
- `user_service`: registro, login, email duplicado, credenciales inválidas, errores del DAO.
- `event_service`: listado, detalle, creación, actualización, cancelación, listado admin, reporte de ocupación, y todas sus ramas de error.
- `ticket_service`: compra (con/sin cupo, evento cancelado/inexistente), "mis entradas", cancelación (con liberación de cupo y con reasignación a lista de espera), transferencia, y ramas de error.
- `waitlist_service`: anotarse (evento agotado/cancelado/con cupo/ya anotado), consulta, errores del DAO.
- `utils/hash`: hashing determinístico, salt distinto, no expone la contraseña.
- `utils/jwt`: round-trip de claims, token malformado/expirado/secreto incorrecto/método de firma inesperado.
- `utils/response`: formato uniforme `{"data"}` y `{"error"}`.

**Pruebas de integración (controllers con `httptest`), validando status codes de éxito y error:**
- `auth`: register y login (éxito, body inválido, email duplicado, credenciales inválidas).
- `event`: catálogo y detalle (público) + ABM admin completo (crear, actualizar, listar, reporte, cancelar) con sus errores.
- `ticket`: compra, "mis entradas", cancelación y transferencia (éxito, id inválido, body inválido, errores de negocio, error de servicio).
- `waitlist`: anotarse y consultar.
- `utils/middleware`: `AuthMiddleware` (sin token / sin Bearer / token inválido / token válido), `AdminMiddleware` (rol no-admin / rol admin) y `CORSMiddleware` (preflight OPTIONS / request normal).
- `health`: rama de DB no disponible (503).

---

## Qué faltó y por qué

El enunciado (§5) autoriza para el backend **únicamente**: librería estándar `testing`, `testify`
y `httptest`. Con ese set, lo que quedó sin cubrir es lo que **no se puede testear sin sumar
herramientas que el enunciado no lista**:

### Paquetes sin cobertura
- **`dao`** — habla directamente con GORM/MySQL. Testearlo requiere una base de datos: o bien un
  driver SQLite en memoria (`glebarez/sqlite`) o un SQL-mock (`go-sqlmock`). Ninguno está autorizado
  por la consigna, así que no se usó.
- **`clients`** — solo arma el DSN y abre la conexión a MySQL; mismo motivo que `dao`.
- **`main.go`** — arranque del servidor (wiring de routers e inyección de dependencias). No tiene
  lógica de negocio aislable; se ejercita al levantar la app.

### Líneas defensivas/infra no cubiertas (4 en total)
Dentro de services/controllers/utils, lo único sin cubrir son ramas defensivas o de infraestructura:

| Ubicación | Línea | Por qué no se cubre |
|-----------|-------|---------------------|
| `controllers/health_controller.go` | rama `200` del health check | Necesita una DB viva o un SQL-mock (no autorizado). La rama de error (503) **sí** se cubre. |
| `services/user_service.go` (Register) | `return nil, err` tras firmar el JWT | Defensiva: firmar con HMAC + clave `[]byte` nunca falla. |
| `services/user_service.go` (Login) | ídem | Igual que arriba. |
| `utils/jwt.go` | `return nil, errors.New("token inválido")` | Defensiva: si `ParseWithClaims` no da error, el token siempre es válido. Inalcanzable por diseño. |

### Bug encontrado y corregido al testear
El test de `/health` reveló que `db.Raw("SELECT 1").Error` **nunca ejecutaba** la query (en GORM
`Raw` solo arma la sentencia; se ejecuta con `Scan`/`Exec`/`Find`), por lo que `/health` devolvía
`200` aunque la DB estuviera caída. Se corrigió a `db.Exec("SELECT 1").Error`, que sí ejecuta.

---

## Herramientas utilizadas y para qué

Todas son las que el enunciado (§5) permite explícitamente. **No se instaló ninguna dependencia nueva.**

| Herramienta | Para qué se usó |
|-------------|-----------------|
| **`testing`** (librería estándar de Go) | Estructura de todos los tests (`func TestXxx(t *testing.T)`), subtests y asserts base. |
| **`testify/assert`** | Aserciones legibles (`assert.Equal`, `assert.Error`, `assert.EqualError`, etc.) en todos los tests. |
| **`testify/mock`** | Mockear los DAOs en los tests unitarios de services y los services en los tests de controllers, para no depender de MySQL. |
| **`net/http/httptest`** | Pruebas de integración de los controllers: simular peticiones HTTP contra el router de Gin y verificar los status codes de éxito y error. |
| **`gin` (modo test)** | Montar routers de prueba aislados (`gin.New()` + `gin.SetMode(TestMode)`) para los handlers y middlewares. |

> Para cubrir `dao` y la rama `200` de `/health` haría falta un driver SQLite de test
> (`glebarez/sqlite`) o un SQL-mock (`go-sqlmock`). **El enunciado no los menciona ni los permite**,
> por lo que no se instalaron. Si se quisieran agregar, debe validarse con el docente.
