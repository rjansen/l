NAME 			:= l
BIN         	:= $(NAME)
REPO        	:= github.com/rjansen/$(NAME)
BUILD       	:= $(shell openssl rand -hex 10)
VERSION     	:= $(shell if [ -f version ]; then awk '{printf $1}' < version; else openssl rand -hex 5; fi)
MAKEFILE    	:= $(word $(words $(MAKEFILE_LIST)), $(MAKEFILE_LIST))
BASE_DIR    	:= $(shell cd $(dir $(MAKEFILE)); pwd)
PKGS        	:= $(shell go list ./... | grep -v /vendor/)

#Test and Benchmark Parameters
TEST_PKGS ?= 
BENCHS ?= .
COVERAGE_FILE := $(NAME).coverage
COVERAGE_HTML := $(NAME).coverage.html
PKG_COVERAGE := $(NAME).pkg.coverage

.PHONY: default
default: build

.PHONY: install
install: deps
	@echo "$(REPO) installed successfully" 

.PHONY: deps
deps:
	brew install go
	go get -u github.com/kardianos/govendor
	govendor sync

.PHONY: sync
sync:
	govendor sync

.PHONY: all
all: build test bench coverage

.PHONY: build
build:
	@echo "Building $(REPO)@$(VERSION)-$(BUILD)"
	go build $(PKGS)

.PHONY: clean
clean: 
	-rm $(NAME)*coverage*
	-rm *.test
	-rm *.pprof

.PHONY: reset
reset: clean 
	-cd vendor; rm -r */

.PHONY: test_all
test_all:
	go test -v -race $(PKGS)

.PHONY: test
test:
	@if [ "$(TEST_PKGS)" == "" ]; then \
	    echo "Test All Pkgs";\
		go test -v -race $(PKGS) || exit 501;\
	else \
	    echo "Test Selected Pkgs=$(TEST_PKGS)";\
		SELECTED_TEST_PKGS="";\
	    for tstpkg in $(TEST_PKGS); do \
			go test -v -race $(REPO)/$$tstpkg || exit 501;\
		done; \
	fi

.PHONY: bench_all
bench_all:
	#go test -bench=. -run="^$$" -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem $(PKGS)
	go test -bench=. -run="^$$" -benchmem $(PKGS)

.PHONY: bench
bench:
	@if [ "$(TEST_PKGS)" == "" ]; then \
	    echo "Benchmark all Pkgs" ;\
	    for tstpkg in $(PKGS); do \
		    go test -bench=$(BENCHS) -run="^$$" -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem $$tstpkg || exit 501;\
		done; \
	else \
	    echo "Benchmark Selected Pkgs=$(TEST_PKGS)" ;\
	    for tstpkg in $(TEST_PKGS); do \
		    go test -bench=$(BENCHS) -run="^$$" -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem $(REPO)/$$tstpkg || exit 501;\
		done; \
	fi

.PHONY: benchcmp
benchcmp: 
ifndef BENCH_BEFORE
	@echo "You must define the BENCH_BEFORE variable!"
	@exit 1
endif
ifndef BENCH_AFTER
	@echo "You must define the BENCH_AFTER variable!"
	@exit 1
endif
	benchcmp $(BENCH_BEFORE) $(BENCH_AFTER)

.PHONY: coverage
coverage:
	@echo "Testing with coverage"
	@echo 'mode: set' > $(COVERAGE_FILE)
	@touch $(PKG_COVERAGE)
	@touch $(COVERAGE_FILE)
	@if [ "$(TEST_PKGS)" == "" ]; then \
		for pkg in $(PKGS); do \
			go test -v -coverprofile=$(PKG_COVERAGE) $$pkg || exit 501; \
			grep -v 'mode: set' $(PKG_COVERAGE) >> $(COVERAGE_FILE); \
		done; \
	else \
	    echo "Testing with covegare the Pkgs=$(TEST_PKGS)" ;\
	    for tstpkg in $(TEST_PKGS); do \
			go test -v -coverprofile=$(PKG_COVERAGE) $(REPO)/$$tstpkg || exit 501; \
			grep -v 'mode: set' $(PKG_COVERAGE) >> $(COVERAGE_FILE); \
		done; \
	fi
	@echo "Generating report"
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	open $(COVERAGE_HTML)
