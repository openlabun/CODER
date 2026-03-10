# CODER — Diagrama Entidad-Relación de Base de Datos

## Diagrama ER por Módulo

```mermaid
erDiagram
    %% ══════════════════════════════════════════
    %% MÓDULO: AUTH
    %% ══════════════════════════════════════════

    USERS {
        UUID id PK "PRIMARY KEY"
        TEXT username UK "UNIQUE, NOT NULL"
        TEXT password "NOT NULL (bcrypt hash)"
        VARCHAR_16 role "CHECK: student | professor"
        TIMESTAMPTZ created_at "DEFAULT NOW()"
        TIMESTAMPTZ updated_at "DEFAULT NOW() — trigger"
    }

    %% ══════════════════════════════════════════
    %% MÓDULO: CHALLENGES
    %% ══════════════════════════════════════════

    CHALLENGES {
        UUID id PK "PRIMARY KEY"
        TEXT title "NOT NULL"
        TEXT description "NOT NULL, DEFAULT ''"
        VARCHAR_16 status "CHECK: draft | published | archived"
        INTEGER time_limit "DEFAULT 1500 (ms)"
        INTEGER memory_limit "DEFAULT 256 (MB)"
        VARCHAR_16 difficulty "CHECK: easy | medium | hard"
        TEXT_ARRAY tags "DEFAULT '{}'"
        TEXT input_format "nullable"
        TEXT output_format "nullable"
        TEXT constraints "nullable"
        TIMESTAMPTZ created_at "DEFAULT NOW()"
        TIMESTAMPTZ updated_at "DEFAULT NOW() — trigger"
    }

    TEST_CASES {
        UUID id PK "PRIMARY KEY"
        UUID challenge_id FK "REFERENCES challenges ON DELETE CASCADE"
        VARCHAR_100 name "NOT NULL"
        TEXT input "NOT NULL"
        TEXT expected_output "NOT NULL"
        BOOLEAN is_sample "DEFAULT false"
        INTEGER points "DEFAULT 10"
        TIMESTAMPTZ created_at "DEFAULT NOW()"
    }

    %% ══════════════════════════════════════════
    %% MÓDULO: SUBMISSIONS
    %% ══════════════════════════════════════════

    SUBMISSIONS {
        UUID id PK "PRIMARY KEY"
        UUID challenge_id FK "NOT NULL"
        TEXT user_id FK "NOT NULL"
        TEXT code "NOT NULL"
        VARCHAR_32 language "NOT NULL"
        VARCHAR_16 status "CHECK: queued | running | accepted | wrong_answer | error"
        INTEGER score "DEFAULT 0 (0-100)"
        INTEGER time_ms_total "DEFAULT 0"
        UUID exam_id FK "nullable, REFERENCES exams ON DELETE SET NULL"
        TIMESTAMPTZ created_at "DEFAULT NOW()"
        TIMESTAMPTZ updated_at "DEFAULT NOW() — trigger"
    }

    %% ══════════════════════════════════════════
    %% MÓDULO: COURSES
    %% ══════════════════════════════════════════

    COURSES {
        UUID id PK "PRIMARY KEY"
        VARCHAR_255 name "NOT NULL"
        VARCHAR_50 code "NOT NULL"
        VARCHAR_20 period "NOT NULL"
        INTEGER group_number "NOT NULL"
        UUID professor_id FK "REFERENCES users"
        VARCHAR_50 enrollment_code UK "UNIQUE — ej: CS101-20251G1"
        TIMESTAMPTZ created_at "DEFAULT NOW()"
        TIMESTAMPTZ updated_at "DEFAULT NOW() — trigger"
    }

    COURSE_STUDENTS {
        UUID course_id PK_FK "REFERENCES courses ON DELETE CASCADE"
        UUID student_id PK_FK "REFERENCES users ON DELETE CASCADE"
        TIMESTAMPTZ enrolled_at "DEFAULT NOW()"
    }

    COURSE_CHALLENGES {
        UUID course_id PK_FK "REFERENCES courses ON DELETE CASCADE"
        UUID challenge_id PK_FK "REFERENCES challenges ON DELETE CASCADE"
        TIMESTAMPTZ assigned_at "DEFAULT NOW()"
    }

    %% ══════════════════════════════════════════
    %% MÓDULO: EXAMS
    %% ══════════════════════════════════════════

    EXAMS {
        UUID id PK "PRIMARY KEY"
        VARCHAR_255 title "NOT NULL"
        TEXT description "nullable"
        UUID course_id FK "REFERENCES courses ON DELETE CASCADE"
        TIMESTAMPTZ start_time "NOT NULL"
        TIMESTAMPTZ end_time "NOT NULL"
        INTEGER duration_minutes "NOT NULL"
        TIMESTAMPTZ created_at "DEFAULT NOW()"
        TIMESTAMPTZ updated_at "DEFAULT NOW() — trigger"
    }

    EXAM_CHALLENGES {
        UUID exam_id PK_FK "REFERENCES exams ON DELETE CASCADE"
        UUID challenge_id PK_FK "REFERENCES challenges ON DELETE CASCADE"
        INTEGER points "NOT NULL, DEFAULT 100"
        INTEGER order "NOT NULL, DEFAULT 0"
    }

    %% ══════════════════════════════════════════
    %% RELACIONES
    %% ══════════════════════════════════════════

    %% Auth ↔ Courses: Un profesor tiene muchos cursos
    USERS ||--o{ COURSES : "professor_id — es profesor de"

    %% Auth ↔ Courses: Muchos estudiantes en muchos cursos
    USERS ||--o{ COURSE_STUDENTS : "student_id — se inscribe en"
    COURSES ||--o{ COURSE_STUDENTS : "course_id — tiene inscritos"

    %% Challenges ↔ Test Cases: Un challenge tiene muchos test cases
    CHALLENGES ||--o{ TEST_CASES : "challenge_id — tiene"

    %% Challenges ↔ Submissions: Un challenge recibe muchas submissions
    CHALLENGES ||--o{ SUBMISSIONS : "challenge_id — recibe"

    %% Auth ↔ Submissions: Un usuario envía muchas submissions
    USERS ||--o{ SUBMISSIONS : "user_id — envía"

    %% Courses ↔ Challenges: Muchos cursos asignan muchos challenges
    COURSES ||--o{ COURSE_CHALLENGES : "course_id — incluye"
    CHALLENGES ||--o{ COURSE_CHALLENGES : "challenge_id — asignado a"

    %% Courses ↔ Exams: Un curso tiene muchos exámenes
    COURSES ||--o{ EXAMS : "course_id — tiene"

    %% Exams ↔ Challenges: Un examen tiene muchos challenges con puntaje
    EXAMS ||--o{ EXAM_CHALLENGES : "exam_id — contiene"
    CHALLENGES ||--o{ EXAM_CHALLENGES : "challenge_id — en examen"

    %% Exams ↔ Submissions: Una submission puede pertenecer a un examen
    EXAMS ||--o{ SUBMISSIONS : "exam_id — recibe entregas de"
```

