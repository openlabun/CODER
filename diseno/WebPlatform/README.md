# Plataforma Web

## Explicación del Tema

Se deberá diseñar, desarrollar y desplegar en producción una plataforma web que permita la evaluación automática de algoritmos para fines educativos. Deberá incluir:

- Un motor de procesamiento basado en colas de tareas y workers especializados por lenguaje de programación.

- Medición de Desempeño

- Calificación de resultados

- Retroalimentación automática

- Rankings por actividad

## Enlaces Relevantes

- **Sistema de CRUD de Usuarios y JWT:** [Documentación ROBLE](https://roble.openlab.uninorte.edu.co/docs)

- **Proyecto Anterior desarrollado en clase de Backend:** [GitHub - DerekPz/juez-online](https://github.com/DerekPz/juez-online.git)

## Código de Prueba

- Backend: apps\api

- Frontend: apps\web

el trabajo seleccionado ya presenta una base funcional adecuada para la implementacion de los requerimientos nescesarios para el desarrollo del proyecto, como los runners en Python, Java, C++ y Javascript, la cola en redis que llama al worker y como tal, es una base solida para la realizacion del proyecto.

## Analisis de la plataforma (proyecto) anteriror
La plataforma analizada corresponde a un sistema de evaluación automática de algoritmos, comúnmente conocido como juez online. Este tipo de aplicaciones permite a los usuarios enviar soluciones en distintos lenguajes de programación para resolver retos algorítmicos, mientras el sistema se encarga de compilar el código, ejecutarlo en un entorno controlado y verificar los resultados contra casos de prueba predefinidos. El proyecto está diseñado para ser escalable, seguro y modular, integrando diversas tecnologías modernas para gestionar tanto la lógica del sistema como la ejecución de código de forma aislada.

El backend del sistema está desarrollado utilizando NestJS, un framework progresivo basado en Node.js y escrito en TypeScript. Esta herramienta facilita la construcción de aplicaciones robustas mediante una estructura modular que permite dividir el sistema en múltiples componentes funcionales. En este proyecto, NestJS se utiliza para implementar varios módulos esenciales como autenticación, retos, envíos de soluciones (submissions), ejecución de código (runners), cursos, calificación, clasificación (leaderboard), observabilidad y asistencia creativa. Gracias al uso de controladores REST y métodos HTTP estándar (GET, POST, PUT y DELETE), el backend puede exponer una API clara y organizada que permite la comunicación con la interfaz de usuario y con otros servicios del sistema. Además, NestJS facilita la integración con sistemas de colas de trabajo que permiten manejar de forma eficiente el procesamiento asíncrono de las soluciones enviadas por los estudiantes.

Desde el punto de vista arquitectónico, el proyecto sigue el enfoque de Clean Architecture (Arquitectura Limpia). Este modelo organiza el código en capas que separan claramente las responsabilidades del sistema, permitiendo que la lógica de negocio sea independiente de tecnologías externas como frameworks, bases de datos o interfaces de usuario. En la capa más interna se encuentran las entidades o dominio, donde se definen los objetos principales del sistema, como los usuarios, los retos de programación, las soluciones enviadas y las clasificaciones. A partir de estas entidades se construyen los casos de uso, que representan las acciones que el sistema puede ejecutar, por ejemplo crear un reto, enviar una solución o procesar un envío de código.

Por encima de estas capas se encuentran las interfaces o adaptadores, que funcionan como puentes entre la lógica interna y el mundo exterior. Aquí se ubican los controladores REST que reciben las peticiones de la API, los repositorios encargados de comunicarse con la base de datos y los adaptadores que interactúan con servicios externos como sistemas de cache o colas de trabajo. Finalmente, en la capa más externa se encuentran los drivers, donde residen las herramientas concretas que utiliza el sistema, como el motor de base de datos, el sistema de mensajería o los contenedores de ejecución.

En cuanto al stack tecnológico, el sistema se apoya en una combinación de herramientas ampliamente utilizadas en aplicaciones de alto rendimiento. La gestión de datos se realiza mediante PostgreSQL, una base de datos relacional que almacena información como usuarios, retos, resultados de evaluaciones y estadísticas del sistema. Para el manejo de tareas asíncronas y colas de procesamiento se utiliza Redis, que permite distribuir las evaluaciones de código entre distintos workers de forma eficiente.

La autenticación y autorización de los usuarios se implementa mediante JSON Web Token, lo cual permite gestionar sesiones seguras sin necesidad de mantener estados en el servidor. Este mecanismo también permite diferenciar distintos roles dentro de la plataforma, como estudiantes y administradores o profesores, garantizando que cada tipo de usuario tenga acceso únicamente a las funcionalidades correspondientes.

Uno de los componentes más importantes del sistema es el sandbox o entorno de ejecución seguro. Para garantizar que el código enviado por los usuarios no represente un riesgo para la infraestructura, el sistema utiliza Docker para ejecutar cada solución dentro de contenedores aislados. Estos contenedores se crean de manera efímera y cuentan con restricciones estrictas, como la desactivación del acceso a la red y límites definidos de CPU y memoria. Una vez que el código termina su ejecución, el contenedor se elimina automáticamente, evitando que cualquier proceso permanezca activo en el sistema anfitrión.

El flujo de procesamiento de soluciones funciona de manera asíncrona. Cuando un usuario envía una solución desde la interfaz del sistema (desarrollada con React y construida con Vite), la API del backend recibe el código y lo registra como una nueva submission. Posteriormente, esta tarea se coloca en una cola gestionada por Redis. Un worker especializado en el lenguaje correspondiente —como Python, Node.js, C++ o Java— toma la tarea de la cola y ejecuta el código dentro de un contenedor Docker. Durante esta ejecución, el sistema compila el programa si es necesario y lo ejecuta contra una serie de casos de prueba ocultos definidos por archivos de entrada y salida. Con base en la comparación de resultados, el sistema determina el estado de la solución, generando veredictos como ACCEPTED, WRONG ANSWER o TIME LIMIT EXCEEDED. Finalmente, estos resultados son almacenados en la base de datos y presentados al usuario en la plataforma.

El despliegue del sistema se realiza mediante Docker Compose, lo cual permite levantar de forma coordinada todos los servicios necesarios para el funcionamiento de la plataforma, incluyendo la API, la base de datos, el sistema de cache y los workers de procesamiento. Esta estrategia también facilita la escalabilidad del sistema, ya que es posible aumentar el número de workers para determinados lenguajes cuando la carga de evaluaciones lo requiera.

Por último, el sistema incorpora mecanismos básicos de observabilidad y monitoreo que permiten analizar el funcionamiento de la plataforma. Esto incluye registros estructurados en formato JSON y métricas que permiten medir aspectos como la cantidad de envíos procesados, los tiempos promedio de ejecución de los programas y los posibles errores internos del sistema. Gracias a esta información, los administradores pueden rastrear el flujo completo de una evaluación, desde el momento en que un estudiante envía su código hasta que el resultado final es generado y almacenado.

En conjunto, la plataforma combina una arquitectura moderna, tecnologías robustas y prácticas de seguridad adecuadas para construir un sistema confiable de evaluación automática de algoritmos, capaz de manejar múltiples usuarios, diferentes lenguajes de programación y un alto volumen de ejecuciones de código de forma segura y eficiente.
