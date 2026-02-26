# Iniciar contenedores LLM

> Ejecutar desde `apps/ai/`. Se recomienda levantar **un modelo a la vez** si la máquina no tiene suficiente RAM/VRAM.

| Modelo | Puerto host |
|---|---|
| CodeLlama 7B | 11434 |
| DeepSeek Coder 6.7B | 11435 |
| Mistral 7B Instruct | 11436 |
| Phi-3 Mini | 11437 |
| Qwen 2.5 Coder 7B | 11438 |

## Comandos

```bash
docker build -t juez_llm_codellama ./apps/ai/models/codellama && docker run -d -p 11434:11434 --name llm_codellama juez_llm_codellama
```

```bash
docker build -t juez_llm_deepseek ./apps/ai/models/deepseek && docker run -d -p 11435:11434 --name llm_deepseek juez_llm_deepseek
```

```bash
docker build -t juez_llm_mistral ./apps/ai/models/mistral && docker run -d -p 11436:11434 --name llm_mistral juez_llm_mistral
```

```bash
docker build -t juez_llm_phi3 ./apps/ai/models/phi3 && docker run -d -p 11437:11434 --name llm_phi3 juez_llm_phi3
```

```bash
docker build -t juez_llm_qwen ./apps/ai/models/qwen && docker run -d -p 11438:11434 --name llm_qwen juez_llm_qwen
```

## Notas

- Cada comando construye la imagen y luego levanta el contenedor en segundo plano (`-d`).
- Para ver los logs de un contenedor: `docker logs -f llm_<modelo>` (ej. `docker logs -f llm_codellama`).
- Para detener un contenedor: `docker stop llm_<modelo>`.
- Para eliminar un contenedor: `docker rm llm_<modelo>`.
