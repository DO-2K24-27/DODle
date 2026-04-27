#!/usr/bin/env python3
"""
RabbitMQ message producer for KEDA scaling demo.
Sends messages to the 'frontend-tasks' queue every 5 seconds.
"""

import pika
import time
import json
import sys
import os
from datetime import datetime

def send_messages():
    """Connect to RabbitMQ and send messages every 5 seconds."""
    
    # RabbitMQ connection parameters
    rabbitmq_host = os.getenv('RABBITMQ_HOST', 'dodle-rabbitmq.polytech.svc.cluster.local')
    rabbitmq_port = int(os.getenv('RABBITMQ_PORT', '5672'))
    rabbitmq_user = os.getenv('RABBITMQ_USER', 'guest')
    rabbitmq_pass = os.getenv('RABBITMQ_PASS', 'password')
    rabbitmq_vhost = os.getenv('RABBITMQ_VHOST', '/')

    credentials = pika.PlainCredentials(rabbitmq_user, rabbitmq_pass)
    parameters = pika.ConnectionParameters(
        host=rabbitmq_host,
        port=rabbitmq_port,
        virtual_host=rabbitmq_vhost,
        credentials=credentials,
        connection_attempts=3,
        retry_delay=2
    )
    
    try:
        # Connect to RabbitMQ
        connection = pika.BlockingConnection(parameters)
        channel = connection.channel()
        
        # Declare the queue (idempotent)
        queue_name = 'frontend-tasks'
        channel.queue_declare(queue=queue_name, durable=True)
        
        print(f"✓ Connected to RabbitMQ on {rabbitmq_host}:{rabbitmq_port}")
        print(f"✓ Publishing to queue: {queue_name}")
        print("✓ Sending messages every 2 seconds (Ctrl+C to stop)\n")
        
        message_counter = 1
        
        while True:
            # Create message payload
            message = {
                'id': message_counter,
                'timestamp': datetime.now().isoformat(),
                'task': 'frontend-processing-task',
                'data': f'Task #{message_counter}'
            }
            
            # Publish message
            channel.basic_publish(
                exchange='',
                routing_key=queue_name,
                body=json.dumps(message),
                properties=pika.BasicProperties(
                    delivery_mode=pika.spec.PERSISTENT_DELIVERY_MODE
                )
            )
            
            print(f"[{datetime.now().strftime('%H:%M:%S')}] Sent message #{message_counter}")
            message_counter += 1
            
            # Wait 2 seconds
            time.sleep(2)
    
    except pika.exceptions.AMQPConnectionError as e:
        print(f"✗ Connection failed: {e}")
        print("Make sure RabbitMQ pod is running: kubectl get pods -n polytech | grep rabbitmq")
        sys.exit(1)
    except KeyboardInterrupt:
        print("\n✓ Stopped sending messages")
        connection.close()
        sys.exit(0)
    except Exception as e:
        print(f"✗ Error: {e}")
        sys.exit(1)

if __name__ == '__main__':
    send_messages()
