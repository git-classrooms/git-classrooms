#!/bin/env bash

set -e

go mod download
go generate
go test -v ./...
