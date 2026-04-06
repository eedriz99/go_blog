#!/bin/bash

# ------------------------------
# Config - Change these values
# ------------------------------
DB_USER="admin"
DB_PASSWORD="adminpassword"
DB_NAME="go_blog"
DB_HOST="localhost"
DB_PORT="5432"
MIGRATIONS_DIR="db/migrations"
# ------------------------------

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

COMMAND=$1
NAME=$2
STEPS=$3

# Ensure migrations folder exists
mkdir -p $MIGRATIONS_DIR

case "$COMMAND" in
  create)
    if [ -z "$NAME" ]; then
      echo "Usage: ./migrate.sh create migration_name"
      exit 1
    fi
    migrate create -ext sql -dir $MIGRATIONS_DIR $NAME

    # Wrap the .down.sql with IF EXISTS automatically
    DOWN_FILE=$(ls -t $MIGRATIONS_DIR/*$NAME.down.sql | head -1)
    if [ -f "$DOWN_FILE" ]; then
      sed -i '1iDROP TABLE IF EXISTS;' "$DOWN_FILE"
    fi
    echo "Migration '$NAME' created successfully."
    ;;
  
  up)
    # Handle dirty database automatically
    STATUS=$(migrate -path $MIGRATIONS_DIR -database $DATABASE_URL version 2>&1)
    if [[ $STATUS == *"dirty"* ]]; then
      DIRTY_VERSION=$(echo $STATUS | grep -o '[0-9]\+')
      echo "Dirty database detected at version $DIRTY_VERSION. Forcing..."
      migrate -path $MIGRATIONS_DIR -database $DATABASE_URL force $DIRTY_VERSION
    fi

    if [ -n "$STEPS" ]; then
      migrate -path $MIGRATIONS_DIR -database $DATABASE_URL up $STEPS
    else
      migrate -path $MIGRATIONS_DIR -database $DATABASE_URL up
    fi
    ;;

  down)
    if [ -n "$STEPS" ]; then
      migrate -path $MIGRATIONS_DIR -database $DATABASE_URL down $STEPS
    else
      migrate -path $MIGRATIONS_DIR -database $DATABASE_URL down
    fi
    ;;

  status)
    migrate -path $MIGRATIONS_DIR -database $DATABASE_URL version
    ;;

  *)
    echo "Usage: ./migrate.sh [create|up|down|status] [name|steps]"
    ;;
esac
