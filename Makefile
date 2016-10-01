NAME 		:= logger
BIN         := $(NAME)
REPO        := farm.e-pedion.com/repo/$(NAME)
BUILD       := $(shell git rev-parse --short HEAD)
#VERSION     := $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
MAKEFILE    := $(word $(words $(MAKEFILE_LIST)), $(MAKEFILE_LIST))
BASE_DIR    := $(shell cd $(dir $(MAKEFILE)); pwd)
PKGS        := $(shell go list)
COVERAGE_FILE   := $(NAME).coverage
COVERAGE_HTML  	:= $(NAME).coverage.html
PKG_COVERAGE   	:= $(NAME).pkg.coverage

ETC_DIR := $(BASE_DIR)/etc
CONF_DIR := $(ETC_DIR)/$(NAME)
CONF := $(CONF_DIR)/$(NAME).yml

ENV := local
TEST_PKGS := 

.PHONY: default
default: build

.PHONY: install
install: install_sw_deps sync
	@echo "Logger installed successfully" 

.PHONY: install_sw_deps
install_sw_deps:
	brew install go
	go get -u github.com/kardianos/govendor

.PHONY: install_deps
install_deps:
	go get -u github.com/op/go-logging
	go get -u github.com/Sirupsen/logrus
	go get -u github.com/uber-go/zap

.PHONY: all
all: build test bench_all coverage

.PHONY: build
build:
	go build $(REPO)

.PHONY: sync
sync:
	govendor sync

.PHONY: reset
reset: 
	-rm $(NAME)*coverage*
	-cd vendor; rm -r */

.PHONY: local
local: 
	@echo "Set enviroment to local"
	$(eval ENV = "local")

.PHONY: dev
dev: 
	@echo "Set enviroment to dev"
	$(eval ENV = "dev")

.PHONY: prod
prod: 
	@echo "Set enviroment to prod"
	$(eval ENV = "prod")

.PHONY: check_env
check_env:
	@if [ "$(ENV)" == "" ]; then \
	    echo "Env is blank: $(ENV)"; \
	    exit 540; \
	fi

.PHONY: filter_conf
filter_conf: check_env
	@echo "Filtering Conf Env=$(ENV)"
	@source $(CONF_DIR)/$(NAME).$(ENV).etv && eval "echo \"`cat $(CONF_DIR)/$(NAME).etf`\"" > $(CONF)

.PHONY: check_conf
check_conf:
	@if [ ! -f $(CONF) ]; then \
	    echo "Config file: $(CONF) not found for Env: $(ENV)"; \
	    exit 541; \
	fi

.PHONY: test_loop
test_loop:
	@if [ "$(TEST_PKGS)" == "" ]; then \
	    echo "Test All Pkgs";\
	    for pkg in $(PKGS); do \
			go test -v -race $$pkg || exit 501;\
		done; \
	else \
	    echo "Test Selected Pkgs=$(TEST_PKGS)";\
	    for tstpkg in $(TEST_PKGS); do \
		    go test -v -race $(REPO)/$$tstpkg || exit 501;\
		done; \
	fi

.PHONY: test_all
test_all:
	go test -v -race 

.PHONY: test
test:
	@if [ "$(TEST_PKGS)" == "" ]; then \
	    echo "Test All Pkgs";\
		go test -v -race || exit 501;\
	else \
	    echo "Test Selected Pkgs=$(TEST_PKGS)";\
		SELECTED_TEST_PKGS="";\
	    for tstpkg in $(TEST_PKGS); do \
			go test -v -race $(REPO)/$$tstpkg || exit 501;\
		done; \
	fi

.PHONY: bench_all
bench_all:
	go test -bench=. -v -race 

.PHONY: bench
bench:
	@if [ "$(TEST_PKGS)" == "" ]; then \
	    echo "Bench All Pkgs" ;\
		go test -bench=. -v -race || exit 501;\
	else \
	    echo "Test Selected Pkgs=$(TEST_PKGS)" ;\
	    for tstpkg in $(TEST_PKGS); do \
		    go test -bench=. -v -race $(REPO)/$$tstpkg || exit 501;\
		done; \
	fi

.PHONY: coverage
coverage:
	@echo "Running tests with coverage report..."
	@echo 'mode: set' > $(COVERAGE_FILE)
	@touch $(PKG_COVERAGE)
	@touch $(COVERAGE_FILE)
	@if [ "$(TEST_PKGS)" == "" ]; then \
		for pkg in $(PKGS); do \
			go test -v -coverprofile=$(PKG_COVERAGE) $$pkg || exit 501; \
			grep -v 'mode: set' $(PKG_COVERAGE) >> $(COVERAGE_FILE); \
		done; \
	else \
	    echo "Covegare Test Selected Pkgs=$(TEST_PKGS)" ;\
	    for tstpkg in $(TEST_PKGS); do \
			go test -v -coverprofile=$(PKG_COVERAGE) $(REPO)/$$tstpkg || exit 501; \
			grep -v 'mode: set' $(PKG_COVERAGE) >> $(COVERAGE_FILE); \
		done; \
	fi
	@echo "Generating HTML report in $(COVERAGE_HTML)..."
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@(which -s open && open $(COVERAGE_HTML)) || (which -s gnome-open && gnome-open $(COVERAGE_HTML)) || (exit 0)
