#!/usr/bin/env bash

if [ ! -f frontend/dist/robots.txt ]; then
  mkdir -p frontend/dist
  cp frontend/public/robots.txt frontend/dist/robots.txt
fi

docker compose -f docker-compose.local.yaml up -d

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
    docker compose -f docker-compose.local.yaml stop
    exit 0
}

trap kill_dev SIGINT SIGTERM

while true; do
  sleep 1
done
