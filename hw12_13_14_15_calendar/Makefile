DIR_API = ./bin/calendar
BIN_API = ${DIR_API}/calendar
CFG_API = ${DIR_API}/config.json
CFG_TMPL_API = ./configs/api-config.json.tmpl
ENV_NAME = ./.env.dev
include ${ENV_NAME}
DOCKER_IMG = "calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

dirs:
	mkdir -p ${DIR_API}

config: dirs
	env `cat ${ENV_NAME}` envsubst < ${CFG_TMPL_API} > ${CFG_API};

build: generate config
	go build -v -o $(BIN_API) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build migrate
	$(BIN_API) --config $(CFG_API)

# build-img:
# 	docker build \
# 		--build-arg=LDFLAGS="$(LDFLAGS)" \
# 		-t $(DOCKER_IMG) \
# 		-f build/Dockerfile .

# run-img: build-img
# 	docker run $(DOCKER_IMG)

version: build
	$(BIN_API) --version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	goose -dir migrations postgres "host=$(STORAGE_HOST) port=$(STORAGE_PORT) user=$(STORAGE_USER) password=$(STORAGE_PASSWORD) dbname=$(STORAGE_DB_NAME) sslmode=disable" up

generate:
	go generate ./...
#	rm -Rf internal/server/grpc/pb
#	mkdir internal/server/grpc/pb
#	protoc -I ./api --go_out=. --go-grpc_out=. --grpc-gateway_out=. --openapiv2_out=. api/eventservice/eventservice.proto

.PHONY: dirs config build run build-img run-img version test lint migrate generate
