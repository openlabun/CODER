# Generador de Ejercicios

## Explicación del Tema

Se deberá integrar un LLM para generación de enunciados de retos, ejemplos y casos de prueba. El modelo deberá ejecutarse localmente.

## Stack Tecnológico

- Despliegue: Se cuenta con dos posibles herramientas para realizar esta función.
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
  
  
  
  Se plantea, se usar la herramienta de **Docker Models** para el proyecto debido a su mayor facilidad para la orquestación y mantenimiento del modelo. Adicionalmente la herramiento de Docker no necesita instalación de drivers para uso de la GPU, pues cuenta con soporte nativo. Sin embargo, se entiende que varios de los modelos considerados para esta funcionalidad no se encuentran disponibles para esta herramienta, por lo que la decisión se tomará teniendo en cuenta la dupla Herramienta + Modelo, donde herramienta será Docker Model en caso de estar disponible, y en caso contrario, será Ollama.
  
  
- Modelo:

| Model                              | Tool          | Strengths                                                                                                                   | Weaknesses                                                              | Ideal For                                             | Complexity to Run |
| ---------------------------------- | ------------- | --------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- | ----------------------------------------------------- | ----------------- |
| **CodeLlama (7B–70B)**             | Ollama        | State-of-the-art open-source for code, good instruction following, multi-language code support. High quality on benchmarks. | Larger variants need GPU; biggest models consume significant resources. | General coding assistant, exercise + test generation. | Medium–High       |
| **WizardCoder (7–15B)**            | Ollama        | Fine-tuned on code tasks with complex instruction training; strong on generation and reasoning.                             | Still heavier than some alternatives; less mainstream tooling.          | Deep code reasoning, elaborated tests, explanations.  | Medium–High       |
| **Mistral / Mixtral (7B / 8×7B)**  | Docker Models | Excellent balance of performance and efficiency; strong code + language.                                                    | Not originally code-specific (but fine-tuned variants exist).           | Lightweight, general purpose coding tasks.            | Low–Medium        |
| **DeepSeek Coder (6.7B+)**         | Ollama        | Strong multi-language performance, solid code generation & context handling.                                                | Mostly community support; less documentation.                           | Local Docker deployments with good code support.      | Medium            |
| **Qwen Coder (up to 72B)**         | Docker Models | Competitive performance on real developer coding tests; strong debugging & logic.                                           | Very large variants require powerful hardware.                          | High-quality generation/teacher scenarios.            | High              |
| **StarCoder / StarCoder2 (7–15B)** | Ollama        | Versatile code model trained on broad programming languages.                                                                | Benchmarks slightly below CodeLlama family.                             | Multi-language coding exercises and simple tests.     | Medium            |
| **Phi-3-mini / lighter models**    | Ollama        | Very efficient, runs on CPUs, cheap.                                                                                        | Less powerful at deep reasoning than larger models.                     | Lightweight coding assistant or supplement.           | Low               |

## Enlaces Relevantes

- [1] Notes on the margins, Rost Glukhov. [Docker Model Runner vs Ollama: Which to Choose? - Rost Glukhov | Personal site and technical blog](https://www.glukhov.org/llm-hosting/comparisons/docker-model-runner-vs-ollama-comparison/) 

- [2] 2024, E2E Networks. https://www.e2enetworks.com/blog/top-8-open-source-llms-for-coding?utm_source=chatgpt.com

- [3]

## Códigos de Prueba
