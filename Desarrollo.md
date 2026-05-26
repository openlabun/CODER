# Manual de desarrollo

## 1. Propósito del documento

Este documento tiene como objetivo servir de guía técnica para comprender, mantener, extender y dar continuidad al desarrollo del proyecto. Está dirigido a futuros equipos de trabajo que necesiten familiarizarse rápidamente con la estructura del repositorio, la organización de la solución, los contenedores, los scripts, las variables de entorno y el flujo de trabajo del sistema.

## 2. Descripción general del proyecto desde la perspectiva de desarrollo

CODER es una plataforma web tipo *online judge* para la creación, despliegue y calificación automática de ejercicios, talleres y exámenes de programación, construida como un monorepo de servicios desacoplados (frontend SPA, backend en Go, worker en Python y runners por lenguaje) que se comunican mediante una cola de mensajería y persistencia externa en Roble. El reto técnico central del proyecto no se encontraba tanto en la diversidad de tecnologías a integrar, sino en sostener una **lógica de negocio robusta** —que abarca cursos, retos, exámenes, sesiones y calificación— junto con una **optimización del código** que garantizara tiempos de respuesta relativamente bajos sin sacrificar consistencia de los datos y una experiencia de usuario fluida incluso bajo carga concurrente durante evaluaciones masivas.

### 2.1 Tecnologías principales

- **Frontend:** React 19 + Vite 7, react-router-dom v7, Monaco Editor, Axios y SweetAlert2 (JavaScript / JSX).
- **Backend:** Go con el framework Fiber v2, JWT de Roble y bcrypt para autenticación, organizado bajo arquitectura hexagonal.
- **Worker de ejecución:** Python 3.11 con `aio-pika` (AMQP) y el Docker SDK para orquestar contenedores runner por lenguaje (Python, Java, C++).
- **Mensajería:** RabbitMQ 3.13 como broker AMQP, con colas independientes por lenguaje (`python.queue`, `java.queue`, `cpp.queue`).
- **Base de datos:** Roble, accedida vía API REST (`roble-api.openlab.uninorte.edu.co`), sin acceso directo a la capa SQL.
- **Inteligencia artificial:** Google Gemini en la nube para generación de contenido, con arquitectura abierta para incorporar Ollama con modelos locales (CodeLlama, DeepSeek, Mistral, Phi-3, Qwen).
- **Contenedores:** Docker Engine y Docker Compose v2; perfiles separados para servicios principales y para los runners (`judge`).
- **Infraestructura adicional:** GitHub Actions (`appleboy/ssh-action`) para CI/CD por SSH, `air` para hot-reload del backend en desarrollo y documentación OpenAPI en `/docs`.

### 2.2 Componentes principales

- **Cliente web (`apps/web`):** SPA en React que entrega la interfaz de docentes y estudiantes, integra el editor Monaco para escribir código y consume la API mediante Axios; concentra autenticación, navegación, edición y visualización de resultados.
- **Servidor/API (`apps/api_v2`):** servicio HTTP en Go que expone los endpoints de autenticación, cursos, retos, exámenes, sesiones, submissions y calificación, valida reglas de negocio y publica trabajos en RabbitMQ.
- **Worker (`apps/worker`):** proceso Python que consume las colas de códigos por revisar, lanza un contenedor runner por lenguaje, supervisa su ejecución y reporta los resultados de vuelta a la API.
- **Runners (`apps/worker/runners/{python,java,cpp}`):** contenedores efímeros y aislados que compilan o interpretan el código del estudiante, ejecutan los casos de prueba bajo límites de CPU, memoria y tiempo, y devuelven la salida estándar para evaluación.
- **Mensajería (RabbitMQ):** broker AMQP que desacopla la API del worker, mantiene una cola por lenguaje y permite absorber picos de carga durante evaluaciones masivas.
- **Base de datos (Roble):** servicio externo accedido vía REST que persiste usuarios, cursos, retos, casos de prueba, sesiones, submissions y resultados; toda la persistencia pasa por adaptadores del backend.
- **Servicios externos (Gemini / Ollama):** integración de IA generativa para asistir al docente en la creación de enunciados y casos de prueba, con la opción a futuro de operar con modelos locales sobre Ollama.

## 3. Estructura del repositorio

### 3.1 Árbol general del repositorio

```text
CODER/
├── apps/
│   ├── api_v2/                        # Backend Go (Fiber) — activo
│   │   ├── cmd/app/                   # Entry point del servidor
│   │   ├── internal/
│   │   │   ├── domain/                # Entidades, repositorios e invariantes
│   │   │   ├── application/           # Casos de uso, puertos, DI
│   │   │   ├── infrastructure/        # Adaptadores (Roble, RabbitMQ, Gemini)
│   │   │   └── interfaces/http/       # Handlers Fiber por ruta
│   │   ├── docs/                      # OpenAPI / Scalar
│   │   ├── test/                      # Tests use_cases / http / performance
│   │   ├── Dockerfile.dev | .prod
│   │   └── go.mod, go.sum
│   ├── api/                           # Backend NestJS (legado, referencia)
│   ├── web/                           # Frontend React 19 + Vite 7
│   │   ├── src/
│   │   │   ├── api/                   # Cliente Axios
│   │   │   ├── components/            # Componentes UI reutilizables
│   │   │   ├── context/               # Providers de estado global
│   │   │   ├── pages/                 # Vistas de la SPA
│   │   │   └── main.jsx, App.jsx
│   │   ├── public/
│   │   ├── Dockerfile.dev | .prod
│   │   ├── nginx.prod.conf
│   │   └── package.json, vite.config.*
│   ├── worker/                        # Worker Python (orquestador)
│   │   ├── app/                       # Lógica del worker
│   │   ├── runners/                   # Imágenes runner por lenguaje
│   │   │   ├── python/
│   │   │   ├── java/
│   │   │   └── cpp/
│   │   ├── test/
│   │   ├── main.py, config.py
│   │   ├── requirements.txt
│   │   └── Dockerfile
│   └── ai/                            # Prueba de despliegue local de LLMs
│       ├── models/                    # codellama, deepseek, mistral, phi3, qwen
│       ├── prompts/
│       └── COMMANDS.md
├── diseno/                            # Activos de diseño (diagramas, mockups)
├── docker-compose.dev.yml             # Stack de desarrollo
├── docker-compose.yml                 # Stack de producción
├── .env.dev | .env.production         # Variables de entorno
├── example.env                        # Plantilla de variables
├── FichaProyecto.md                   # Ficha académica del proyecto
├── README.md                          # Informe principal / avance
├── Instalacion.md                     # Guía de instalación y despliegue
└── Desarrollo.md                      # Manual de desarrollo (este documento)
```

