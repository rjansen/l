GOLINT_VERSION := master
GOLINT_URL     := golang.org/x/lint/golint@$(GOLINT_VERSION)
GOLINT_BIN     := golint
GOLINT         := $(or $(shell which $(GOLINT_BIN)),$(GOLINT_BIN))

.PHONY: golint.install
golint.install: $(GOLINT)
	@echo "$(NAME)@$(BUILD) golint.install $(GOLINT)"

.PHONY: golint.uninstall
golint.uninstall:
	@echo "$(NAME)@$(BUILD) golint.uninstall $(GOTESTSUM)"
	-rm -rf $(GOLINT)

$(GOLINT):
	@echo "$(NAME)@$(BUILD) $(GOLINT_URL)"
	go get -u $(GOLINT_URL)
	which $(GOLINT) > /dev/null 2>&1
