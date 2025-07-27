#!/bin/bash
set -e

# Wait for Cassandra to be ready
echo "Waiting for Cassandra to start..."
until cqlsh -e "DESC KEYSPACES;" > /dev/null 2>&1; do
    echo "Cassandra is unavailable - sleeping"
    sleep 2
done

echo "Cassandra is up - executing migrations"

# Run all CQL files in order
for f in /migrations/*.cql; do
    if [ -f "$f" ]; then
        echo "Running $f"
        cqlsh -f "$f"
    fi
done

echo "Cassandra migrations completed"