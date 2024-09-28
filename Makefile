BUF_VERSION=v1.42.0
SQL_BOILER_VERSION=v4.16.2
GO_LDFLAGS := -s -w -X github.com/averak/gamebox/app/core/build_info.serverVersion=$(shell git describe --tags --always)

.PHONY: install-tools
install-tools:
	go install ./cmd/protoc-gen-gamebox-server
	go install github.com/bufbuild/buf/cmd/buf@${BUF_VERSION}
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/arch-go/arch-go@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/volatiletech/sqlboiler/v4@${SQL_BOILER_VERSION}
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

.PHONY: format
format:
	go fmt ./...
	golangci-lint run --issues-exit-code=0 --fix ./...
	buf format --write

BREAKING_CHANGE_BASE_BRANCH?=develop
.PHONY: lint
lint:
	golangci-lint run --issues-exit-code=1 ./...
	arch-go
	buf lint
	buf breaking --against '.git#branch=$(BREAKING_CHANGE_BASE_BRANCH)'

.PHONY: codegen
codegen:
	find . -type f \( -name 'wire_gen.go' \) -delete
	wire ./...
	find . -type f \( -name '*.connect.go' -or -name '*.pb.go' -or -name '*.gamebox.go' \) -delete
	buf generate
	sqlboiler psql --wipe --templates=templates/sqlboiler,$(shell go env GOPATH)/pkg/mod/github.com/volatiletech/sqlboiler/v4@${SQL_BOILER_VERSION}/templates/main

.PHONY: test
test:
	mkdir -p tmp/coverage
	GAMEBOX_CONFIG_FILEPATH=$(shell pwd)/config/default.json TZ=UTC go test -p=1 -coverpkg=./... -coverprofile=tmp/coverage/cover.out ./...
	go tool cover -html=tmp/coverage/cover.out -o tmp/coverage/cover.html

.PHONY: db-migrate
db-migrate:
	docker-compose run --rm --build db-migrate

.PHONY: db-clean
db-clean:
	docker-compose down -v postgres
	docker-compose up -d postgres

.PHONY: build
build: build-api-server

.PHONE: build-api-server
build-api-server:
	CGO_ENABLED=0 go build -ldflags="$(GO_LDFLAGS)" -o tmp/build/api_server ./entrypoint/api_server

.PHONY: run-api-server
run-api-server:
	go run -ldflags="$(GO_LDFLAGS)" ./entrypoint/api_server
