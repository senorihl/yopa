#!/usr/bin/env bash
set -e

already_runned=$(psql -qtAX --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -c "SELECT run_at FROM migrations WHERE name = '01_pixel' LIMIT 1;")

if [ ! -z "$already_runned" ]; then
  exit 0
fi

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /docker-entrypoint-migrations.d/01_pixel.sql