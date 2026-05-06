import asyncio
from app.application.process_submission import ProcessSubmission
from app.adapters.driven.docker_executor import DockerExecutor
from app.adapters.driven.http_api_client import HTTPAPIClient
from app.adapters.driving.rabbitmq_consumer import RabbitMQConsumer


async def main():
    print("[Worker] Starting Python worker...", flush=True)
    executor = DockerExecutor()
    await executor.ensure_image_cached()
    api_client = HTTPAPIClient()

    use_case = ProcessSubmission(executor, api_client)

    consumer = RabbitMQConsumer(handler=use_case.execute)
    await consumer.start()


if __name__ == "__main__":
    asyncio.run(main())