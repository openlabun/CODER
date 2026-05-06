import asyncio
import json

import aio_pika

from config import RABBITMQ_URL, QUEUE_NAME, CONCURRENCY
from app.domain.mapper import MapSubmissionResult
from app.domain.errors import RetryableSubmissionUpdateError, PermanentSubmissionUpdateError


class RabbitMQConsumer:

    def __init__(self, handler):
        self.handler = handler

    async def start(self):
        connection = await aio_pika.connect_robust(RABBITMQ_URL)
        channel = await connection.channel()

        await channel.set_qos(prefetch_count=CONCURRENCY)
        queue = await channel.declare_queue(QUEUE_NAME, durable=True)

        print(f"[Worker] Started. Listening on queue: {QUEUE_NAME} (concurrency={CONCURRENCY})", flush=True)

        async with queue.iterator() as queue_iter:
            async for message in queue_iter:
                task = asyncio.create_task(self._handle_message(message))
                task.add_done_callback(self._on_task_done)

    async def _handle_message(self, message: aio_pika.IncomingMessage):
        try:
            print(f"Received message: {message.body}", flush=True)
            data = json.loads(message.body)
            submission = MapSubmissionResult(data)
            await self.handler(submission)
            await message.ack()
        except PermanentSubmissionUpdateError as e:
            print(f"Permanent error processing message: {e}. Dropping message.", flush=True)
            await message.ack()
        except Exception as e:
            print(f"Error processing message: {e}. Requeueing.", flush=True)
            await message.nack(requeue=True)

    def _on_task_done(self, task: asyncio.Task):
        if not task.cancelled() and task.exception():
            print(f"[Worker] Unhandled task error: {task.exception()}", flush=True)
