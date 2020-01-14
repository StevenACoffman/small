SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

GO11MODULES=on
APP?=application
REGISTRY?=docker.io/stevenacoffman
COMMIT_SHA=$(shell git rev-parse --short HEAD)
VERSION=`git rev-parse HEAD`
BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

.PHONY: build
## build: build the application
build: clean
	@echo "Building..."
	@go build -o ${APP} main.go

.PHONY: run
## run: runs go run main.go
run:
	go run -race main.go

.PHONY: clean
## clean: cleans the binary
clean:
	@echo "Cleaning"
	@rm -rf ${APP}
	@go get -u -v golang.org/x/tools/cmd/goimports
	@gofmt -l -w -s .
	@goimports -l -w .

.PHONY: test
## test: runs go test with default values
test:
	go test -v -count=1 -race ./...

.PHONY: lint
## lint: lint the application code for problems
lint:
	@golangci-lint run

.PHONY: docker-build
docker-build:	## - Build the smallest and secured golang docker image based on distroless static
	@printf "\033[32m\xE2\x9c\x93 Build the smallest and secured golang docker image based on distroless static\n\033[0m"
	@export DOCKER_CONTENT_TRUST=1 && docker build -f Dockerfile -t ${REGISTRY}/${APP}:${COMMIT_SHA} .

.PHONY: docker-build-no-cache
docker-build-no-cache:	## - Build the smallest and secured golang docker image based on scratch with no cache
	@printf "\033[32m\xE2\x9c\x93 Build the smallest and secured golang docker image based on scratch\n\033[0m"
	@export DOCKER_CONTENT_TRUST=1 && docker build --no-cache -f Dockerfile -t ${REGISTRY}/${APP}:${COMMIT_SHA} .

.PHONY: ls
ls: ## - List 'smallest-secured-golang' docker images
	@printf "\033[32m\xE2\x9c\x93 Look at the size dude !\n\033[0m"
	@docker image ls ${REGISTRY}/${APP}

.PHONY: docker-run
docker-run:	## - Run the smallest and secured golang docker image based on distroless static
	@printf "\033[32m\xE2\x9c\x93 Run the smallest and secured golang docker image based on scratch\n\033[0m"
	@docker run ${REGISTRY}/${APP}:${COMMIT_SHA}

.PHONY: docker-push
## docker-push: pushes the stringifier docker image to registry
docker-push: docker-build
	docker push ${REGISTRY}/${APP}:${COMMIT_SHA}

.PHONY: help
## help: Prints this help message
help: ## - Show help message
	@printf "\033[32m\xE2\x9c\x93 usage: make [target]\n\n\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
