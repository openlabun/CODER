# Instalación y despliegue

## 1. Descripción general de la solución

### 1.1 Lenguajes y tecnologías utilizadas

| Componente                | Lenguaje            | Framework / Tecnología principal                                    |
| ------------------------- | ------------------- | ------------------------------------------------------------------- |
| Frontend Web              | JavaScript          | React 19 + Vite 7, react-router-dom v7, Monaco Editor, Axios        |
| Backend API               | Go 1.25             | Fiber v2, JWT, bcrypt                                               |
| Worker (orquestador)      | Python 3            | aio-pika (AMQP), Docker SDK                                         |
| Runners de ejecución      | Python / Java / C++ | Imágenes Docker aisladas (`coder-runner-python`, `-java`, `-cpp`)   |
| Broker de mensajería      | —                   | RabbitMQ 3.13                                                       |
| Base de datos             | —                   | Roble (API REST `roble-api.openlab.uninorte.edu.co`)                |
| Asistente de IA           | —                   | Ollama (modelos locales: CodeLlama, DeepSeek, Mistral, Phi-3, Qwen) |
| Orquestación de servicios | —                   | Docker Compose                                                      |

**Tabla:** Lenguajes y tecnologías por componente

### 1.2 Componentes de la solución

**Frontend Web (`apps/web`)**
Aplicación SPA que sirve la interfaz de estudiantes y docentes, con editor de código integrado y consumo de la API por HTTP.

- Interfaz usuario
- Editor código
- Consumo API
- Visualización resultados
- Gestión sesión

**Backend API (`apps/api_v2`)**
Servicio HTTP que expone la lógica de negocio bajo arquitectura hexagonal y coordina autenticación, persistencia y envío de tareas.

- Autenticación JWT
- Lógica negocio
- Persistencia Roble
- Publicación submissions
- Documentación OpenAPI

**Worker (`apps/worker`)**
Proceso Python que consume mensajes de RabbitMQ y orquesta contenedores runner para ejecutar el código del estudiante.

- Consumo cola
- Orquestación runners
- Aislamiento ejecución
- Reporte resultados
- Salud runners

**Runners (`apps/worker/runners/{python,java,cpp}`)**
Contenedores efímeros por lenguaje donde se compila y ejecuta el código del usuario bajo límites estrictos de recursos.

- Compilación código
- Ejecución sandbox
- Casos prueba
- Límites recursos
- Captura salida

**RabbitMQ**
Broker que desacopla la API de los workers mediante colas por lenguaje (`python.queue`, `java.queue`, `cpp.queue`).

- Colas submissions
- Desacople servicios
- Tolerancia carga

**Roble (Base de datos)**
Servicio externo accedido vía API REST que provee la persistencia de usuarios, cursos, retos, exámenes y resultados.

- Persistencia entidades
- Consultas REST
- Almacenamiento resultados

**Asistente IA (`apps/ai`)**
Despliegue local de modelos LLM mediante Ollama para asistir al docente en la generación de enunciados y casos de prueba.

- Generación ejercicios
- Sugerencia casos
- Asistencia docente

## 2. Requisitos previos

### 2.1 Software requerido

#### Servidor (host que ejecuta la plataforma)

Para el flujo recomendado (con contenedores), todos los servicios se construyen y orquestan vía Docker Compose, por lo que el host únicamente necesita:

| Software                        | Versión mínima                  | Uso                                                                          |
| ------------------------------- | ------------------------------- | ---------------------------------------------------------------------------- |
| Docker Engine                   | 20.10+                          | Construcción y ejecución de las imágenes (`apiv2`, `worker`, `web`, runners) |
| Docker Compose                  | v2 (plugin)                     | Orquestación de los archivos `docker-compose.dev.yml` y `docker-compose.yml` |
| Git                             | 2.30+                           | Clonado del repositorio                                                      |
| Sistema operativo               | Linux / Windows 10+ / macOS 12+ | Soporte de Docker Desktop o Docker Engine nativo                             |
| Acceso a `/var/run/docker.sock` | —                               | Requerido por el `worker`, que crea contenedores runner dinámicamente        |
| Conectividad saliente           | —                               | Acceso HTTPS a `roble-api.openlab.uninorte.edu.co` (base de datos)           |
| Puertos disponibles             | —                               | `4000` (API), `5173` (web), `5672` y `15672` (RabbitMQ)                      |

