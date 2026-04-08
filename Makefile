APP=docs-ssot

###############################################################################
# Install
###############################################################################

.PHONY: install
install:
	go install ./cmd/docs-ssot

.PHONY: install-dev
install-dev:
	brew install lefthook
	lefthook install
	go get -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4

###############################################################################
# Development
###############################################################################

APP_VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")

.PHONY: build
build:
	go build -trimpath \
		-ldflags="-s -w -X main.appVersion=$(APP_VERSION)" \
		-o bin/$(APP) ./cmd/docs-ssot

.PHONY: run
run:
	go run ./cmd/docs-ssot build

.PHONY: clean
clean:
	rm -rf bin
	rm -f README.md CLAUDE.md

###############################################################################
# Golang Linting
###############################################################################
.PHONY: go-fmt
go-fmt:
	go tool golangci-lint fmt

# lint
.PHONY: go-lint-check
go-lint-check:
	go tool golangci-lint run

# lint and fix
.PHONY: go-lint
go-lint:
	go tool golangci-lint run --fix

.PHONY: go-lint-fast-check
go-lint-fast-check:
	go tool golangci-lint run --fast-only

.PHONY: go-lint-fast
go-lint-fast:
	go tool golangci-lint run --fast-only --fix

# verify golangci-lint configuration
# Note: run after modifying .golangci.yml
.PHONY: go-lint-verify-config
go-lint-verify-config:
	go tool golangci-lint config verify

# clean golangci-lint cache
.PHONY: go-clean-lint-cache
go-clean-lint-cache:
	go tool golangci-lint cache clean

###############################################################################
# Golang Test
###############################################################################

.PHONY: test
test:
	go test ./...

###############################################################################
# Generate docs
###############################################################################

.PHONY: docs
docs:
	go run ./cmd/docs-ssot build

.PHONY: docs-index
docs-index:
	go run ./cmd/docs-ssot index

.PHONY: docs-validate
docs-validate:
	go run ./cmd/docs-ssot validate

# Usage: make docs-include FILE=template/README.tpl.md
.PHONY: docs-include
docs-include:
	go run ./cmd/docs-ssot include $(FILE)

# Usage: make docs-check
# Usage: make docs-check ARGS="--threshold 0.75 --section-level 3"
.PHONY: docs-check
docs-check:
	go run ./cmd/docs-ssot check --root template/docs $(ARGS)

# Usage: make docs-migrate FILES="README.md CLAUDE.md"
# Usage: make docs-migrate FILES="README.md" ARGS="--dry-run"
.PHONY: docs-migrate
docs-migrate:
	go run ./cmd/docs-ssot migrate $(ARGS) $(FILES)

# Usage: make docs-migrate-agents
# Usage: make docs-migrate-agents ARGS="--dry-run"
# Usage: make docs-migrate-agents ARGS="--target-tools cursor,copilot"
.PHONY: docs-migrate-agents
docs-migrate-agents:
	go run ./cmd/docs-ssot migrate --agents $(ARGS)

.PHONY: docs-version
docs-version:
	go run ./cmd/docs-ssot version

###############################################################################
# VitePress documentation site
###############################################################################

.PHONY: install-docs
install-docs:
	cd docs && bun install

.PHONY: vitepress-dev
vitepress-dev:
	cd docs && bun run dev

.PHONY: vitepress-build
vitepress-build:
	cd docs && bun run build

.PHONY: vitepress-preview
vitepress-preview:
	cd docs && bun run preview

#------------------------------------------------------------------------------
# Release
#------------------------------------------------------------------------------

# update-git-tag: Create and push a git tag for the new version
# e.g. make update-git-tag new=0.1.0
.PHONY: update-git-tag
update-git-tag:
	@echo "Creating git tag v${new}"
	@git tag -a "v${new}" -m "Release version ${new}"
	@echo "Git tag v${new} created"
	@echo "Pushing git tag v${new} to origin"
	@git push origin "v${new}"
	@echo "Git tag v${new} pushed to origin"

COMMIT ?= HEAD

# e.g. make retag TAG=v0.1.0
.PHONY: retag
retag:
	git tag -d $(TAG) 2>/dev/null || true
	git push --delete origin $(TAG) 2>/dev/null || true
	git tag -a $(TAG) $(COMMIT) -m "retag $(TAG)"
	git push origin $(TAG)
