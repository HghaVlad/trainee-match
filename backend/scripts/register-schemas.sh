#!/bin/bash

set -e

until curl -fsS "${SCHEMA_REGISTRY_URL}/subjects" >/dev/null 2>&1; do
  echo "Waiting for Schema Registry at ${SCHEMA_REGISTRY_URL}..."
  sleep 2
done

# Validates all avro schemes and registers them in registry

REGISTRY_URL=${SCHEMA_REGISTRY_URL:-http://localhost:8081}
SCHEMA_DIR=${SCHEMA_DIR:-/schemas}

echo "Using registry: $REGISTRY_URL"
echo "Scanning: $SCHEMA_DIR"

for file in $(find "$SCHEMA_DIR" -name "*.avsc"); do
  echo "----------------------------------------"
  echo "Processing: $file"

  name=$(basename "$file" .avsc)

  subject="${name}-value"

  echo "Subject: $subject"

  schema=$(jq -c . < "$file" | jq -Rs .)

  # check if subject exists
  exists=$(curl -s "$REGISTRY_URL/subjects/$subject/versions" | jq 'type != "null"' || echo "false")

  if [[ "$exists" == "true" ]]; then
    echo "Checking compatibility..."

    result=$(curl -s -X POST \
      "$REGISTRY_URL/compatibility/subjects/$subject/versions/latest" \
      -H "Content-Type: application/vnd.schemaregistry.v1+json" \
      -d "{\"schema\": $schema}")

    compatible=$(echo "$result" | jq -r .is_compatible)

    if [[ "$compatible" != "true" ]]; then
      echo "Incompatible schema for subject: $subject"
      echo "$result"
      exit 1
    fi

    echo "Compatible"
  else
    echo "Subject does not exist yet, will create"
  fi

  echo "Registering schema..."

  register_result=$(curl -s -X POST \
    "$REGISTRY_URL/subjects/$subject/versions" \
    -H "Content-Type: application/vnd.schemaregistry.v1+json" \
    -d "{\"schema\": $schema}")

  id=$(echo "$register_result" | jq -r .id)

  echo "Registered with id: $id"
done

echo "----------------------------------------"
echo "Done"