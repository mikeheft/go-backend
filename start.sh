#!/bin/sh

set -e

echo "[INFO] Running db migration"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "[INFO] Starting the app"
exec "$@"