include Makefile.vars

.PHONY: ci
ci: vet coverage.text bench

.PHONY: clean
clean:
	@echo "$(REPO) clean"
	-rm $(NAME)*coverage* > /dev/null 2>&1
	-rm *.test > /dev/null 2>&1
	-rm *.pprof > /dev/null 2>&1

.PHONY: clearcache
clearcache:
	@echo "$(REPO) clearcache"
	-rm -Rf $(BASE_DIR)/on > /dev/null 2>&1
	-rm -Rf $(BASE_DIR)/vendor > /dev/null 2>&1
	-rm -Rf $(TMP_DIR) > /dev/null 2>&1

.PHONY: install.gvm
install.gvm:
	@echo "$(REPO) install.gvm"
	which gvm || \
		curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer | bash

.PHONY: install
install: deps
	@echo "$(REPO) install"

$(TMP_DIR):
	mkdir -p $(TMP_DIR)

.PHONY: deps
deps: $(TMP_DIR)
	@echo "$(REPO) deps"
	which gotestsum || (\
		cd $(TMP_DIR) && \
		curl -O -L https://github.com/gotestyourself/gotestsum/releases/download/v0.3.2/gotestsum_0.3.2_linux_amd64.tar.gz && \
		tar xf gotestsum_0.3.2_linux_amd64.tar.gz && \
		mv -f gotestsum /usr/local/bin \
	)
	gotestsum --help > /dev/null 2>&1
	which codecov || (\
		cd $(TMP_DIR) && \
		curl -L -o codecov https://codecov.io/bash && \
		chmod a+x codecov && \
		mv -f codecov /usr/local/bin \
	)
	codecov -h > /dev/null 2>&1

.PHONY: debug.deps
debug.deps:
	@echo "$(REPO) deps"
	which dlv || \
		go get -u github.com/derekparker/delve/cmd/dlv
	dlv version

.PHONY: vendor
vendor:
	@echo "$(REPO) vendor"
	go mod vendor
	go mod verify

.PHONY: fmt
fmt:
	@echo "$(REPO) fmt"
	go fmt $(PKGS)

.PHONY: vet
vet:
	@echo "$(REPO) vet: $(PKGS)"
	go vet $(PKGS)

.PHONY: debug
debug:
	@echo "$(REPO) debug"
	dlv debug $(REPO)

.PHONY: debugtest
debugtest:
	@echo "$(REPO) debugtest"
ifeq ($(TEST_PKGS),)
	@echo "debugtest: pkgs=*"
	dlv test --build-flags='$(REPO)' -- -test.run $(TESTS)
else
	@echo "debugtest: pkgs=$(TEST_PKGS)"
	dlv test --build-flags='$(REPO)/$(TEST_PKGS)' -- -test.run $(TESTS)
endif

.PHONY: test
test:
	@echo "$(REPO) test"
ifeq ($(TEST_PKGS),)
	@echo "test: pkgs=*"
	gotestsum -f short-verbose -- -v -race -run $(TESTS) $(PKGS)
else
	@echo "test: pkgs=$(TEST_PKGS)"
	$(foreach pkg,$(TEST_PKGS),\
		gotestsum -f short-verbose -- -v -race -run $(TESTS) $(REPO)/$(pkg);\
	)
endif

.PHONY: itest
itest:
	@echo "$(REPO) itest"
ifeq ($(TEST_PKGS),)
	@echo "itest: pkgs=*"
	gotestsum -f short-verbose -- -tags=integration -v -race -run $(TESTS) $(PKGS)
else
	@echo "itest: pkgs=$(TEST_PKGS)"
	$(foreach pkg,$(TEST_PKGS),\
		gotestsum -f short-verbose -- -tags=integration -v -race -run $(TESTS) $(REPO)/$(pkg);\
	)
endif

.PHONY: bench
bench:
	@echo "$(REPO) bench"
ifeq ($(TEST_PKGS),)
	@echo "bench: pkgs=*"
	gotestsum -f short-verbose -- -bench=. -run="^$$" -benchmem $(PKGS)
else
	@echo "bench: pkgs=$(TEST_PKGS)"
	$(foreach pkg,$(TEST_PKGS),\
		gotestsum -f short-verbose -- -bench=. -run="^$$" -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem $(REPO)/$(pkg);\
	)
endif

.PHONY: coverage
coverage:
	@echo "$(REPO) coverage"
	@touch $(COVERAGE_FILE)
ifeq ($(TEST_PKGS),)
	@echo "coverage: pkgs=*"
	gotestsum -f short-verbose -- -tags=integration -v -run $(TESTS) \
			  -covermode=atomic -coverpkg=./... -coverprofile=$(COVERAGE_FILE) $(PKGS)
else
	@echo "bench: pkgs=$(TEST_PKGS)"
	@touch $(COVERAGE_PKG_FILE)
	@echo 'mode: atomic' > $(COVERAGE_FILE)
	$(foreach pkg,$(TEST_PKGS),\
		gotestsum -f short-verbose -- -tags=integration -v -run $(TESTS) \
				  -covermode=atomic -coverpkg=./... -coverprofile=$(COVERAGE_PKG_FILE) $(REPO)/$(pkg);\
		grep -v 'mode: atomic' $(COVERAGE_PKG_FILE) >> $(COVERAGE_FILE);\
	)
endif

.PHONY: coverage.text
coverage.text: coverage
	@echo "$(REPO) coverage.text"
	go tool cover -func=$(COVERAGE_FILE)

.PHONY: coverage.html
coverage.html: coverage
	@echo "$(REPO) coverage.html"
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	open $(COVERAGE_HTML) || google-chrome $(COVERAGE_HTML) || google-chrome-stable $(COVERAGE_HTML)

.PHONY: coverage.push
coverage.push:
	@echo "$(REPO) coverage.push"
	@#curl -sL https://codecov.io/bash | bash -s - -f $(COVERAGE_FILE)$(if $(CODECOV_TOKEN), -t $(CODECOV_TOKEN),)
	@codecov -f $(COVERAGE_FILE)$(if $(CODECOV_TOKEN), -t $(CODECOV_TOKEN),)

.PHONY: docker
docker.build:
	@echo "$(REPO)@$(BUILD) docker"
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) \
		         -t $(DOCKER_NAME) -t $(DOCKER_NAME):$(VERSION) -f ./etc/docker/Dockerfile .

.PHONY: docker.bash
docker.bash:
	@echo "$(REPO)@$(BUILD) docker.bash"
	docker run --rm --name $(NAME)-bash --entrypoint bash -it -u $(shell id -u):$(shell id -g) \
			   -v `pwd`:/go/src/$(REPO) $(DOCKER_NAME)

docker.%:
	@echo "$(REPO)@$(BUILD) docker.$*"
	@docker run --rm --name $(NAME)-run -u $(shell id -u):$(shell id -g) \
    		    -v `pwd`:/go/src/$(REPO) $(DOCKER_NAME) $*
