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

.PHONY: build
build:
	go build -o bin/$(APP) ./cmd/docs-ssot

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