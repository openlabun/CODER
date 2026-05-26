# Plan de Pruebas

### Pruebas de Seguridad

- **Pruebas de Acceso de Usuario:** 



```
go test -v ./test/roble_auth -run TestUserLogin
```

- **Pruebas de Registro:** 



```
go test -v ./test/roble_auth -run TestUserRegistration
```

### Pruebas de Persistencia

- **Persistencia de Cursos:** 

```
go test -v ./test/roble_persistence -run TestCourseCRUD
```

- **Persistencia de Exámenes:** 

```
go test -v ./test/roble_persistence -run TestExamCRUD
```

- **Persistencia de Revisiones:**

```
go test -v ./test/roble_persistence -run TestSubmissionCRUD
```

### Pruebas Funcionales

##### Módulo de Autenticación

- **Autenticación de Usuario Estudiante:**



```
go test -v ./test/use_cases/user-run TestStudentAuth
```

- **Autenticación de Usuario Docente:**
  
  

```
go test -v ./test/use_cases/user-run TestTeacherAuth
```

##### Módulo de Cursos

- **CRUD de un Curso:**



```
go test -v ./test/use_cases/exam -run TestCourseCRUD
```

- **Inscripción de Estudiante:**



```
go test -v ./test/use_cases/exam -run TestCourseEnrollment
```

- **Curso desde la vista de Estudiante:**



```
go test -v ./test/use_cases/course -run TestCourseFromStudentView
```

##### Módulo de Exámenes

- **CRUD de un Examen:**



```
go test -v ./test/use_cases/exam -run TestExamCRUD
```

- **Examen (`public`) desde la vista de Estudiante:**



```
go test -v ./test/use_cases/exam -run TestExamFromStudentView
```

- **Examen (`teachers`) desde la vista de Docente:**



```
go test -v ./test/use_cases/exam -run TestExamFromTeacherView
```

- **CRUD de un Reto (`Challenge`):**



```
go test -v ./test/use_cases/exam -run TestChallengeCRUD
```

- **Estados de un Reto:**



```
go test -v ./test/use_cases/exam -run TestChallengeStates
```

- **Fork de un Reto:**



```
go test -v ./test/use_cases/exam -run TestChallengeFork
```

- **Visibilidad de un Reto para Docentes:**



```
go test -v ./test/use_cases/exam -run TestChallengeFromTeachersView
```

- **CRUD de Casos de uso (`TestCase`):**



```
go test -v ./test/use_cases/exam -run TestTestCaseCRUD
```

- **Casos de Uso desde la vista de Estudiante:**



```
go test -v ./test/use_cases/exam -run TestTestCaseFromStudentView
```

- **CRUD de puntos de examen (`ExamItem`):**



```
go test -v ./test/use_cases/exam -run TestExamItemCRUD
```

- **Privacidad de Retos para Puntos de Examen:**



```
go test -v ./test/use_cases/exam -run 
```

- **Fork automático para Punto de Examen:**



```
go test -v ./test/use_cases/exam -run TestExamItemChallengeFork
```

- **Privacidad de Puntos de Examen:**



```
go test -v ./test/use_cases/exam -run TestExamItemPrivacy
```

- **Generación de Plantilla de Código:**



```
go test -v ./test/use_cases/submission -run TestChallengeDefaultCodeTemplates
```

##### Módulo de Revisiones

- **CRUD de Sesiones (`Session`):**



```
go test -v ./test/use_cases/submission -run TestSessionCRUD
```

- **Heartbeat de Sesiones:**



```
go test -v ./test/use_cases/submission -run TestSessionHeartbeat
```

- **Congelamiento y Bloqueo de Sesiones:**



```
go test -v ./test/use_cases/submission -run TestSessionFreezeAndBlock
```

- **Creación y obtención de Revisiones (`Submission`):**



```
go test -v ./test/use_cases/submission -run TestSubmissionCreateAndRead
```

- **Límite de intentos para examen:**



```
go test -v ./test/use_cases/submission -run TestExamTryLimit
```

- **Límite de intentos para revisiones:**



```
go test -v ./test/use_cases/submission -run TestSubmissionTryLimit
```

- **Revisiones Inválidas:**



```
go test -v ./test/use_cases/submission -run TestInvalidSubmissions
```

- **Ejecución de Revisiones:**



```
go test -v ./test/use_cases/submission -run TestSubmissionExecution
```

- **Ejecución de Código sin revisión:**



```
go test -v ./test/use_cases/submission -run TestSubmissionWithoutScore
```

- **Ejecución de Código con caso personalizado:**



```
go test -v ./test/use_cases/submission -run 
```

- **Puntaje de Revisiones:**



```
go test -v ./test/use_cases/submission -run TestSubmissionScoring
```

- **Puntaje del examen:**



### Pruebas de Rendimiento

- **Inicios de sesión simultáneos:**



```
go test -v ./test/use_cases/performance -run TestLoginPerformance
```

- **Lista de exámenes (`Exam`) públicos:**



```
go test -v ./test/use_cases/performance -run TestExamListPerformance
```

- **Activación y Heartbeat de sesiones (`Session`):**



```
go test -v ./test/use_cases/performance -run TestSessionsPerformance
```

- **Revisiones (`Submission`) simultáneas:**



```
go test -v ./test/use_cases/performance -run TestSubmissionPerformance
```
