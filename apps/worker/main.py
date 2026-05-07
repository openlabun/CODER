import asyncio
import signal

from app.application.process_submission import ProcessSubmission
from app.adapters.driven.container_pool.adapter import ContainerPoolAdapter
from app.adapters.driven.http_api_client import HTTPAPIClient
from app.adapters.driving.rabbitmq_consumer import RabbitMQConsumer
from config import (
    PYTHON_RUNNERS, PYTHON_RUNNER_IMAGE,
    RUNNER_NETWORK, RUNNER_NAME_PREFIX,
    RUNNER_HEALTH_CHECK_TIMEOUT, RUNNER_RESTART_TIMEOUT, RUNNER_REPAIR_INTERVAL,
)


async def main():
    loop = asyncio.get_running_loop()

    pool = ContainerPoolAdapter(
        language="python",
        image=PYTHON_RUNNER_IMAGE,
        count=PYTHON_RUNNERS,
        network=RUNNER_NETWORK,
        name_prefix=RUNNER_NAME_PREFIX,
        health_check_timeout=RUNNER_HEALTH_CHECK_TIMEOUT,
        restart_timeout=RUNNER_RESTART_TIMEOUT,
        repair_interval=RUNNER_REPAIR_INTERVAL,
    )
    await pool.startup()

    api_client = HTTPAPIClient()
    use_case = ProcessSubmission(pool, api_client)
    consumer = RabbitMQConsumer(handler=use_case.execute)

    stop_event = asyncio.Event()
    loop.add_signal_handler(signal.SIGTERM, stop_event.set)
    loop.add_signal_handler(signal.SIGINT, stop_event.set)

    consumer_task = asyncio.create_task(consumer.start())
    shutdown_task = asyncio.create_task(stop_event.wait())

    print("[Worker] Running. Send SIGTERM or SIGINT to shut down.", flush=True)
    await asyncio.wait({consumer_task, shutdown_task}, return_when=asyncio.FIRST_COMPLETED)
    shutdown_task.cancel()

    print("[Worker] Shutting down...", flush=True)
    await consumer.stop()
    await consumer.wait_for_tasks(timeout=30.0)

    consumer_task.cancel()
    try:
        await consumer_task
    except (asyncio.CancelledError, Exception):
        pass

    await pool.shutdown()
    print("[Worker] Shutdown complete.", flush=True)


if __name__ == "__main__":
    asyncio.run(main())
