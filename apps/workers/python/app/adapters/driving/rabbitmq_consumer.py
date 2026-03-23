import json
import pika
from config import RABBITMQ_URL, QUEUE_NAME
from app.domain.models import SubmissionResult
from app.domain.mapper import MapSubmissionResult

class RabbitMQConsumer:

    def __init__(self, handler):
        self.handler = handler

    def start(self):
        connection = pika.BlockingConnection(pika.URLParameters(RABBITMQ_URL))
        channel = connection.channel()

        channel.queue_declare(queue=QUEUE_NAME, durable=True)

        channel.basic_qos(prefetch_count=1)

        def callback(ch, method, properties, body):
            try:
                data = json.loads(body)

                submission = MapSubmissionResult(data)
        
                self.handler(submission)

                ch.basic_ack(delivery_tag=method.delivery_tag)

            except Exception as e:
                print(f"Error: {e}")
                ch.basic_nack(delivery_tag=method.delivery_tag, requeue=True)

        channel.basic_consume(queue=QUEUE_NAME, on_message_callback=callback)

        print("Worker started...")
        channel.start_consuming()