#!/bin/bash

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
until PGPASSWORD=postgres psql -h postgres -U postgres -d chatbox -c '\q'; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done

echo "PostgreSQL is up - executing schema"
PGPASSWORD=postgres psql -h postgres -U postgres -d chatbox -f /app/scripts/init-db.sql

echo "Database initialization completed"
