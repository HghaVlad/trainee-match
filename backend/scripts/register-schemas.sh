#!/bin/bash

set -e

# defaults
REGISTRY_URL=${SCHEMA_REGISTRY_URL:-http://localhost:8081}
SCHEMA_DIR=${SCHEMA_DIR:-/schemas}

echo "Using registry: $REGISTRY_URL"
echo "Scanning: $SCHEMA_DIR"

# wait for registry
until curl -fsS "${REGISTRY_URL}/subjects" >/dev/null 2>&1; do
  echo "Waiting for Schema Registry at ${REGISTRY_URL}..."
  sleep 2
done

# process schemas
find "$SCHEMA_DIR" -name "*.avsc" | while read -r file; do
  echo "----------------------------------------"
  echo "Processing: $file"

  name=$(basename "$file" .avsc)
  subject="${name}-value"

  echo "Subject: $subject"

  schema=$(jq -c . < "$file" | jq -Rs .)

  # check if subject exists
  http_code=$(curl -s -o /dev/null -w "%{http_code}" \
    "$REGISTRY_URL/subjects/$subject/versions")

  if [[ "$http_code" == "200" ]]; then
    echo "Subject exists, checking compatibility..."

    result=$(curl -s -X POST \
      "$REGISTRY_URL/compatibility/subjects/$subject/versions/latest" \
      -H "Content-Type: application/vnd.schemaregistry.v1+json" \
      -d "{\"schema\": $schema}")

    compatible=$(echo "$result" | jq -r '.is_compatible // "false"')

    if [[ "$compatible" != "true" ]]; then
      echo "Incompatible schema for subject: $subject"
      echo "$result"
      exit 1
    fi

    echo "Compatible"
  elif [[ "$http_code" == "404" ]]; then
    echo "Subject does not exist yet, will create"
  else
    echo "Unexpected response when checking subject: HTTP $http_code"
    exit 1
  fi

  echo "Registering schema..."

  register_result=$(curl -s -X POST \
    "$REGISTRY_URL/subjects/$subject/versions" \
    -H "Content-Type: application/vnd.schemaregistry.v1+json" \
    -d "{\"schema\": $schema}")

  id=$(echo "$register_result" | jq -r '.id // empty')

  if [[ -z "$id" ]]; then
    echo "Failed to register schema:"
    echo "$register_result"
    exit 1
  fi

  echo "Registered with id: $id"
done

echo "----------------------------------------"
echo "Done"