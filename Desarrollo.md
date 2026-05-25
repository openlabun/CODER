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
│   │   ├── scripts/                   # Utilidades de soporte
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

## 6. Variables de entorno

Documente las variables de entorno necesarias para el funcionamiento del sistema.

### 6.1 Variables requeridas

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

### 6.2 Variables por ambiente

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

### 6.3 Archivos de configuración

Todos los archivos de entorno residen en la **raíz del repositorio**:

- `./.env.dev` — variables para el ambiente de desarrollo.
- `./.env.production` — variables para el ambiente de producción.
- `./example.env` — plantilla pública con los nombres de variable (sin valores reales) versionada en Git.

### 6.4 Manejo seguro de secretos

Los secretos de desarrollo (`WORKER_KEY`, `INTERNAL_USER_PASSWORD`, `GEMINI_API_KEY`, credenciales de Roble) se gestionan en el `.env.dev` local de cada desarrollador, que está incluido en `.gitignore` y nunca debe subirse al repositorio. La referencia compartida es `example.env`: un archivo plantilla con los nombres de variable y valores marcadores (`CHANGE_THIS_…`) que cada miembro del equipo copia a `.env.dev` y completa con sus propias credenciales. Cualquier rotación de un secreto se comunica por canal seguro (no por commits ni issues) y se reemplaza únicamente en el `.env.dev` local. El despliegue —con su propio juego de secretos en `.env.production` y en los *Actions secrets* del repositorio— se documenta en [`Instalacion.md` § 4.4.3](Instalacion.md#443-variables-de-entorno-y-secretos).

Los siguientes secretos deben definirse manualmente por cada desarrollador (no usar los valores marcadores del `example.env`):

| Variable                 | Tipo de secreto                 | Fuente para obtenerlo                                                  |
| ------------------------ | ------------------------------- | ---------------------------------------------------------------------- |
| `WORKER_KEY`             | Token compartido API ↔ Worker   | Generar localmente (cadena aleatoria); mismo valor en API y worker     |
| `INTERNAL_USER_PASSWORD` | Contraseña del usuario interno  | Definir localmente y registrarla en el proyecto de Roble de desarrollo |
| `GEMINI_API_KEY`         | API key de Google Gemini        | Google AI Studio (cuenta personal o del proyecto)                      |
| `ROBLE_PROJECT`          | Identificador de proyecto Roble | Panel de Roble (proyecto de desarrollo asignado al equipo)             |

> Para el funcionamiento de las variables `INTERNAL_USER_EMAIL` y `INTERNAL_USER_PASSWORD` se deberá crear un usuario administrador en Roble con esas mismas credenciales

## 7. Flujo de trabajo de desarrollo

Describa el proceso recomendado para trabajar sobre el proyecto.

### 7.1 Preparación del entorno

1. **Clonar el repositorio.**
   
   ```bash
   git clone https://github.com/openlabun/CODER.git
   cd CODER
   ```

2. **Configurar variables de entorno.** Copiar la plantilla a la raíz del repositorio y completar los valores locales (ver [6.4](#64-manejo-seguro-de-secretos)).
   
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

### 7.2 Desarrollo de nuevas funcionalidades

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

### 7.3 Ejecución de pruebas y validaciones

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

### 7.4 Integración de cambios

- **Publicación de la rama:** se publica la rama `feature/...` dividiendo cada commit por **acciones específicas**, usando la nomenclatura estándar de tipo de cambio: `feat/FEATURE` para nuevas funcionalidades, `fix/FIX` para correcciones, `docs/DOCUMENTATION` para documentación, `refactor/REFACTOR` para reestructuración interna, `test/TEST` para pruebas, etc. Cada commit debe representar una acción atómica y trazable.
- **Creación del Pull Request:** se abre un PR contra la rama de integración describiendo los cambios principales, las decisiones técnicas tomadas y adjuntando **imágenes de la ejecución de las pruebas** (capturas del `go test`, screenshots del frontend, etc.) que evidencien el correcto funcionamiento.
- **Manejo de correcciones:** las **correcciones del trabajo** se hacen sobre el propio PR (nuevos commits o ajustes a los existentes), mientras que las **correcciones a una *issue*** se gestionan en la página de la issue. En ambos casos, las observaciones deben acompañarse de **outputs, logs o capturas** que demuestren el funcionamiento incorrecto reportado o la solución aplicada.
- **Criterios mínimos para integrar:**
  - Funcionamiento correcto de la nueva funcionalidad, demostrado mediante los tests añadidos en [7.2](#72-desarrollo-de-nuevas-funcionalidades).
  - Evidencia de que `docker-compose.dev.yml` sigue levantando el stack sin errores tras los cambios.
  - El resto de la batería de pruebas de funcionalidad (`./test/use_cases/...`, `./test/http/...`) sigue pasando sin regresiones.

## 8. Dependencias y servicios externos

Documente las dependencias técnicas y servicios de terceros utilizados por el proyecto.

### 8.1 Servicios externos integrados

- **Roble (Openlab, Universidad del Norte):** servicio institucional accedido por API REST que provee tanto la autenticación de usuarios como la persistencia de todas las entidades del sistema (cursos, retos, exámenes, sesiones, submissions y resultados).
- **API de Google Gemini:** modelo de lenguaje en la nube usado por el backend para asistir al docente en la generación de enunciados, casos de prueba y demás contenido académico.

### 8.2 Requisitos de acceso

Un equipo nuevo solo requiere dos accesos para empezar a trabajar:

- **Repositorio del proyecto** en GitHub (`openlabun/CODER`), con permisos de colaborador.
- **Proyectos de Roble** correspondientes a las bases de datos de **desarrollo** y **producción**, asignados por el administrador de Openlab.

### 8.3 Consideraciones de desarrollo y pruebas

- **Helpers de pruebas:** todas las pruebas de la API deben apoyarse en las funciones definidas en [`apps/api_v2/test/http/utils.go`](apps/api_v2/test/http/utils.go), particularmente para la **medición de tiempos** de ejecución de los casos de uso. Esto garantiza coherencia entre todas las pruebas y permite comparar tiempos de respuesta de forma consistente.

- **Base de datos en desarrollo:** los despliegues de desarrollo deben usar **el proyecto de Roble destinado a desarrollo**, no el de producción. Las pruebas suelen generar un alto volumen de modificaciones (creación y borrado de cursos, exámenes, retos, sesiones y submissions) y operar contra el proyecto productivo contaminaría los datos reales.

- **Usuarios de prueba predeterminados:** las pruebas funcionales asumen la existencia de dos cuentas en el Roble de desarrollo, con contraseña `Password123!`:
  
  - `test@test.com` — usuario con permisos de **docente**.
  - `stud@test.com` — usuario con permisos de **estudiante**.
  
  Estas cuentas deben crearse en el proyecto de Roble antes de ejecutar las pruebas, tal como lo indica [`apps/api_v2/docs/TestsPlan.md`](apps/api_v2/docs/TestsPlan.md).

## 9. Convenciones del proyecto

### 9.1 Convenciones de código

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

### 9.2 Convenciones de repositorio

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

### 9.3 Convenciones de documentación

La documentación técnica del proyecto se concentra en cuatro ubicaciones, que deben mantenerse actualizadas a medida que el código evoluciona:

- [`apps/api_v2/docs/Bussiness_rules.md`](apps/api_v2/docs/Bussiness_rules.md): reglas de negocio del sistema (entidades, estados, restricciones e invariantes por módulo).
- [`apps/api_v2/docs/TestsPlan.md`](apps/api_v2/docs/TestsPlan.md): plan de pruebas con la descripción y los comandos `go test` para cada caso.
- [`apps/api_v2/docs/openapi.yaml`](apps/api_v2/docs/openapi.yaml): especificación OpenAPI de la API, fuente de la UI publicada en `/docs` con Scalar.
- [`apps/api_v2/docs/informs/`](apps/api_v2/docs/informs/): carpeta con los **informes de todos los cambios mayores** del proyecto (un archivo por hito o cambio relevante).

Cualquier funcionalidad nueva, ajuste de reglas o cambio mayor debe reflejarse en el archivo correspondiente como parte del mismo PR.

## 10. Problemas frecuentes y recomendaciones

Documente errores comunes, limitaciones conocidas, deuda técnica o advertencias importantes para futuros equipos.

### 11.1 Problemas frecuentes

Ejemplo:

- errores de puertos;
- variables de entorno faltantes;
- problemas de conexión a la base de datos;
- conflictos entre versiones;
- errores de permisos.

### 11.2 Deuda técnica conocida

Liste componentes incompletos, decisiones provisionales, refactors pendientes o limitaciones actuales del sistema.

### 11.3 Recomendaciones para continuidad

Indique sugerencias concretas para futuros grupos que deban continuar el proyecto.

## 12. Historial de decisiones técnicas relevantes

Documente decisiones importantes tomadas durante el desarrollo y la razón detrás de ellas.

Ejemplo:

- elección de framework;
- cambio de base de datos;
- adopción o descarte de contenedores;
- reestructuración de módulos;
- cambio en estrategia de autenticación.

## 13. Referencias relacionadas

Incluya enlaces o referencias útiles para continuar el desarrollo.

Ejemplo:

- documentación oficial de tecnologías usadas;
- enlaces a servicios externos;
- documentación de arquitectura;
- instalación del proyecto;
- informe principal.