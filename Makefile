SHELL := /bin/bash

.PHONY: help up down clear lint tests generate-mocks

## help: shows this help message
help:
	@ echo "Usage: make [target]"
	@ sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## up: starts the application test
up: 
	docker-compose up parser
	docker-compose down

## down: put down all docker containers
down:
	docker-compose down

## clean: clean up all docker containers and volumes
clean: down
	docker ps -aq | xargs docker stop | xargs docker rm
	rm -rf ./dbdata

## lint: runs linter for all packages
lint:
	@ docker run  --rm -v "`pwd`:/workspace:cached" -w "/workspace/." golangci/golangci-lint:latest golangci-lint run  --timeout 5m

## tests: runs all tests, with no Docker
tests:
	go test -v -race -tags=integration ./pkg/...

## generate-mocks: generates mocks for all interfaces through mockery
generate-mocks:
	- mockery --all --dir ./pkg --output ./pkg/mocks --exported --case=underscore
