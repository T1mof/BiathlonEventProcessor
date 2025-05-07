ifeq ($(OS),Windows_NT)
    EXEEXT := .exe
else
    EXEEXT :=
endif

BINARY := biathlon$(EXEEXT)
PKGS   := ./internal/... ./cmd/biathlon

.PHONY: all fmt vet lint test build run clean docker

all: fmt vet test build

fmt:
	@go fmt ./...

vet:
	@go vet ./...

lint:
	golangci-lint run

test:
	@go test $(PKGS)

build:
	@echo ">> go build"
	@go build -o $(BINARY) ./cmd/biathlon

run: build
	@echo ">> ./$(BINARY)"
	@./$(BINARY)

