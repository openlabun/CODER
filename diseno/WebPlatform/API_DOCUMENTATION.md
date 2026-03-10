# CODER API — Documentación Técnica

## Índice

1. [Resumen General](#resumen-general)
2. [Endpoints](#endpoints)
3. [Contratos (Inputs / Outputs)](#contratos-inputs--outputs)
4. [Funcionalidades Implementadas](#funcionalidades-implementadas)
   - [Autenticación](#1-autenticación-auth)
   - [Challenges (Desafíos)](#2-challenges-desafíos)
   - [Test Cases (Casos de Prueba)](#3-test-cases-casos-de-prueba)
   - [Submissions (Entregas)](#4-submissions-entregas)
   - [Courses (Cursos)](#5-courses-cursos)
   - [Exams (Exámenes)](#6-exams-exámenes)
   - [Leaderboard (Tabla de posiciones)](#7-leaderboard-tabla-de-posiciones)
   - [AI (Generación con IA)](#8-ai-generación-con-ia)
   - [Metrics (Métricas)](#9-metrics-métricas)
   - [Health Checks](#10-health-checks)

---

## Resumen General

API backend construida con **NestJS** (TypeScript) que implementa un **juez online** para programación. Soporta ejecución de código en **Python, Node.js, C++ y Java** dentro de contenedores Docker aislados. Utiliza **PostgreSQL** como base de datos, **Redis** como cola de trabajos, y **Google Gemini** para generación de contenido con IA.

**Stack tecnológico:**

- Framework: NestJS
- Base de datos: PostgreSQL (pg)
- Cola: Redis (ioredis)
- Autenticación: JWT + bcrypt
- Ejecución: Docker containers aislados
- IA: Google Generative AI (Gemini Flash)

---

## Endpoints

### Autenticación

| Método | Ruta             | Auth | Rol        | Descripción            |
| ------ | ---------------- |:----:| ---------- | ---------------------- |
| `POST` | `/auth/register` | ❌    | —          | Registrar usuario      |
| `POST` | `/auth/login`    | ❌    | —          | Iniciar sesión (JWT)   |
| `GET`  | `/auth/me`       | ✅    | cualquiera | Obtener usuario actual |

### Challenges

| Método  | Ruta                      | Auth | Rol            | Descripción                    |
| ------- | ------------------------- |:----:| -------------- | ------------------------------ |
| `POST`  | `/challenges`             | ✅    | profesor/admin | Crear challenge con test cases |
| `GET`   | `/challenges`             | ✅    | cualquiera     | Listar challenges públicos     |
| `GET`   | `/challenges/:id`         | ✅    | cualquiera     | Obtener detalle de challenge   |
| `PATCH` | `/challenges/:id`         | ✅    | profesor/admin | Actualizar challenge           |
| `POST`  | `/challenges/:id/publish` | ✅    | profesor/admin | Publicar challenge             |
| `POST`  | `/challenges/:id/archive` | ✅    | profesor/admin | Archivar challenge             |

### Test Cases

| Método   | Ruta                                 | Auth | Rol            | Descripción                                  |
| -------- | ------------------------------------ |:----:| -------------- | -------------------------------------------- |
| `POST`   | `/test-cases`                        | ✅    | profesor/admin | Crear caso de prueba                         |
| `GET`    | `/test-cases/challenge/:challengeId` | ✅    | cualquiera     | Listar test cases (alumnos solo ven samples) |
| `DELETE` | `/test-cases/:id`                    | ✅    | profesor/admin | Eliminar caso de prueba                      |

### Submissions

| Método | Ruta               | Auth | Rol        | Descripción                    |
| ------ | ------------------ |:----:| ---------- | ------------------------------ |
| `POST` | `/submissions`     | ✅    | estudiante | Enviar código                  |
| `GET`  | `/submissions/:id` | ✅    | cualquiera | Obtener detalle de submission  |
| `GET`  | `/submissions`     | ✅    | cualquiera | Listar submissions del usuario |

### Courses

| Método   | Ruta                               | Auth | Rol            | Descripción                  |
| -------- | ---------------------------------- |:----:| -------------- | ---------------------------- |
| `POST`   | `/courses`                         | ✅    | profesor/admin | Crear curso                  |
| `GET`    | `/courses`                         | ✅    | cualquiera     | Listar cursos propios        |
| `GET`    | `/courses/browse`                  | ✅    | cualquiera     | Navegar todos los cursos     |
| `GET`    | `/courses/:id`                     | ✅    | cualquiera     | Detalle de curso             |
| `POST`   | `/courses/:id`                     | ✅    | profesor       | Actualizar curso             |
| `POST`   | `/courses/enroll`                  | ✅    | estudiante     | Inscribirse con código       |
| `POST`   | `/courses/:id/students`            | ✅    | profesor/admin | Agregar estudiante           |
| `DELETE` | `/courses/:id/students/:studentId` | ✅    | profesor/admin | Remover estudiante           |
| `POST`   | `/courses/:id/challenges`          | ✅    | profesor/admin | Asignar challenge a curso    |
| `GET`    | `/courses/:id/students`            | ✅    | cualquiera     | Listar estudiantes del curso |
| `GET`    | `/courses/:id/challenges`          | ✅    | cualquiera     | Listar challenges del curso  |

### Exams

| Método | Ruta                      | Auth | Rol            | Descripción                      |
| ------ | ------------------------- |:----:| -------------- | -------------------------------- |
| `POST` | `/exams`                  | ✅    | profesor/admin | Crear examen                     |
| `GET`  | `/exams/course/:courseId` | ✅    | cualquiera     | Listar exámenes de un curso      |
| `GET`  | `/exams/:id`              | ✅    | cualquiera     | Detalle de examen con challenges |

### Leaderboard

| Método | Ruta                         | Auth | Rol        | Descripción           |
| ------ | ---------------------------- |:----:| ---------- | --------------------- |
| `GET`  | `/leaderboard/challenge/:id` | ✅    | cualquiera | Ranking por challenge |
| `GET`  | `/leaderboard/course/:id`    | ✅    | cualquiera | Ranking por curso     |

### AI

| Método | Ruta                           | Auth | Rol            | Descripción                        |
| ------ | ------------------------------ |:----:| -------------- | ---------------------------------- |
| `POST` | `/ai/generate-challenge-ideas` | ✅    | profesor/admin | Generar ideas de challenges con IA |
| `POST` | `/ai/generate-test-cases`      | ✅    | profesor/admin | Generar test cases con IA          |

### Metrics & Health

| Método | Ruta            | Auth | Rol | Descripción                              |
| ------ | --------------- |:----:| --- | ---------------------------------------- |
| `GET`  | `/metrics`      | ❌    | —   | Métricas del sistema (JSON + Prometheus) |
| `GET`  | `/health`       | ❌    | —   | Health check general                     |
| `GET`  | `/cache/health` | ❌    | —   | Health check Redis                       |
| `GET`  | `/db/health`    | ❌    | —   | Health check PostgreSQL                  |

---

## Contratos (Inputs / Outputs)

### Auth

#### `POST /auth/register`

**Input:**

```json
{
  "username": "string",
  "password": "string",
  "role": "student | professor | admin"
}
```

**Output:**

```json
{
  "id": "uuid",
  "username": "string",
  "role": "string",
  "createdAt": "ISO 8601"
}
```

#### `POST /auth/login`

**Input:**

```json
{
  "username": "string",
  "password": "string"
}
```

**Output:**

```json
{
  "access_token": "string (JWT)"
}
```

#### `GET /auth/me`

**Headers:** `Authorization: Bearer <JWT>`

**Output:**

```json
{
  "id": "uuid",
  "username": "string",
  "role": "student | professor | admin"
}
```

---

### Challenges

#### `POST /challenges`

**Input:**

```json
{
  "title": "string",
  "description": "string",
  "difficulty": "easy | medium | hard (opcional)",
  "timeLimit": "number (ms, opcional)",
  "memoryLimit": "number (MB, opcional)",
  "tags": ["string"] ,
  "inputFormat": "string (opcional)",
  "outputFormat": "string (opcional)",
  "constraints": "string (opcional)",
  "publicTestCases": [
    { "name": "string", "input": "string", "expectedOutput": "string", "points": "number" }
  ],
  "hiddenTestCases": [
    { "name": "string", "input": "string", "expectedOutput": "string", "points": "number" }
  ]
}
```

**Output:**

```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "status": "draft",
  "difficulty": "string",
  "timeLimit": "number",
  "memoryLimit": "number",
  "tags": ["string"],
  "inputFormat": "string",
  "outputFormat": "string",
  "constraints": "string",
  "createdAt": "ISO 8601",
  "updatedAt": "ISO 8601"
}
```

#### `GET /challenges`

**Output:**

```json
[
  {
    "id": "uuid",
    "title": "string",
    "description": "string",
    "status": "published",
    "difficulty": "string",
    "tags": ["string"],
    "createdAt": "ISO 8601"
  }
]
```

#### `GET /challenges/:id`

**Output:**

```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "status": "string",
  "difficulty": "string",
  "timeLimit": "number",
  "memoryLimit": "number",
  "tags": ["string"],
  "inputFormat": "string",
  "outputFormat": "string",
  "constraints": "string",
  "testCases": [
    { "id": "uuid", "name": "string", "input": "string", "expectedOutput": "string", "isSample": "boolean", "points": "number" }
  ]
}
```

---

### Test Cases

#### `POST /test-cases`

**Input:**

```json
{
  "challengeId": "uuid",
  "name": "string",
  "input": "string",
  "expectedOutput": "string",
  "isSample": "boolean (opcional, default false)",
  "points": "number (opcional)"
}
```

**Output:**

```json
{
  "id": "uuid",
  "challengeId": "uuid",
  "name": "string",
  "input": "string",
  "expectedOutput": "string",
  "isSample": "boolean",
  "points": "number",
  "createdAt": "ISO 8601"
}
```

---

### Submissions

#### `POST /submissions`

**Input:**

```json
{
  "challengeId": "uuid",
  "code": "string (código fuente)",
  "language": "python | node | cpp | java",
  "examId": "uuid (opcional)"
}
```

**Output:**

```json
{
  "id": "uuid",
  "status": "queued",
  "createdAt": "ISO 8601"
}
```

#### `GET /submissions/:id`

**Output:**

```json
{
  "id": "uuid",
  "challengeId": "uuid",
  "userId": "uuid",
  "code": "string",
  "language": "string",
  "status": "queued | running | accepted | wrong_answer | error",
  "score": "number (0-100)",
  "timeMsTotal": "number",
  "examId": "uuid | null",
  "createdAt": "ISO 8601",
  "updatedAt": "ISO 8601"
}
```

#### `GET /submissions?challengeId=&status=&limit=20&offset=0`

**Output:**

```json
[
  {
    "id": "uuid",
    "challengeId": "uuid",
    "language": "string",
    "status": "string",
    "score": "number",
    "timeMsTotal": "number",
    "createdAt": "ISO 8601"
  }
]
```

---

### Courses

#### `POST /courses`

**Input:**

```json
{
  "name": "string",
  "code": "string",
  "period": "string",
  "groupNumber": "number"
}
```

**Output:**

```json
{
  "id": "uuid",
  "name": "string",
  "code": "string",
  "period": "string",
  "groupNumber": "number",
  "enrollmentCode": "string (ej: CS101-20251G1)",
  "professorId": "uuid",
  "createdAt": "ISO 8601"
}
```

#### `POST /courses/enroll`

**Input:**

```json
{
  "enrollmentCode": "string"
}
```

**Output:**

```json
{
  "message": "Enrolled successfully"
}
```

#### `POST /courses/:id/challenges`

**Input:**

```json
{
  "challengeId": "uuid"
}
```

#### `POST /courses/:id/students`

**Input:**

```json
{
  "studentId": "uuid"
}
```

#### `GET /courses/:id/students`

**Output:**

```json
[
  { "id": "uuid", "username": "string" }
]
```

#### `GET /courses/:id/challenges`

**Output:**

```json
[
  {
    "id": "uuid",
    "title": "string",
    "description": "string",
    "difficulty": "string",
    "status": "string"
  }
]
```

---

### Exams

#### `POST /exams`

**Input:**

```json
{
  "title": "string",
  "description": "string",
  "courseId": "uuid",
  "startTime": "ISO 8601",
  "endTime": "ISO 8601",
  "durationMinutes": "number",
  "challenges": [
    { "challengeId": "uuid", "points": "number", "order": "number" }
  ]
}
```

**Output:**

```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "courseId": "uuid",
  "startTime": "ISO 8601",
  "endTime": "ISO 8601",
  "durationMinutes": "number",
  "createdAt": "ISO 8601"
}
```

#### `GET /exams/:id`

**Output:**

```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "courseId": "uuid",
  "startTime": "ISO 8601",
  "endTime": "ISO 8601",
  "durationMinutes": "number",
  "challenges": [
    { "challengeId": "uuid", "title": "string", "points": "number", "order": "number" }
  ]
}
```

---

### Leaderboard

#### `GET /leaderboard/challenge/:id`

**Output:**

```json
[
  {
    "rank": "number",
    "userId": "uuid",
    "username": "string",
    "score": "number",
    "timeMs": "number",
    "submittedAt": "ISO 8601"
  }
]
```

#### `GET /leaderboard/course/:id`

**Output:**

```json
[
  {
    "rank": "number",
    "userId": "uuid",
    "username": "string",
    "totalScore": "number",
    "challengesSolved": "number",
    "totalTimeMs": "number"
  }
]
```

---

### AI

#### `POST /ai/generate-challenge-ideas`

**Input:**

```json
{
  "topic": "string",
  "difficulty": "easy | medium | hard (opcional)",
  "count": "number (opcional)"
}
```

**Output:**

```json
{
  "ideas": [
    {
      "title": "string",
      "description": "string",
      "difficulty": "string",
      "inputFormat": "string",
      "outputFormat": "string",
      "constraints": "string"
    }
  ]
}
```

#### `POST /ai/generate-test-cases`

**Input:**

```json
{
  "challengeDescription": "string",
  "inputFormat": "string",
  "outputFormat": "string",
  "publicCount": "number (opcional)",
  "hiddenCount": "number (opcional)"
}
```

**Output:**

```json
{
  "publicTestCases": [
    { "name": "string", "input": "string", "expectedOutput": "string" }
  ],
  "hiddenTestCases": [
    { "name": "string", "input": "string", "expectedOutput": "string" }
  ]
}
```

---

### Metrics

#### `GET /metrics`

**Output:**

```json
{
  "submissions_total": "number",
  "submissions_accepted": "number",
  "submissions_rejected": "number",
  "submissions_failed": "number",
  "average_execution_time_ms": "number",
  "challenges_total": "number",
  "courses_total": "number",
  "users_total": "number"
}
```

También retorna formato Prometheus text en la misma respuesta.

---

### Health

#### `GET /health`

```json
{ "status": "ok", "ts": "ISO 8601" }
```

#### `GET /cache/health`

```json
{ "ok": true, "durationMs": "number" }
```

#### `GET /db/health`

```json
{ "ok": true, "durationMs": "number" }
```

---

## Funcionalidades Implementadas

### 1. Autenticación (Auth)

Registro de usuarios con roles (student/professor/admin), login con JWT, y protección de rutas mediante guards. Las contraseñas se hashean con bcrypt.

```mermaid
sequenceDiagram
    actor U as Usuario
    participant API as Auth Controller
    participant S as Auth Service
    participant DB as PostgreSQL

    Note over U,DB: === REGISTRO ===
    U->>API: POST /auth/register {username, password, role}
    API->>S: register(dto)
    S->>S: bcrypt.hash(password)
    S->>DB: INSERT INTO users (username, password_hash, role)
    DB-->>S: user row
    S-->>API: { id, username, role }
    API-->>U: 201 Created

    Note over U,DB: === LOGIN ===
    U->>API: POST /auth/login {username, password}
    API->>S: login(dto)
    S->>DB: SELECT * FROM users WHERE username = ?
    DB-->>S: user row
    S->>S: bcrypt.compare(password, hash)
    S->>S: jwt.sign({ sub: user.id, role })
    S-->>API: { access_token: JWT }
    API-->>U: 200 OK

    Note over U,DB: === RUTA PROTEGIDA ===
    U->>API: GET /auth/me [Authorization: Bearer JWT]
    API->>API: JwtAuthGuard: verificar token
    API->>S: verify(token)
    S-->>API: { id, username, role }
    API-->>U: 200 OK { user data }
```

---

### 2. Challenges (Desafíos)

CRUD de problemas de programación con ciclo de vida (draft → published → archived). Soporta test cases públicos y ocultos desde la creación.

```mermaid
sequenceDiagram
    actor P as Profesor
    participant API as Challenges Controller
    participant UC as Use Cases
    participant DB as PostgreSQL

    Note over P,DB: === CREAR CHALLENGE ===
    P->>API: POST /challenges {title, description, publicTestCases, hiddenTestCases, ...}
    API->>API: JwtAuthGuard + RolesGuard(professor, admin)
    API->>UC: CreateChallengeUseCase.execute(dto)
    UC->>UC: Challenge.create(props) — status = draft
    UC->>DB: INSERT INTO challenges (...)
    DB-->>UC: challenge row
    loop Para cada test case
        UC->>UC: TestCase.create(props)
        UC->>DB: INSERT INTO test_cases (...)
    end
    UC-->>API: challenge entity
    API-->>P: 201 Created

    Note over P,DB: === PUBLICAR ===
    P->>API: POST /challenges/:id/publish
    API->>UC: PublishChallengeUseCase.execute(id)
    UC->>DB: SELECT challenge WHERE id = ?
    UC->>UC: challenge.publish() — status → published
    UC->>DB: UPDATE challenges SET status = 'published'
    UC-->>API: updated challenge
    API-->>P: 200 OK

    Note over P,DB: === LISTAR (ESTUDIANTE) ===
    actor E as Estudiante
    E->>API: GET /challenges
    API->>UC: ListChallengesUseCase.execute()
    UC->>DB: SELECT * FROM challenges WHERE status = 'published'<br/>AND id NOT IN (SELECT challenge_id FROM course_challenges)
    DB-->>UC: [challenges]
    UC-->>API: challenge list
    API-->>E: 200 OK
```

---

### 3. Test Cases (Casos de Prueba)

Gestión de casos de prueba por challenge. Los estudiantes solo pueden ver los sample (públicos); los profesores ven todos.

```mermaid
sequenceDiagram
    actor P as Profesor
    actor E as Estudiante
    participant API as TestCases Controller
    participant UC as Use Cases
    participant DB as PostgreSQL

    Note over P,DB: === CREAR TEST CASE ===
    P->>API: POST /test-cases {challengeId, name, input, expectedOutput, isSample, points}
    API->>API: JwtAuthGuard + RolesGuard(professor, admin)
    API->>UC: CreateTestCaseUseCase.execute(dto)
    UC->>DB: INSERT INTO test_cases (...)
    DB-->>UC: test_case row
    UC-->>API: test case entity
    API-->>P: 201 Created

    Note over E,DB: === LISTAR (ESTUDIANTE — solo samples) ===
    E->>API: GET /test-cases/challenge/:challengeId
    API->>UC: ListTestCasesUseCase.execute(challengeId, user)
    UC->>DB: SELECT * FROM test_cases WHERE challenge_id = ?
    UC->>UC: Filtrar: si role=student → solo is_sample=true
    UC-->>API: [test cases filtrados]
    API-->>E: 200 OK

    Note over P,DB: === ELIMINAR ===
    P->>API: DELETE /test-cases/:id
    API->>UC: DeleteTestCaseUseCase.execute(id)
    UC->>DB: DELETE FROM test_cases WHERE id = ?
    UC-->>API: OK
    API-->>P: 200 OK
```

---

### 4. Submissions (Entregas)

Envío de código fuente que se encola en Redis y se ejecuta en contenedores Docker aislados. El worker procesa la cola, ejecuta el código contra todos los test cases, calcula el score y actualiza el estado.

```mermaid
sequenceDiagram
    actor E as Estudiante
    participant API as Submissions Controller
    participant UC as CreateSubmissionUseCase
    participant DB as PostgreSQL
    participant Redis as Redis Queue
    participant W as Worker
    participant Docker as Docker Container

    Note over E,Docker: === ENVÍO DE CÓDIGO ===
    E->>API: POST /submissions {challengeId, code, language, examId?}
    API->>API: JwtAuthGuard + RolesGuard(student)
    API->>UC: execute(dto, userId)
    UC->>DB: Verificar challenge existe
    UC->>UC: Submission.create(props) — status = queued
    UC->>DB: INSERT INTO submissions (...)
    UC->>UC: Escribir código en /code/{id}/
    UC->>Redis: LPUSH queue:submissions submissionId
    UC-->>API: { id, status: queued }
    API-->>E: 201 Created

    Note over W,Docker: === PROCESAMIENTO ASÍNCRONO ===
    W->>Redis: BRPOP queue:submissions (blocking)
    Redis-->>W: submissionId
    W->>DB: SELECT submission WHERE id = ?
    W->>DB: UPDATE status = 'running'
    W->>DB: SELECT test_cases WHERE challenge_id = ?
    W->>W: Crear /temp_tests/sub-{id}/<br/>Escribir input*.in, output*.out
    W->>W: Seleccionar imagen Docker<br/>(juez_runner_{lang}:local)
    W->>Docker: docker run --rm --network none<br/>--memory 512m --cpus 0.5<br/>-v code:/code:ro -v tests:/tests:ro
    Docker->>Docker: Compilar (si aplica)
    loop Para cada test case
        Docker->>Docker: Ejecutar código con input
        Docker->>Docker: Comparar output vs expected
    end
    Docker-->>W: JSON { status, timeMsTotal, cases[] }
    W->>W: Calcular score = (passedPoints / totalPoints) × 100
    W->>DB: UPDATE submission SET status, score, time_ms_total
    W->>W: Registrar métricas

    Note over E,Docker: === CONSULTAR RESULTADO ===
    E->>API: GET /submissions/:id
    API->>DB: SELECT * FROM submissions WHERE id = ?
    DB-->>API: submission row
    API-->>E: 200 OK { status, score, ... }
```

---

### 5. Courses (Cursos)

Gestión de cursos académicos con código de inscripción autogenerado, matriculación de estudiantes y asignación de challenges.

```mermaid
sequenceDiagram
    actor P as Profesor
    actor E as Estudiante
    participant API as Courses Controller
    participant UC as Use Cases
    participant DB as PostgreSQL

    Note over P,DB: === CREAR CURSO ===
    P->>API: POST /courses {name, code, period, groupNumber}
    API->>API: JwtAuthGuard + RolesGuard(professor, admin)
    API->>UC: CreateCourseUseCase.execute(dto, professorId)
    UC->>UC: Generar enrollmentCode:<br/>"{code}-{period}G{groupNumber}"
    UC->>DB: INSERT INTO courses (...)
    DB-->>UC: course row
    UC-->>API: course entity con enrollmentCode
    API-->>P: 201 Created

    Note over E,DB: === INSCRIPCIÓN CON CÓDIGO ===
    E->>API: POST /courses/enroll {enrollmentCode}
    API->>UC: EnrollStudentUseCase.execute(code, studentId)
    UC->>DB: SELECT * FROM courses<br/>WHERE enrollment_code = ?
    UC->>DB: INSERT INTO course_students<br/>(course_id, student_id)
    UC-->>API: success
    API-->>E: 200 OK

    Note over P,DB: === ASIGNAR CHALLENGE ===
    P->>API: POST /courses/:id/challenges {challengeId}
    API->>UC: AssignChallengeUseCase.execute(courseId, challengeId)
    UC->>DB: INSERT INTO course_challenges<br/>(course_id, challenge_id)
    UC-->>API: success
    API-->>P: 200 OK

    Note over E,DB: === LISTAR CURSOS (por rol) ===
    E->>API: GET /courses
    API->>UC: ListCoursesUseCase.execute(userId, role)
    alt role = student
        UC->>DB: SELECT courses JOIN course_students<br/>WHERE student_id = ?
    else role = professor
        UC->>DB: SELECT courses WHERE professor_id = ?
    else role = admin
        UC->>DB: SELECT * FROM courses
    end
    DB-->>UC: [courses]
    UC-->>API: course list
    API-->>E: 200 OK
```

---

### 6. Exams (Exámenes)

Creación de exámenes con tiempo limitado, asociados a cursos y con challenges ordenados por puntos.

```mermaid
sequenceDiagram
    actor P as Profesor
    actor E as Estudiante
    participant API as Exams Controller
    participant UC as Use Cases
    participant DB as PostgreSQL

    Note over P,DB: === CREAR EXAMEN ===
    P->>API: POST /exams {title, description, courseId,<br/>startTime, endTime, durationMinutes,<br/>challenges: [{challengeId, points, order}]}
    API->>API: JwtAuthGuard + RolesGuard(professor, admin)
    API->>UC: CreateExamUseCase.execute(dto)
    UC->>DB: INSERT INTO exams (...)
    DB-->>UC: exam row
    loop Para cada challenge
        UC->>DB: INSERT INTO exam_challenges<br/>(exam_id, challenge_id, points, order)
    end
    UC-->>API: exam entity
    API-->>P: 201 Created

    Note over E,DB: === OBTENER DETALLE ===
    E->>API: GET /exams/:id
    API->>UC: GetExamDetailsUseCase.execute(id)
    UC->>DB: SELECT exam + JOIN exam_challenges<br/>+ JOIN challenges
    UC->>UC: exam.isActive(now) — verificar vigencia
    DB-->>UC: exam + challenges con puntos y orden
    UC-->>API: exam detail
    API-->>E: 200 OK
```

---

### 7. Leaderboard (Tabla de posiciones)

Rankings por challenge (mejor submission por usuario) y por curso (agregado de todos los challenges).

```mermaid
sequenceDiagram
    actor U as Usuario
    participant API as Leaderboard Controller
    participant UC as Use Cases
    participant DB as PostgreSQL

    Note over U,DB: === RANKING POR CHALLENGE ===
    U->>API: GET /leaderboard/challenge/:id
    API->>UC: GetChallengeLeaderboardUseCase.execute(challengeId)
    UC->>DB: SELECT DISTINCT ON (user_id)<br/>s.*, u.username<br/>FROM submissions s JOIN users u<br/>WHERE challenge_id = ?<br/>ORDER BY user_id, score DESC, time_ms_total ASC
    DB-->>UC: [entries]
    UC->>UC: Asignar rank por posición
    UC-->>API: [{ rank, userId, username, score, timeMs, submittedAt }]
    API-->>U: 200 OK

    Note over U,DB: === RANKING POR CURSO ===
    U->>API: GET /leaderboard/course/:id
    API->>UC: GetCourseLeaderboardUseCase.execute(courseId)
    UC->>DB: — Obtener challenges del curso<br/>— Para cada estudiante: mejor score por challenge<br/>— SUM(scores), COUNT(solved), SUM(time)
    DB-->>UC: [aggregated entries]
    UC->>UC: Ordenar por totalScore DESC, totalTimeMs ASC
    UC-->>API: [{ rank, userId, username, totalScore, challengesSolved, totalTimeMs }]
    API-->>U: 200 OK
```

---

### 8. AI (Generación con IA)

Generación de ideas de challenges y test cases usando Google Gemini Flash. Disponible solo para profesores y admins.

```mermaid
sequenceDiagram
    actor P as Profesor
    participant API as AI Controller
    participant G as GeminiService
    participant Gemini as Google Gemini API

    Note over P,Gemini: === GENERAR IDEAS DE CHALLENGES ===
    P->>API: POST /ai/generate-challenge-ideas<br/>{topic, difficulty?, count?}
    API->>API: JwtAuthGuard + RolesGuard(professor, admin)
    API->>G: generateChallengeIdeas(dto)
    G->>G: Construir prompt con topic,<br/>difficulty, count
    G->>Gemini: generateContent(prompt)
    Gemini-->>G: Respuesta con JSON estructurado
    G->>G: Parsear JSON → ideas[]
    G-->>API: { ideas: [...] }
    API-->>P: 200 OK

    Note over P,Gemini: === GENERAR TEST CASES ===
    P->>API: POST /ai/generate-test-cases<br/>{challengeDescription, inputFormat,<br/>outputFormat, publicCount?, hiddenCount?}
    API->>G: generateTestCases(dto)
    G->>G: Construir prompt con descripción<br/>y formatos I/O
    G->>Gemini: generateContent(prompt)
    Gemini-->>G: Respuesta con JSON estructura
    G->>G: Parsear → publicTestCases[], hiddenTestCases[]
    G-->>API: { publicTestCases, hiddenTestCases }
    API-->>P: 200 OK
```

---

### 9. Metrics (Métricas)

Recolección y exposición de métricas del sistema en formato JSON y Prometheus. Tracked: submissions (total, accepted, rejected, failed), tiempo de ejecución promedio, totales de challenges/cursos/usuarios.

```mermaid
sequenceDiagram
    actor M as Monitor/Prometheus
    participant API as Metrics Controller
    participant MC as MetricsCollector
    participant DB as PostgreSQL

    Note over M,DB: === CONSULTAR MÉTRICAS ===
    M->>API: GET /metrics
    API->>MC: getMetrics()
    MC->>MC: Leer contadores internos:<br/>submissions_total, accepted,<br/>rejected, failed, avg_time
    MC->>DB: COUNT(*) FROM challenges
    MC->>DB: COUNT(*) FROM courses
    MC->>DB: COUNT(*) FROM users
    DB-->>MC: totals
    MC-->>API: Métricas JSON + formato Prometheus
    API-->>M: 200 OK

    Note over MC,MC: === REGISTRO (interno, desde Worker) ===
    Note right of MC: Worker llama tras cada submission:<br/>- incrementSubmissionsTotal()<br/>- incrementSubmissionsAccepted()<br/>- recordExecutionTime(ms)<br/>Promedios rolling de últimas 1000 ejecuciones
```

---

### 10. Health Checks

Endpoints de salud para monitoreo de infraestructura: aplicación, base de datos y caché.

```mermaid
sequenceDiagram
    actor M as Monitor
    participant API as Health/DB/Cache Controllers
    participant DB as PostgreSQL
    participant Redis as Redis

    M->>API: GET /health
    API-->>M: { status: "ok", ts: "2026-03-09T..." }

    M->>API: GET /db/health
    API->>DB: SELECT 1 (con timer)
    DB-->>API: resultado
    API-->>M: { ok: true, durationMs: 2 }

    M->>API: GET /cache/health
    API->>Redis: PING (con timer)
    Redis-->>API: PONG
    API-->>M: { ok: true, durationMs: 1 }
```

---

## Arquitectura General

```mermaid
graph TB
    subgraph Cliente
        WEB[Web App - React]
    end

    subgraph API["API (NestJS)"]
        AUTH[Auth Module]
        CHAL[Challenges Module]
        SUB[Submissions Module]
        COURSE[Courses Module]
        EXAM[Exams Module]
        LB[Leaderboard Module]
        AI[AI Module]
        MET[Metrics Module]
        HC[Health Checks]
    end

    subgraph Infra["Infraestructura"]
        PG[(PostgreSQL)]
        RD[(Redis)]
        GEMINI[Google Gemini API]
    end

    subgraph Worker["Worker Process"]
        W[Submission Worker]
        subgraph Runners["Docker Runners"]
            PY[Python Runner]
            NODE[Node.js Runner]
            CPP[C++ Runner]
            JAVA[Java Runner]
        end
    end

    WEB --> AUTH
    WEB --> CHAL
    WEB --> SUB
    WEB --> COURSE
    WEB --> EXAM
    WEB --> LB
    WEB --> AI
    WEB --> MET

    AUTH --> PG
    CHAL --> PG
    SUB --> PG
    SUB --> RD
    COURSE --> PG
    EXAM --> PG
    LB --> PG
    AI --> GEMINI
    MET --> PG
    HC --> PG
    HC --> RD

    W --> RD
    W --> PG
    W --> PY
    W --> NODE
    W --> CPP
    W --> JAVA
```

---

## Modelo de Datos

```mermaid
erDiagram
    USERS {
        uuid id PK
        string username UK
        string password_hash
        enum role "student|professor|admin"
        timestamp created_at
        timestamp updated_at
    }

    CHALLENGES {
        uuid id PK
        string title
        text description
        enum status "draft|published|archived"
        enum difficulty "easy|medium|hard"
        int time_limit
        int memory_limit
        text[] tags
        text input_format
        text output_format
        text constraints
        timestamp created_at
        timestamp updated_at
    }

    TEST_CASES {
        uuid id PK
        uuid challenge_id FK
        string name
        text input
        text expected_output
        boolean is_sample
        int points
        timestamp created_at
    }

    SUBMISSIONS {
        uuid id PK
        uuid challenge_id FK
        uuid user_id FK
        text code
        string language
        enum status "queued|running|accepted|wrong_answer|error"
        int score
        int time_ms_total
        uuid exam_id FK
        timestamp created_at
        timestamp updated_at
    }

    COURSES {
        uuid id PK
        string name
        string code
        string period
        int group_number
        uuid professor_id FK
        string enrollment_code UK
        timestamp created_at
        timestamp updated_at
    }

    EXAMS {
        uuid id PK
        string title
        text description
        uuid course_id FK
        timestamp start_time
        timestamp end_time
        int duration_minutes
        timestamp created_at
        timestamp updated_at
    }

    COURSE_STUDENTS {
        uuid course_id FK
        uuid student_id FK
    }

    COURSE_CHALLENGES {
        uuid course_id FK
        uuid challenge_id FK
    }

    EXAM_CHALLENGES {
        uuid exam_id FK
        uuid challenge_id FK
        int points
        int order
    }

    USERS ||--o{ SUBMISSIONS : "envía"
    USERS ||--o{ COURSES : "es profesor de"
    CHALLENGES ||--o{ TEST_CASES : "tiene"
    CHALLENGES ||--o{ SUBMISSIONS : "recibe"
    COURSES ||--o{ COURSE_STUDENTS : "tiene"
    USERS ||--o{ COURSE_STUDENTS : "se inscribe"
    COURSES ||--o{ COURSE_CHALLENGES : "incluye"
    CHALLENGES ||--o{ COURSE_CHALLENGES : "asignado a"
    COURSES ||--o{ EXAMS : "tiene"
    EXAMS ||--o{ EXAM_CHALLENGES : "contiene"
    CHALLENGES ||--o{ EXAM_CHALLENGES : "en examen"
    EXAMS ||--o{ SUBMISSIONS : "pertenece a"
```
