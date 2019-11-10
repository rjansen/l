GO_MODULE          := $(REPO)/$(NAME)
GO_VERSION         := 1.13.1
export GO111MODULE := on
GO_PKGS            := ./...
GO_TEST_PKGS       := $(if $(GO_TEST_PKGS),$(addprefix $(GO_MODULE)/,$(GO_TEST_PKGS)),$(GO_PKGS))
GO_TESTS           ?= .
GO_TEST_VERBOSITY  ?= $(if $(DEBUG),standard-verbose,short-verbose)
GO_COVERAGE_FILE   := $(TMP_DIR)/$(subst /,_,$(NAME)).coverage
GO_COVERAGE_HTML   := $(TMP_DIR)/$(subst /,_,$(NAME)).coverage.html

.PHONY: go.clearcache
go.clearcache:
	@echo "$(NAME)@$(BUILD) go.clearcache"
	-rm -f $(GO_COVERAGE_FILE)
	-rm -f $(GO_COVERAGE_HTML)
	-rm -rf vendor

.PHONY: go.clearmodules
go.clearmodules:
	@echo "$(NAME)@$(BUILD) go.clearmodules"
	-rm -Rf $(HOME)/.cache/go-build
	-rm -Rf $(HOME)/gopath/pkg/mod

.PHONY: go.vendor
go.vendor:
	@echo "$(NAME)@$(BUILD) go.vendor"
	go mod verify && go mod vendor

.PHONY: go.vet
go.vet:
	@echo "$(NAME)@$(BUILD) go.vet"
	go vet $(GO_TEST_PKGS)

.PHONY: go.lint
go.lint:
	@echo "$(NAME)@$(BUILD) go.lint"
	golint $(shell go list $(GO_TEST_PKGS))

.PHONY: go.test
go.test:
	@echo "$(NAME)@$(BUILD) go.test $(GO_TEST_PKGS)"
	$(GOTESTSUM) -f $(GO_TEST_VERBOSITY) -- -v -race -count 1 -run $(GO_TESTS) $(GO_TEST_PKGS)

.PHONY: go.itest
go.itest:
	@echo "$(NAME)@$(BUILD) go.itest $(GO_TEST_PKGS)"
	$(GOTESTSUM) -f $(GO_TEST_VERBOSITY) -- -tags="integration" -v -race -count 1 -run $(GO_TESTS) $(GO_TEST_PKGS)

.PHONY: go.bench
go.bench:
	@echo "$(NAME)@$(BUILD) go.bench $(GO_TEST_PKGS)"
	$(GOTESTSUM) -f $(TEST_VERBOSITY) -- -bench=. -run="^$$" -benchmem $(GO_TEST_PKGS)

.PHONY: go.coverage
go.coverage: $(TMP_DIR)
	@echo "$(NAME)@$(BUILD) go.coverage $(GO_TEST_PKGS)"
	@touch $(GO_COVERAGE_FILE)
	$(GOTESTSUM) -f $(GO_TEST_VERBOSITY) -- -tags="integration" -v -race -count 1 -run $(GO_TESTS) \
		-covermode=atomic -coverpkg=$(GO_PKGS) -coverprofile=$(GO_COVERAGE_FILE) $(GO_TEST_PKGS)

.PHONY: go.coverage.text
go.coverage.text: go.coverage
	@echo "$(NAME)@$(BUILD) go.coverage.text"
	go tool cover -func=$(GO_COVERAGE_FILE)

.PHONY: go.coverage.html
go.coverage.html: go.coverage
	@echo "$(NAME)@$(BUILD) go.coverage.html"
	go tool cover -html=$(GO_COVERAGE_FILE) -o $(GO_COVERAGE_HTML)
