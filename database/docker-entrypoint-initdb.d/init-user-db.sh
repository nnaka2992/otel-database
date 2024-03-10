#!/usr/bin/env sh

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER app WITH PASSWORD 'otel_password';
    CREATE DATABASE otel;
    \c otel;
    GRANT ALL ON SCHEMA public TO app;
EOSQL

PG_PASSWORD='otel_password' psql -v ON_ERROR_STOP=1 --username "app"  --dbname "otel" < /db/schema/*.sql
