# Orthanc CLI - Production Makefile
# =====================================
# A production-grade build system for the Orthanc CLI tool.
#
# Quick Start:
#   make deps      - Download dependencies (first time setup)
#   make build     - Build for your platform
#   make install   - Install to your system
#   make help      - Show all available commands

# Project Configuration
# ---------------------
BINARY_NAME := orthanc
MODULE_NAME := github.com/proencaj/orthanc-cli
MAIN_PACKAGE := ./cmd/orthanc

# Build Configuration
# -------------------
BUILD_DIR := bin
DIST_DIR := dist
COVERAGE_DIR := coverage
INSTALL_DIR := /usr/local/bin

# Go Configuration
# ----------------
GO := go
GOFLAGS := -buildvcs=false
GOTEST_FLAGS := -race -coverprofile=$(COVERAGE_DIR)/coverage.out

# Version Information
# -------------------
# Automatically detect version from git, fallback to "dev" if not in a git repo
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell $(GO) version | awk '{print $$3}')

# Linker Flags
# ------------
# Inject version information into the binary
LDFLAGS := -ldflags "\
	-X main.Version=$(VERSION) \
	-X main.Commit=$(COMMIT) \
	-X main.BuildTime=$(BUILD_TIME) \
	-s -w"

# Platform Detection
# ------------------
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	DETECTED_OS := linux
endif
ifeq ($(UNAME_S),Darwin)
	DETECTED_OS := darwin
endif

# Build Targets
# -------------
PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64

# Color Output (for better UX)
# -----------------------------
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Default Target
# --------------
.DEFAULT_GOAL := help

.PHONY: all
all: clean fmt vet test build ## Run all checks and build (CI/CD)

# Help Target
# -----------
.PHONY: help
help: ## Show this help message
	@echo "$(BLUE)Orthanc CLI - Build System$(NC)"
	@echo ""
	@echo "$(YELLOW)Quick Start for New Users:$(NC)"
	@echo "  $(GREEN)make deps$(NC)       - Download dependencies (run this first!)"
	@echo "  $(GREEN)make build$(NC)      - Build for your current platform"
	@echo "  $(GREEN)make install$(NC)    - Install to your system"
	@echo "  $(GREEN)make test$(NC)       - Run all tests"
	@echo ""
	@echo "$(YELLOW)Available Targets:$(NC)"
	@awk 'BEGIN {FS = ":.*##"; printf ""} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-18s$(NC) %s\n", $$1, $$2 } \
		/^##@/ { printf "\n$(YELLOW)%s$(NC)\n", substr($$0, 5) }' $(MAKEFILE_LIST)

##@ Build Targets

.PHONY: build
build: ## Build binary for current platform
	@echo "$(BLUE)Building $(BINARY_NAME) for $(DETECTED_OS)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "$(GREEN)✓ Binary built: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: build-all
build-all: ## Build binaries for all platforms
	@echo "$(BLUE)Building for all platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		output_name=$(BINARY_NAME)-$$GOOS-$$GOARCH; \
		echo "  $(YELLOW)→$(NC) Building for $$GOOS/$$GOARCH..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH $(GO) build $(GOFLAGS) $(LDFLAGS) \
			-o $(BUILD_DIR)/$$output_name $(MAIN_PACKAGE); \
		if [ $$? -eq 0 ]; then \
			echo "    $(GREEN)✓$(NC) $$output_name"; \
		else \
			echo "    $(RED)✗$(NC) Failed to build $$output_name"; \
			exit 1; \
		fi; \
	done
	@echo ""
	@echo "$(GREEN)✓ All binaries built successfully:$(NC)"
	@ls -lh $(BUILD_DIR)/

.PHONY: dev
dev: ## Quick development build (no version info, faster)
	@echo "$(BLUE)Building development version...$(NC)"
	@$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "$(GREEN)✓ Development build: ./$(BINARY_NAME)$(NC)"

##@ Installation

.PHONY: install
install: build ## Install binary to system (default: /usr/local/bin)
	@echo "$(BLUE)Installing $(BINARY_NAME) to $(INSTALL_DIR)...$(NC)"
	@install -d $(INSTALL_DIR)
	@install -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(GREEN)✓ Installed to $(INSTALL_DIR)/$(BINARY_NAME)$(NC)"
	@echo "  Run '$(BINARY_NAME) --help' to get started"

.PHONY: uninstall
uninstall: ## Remove binary from system
	@echo "$(BLUE)Uninstalling $(BINARY_NAME)...$(NC)"
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(GREEN)✓ Uninstalled$(NC)"

##@ Testing & Quality

.PHONY: test
test: ## Run all tests with race detector
	@echo "$(BLUE)Running tests...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@$(GO) test $(GOTEST_FLAGS) ./...
	@echo "$(GREEN)✓ All tests passed$(NC)"

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo "$(BLUE)Running tests (verbose)...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@$(GO) test -v $(GOTEST_FLAGS) ./...

