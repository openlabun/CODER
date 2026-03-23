import json
import time

import pika
from pika.exceptions import AMQPConnectionError

from config import RABBITMQ_URL, QUEUE_NAME
from app.domain.mapper import MapSubmissionResult
from app.domain.errors import RetryableSubmissionUpdateError, PermanentSubmissionUpdateError

class RabbitMQConsumer:

    def __init__(self, handler):
        self.handler = handler

    def _connect_with_retry(self, max_attempts=20, delay_seconds=2):
        last_error = None

        for attempt in range(1, max_attempts + 1):
            try:
                print(f"[RabbitMQ] Connecting to {RABBITMQ_URL} (attempt {attempt}/{max_attempts})", flush=True)
                return pika.BlockingConnection(pika.URLParameters(RABBITMQ_URL))
            except AMQPConnectionError as err:
                last_error = err
                if attempt < max_attempts:
                    time.sleep(delay_seconds)

        raise AMQPConnectionError(
            f"Could not connect to RabbitMQ after {max_attempts} attempts: {last_error}"
        )

    def start(self):
        connection = self._connect_with_retry()
        channel = connection.channel()

        channel.queue_declare(queue=QUEUE_NAME, durable=True)

        channel.basic_qos(prefetch_count=1)

        def callback(ch, method, properties, body):
            try:
                print(f"Received message: {body}", flush=True)
                data = json.loads(body)

                submission = MapSubmissionResult(data)
        
                self.handler(submission)

                ch.basic_ack(delivery_tag=method.delivery_tag)

            except PermanentSubmissionUpdateError as e:
                print(f"Permanent error processing message: {e}. Dropping message.", flush=True)
                ch.basic_ack(delivery_tag=method.delivery_tag)

            except RetryableSubmissionUpdateError as e:
                print(f"Retryable error processing message: {e}. Requeueing message.", flush=True)
                ch.basic_nack(delivery_tag=method.delivery_tag, requeue=True)

            except Exception as e:
                print(f"Error: {e}", flush=True)
                ch.basic_nack(delivery_tag=method.delivery_tag, requeue=True)

        channel.basic_consume(queue=QUEUE_NAME, on_message_callback=callback)

        print(f"Worker started. Listening queue: {QUEUE_NAME}", flush=True)
        channel.start_consuming()