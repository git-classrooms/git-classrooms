#!/usr/bin/env bash

set -e

docker compose -f docker-compose.local.yml up -d

air -c .air.toml &
AIR_PID=$!

cd frontend
yarn

yarn dev &
YARN_PID=$!

cd ..

function kill_dev() {
  kill $AIR_PID
  kill $YARN_PID
  docker compose -f docker-compose.local.yml stop
  exit 0
}

trap kill_dev SIGINT SIGTERM

while true; do
  sleep 1
done
