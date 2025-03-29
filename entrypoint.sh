#!/bin/sh

# GOOSE=$(which goose)

# echo "Running database migrations..."
# $GOOSE -dir ./sql/schema postgres "$DATABASE_URL" up || (echo "Migration failed, retrying in 5s..." && sleep 5 && $GOOSE -dir ./sql/schema postgres "$DATABASE_URL" up)

# echo "DATABASE_URL is: $DATABASE_URL"

# echo "Starting Webcrawler..."
# exec ./crawler

GOOSE=$(which goose)

# Compose database URL
export DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

# Compose RabbitMQ URL
export RABBITMQ_URL="amqp://${RABBITMQ_USER}:${RABBITMQ_PASSWORD}@${RABBITMQ_HOST}:${RABBITMQ_PORT}/"

echo "Running database migrations..."
$GOOSE -dir ./sql/schema postgres "$DATABASE_URL" up || (echo "Migration failed, retrying in 5s..." && sleep 5 && $GOOSE -dir ./sql/schema postgres "$DATABASE_URL" up)

echo "DATABASE_URL is: $DATABASE_URL"

echo "Starting Webcrawler..."
exec ./crawler