## Leyenda de Módulos

| Color / Grupo   | Módulo                    | Tablas                                            |
| --------------- | ------------------------- | ------------------------------------------------- |
| **Auth**        | Autenticación y usuarios  | `users`                                           |
| **Challenges**  | Problemas de programación | `challenges`, `test_cases`                        |
| **Submissions** | Entregas de código        | `submissions`                                     |
| **Courses**     | Cursos académicos         | `courses`, `course_students`, `course_challenges` |
| **Exams**       | Exámenes                  | `exams`, `exam_challenges`                        |

## Resumen de Relaciones

```
┌─────────────────────────────────────────────────────────────────────┐
│                        MÓDULO AUTH                                  │
│  ┌──────────┐                                                      │
│  │  USERS   │──────────────────────────────┐                       │
│  └──────────┘                              │                       │
│       │ 1                                  │ 1                     │
│       │                                    │                       │
└───────┼────────────────────────────────────┼───────────────────────┘
        │                                    │
        │ envía (N)                           │ es profesor de (N)
        ▼                                    ▼
┌───────────────────────┐   ┌──────────────────────────────────────┐
│   MÓDULO SUBMISSIONS  │   │          MÓDULO COURSES              │
│  ┌──────────────┐     │   │  ┌─────────┐                        │
│  │ SUBMISSIONS  │◄────┼───┼──│ COURSES │                        │
│  └──────────────┘     │   │  └─────────┘                        │
│       ▲               │   │      │ 1          │ 1      │ 1      │
│       │               │   │      ▼ N          ▼ N      ▼ N      │
│       │               │   │ ┌──────────┐ ┌─────────┐ ┌───────┐ │
│       │               │   │ │ COURSE_  │ │ COURSE_ │ │       │ │
│       │               │   │ │ STUDENTS │ │CHALLENG.│ │ EXAMS │ │
│       │               │   │ └──────────┘ └─────────┘ └───────┘ │
│       │               │   │      ▲ N          ▲ N      │ 1     │
│       │               │   └──────┼────────────┼────────┼───────┘
│       │               │          │            │        │
│       │               │          │            │        ▼ N
│       │               │   ┌──────┘    ┌───────┘  ┌───────────┐
│       │               │   │           │          │   MÓDULO   │
│       │               │   │   ┌───────┘          │   EXAMS    │
└───────┼───────────────┘   │   │          ┌───────┴──────────┐ │
        │ challenge_id (N)  │   │          │ EXAM_CHALLENGES  │ │
        │                   │   │          └──────────────────┘ │
        │ exam_id (N)       │   │                  ▲ N          │
┌───────┴───────────────────┴───┴──────────────────┼────────────┘
│              MÓDULO CHALLENGES                    │
│  ┌──────────────┐                                │
│  │  CHALLENGES  │────────────────────────────────┘
│  └──────────────┘
│       │ 1
│       ▼ N
│  ┌──────────────┐
│  │  TEST_CASES  │
│  └──────────────┘
└───────────────────────────────────────────────────┘
```

