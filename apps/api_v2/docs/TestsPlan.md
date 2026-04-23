# Plan de Pruebas

### Pruebas de Seguridad

Actualmente el sistema soporta la autenticación de usuario a través de la plataforma institucional Roble [[1]](https://roble.openlab.uninorte.edu.co/docs). Existen pruebas directas de acceso a través de este servicio.

- **Pruebas de Acceso de Usuario:** Accede al usuario `test@test.com` creado por defecto anteriormente con un hash de la contraseña `Password123!`

```
go test -v ./test/roble_auth -run TestUserLogin
```

- **Pruebas de Registro:** Crea un usuario con el correo `test{random_value}@test.com` y contraseña `Password123!`

```
go test -v ./test/roble_auth -run TestUserRegistration
```

### Pruebas de Persistencia

El sistema de almacenamiento de datos, de igualmanera está soportado por la plataforma institucional Roble. En todos se usa el acceso del usuario `test@test.com`, creado con **permisos de docente**. Se realizó pruebas del **CRUD completo de datos** para cada uno de los módulos de la aplicación:

- **Persistencia de Cursos:** Se crea, actualiza, verifica y se elimina un curso (`Course`).

```
go test -v ./test/roble_persistence -run TestCourseCRUD
```

- **Persistencia de Exámenes:** Inicialmente se crea un curso (`Course`) al que vincular un examen (`Exam`), para posteriormente hacer un reto (`Challenge`) con sus respectivos casos de prueba (`TestCase`), el cual se vinculará al examen a través de un punto (`ExamItem`). En cada caso se hará una prueba de CRUD y verificación que el reto esté vinculado al examen. Para los casos de retos y casos de prueba se deberán crear variables (`IOVariable`).

```
go test -v ./test/roble_persistence -run TestExamCRUD
```

- **Persistencia de Revisiones:** En este proceso se crean un curso (`Course`), examen (`Exam`) y reto (`Challenge`) auxiliares. Posteriormente se procede a crear una sesión (`Session`), se actualizará y verificará su estado y valores. A través de la sesión se creará una revisión (`Submission`) y finalmente los resultados de la revisión (`SubmissionResults`), en ambos casos de validará todo el proceso de CRUD.

```
go test -v ./test/roble_persistence -run TestSubmissionCRUD
```

### Pruebas Funcionales

Este es un set de pruebas de construye previamente una instancia de la capa de aplicación de la API y posteriormente usa los casos de uso para verificar la integración de las dependencias, su éxito al resolver la funcionalidad y el tiempo que toma por cada caso de uso. Los accesos de usuario se hacen a través del usuario estudiante (`stud@test.com`, contraseña sin hash: `Password123!`) y el usuario docente (`test@test.com`, contraseña sin hash: `Password123!`), los cuales deben ser previamente creados con sus respectivos permisos en la base de datos.

##### Módulo de Autenticación

- **Autenticación de Usuario Estudiante:**
  
  - Paso 1: Intenta hacer login al usuario de estudiante
  
  - Paso 2: Si no lo logra, crea el usuario
  
  - Paso 3: Obtiene la información del estudiante
  
  ```
  go test -v ./test/use_cases/user-run TestStudentAuth
  ```

- **Autenticación de Usuario Docente:**
  
  - Paso 1: Intenta hacer login al usuario de docente
  
  - Paso 2: Intenta registrar el usuario (espera `error`)
  
  - Paso 3: Obtiene la información del docente
  
  - Paso 4: Verifica su rol
  
  ```
  go test -v ./test/use_cases/user-run TestTeacherAuth
  ```

##### Módulo de Cursos

- **CRUD de un Curso:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crea un curso
  
  - Paso 3: Actualiza datos del curso
  
  - Paso 4: Obtiene datos del curso y los valida
  
  - Paso 5: Elimina el curso
  
  - Paso 6: Verifica eliminación
  
  ```
  go test -v ./test/use_cases/exam -run TestCourseCRUD
  ```

- **Inscripción de Estudiante:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crea un curso
  
  - Paso 3: Inscribe un estudiante al curso
  
  - Paso 4: Obtiene datos del curso y verifica inscripción
  
  - Paso 5: Retira un estudiante del curso
  
  - Paso 6: Obtiene datos del curso y verifica retiro
  
  ```
  go test -v ./test/use_cases/exam -run TestCourseEnrollment
  ```

- **Curso desde la vista de Estudiante:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crea un curso
  
  - Paso 3: Iniciar sesión con usuario de estudiante
  
  - Paso 4: Obtiene el curso desde usuario estudiante (espera `error`)
  
  - Paso 5: Inscribe al estudiante al curso
  
  - Paso 6: Obtiene el curso desde usuario estudiante
  
  - Paso 7: Retira un estudiante del curso
  
  - Paso 8: Obtiene el curso desde usuario estudiante (espera `error`)
  
  ```
  go test -v ./test/use_cases/course -run TestCourseFromStudentView
  ```

##### Módulo de Exámenes

- **CRUD de un Examen:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crea un curso
  
  - Paso 3: Crea un examen
  
  - Paso 4: Actualiza el examen
  
  - Paso 5: Obtiene los datos del examen y valida el cambio
  
  - Paso 6: Hace un cambio erróneo en el examen (espera `error`)
  
  - Paso 7: Elimina el examen
  
  - Paso 8: Verifica eliminación
  
  ```
  go test -v ./test/use_cases/exam -run TestExamCRUD
  ```

- **Examen (`public`) desde la vista de Estudiante:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crea un curso
  
  - Paso 3: Crea un examen con visibilidad `private`
  
  - Paso 4: Iniciar sesión con usuario de estudiante
  
  - Paso 5: Obtener datos del examen (espera `error`)
  
  - Paso 6: Crea un examen con visibilidad `course`
  
  - Paso 7: Obtiene datos del examen (espera `error`)
  
  - Paso 8: Inscribe al estudiante en el curso
  
  - Paso 9: Obtiene datos del examen
  
  - Paso 10: Crea un examen con visibilidad `public` y sin curso
  
  - Paso 11: Obtiene datos del examen
  
  ```
  go test -v ./test/use_cases/exam -run TestExamFromStudentView
  ```

- **Examen (`teachers`) desde la vista de Docente:**
  
  - Paso 1: Iniciar sesión con usuario de docente (creador)
  
  - Paso 2: Crea un examen con visibilidad `private`
  
  - Paso 3: Iniciar sesión con segundo usuario de docente (observador)
  
  - Paso 4: Obtener el examen desde el docente observador (espera `error`)
  
  - Paso 5: Crea un examen con visibilidad `teachers`
  
  - Paso 6: Obtener el examen desde el docente observador
  
  - Paso 7: Obtener el examen desde la vista de estudiante (espera `error`)
  
  - Paso 8: Crea un curso
  
  - Paso 9: Crea un examen con visibilidad `teachers` asociado al curso
  
  - Paso 10: Obtener el examen desde el docente observador
  
  ```
  go test -v ./test/use_cases/exam -run TestExamFromTeacherView
  ```

- **CRUD de un Reto (`Challenge`):**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crea un reto 
  
  - Paso 3: Actualiza el reto
  
  - Paso 4: Obtiene los datos del reto y valida los cambios
  
  - Paso 5: Elimina el reto
  
  - Paso 6: Verifica eliminación
  
  ```
  go test -v ./test/use_cases/exam -run TestChallengeCRUD
  ```

- **Estados de un Reto:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crea un reto en estado `draft`
  
  - Paso 3: Actualiza el reto
  
  - Paso 4: Obtiene los datos del reto y valida los cambios
  
  - Paso 5: Actualiza el reto a estado `published`
  
  - Paso 6: Actualiza el reto (espera `error`)
  
  - Paso 7: Actualiza el reto a estado `private`
  
  - Paso 8: Actualiza el reto (espera `error`)
  
  - Paso 9: Actualiza el reto a estado `archived`
  
  - Paso 10: Actualiza el reto (espera `error`)
  
  - Paso 11: Actualiza el reto a estado `published`
  
  - Paso 12: Actualiza el reto a estado `archived`
  
  - Paso 13: Actualiza el reto a estado `private`
  
  - Paso 14: Actualiza el reto a estado `draft` (espera `error`)
  
  - Paso 15: Actualiza el reto a estado `published`
  
  - Paso 16: Actualiza el reto a estado `draft` (espera `error`)
  
  ```
  go test -v ./test/use_cases/exam -run TestChallengeStates
  ```

- **Fork de un Reto:**
  
  - Paso 1: Iniciar sesión con usuario de docente (creador)
  
  - Paso 2: Crea un reto en estado `private`
  
  - Paso 3: Iniciar sesión con usuario de docente (observador)
  
  - Paso 4: Hace fork al reto (espera `error`)
  
  - Paso 5: Actualiza 
  
  - Paso 4: Hace fork al reto
  
  - Paso 5: Actualiza el reto copiado
  
  - Paso 6: Verifica que no haya cambios en el reto original
  
  - Paso 7: Verifica los cambios en el reto copiado
  
  ```
  go test -v ./test/use_cases/exam -run TestChallengeFork
  ```

- **Visibilidad de un Reto para Docentes:**
  
  - Paso 1: Iniciar sesión con usuario de docente (creador)
  
  - Paso 2: Crea un reto (`private`)
  
  - Paso 3: Iniciar sesión con usuario de docente (observador)
  
  - Paso 4: Obtiene datos del reto con docente observador (espera `error`)
  
  - Paso 5: Actualiza el reto a visibilidad `published`
  
  - Paso 6: Obtiene datos del reto con docente observador
  
  - Paso 7: Actualiza el reto a visibilidad `archived`
  
  - Paso 8: Obtiene datos del reto con docente observador (espera `error`)
  
  ```
  go test -v ./test/use_cases/exam -run TestChallengeFromTeachersView
  ```

- **CRUD de Casos de uso (`TestCase`):**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear un reto
  
  - Paso 3: Crear un caso de uso sin valores de entrada (espera `error`)
  
  - Paso 4: Crear un caso de uso válido
  
  - Paso 4: Actualizar el caso de uso
  
  - Paso 5: Eliminar el caso de uso
  
  - Paso 6: Verificar eliminación
  
  ```
  go test -v ./test/use_cases/exam -run TestTestCaseCRUD
  ```

- **Casos de Uso desde la vista de Estudiante:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear un reto
  
  - Paso 3: Crear un caso de prueba (`isSample == true`)
  
  - Paso 4: Crear un caso de prueba (`isSample == false`)
  
  - Paso 5: Obtener casos de prueba con vista de Docente
  
  - Paso 6: Iniciar sesión con usuario de estudiante
  
  - Paso 7: Obtener casos de prueba con vista de Estudiante (espera solo 1)
  
  ```
  go test -v ./test/use_cases/exam -run TestTestCaseFromStudentView
  ```

- **CRUD de puntos de examen (`ExamItem`):**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear un examen
  
  - Paso 3: Crear un reto
  
  - Paso 4: Crear un punto de examen
  
  - Paso 5: Crear otro punto de examen con mismo reto (espera `error`)
  
  - Paso 6: Actualizar punto de examen
  
  - Paso 7: Obtener punto de examen y validar datos
  
  - Paso 8: Eliminar punto de examen
  
  - Paso 9: Verificar eliminación
  
  ```
  go test -v ./test/use_cases/exam -run TestExamItemCRUD
  ```

- **Privacidad de Retos para Puntos de Examen:**
  
  - Paso 1: Iniciar sesión con usuario de docente (creador)
  
  - Paso 2: Crear un reto `private`
  
  - Paso 3: Iniciar sesión con usuario de docente (publicador)
  
  - Paso 4: Crear un reto `private` con docente publicador
  
  - Paso 5: Crear un reto `published `con docente publicador
  
  - Paso 6: Crear examen con docente creador
  
  - Paso 7: Crear punto de examen con reto `private` del docente creador
  
  - Paso 8: Crear punto de examen con reto `private` del docente publicador (espera `error`)
  
  - Paso 9: Creara punto de examen con reto `published`
  
  ```
  go test -v ./test/use_cases/exam -run 
  ```

- **Fork automático para Punto de Examen:**
  
  - Paso 1: Iniciar sesión con usuario de docente (creador)
  
  - Paso 2: Crear un examen
  
  - Paso 3: Iniciar sesión con usuario de docente (publicador)
  
  - Paso 4: Crear un reto `published` con docente publicador
  
  - Paso 5: Crea un punto de examen con docente creador
  
  - Paso 6: Actualizar reto con docente publicador
  
  - Paso 7: Verificar cambios en el reto
  
  - Paso 8: Obtener punto de examen
  
  - Paso 9: Verificar que no se hayan producido los cambios en el reto del punto de examen
  
  ```
  go test -v ./test/use_cases/exam -run TestExamItemChallengeFork
  ```

- **Privacidad de Puntos de Examen:**
  
  - Paso 1: Iniciar sesión con usuario de docente (creador)
  
  - Paso 2: Crear un examen
  
  - Paso 3: Crear un reto `published`
  
  - Paso 4: Iniciar sesión con usuario de docente (observador)
  
  - Paso 5: Crear un punto de examen desde docente observador (espera `error`)
  
  - Paso 6: Crear un punto de examen desde docente creador
  
  - Paso 7: Actualiza punto de examen desde docente observador (espera `error`)
  
  - Paso 8: Eliminar punto de examen desde docente observador (espera `error`)
  
  ```
  go test -v ./test/use_cases/exam -run TestExamItemPrivacy
  ```

- **Generación de Plantilla de Código:**
  
  - Paso 1: Iniciar sesión con usuario de docente (creador)
  
  - Paso 2: Crear un reto
  
  - Paso 3: Crear casos de prueba
  
  - Paso 4: Obtener plantillas por defecto para el reto
  
  - Paso 5: Validar que se reciban todas las variables esperadas y el print con el output

##### Módulo de Revisiones

- **CRUD de Sesiones (`Session`):**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear examen público (visibilidad `public` y sin curso)
  
  - Paso 3: Crear otro examen público
  
  - Paso 4: Iniciar sesión con usuario de estudiante
  
  - Paso 5: Crear sesión con examen
  
  - Paso 6: Crear sesión con el otro examen (espera `error`)
  
  - Paso 7: Obtener la sesión
  
  - Paso 8: Cerrar la sesión
  
  - Paso 9: Obtener la sesión y confirmar cierre
  
  ```
  go test -v ./test/use_cases/submission -run TestSessionCRUD
  ```

- **Heartbeat de Sesiones:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear examen público (visibilidad `public` y sin curso)
  
  - Paso 3: Iniciar sesión con usuario de estudiante
  
  - Paso 4: Crear sesión con examen
  
  - Paso 5: Obtener la sesión
  
  - Paso 6: Hacer heartbeat a la sesión
  
  - Paso 7: Obtener la sesión
  
  ```
  go test -v ./test/use_cases/submission -run TestSessionHeartbeat
  ```

- **Congelamiento y Bloqueo de Sesiones:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear examen público (visibilidad `public` y sin curso)
  
  - Paso 3: Iniciar sesión con usuario de estudiante
  
  - Paso 4: Crear sesión con examen
  
  - Paso 5: Bloquear sesión desde cuenta de docente
  
  - Paso 6: Obtener la sesión
  
  - Paso 7: Crear sesión con examen
  
  - Paso 8: Esperar tiempo para congelamiento de examen
  
  - Paso 9: Obtener la sesión y comprobar que está congelada
  
  - Paso 10: Hacer heartbeat
  
  - Paso 11: Comprobar que se volvió a activar
  
  - Paso 12: Bloquear sesión desde vista de docente
  
  - Paso 13: Obtener la sesión y comprobar que está bloqueada
  
  ```
  go test -v ./test/use_cases/submission -run TestSessionFreezeAndBlock
  ```

- **Creación y obtención de Revisiones (`Submission`):**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear examen público (visibilidad `public` y sin curso)
  
  - Paso 3: Crear un reto
  
  - Paso 4: Crear casos de prueba
  
  - Paso 5: Crear un punto de examen
  
  - Paso 6: Iniciar sesión con usuario de estudiante
  
  - Paso 7: Crear una sesión en el examen
  
  - Paso 8: Crear una revisión
  
  - Paso 9: Obtener revisiones a partir del ID del reto
  
  - Paso 10: Obtener revisiones a partir del ID de la sesión
  
  - Paso 11: Obtener revisiones a partir del ID del usuario
  
  - Paso 12: Obtener el `status` de la revisión
  
  ```
  go test -v ./test/use_cases/submission -run TestSubmissionCreateAndRead
  ```

- **Límite de intentos:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear examen público (con solo 2 intentos `try_limit`)
  
  - Paso 3: Iniciar sesión con usuario de estudiante
  
  - Paso 4: Crear una sesión en el examen
  
  - Paso 5: Cerrar la sesión
  
  - Paso 6: Crear una sesión en el examen
  
  - Paso 7: Cerrar la sesión
  
  - Paso 8: Crear una sesión en el examen (espera `error`)

- **Revisiones Inválidas:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear examen público (visibilidad `public`, sin curso y 60 segundos de tiempo para resolver)
  
  - Paso 3: Crear un reto
  
  - Paso 4: Crear casos de prueba
  
  - Paso 5: Crear un punto de examen
  
  - Paso 6: Iniciar sesión con usuario de estudiante
  
  - Paso 7: Crear una revisión sin sesión (espera `error`)
  
  - Paso 8: Crear una sesión en el examen
  
  - Paso 9: Cerrar el examen desde la vista de docente
  
  - Paso 10: Crear una revisión (espera `error`)
  
  - Paso 11: Esperar 61 segundos y crear un revisión (espera `error`)
  
  - Paso 12: Obtener datos de sesión
  
  - Paso 13: Confirmar que la sesión tiene estado `expired`
  
  - Paso 14: Crear una sesión en el examen
  
  - Paso 15: Bloquear la sesión desde la vista de docente
  
  - Paso 16: Crear una revisión (espera `error`)
  
  - Paso 17: Crear una sesión en el examen
  
  - Paso 18: Cerrar el examen
  
  - Paso 19: Crear una revisión (espera `error`)
  
  ```
  go test -v ./test/use_cases/submission -run TestInvalidSubmissions
  ```

- **Ejecución de Revisiones:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear examen público (visibilidad `public` y sin curso)
  
  - Paso 3: Crear un reto
  
  - Paso 4: Crear casos de prueba
  
  - Paso 5: Crear un punto de examen
  
  - Paso 6: Iniciar sesión con usuario de estudiante
  
  - Paso 7: Crear una revisión
  
  - Paso 8: Obtener el `status` de la revisión hasta que su estado sea `accepted`
  
  ```
  go test -v ./test/use_cases/submission -run TestSubmissionExecution
  ```

- **Ejecución de Código sin revisión:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear examen público (visibilidad `public` y sin curso)
  
  - Paso 3: Crear un reto
  
  - Paso 4: Crear 2 casos de prueba con valor de 3 puntos
  
  - Paso 5: Crear un caso de prueba con valor de 6 puntos (debe ser imposible de cumplir)
  
  - Paso 6: Iniciar sesión con usuario de estudiante
  
  - Paso 7: Crear una revisión
  
  - Paso 8: Obtener el `status` de la revisión hasta que su estado sea `accepted` o `wrong_answer`
  
  - Paso 9: Confirmar valor del atributo `Score` de la revisión corresponde a 0

- **Ejecución de Código con caso personalizado:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear examen público (visibilidad `public` y sin curso)
  
  - Paso 3: Crear un reto
  
  - Paso 4: Crear 2 casos de prueba con valor de 3 puntos
  
  - Paso 5: Crear un caso de prueba con valor de 6 puntos (debe ser imposible de cumplir)
  
  - Paso 6: Iniciar sesión con usuario de estudiante
  
  - Paso 7: Crear una revisión con un caso de prueba personalizado
  
  - Paso 8: Obtener el `status` de la revisión hasta que su estado sea `accepted` o `wrong_answer`
  
  - Paso 9: Confirmar valor del atributo `Score` de la revisión corresponde a 0

- **Puntaje de Revisiones:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear examen público (visibilidad `public` y sin curso)
  
  - Paso 3: Crear un reto
  
  - Paso 4: Crear 2 casos de prueba con valor de 3 puntos
  
  - Paso 5: Crear un caso de prueba con valor de 6 puntos (debe ser imposible de cumplir)
  
  - Paso 5: Crear un punto de examen
  
  - Paso 6: Iniciar sesión con usuario de estudiante
  
  - Paso 7: Crear una revisión
  
  - Paso 8: Obtener el `status` de la revisión hasta que su estado sea `accepted` o `wrong_answer`
  
  - Paso 9: Confirmar valor del atributo `Score` de la revisión corresponde a 6
  
  ```
  go test -v ./test/use_cases/submission -run TestSubmissionScoring
  ```

### Pruebas de Rendimiento

Este es un conjunto de pruebas desarrolladas para **medir la resiliencia del sistema** ante solicitudes recurrentes en secciones específicas de los casos de uso contemplados para la aplicación. Para cada sección marcada como **crítica** se ejecutarán *n* solicitudes simultáneas y se medirá el tiempo de respuesta de cada una y el tiempo promedio. 

Posterior a la finalización del paso se mostrarán los datos individuales de cada medición. Y al finalizar la prueba se mostrará la solicitud que más demoró, la que menos y el promedio.

- **Inicios de sesión simultáneos:**
  
  - Paso 1: Registro de usuarios de estudiantes (**crítica**)
  
  - Paso 2: Inicio de sesión de usuarios de estudiantes (**crítica**)
  
  - Paso 3: Refrescar token (**crítica**)
  
  - Paso 4: Obtener datos del usuario (**crítica**)

- **Lista de exámenes (`Exam`) públicos:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear exámenes públicos (**crítica**)
  
  - Paso 3: Iniciar sesión con usuario de estudiante
  
  - Paso 4: Obtener lista de exámenes públicos (**crítica**)

- **Activación y Heartbeat de sesiones (`Session`):**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear examen público
  
  - Paso 3: Iniciar sesión con usuarios de estudiantes (**crítica**)
  
  - Paso 4: Cada usuario creará una sesión (**crítica**)
  
  - Paso 5: Cada usuario hará heartbeat a la sesión (**crítica**)
  
  - Paso 6: Esperar el tiempo de `FREEZE_TIME`
  
  - Paso 7: Cada usuario hará heartbeat a la sesión (**crítica**)

- **Revisiones (`Submission`) simultáneas:**
  
  - Paso 1: Iniciar sesión con usuario de docente
  
  - Paso 2: Crear examen público
  
  - Paso 3: Crear reto
  
  - Paso 4: Crear casos de prueba
  
  - Paso 5: Crear punto de examen
  
  - Paso 6: Iniciar sesión con usuarios de estudiantes (**crítica**)
  
  - Paso 7: Cada usuario subirá una revisión (**crítica**)
  
  - Paso 8: Cada usuario revisará el estado de la revisión hasta obtener `accepted` (**crítica**)
