SHELL=/bin/bash -o pipefail
$( shell mkdir -p bin )

COVERAGE_PROFILE ?= coverage.out
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
GOLANGCI_VERSION = 1.49.0
TIME_TESTS_PATH = time_tests

ifeq ($(GOARCH),arm)
	ARCH=armv7
else
	ARCH=$(GOARCH)
endif

COMMIT=$(shell git rev-parse --verify HEAD)

###########
# LINTING
###########
bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint

bin/golangci-lint-${GOLANGCI_VERSION}:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

.PHONY: lint fix
lint: bin/golangci-lint
	bin/golangci-lint run

fix: bin/golangci-lint
	bin/golangci-lint run --fix

###########
# TESTING
###########
test:
	go test -v ./... -covermode=count -coverprofile=${COVERAGE_PROFILE}
	go tool cover -func=${COVERAGE_PROFILE} -o=${COVERAGE_PROFILE}

time-tests:
	cd ${TIME_TESTS_PATH} && $(MAKE) generate-time-table
