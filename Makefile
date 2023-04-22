PROJECT_NAME := "user-service"
PKG_LIST := $(shell go list ./... | grep -v /vendor/)

.PHONY: proto all dep build clean test coverage coverhtml lint gotidy migrateup migratedown sqlc

all: build

init: 

format: ## Format the files
	@go fmt ./...

lint: ## Lint the files
	golangci-lint run ./... --verbose

test: ## Run unittests
	@go test -short ./... -coverprofile=cov.out `go list ./... | grep -v vendor/` fmt

race:  go-modules ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan:  go-modules ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	./scripts/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	./scripts/coverage.sh html;

go-modules: ## Get the dependencies
	@go mod tidy
	@go mod vendor

gotidy: ## Get the latest dependencies
	@rm -fr go.sum
	@go mod tidy -compat=${GOVERSION}
	@go mod vendor

goclean: ## Get the latest dependencies
	@rm -fr go.sum
	@go clean -modcache
	@go mod tidy -compat=${GOVERSION}
	@go mod vendor

build:  go-modules ## Build the binary file
	@go build -v -o bin/${PROJECT_NAME} .

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

server:
	@go run main.go

docker:
	@make linux_amd64
	@docker build -t ${PROJECT_NAME}:latest .

.PHONY: all-sys
all-sys: darwin_amd64 darwin_arm64 linux_amd64 linux_arm64 windows_amd64

.PHONY: darwin_amd64
darwin_amd64:
	GOOS=darwin  GOARCH=amd64 $(MAKE) build

.PHONY: darwin_arm64
darwin_arm64:
	GOOS=darwin  GOARCH=arm64 $(MAKE) build

.PHONY: linux_amd64
linux_amd64:
	GOOS=linux   GOARCH=amd64 $(MAKE) build

.PHONY: linux_arm64
linux_arm64:
	GOOS=linux   GOARCH=arm64 $(MAKE) build

.PHONY: windows_amd64
windows_amd64:
	GOOS=windows GOARCH=amd64 EXTENSION=.exe $(MAKE) build
