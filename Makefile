.PHONY: build test lint fmt install clean run

BINARY := bin/vaultenv
MAIN   := ./main.go

build:
	go build -o $(BINARY) $(MAIN)

test:
	go test ./... -v -cover

lint:
	golangci-lint run

fmt:
	gofmt -w .

install:
	go install ./...

clean:
	rm -rf bin/

run:
	go run $(MAIN)
