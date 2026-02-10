#!/usr/bin/env bash
# Ensure the script file has executable permissions so it can be run directly.
# If the file is executed via `sh create_and_enqueue.sh` this block will still
# try to make the file executable for subsequent runs.
if [ -w "$0" ]; then
  chmod +x "$0" || true
fi
# create_and_enqueue.sh
#
# Crea un challenge y una submission de prueba en Postgres, y empuja la submission
# a la cola Redis (key: queue:submissions). Funciona tanto con clientes locales
# (`psql`, `redis-cli`) como con contenedores Docker (usando `docker exec`).
#
# Uso:
#   ./create_and_enqueue.sh [--challenge-id ID] [--submission-id ID] [--user-id USER]
#                           [--pg-host HOST] [--pg-port PORT] [--pg-user USER]
#                           [--pg-db DB] [--pg-container NAME]
#                           [--redis-host HOST] [--redis-port PORT] [--redis-container NAME]
#
# Ejemplos:
#   # usar psql/redis-cli locales (valores por defecto)
#   ./create_and_enqueue.sh
#
#   # usar contenedores docker
#   ./create_and_enqueue.sh --pg-container postgres --redis-container redis
#
# Variables por defecto:
#   PG_HOST=127.0.0.1
#   PG_PORT=5432
#   PG_USER=postgres
#   PG_DB=postgres
#   REDIS_HOST=127.0.0.1
#   REDIS_PORT=6379
#
# Nota:
# - El script asumirá que la tabla `public.challenges` y `public.submissions`
#   existen (las migraciones están en migrations/001_create_challenges.sql y 002).
# - El `status` de la submission se insertará como `queued` (el worker espera esto).
#
set -euo pipefail

# ---- utilidades ----
err() { echo "$@" >&2; }
usage() {
  cat <<EOF
Usage: $0 [options]

Options:
  --challenge-id ID       UUID para el challenge a crear (si ya existe, se actualizará)
  --submission-id ID      UUID para la submission (si ya existe, se actualizará)
  --user-id USER          user_id para la submission (default: test-user)
  --pg-host HOST          Postgres host (default: 127.0.0.1)
  --pg-port PORT          Postgres port (default: 5432)
  --pg-user USER          Postgres user (default: postgres)
  --pg-db DB              Postgres database (default: postgres)
  --pg-container NAME     Si se proporciona, usa "docker exec -i NAME psql ..." en vez de psql local
  --redis-host HOST       Redis host (default: 127.0.0.1)
  --redis-port PORT       Redis port (default: 6379)
  --redis-container NAME  Si se proporciona, usa "docker exec -i NAME redis-cli ..." en vez de redis-cli local
  --help                  Mostrar este mensaje
EOF
  exit 1
}

# ---- defaults ----
CHALLENGE_ID=""
SUBMISSION_ID=""
USER_ID="test-user"

PG_HOST="127.0.0.1"
PG_PORT="5432"
PG_USER="postgres"
PG_DB="postgres"
PG_CONTAINER=""

REDIS_HOST="127.0.0.1"
REDIS_PORT="6379"
REDIS_CONTAINER=""

# ---- parse args ----
while [[ $# -gt 0 ]]; do
  case "$1" in
    --challenge-id) CHALLENGE_ID="$2"; shift 2;;
    --submission-id) SUBMISSION_ID="$2"; shift 2;;
    --user-id) USER_ID="$2"; shift 2;;
    --pg-host) PG_HOST="$2"; shift 2;;
    --pg-port) PG_PORT="$2"; shift 2;;
    --pg-user) PG_USER="$2"; shift 2;;
    --pg-db) PG_DB="$2"; shift 2;;
    --pg-container) PG_CONTAINER="$2"; shift 2;;
    --redis-host) REDIS_HOST="$2"; shift 2;;
    --redis-port) REDIS_PORT="$2"; shift 2;;
    --redis-container) REDIS_CONTAINER="$2"; shift 2;;
    --help) usage;;
    *) err "Unknown arg: $1"; usage;;
  esac
done

# ---- helpers para generar UUID ----
gen_uuid() {
  if command -v uuidgen >/dev/null 2>&1; then
    uuidgen
  elif command -v openssl >/dev/null 2>&1; then
    # openssl rand -hex 16 -> 32 hex chars; format like 8-4-4-4-12
    local hex
    hex=$(openssl rand -hex 16)
    echo "${hex:0:8}-${hex:8:4}-${hex:12:4}-${hex:16:4}-${hex:20:12}"
  else
    # fallback simple (pseudo-uuid)
    date +%s%N | sha1sum | awk '{print substr($1,1,8) "-" substr($1,9,4) "-" substr($1,13,4) "-" substr($1,17,4) "-" substr($1,21,12)}'
  fi
}

