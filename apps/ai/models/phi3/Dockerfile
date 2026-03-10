FROM ollama/ollama:latest

EXPOSE 11434

COPY Modelfile /Modelfile

RUN ollama serve & \
    sleep 5 && \
    ollama pull phi3:mini && \
    ollama create phi3-custom -f /Modelfile

CMD ["serve"]