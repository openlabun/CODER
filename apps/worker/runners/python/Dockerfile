FROM python:3.12-slim

RUN pip install --no-cache-dir aiohttp==3.9.5

WORKDIR /app

COPY server.py .

EXPOSE 8080

CMD ["python", "server.py"]
