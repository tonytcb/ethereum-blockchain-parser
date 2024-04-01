SHELL := /bin/bash

.PHONY: help up down clear

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
