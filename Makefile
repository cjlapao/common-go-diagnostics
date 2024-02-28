NAME ?= common-go-diagnostics
CONTAINER_ID ?= test
export PACKAGE_NAME ?= $(NAME)
ifeq ($(OS),Windows_NT)
	export VERSION=$(shell type VERSION)
else
	export VERSION=$(shell cat VERSION)
endif

COBERTURA = cobertura

GOX = gox

GOLANGCI_LINT = golangci-lint

START_SUPER_LINTER_CONTAINER = start_super_linter_container

DEVELOPMENT_TOOLS = $(GOX) $(COBERTURA) $(GOLANGCI_LINT)
SECURITY_TOOLS = $(GOSEC)

.PHONY: help
help:
  # make version:
	# make test
	# make lint

.PHONY: version
version:
	@echo Version: $(VERSION)

.PHONY: test
test:
	@echo "Running tests..."
	@cd pkg && go test -v -covermode count ./...
	@echo "Tests finished."

.PHONY: coverage
coverage:
	@echo "Running coverage report..."
ifeq ("$(wildcard coverage)","")
	@echo "Creating coverage directory..."
	@mkdir coverage
endif
	@cd pkg && go test -coverprofile coverage.txt -covermode count -v ./...
	@cd pkg && gocov convert coverage.txt | gocov-xml >../coverage/cobertura-coverage.xml
	@cd pkg && rm coverage.txt

.PHONY: lint
lint: $(START_SUPER_LINTER_CONTAINER)
	@echo "Running linter..."
	@docker cp $(PACKAGE_NAME)-linter:/tmp/lint/super-linter.log .
	@echo "Linter report saved to super-linter.log"
	@docker stop $(PACKAGE_NAME)-linter
	@echo "Linter finished."

.PHONY: build
build:
	@echo "Building..."
ifneq ("$(wildcard out)","")
	@echo "Creating out directory..."
	@mkdir out
	@mkdir out/binaries
endif

	@cd pkg && go build -o ../out/binaries/$(PACKAGE_NAME)
	@echo "Build finished."

.PHONY: clean
clean:
	@echo "Cleaning..."
ifneq ("$(wildcard bin)","")
	@echo "Removing bin directory..."
	@rm -rf bin
endif
ifneq ("$(wildcard out)","")
	@echo "Removing out directory..."
	@rm -rf out
endif
ifneq ("$(wildcard coverage)","")
	@echo "Removing coverage directory..."
	@rm -rf out
endif
ifneq ("$(wildcard tmp)","")
	@echo "Removing tmp directory..."
	@rm -rf out
endif
	@echo "Clean finished."

.PHONY: security-check
security-check:
	@echo "Running security check..."
	@cd pkg && gosec ./...
	@echo "Security check finished."

.PHONY: deps
deps: $(DEVELOPMENT_TOOLS)

$(COBERTURA):
	@echo "Installing cobertura..."
	@go install github.com/axw/gocov/gocov@latest
	@go install github.com/AlekSi/gocov-xml@latest
	@go install github.com/matm/gocov-html/cmd/gocov-html@latest

$(GOX):
	@echo "Installing gox..."
	@go install github.com/mitchellh/gox@latest

$(GOLANGCI_LINT):
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

$(GOSEC):
	@echo "Installing gosec..."
	@go install github.com/securego/gosec/v2/cmd/gosec@latest

$(GOREPORTCARD):
	@echo "Installing goreportcard-cli..."
	@go install github.com/gojp/goreportcard/cmd/goreportcard-cli@latest

$(GOCYCLO):
	@echo "Installing gocylco..."
	@go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

$(START_SUPER_LINTER_CONTAINER):
ifeq ($(OS), Windows_NT)
	$(eval CONTAINER_ID := $(shell  docker ps -a -q -f "name=$(PACKAGE_NAME)-linter"))
	@IF "$(CONTAINER_ID)" EQU "" (\
	docker run --name $(PACKAGE_NAME)-linter -e DEFAULT_BRANCH=main -e RUN_LOCAL=true -e VALIDATE_ALL_CODEBASE=true -e VALIDATE_JSCPD=false -e CREATE_LOG_FILE=true -e VALIDATE_GO=false -v .:/tmp/lint ghcr.io/super-linter/super-linter:latest \
	) \
	ELSE (\
	docker start $(PACKAGE_NAME)-linter --attach \
	);
else
	$(eval CONTAINER_ID := $(shell docker ps -a | grep $(PACKAGE_NAME)-linter | awk '{print $$1}'))
	@if [ -z $(CONTAINER_ID) ]; then \
	echo "Linter container does not exist, creating it..."; \
	docker run --platform linux/amd64 --name $(PACKAGE_NAME)-linter -e DEFAULT_BRANCH=main -e RUN_LOCAL=true -e VALIDATE_ALL_CODEBASE=true -e VALIDATE_JSCPD=false -e CREATE_LOG_FILE=true -e VALIDATE_GO=false -v .:/tmp/lint ghcr.io/super-linter/super-linter:latest; \
	else \
	echo "Linter container already exists $(CONTAINER_ID), starting it..."; \
	docker start $(PACKAGE_NAME)-linter --attach; \
	fi
endif