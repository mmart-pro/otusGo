DIR_API = ./bin/calendar
BIN_API = ${DIR_API}/calendar
CFG_API = ${DIR_API}/config.json
CFG_TMPL_API = ./configs/calendar-config.json.tmpl
DIR_SNDR = ./bin/calendar_sender
BIN_SNDR = ${DIR_SNDR}/calendar_sender
CFG_SNDR = ${DIR_SNDR}/config.json
CFG_TMPL_SNDR = ./configs/sender-config.json.tmpl
DIR_SCHD = ./bin/calendar_scheduler
BIN_SCHD = ${DIR_SCHD}/calendar_scheduler
CFG_SCHD = ${DIR_SCHD}/config.json
CFG_TMPL_SCHD = ./configs/scheduler-config.json.tmpl
ENV_NAME = ./.env.dev
include ${ENV_NAME}
DOCKER_IMG = "calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

dirs:
	mkdir -p ${DIR_API}
	mkdir -p ${DIR_SCHD}
	mkdir -p ${DIR_SNDR}

config: dirs
	env `cat ${ENV_NAME}` envsubst < ${CFG_TMPL_API} > ${CFG_API};
	env `cat ${ENV_NAME}` envsubst < ${CFG_TMPL_SCHD} > ${CFG_SCHD};
	env `cat ${ENV_NAME}` envsubst < ${CFG_TMPL_SNDR} > ${CFG_SNDR};

build: generate config
	go build -v -o $(BIN_API) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHD) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN_SNDR) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

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
	$(BIN_SCHD) --version
	$(BIN_SNDR) --version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	goose -dir migrations postgres "host=$(CALENDAR_STORAGE_HOST) port=$(CALENDAR_STORAGE_PORT) user=$(CALENDAR_STORAGE_USER) password=$(CALENDAR_STORAGE_PASSWORD) dbname=$(CALENDAR_STORAGE_DB_NAME) sslmode=disable" up

generate:
	go generate ./...
#	rm -Rf internal/server/grpc/pb
#	mkdir internal/server/grpc/pb
#	protoc -I ./api --go_out=. --go-grpc_out=. --grpc-gateway_out=. --openapiv2_out=. api/eventservice/eventservice.proto

.PHONY: dirs config build run build-img run-img version test lint migrate generate