.PHONY: test-short
test-short: ## Run only short tests (skip long-running tests)
	@echo "$(BLUE)Running short tests...$(NC)"
	@$(GO) test -short ./...
	@echo "$(GREEN)✓ Short tests passed$(NC)"

.PHONY: coverage
coverage: test ## Generate HTML coverage report
	@echo "$(BLUE)Generating coverage report...$(NC)"
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "$(GREEN)✓ Coverage report: $(COVERAGE_DIR)/coverage.html$(NC)"

.PHONY: coverage-view
coverage-view: coverage ## Generate and open coverage report in browser
	@echo "$(BLUE)Opening coverage report...$(NC)"
	@which open >/dev/null 2>&1 && open $(COVERAGE_DIR)/coverage.html || \
		which xdg-open >/dev/null 2>&1 && xdg-open $(COVERAGE_DIR)/coverage.html || \
		echo "$(YELLOW)Please open $(COVERAGE_DIR)/coverage.html manually$(NC)"

.PHONY: bench
bench: ## Run benchmarks
	@echo "$(BLUE)Running benchmarks...$(NC)"
	@$(GO) test -bench=. -benchmem ./...

##@ Code Quality

.PHONY: fmt
fmt: ## Format code with gofmt
	@echo "$(BLUE)Formatting code...$(NC)"
	@$(GO) fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

.PHONY: fmt-check
fmt-check: ## Check if code is formatted
	@echo "$(BLUE)Checking code formatting...$(NC)"
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "$(RED)✗ The following files are not formatted:$(NC)"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	@echo "$(GREEN)✓ All files are properly formatted$(NC)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@$(GO) vet ./...
	@echo "$(GREEN)✓ Vet passed$(NC)"

.PHONY: lint
lint: ## Run golangci-lint (requires golangci-lint)
	@echo "$(BLUE)Running golangci-lint...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(GREEN)✓ Lint passed$(NC)"; \
	else \
		echo "$(RED)✗ golangci-lint not installed$(NC)"; \
		echo "  Install: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

.PHONY: staticcheck
staticcheck: ## Run staticcheck (requires staticcheck)
	@echo "$(BLUE)Running staticcheck...$(NC)"
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
		echo "$(GREEN)✓ Staticcheck passed$(NC)"; \
	else \
		echo "$(RED)✗ staticcheck not installed$(NC)"; \
		echo "  Install: go install honnef.co/go/tools/cmd/staticcheck@latest"; \
		exit 1; \
	fi

.PHONY: check
check: fmt-check vet test ## Run all checks (fmt, vet, test)
	@echo "$(GREEN)✓ All checks passed!$(NC)"

##@ Dependencies

.PHONY: deps
deps: ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@$(GO) mod download
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

.PHONY: deps-update
deps-update: ## Update all dependencies
	@echo "$(BLUE)Updating dependencies...$(NC)"
	@$(GO) get -u ./...
	@$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

.PHONY: deps-verify
deps-verify: ## Verify dependencies
	@echo "$(BLUE)Verifying dependencies...$(NC)"
	@$(GO) mod verify
	@echo "$(GREEN)✓ Dependencies verified$(NC)"

.PHONY: tidy
tidy: ## Tidy go.mod and go.sum
	@echo "$(BLUE)Tidying go modules...$(NC)"
	@$(GO) mod tidy
	@echo "$(GREEN)✓ Modules tidied$(NC)"

##@ Release & Distribution

.PHONY: release
release: clean build-all ## Create release archives for distribution
	@echo "$(BLUE)Creating release archives...$(NC)"
	@mkdir -p $(BUILD_DIR)/release
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		binary_name=$(BINARY_NAME)-$$GOOS-$$GOARCH; \
		archive_name=$(BINARY_NAME)-$(VERSION)-$$GOOS-$$GOARCH.tar.gz; \
		echo "  $(YELLOW)→$(NC) Creating $$archive_name..."; \
		cd $(BUILD_DIR) && tar -czf release/$$archive_name $$binary_name && cd ..; \
	done
	@echo ""
	@echo "$(GREEN)✓ Release archives created:$(NC)"
	@ls -lh $(BUILD_DIR)/release/

.PHONY: checksums
checksums: ## Generate SHA256 checksums for release files
	@echo "$(BLUE)Generating checksums...$(NC)"
	@cd $(BUILD_DIR)/release && \
		shasum -a 256 * > SHA256SUMS.txt && \
		cd ../..
	@echo "$(GREEN)✓ Checksums generated: $(BUILD_DIR)/release/SHA256SUMS.txt$(NC)"

##@ Cleanup

.PHONY: clean
clean: ## Remove build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR) $(DIST_DIR) $(COVERAGE_DIR)
	@rm -f $(BINARY_NAME)
	@$(GO) clean
	@echo "$(GREEN)✓ Clean complete$(NC)"

