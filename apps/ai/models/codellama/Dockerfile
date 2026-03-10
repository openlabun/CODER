FROM ollama/ollama:latest

EXPOSE 11434

COPY Modelfile /Modelfile

RUN ollama serve & \
    sleep 5 && \
    ollama pull codellama:7b && \
    ollama create codellama-custom -f /Modelfile

CMD ["serve"]