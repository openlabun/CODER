FROM ollama/ollama:latest

EXPOSE 11434

COPY Modelfile /Modelfile

RUN ollama serve & \
    sleep 5 && \
    ollama pull mistral:7b-instruct && \
    ollama create mistral-custom -f /Modelfile

CMD ["serve"]