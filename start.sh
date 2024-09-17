#!/bin/sh

set -e

echo "[INFO] Running db migration"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "[INFO] Starting the app"
exec "$@"