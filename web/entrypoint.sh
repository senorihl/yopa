#!/bin/sh
set -ex

if [ "$1" = "node" ] || [ "$1" = "yarn" ]; then
  yarn install --no-progress --frozen-lockfile
fi

exec "$@"