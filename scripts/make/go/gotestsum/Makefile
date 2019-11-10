GOTESTSUM_TMP_DIR  := $(TMP_DIR)/gotestsum
GOTESTSUM_VERSION  := 0.3.5
# TODO: Create a map between os + arch filename and our scripts convention
GOTESTSUM_NAME     := gotestsum_$(GOTESTSUM_VERSION)_$(OS)_$(ARCH).tar.gz
GOTESTSUM_FILE     := $(GOTESTSUM_TMP_DIR)/$(GOTESTSUM_NAME)
GOTESTSUM_URL      := https://github.com/gotestyourself/gotestsum/releases/download/v$(GOTESTSUM_VERSION)/$(GOTESTSUM_NAME)
GOTESTSUM_BIN      := gotestsum
GOTESTSUM          := $(or $(shell which $(GOTESTSUM_BIN)),$(GOTESTSUM_TMP_DIR)/$(GOTESTSUM_BIN))

.PHONY: gotestsum.clearcache
gotestsum.clearcache:
	@echo "$(NAME)@$(BUILD) gotestsum.clearcache"
	-rm -Rf $(GOTESTSUM_TMP_DIR)

$(GOTESTSUM_TMP_DIR):
	mkdir -p $(GOTESTSUM_TMP_DIR)

.PHONY: gotestsum.install
gotestsum.install: $(GOTESTSUM)
	@echo "$(NAME)@$(BUILD) gotestsum.install $(GOTESTSUM)"

.PHONY: gotestsum.uninstall
gotestsum.uninstall:
	@echo "$(NAME)@$(BUILD) gotestsum.uninstall $(GOTESTSUM)"
	-rm -f $(GOTESTSUM)

$(GOTESTSUM_FILE): | $(GOTESTSUM_TMP_DIR)
	@echo "$(NAM)@$(BUILD) $(GOTESTSUM_FILE)"
	curl -o $(GOTESTSUM_FILE) -L $(GOTESTSUM_URL)

$(GOTESTSUM): | $(GOTESTSUM_FILE)
	@echo "$(NAME)@$(BUILD) $(GOTESTSUM)"
	cd $(GOTESTSUM_TMP_DIR) && tar xf $(GOTESTSUM_NAME)
	$(GOTESTSUM) --help > /dev/null 2>&1
