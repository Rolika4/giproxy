ENTRY_FILE = ./cmd/main.go
CURRENT_DIR=$(shell pwd)
DIST_DIR=${CURRENT_DIR}/bin
BIN_NAME=giproxy
DOCKER_WORKSPACE=rolliku
VERSION?=0.1.0

HOST_OS?=$(shell go env GOOS)
HOST_ARCH?=$(shell go env GOARCH)

.DEFAULT_GOAL:=help


# set default shell
SHELL=/bin/bash -o pipefail -o errexit

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: install
install: ## Installing dependencies
	@echo ">>> Installing dependencies..."
	go mod tidy
	go mod download

.PHONY: install-nodemon
install-nodemon: ## Installing nodemon tool for hot reload
	@echo ">>> Checking if nodemon is installed..."
	@if ! [ -x "$$(command -v nodemon)" ]; then \
		echo ">>> Installing nodemon globally with npm..."; \
		npm install -g nodemon; \
	else \
		echo ">>> nodemon is already installed."; \
	fi
.PHONY: run
run: ## Running the application
	@echo ">>> Starting application..."
	go run $(ENTRY_FILE)

.PHONY: clean
clean: ## Cleaning up temporary files
	@echo ">>> Cleaning up temporary files..."
	rm -rf ./tmp ./bin

.PHONY: watch
watch: ## Running the application with nodemon (hot reload)
	@echo ">>> Starting application with nodemon..."
	nodemon --exec "go run $(ENTRY_FILE)" --ext go --signal SIGKILL

.PHONY: build
build: ## Building the application
	@echo ">>> Building application..."
	CGO_ENABLED=0 GOOS=${HOST_OS} GOARCH=${HOST_ARCH} go build -v -o ${DIST_DIR}/${BIN_NAME} ${ENTRY_FILE}

.PHONY: lint
lint: golangci-lint ## Run go lint
	${GOLANGCILINT} run


.PHONY: build-image
build-image: ## Building the docker image
	@echo ">>> Building docker image..."
## GOOS=darwin GOARCH=arm64
	docker buildx build --platform linux/arm64,linux/amd64 -t ${DOCKER_WORKSPACE}/${BIN_NAME}:${VERSION} .

.PHONY: reset
reset: ## Resetting dependencies
	@echo ">>> Resetting dependencies..."
	rm -rf go.sum
	$(MAKE) install


GOLANGCILINT = ${CURRENT_DIR}/bin/golangci-lint

.PHONY: golangci-lint
golangci-lint: ## Download golangci-lint locally if necessary.
	$(call go-get-tool,$(GOLANGCILINT),github.com/golangci/golangci-lint/cmd/golangci-lint,v1.62.2)

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Downloading $(2)" ;\
go get -d $(2)@$(3) ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef