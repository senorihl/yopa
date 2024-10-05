#!/usr/bin/env sh
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE TABLE IF NOT EXISTS migrations ( name VARCHAR, run_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), PRIMARY KEY (name) );
EOSQL