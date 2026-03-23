import os

RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
QUEUE_NAME = os.getenv("QUEUE_NAME", "python.queue")

API_BASE_URL = os.getenv("API_BASE_URL", "http://api:8080")
API_TOKEN = os.getenv("API_TOKEN", "secret-token")

DOCKER_IMAGE = "python:3.11"
EXECUTION_TIMEOUT = 5  # seconds