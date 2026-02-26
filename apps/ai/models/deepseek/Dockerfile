FROM ollama/ollama:latest

# Expose Ollama port
EXPOSE 11434

# Copy custom model definition
COPY Modelfile /Modelfile

# Start Ollama and create model
RUN ollama serve & \
    sleep 5 && \
    ollama create deepseek-custom -f /Modelfile

CMD ["serve"]