Si se desea ejecutar los servicios **sin contenedores** (modo desarrollo nativo), adicionalmente se requiere:

| Software | Versión mínima  | Componente                                          |
| -------- | --------------- | --------------------------------------------------- |
| Go       | 1.25+           | Backend `api_v2` (`go run` / `air` para hot-reload) |
| Node.js  | 20.19+ / 22.12+ | Frontend `web` (requerido por Vite 7)               |
| npm      | 10+             | Instalación de dependencias del frontend            |
| Python   | 3.11+           | Worker (`apps/worker`)                              |
| RabbitMQ | 3.13            | Broker de mensajería                                |
| `air`    | última          | Hot-reload del backend Go                           |

#### Cliente (usuario final)

La plataforma se entrega como SPA, por lo que el cliente solo necesita:

| Software            | Versión mínima                             | Uso                                       |
| ------------------- | ------------------------------------------ | ----------------------------------------- |
| Navegador web       | Chrome/Edge 110+, Firefox 110+, Safari 16+ | Renderizado de la SPA y del editor Monaco |
| Conexión a internet | —                                          | Acceso al frontend y consumo de la API    |

### 2.2 Variables de entorno

La configuración se centraliza en archivos `.env` (`.env.dev` para desarrollo y `.env.production` para producción). A continuación se listan únicamente las variables actualmente consumidas por el código.

#### Frontend (`apps/web`)

| Variable       | Descripción                                                  | Ejemplo                 |
| -------------- | ------------------------------------------------------------ | ----------------------- |
| `VITE_API_URL` | URL base del backend que consume el cliente Axios            | `http://localhost:4000` |
| `WEB_PORT`     | Puerto en el que se expone el servidor de desarrollo de Vite | `5173`                  |

#### Backend API (`apps/api_v2`)

| Variable                   | Descripción                                                              | Ejemplo                                     |
| -------------------------- | ------------------------------------------------------------------------ | ------------------------------------------- |
| `APIV2_PORT`               | Puerto HTTP en el que escucha Fiber                                      | `4000`                                      |
| `APIV2_CORS_ALLOW_ORIGINS` | Orígenes permitidos por CORS                                             | `*`                                         |
| `ROBLE_PROJECT`            | Identificador del proyecto en Roble                                      | `coder_tests_6449a7e839`                    |
| `ROBLE_BASE_URL`           | URL base de la API de Roble                                              | `https://roble-api.openlab.uninorte.edu.co` |
| `INTERNAL_USER_EMAIL`      | Correo del usuario interno usado para escribir resultados de submissions | `admin@coder.com`                           |
| `INTERNAL_USER_PASSWORD`   | Contraseña del usuario interno (se hashea internamente)                  | `ChangeThisAdminPassword123!`               |
| `WORKER_KEY`               | Token compartido que autentica al worker al reportar resultados          | `CHANGE_THIS_WORKER_KEY_123`                |
| `SESSION_FREEZE_TIME`      | Segundos de congelamiento de una sesión tras un evento anti-cheat        | `60`                                        |
| `GEMINI_API_KEY`           | API key de Google Gemini para la generación de contenido                 | `AIza...`                                   |
| `RABBITMQ_URL`             | Cadena de conexión al broker de mensajería                               | `amqp://guest:guest@rabbitmq:5672/`         |

#### Worker (`apps/worker`)

