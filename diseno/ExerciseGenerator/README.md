# Generador de Ejercicios

## Explicación del Tema

El sistema integrará un **Modelo de Lenguaje de Gran Escala (LLM)** ejecutado localmente con el objetivo de generar automáticamente enunciados de ejercicios de programación, ejemplos explicativos y casos de prueba asociados. La ejecución local del modelo garantiza independencia de servicios externos, menor latencia en la generación de contenido y mayor control sobre la privacidad y el uso de los datos.

El generador de ejercicios integrará un LLM ejecutado de manera local como componente central del sistema. Este modelo será responsable de producir, a partir de un conjunto de parámetros de entrada, un ejercicio completo compuesto por:

1. Enunciado del problema.

2. Especificación de entradas y salidas.

3. Restricciones formales.

4. Ejemplos explicativos.

5. Casos de prueba automatizados.

## Stack Tecnológico

- Despliegue: Se plantea el uso de Docker, dado que permite la contenerización del modelo, además, que facilita su escalamiento a través de Kubernetes.

- Soporte: Se cuenta con dos posibles herramientas para realizar esta función.
  
  - Docker Model
  - Ollama
  
  En el siguiente cuadro comparativo extraído de la fuente [1] se puede evidenciar las ventajas de usar una herramienta u otra:
  
  | Feature                | Docker Model Runner                                            | Ollama                                                         |
  | ---------------------- | -------------------------------------------------------------- | -------------------------------------------------------------- |
  | **Installation**       | Docker Desktop AI tab or `docker-model-plugin`                 | Single command: `curl \| sh`                                   |
  | **Command Style**      | `docker model pull/run/package`                                | `ollama pull/run/list`                                         |
  | **Model Format**       | GGUF (OCI Artifacts)                                           | GGUF (native)                                                  |
  | **Model Distribution** | Docker Hub, OCI registries                                     | Ollama registry                                                |
  | **GPU Setup**          | Automatic (simpler than traditional Docker)                    | Automatic                                                      |
  | **API**                | OpenAI-compatible                                              | OpenAI-compatible                                              |
  | **Docker Integration** | Native (is Docker)                                             | Runs in Docker if needed                                       |
  | **Compose Support**    | Native                                                         | Via Docker image                                               |
  | **Learning Curve**     | Low (for Docker users)                                         | Lowest (for everyone)                                          |
  | **Ecosystem Partners** | Google, Hugging Face, VMware                                   | LangChain, CrewAI, Open WebUI                                  |
  | **Best For**           | Docker-native workflows                                        | Standalone simplicity                                          |
  | **Inference Speed**    | 20-30 tokens/sec on CPU and 50-80 tokens/sec on mid-range GPUs | 20-30 tokens/sec on CPU and 50-80 tokens/sec on mid-range GPUs |
  | **Memory Usage**       | 4-6 GB RAM. Container overhead is minimal                      | 4-6 GB RAM                                                     |
  | **Startup Time**       | Container adds ~1 sec, model loading 2-5 secs                  | Model loading 2-5 secs                                         |
  
  Teniendo en cuenta el cuadro comparativo, en esta línea ambas opciones resultan viables para el desarrollo del proyecto. 
  
  Dado que ambas soluciones cuentan con soporte para contenerización con Docker, la única diferencia de Docker Modelos con Ollama para nuestro caso de uso sería que en este primero no se requerie instalación de drivers para uso de la GPU, debido a que cuenta con soporte nativo desde Docker. Sin embargo, Docker Models cuenta con una gran desventaja y es su reducido catálogo de modelos, sobretodo aquellos que requerimos para nuestro proyecto, por lo que la decisión se tomará teniendo en cuenta la dupla Herramienta + Modelo, donde herramienta será Docker Model en caso de estar disponible, y en caso contrario, será Ollama.

- Modelo:

| Modelo                  | Parámetros | Contexto (tokens) | RAM aprox (Q4) | Fortaleza principal              |
| ----------------------- | ---------- | ----------------- | -------------- | -------------------------------- |
| **DeepSeek Coder 6.7B** | 6.7B       | ~16K              | 6–8 GB         | Excelente en código estructurado |
| **Qwen2.5-Coder 7B**    | 7B         | ~32K              | 8–10 GB        | Mejor razonamiento largo         |
| **Code Llama 7B**       | 7B         | ~16K              | 8–10 GB        | Generación limpia de código      |
| **Mistral 7B Instruct** | 7B         | ~8K               | 7–9 GB         | Buen razonamiento general        |
| **Phi-3 Mini**          | ~3.8B      | ~8K–16K           | 4–6 GB         | Muy eficiente                    |

## Enlaces Relevantes

- [1] Notes on the margins, Rost Glukhov. [Docker Model Runner vs Ollama: Which to Choose? - Rost Glukhov | Personal site and technical blog](https://www.glukhov.org/llm-hosting/comparisons/docker-model-runner-vs-ollama-comparison/) 

- [2] 2024, E2E Networks. https://www.e2enetworks.com/blog/top-8-open-source-llms-for-coding?utm_source=chatgpt.com

- [3]

## Códigos de Prueba