.PHONY: clean-deps
clean-deps: ## Clean module cache
	@echo "$(BLUE)Cleaning module cache...$(NC)"
	@$(GO) clean -modcache
	@echo "$(GREEN)✓ Module cache cleaned$(NC)"

##@ Running

.PHONY: run
run: build ## Build and run the binary
	@echo "$(BLUE)Running $(BINARY_NAME)...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME)

.PHONY: run-help
run-help: build ## Build and show help
	@$(BUILD_DIR)/$(BINARY_NAME) --help

##@ Information

.PHONY: version
version: ## Show version information
	@echo "$(BLUE)Version Information:$(NC)"
	@echo "  Version:    $(VERSION)"
	@echo "  Commit:     $(COMMIT)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Go Version: $(GO_VERSION)"

.PHONY: info
info: ## Show project information
	@echo "$(BLUE)Project Information:$(NC)"
	@echo "  Project:      Orthanc CLI"
	@echo "  Binary:       $(BINARY_NAME)"
	@echo "  Module:       $(MODULE_NAME)"
	@echo "  Go Version:   $(GO_VERSION)"
	@echo "  Build Dir:    $(BUILD_DIR)"
	@echo "  Install Dir:  $(INSTALL_DIR)"
	@echo "  Main Package: $(MAIN_PACKAGE)"
	@echo "  Platform:     $(DETECTED_OS)"

.PHONY: doctor
doctor: ## Check if all required tools are installed
	@echo "$(BLUE)Checking development environment...$(NC)"
	@echo ""
	@echo "$(YELLOW)Required Tools:$(NC)"
	@command -v go >/dev/null 2>&1 && \
		echo "  $(GREEN)✓$(NC) Go ($(shell go version))" || \
		echo "  $(RED)✗$(NC) Go (not found - install from https://go.dev/dl/)"
	@command -v git >/dev/null 2>&1 && \
		echo "  $(GREEN)✓$(NC) Git ($(shell git --version))" || \
		echo "  $(RED)✗$(NC) Git (not found)"
	@command -v make >/dev/null 2>&1 && \
		echo "  $(GREEN)✓$(NC) Make ($(shell make --version | head -n 1))" || \
		echo "  $(RED)✗$(NC) Make (not found)"
	@echo ""
	@echo "$(YELLOW)Optional Tools:$(NC)"
	@command -v golangci-lint >/dev/null 2>&1 && \
		echo "  $(GREEN)✓$(NC) golangci-lint ($(shell golangci-lint --version 2>&1 | head -n 1))" || \
		echo "  $(YELLOW)○$(NC) golangci-lint (not found - install from https://golangci-lint.run/)"
	@command -v staticcheck >/dev/null 2>&1 && \
		echo "  $(GREEN)✓$(NC) staticcheck" || \
		echo "  $(YELLOW)○$(NC) staticcheck (not found - run: go install honnef.co/go/tools/cmd/staticcheck@latest)"
	@echo ""
	@echo "$(BLUE)Module Status:$(NC)"
	@$(GO) mod verify >/dev/null 2>&1 && \
		echo "  $(GREEN)✓$(NC) Go modules verified" || \
		echo "  $(YELLOW)○$(NC) Run 'make deps' to download dependencies"

##@ Docker (if Dockerfile exists)

.PHONY: docker-build
docker-build: ## Build Docker image
	@if [ -f Dockerfile ]; then \
		echo "$(BLUE)Building Docker image...$(NC)"; \
		docker build -t $(BINARY_NAME):$(VERSION) -t $(BINARY_NAME):latest .; \
		echo "$(GREEN)✓ Docker image built: $(BINARY_NAME):$(VERSION)$(NC)"; \
	else \
		echo "$(YELLOW)No Dockerfile found$(NC)"; \
	fi

.PHONY: docker-run
docker-run: ## Run Docker container
	@if [ -f Dockerfile ]; then \
		echo "$(BLUE)Running Docker container...$(NC)"; \
		docker run --rm -it $(BINARY_NAME):latest; \
	else \
		echo "$(YELLOW)No Dockerfile found$(NC)"; \
	fi

##@ CI/CD

.PHONY: ci
ci: deps check build ## Run CI pipeline (deps, check, build)
	@echo "$(GREEN)✓ CI pipeline completed successfully!$(NC)"

.PHONY: ci-full
ci-full: clean deps check build-all test ## Run full CI pipeline
	@echo "$(GREEN)✓ Full CI pipeline completed successfully!$(NC)"

##@ Tools Installation

.PHONY: install-tools
install-tools: ## Install development tools (golangci-lint, staticcheck)
	@echo "$(BLUE)Installing development tools...$(NC)"
	@echo "  Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "  Installing staticcheck..."
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "$(GREEN)✓ Development tools installed$(NC)"
	@echo "  Make sure $(shell go env GOPATH)/bin is in your PATH"