### 3.2 Descripción de directorios y archivos relevantes

- **`apps/api_v2/`**: backend principal en Go (Fiber) con arquitectura hexagonal.
- **`apps/api_v2/cmd/app/`**: entry point del servidor HTTP.
- **`apps/api_v2/internal/domain/`**: entidades, repositorios e invariantes de dominio.
- **`apps/api_v2/internal/application/`**: casos de uso, puertos y contenedor de inyección de dependencias.
- **`apps/api_v2/internal/infrastructure/`**: adaptadores concretos (Roble, RabbitMQ, Gemini, JWT).
- **`apps/api_v2/internal/interfaces/http/`**: handlers Fiber organizados por ruta.
- **`apps/api_v2/docs/`**: especificación OpenAPI publicada en `/docs` con Scalar.
- **`apps/api_v2/test/`**: pruebas de casos de uso, integración HTTP y rendimiento.
- **`apps/api/scripts/`**: scripts auxiliares del backend legado (NestJS), útiles como referencia operativa y pruebas rápidas.
- **`apps/api/migrations/`**: migraciones SQL históricas del backend legado.
- **`apps/web/`**: frontend SPA en React 19 + Vite 7.
- **`apps/web/src/api/`**: cliente Axios y configuración de endpoints.
- **`apps/web/src/components/`**: componentes UI reutilizables.
- **`apps/web/src/context/`**: providers de estado global (auth, sesión, etc.).
- **`apps/web/src/pages/`**: vistas/rutas de la SPA.
- **`apps/worker/`**: orquestador Python que consume RabbitMQ y lanza runners.
- **`apps/worker/runners/`**: imágenes Docker efímeras por lenguaje (`python`, `java`, `cpp`).
- **`apps/ai/`**: prueba de despliegue local de LLMs vía Ollama (modelos y prompts).
- **`diseno/`**: activos de diseño (diagramas y mockups del sistema).
- **`.github/`**: workflows de CI/CD que despliegan a producción por SSH.
- **`docker-compose.dev.yml` / `docker-compose.yml`**: orquestación de los stacks de desarrollo y producción.
- **`.env.dev` / `.env.production` / `example.env`**: variables de entorno por ambiente y plantilla pública.

## 4. Organización de la solución a nivel de código

### 4.1 Organización por módulos o capas

El backend (`apps/api_v2`) sigue **arquitectura hexagonal**, con cuatro capas claramente separadas:

- **Dominio (`internal/domain`):** entidades, reglas de negocio y contratos de repositorio.
- **Aplicación (`internal/application`):** casos de uso y puertos que orquestan el dominio.
- **Infraestructura (`internal/infrastructure`):** adaptadores hacia Roble, RabbitMQ, Gemini y proveedores de seguridad.
- **Interfaces (`internal/interfaces/http`):** capa de presentación HTTP (handlers Fiber).

El frontend (`apps/web`) se organiza por **responsabilidades de UI**: `pages` (vistas), `components` (UI reutilizable), `context` (estado global) y `api` (cliente HTTP). El worker (`apps/worker`) es plano por su rol acotado: `main.py` + `config.py` + módulos auxiliares en `app/`.

### 4.2 Relación entre componentes del sistema y código fuente

- **Cliente web:** `apps/web/src`
- **Servidor/API:** `apps/api_v2/`
- **Worker:** `apps/worker/`
- **Runners:** `apps/worker/runners/{python,java,cpp}/`
- **Mensajería (RabbitMQ):** definida como servicio en `docker-compose.dev.yml` y `docker-compose.yml`; cliente AMQP en `apps/api_v2/internal/infrastructure/publisher/rabbitmq/` (publicación) y `apps/worker/config.py` + `apps/worker/main.py` (consumo).
- **Base de datos (Roble):** adaptadores HTTP en `apps/api_v2/internal/infrastructure/persistance/roble/`.
- **Servicios externos (Gemini / Ollama):** integraciones en `apps/api_v2/internal/infrastructure/generative-ai/` y plantillas de prueba en `apps/ai/`.

## 5. Contenedores

### 5.1 Contenedores utilizados

La siguiente tabla resume los servicios declarados en `docker-compose.dev.yml` y `docker-compose.yml`:

| Componente           | Contenedor      | Responsabilidad                                                                | Puerto                         |
| -------------------- | --------------- | ------------------------------------------------------------------------------ | ------------------------------ |
| Backend API          | `apiv2`         | Sirve la API HTTP (Fiber) y publica trabajos en RabbitMQ.                      | `4000`                         |
| Frontend Web         | `web`           | Sirve la SPA de React (Vite en dev, Nginx en prod).                            | `5173`                         |
| Broker de mensajería | `rabbitmq`      | Encola submissions por lenguaje entre la API y el worker.                      | `5672` (AMQP), `15672` (UI)    |
| Worker               | `worker`        | Consume colas de RabbitMQ y orquesta el lanzamiento de runners vía Docker SDK. | —                              |
| Runner Python        | `runner-python` | Ejecuta soluciones en Python en un contenedor aislado por submission.          | — (perfil `runners` / `judge`) |
| Runner Java          | `runner-java`   | Compila y ejecuta soluciones en Java en un contenedor aislado por submission.  | — (perfil `runners` / `judge`) |
| Runner C++           | `runner-cpp`    | Compila y ejecuta soluciones en C++ en un contenedor aislado por submission.   | — (perfil `runners` / `judge`) |

### 5.2 Archivos relacionados con contenedores

**Orquestadores (raíz del repositorio):**

- Desarrollo: `docker-compose.dev.yml`
- Producción: `docker-compose.yml`

**Backend API (`apps/api_v2`):**

- Desarrollo: `apps/api_v2/Dockerfile.dev`
- Producción: `apps/api_v2/Dockerfile.prod`

**Frontend Web (`apps/web`):**

- Desarrollo: `apps/web/Dockerfile.dev`
- Producción: `apps/web/Dockerfile.prod` (sirve el build estático con `nginx.prod.conf`)

**Worker (`apps/worker`):**

- Desarrollo y producción: `apps/worker/Dockerfile` (misma imagen para ambos entornos)

**Runners (`apps/worker/runners/`):**

- Python — Desarrollo y producción: `apps/worker/runners/python/Dockerfile`
- Java — Desarrollo y producción: `apps/worker/runners/java/Dockerfile`
- C++ — Desarrollo y producción: `apps/worker/runners/cpp/Dockerfile`

> Las imágenes de los runners son idénticas en ambos entornos; lo que cambia es el perfil de Compose que las construye: `runners` en desarrollo y `judge` en producción.

### 5.3 Construcción y ejecución de contenedores

