FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    g++ \
    python3 \
    python3-pip \
    && pip3 install --no-cache-dir aiohttp==3.9.5 --break-system-packages \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY server.py .

EXPOSE 8080

CMD ["python3", "server.py"]
