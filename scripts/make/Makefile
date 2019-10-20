TOOLS_DIR             ?= /usr/local/bin
OS        	          ?= darwin_amd64
include scripts/make/go/gotestsum/Makefile
include scripts/make/go/lint/Makefile
include scripts/make//go/lib/Makefile
CODECOV_COVERAGE_FILE := $(GO_COVERAGE_FILE)
include scripts/make/codecov/Makefile

.PHONY: ci
ci: go.vet go.coverage.text
