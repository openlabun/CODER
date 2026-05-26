# CODER - Uninorte Online Judge

## Resumen ejecutivo

CODER es una plataforma web tipo online judge desarrollada para el Departamento de Ingenieria de Sistemas y Computacion de la Universidad del Norte, que surge como respuesta a la necesidad de mejorar los procesos de evaluacion de programacion en contextos con alta carga academica. En el esquema tradicional, la revision manual del codigo incrementa la carga operativa del docente, retrasa la retroalimentacion al estudiante y dificulta la realizacion de evaluaciones frecuentes y formativas en grupos numerosos; esta situacion impacta la oportunidad del aprendizaje y la consistencia de los criterios de calificacion.

La solucion propuesta automatiza la creacion, ejecucion y calificacion de ejercicios, talleres y examenes de programacion mediante una arquitectura desacoplada y escalable: frontend web para docentes y estudiantes, backend en Go con arquitectura hexagonal, cola de mensajeria con RabbitMQ, worker de ejecucion en Python y runners aislados por lenguaje para ejecutar codigo en entornos controlados. La plataforma permite gestionar cursos, retos, casos de prueba, sesiones e intentos de evaluacion, y entregar resultados de forma automatica y trazable.

Su alcance cubre el ciclo principal de evaluacion academica, desde la definicion de actividades por parte del docente hasta el envio de soluciones y consulta de resultados por parte del estudiante, incluyendo mecanismos de autenticacion, control de permisos y analitica de desempeno. El valor del proyecto radica en habilitar evaluaciones a mayor escala con retroalimentacion mas oportuna, reduciendo la carga manual docente y fortaleciendo la calidad del proceso de ensenanza y aprendizaje en asignaturas de programacion.

## Documentación del repositorio

| Documento                          | Descripción                                  |
| ---------------------------------- | -------------------------------------------- |
| [Informe.md](./Informe.md)         | Documento principal del proyecto             |
| [Instalación.md](./Instalacion.md) | Guía de instalación, desarrollo y despliegue |
| [Desarrollo.md](./Desarrollo.md)   | Detalles técnicos del desarrollo             |

## Estudiantes

| Nombre              | GitHub                                                       |
| ------------------- | ------------------------------------------------------------ |
| Jhon Jimenez        | [@Jhon-jimenez](https://github.com/Jhon-Jimenez)             |
| Miguel Julio        | [@mjulioa-m](https://github.com/mjulioa-m)                   |
| Sebastián Maldonado | [@SebastianMaldonado](https://github.com/SebastianMaldonado) |
| Raúl Morales        | [@Azkalagon](https://github.com/Azkalagon)                   |

## Tutores

- Augusto Salazar
- Daniel Romero