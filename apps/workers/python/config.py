import os

RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
QUEUE_NAME = os.getenv("PYTHON_QUEUE", "python.queue")

API_BASE_URL = os.getenv("APIV2_BASE_URL", "http://api:8080")
API_TOKEN = os.getenv("WORKER_KEY", "secret-token")

DOCKER_IMAGE = "python:3.11"
EXECUTION_TIMEOUT = 5  # seconds