#!/bin/bash
set -euo pipefail

until /opt/kafka/bin/kafka-broker-api-versions.sh --bootstrap-server kafka-1:9092 >/dev/null 2>&1 &&
      /opt/kafka/bin/kafka-broker-api-versions.sh --bootstrap-server kafka-2:9092 >/dev/null 2>&1 &&
      /opt/kafka/bin/kafka-broker-api-versions.sh --bootstrap-server kafka-3:9092 >/dev/null 2>&1; do
  echo "Waiting for Kafka broker..."
  sleep 2
done

/opt/kafka/bin/kafka-topics.sh --create --if-not-exists \
  --topic vacancy.events \
  --bootstrap-server kafka-1:9092,kafka-2:9092,kafka-3:9092 \
  --partitions 3 \
  --replication-factor 3

echo "Kafka topic vacancy.events is ready."

/opt/kafka/bin/kafka-topics.sh --create --if-not-exists \
  --topic companymember.events \
  --bootstrap-server kafka-1:9092,kafka-2:9092,kafka-3:9092 \
  --partitions 3 \
  --replication-factor 3

/opt/kafka/bin/kafka-topics.sh --create --if-not-exists \
  --topic company.events \
  --bootstrap-server kafka-1:9092,kafka-2:9092,kafka-3:9092 \
  --partitions 3 \
  --replication-factor 3

echo "Kafka topic company.events is ready."

/opt/kafka/bin/kafka-topics.sh --create --if-not-exists \
  --topic resume.events \
  --bootstrap-server kafka-1:9092,kafka-2:9092,kafka-3:9092 \
  --partitions 6 \
  --replication-factor 3

echo "Kafka topic resume.events is ready."

/opt/kafka/bin/kafka-topics.sh --create --if-not-exists \
  --topic candidate.events \
  --bootstrap-server kafka-1:9092,kafka-2:9092,kafka-3:9092 \
  --partitions 3 \
  --replication-factor 3

echo "Kafka topic candidate.events is ready."

/opt/kafka/bin/kafka-topics.sh --create --if-not-exists \
  --topic user.events \
  --bootstrap-server kafka-1:9092,kafka-2:9092,kafka-3:9092 \
  --partitions 3 \
  --replication-factor 3

echo "Kafka topic user.events is ready."

/opt/kafka/bin/kafka-topics.sh --create --if-not-exists \
  --topic dlq \
  --bootstrap-server kafka-1:9092,kafka-2:9092,kafka-3:9092 \
  --partitions 3 \
  --replication-factor 3

echo "Kafka topic dlq is ready."