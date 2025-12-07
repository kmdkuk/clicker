LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

.PHONY: build-wasm
build-wasm:
	GOOS=js GOARCH=wasm go build -o pages/main.wasm cmd/clicker/main.go

build: ## Build the Go application.
	go build -o bin/clicker cmd/clicker/main.go

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet staticcheck ginkgo lint## Run tests.
	$(STATICCHECK) ./...
	$(GINKGO) -p -v -r --trace --cover --coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

GOLANGCI_LINT = $(LOCALBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v2.7.2 # renovate: golangci/golangci-lint
golangci-lint:
	@[ -f $(GOLANGCI_LINT) ] || { \
	set -e ;\
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell dirname $(GOLANGCI_LINT)) $(GOLANGCI_LINT_VERSION) ;\
	}

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter & yamllint
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix

STATICCHECK ?= $(LOCALBIN)/staticcheck
.PHONY: staticcheck
staticcheck: $(STATICCHECK)
$(STATICCHECK): $(LOCALBIN)
	test -s $(LOCALBIN)/staticcheck || GOBIN=$(LOCALBIN) go install honnef.co/go/tools/cmd/staticcheck@latest

GINKGO ?= $(LOCALBIN)/ginkgo
.PHONY: ginkgo
ginkgo: $(GINKGO)
$(GINKGO): $(LOCALBIN)
	test -s $(GINKGO) || GOBIN=$(LOCALBIN) go install github.com/onsi/ginkgo/v2/ginkgo@latest
