MAKEFLAGS += --no-print-directory
MAIN_PACKAGE_PATH := .
BINARY_NAME := git-classrooms
APP_VERSION ?= $(shell git describe --tags --always --dirty)
APP_GIT_COMMIT ?= $(shell git rev-parse HEAD)
APP_GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
APP_GIT_REPOSITORY ?= https://github.com/git-classrooms/git-classrooms
APP_BUILD_TIME ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
define GOFLAGS
-ldflags=" \
  -s -w \
  -X main.version=$(APP_VERSION) \
"
endef
#   -X github.com/git-classrooms/git-classrooms/internal/storage/local/info.version=$(APP_VERSION) \
#   -X github.com/git-classrooms/git-classrooms/internal/storage/local/info.gitCommit=$(APP_GIT_COMMIT) \
#   -X github.com/git-classrooms/git-classrooms/internal/storage/local/info.gitBranch=$(APP_GIT_BRANCH) \
#   -X github.com/git-classrooms/git-classrooms/internal/storage/local/info.gitRepository=$(APP_GIT_REPOSITORY) \
#   -X github.com/git-classrooms/git-classrooms/internal/storage/local/info.buildTime=$(APP_BUILD_TIME) \

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build                             Build"
	@echo "  build/frontend                    Build frontend"
	@echo "  generate                          Generate"
	@echo "  generate/client                   Generate client pkg"
	@echo "  setup                             Install dependencies"
	@echo "  setup/ci                          Install dependencies for CI"
	@echo "  clean                             Clean"
	@echo "  run                               Run"
	@echo "  run/dev                           Run development environment"
	@echo "  build/docker                      Build docker image"
	@echo "  run/docker                        Run docker container"
	@echo "  infra/up                          Run infrastructure in docker compose (postgres, (pgadmin) and mailpit)"
	@echo "  infra/stop                        Run infrastructure stop"
	@echo "  infra/down                        Run infrastructure down (delete)"
	@echo "  migrate/new name=<name>           Create new migration"
	@echo "  migrate/check                     Check and print outstanding migrations"
	@echo "  migrate/status                    Migrate status"
	@echo "  seed/up                           Seed up"
	@echo "  seed/reset                        Seed reset"
	@echo "  tidy                              Fmt and Tidy"
	@echo "  lint                              Lint"
	@echo "  test                              Test"
	@echo "  test/verbose                      Test verbose"
	@echo "  config                            TODO: Edit config"
	@echo "  debug                             Debug"


.PHONY: run/dev
run/dev:
	@echo "Starting development environment..."
	@go tool concur || true

.PHONY: run
run: build
	@echo "Running binary..."
	@./bin/$(BINARY_NAME)

.PHONY: build
build:
	@if [ -z "$(CI)" ]; then \
		$(MAKE) generate; \
	fi
	@if ! [ -d ./frontend/dist ]; then \
		$(MAKE) build/frontend; \
	fi
	@echo "Building binary..."
	@CGO_ENABLED=0 go build $(GOFLAGS) -o ./bin/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)

.PHONY: build/frontend
build/frontend:
	@echo "Building frontend..."
	@cd frontend && yarn build

.PHONY: setup
setup: setup/frontend
	@echo "Setting up environment..."
	@# TODO: maybe go install github.com/mikefarah/yq/v4@latest
	go mod download

.PHONY: setup/ci
setup/ci:
	@echo "Installing..."
	go mod download

.PHONY: setup/frontend
setup/frontend:
	@cd frontend &&	yarn install

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf ./bin ./tmp
	@rm -f docs/docs.go docs/swagger.json docs/swagger.yaml
	@rm -rf $(shell yq '.packages | to_entries | map(.value.config.dir) | .[]' .mockery.yaml)
	@cd frontend && yarn clean

.PHONY: generate
generate:
	@echo "Generating code..."
	@go generate ./...
	@if [ -z "$(CI)" ]; then \
		$(MAKE) generate/client; \
	fi

.PHONY: generate/client
generate/client:
	@echo "Generating client..."
	@if ! [ -f "docs/swagger.yaml" ]; then \
		$(MAKE) generate; \
	fi
	@cd frontend && yarn generate

.PHONY: build/docker
build/docker:
	@echo "Running docker..."
	docker build -t $(BINARY_NAME)-docker \
		--build-arg APP_VERSION=$(APP_VERSION) \
		--build-arg APP_GIT_COMMIT=$(APP_GIT_COMMIT) \
		--build-arg APP_GIT_BRANCH=$(APP_GIT_BRANCH) \
		--build-arg APP_GIT_REPOSITORY=$(APP_GIT_REPOSITORY) \
		--build-arg APP_BUILD_TIME=$(APP_BUILD_TIME) \
		.

.PHONY: run/docker
run/docker:
	@echo "Running docker..."
	docker run -it --env-file .env --rm -p 3000:3000 $(BINARY_NAME)-docker

.PHONY: infra/up
infra/up:
	@docker compose -f docker-compose.local.yaml up -d

.PHONY: migrate/new
migrate/new: migrate/build
	@echo "Migrating up..."
	@if [ -z "$(name)" ]; then \
		echo "error: name is required"; \
		echo "usage: make migrate/new name=name_of_migration"; \
		exit 1; \
	fi
	go tool goose normal create $(name) sql

.PHONY: migrate/status
migrate/status: migrate/build
	@echo "Migrating status..."
	go tool goose normal status

.PHONY: seed/up
seed/up: migrate/build
	@echo "Seeding up..."
	go tool goose seed -no-versioning up

.PHONY: seed/reset
seed/reset: migrate/build
	@echo "Seeding reset..."
	go tool goose seed -no-versioning reset

.PHONY: migrate/check
migrate/check:
	@echo "Building migrate tool..."
	@go build -o ./bin/migrate $(MAIN_PACKAGE_PATH)/code_gen/migrations/.
	go tool migrations

.PHONY: tidy
tidy:
	@echo "Tidying up..."
	go fmt ./...
	go tool swag fmt --exclude frontend
	go mod tidy

.PHONY: lint
lint:
	@echo "Linting..."
	go tool golangci-lint run

.PHONY: lint/frontend
lint/frontend:
	@echo "Linting frontend..."
	cd frontend && yarn lint

.PHONY: test/frontend
test/frontend:
	@echo "Testing frontend..."
	cd frontend && yarn test

.PHONY: test
test:
	@echo "Testing..."
	go test -cover ./...

.PHONY: test/verbose
test/verbose:
	@echo "Testing..."
	go test -v -cover ./...

.PHONY: infra/logs
infra/logs:
	@docker compose -f docker-compose.local.yaml logs -n 10 -f || true

.PHONY: infra/stop
infra/stop:
	@docker compose -f docker-compose.local.yaml stop

.PHONY: infra/down
infra/down:
	@docker compose -f docker-compose.local.yaml down --volumes

.PHONY: infra/status
infra/status:
	@docker compose -f docker-compose.local.yaml ps -a --format="table {{.Service}}\t{{.State}}\t{{.Status}}"

.PHONY: debug
debug:
	@echo "Debugging..."
	dlv debug

.PHONY: config
config:
	@echo "TODO: Implement config with sops"
