.PHONY: help build format lint test clean

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/  /'

## build: Build the package and examples
build:
	@echo "==> Building..."
	go build -v ./...

## format: Format Go source code
format:
	@echo "==> Formatting..."
	go fmt ./...

## lint: Run Go vet
lint:
	@echo "==> Linting..."
	go vet ./...

## test: Run tests
test:
	@echo "==> Testing..."
	go test -v ./...

## clean: Clean build cache
clean:
	@echo "==> Cleaning..."
	go clean

