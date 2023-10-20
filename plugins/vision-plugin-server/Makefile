include .env
export

all: run
.PHONY: all

run: 
	@go run cmd/scigateway/main.go 
.PHONY: run

build: test
	@go build cmd/scigateway/main.go 
.PHONY: build

test: 
	@go test -v ./...
.PHONY: test