import os

RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
QUEUE_NAME = os.getenv("PYTHON_QUEUE", "python.queue")

API_BASE_URL = os.getenv("APIV2_BASE_URL", "http://api:8080")
API_TOKEN = os.getenv("WORKER_KEY", "secret-token")

DOCKER_IMAGE = "python:3.11"
EXECUTION_TIMEOUT = 5  # seconds
CONCURRENCY = int(os.getenv("WORKER_CONCURRENCY", "8"))

# Runner pool
PYTHON_RUNNERS = int(os.getenv("PYTHON_RUNNERS", "3"))
CPP_RUNNERS = int(os.getenv("CPP_RUNNERS", "3"))
JAVA_RUNNERS = int(os.getenv("JAVA_RUNNERS", "3"))

PYTHON_RUNNER_IMAGE = os.getenv("PYTHON_RUNNER_IMAGE", "coder-runner-python")
CPP_RUNNER_IMAGE = os.getenv("CPP_RUNNER_IMAGE", "coder-runner-cpp")
JAVA_RUNNER_IMAGE = os.getenv("JAVA_RUNNER_IMAGE", "coder-runner-java")

RUNNER_NETWORK = os.getenv("RUNNER_NETWORK", "judge_net")
RUNNER_NAME_PREFIX = os.getenv("RUNNER_NAME_PREFIX", "coder-runner")

RUNNER_HEALTH_CHECK_TIMEOUT = float(os.getenv("RUNNER_HEALTH_CHECK_TIMEOUT", "2.0"))
RUNNER_RESTART_TIMEOUT = float(os.getenv("RUNNER_RESTART_TIMEOUT", "120.0"))
RUNNER_REPAIR_INTERVAL = float(os.getenv("RUNNER_REPAIR_INTERVAL", "30.0"))