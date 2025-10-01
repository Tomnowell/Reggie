#!/bin/sh

set -e # Exit early if any commands fail

if [ ! -d "$(dirname "$0")/app" ] || [ -z "$(ls -A "$(dirname "$0")/app"/*.go 2>/dev/null)" ]; then
  echo "Error: 'app' directory or Go files not found."
  exit 1
fi

(
  cd "$(dirname "$0")" # Ensure compile steps are run within the repository directory
  go build -o /tmp/reggie app/*.go
)
exec /tmp/reggie "$@"
