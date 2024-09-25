#!/bin/env bash

go install github.com/vektra/mockery/v2@v2.42.2
go install github.com/swaggo/swag/cmd/swag@latest

go mod download
go generate
