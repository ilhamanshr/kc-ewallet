BIN_DIR := ./bin
MIGRATIONS_DIR=./migrations
DOCKER_COMPOSE_FILE=./tools/docker/docker-compose.yaml

GOCMD=go

GOPATH  := $(shell $(GOCMD) env GOPATH)
AIRPATH := $(GOPATH)/bin/air

PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

PHONY: docker/up
docker/up:
	docker compose -f ${DOCKER_COMPOSE_FILE} up -d

PHONY: docker/down
docker/down:
	docker compose -f ${DOCKER_COMPOSE_FILE} down

PHONY: run-migrate-create
run-migrate-create:
	@echo 'Creating migrations file for ${name}'
	migrate create -ext sql -dir ${MIGRATIONS_DIR} ${name}

run-migrate-up: 
	migrate -database "postgres://sonar:sonar@127.0.0.1:5432/backendservice?sslmode=disable" -path "./migrations" up

run-migrate-down-by-1: 
	migrate -database "postgres://sonar:sonar@127.0.0.1:5432/backendservice?sslmode=disable" -path "./migrations" down 1

run-migrate-down-all: 
	migrate -database "postgres://sonar:sonar@127.0.0.1:5432/backendservice?sslmode=disable" -path "./migrations" down -all

run-migrate-drop-force: 
	migrate -database "postgres://sonar:sonar@127.0.0.1:5432/backendservice?sslmode=disable" -path "./migrations"  drop -f

PHONY: run/migrate
run/migrate:
	make run-migrate-drop-force
	make run-migrate-up

.PHONY: sqlc
sqlc:
	sqlc generate --file="./tools/sqlc/sqlc.yaml"

.PHONY: build-services
build-services:
	docker build -f ./Dockerfile.http -t http-server:0.1.0 .

.PHONY: test
test:
	GOARCH=amd64 go test -v -race -count=1 ./...

.PHONY: watch
watch:
	test -s ${AIRPATH} || curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(GOPATH)/bin
	GO_APP_ENV=dev GIN_MODE=debug SQLCDEBUG=dumpast=1 ${AIRPATH}

.PHONY: generate
generate:
	go generate ./...