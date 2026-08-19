# Ecommerce Platform

**Proyecto personal** en evolución: la idea a medio plazo es convertirlo en una **plataforma de comercio electrónico** completa. Hoy el foco sigue apoyándose en un catálogo pensado al inicio para relojes: API en Go, frontend estático (HTML/CSS), usuarios, productos y comentarios.

**Sobre los nombres y el alcance:** el proyecto **empezó** con la idea de centrarse **solo en relojes**; por eso muchas cosas aún hablan de relojes (por ejemplo el nombre de la base `store_watches`, textos del frontend, ejemplos en configuración). **Decidí ampliar** el objetivo hacia una **plataforma de comercio** más general. Ese cambio de rumbo es reciente: **con el tiempo** se irá **actualizando** código, copys, esquema de datos y documentación para que reflejen mejor un ecommerce genérico, sin prisa de romper lo que ya funciona.

## Contenido

- [Qué es este proyecto](#qué-es-este-proyecto)
- [Stack y herramientas](#stack-y-herramientas)
- [Arquitectura del código](#arquitectura-del-código)
- [Requisitos previos](#requisitos-previos)
- [Puesta en marcha con Docker](#puesta-en-marcha-con-docker)
- [Ficheros de Docker Compose](#ficheros-de-docker-compose)
- [Variables de entorno](#variables-de-entorno)
- [Puertos y servicios](#puertos-y-servicios)
- [Base de datos](#base-de-datos)
- [API HTTP](#api-http)
- [Desarrollo local sin Docker](#desarrollo-local-sin-docker)
- [Tests](#tests)
- [CI en GitHub Actions](#ci-en-github-actions)
- [Seguridad y secretos](#seguridad-y-secretos)
- [Solución de problemas](#solución-de-problemas)

## Qué es este proyecto

Es una aplicación web monolito **API + estáticos**: el backend en **Go** expone rutas REST y entrega los archivos del directorio `frontend` (HTML y CSS). Los datos persistentes viven en **MySQL**; **Redis** se usa para soporte de la aplicación (por ejemplo tokens CSRF asociados a la sesión del usuario). El esquema inicial y datos de ejemplo se cargan automáticamente la primera vez que arranca el contenedor de MySQL.

No es un producto comercial cerrado: es un **laboratorio** para ir sumando piezas (pagos, carrito, panel de administración, etc.) hasta tener una plataforma de ecommerce coherente. Lo que aún suene muy a “tienda de relojes” en nombres o textos es herencia de ese arranque; se irá unificando con el tiempo (como se indica arriba).

Funcionalidades principales hoy:

- Página principal y archivos estáticos.
- Registro y login (JWT en cookies HTTP-only; rutas sensibles protegidas).
- Listado de productos, filtro por marca y detalle por ID.
- Listado de comentarios y alta de comentarios (ruta protegida con autenticación y CSRF).
- Rate limiting por IP, CORS diferenciado para rutas públicas y privadas, cabeceras de seguridad HTTP.

En el esquema SQL existen tablas relacionadas con **pagos y webhooks** (PayPal) y el archivo `.env.example` documenta variables para PayPal; la integración activa de pagos puede ampliarse en el backend según evolucione el proyecto.

## Stack y herramientas

| Área | Tecnología |
|------|------------|
| Lenguaje runtime | Go 1.23 |
| HTTP | `net/http`, enrutador [gorilla/mux](https://github.com/gorilla/mux) |
| Configuración | [Viper](https://github.com/spf13/viper) y **variables de entorno** (fichero `.env`) |
| MySQL | Driver [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql), acceso con [sqlx](https://github.com/jmoiron/sqlx) |
| Redis | [go-redis v9](https://github.com/redis/go-redis) |
| Autenticación | JWT ([golang-jwt/jwt](https://github.com/golang-jwt/jwt)), hashing con [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) |
| Contenedores | Docker, Docker Compose |
| Base de datos (contenedor) | MySQL 8.0 |
| Caché / sesión (contenedor) | Redis 7 (Alpine) |
| Frontend | HTML y CSS estático (sin bundler en el repositorio) |

## Arquitectura del código

El backend sigue una organización por capas inspirada en **hexagonal / puertos y adaptadores**:

- `cmd/api`: punto de entrada del binario.
- `internal/core`: dominio, servicios y puertos (interfaces).
- `internal/adapters/primary/http`: handlers HTTP y middleware.
- `internal/adapters/secondary`: repositorios MySQL, Redis, archivos estáticos.
- `internal/bootstrap`: cableado de infraestructura y middlewares.
- `internal/config`: carga y validación de configuración.
- `pkg`: utilidades compartidas (cookies, rate limiter, errores, etc.).

## Requisitos previos

- [Docker](https://docs.docker.com/get-docker/) y [Docker Compose](https://docs.docker.com/compose/) (plugin V2).
- Para desarrollo fuera de contenedor: Go 1.23+, MySQL 8 y Redis 7 accesibles según tu configuración.

## Puesta en marcha con Docker

1. **Clonar el repositorio** (si aún no lo tienes).

2. **Crear el fichero de entorno** en la raíz del proyecto a partir del ejemplo:

   ```bash
   cp .env.example .env
   ```

   Edita `.env` y define al menos:

   - `MYSQL_ROOT_PASSWORD`: contraseña del usuario `root` de MySQL.
   - `REDIS_PASSWORD`: contraseña que usará Redis con `requirepass`.
   - `SECURITY_JWT_JWT_SECRET`: cadena larga y aleatoria para firmar JWT.

3. **Levantar los servicios** desde la raíz del repositorio:

   ```bash
   docker compose up --build
   ```

   La primera vez, MySQL ejecutará el script `backend/migrations/init.sql` (montado en `/docker-entrypoint-initdb.d`) y creará la base `store_watches` con tablas y datos de ejemplo.

4. **Abrir la aplicación** en el navegador:

   - Aplicación: [http://localhost:8080](http://localhost:8080)

Para parar los contenedores: `Ctrl+C` o `docker compose down`. Los datos de MySQL persisten en el volumen Docker `mysql_data`. Si necesitas **reinicializar la base desde cero**, elimina el volumen: `docker compose down -v` (esto borra los datos de MySQL).

### Ficheros de Docker Compose

Además de `docker-compose.yml` en la raíz del repositorio hay dos ficheros opcionales que conviene conocer:

| Fichero | Rol |
|---------|-----|
| `docker-compose.yml` | Definición principal: MySQL, Redis, backend, red `app-network`, volúmenes, límites opcionales de CPU/memoria y **healthcheck del backend** (HTTP vía `wget`). |
| `docker-compose.override.yml` | Si existe en la raíz, **Docker Compose lo fusiona automáticamente** con el fichero principal (no hace falta pasar `-f`). En este proyecto monta `./frontend` en el contenedor del backend para **editar estáticos sin reconstruir la imagen**. Si prefieres servir solo lo copiado en la imagen, renombra o elimina este fichero antes de levantar el stack. |
| `docker-compose.prod.yml` | Fragmento con ajustes orientados a un despliegue más parecido a producción (reinicio, recursos, `ENV=production`, healthcheck del backend). **No** se aplica solo: indica ambos ficheros al ejecutar Compose, por ejemplo: |

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up --build -d
```

Con `docker-compose.prod.yml`, MySQL y Redis suelen publicarse en el host en los puertos **3306** y **6379** (no en 3307/6380). Comprueba que no choquen con instancias locales.

En el **compose base** y en **`docker-compose.prod.yml`**, el healthcheck del backend usa `wget`. La imagen final (`backend/Dockerfile`, etapa Alpine) **no** instala `wget` por defecto: el contenedor puede quedar *unhealthy* hasta que añadas la herramienta al `Dockerfile` o cambies el `test` del healthcheck (por ejemplo con un binario ya presente en la imagen).

## Variables de entorno

La configuración se basa en **variables de entorno**, documentadas en [`.env.example`](.env.example). Docker Compose inyecta en el servicio `backend` las credenciales de MySQL y Redis y el secreto JWT; el resto (rate limiting, PayPal, etc.) puedes definirla en tu `.env` siguiendo ese mismo ejemplo.

En código, las claves que lee Viper coinciden con el nombre de las variables de entorno en mayúsculas (por ejemplo `DATABASE_HOST`, `REDIS_PORT`). `AutomaticEnv()` permite sobrescribir lo definido en `backend/internal/config/.env` con variables del sistema operativo (incluidas las que inyecta Docker Compose en el contenedor).

## Puertos y servicios

| Servicio | Puerto en tu máquina | Puerto dentro de la red Docker |
|----------|----------------------|--------------------------------|
| Backend HTTP | **8080** | 8080 |
| MySQL | **3307** | 3306 |
| Redis | **6380** | 6379 |

Los puertos **3307** y **6380** evitan choques con instalaciones locales de MySQL y Redis en los puertos por defecto.

Si levantas el stack con `docker-compose.prod.yml`, en el host pueden publicarse **3306** y **6379**; los detalles están en [Ficheros de Docker Compose](#ficheros-de-docker-compose).

## Base de datos

- **Motor**: MySQL 8.0, base `store_watches`, charset `utf8mb4`.
- **Inicialización**: `backend/migrations/init.sql` (creación de tablas + datos semilla de relojes).

Tablas principales:

| Tabla | Descripción breve |
|-------|-------------------|
| `user_registration` | Usuarios (nombre de usuario y contraseña almacenada de forma segura vía la aplicación). |
| `Products` | Catálogo (nombre, descripción, precio, stock, marca, imagen, etc.). |
| `comments` | Comentarios y valoraciones ligados a usuarios. |
| `payments` | Registros de órdenes PayPal (esquema preparado para integración). |
| `webhook_events` | Eventos de webhook asociados a pagos. |

**Redis** almacena información efímera necesaria para el flujo seguro de la API (p. ej. tokens CSRF), no sustituye a MySQL como fuente de verdad del negocio.

## API HTTP

Todas las rutas bajo el mismo origen que sirve el backend (por defecto `http://localhost:8080`). Resumen:

| Método | Ruta | Autenticación | Descripción |
|--------|------|---------------|-------------|
| GET | `/` | No | Página principal. |
| GET | `/comments` | No | Lista de comentarios. |
| GET | `/products` | No | Lista de productos. |
| GET | `/product-id/{id}` | No | Detalle de un producto. |
| GET | `/products-brand/{brand}` | No | Productos por marca. |
| POST | `/register` | No | Registro de usuario (emisión de cookies/JWT según implementación). |
| POST | `/login` | No | Inicio de sesión. |
| POST | `/comments/newComments` | Sí (JWT) + CSRF | Alta de comentario. |

Las rutas `POST` que reciben JSON validan `Content-Type: application/json`
(también con parámetros como `charset=utf-8`), rechazan campos desconocidos,
aceptan un solo valor JSON por request y limitan el tamaño del cuerpo para
evitar payloads excesivos: autenticación usa 8 KiB y comentarios usa 16 KiB.

Los ficheros estáticos (CSS, imágenes, etc.) se sirven bajo rutas registradas por `StaticFileHandler` (mismo servidor).

## Desarrollo local sin Docker

1. Instala y arranca **MySQL** y **Redis** localmente, o usa solo los servicios de datos en Docker:

   ```bash
   docker compose up mysql redis
   ```

2. Si el backend corre **en tu máquina** y las bases de datos en Docker, apunta el host a `127.0.0.1` y los puertos publicados (**3307** para MySQL, **6380** para Redis) mediante variables de entorno (por ejemplo un `.env` en la raíz o en `backend/internal/config`, según cómo ejecutes la app).

3. Desde el directorio `backend`:

   ```bash
   go run ./cmd/api
   ```

   Asegúrate de cumplir la validación de configuración (JWT, credenciales MySQL y Redis, timeouts de Redis, etc.); revisa `internal/config/config.go` y `.env.example`.

## Tests

Por ahora la cobertura de pruebas es **reducida**; la intención es **ir añadiendo tests a medida que avance el proyecto** (unitarios, de integración con la base de datos, etc.).

Cuando quieras ejecutar los que ya existan, desde `backend`:

```bash
go test ./...
```

### CI en GitHub Actions

En [`.github/workflows/ci.yml`](.github/workflows/ci.yml) hay un workflow **CI** que, en **push y pull request hacia la rama `main`**, ejecuta en paralelo:

- **Test**: `go build ./cmd/api` y `go test ./...` dentro de `backend/`.
- **Docker Build**: `docker build -f backend/Dockerfile` desde la raíz del repositorio (respeta [`.dockerignore`](.dockerignore) para el contexto de build).

Si trabajas sobre **`develop`** u otras ramas, ese workflow **no** se dispara hasta que abras un PR contra `main` (o amplíes la lista `branches` en el YAML para incluir `develop`).

## Seguridad y secretos

- No subas **contraseñas, JWT secrets ni claves de PayPal** al repositorio. Los ficheros con secretos deben permanecer solo en tu máquina: revisa [`.gitignore`](.gitignore) (incluye `.env` y `backend/internal/config/.env`, entre otros).
- En producción usa secretos rotados, HTTPS y revisa CORS y rate limits según tu despliegue.

## Solución de problemas

- **El backend no arranca tras `docker compose up`**: espera a que los healthchecks de MySQL y Redis pasen a “healthy”; el servicio `backend` depende de ellos. Si el backend queda *unhealthy* pero responde en el navegador, revisa el healthcheck con `wget` (ver [Ficheros de Docker Compose](#ficheros-de-docker-compose)).
- **MySQL vacío o sin tablas**: el script `init.sql` solo se aplica en la **primera** creación del volumen; si ya existía un volumen antiguo, usa `docker compose down -v` y vuelve a levantar (perderás datos).
- **Error de conexión desde el host a la base**: recuerda usar el puerto **3307** (no 3306) cuando te conectes desde fuera de Docker Compose.