| Variable                      | Descripción                                                          | Ejemplo                             |
| ----------------------------- | -------------------------------------------------------------------- | ----------------------------------- |
| `RABBITMQ_URL`                | Cadena de conexión al broker                                         | `amqp://guest:guest@rabbitmq:5672/` |
| `APIV2_BASE_URL`              | URL del backend a la que el worker reporta resultados                | `http://apiv2:4000`                 |
| `WORKER_KEY`                  | Token compartido con la API para autenticar el reporte de resultados | `CHANGE_THIS_WORKER_KEY_123`        |
| `WORKER_CONCURRENCY`          | Número máximo de tareas en paralelo dentro del worker                | `8`                                 |
| `PYTHON_QUEUE`                | Nombre de la cola de submissions de Python                           | `python.queue`                      |
| `PYTHON_RUNNERS`              | Cantidad de contenedores runner de Python a mantener vivos           | `3`                                 |
| `JAVA_RUNNERS`                | Cantidad de contenedores runner de Java a mantener vivos             | `3`                                 |
| `CPP_RUNNERS`                 | Cantidad de contenedores runner de C++ a mantener vivos              | `3`                                 |
| `PYTHON_RUNNER_IMAGE`         | Imagen Docker del runner de Python                                   | `coder-runner-python`               |
| `JAVA_RUNNER_IMAGE`           | Imagen Docker del runner de Java                                     | `coder-runner-java`                 |
| `CPP_RUNNER_IMAGE`            | Imagen Docker del runner de C++                                      | `coder-runner-cpp`                  |
| `RUNNER_NETWORK`              | Red Docker en la que se conectan los runners                         | `judge_net`                         |
| `RUNNER_NAME_PREFIX`          | Prefijo aplicado a los nombres de los contenedores runner            | `coder-runner`                      |
| `RUNNER_HEALTH_CHECK_TIMEOUT` | Timeout (s) del healthcheck contra un runner                         | `2.0`                               |
| `RUNNER_RESTART_TIMEOUT`      | Timeout (s) para reiniciar un runner caído                           | `120.0`                             |
| `RUNNER_REPAIR_INTERVAL`      | Intervalo (s) entre ciclos de reparación de la pool de runners       | `30.0`                              |

> **Nota:** Las variables relacionadas con la capa de caché (`CACHE_DB_DSN`, `CACHE_SYNC_QUEUE_NAME`, `CACHE_*_TABLE_TTL`), las colas alternativas `JAVA_QUEUE` y `CPP_QUEUE`, y las del bloque `DEPRECATED` (`POSTGRES_*`, `REDIS_*`, `HOST_SUBMISSIONS_DIR`, `API_PORT`) están declaradas en `.env.dev` pero no son consumidas por el código actual; se omiten hasta que sus módulos sean implementados.

## 3. Instalación para ambiente de desarrollo

### 3.1 Desarrollo sin contenedores

#### 3.1.1 Clonar el repositorio

Clonar el repositorio desde GitHub y posicionarse en la carpeta raíz del proyecto.

```bash
git clone https://github.com/openlabun/CODER.git
cd CODER
```

#### 3.1.2 Instalar dependencias

Cada componente del monorepo gestiona sus dependencias por separado.

**Backend (`apps/api_v2`)** — requiere Go 1.25+ y opcionalmente `air` para hot-reload:

```bash
cd apps/api_v2
go mod download
go install github.com/air-verse/air@latest   # opcional, solo si se quiere hot-reload
cd ../..
```

**Frontend (`apps/web`)** — requiere Node.js 20.19+ (o 22.12+):

```bash
cd apps/web
npm install
cd ../..
```

**Worker (`apps/worker`)** — requiere Python 3.11+:

```bash
cd apps/worker
python -m venv .venv
# Linux/macOS:
source .venv/bin/activate
# Windows (PowerShell):
.\.venv\Scripts\Activate.ps1
pip install -r requirements.txt
cd ../..
```

#### 3.1.3 Configurar variables de entorno

Crear el archivo `.env.dev` en la raíz del repositorio tomando como base las variables documentadas en la sección [2.2 Variables de entorno](#22-variables-de-entorno). Como mínimo, deben definirse:

```bash
# Backend
APIV2_PORT=4000
APIV2_CORS_ALLOW_ORIGINS=*
ROBLE_PROJECT=<tu_proyecto_en_roble>
ROBLE_BASE_URL=https://roble-api.openlab.uninorte.edu.co
INTERNAL_USER_EMAIL=admin@coder.com
INTERNAL_USER_PASSWORD=ChangeThisAdminPassword123!
WORKER_KEY=CHANGE_THIS_WORKER_KEY_123
SESSION_FREEZE_TIME=60
GEMINI_API_KEY=<tu_api_key>
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Worker
APIV2_BASE_URL=http://localhost:4000
PYTHON_QUEUE=python.queue
WORKER_CONCURRENCY=8
PYTHON_RUNNERS=3
JAVA_RUNNERS=3
CPP_RUNNERS=3
PYTHON_RUNNER_IMAGE=coder-runner-python
JAVA_RUNNER_IMAGE=coder-runner-java
CPP_RUNNER_IMAGE=coder-runner-cpp
RUNNER_NETWORK=judge_net
RUNNER_NAME_PREFIX=coder-runner
RUNNER_HEALTH_CHECK_TIMEOUT=2.0
RUNNER_RESTART_TIMEOUT=120.0
RUNNER_REPAIR_INTERVAL=30.0

# Frontend
VITE_API_URL=http://localhost:4000
WEB_PORT=5173
```

#### 3.1.4 Ejecutar servicios requeridos

El worker necesita un broker **RabbitMQ** corriendo y, además, las imágenes Docker de los runners para poder lanzarlos dinámicamente.

**RabbitMQ** (vía Docker, opción más sencilla):

```bash
docker run -d --name rabbitmq \
  -p 5672:5672 -p 15672:15672 \
  --network judge_net \
  rabbitmq:3.13-management
```

> Crear la red `judge_net` si no existe: `docker network create judge_net`.

**Imágenes de los runners** (construir una sola vez):

```bash
docker build -t coder-runner-python ./apps/worker/runners/python
docker build -t coder-runner-java   ./apps/worker/runners/java
docker build -t coder-runner-cpp    ./apps/worker/runners/cpp
```

#### 3.1.5 Iniciar la aplicación

Cada servicio se levanta en una terminal distinta, en el siguiente orden:

**1) Backend API** (`apps/api_v2`):

