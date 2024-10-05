#!/bin/sh
set -e

if [ "$1" = "node" ] || [ "$1" = "yarn" ]; then
  yarn install --no-progress
fi

exec "$@"