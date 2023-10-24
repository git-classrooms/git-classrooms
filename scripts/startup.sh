#! /bin/sh

air -c .air.toml &

cd assets

yarn
yarn gen
yarn dev

