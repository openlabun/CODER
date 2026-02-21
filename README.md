# CODER - Proyecto Final

## Introducción

En un contexto institucional la evaluación de algoritmos es un componente esencial en la formación de los ingenieros de sistemas, pero su implementación manual presenta limitaciones de escalabilidad, objetividad y rapidez en la retroalimentación. En contextos de educación digital y grupos numerosos, estas dificultades se acentúan, afectando la calidad del proceso evaluativo. Este proyecto propone el diseño y desarrollo de **CODER**, una plataforma web de evaluación automática de algoritmos que integra ejecución segura de código, retroalimentación inmediata, calificación de pruebas e integridad académica para fortalecer los procesos de enseñanza y evaluación.

## Planteamiento del Problema

Los procesos tradicionales de evaluación de algoritmos dependen en gran medida de la revisión manual del código por parte del docente, lo que dificulta la aplicación de criterios claros, incrementa la carga operativa para el docente y retrasa la retroalimentación al estudiante. En cursos donde hay una gran cantidad de estudiantes, esta situación impacta negativamente la calidad del aprendizaje y limita la posibilidad de realizar evaluaciones frecuentes y formativas, dado que la alta carga de trabajo para el docente lo lleva a reducir la cantidad de evaluaciones o en su defecto, entregar de manera retrasada los resultados y retroalimentación.

## Restricciones y supuestos de diseño

#### Escalabilidad y Despliegue

- La plataforma deberá ejecutarse sobre una infraestructura de contenedores que permitan su rápida escalabilidad

- El sistema deberá tolerar picos de carga durante evaluaciones masivas.

- El sandbox de ejecución tendrá límites estrictos de CPU, memoria y tiempo.

- El modelo de IA para generación de contenido académico se ejecutará de forma local. Deberá usarse un modelo en versión reducida.

#### Inteligencia Artificial

- El sistema de detección de fraudes será utilizado únicamente como sistema de alerta, no de sanción automática.

- Todo contenido generado por IA deberá ser validado manualmente por el docente.

- La IA deberá realizar y ejecutar test de prueba y casos de uso para los ejercicios propuestos.

#### Restricciones de Usuario

- Se asumirá conectividad a internet estable en los entornos institucionales de evaluación.

- El sistema deberá soportar múltiples lenguajes de programación definidos previamente.

- El acceso al sistema estará controlado mediante autenticación y roles.

## Alcance

El proyecto comprende el diseño, desarrollo, integración y despliegue en producción de la plataforma web **CODER**, orientada a la evaluación automática de algoritmos en cursos de Ingeniería de Sistemas. La solución permitirá a los docentes crear y administrar retos de programación, definir casos de prueba, configurar parciales y gestionar cursos académicos, mientras que los estudiantes podrán enviar soluciones, recibir retroalimentación automática y consultar resultados y métricas de desempeño.

La plataforma incluirá un motor de evaluación basado en colas de tareas y *workers* especializados por lenguaje, con capacidad de escalamiento dinámico. Asimismo, integrará módulos de detección de plagio, control de intentos, ventanas temporales de evaluación y auditoría completa de actividades. Adicionalmente, se incorporará un componente de inteligencia artificial ejecutado localmente para asistir a los docentes en la generación de contenido académico, sujeto siempre a validación humana.