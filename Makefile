FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp

## up: starts all containers in the background qitout forcing build
up:
	@echo "Starting Docker images"
	docker-compose up -d
	@echo "Docker images started"

## down: stop docker compose
down:
	@echo "Stopping docker compose"
	docker-compose down
	@echo "docker compose down done"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build:
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building when required and starting docker images"
	docker-compose up --build -d
	@echo "Docker images built and started"

# build_broker: builds the broker binary
# build_broker:
# 	@echo "Building broker binary"
# 	cd ./broker-service && env GOOS=linux CGO_ENABLE=0 go build -o ${BROKER_BINARY} main.go
# 	@echo "Done"

## build_front: builds the front end binary
build_front:
	@echo "Building the front end binary"
	cd ./front-end/front-end/cmd/web && env GOOS=linux CGO_ENABLE=0 go build -o ${FRONT_END_BINARY} main.go
	@echo "Done"

## start: starts the front end
start: build_front
	@echo "Starting the front end"
	cd ./front-end/front-end/cmd/web && ./${FRONT_END_BINARY}

## stop: stops the front end
stop:
	@echo "Stopping the front end"
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end"

## postgres: creates a docker container with the credentials below
postgres:
	docker run --name auth-service -e POSTGRES_PASSWORD=verysecret -e POSTGRES_USER=admin -p 5432:5432 -d postgres:alpine3.19

## createdb: creates a new db in our container
createdb:
	docker exec -it auth-service createdb --username=root --owner=root authdb

## dropdb: drops all data in our postgres container
dropdb:
	docker exec -it auth-service dropdb authdb
