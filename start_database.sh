#!/usr/bin/env sh

docker run --name otel-database -e POSTGRES_PASSWORD="test123" \
    -p 5432:5432 \
    -v "$(pwd)/docker-entrypoint-initdb.d":/docker-entrypoint-initdb.d \
    -v "$(pwd)/../backend/db/schema":/db/schema \
    postgres:16.2

docker ps -al | grep 'otel-database' | awk '{print $1}' | xargs docker rm
