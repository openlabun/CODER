FROM ollama/ollama:latest

EXPOSE 11434

COPY Modelfile /Modelfile

RUN ollama serve & \
    sleep 5 && \
    ollama pull qwen2.5-coder:7b && \
    ollama create qwen-custom -f /Modelfile

CMD ["serve"]