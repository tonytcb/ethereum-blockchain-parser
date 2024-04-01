SHELL := /bin/bash

.PHONY: help lint tests generate-mocks

## help: shows this help message
help:
	@ echo "Usage: make [target]"
	@ sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## lint: runs linter for all packages
lint:
	@ docker run  --rm -v "`pwd`:/workspace:cached" -w "/workspace/." golangci/golangci-lint:latest golangci-lint run  --timeout 5m

## tests: runs all tests, with no Docker
tests:
	go test -v -race -tags=integration ./pkg/...

## generate-mocks: generates mocks for all interfaces through mockery
generate-mocks:
	- mockery --all --dir ./pkg --output ./pkg/mocks --exported --case=underscore
