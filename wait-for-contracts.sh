#!/bin/sh
# wait-for-contracts.sh

set -e

# Path to the contract addresses file
CONTRACTS_FILE="/artifacts/deployed-contracts.json"

# Wait for the file to exist
until [ -f "$CONTRACTS_FILE" ]; do
  >&2 echo "Contracts file not found - waiting"
  sleep 1
done

>&2 echo "Contracts file found - starting server"
# Execute the command passed to the script
exec "$@"