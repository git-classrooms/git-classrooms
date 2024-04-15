#!/bin/env bash

set -e

COMMAND=$1

if [ -z $COMMAND ]; then 
    echo "Usage: $0 one|air"
    exit 1
fi

if [ $COMMAND == "one" ]; then
    go mod download 
    go generate
    go test -v ./...
fi

if [ $COMMAND == "air" ]; then
    go mod download 
    air -c .air.test.toml
fi