```bash
cd apps/api_v2
air -c .air.toml          # con hot-reload
# o, sin air:
go run ./cmd/app
```

**2) Worker** (`apps/worker`):

```bash
cd apps/worker
source .venv/bin/activate           # Linux/macOS
# .\.venv\Scripts\Activate.ps1      # Windows
python main.py
```

**3) Frontend Web** (`apps/web`):

```bash
cd apps/web
npm run dev -- --host
```

> Se usa `--host` para exponer el servidor de Vite en la red local (no solo en `localhost`) y facilitar pruebas desde otros dispositivos o desde servicios externos al proceso actual.

Una vez los tres servicios estén arriba:

- Frontend: `http://localhost:5173`
- API: `http://localhost:4000`
- Documentación API (Scalar): `http://localhost:4000/docs`
- Consola de RabbitMQ: `http://localhost:15672` (usuario/contraseña `guest`/`guest`)

### 3.2 Desarrollo con contenedores

#### 3.2.1 Construcción de contenedores

Construye las imágenes de `apiv2`, `worker` y `web` definidas en `docker-compose.dev.yml`.

```bash
docker compose -f docker-compose.dev.yml build
```

Adicionalmente, construir las imágenes de los runners (no se construyen automáticamente porque pertenecen al perfil `runners`):

```bash
docker compose -f docker-compose.dev.yml --profile runners build
```

#### 3.2.2 Ejecución del entorno

Levanta todos los servicios en segundo plano usando las variables de `.env.dev`.

```bash
docker compose -f docker-compose.dev.yml --env-file .env.dev up -d
```

Para seguir los logs en tiempo real:

```bash
docker compose -f docker-compose.dev.yml logs -f
```

#### 3.2.3 Servicios disponibles

Verifica que los contenedores estén corriendo:

```bash
docker compose -f docker-compose.dev.yml ps
```

Endpoints expuestos:

| Servicio        | URL                                    |
| --------------- | -------------------------------------- |
| Frontend Web    | `http://localhost:5173`                |
| Backend API     | `http://localhost:4000`                |
| Documentación   | `http://localhost:4000/docs`           |
| RabbitMQ (AMQP) | `localhost:5672`                       |
| RabbitMQ (UI)   | `http://localhost:15672` (guest/guest) |

#### 3.2.4 Apagado del entorno

Detiene y elimina los contenedores, manteniendo las imágenes construidas.

```bash
docker compose -f docker-compose.dev.yml down
```

Para eliminar también volúmenes y la red `judge_net`:

```bash
docker compose -f docker-compose.dev.yml down -v --remove-orphans
```

## 4. Despliegue

### 4.1 Arquitectura de despliegue

El despliegue se organiza en tres planos: el **cliente** (navegador del usuario), el **servidor** (host con Docker que ejecuta la API, el broker, el worker y los runners) y los **servicios externos** consumidos por la API (Roble y Gemini).

