#!/usr/bin/env sh
# fetches from the base image
MIGRATE_CMD=$MIGRATE_CLI
if [ "$MIGRATE_CMD" = "" ] || [ ! -f "$MIGRATE_CMD" ]; then
    MIGRATE_CMD=/monkey/migrate
    if [ -f "/monkey/migrate" ]; then
        chmod +x  /monkey/migrate
    fi;
    if [ -f "/usr/local/bin/migrate" ]; then
        MIGRATE_CMD=/usr/local/bin/migrate
    fi;
fi;
chmod +x $MIGRATE_CMD
# make sure the migrate command is not being used by the filesytem
sync

MIGRATIONS_PATH="/monkey/migrations"
MYSQL_ENGINE="mysql"
# setting migration table name
MIGRATION_TABLE=$COMPONENT"_schema_migrations"

if [ "$DB_ENGINE" = "$MYSQL_ENGINE" ]; then
    echo "Running mysql migration..."
    LINK="mysql://$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?x-migrations-table=$MIGRATION_TABLE"
    $MIGRATE_CMD -verbose -database $LINK -path $MIGRATIONS_PATH/mysql up
else
    echo "Running postgres migration..."
    LINK="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable&x-migrations-table=$MIGRATION_TABLE"
    $MIGRATE_CMD -verbose -database $LINK -path $MIGRATIONS_PATH/psql up
fi
