import random
import string
import time
from datetime import datetime

from kafka import KafkaProducer
import json


def random_string(length):
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))


def random_phone():
    return '+972' + ''.join(random.choices(string.digits, k=7))


def random_email():
    return random_string(5).lower() + '@' + random_string(5).lower() + '.com'


def random_date_iso():
    now = datetime.utcnow()
    return now.isoformat() + 'Z'


def generate_item(track_number):
    return {
        "chrt_id": random.randint(1000000, 9999999),
        "track_number": track_number,
        "price": random.randint(100, 500),
        "rid": random_string(16),
        "name": random_string(10),
        "sale": random.randint(0, 25),
        "size": str(random.randint(5, 100)),
        "total_price": random.randint(100, 5000),
        "nm_id": random.randint(1000000, 9999999),
        "brand": random_string(20),
        "status": random.randint(-100, 250)
    }


def generate_payment():
    return {
        "transaction": random_string(16),
        "request_id": random_string(20),
        "currency": "RUB",
        "provider": "wbpay",
        "amount": random.randint(100, 5000),
        "payment_dt": int(time.time()),
        "bank": random_string(10),
        "delivery_cost": random.randint(0, 2000),
        "goods_total": random.randint(100, 1000),
        "custom_fee": random.randint(0, 50)
    }


def generate_delivery():
    return {
        "name": f"{random_string(5)} {random_string(10)}",
        "phone": random_phone(),
        "zip": ''.join(random.choices(string.digits, k=7)),
        "city": random_string(15),
        "address": f"{random_string(10)} {random.randint(1, 20)}",
        "region": random_string(8),
        "email": random_email()
    }


def generate_message():
    track_number = 'WBILM' + random_string(5).upper() + 'TRACK'
    return {
        "order_uid": random_string(16),
        "track_number": track_number,
        "entry": 'WBIL',
        "delivery": generate_delivery(),
        "payment": generate_payment(),
        "items": [generate_item(track_number) for i in range(random.randint(1, 5))],
        "locale": "en",
        "internal_signature": random_string(10),
        "customer_id": random_string(20),
        "delivery_service": random_string(5),
        "shardkey": str(random.randint(0, 10)),
        "sm_id": random.randint(0, 10000),
        "date_created": random_date_iso(),
        "oof_shard": str(random.randint(0, 10))
    }


def main():
    producer = KafkaProducer(
        bootstrap_servers=['localhost:20092'],
        value_serializer=lambda v: json.dumps(v).encode('utf-8')
    )
    while True:
        producer.send('orders', value=generate_message())
        producer.flush()
        time.sleep(1)


if __name__ == '__main__':
    main()