```mermaid
flowchart LR
    Cliente["Cliente<br/>(Navegador / SPA)"]

    subgraph Servidor["Servidor (Docker host)"]
        API["Backend API<br/>(Fiber - Go)"]
        MQ["RabbitMQ<br/>(broker)"]
        Worker["Worker<br/>(Python)"]
        subgraph Runners["Runners"]
            RP["Python"]
            RJ["Java"]
            RC["C++"]
        end
    end

    subgraph Externo["Servicios externos"]
        Roble[("Roble<br/>(Base de datos)")]
        Gemini["Gemini API<br/>(LLM)"]
    end

    Cliente -- HTTPS --> API
    API -- AMQP --> MQ
    MQ -- AMQP --> Worker
    Worker -- "docker run" --> Runners
    Worker -- "POST resultados" --> API
    API -- REST --> Roble
    API -- HTTPS --> Gemini
```

**Gráfico:** Arquitectura de despliegue de CODER

### 4.2 Proceso de actualización

El proceso de actualización del entorno se ejecuta mediante un pipeline de **CI/CD** definido en GitHub Actions ([`.github/workflows/main.yml`](.github/workflows/main.yml)). El flujo se dispara automáticamente con cada `push` a la rama `release` y aplica los cambios sobre el servidor de producción a través de una sesión SSH.

```yaml
name: Deploy

on:
    push:
        branches: [ release ]

jobs:
    build:
        runs-on: ubuntu-latest
        steps:
            - name: Deploy using ssh
              uses: appleboy/ssh-action@master
              with:
                  host: ${{ secrets.HOST }}
                  username: ${{ secrets.USERNAME }}
                  key: ${{ secrets.SSH_PRIVATE_KEY }}
                  port: 22
                  script: |
                      cd /home/proyecto/CODER/
                      git pull origin release
                      git status
                      COMPOSE_BAKE=true docker compose --env-file .env.production build
                      docker compose down
                      docker compose -f docker-compose.yml --profile judge build runner-python runner-cpp runner-java
                      docker compose -f docker-compose.yml --env-file .env.production up -d
```

El workflow realiza las siguientes acciones sobre el servidor:

1. Se conecta por SSH usando los secretos `HOST`, `USERNAME` y `SSH_PRIVATE_KEY` configurados en el repositorio.
2. Se posiciona en el directorio del proyecto (`/home/proyecto/CODER/`) y trae los últimos cambios con `git pull origin release`.
3. Reconstruye las imágenes principales (`apiv2`, `worker`, `web`) usando las variables de `.env.production` y la optimización `COMPOSE_BAKE=true`.
4. Detiene los contenedores existentes (`docker compose down`) y reconstruye las imágenes de los runners (`runner-python`, `runner-cpp`, `runner-java`) bajo el perfil `judge`.
5. Levanta el stack actualizado en segundo plano con `docker compose up -d`.

> **Sobre `COMPOSE_BAKE=true`:** habilita Docker Buildx Bake en `docker compose build`, permitiendo optimizar la construcción (incluido paralelismo) y reducir tiempos de despliegue en pipelines CI/CD.

> **Nota:** Actualmente el pipeline contempla un único entorno de despliegue (producción). Para mejorar la gestión del ciclo de vida del software sería conveniente disponer de **dos despliegues separados**, uno de **desarrollo** (por ejemplo, disparado por `push` a `main` o a una rama `develop`) y otro de **producción** (disparado por `push` a `release`), de modo que los cambios puedan validarse en un entorno equivalente antes de llegar a los usuarios finales. A favor de esta separación, el proyecto ya cuenta con **bases de datos independientes** para cada entorno en Roble (`coder_tests_*` para desarrollo y un proyecto distinto para producción), por lo que únicamente faltaría replicar el job de despliegue con su propio `.env`, sus propios secretos SSH y su propio host de destino.

### 4.3 Despliegue sin contenedores

