# Go Daemon — Makefile

BINARY := bin/daemon
DAEMON_PORT ?= 2222
SIDE_PORT ?= 8080

.PHONY: run build test vet lint docker clean install

run: build
	./$(BINARY) --verbose

build:
	go build -o $(BINARY) ./cmd/daemon

test:
	go test -race -cover ./...

vet:
	go vet ./...

lint: vet
	golangci-lint run ./...

install:
	go install ./cmd/daemon

docker:
	docker build -t go-daemon-template .

clean:
	rm -rf bin/