# CODER - Uninorte Online Judge

## Resumen ejecutivo

### Contexto y problemática

En un contexto institucional, la evaluación de algoritmos es un componente esencial en la formación de ingenieros de sistemas. Sin embargo, los procesos tradicionales dependen en gran medida de la revisión manual del código por parte del docente, lo que genera **sobrecarga operativa significativa**, **retroalimentación lenta** a los estudiantes y **dificultades para realizar evaluaciones frecuentes y masivas**. En contextos de educación digital y grupos numerosos, estas limitaciones se acentúan, impactando negativamente la calidad del aprendizaje y la posibilidad de llevar a cabo evaluaciones formativas de manera oportuna.

### Solución propuesta

**CODER** es una plataforma web tipo *online judge* desarrollada específicamente para el Departamento de Ingeniería de Sistemas y Computación de la Universidad del Norte. La solución **automatiza la creación, ejecución y calificación de ejercicios, talleres y exámenes de programación**, permitiendo a docentes definir retos con casos de prueba predefinidos y a estudiantes enviar soluciones que se evalúan automáticamente en tiempo real.

### Arquitectura y escalabilidad

La plataforma está construida bajo una **arquitectura desacoplada y escalable**:
- **Backend robusto** en Go (Fiber) que implementa lógica de negocio bajo arquitectura hexagonal
- **Cola de tareas** con RabbitMQ que desacopla el procesamiento de evaluaciones
- **Pool de runners aislados** por lenguaje (Python, Java, C++) que ejecutan código en sandboxes seguros
- Capacidad de procesar **hasta 25 envíos concurrentes** sin degradación del servicio
- **Persistencia externa** mediante Roble, garantizando integridad de datos
- **Integración de IA generativa** (Google Gemini) para asistir a docentes en la creación de contenido

### Alcance y características

La plataforma cubre el ciclo completo de evaluación:
- **Para docentes:** Crear y administrar cursos, diseñar retos y exámenes, configurar límites de tiempo e intentos, visualizar resultados agregados por estudiante y por ejercicio
- **Para estudiantes:** Resolver retos de forma autónoma, enviar soluciones a exámenes con retroalimentación automática, consultar históricos de intentos y métricas de desempeño
- **Funcionalidades transversales:** Autenticación segura con JWT, control granular de permisos, análisis educativo mediante dashboards y leaderboards, auditoría completa de actividades

### Valor y beneficios

El proyecto habilita **evaluaciones frecuentes y formativas a gran escala** manteniendo:
- ✅ **Retroalimentación inmediata** — Los estudiantes reciben resultados al instante
- ✅ **Escalabilidad horizontal** — Soporta picos de carga durante evaluaciones masivas
- ✅ **Seguridad y aislamiento** — Código ejecutado en entornos contenedorizados sin riesgo para el servidor
- ✅ **Integridad académica** — Registro completo de intentos, tiempos y resultados
- ✅ **Reducción de carga docente** — Libera al docente de revisión manual y permite enfocarse en análisis de conceptos débiles
- ✅ **Experiencia de usuario fluida** — Interfaz intuitiva con editor de código integrado (Monaco Editor)

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