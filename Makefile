.DEFAULT_GOAL := init

.PHONY: init
init:
ifeq ($(shell basename $(shell pwd)),evergram-identity)
	@echo "Already in the correct directory"
else
	cd $(shell dirname $(shell realpath $(lastword $(MAKEFILE_LIST))))
endif

.PHONY: build
build:
	@go build -o /dev/null ./cmd

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: gen
gen:
	@protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.

.PHONY: check
check:
	@go build -o /dev/null ./...

.PHONY: test
test:
	@go test -cover -race ./... -v -short 

.PHONY: run
run:
	go mod tidy
	go mod vendor
	docker-compose -p evergram up -d  --build && docker-compose logs -f

.PHONY: restart
restart:
	docker-compose -p evergram restart -t 10 identity