## Índices

| Tabla               | Índice                            | Columnas                                   | Notas                                 |
| ------------------- | --------------------------------- | ------------------------------------------ | ------------------------------------- |
| `challenges`        | `idx_challenges_created_at`       | `created_at DESC`                          | Listados recientes                    |
| `challenges`        | `idx_challenges_difficulty`       | `difficulty`                               | Filtrado por dificultad               |
| `submissions`       | `idx_submissions_challenge`       | `challenge_id`                             | Lookup por challenge                  |
| `submissions`       | `idx_submissions_created`         | `created_at DESC`                          | Listados recientes                    |
| `submissions`       | `idx_submissions_challenge_score` | `challenge_id, score DESC, created_at ASC` | Leaderboard (WHERE status = accepted) |
| `submissions`       | `idx_submissions_user_challenge`  | `user_id, challenge_id, score DESC`        | Mejor score por usuario               |
| `submissions`       | `idx_submissions_exam`            | `exam_id`                                  | Submissions de examen                 |
| `test_cases`        | `idx_test_cases_challenge`        | `challenge_id`                             | Lookup por challenge                  |
| `test_cases`        | `idx_test_cases_sample`           | `challenge_id, is_sample`                  | Filtrado samples                      |
| `courses`           | `idx_courses_professor`           | `professor_id`                             | Cursos del profesor                   |
| `courses`           | `idx_courses_enrollment_code`     | `enrollment_code`                          | Lookup por código                     |
| `course_students`   | `idx_course_students_student`     | `student_id`                               | Cursos del estudiante                 |
| `course_challenges` | `idx_course_challenges_challenge` | `challenge_id`                             | Cursos del challenge                  |
| `exams`             | `idx_exams_course`                | `course_id`                                | Exámenes del curso                    |

## Triggers

Todas las tablas con `updated_at` tienen un trigger `trg_{tabla}_updated_at` que ejecuta la función `set_updated_at()` antes de cada `UPDATE`, actualizando automáticamente el timestamp.

| Tabla         | Trigger                      |
| ------------- | ---------------------------- |
| `challenges`  | `trg_challenges_updated_at`  |
| `submissions` | `trg_submissions_updated_at` |
| `users`       | `trg_users_updated_at`       |
| `courses`     | `trg_courses_updated_at`     |
| `exams`       | `trg_exams_updated_at`       |
