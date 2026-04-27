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
  
  - `10`: Primer Semestre
  
  - `20`: Intersemestral
  
  - `30`: Segundo Semestre

- Si no se ingresa un periodo el curso será accesible de manera indeterminada

## Módulo Exámenes

#### Exámenes (`Exam`)

- Un examen (`Exam`) contiene varios retos (`Challenge`) a través de un `ExamItem`. Entre los exámenes y los retos se manejará una relación muchos a muchos.

- Un examen (`Exam`) puede estar o no vinculado a un curso (`Course`), si no está vinculado, podrá ser visto por todos los usuarios estudiantes (según configuración de visibilidad).

- La configuración de visbilidad se manejará de la siguiente manera:
  
  - `public`: accesible como un reto (`Challenge`) disponible para todos los usuarios en la galería de retos
  
  - `course`: solo visible para el curso (`Course`) al que se encuentra vinculado el examen, si no está vinculado esta opción no deberá ser válida
  
  - `teachers`: solo visible para otros usuarios profesores en la galería de retos.
  
  - `private`: solo será visible para el usuario que lo creó.

- El examen podrá configurarse con una fecha de inicio (atributo `StartTime`) y con una fecha de cierre (atributo `EndTime`). Después de la fecha de cierre el examen no recibirá revisiones (`Submission`), a menos que se encuentre activada la configuración para aceptar respuestas tardías (`AllowLateSubmissions == True`).

- El examen podrá tener un tiempo límite (atributo `TimeLimit`) y un límite de intentos (atributo `TryLimit`). Si alguna revisión (`Submission`) se carga después del tiempo límite **no será recibida**, y si se excede el límite de intentos tampoco. Si el valor de estos atributos es `-1` entonces se tomará como **tiempo o intentos ilimitados**.

- Solo los exámenes (`Exam`) pueden resolverse. Los retos (`Challenge`) deben estar asociados a uno para poder resolverse.

#### Retos (`Challenge`)

- Son enunciados de retos de programación, los usuarios profesores pueden crearlos y alojarlos en su repositorio, estos podrán ser compartidos o editados según su configuración de guardado.

- Los retos (`Challenge`) pueden guardarse como:
  
  - `draft`: permiten modificación, no pueden ser vinculados a un examen.
  
  - `published`: pueden ser vinculados a un examen, y cualquier otro usuario profesor lo puede buscar y hacer copia (*fork*).
  
  - `private`: solo serán visibles en el repositorio para su dueño pero pueden ser vinculados a un examen.
  
  - `archived`: no permiten modificación, tampoco estarán visibles en el examen.

- Una vez publicado un reto **no podrá modificarse**, solo podrá archivarse. Y una vez archivado, puede volver a publicarse.

- Se puede hacer copia (*fork*) a un reto existente de otro usuario (si está disponible para su vista) o de uno del que el usuario es dueño.

- Cada reto deberá contar con un time límite de ejecución y una memoria límite para el worker. **(No podrá exceder determinados valores fijados en las variables en producción)**.

- Los retos contarán con `Templates` para el código, estos se podrán solicitar a la API para usar la plantilla por defecto o se podrá modificar para brindar soluciones predeterminadas, tener en cuenta que **si no se suministra un lenguaje** junto a su template, se considerará que el lenguaje **no está permitido** para resolver el reto.

- Los retos deberán contar con unas consideraciones de entrada/salida. Para ello deberán ingresar por cada variable los siguientes valores:
  
  - Nombre (`Name`)
  
  - Tipo (`Type`: `string`, `int`, `float`)
  
  - Valor esperado/ejemplo (`Value`)

#### Punto de Examen (`ExamItem`)

- Son asociaciones de muchos a muchos entre exámenes (`Exam`) y retos (`Challenge`)

- Solo se pueden crear a partir de retos (`Challenge`) que se encuentren en tu repositorio (visibilidad `private`) o que otros usuarios hayan publicado (visibilidad `published`).

- En caso de usar un reto de otro usuario se realizará *fork* automáticamente hacia tu propio repositorio.

- Por cada punto de examen (`ExamItem`) se deberá incluir valores para el orden del punto en el examen (atributo `Order`) y el valor del punto (atributo `Points`). El orden se calcula automáticamente cuando se crea el `ExamItem`.

- Solo el dueño de un examen (`Exam`) puede crear puntos (`ExamItem`), modificarlos o borrarlos en él.

- Al crear un nuevo punto (`ExamItem`) validará si ese reto (`Challenge`) no se encontraba ya en el examen (`Exam`). Si lo estaba arrojará error.

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

- Deberán siempre estar vinculadas a un examen (`Exam`), a una sesión (`Session`) y a un usuario (`User`).

- Estará vinculada a un Puntaje de Punto de Examen (`ExamItemScore`) para calificar el puntaje generado, esto solo se hará si el atributo `Scorable` está activo.

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

#### Puntaje de Examen (`ExamScore`)

- El puntaje de examen depende de la sesión del estudiante (relación 1 a 1).

- La estructura de este modelo será desnormalizada por temas de rendimiento, el puntaje estará asociado al examen al que pertenece y al usuario (esta información también la tiene la sesión).

- Este modelo también contará con un valor numérico del resultado del examen y el tiempo que se demoró en completar.

- Se creará al iniciar la sesión pero no se calculará el resultado hasta completar o expirar una sesión, si la sesión finaliza por bloqueo no será procesado el resultado.

#### Puntaje de Punto de Examen (`ExamItemScore`)

- Este modelo dependerá del Puntaje de Examen (relación 1 a muchos).

- También tendrá relación (1 a muchos) con el Punto de Examen (`ExamItem`) al que hace referencia.

- Se crearán todos al crear un Puntaje de Examen (`ExamScore`).

- Contará con la información del puntaje resultado de ese punto en específico, calculado sobre la base del atributo `Points` del modelo de Punto de Examen. Solo se modificará si alguna de las revisiones consigue un puntaje mayor.

- Llevará el contador de intentos, si este supera el límite permitido en el examen, se bloquearán todas las revisiones siguientes (`Submission`).

# 
