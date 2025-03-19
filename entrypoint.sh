#!/bin/sh

GOOSE=$(which goose)

echo "Running database migrations..."
$GOOSE -dir ./sql/schema postgres "$DATABASE_URL" up || (echo "Migration failed, retrying in 5s..." && sleep 5 && $GOOSE -dir ./sql/schema postgres "$DATABASE_URL" up)

echo "DATABASE_URL is: $DATABASE_URL"

echo "Starting Webcrawler..."
exec ./crawler