```bash
# Construir imágenes principales (apiv2, worker, web)
docker compose -f docker-compose.dev.yml build

# Construir imágenes de los runners (perfil aparte)
docker compose -f docker-compose.dev.yml --profile runners build

# Levantar el stack
docker compose -f docker-compose.dev.yml --env-file .env.dev up -d

# Detener el stack
docker compose -f docker-compose.dev.yml down
```

Para detalle paso a paso, ver [`Instalacion.md` § 3.2](Instalacion.md#32-desarrollo-con-contenedores).

### 5.4 Redes, puertos y volúmenes

Todos los servicios comparten la red bridge `judge_net`, lo que les permite resolverse por el nombre del servicio (`apiv2`, `rabbitmq`, `worker`, `web`).

| Componente           | Contenedor      | Volumen                                                            | Puerto                     |
| -------------------- | --------------- | ------------------------------------------------------------------ | -------------------------- |
| Backend API          | `apiv2`         | `./apps/api_v2:/app` (bind-mount, hot-reload)                      | `4000:4000`                |
| Frontend Web         | `web`           | `./apps/web:/usr/src/app` + `/usr/src/app/node_modules`            | `5173:5173`                |
| Broker de mensajería | `rabbitmq`      | —                                                                  | `5672:5672`, `15672:15672` |
| Worker               | `worker`        | `./apps/worker:/app` + `/var/run/docker.sock:/var/run/docker.sock` | —                          |
| Runner Python        | `runner-python` | —                                                                  | — (instanciado on-demand)  |
| Runner Java          | `runner-java`   | —                                                                  | — (instanciado on-demand)  |
| Runner C++           | `runner-cpp`    | —                                                                  | — (instanciado on-demand)  |

**Conexiones entre servicios:**

- `web` → `apiv2`: HTTP en el puerto `4000` (configurado vía `VITE_API_URL`).
- `apiv2` → `rabbitmq`: AMQP en el puerto `5672` para publicar submissions.
- `worker` → `rabbitmq`: AMQP en el puerto `5672` para consumir las colas por lenguaje.
- `worker` → `apiv2`: HTTP en el puerto `4000` para reportar los resultados de ejecución.
- `worker` → daemon Docker: a través del socket Unix `/var/run/docker.sock` para crear contenedores runner.
- `worker` → `runner-*`: HTTP interno dentro de `judge_net` contra el `server.py` de cada runner.
- `apiv2` → Roble (externo): HTTPS al puerto `443` de `roble-api.openlab.uninorte.edu.co`.
- `apiv2` → Gemini API (externo): HTTPS al puerto `443` de los endpoints de Google.
- Operador → `rabbitmq` UI: HTTP en el puerto `15672` (consola de administración).

### 5.5 Recomendaciones para modificar contenedores

- **Mantener `Dockerfile.dev` y `Dockerfile.prod` separados.** No fusionarlos: el modo desarrollo depende de bind-mounts y hot-reload, mientras que producción busca imágenes mínimas y autocontenidas; mezclar ambos enfoques degrada el rendimiento y reintroduce vulnerabilidades.
- **Reutilizar la red `judge_net`.** Cualquier servicio nuevo debe unirse explícitamente a esta red para que la resolución por nombre (`apiv2`, `rabbitmq`, `worker`, etc.) siga funcionando; crear redes paralelas obliga a hardcodear IPs o exponer puertos innecesariamente.
- **Cuidar el montaje de `/var/run/docker.sock`.** Solo el `worker` debe tenerlo, porque lo necesita para crear runners dinámicamente; extender ese montaje a otros contenedores equivale a otorgarles acceso root al host Docker y abre la puerta a *container escape*.
- **Conservar el `healthcheck` de RabbitMQ y los `depends_on: condition: service_healthy`.** Estas condiciones evitan que la API y el worker arranquen antes que el broker; quitarlas provoca reintentos, mensajes perdidos y errores intermitentes al levantar el stack.
- **Agregar un nuevo lenguaje runner implica tres cambios coordinados.** Hay que crear el `Dockerfile` en `apps/worker/runners/<lenguaje>/`, registrar el servicio en `docker-compose.yml`/`docker-compose.dev.yml` bajo el perfil `runners`/`judge`, y declarar las variables `*_RUNNER_IMAGE` y `*_RUNNERS` para que el worker pueda construir y orquestar la pool.
- **No exponer puertos innecesarios en producción.** Usar `expose:` en lugar de `ports:` cuando un servicio solo requiere comunicación interna —como ya se hace con `apiv2` y `rabbitmq` en `docker-compose.yml`— reduce la superficie de ataque sin perder funcionalidad.
- **Las variables siempre vía archivos `.env`.** No incrustar valores en el `docker-compose.yml`: cualquier secreto hardcodeado allí queda en el historial de Git, y la sustitución por entorno permite reutilizar el mismo Compose en distintos despliegues sin tocar el repositorio.

## 6. Scripts y automatizaciones

Esta sección resume los comandos y scripts más usados en desarrollo, validación y soporte.

### 6.1 Scripts principales

- **Frontend (`apps/web/package.json`):**
   - `npm run dev`: inicia Vite en modo desarrollo.
   - `npm run build`: genera el build de producción del frontend.
   - `npm run lint`: ejecuta ESLint sobre el código del frontend.
   - `npm run preview`: sirve localmente el build generado para validación rápida.

- **Backend API v2 (Go):**
   - `go test ./test/...`: ejecuta toda la suite de pruebas del backend Go.
   - `go test -v ./test/use_cases/...`: ejecuta pruebas de casos de uso con salida detallada.
   - `go test -bench=. ./test/performance/...`: ejecuta benchmarks de rendimiento.

- **Backend legado (`apps/api/package.json`, referencia):**
   - `npm run dev`: levanta la API legada con recarga en caliente.
   - `npm run worker:submissions`: levanta el worker legado de submissions.
   - `npm run test`: ejecuta pruebas del backend legado.

### 6.2 Ubicación de scripts auxiliares

- **`apps/api/scripts/`**: scripts de soporte del backend legado (por ejemplo, `fix-enrollment-codes.ts`, pruebas de colas y worker en `.ps1` y `.sh`).
- **`apps/api/migrations/`**: scripts SQL de migración histórica.
- **`apps/api_v2/test/`**: comandos de prueba automatizada del backend actual (vía `go test`).

### 6.3 Consideraciones para su uso

- Los scripts en PowerShell (`.ps1`) están orientados a Windows; los scripts `.sh` a Linux/macOS.
- Los scripts de `apps/api/scripts/` pertenecen al backend legado y no sustituyen los flujos principales de `apps/api_v2`.
- Para pruebas del backend actual se recomienda ejecutar dentro del contenedor `apiv2` para asegurar dependencias y red consistentes.
- Antes de ejecutar scripts que interactúan con colas o runners, validar variables de entorno críticas: `RABBITMQ_URL`, `WORKER_KEY`, `RUNNER_NETWORK` y `*_RUNNER_IMAGE`.

## 7. Variables de entorno

Esta sección documenta las variables de entorno necesarias para el funcionamiento del sistema.

### 7.1 Variables requeridas

A continuación se listan únicamente las variables consumidas por el código vigente, agrupadas por componente.

#### Frontend (`apps/web`)

- `VITE_API_URL`: URL base del backend que consume el cliente Axios.
- `WEB_PORT`: puerto en el que se expone el servidor de Vite.

#### Backend API (`apps/api_v2`)

- `APIV2_PORT`: puerto HTTP en el que escucha Fiber.
- `APIV2_CORS_ALLOW_ORIGINS`: orígenes permitidos por CORS.
- `ROBLE_PROJECT`: identificador del proyecto en Roble.
- `ROBLE_BASE_URL`: URL base de la API de Roble.
- `INTERNAL_USER_EMAIL`: correo del usuario interno usado para escribir resultados de submissions.
- `INTERNAL_USER_PASSWORD`: contraseña del usuario interno (se hashea internamente).
- `WORKER_KEY`: token compartido que autentica al worker al reportar resultados.
- `SESSION_FREEZE_TIME`: segundos de congelamiento de una sesión tras un evento anti-cheat.
- `GEMINI_API_KEY`: API key de Google Gemini para la generación de contenido.
- `RABBITMQ_URL`: cadena de conexión al broker AMQP.

#### Worker (`apps/worker`)

- `RABBITMQ_URL`: cadena de conexión al broker.
- `APIV2_BASE_URL`: URL del backend a la que el worker reporta resultados.
- `WORKER_KEY`: token compartido con la API para autenticar el reporte de resultados.
- `WORKER_CONCURRENCY`: número máximo de tareas en paralelo dentro del worker.
- `PYTHON_QUEUE`: nombre de la cola de submissions de Python.
- `PYTHON_RUNNERS` / `JAVA_RUNNERS` / `CPP_RUNNERS`: cantidad de runners por lenguaje a mantener vivos.
- `PYTHON_RUNNER_IMAGE` / `JAVA_RUNNER_IMAGE` / `CPP_RUNNER_IMAGE`: nombres de las imágenes Docker que el worker lanza por lenguaje.
- `RUNNER_NETWORK`: red Docker en la que se conectan los runners.
- `RUNNER_NAME_PREFIX`: prefijo aplicado a los nombres de los contenedores runner.
- `RUNNER_HEALTH_CHECK_TIMEOUT`: timeout (s) del healthcheck contra un runner.
- `RUNNER_RESTART_TIMEOUT`: timeout (s) para reiniciar un runner caído.
- `RUNNER_REPAIR_INTERVAL`: intervalo (s) entre ciclos de reparación de la pool.

### 7.2 Variables por ambiente

Las mismas variables aplican en desarrollo y producción; lo que cambia son los valores concretos. Las diferencias observadas entre `.env.dev` y `.env.production` son:

| Variable                 | Desarrollo (`.env.dev`)               | Producción (`.env.production`)     |
| ------------------------ | ------------------------------------- | ---------------------------------- |
| `NODE_ENV`               | `development`                         | `production`                       |
| `VITE_API_URL`           | `http://localhost:4000`               | URL pública del backend            |
| `ROBLE_PROJECT`          | proyecto de pruebas (`coder_tests_…`) | proyecto de producción (`coder_…`) |
| `WORKER_CONCURRENCY`     | `8`                                   | `10`                               |
| `WORKER_KEY`             | valor de desarrollo                   | valor rotado para producción       |
| `INTERNAL_USER_PASSWORD` | valor de desarrollo                   | contraseña rotada                  |
| `GEMINI_API_KEY`         | clave de desarrollo                   | clave de producción                |

**Variables comunes a ambos ambientes** (mismo valor en `.env.dev` y `.env.production`): `APIV2_PORT`, `APIV2_BASE_URL`, `ROBLE_BASE_URL`, `PYTHON_QUEUE`, `PYTHON_RUNNERS` / `JAVA_RUNNERS` / `CPP_RUNNERS`, `*_RUNNER_IMAGE`, `RUNNER_NETWORK`, `RUNNER_NAME_PREFIX`, los tres timeouts de runners, `RABBITMQ_URL`, `SESSION_FREEZE_TIME` e `INTERNAL_USER_EMAIL`.

### 7.3 Archivos de configuración

Todos los archivos de entorno residen en la **raíz del repositorio**:

- `./.env.dev` — variables para el ambiente de desarrollo.
- `./.env.production` — variables para el ambiente de producción.
- `./example.env` — plantilla pública con los nombres de variable (sin valores reales) versionada en Git.

### 7.4 Manejo seguro de secretos

Los secretos de desarrollo (`WORKER_KEY`, `INTERNAL_USER_PASSWORD`, `GEMINI_API_KEY`, credenciales de Roble) se gestionan en el `.env.dev` local de cada desarrollador, que está incluido en `.gitignore` y nunca debe subirse al repositorio. La referencia compartida es `example.env`: un archivo plantilla con los nombres de variable y valores marcadores (`CHANGE_THIS_…`) que cada miembro del equipo copia a `.env.dev` y completa con sus propias credenciales. Cualquier rotación de un secreto se comunica por canal seguro (no por commits ni issues) y se reemplaza únicamente en el `.env.dev` local. El despliegue —con su propio juego de secretos en `.env.production` y en los *Actions secrets* del repositorio— se documenta en [`Instalacion.md` § 4.4.3](Instalacion.md#443-variables-de-entorno-y-secretos).

Los siguientes secretos deben definirse manualmente por cada desarrollador (no usar los valores marcadores del `example.env`):

| Variable                 | Tipo de secreto                 | Fuente para obtenerlo                                                  |
| ------------------------ | ------------------------------- | ---------------------------------------------------------------------- |
| `WORKER_KEY`             | Token compartido API ↔ Worker   | Generar localmente (cadena aleatoria); mismo valor en API y worker     |
| `INTERNAL_USER_PASSWORD` | Contraseña del usuario interno  | Definir localmente y registrarla en el proyecto de Roble de desarrollo |
| `GEMINI_API_KEY`         | API key de Google Gemini        | Google AI Studio (cuenta personal o del proyecto)                      |
| `ROBLE_PROJECT`          | Identificador de proyecto Roble | Panel de Roble (proyecto de desarrollo asignado al equipo)             |

> Para el funcionamiento de las variables `INTERNAL_USER_EMAIL` y `INTERNAL_USER_PASSWORD` debe crearse un usuario administrador en Roble con esas mismas credenciales.

## 8. Flujo de trabajo de desarrollo

Esta sección describe el proceso recomendado para trabajar sobre el proyecto.

### 8.1 Preparación del entorno

1. **Clonar el repositorio.**
   
   ```bash
   git clone https://github.com/openlabun/CODER.git
   cd CODER
   ```

2. **Configurar variables de entorno.** Copiar la plantilla a la raíz del repositorio y completar los valores locales (ver [7.4](#74-manejo-seguro-de-secretos)).
   
   ```bash
   cp example.env .env.dev
   ```
   
   Ubicación esperada: `./.env.dev` en la raíz del repositorio.

3. **Construir las imágenes principales.**
   
   ```bash
   docker compose -f docker-compose.dev.yml build
   ```

4. **Construir las imágenes de los runners.**
   
   ```bash
   docker compose -f docker-compose.dev.yml --profile runners build
   ```

5. **Levantar el stack de desarrollo.**
   
   ```bash
   docker compose -f docker-compose.dev.yml --env-file .env.dev up -d
   ```

6. **Verificar los servicios disponibles:** `http://localhost:5173` (frontend), `http://localhost:4000/docs` (API + Scalar) y `http://localhost:15672` (consola de RabbitMQ).

### 8.2 Desarrollo de nuevas funcionalidades

Para incorporar una nueva funcionalidad se debe seguir el siguiente orden de trabajo:

1. **Crear una nueva rama** con el formato `feature/[nombre-de-la-funcionalidad]` a partir de la rama de integración actual.
   
   ```bash
   git checkout -b feature/nombre-de-la-funcionalidad
   ```

2. **Modificar el dominio y la capa de aplicación.** Si la funcionalidad afecta la lógica de negocio, los cambios se realizan en `apps/api_v2/internal/domain/` (entidades, invariantes, contratos de repositorio) y `apps/api_v2/internal/application/` (casos de uso, puertos y DI).

3. **Modificar la capa de infraestructura.** Si la funcionalidad introduce una nueva conexión externa (Roble, RabbitMQ, Gemini, etc.) o cambia el comportamiento de las existentes, se ajustan los adaptadores en `apps/api_v2/internal/infrastructure/`.

4. **Modificar la capa de interfaces.** Si la funcionalidad expone un nuevo endpoint, se crea un paquete bajo `apps/api_v2/internal/interfaces/http/<modulo>/<verbo-recurso>/` (con sus archivos `route.go`, `dtos.go` y `mapper.go` siguiendo la convención existente) y se registra el handler en `apps/api_v2/internal/interfaces/http/routes.go`.

5. **Agregar un test de funcionamiento.** Crear pruebas en `apps/api_v2/test/` (preferentemente bajo `test/http/` para validar el endpoint completo) reutilizando los helpers expuestos en `test/http/utils.go`.
   
   ```bash
   cd apps/api_v2 && go test ./test/http/...
   ```

6. **Actualizar la interfaz.** Una vez verificado el backend, ajustar el frontend en `apps/web/src/` (cliente Axios en `src/api/`, componentes en `src/components/` y vistas en `src/pages/`) para consumir la nueva funcionalidad.

### 8.3 Ejecución de pruebas y validaciones

Las pruebas de la API se ejecutan **dentro del contenedor `apiv2`**, que ya tiene el toolchain de Go y el acceso a la red `judge_net` necesarios para alcanzar RabbitMQ y Roble.

```bash
# Stack de desarrollo arriba
docker compose -f docker-compose.dev.yml --env-file .env.dev up -d

# Ejecutar toda la batería de pruebas dentro del contenedor
docker compose -f docker-compose.dev.yml exec apiv2 go test ./test/...

# Ejecutar un grupo específico (ej. casos de uso)
docker compose -f docker-compose.dev.yml exec apiv2 go test -v ./test/use_cases/...

# Ejecutar una prueba puntual por nombre
docker compose -f docker-compose.dev.yml exec apiv2 go test -v ./test/use_cases/exam -run TestCourseCRUD

# Benchmarks de rendimiento
docker compose -f docker-compose.dev.yml exec apiv2 go test -bench=. ./test/performance/...
```

Validaciones del frontend (lint y build), también dentro del contenedor:

```bash
docker compose -f docker-compose.dev.yml exec web npm run lint
docker compose -f docker-compose.dev.yml exec web npm run build
```

El catálogo completo de pruebas disponibles, junto con la descripción paso a paso de cada caso y los comandos `go test -run` correspondientes, se encuentra en [`apps/api_v2/docs/TestsPlan.md`](apps/api_v2/docs/TestsPlan.md). Categorías cubiertas en ese plan:

- **Pruebas de Seguridad** — acceso y registro de usuarios contra Roble (`./test/roble_auth`).
- **Pruebas de Persistencia** — CRUD completo de cursos, exámenes y revisiones (`./test/roble_persistence`).
- **Pruebas Funcionales** — autenticación, cursos, exámenes y revisiones a nivel de casos de uso (`./test/use_cases/...`).
- **Pruebas de Rendimiento** — benchmarks con `testing.B` y goroutines paralelas (`./test/performance/...`).

### 8.4 Integración de cambios

- **Publicación de la rama:** se publica la rama `feature/...` dividiendo cada commit por **acciones específicas**, usando la nomenclatura estándar de tipo de cambio: `feat/FEATURE` para nuevas funcionalidades, `fix/FIX` para correcciones, `docs/DOCUMENTATION` para documentación, `refactor/REFACTOR` para reestructuración interna, `test/TEST` para pruebas, etc. Cada commit debe representar una acción atómica y trazable.
- **Creación del Pull Request:** se abre un PR contra la rama de integración describiendo los cambios principales, las decisiones técnicas tomadas y adjuntando **imágenes de la ejecución de las pruebas** (capturas del `go test`, screenshots del frontend, etc.) que evidencien el correcto funcionamiento.
- **Manejo de correcciones:** las **correcciones del trabajo** se hacen sobre el propio PR (nuevos commits o ajustes a los existentes), mientras que las **correcciones a una *issue*** se gestionan en la página de la issue. En ambos casos, las observaciones deben acompañarse de **outputs, logs o capturas** que demuestren el funcionamiento incorrecto reportado o la solución aplicada.
- **Criterios mínimos para integrar:**
   - Funcionamiento correcto de la nueva funcionalidad, demostrado mediante los tests añadidos en [8.2](#82-desarrollo-de-nuevas-funcionalidades).
  - Evidencia de que `docker-compose.dev.yml` sigue levantando el stack sin errores tras los cambios.
  - El resto de la batería de pruebas de funcionalidad (`./test/use_cases/...`, `./test/http/...`) sigue pasando sin regresiones.

## 9. Dependencias y servicios externos

Documente las dependencias técnicas y servicios de terceros utilizados por el proyecto.

### 9.1 Servicios externos integrados

- **Roble (Openlab, Universidad del Norte):** servicio institucional accedido por API REST que provee tanto la autenticación de usuarios como la persistencia de todas las entidades del sistema (cursos, retos, exámenes, sesiones, submissions y resultados).
- **API de Google Gemini:** modelo de lenguaje en la nube usado por el backend para asistir al docente en la generación de enunciados, casos de prueba y demás contenido académico.

### 9.2 Requisitos de acceso

Un equipo nuevo solo requiere dos accesos para empezar a trabajar:

- **Repositorio del proyecto** en GitHub (`openlabun/CODER`), con permisos de colaborador.
- **Proyectos de Roble** correspondientes a las bases de datos de **desarrollo** y **producción**, asignados por el administrador de Openlab.

### 9.3 Consideraciones de desarrollo y pruebas

- **Helpers de pruebas:** todas las pruebas de la API deben apoyarse en las funciones definidas en [`apps/api_v2/test/http/utils.go`](apps/api_v2/test/http/utils.go), particularmente para la **medición de tiempos** de ejecución de los casos de uso. Esto garantiza coherencia entre todas las pruebas y permite comparar tiempos de respuesta de forma consistente.

- **Base de datos en desarrollo:** los despliegues de desarrollo deben usar **el proyecto de Roble destinado a desarrollo**, no el de producción. Las pruebas suelen generar un alto volumen de modificaciones (creación y borrado de cursos, exámenes, retos, sesiones y submissions) y operar contra el proyecto productivo contaminaría los datos reales.

- **Usuarios de prueba predeterminados:** las pruebas funcionales asumen la existencia de dos cuentas en el Roble de desarrollo, con contraseña `Password123!`:
  
  - `test@test.com` — usuario con permisos de **docente**.
  - `stud@test.com` — usuario con permisos de **estudiante**.
  
  Estas cuentas deben crearse en el proyecto de Roble antes de ejecutar las pruebas, tal como lo indica [`apps/api_v2/docs/TestsPlan.md`](apps/api_v2/docs/TestsPlan.md).

## 10. Convenciones del proyecto

### 10.1 Convenciones de código

La nomenclatura del dominio se mantiene en **PascalCase en inglés** y debe respetarse en cualquier nueva entidad o atributo. El detalle completo de cada entidad y sus reglas vive en [`apps/api_v2/docs/Bussiness_rules.md`](apps/api_v2/docs/Bussiness_rules.md); aquí se resume la nomenclatura por módulo:

**Módulo Usuarios:**

- `User`: entidad base de cualquier persona registrada en la plataforma.
- `Role`: rol del usuario, con valores `student`, `professor` o `admin`.

**Módulo Cursos:**

- `Course`: curso académico al que pertenecen exámenes y estudiantes.
- `CourseStudent`: relación de inscripción entre un curso y un estudiante.
- `EnrollmentCode` / `EnrollmentURL`: credenciales de acceso a un curso privado.
- `Visibility`: configuración con valores `public`, `private` o `blocked`.
- `Period`: semestre del curso con valores `10` (Primer Semestre), `20` (Intersemestral) y `30` (Segundo Semestre).

**Módulo Exámenes:**

- `Exam`: agrupador de retos que el estudiante resolverá como evaluación.
- `Challenge`: enunciado individual de programación, con estado `draft`, `published`, `private` o `archived`.
- `ExamItem`: punto que vincula un `Exam` con un `Challenge` (relación muchos a muchos) e indica `Order` y `Points`.
- `TestCase`: caso de prueba asociado a un `Challenge`, con atributo `IsSample` para distinguir entre ejemplos visibles y casos secretos.
- `IOVariable`: variable de entrada/salida con `Name`, `Type` (`string`, `int`, `float`) y `Value`.
- `Template`: plantilla de código por lenguaje (`cpp`, `python`, `java`) entregada al estudiante.

**Módulo Revisiones:**

- `Submission`: envío de código de un estudiante para un `ExamItem`, vinculado siempre a `Session` y `User`.
- `SubmissionResult`: resultado por cada `TestCase` ejecutado, con estado `queued`, `running`, `accepted`, `wrong_answer` o `error`.
- `Session`: sesión de examen con estados `active`, `frozen`, `completed`, `expired` o `blocked`; usa `Heartbeat` para mantenerse viva.
- `ExamScore`: puntaje agregado del examen por usuario, calculado al cerrar la sesión (relación 1 a 1 con `Session`).
- `ExamItemScore`: puntaje por punto del examen, dependiente de `ExamScore` y vinculado al `ExamItem` correspondiente.

Adicionalmente, los archivos de handlers HTTP siguen la convención `<verbo>-<recurso>` para el nombre del paquete (ej. `post-create`, `get-by-id`, `delete-exam`), y cada paquete expone `route.go`, `dtos.go` y `mapper.go` por separado.

### 10.2 Convenciones de repositorio

#### Nombres de ramas

Cada rama debe describir su intención mediante uno de los siguientes prefijos:

- `feature/[nombre-de-la-funcionalidad]`: nueva funcionalidad.
- `documentation/[informes]`: cambios en documentación o informes.
- `fix/[error-corregido]`: corrección de un bug.
- `updates/[sistema-actualizado]`: actualización de un componente, dependencia o módulo existente.

#### Mensajes de commit

Los mensajes de commit deben iniciar con la palabra clave que clasifica el cambio:

- `FEATURE`: implementación de una funcionalidad nueva.
- `ADD`: archivos o recursos auxiliares que soportan una funcionalidad (assets, helpers, tipos, mocks).
- `DOCUMENTATION`: cambios en la documentación del proyecto.
- `FIX`: corrección de errores.
- `REFACTOR`: reestructuración interna del código sin cambiar su comportamiento.

Cada commit debe representar una acción atómica, lo que permite mantener un historial limpio y trazable cuando se publica la rama.

#### Manejo de issues

Las issues se clasifican en tres categorías según su propósito:

- **Funcionalidades:** describen qué se quiere implementar. Deben incluir una **breve descripción** del objetivo, las **I/O considerations** (entradas, salidas y validaciones esperadas) y los **pasos a seguir** para llegar al resultado.
- **Pruebas:** documentan los archivos o la **secuencia a seguir** para probar una funcionalidad, incluyendo precondiciones, comandos y resultados esperables.
- **Errores:** registran un bug con el **paso a paso de cómo se consiguió**, el **resultado esperado** y el **resultado obtenido**, acompañados de logs, capturas o trazas que evidencien el comportamiento incorrecto.

### 10.3 Convenciones de documentación

La documentación técnica del proyecto se concentra en cuatro ubicaciones, que deben mantenerse actualizadas a medida que el código evoluciona:

- [`apps/api_v2/docs/Bussiness_rules.md`](apps/api_v2/docs/Bussiness_rules.md): reglas de negocio del sistema (entidades, estados, restricciones e invariantes por módulo).
- [`apps/api_v2/docs/TestsPlan.md`](apps/api_v2/docs/TestsPlan.md): plan de pruebas con la descripción y los comandos `go test` para cada caso.
- [`apps/api_v2/docs/openapi.yaml`](apps/api_v2/docs/openapi.yaml): especificación OpenAPI de la API, fuente de la UI publicada en `/docs` con Scalar.
- [`apps/api_v2/docs/informs/`](apps/api_v2/docs/informs/): carpeta con los **informes de todos los cambios mayores** del proyecto (un archivo por hito o cambio relevante).

Cualquier funcionalidad nueva, ajuste de reglas o cambio mayor debe reflejarse en el archivo correspondiente como parte del mismo PR.

## 11. Problemas frecuentes y recomendaciones

Documente errores comunes, limitaciones conocidas, deuda técnica o advertencias importantes para futuros equipos.

### 11.1 Problemas frecuentes

Durante el desarrollo, los dos problemas que más recurrentemente bloquean al equipo son los relacionados con la **política de CORS** y con la **configuración de los runners**.

#### Errores de CORS Policy

El navegador rechaza las peticiones del frontend al backend con un mensaje del tipo `blocked by CORS policy` cuando el origen desde el que se sirve la SPA no está incluido en los orígenes permitidos por la API.

- **Causa típica:** la variable `APIV2_CORS_ALLOW_ORIGINS` no incluye el origen real desde el que se carga el frontend (por ejemplo, `http://localhost:5173` en desarrollo o el dominio público en producción), o quedó en `*` mientras el cliente envía cabeceras autenticadas (combinación inválida según la spec de CORS).
- **Cómo verificarlo:** abrir la consola del navegador (`F12 → Network`), revisar la respuesta de la petición fallida y validar las cabeceras `Access-Control-Allow-Origin` y `Access-Control-Allow-Credentials` devueltas por la API.
- **Solución:** ajustar `APIV2_CORS_ALLOW_ORIGINS` en el `.env.*` correspondiente al origen exacto del frontend y reiniciar `apiv2`. En producción, evitar `*` y listar explícitamente los dominios autorizados.

#### Configuración de los runners

El worker no logra ejecutar submissions porque los runners no están disponibles o no pueden ser lanzados por Docker.

- **Imágenes no construidas:** las imágenes `coder-runner-python`, `coder-runner-java` y `coder-runner-cpp` viven bajo el perfil `runners` (dev) o `judge` (prod) y **no se construyen automáticamente** con `docker compose build`. Solución: ejecutar `docker compose -f docker-compose.dev.yml --profile runners build` (o `--profile judge` en producción).
- **Variables desalineadas con las imágenes:** si los nombres de imagen en `PYTHON_RUNNER_IMAGE`, `JAVA_RUNNER_IMAGE` o `CPP_RUNNER_IMAGE` no coinciden con los nombres construidos por Docker, el worker falla al lanzar el contenedor. Solución: verificar que las variables del `.env` coincidan exactamente con las imágenes listadas por `docker images | grep coder-runner`.
- **Red `judge_net` ausente:** los runners se lanzan en la red declarada por `RUNNER_NETWORK`; si esa red no existe (o se eliminó con `docker compose down -v`), el worker no puede conectarlos. Solución: volver a levantar el stack con `docker compose up -d` o crear la red manualmente con `docker network create judge_net`.
- **Acceso al socket Docker:** el worker requiere `/var/run/docker.sock` montado para crear los runners; en Windows/macOS con Docker Desktop, ese montaje funciona transparentemente, pero en hosts Linux endurecidos puede requerir permisos adicionales. Solución: confirmar que el contenedor `worker` tenga el volumen `/var/run/docker.sock:/var/run/docker.sock` y que el usuario que ejecuta Docker pertenezca al grupo `docker`.
- **Timeouts demasiado bajos:** valores agresivos en `RUNNER_HEALTH_CHECK_TIMEOUT`, `RUNNER_RESTART_TIMEOUT` o `RUNNER_REPAIR_INTERVAL` provocan que el worker considere muertos a runners que solo están bajo carga. Solución: ajustar estos valores en el `.env` cuando se observen reinicios continuos en los logs del worker.

#### Inicio lento de RabbitMQ

Al levantar el stack, RabbitMQ puede tardar más de lo previsto en estar listo, lo que hace que `apiv2` y `worker` —que esperan al broker mediante `depends_on: condition: service_healthy`— fallen al arrancar y reporten errores de conexión AMQP.

- **Causa típica:** el `healthcheck` del servicio `rabbitmq` tiene un `interval`/`timeout` demasiado corto para la máquina donde se está desplegando; el contenedor de RabbitMQ alcanza a iniciar el proceso Erlang pero aún no responde a `rabbitmq-diagnostics -q ping` cuando el healthcheck se evalúa, por lo que Docker lo marca como `unhealthy` y los demás contenedores no arrancan.
- **Cómo verificarlo:** ejecutar `docker compose ps` y observar el estado `unhealthy` del contenedor `rabbitmq`; revisar `docker compose logs rabbitmq` para confirmar que el proceso terminó de inicializar.
- **Solución:** subir el `timeout` (y, si es necesario, el `interval` y `retries`) del `healthcheck` de `rabbitmq` en `docker-compose.yml` / `docker-compose.dev.yml` a valores holgados (por ejemplo, `timeout: 60s`, `interval: 60s`, `retries: 20`) para darle tiempo al servicio a estar realmente listo antes de que se valide su estado. Sin este ajuste, el despliegue falla de forma intermitente y los contenedores dependientes no llegan a iniciarse.

### 11.2 Deuda técnica conocida

- **Rehacer el frontend desde cero.** La SPA actual no cumple con los estándares de calidad ni con las buenas prácticas esperadas: presenta un fuerte acoplamiento entre componentes, lógica de negocio embebida y reglas hardcodeadas que **duplican —y a veces contradicen— las que ya existen en la API**. Se recomienda reescribirlo completamente, manteniendo a la API como única fuente de verdad para reglas de negocio y limitando el frontend a responsabilidades de UI y consumo HTTP.
- **Migrar o aliviar la persistencia.** Roble se está volviendo un cuello de botella: cada consulta es una llamada HTTP externa, lo que limita el rendimiento bajo carga. Existen dos caminos:
  - **Migrar a un servicio de base de datos más eficiente** (PostgreSQL gestionado, por ejemplo) con acceso directo desde la API.
  - **Implementar una capa de caché local** que reduzca el número de llamadas a Roble manteniéndolo como sistema de registro. Hay un avance significativo de esta infraestructura de caché en la rama [`feature/cache-infrastructure`](https://github.com/openlabun/CODER/tree/feature/cache-infrastructure), que define el adaptador, las TTL por tabla (`CACHE_*_TABLE_TTL`) y la cola de sincronización (`CACHE_SYNC_QUEUE_NAME`); falta integrar y validar este trabajo con el resto del sistema.

### 11.3 Recomendaciones para continuidad

Para los siguientes equipos que tomen el proyecto, se sugieren las siguientes líneas de trabajo:

- **Mejorar los tiempos de respuesta de la API.** Optimizar el camino crítico de las operaciones más frecuentes (login, carga de examen, envío de submission y consulta de resultados), apoyándose en las pruebas de rendimiento existentes (`./test/performance/...`) para medir el impacto de cada cambio. La capa de caché de [`feature/cache-infrastructure`](https://github.com/openlabun/CODER/tree/feature/cache-infrastructure) es un buen punto de partida.
- **Reconstruir completamente el frontend.** Tomar como base el manual de desarrollo y la documentación de la API (`/docs`) para producir una nueva SPA limpia, sin lógica de negocio hardcodeada y con un acoplamiento bajo entre componentes (ver [11.2](#112-deuda-técnica-conocida)).
- **Implementar runners para Java y C++.** La infraestructura ya contempla las imágenes y las variables (`JAVA_RUNNER_IMAGE`, `CPP_RUNNER_IMAGE`, `JAVA_RUNNERS`, `CPP_RUNNERS`); falta completar la ejecución, parsing de salida y manejo de errores de compilación dentro del worker para cada lenguaje.
- **Implementar detección de plagio por comparación de código.** Diseñar un módulo que compare las submissions de un mismo `ExamItem` mediante técnicas como tokenización con n-grams (estilo Moss/JPlag) o comparación estructural por AST, y exponer alertas para revisión manual por parte del docente. El informe principal ya describe el enfoque sugerido.
- **Migrar el flujo del examen de HTTP request a WebSocket.** El modelo actual basado en *heartbeats* HTTP y `polling` introduce latencia y carga innecesaria sobre la API; migrar la sesión de examen a una conexión WebSocket permitiría notificaciones en tiempo real (bloqueo manual del docente, expiración de la sesión, resultados de submission), reducir el tráfico HTTP y mejorar la experiencia del estudiante.

## 12. Historial de decisiones técnicas relevantes

- **Reconstrucción total del backend en Go.** Se descartó el backend original en NestJS/TypeScript (`apps/api`) y se reescribió desde cero en Go con Fiber v2 bajo arquitectura hexagonal (`apps/api_v2`), buscando mayor rendimiento bajo carga, una separación más estricta de capas y un binario autocontenido más sencillo de desplegar.
- **Persistencia y autenticación delegadas a Roble.** Se eliminó la dependencia directa de PostgreSQL en el backend y se integró Roble como única fuente de verdad para persistencia y autenticación (JWT emitido por Roble), reduciendo la infraestructura propia a mantener y aprovechando el servicio institucional.
- **Worker en Python con runners on-demand.** Se introdujo un worker dedicado en Python que orquesta contenedores runner por lenguaje, aislando la ejecución del código del estudiante del proceso de la API y permitiendo evolucionar la lógica de ejecución sin tocar el backend.
- **Adopción de RabbitMQ como Message Broker.** Se incorporó RabbitMQ entre la API y el worker para desacoplar la publicación de las submissions de su ejecución, manejar concurrencia de forma natural y absorber picos de carga durante evaluaciones masivas.
- **Infraestructura de comunicación con la API de Gemini.** Se diseñó una capa de adaptadores (`internal/infrastructure/generative-ai/cloud/gemini`) para integrar Gemini como proveedor de IA generativa, dejando el contrato listo para soportar otros proveedores (Ollama local) sin tocar los casos de uso.
- **Migración de runners on-demand a *pool* precargado.** Para alcanzar el objetivo de **25 solicitudes concurrentes**, se reemplazó el esquema de crear un runner por submission por una *pool* de runners precargados (`PYTHON_RUNNERS`, `JAVA_RUNNERS`, `CPP_RUNNERS`) que el worker mantiene vivos y reutiliza, eliminando el costo de arranque de un contenedor por cada ejecución y permitiendo procesar solicitudes en paralelo.

## 13. Referencias relacionadas

**Documentación interna:**

- [`apps/api_v2/docs/Bussiness_rules.md`](apps/api_v2/docs/Bussiness_rules.md) — reglas de negocio del sistema, entidades y estados por módulo.
- [`apps/api_v2/docs/TestsPlan.md`](apps/api_v2/docs/TestsPlan.md) — plan de pruebas con descripción y comandos `go test` por caso.
- [`apps/api_v2/docs/openapi.yaml`](apps/api_v2/docs/openapi.yaml) — especificación OpenAPI de la API (UI servida en `/docs` con Scalar).
- [`apps/api_v2/docs/informs/Worker-changes.md`](apps/api_v2/docs/informs/Worker-changes.md) — informe de la migración del worker a *pool* de runners precargados.
- [`README.md`](README.md) — informe principal del proyecto.
- [`Instalacion.md`](Instalacion.md) — guía de instalación y despliegue.
- [`FichaProyecto.md`](FichaProyecto.md) — ficha académica del proyecto.

**Servicios externos:**

- [Roble — Documentación oficial](https://roble.openlab.uninorte.edu.co/docs).
- [Google Gemini API](https://ai.google.dev/gemini-api/docs).

**Tecnologías del stack:**

- [Fiber v2 — Documentación oficial](https://docs.gofiber.io/) — framework HTTP usado por el backend.
- [RabbitMQ — Documentación oficial](https://www.rabbitmq.com/docs) — broker AMQP usado entre la API y el worker.

**Repositorios y ramas relacionadas:**

- [juez-online — Repositorio base](https://github.com/DerekPz/juez-online) — proyecto del que se partió originalmente.
- [`feature/cache-infrastructure`](https://github.com/openlabun/CODER/tree/feature/cache-infrastructure) — rama con el avance de la capa de caché sobre Roble.
