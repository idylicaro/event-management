#!/bin/bash

# Load environment variables from .env if it exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Environment variables configuration
DB_HOST=${POSTGRES_HOST:-"localhost"}
DB_PORT=${POSTGRES_PORT:-"5432"}
DB_USER=${POSTGRES_USER:-"default_user"}
DB_PASSWORD=${POSTGRES_PASSWORD:-"default_password"}
DB_NAME=${POSTGRES_DB:-"default_db"}
SSL_MODE=${SSL_MODE:-"disable"}

# Build the connection URL
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${SSL_MODE}"

# Path to migration files
MIGRATIONS_PATH="migrations"

# Check if the action is provided
ACTION=$1
if [ -z "$ACTION" ]; then
  echo "Error: No action specified (use 'up', 'down', 'version', etc.)"
  exit 1
fi

# Run the migration
case $ACTION in
  up)
    echo "Running migrations up..."
    migrate -path $MIGRATIONS_PATH -database $DATABASE_URL up
    ;;
  down)
    echo "Reverting migrations down..."
    migrate -path $MIGRATIONS_PATH -database $DATABASE_URL down
    ;;
  version)
    echo "Getting the current migration version..."
    migrate -path $MIGRATIONS_PATH -database $DATABASE_URL version
    ;;
  *)
    echo "Unknown action: $ACTION. Use 'up', 'down', 'version', etc."
    exit 1
    ;;
esac
