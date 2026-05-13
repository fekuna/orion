#!/usr/bin/env bash
# =============================================================================
# docker/postgres/init/01_create_databases.sh
#
# Runs once on first container startup (via docker-entrypoint-initdb.d).
# Creates one database per service so each service is fully isolated.
# Add a new CREATE DATABASE line for every new service you add.
# =============================================================================
set -euo pipefail

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- product-service
    CREATE DATABASE product_db;

    -- future services (uncomment when the service is added)
    -- CREATE DATABASE transaction_db;
    -- CREATE DATABASE payment_db;
EOSQL

echo "orion-v2: databases created."
