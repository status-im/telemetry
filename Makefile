.PHONY: all build lint test

all: build

deps: lint-install

build:
	go build -o build/server cmd/server/main.go

vendor:
	go mod tidy

lint-install:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		bash -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint:
	@echo "lint"
	@golangci-lint --exclude=SA1019 run ./... --deadline=5m

test: postgres
	go test -v -failfast ./...

generate:
	go generate ./telemetry/sql

run: build postgres
	./build/server --data-source-name=postgres://telemetry:newPassword@127.0.0.1:5432/telemetrydb?sslmode=disable

postgres:
	docker inspect telemetry-postgres > /dev/null ||\
		docker run --name telemetry-postgres -e POSTGRES_USER=telemetry -e POSTGRES_PASSWORD=newPassword -e POSTGRES_DB=telemetrydb -p 5432:5432 -d postgres &&\
		sleep 3

postgres-clean:
	docker stop telemetry-postgres &&\
		docker rm telemetry-postgres
