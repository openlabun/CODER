# Reglas de Negocio

## Módulo Usuarios

- Los usuarios solo podrán manejar los siguientes **roles** (valores en el atributo `Role`):
  
  - `student`: Acceso solo a las entidades `CourseStudent` y `Submission`
  
  - `professor`: acceso a las entidades `Course` y `Exam`, pero solo a aquellas que le pertenecen, es decir, su `ID` sea igual al valor de `ProfessorID`. También tendrá acceso a las entidades de `Challenge` y `TestCase` que hagan parte de un `Exam` que le pertenezca.
  
  - `admin`: acceso a todas las demás entidades

## Módulo Cursos

- Los cursos tendrán un control de visibilidad:
  
  - `public`: cualquier estudiante puede acceder a él con el código de acceso.
  
  - `private`: solo se puede acceder a través de `EnrollmentCode` o `EnrollmentURL`.
  
  - `blocked`: ningún estudiante puede acceder al curso.

- Los cursos estarán disponibles durante un periodo determinado, deberá ingresarse el año y el semestre, la API maneja los siguiente valores para el semestre:
  
  - `01`: Primer Semestre
  
  - `02`: Intersemestral
  
  - `03`: Segundo Semestre

- Si no se ingresa un periodo el curso será accesible de manera indeterminada

## Módulo Exámenes

#### Exámenes (`Exam`)

- Un examen (`Exam`) contiene varios retos (`Challenge`)

- Un examen (`Exam`) puede estar o no vinculado a un curso (`Course`), si no está vinculado, podrá ser visto por todos los usuarios estudiantes (según configuración de visibilidad).

- La configuración de visbilidad se manejará de la siguiente manera:
  
  - `public`: accesible como un reto (`Challenge`) disponible para todos los usuarios en la galería de retos
  
  - `course`: solo visible para el curso (`Course`) al que se encuentra vinculado el examen, si no está vinculado esta opción no deberá ser válida
  
  - `teachers`: solo visible para otros usuarios profesores en la galería de retos.
  
  - `private`: solo será visible para el usuario que lo creó.

- El examen podrá configurarse con una fecha de inicio (atributo `StartTime`) y con una fecha de cierre (atributo `EndTime`). Después de la fecha de cierre el examen no recibirá revisiones (`Submission`), a menos que se encuentre activada la configuración para aceptar respuestas tardías (`AllowLateSubmissions == True`).

- El examen podrá tener un tiempo límite (atributo `TimeLimit`) y un límite de intentos (atributo `TryLimit`). Si alguna revisión (`Submission`) se carga después del tiempo límite no será recibida, y si se excede el límite de intentos tampoco.

#### Retos (`Challenge`)

- Los retos (`Challenge`) pueden guardarse como:
  
  - `draft`: permiten modificación, estarán en el examen pero no serán visibles
  
  - `published`: no permiten modificación, estarán visibles en el examen y permitirán ejecuciones.
  
  - `archived`: no permiten modificación, tampoco estarán visibles en el examen.

- Una vez publicado un reto no podrá modificarse, solo podrá archivarse. Y una vez archivado, puede volver a publicarse.

- Cada reto deberá contar con un time límite de ejecución y una memoria límite para el worker. **(No podrá exceder determinados valores fijados en las variables en producción)**.

- Los retos deberán contar con unas consideraciones de entrada/salida. Para ello deberán ingresar por cada variable los siguientes valores:
  
  - Nombre (`Name`)
  
  - Tipo (`Type`: `string`, `int`, `float`)
  
  - Valor esperado/ejemplo (`Value`)

#### Casos de Prueba (`TestCase`)

- Deberán estar vinculados a un reto

- Cada uno deberá contar con **uno o más** valores de entrada (con la misma estructura que la definida anteriormente) y un valor de salida esperado.

- Existirán casos de prueba de ejemplo para el estudiante (serán visibles) y existirán casos de prueba secretos (atributo `IsSample == False`). Podrá asignarse el puntaje por pasar el caso de prueba con el atributo `Points`.

## Módulo Revisiones

#### Revisiones (`Submission`)

- Pueden ejecutarse con varios lenguajes de programación, hasta ahora:
  
  - `cpp`: C++
  
  - `python`: Python
  
  - `java`: Java

- Las revisiones ejecutan todos los casos de prueba y generan un puntaje (atributo `Score`) a partir de ellos.

- Deberán siempre estar vinculadas a un reto (`Challenge`), a una sesión (`Session`) y a un usuario (`User`).

#### Resultado de Revisiones (`SubmissionResult`)

- Cada uno estará vinculado a un caso de prueba (`TestCase`) que se ejecutará.

- Pueden tener 5 estados posibles:
  
  - `queued`: está en la cola de ejecución
  
  - `running`: en la última actualización ya se estaba ejecutando
  
  - `accepted`: pasó exitosamente el caso de prueba
  
  - `wrong_answer`: ejecutó pero la respuesta fue incorrecta
  
  - `error`: no se pudo compilar el código

- Dependiendo del estado resultante pueden guardarse los valores para:
  
  - `ActualOutput`: si el estado fue `accepted` o `wrong_answer`, muestra el resultado generado.
  
  - `ErrorMessage`: mensaje de error resultante cuando el estado fue `error`.

#### Sesión de Usuario (`Session`)

- Las sesiones se abren cuando el usuario inicia un examen (`Exam`), solo a través de una sesión activa se permite la ejecución de revisiones.

- Las sesiones varían entre los siguientes estados:
  
  - `active`: permite hacer ejecuciones, si el temporizador está activado seguirá contando el tiempo.
  
  - `frozen`: no permite hacer ejecuciones, se genera automáticamente cuando el usuario lleva más de un tiempo determinado de actividad (se programará desde variables de entorno). El navegador deberá enviar `Heartbeats` continuamente para evitar este estado. En este estado el contador no tiene efecto hasta que el usuario se reconecte.
  
  - `completed`:  el usuario finalizó exitosamente el examen, se cierra la sesión.
  
  - `expired`: se cierra la sesión cuando el tiempo se agota y no se permiten revisiones atrasadas (atributo `AllowLateSubmissions == False`). También se cierra si el examen termina antes que el usuario le de a finalizar.
  
  - `blocked`: diseñado para casos de plagio, el usuario profesor puede activarlo manualmente a un estudiante en tiempo real. Bloqueará la sesión del usuario y no le permitirá hacer ninguna revisión.


