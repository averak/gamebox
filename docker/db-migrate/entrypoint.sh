#!/usr/bin/env sh

psqldef --host="$POSTGRES_HOST" --port="$POSTGRES_PORT" --user="$POSTGRES_USER" --password="$POSTGRES_PASSWORD" --file=/ddl.sql --enable-drop-table "$POSTGRES_DB"
