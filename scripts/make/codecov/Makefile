CODECOV_TMP_DIR    := $(TMP_DIR)/codecov
CODECOV_NAME       := codecov.$(OS)_$(ARCH)
CODECOV_FILE       := $(CODECOV_TMP_DIR)/$(CODECOV_NAME)
CODECOV_URL        := https://codecov.io/bash
CODECOV_BIN        := codecov
CODECOV            := $(or $(shell which $(CODECOV_BIN)),$(CODECOV_TMP_DIR)/$(CODECOV_BIN))

.PHONY: codecov.clearcache
codecov.clearcache:
	@echo "$(NAME)@$(BUILD) codecov.clearcache"
	-rm -Rf $(CODECOV_TMP_DIR)

$(CODECOV_TMP_DIR):
	mkdir -p $(CODECOV_TMP_DIR)

.PHONY: codecov.install
codecov.install: $(CODECOV)
	@echo "$(NAME)@$(BUILD) codecov.install $(CODECOV)"

.PHONY: codecov.uninstall
codecov.uninstall:
	@echo "$(NAME)@$(BUILD) codecov.uninstall $(CODECOV)"
	-rm -f $(CODECOV)

$(CODECOV_FILE): | $(CODECOV_TMP_DIR)
	@echo "$(NAME)@$(BUILD) $(CODECOV_FILE)"
	curl -o $(CODECOV_FILE) -L $(CODECOV_URL)

$(CODECOV): | $(CODECOV_FILE)
	@echo "$(NAME)@$(BUILD) $(CODECOV)"
	cd $(CODECOV_TMP_DIR) && chmod a+x $(CODECOV_NAME) && cp -vf $(CODECOV_NAME) $(CODECOV)
	$(CODECOV) -h > /dev/null 2>&1

.PHONY: codecov.push
codecov.push:
	@echo "$(NAME) codecov.push"
	$(CODECOV) -f $(CODECOV_COVERAGE_FILE)$(if $(CODECOV_TOKEN), -t $(CODECOV_TOKEN),)