# ---- asegurar IDs ----
if [[ -z "$CHALLENGE_ID" ]]; then
  CHALLENGE_ID=$(gen_uuid)
fi
if [[ -z "$SUBMISSION_ID" ]]; then
  SUBMISSION_ID=$(gen_uuid)
fi

echo "[info] challenge_id: $CHALLENGE_ID"
echo "[info] submission_id: $SUBMISSION_ID"
echo "[info] user_id: $USER_ID"

# ---- construir SQL ----
# Usamos INSERT ... ON CONFLICT para ser idempotentes: si existe, actualizamos.
SQL=$(cat <<SQL
BEGIN;

-- Crear/actualizar challenge de prueba (status 'published')
INSERT INTO public.challenges (id, title, description, status, created_at, updated_at)
VALUES ('$CHALLENGE_ID', 'Reto de prueba (auto)', 'Generado por create_and_enqueue.sh', 'published', NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
  title = EXCLUDED.title,
  description = EXCLUDED.description,
  status = EXCLUDED.status,
  updated_at = NOW();

-- Crear/actualizar submission en estado queued
INSERT INTO public.submissions (id, challenge_id, user_id, status, created_at, updated_at)
VALUES (
  '$SUBMISSION_ID',
  '$CHALLENGE_ID',
  '$USER_ID',
  'queued',
  NOW(),
  NOW()
)
ON CONFLICT (id) DO UPDATE SET
  challenge_id = EXCLUDED.challenge_id,
  user_id = EXCLUDED.user_id,
  status = EXCLUDED.status,
  updated_at = NOW();

COMMIT;
SQL
)

# ---- ejecutar SQL en Postgres ----
echo "[info] Ejecutando SQL en Postgres..."

if [[ -n "$PG_CONTAINER" ]]; then
  # Usar docker exec - permitimos que la imagen tenga psql disponible.
  echo "$SQL" | docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB"
else
  # Usar psql local. Requerimientos: psql en PATH y accesible con las credenciales.
  if ! command -v psql >/dev/null 2>&1; then
    err "psql no encontrado en PATH. Instala psql o usa --pg-container para ejecutar dentro de un contenedor."
    exit 2
  fi
  PGPASSWORD="${PG_PASSWORD:-}" psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d "$PG_DB" -v ON_ERROR_STOP=1 -q -c "$SQL"
fi

echo "[ok] SQL ejecutado."

# ---- comprobar que la submission está en la tabla ----
echo "[info] Verificando submission en la base de datos..."
CHECK_SQL="SELECT id, challenge_id, user_id, status, created_at FROM public.submissions WHERE id = '$SUBMISSION_ID';"

if [[ -n "$PG_CONTAINER" ]]; then
  docker exec -i "$PG_CONTAINER" psql -U "$PG_USER" -d "$PG_DB" -c "$CHECK_SQL"
else
  psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d "$PG_DB" -c "$CHECK_SQL"
fi

# ---- empujar ID a Redis (LPUSH queue:submissions $SUBMISSION_ID) ----
echo "[info] Encolando submission en Redis (key: queue:submissions)..."

if [[ -n "$REDIS_CONTAINER" ]]; then
  docker exec -i "$REDIS_CONTAINER" redis-cli LPUSH queue:submissions "$SUBMISSION_ID"
else
  if ! command -v redis-cli >/dev/null 2>&1; then
    err "redis-cli no encontrado en PATH. Instala redis-tools o usa --redis-container para ejecutar en un contenedor."
    exit 3
  fi
  redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" LPUSH queue:submissions "$SUBMISSION_ID"
fi

echo "[ok] Encolado: $SUBMISSION_ID"

# ---- mostrar estado de la cola y fin ----
if [[ -n "$REDIS_CONTAINER" ]]; then
  echo "[info] LLEN queue:submissions:"
  docker exec -i "$REDIS_CONTAINER" redis-cli LLEN queue:submissions
  echo "[info] LRANGE 0 -5 (últimos 6 elementos):"
  docker exec -i "$REDIS_CONTAINER" redis-cli LRANGE queue:submissions 0 5
else
  echo "[info] LLEN queue:submissions:"
  redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" LLEN queue:submissions
  echo "[info] LRANGE 0 -5 (últimos 6 elementos):"
  redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" LRANGE queue:submissions 0 5
fi

echo
echo "Hecho. El worker debería detectar el job y procesarlo."
echo "Submission ID: $SUBMISSION_ID"
echo "Challenge ID: $CHALLENGE_ID"