Sigue el mismo flujo que el descrito en [3.1 Desarrollo sin contenedores](#31-desarrollo-sin-contenedores), pero ejecutado sobre un servidor Linux con un usuario de servicio y con un gestor de procesos (`systemd`) que mantenga los servicios vivos tras reinicios.

#### 4.3.1 Preparación del servidor

Servidor Linux (Ubuntu 22.04+ recomendado) con acceso SSH y un usuario dedicado para el despliegue.

```bash
sudo adduser --system --group coder
sudo mkdir -p /home/proyecto/CODER && sudo chown -R coder:coder /home/proyecto/CODER
sudo apt update && sudo apt install -y git curl build-essential
```

Instalar los runtimes requeridos en la sección [2.1](#21-software-requerido): Go 1.25+, Node.js 20.19+, Python 3.11+ y RabbitMQ 3.13.

```bash
# RabbitMQ
sudo apt install -y rabbitmq-server
sudo systemctl enable --now rabbitmq-server
```

#### 4.3.2 Instalación de dependencias

Clonado del repositorio y resolución de dependencias por componente (idéntico a [3.1.1](#311-clonar-el-repositorio) y [3.1.2](#312-instalar-dependencias)).

```bash
sudo -iu coder
cd /home/proyecto/CODER
git clone https://github.com/openlabun/CODER.git .
git checkout release

# Backend
(cd apps/api_v2 && go mod download && go build -o /home/proyecto/CODER/bin/apiv2 ./cmd/app)

# Frontend (build estático para producción)
(cd apps/web && npm ci && npm run build)

# Worker
(cd apps/worker && python -m venv .venv && .venv/bin/pip install -r requirements.txt)
```

#### 4.3.3 Configuración de la aplicación

Crear el archivo `/home/proyecto/CODER/.env.production` con las variables documentadas en [2.2](#22-variables-de-entorno), apuntando al proyecto de Roble de producción y al RabbitMQ local. Restringir permisos al usuario de servicio:

```bash
sudo chown coder:coder /home/proyecto/CODER/.env.production
sudo chmod 600 /home/proyecto/CODER/.env.production
```

#### 4.3.4 Ejecución de la aplicación

Cada servicio se administra como una unidad de `systemd` para garantizar reinicio automático y trazas centralizadas.

```ini
# /etc/systemd/system/coder-apiv2.service
[Unit]
Description=CODER API v2
After=network.target rabbitmq-server.service

[Service]
User=coder
WorkingDirectory=/home/proyecto/CODER
EnvironmentFile=/home/proyecto/CODER/.env.production
ExecStart=/home/proyecto/CODER/bin/apiv2
Restart=always

[Install]
WantedBy=multi-user.target
```

Crear unidades análogas `coder-worker.service` (ejecuta `apps/worker/.venv/bin/python apps/worker/main.py`) y `coder-web.service` (sirve `apps/web/dist` con `npx serve` o un Nginx).

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now coder-apiv2 coder-worker coder-web
```

#### 4.3.5 Actualización de versiones

Traer los cambios desde la rama `release`, recompilar artefactos y reiniciar las unidades afectadas.

```bash
cd /home/proyecto/CODER
git pull origin release
(cd apps/api_v2 && go build -o /home/proyecto/CODER/bin/apiv2 ./cmd/app)
(cd apps/web && npm ci && npm run build)
(cd apps/worker && .venv/bin/pip install -r requirements.txt)

sudo systemctl restart coder-apiv2 coder-worker coder-web
```

### 4.4 Despliegue con contenedores

Reutiliza el flujo de [3.2 Desarrollo con contenedores](#32-desarrollo-con-contenedores), pero usando `docker-compose.yml` y `.env.production`. Es el camino que sigue el pipeline de CI/CD descrito en [4.2](#42-proceso-de-actualización).

#### 4.4.1 Construcción de imágenes

Construye las imágenes principales (`apiv2`, `worker`, `web`) y, bajo el perfil `judge`, las imágenes de los runners.

```bash
cd /home/proyecto/CODER
COMPOSE_BAKE=true docker compose -f docker-compose.yml --env-file .env.production build
docker compose -f docker-compose.yml --profile judge build runner-python runner-cpp runner-java
```

#### 4.4.2 Ejecución en servidor

Levanta el stack en segundo plano con la política de reinicio definida en `docker-compose.yml` (`restart: always`).

```bash
docker compose -f docker-compose.yml --env-file .env.production up -d
docker compose -f docker-compose.yml ps
```

#### 4.4.3 Variables de entorno y secretos

El archivo `.env.production` (ver [2.2](#22-variables-de-entorno)) debe permanecer fuera del repositorio, con permisos `600` y propiedad del usuario que ejecuta Docker.

```bash
sudo chown root:root /home/proyecto/CODER/.env.production
sudo chmod 600 /home/proyecto/CODER/.env.production
```

Los secretos consumidos por el pipeline (`HOST`, `USERNAME`, `SSH_PRIVATE_KEY`) se configuran en *Settings → Secrets and variables → Actions* del repositorio en GitHub.

#### 4.4.4 Persistencia y redes

- Todos los servicios comparten la red bridge `judge_net`, declarada en el `docker-compose.yml`.
- El worker monta `/var/run/docker.sock` para poder lanzar contenedores runner dinámicamente.
- La persistencia de datos de negocio reside en **Roble** (externo), por lo que no se requieren volúmenes adicionales para la API. Solo RabbitMQ podría requerir un volumen si se desea conservar el estado de las colas entre reinicios.

#### 4.4.5 Actualización del despliegue

Equivale a re-ejecutar los pasos [4.4.1](#441-construcción-de-imágenes) y [4.4.2](#442-ejecución-en-servidor) tras un `git pull`. Es exactamente el script que ejecuta el workflow de GitHub Actions:

```bash
cd /home/proyecto/CODER/
git pull origin release
COMPOSE_BAKE=true docker compose --env-file .env.production build
docker compose down
docker compose -f docker-compose.yml --profile judge build runner-python runner-cpp runner-java
docker compose -f docker-compose.yml --env-file .env.production up -d
```

## 5. Verificación de funcionamiento

Una vez desplegado, se recomienda validar el sistema en el siguiente orden, de abajo hacia arriba en la cadena de dependencias.

1. **Estado de los contenedores.** Verificar que todos los servicios estén en estado `Up` y sin reinicios recientes.
   
   ```bash
   docker compose -f docker-compose.yml ps
   ```

2. **Broker RabbitMQ.** Acceder a la consola de administración (`http://<servidor>:15672`, credenciales `guest/guest` por defecto) y comprobar que existan las colas `python.queue`, `java.queue` y `cpp.queue` cuando se envíe la primera submission. Alternativamente:
   
   ```bash
   docker exec -it $(docker ps -qf name=rabbitmq) rabbitmqctl list_queues name messages consumers
   ```
   
   Cada cola debe mostrar `consumers ≥ 1` (el worker conectado).

3. **API Backend.** Comprobar que responde en su puerto y que la documentación OpenAPI carga:
   
   ```bash
   curl -i http://<servidor>:4000/docs/openapi.yaml
   ```
   
   Abrir `http://<servidor>:4000/docs` en el navegador para ver la UI de Scalar.

4. **Worker y runners.** Tras enviar una submission desde el frontend, validar que el worker recibió el mensaje y que se creó un contenedor runner efímero:
   
   ```bash
   docker compose logs -f worker
   docker ps --filter "name=coder-runner"
   ```

5. **Frontend.** Abrir `http://<servidor>:5173` y autenticarse. La consola del navegador (`F12`) no debería mostrar errores de CORS ni de red contra `VITE_API_URL`.

6. **Flujo end-to-end.** Crear un reto desde la cuenta docente, enviar una solución desde una cuenta estudiante y verificar que la calificación queda registrada y visible en el panel de resultados. Este es el indicador definitivo de que toda la cadena (API → RabbitMQ → Worker → Runner → API → Roble) opera correctamente.

## 6. Solución de problemas frecuentes

1. **La API no arranca y termina con error de conexión a Roble.**
   Suele deberse a variables `ROBLE_PROJECT` o `ROBLE_BASE_URL` ausentes/incorrectas en `.env.*`. Comprobar con `docker compose logs apiv2` y validar que `roble-api.openlab.uninorte.edu.co` sea alcanzable desde el host (`curl -I https://roble-api.openlab.uninorte.edu.co`).

2. **El worker no consume mensajes de la cola.**
   Generalmente RabbitMQ no terminó de arrancar antes que el worker, o `RABBITMQ_URL` apunta al host equivocado. Verificar `docker compose logs rabbitmq` y reiniciar el worker (`docker compose restart worker`). Confirmar en la consola que la cola tiene `consumers ≥ 1`.

3. **Las submissions quedan en cola pero no se ejecutan.**
   Implica que el worker no logró lanzar un runner. Causas habituales:
   
   - Las imágenes `coder-runner-python`, `coder-runner-java`, `coder-runner-cpp` no fueron construidas (correr `docker compose --profile judge build`).
   - El worker no tiene acceso al socket Docker; revisar que `/var/run/docker.sock` esté montado.
   - El servicio worker no comparte la red `judge_net` con los runners; revisar `RUNNER_NETWORK`.

4. **El backend rechaza el reporte de resultados del worker con Error HTTP 401/403.**
   `WORKER_KEY` debe ser idéntico en el `.env` que carga la API y en el que carga el worker. Cualquier divergencia bloquea el endpoint de resultados.

5. **El frontend no se comunica con la API (errores de CORS o `net::ERR_CONNECTION_REFUSED`).**
   
   - `VITE_API_URL` no fue inyectado en build time del contenedor `web` (recordar que Vite hornea variables al hacer `npm run build`).
   - `APIV2_CORS_ALLOW_ORIGINS` no incluye el origen del frontend (en producción, evitar `*` y listar el dominio real).

6. **Las funcionalidades de IA fallan con error 401/429.**
   `GEMINI_API_KEY` inválida, vencida o con cuota agotada. Validar en Google AI Studio y rotar la clave si fue expuesta.

7. **Conflicto de puertos al levantar el stack.**
   Otro proceso ocupa `4000`, `5173`, `5672` o `15672`. Identificar con `ss -tlnp | grep <puerto>` (Linux) o `Get-NetTCPConnection -LocalPort <puerto>` (Windows) y liberar el puerto o cambiar el mapeo en `docker-compose.yml`.

8. **El despliegue automático falla en el paso SSH del workflow.**
   Revisar los secretos `HOST`, `USERNAME` y `SSH_PRIVATE_KEY` en GitHub y que el directorio `/home/proyecto/CODER/` exista en el servidor con permisos para el usuario configurado.

9. **Sesiones de examen se congelan inesperadamente o no se reactivan.**
   `SESSION_FREEZE_TIME` está en segundos; un valor demasiado bajo provoca bloqueos seguidos. Ajustar y reiniciar la API.

10. **Tras un `docker compose down -v`, no se reciben mensajes.**
    Al eliminar volúmenes se pierde el estado de RabbitMQ; cualquier mensaje en cola desaparece. No usar `-v` en producción salvo mantenimiento explícito.

## 7. Mantenimiento y actualización

**Actualización de código.** Es el flujo descrito en [4.4.5](#445-actualización-del-despliegue): `git pull` sobre la rama `release`, rebuild de imágenes y `up -d`. En producción este paso lo dispara el workflow de GitHub Actions cada vez que se mergea hacia `release`.

**Actualización de imágenes base.** Mensualmente reconstruir las imágenes para recoger parches de seguridad de `rabbitmq:3.13-management`, las imágenes de los runners y las bases de Go/Node/Python usadas en los `Dockerfile.prod`.

```bash
docker compose -f docker-compose.yml --env-file .env.production build --pull --no-cache
docker compose -f docker-compose.yml --env-file .env.production up -d
```

**Monitoreo y logs.** Centralizar la inspección con:

```bash
docker compose -f docker-compose.yml logs -f --tail=200 apiv2 worker
docker stats
```

Para un seguimiento más fino, se recomienda integrar un colector (Loki/Promtail, ELK o similar) que recoja los logs de los contenedores.

**Higiene del host Docker.** Los runners son efímeros y, ante fallos, pueden quedar contenedores e imágenes huérfanas. Programar una limpieza semanal:

```bash
docker container prune -f
docker image prune -f
docker network prune -f
```

**Respaldo de configuración.** Los archivos `.env.production`, las unidades de `systemd` (si se usa el modo sin contenedores) y los secretos de GitHub son los activos críticos a respaldar. Roble gestiona los datos de negocio, por lo que el respaldo de la base recae en ese servicio externo.

**Rotación de credenciales.** Rotar periódicamente `WORKER_KEY`, `INTERNAL_USER_PASSWORD` y `GEMINI_API_KEY`. Tras cada rotación, redeplegar el stack para que API y worker queden con valores sincronizados.

**Verificación post-actualización.** Tras cada despliegue, repetir el checklist de la sección [5. Verificación de funcionamiento](#5-verificación-de-funcionamiento), prestando especial atención al flujo end-to-end de envío y calificación.

## 8. Referencias relacionadas

- [1] **`air` — Live reload para aplicaciones Go.** Repositorio y documentación oficial: <https://github.com/air-verse/air>
- [2] **Roble — Documentación oficial.** Universidad del Norte, Openlab: <https://roble.openlab.uninorte.edu.co/docs>
- [3] **GitHub Actions — Documentación oficial.** GitHub: <https://docs.github.com/actions>
- [4] **Uso de secretos en GitHub Actions.** GitHub: <https://docs.github.com/actions/security-guides/using-secrets-in-github-actions>
