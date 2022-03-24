SHELL := bash
.ONESHELL:
MAKEFLAGS += --no-builtin-rules

NOCACHE := $(if $(NOCACHE),"--no-cache")
PROTO_GENERATED_PACKAGE := "pkg/proto"

export APP_NAME := exchange-rate-test-task
export DOCKER_REPOSITORY := harzu
export VERSION := $(if $(VERSION),$(VERSION),$(if $(COMMIT_SHA),$(COMMIT_SHA),$(shell git rev-parse --verify HEAD)))

.PHONY: help
help: ## List all available targets with help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: build-helper build-prod ## Build all containers

.PHONY: build-helper
build-helper:
	@docker build ${NOCACHE} --pull -f ./build/helper.Dockerfile -t ${DOCKER_REPOSITORY}/${APP_NAME}-helper:${VERSION} .

.PHONY: build-prod
build-prod:
	@docker build ${NOCACHE} --pull -f ./build/Dockerfile -t ${DOCKER_REPOSITORY}/${APP_NAME}-helper:${VERSION} .

.PHONY: run
run: ## Run develop docker-compose
	@docker-compose up app

.PHONY: stop ## Stop all develop containers
stop:
	@docker-compose down -v

.PHONY: test-short
test-short: ## Run unit tests
	@go test ./... -cover -short

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run