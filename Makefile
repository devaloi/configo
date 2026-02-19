.PHONY: build test lint fmt clean

build:
	go build ./...

test:
	go test -race ./...

lint:
	golangci-lint run

fmt:
	gofmt -w .

clean:
	go clean -testcache

all: fmt lint build test
