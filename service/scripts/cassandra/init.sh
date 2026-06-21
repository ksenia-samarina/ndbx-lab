#!/bin/bash
set -e

echo "Starting Cassandra Initialization"

if [ -z "$CASSANDRA_HOSTS" ] || [ -z "$CASSANDRA_KEYSPACE" ]; then
    echo "Error: CASSANDRA_HOSTS or CASSANDRA_KEYSPACE is not set in environment."
    exit 1
fi

CASSANDRA_FIRST_HOST=$(echo "$CASSANDRA_HOSTS" | cut -d',' -f1)

echo "Connecting to Cassandra at ${CASSANDRA_FIRST_HOST}:${CASSANDRA_PORT}..."

CQLSH_ARGS=("${CASSANDRA_FIRST_HOST}" "${CASSANDRA_PORT}")

if [ -n "${CASSANDRA_USERNAME}" ] && [ -n "${CASSANDRA_PASSWORD}" ]; then
    CQLSH_ARGS+=("-u" "${CASSANDRA_USERNAME}" "-p" "${CASSANDRA_PASSWORD}")
fi

sed "s/\${CASSANDRA_KEYSPACE}/${CASSANDRA_KEYSPACE}/g" /schema.cql | cqlsh "${CQLSH_ARGS[@]}"

echo "Cassandra Schema Successfully Initialized"
