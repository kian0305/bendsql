.PHONY: build

default: build

build: fmt vet
	go build -o bin/bendctl cmd/bendctl/main.go

test:
	GO111MODULE=on go test -p 1 -v -race ./...
	go vet ./...

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...
clean:
	rm bin/*
