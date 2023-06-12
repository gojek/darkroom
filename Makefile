BUILD_DIR := ".out"
APP_EXECUTABLE="$(BUILD_DIR)/darkroom"

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

BUILD_INFO_GIT_TAG ?= $(shell git describe --tags 2>/dev/null || echo unknown)
BUILD_INFO_GIT_COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo unknown)
BUILD_INFO_BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" || echo unknown)
BUILD_INFO_VERSION ?= $(shell prefix=$$(echo $(BUILD_INFO_GIT_TAG) | cut -c 1); if [ "$${prefix}" = "v" ]; then echo $(BUILD_INFO_GIT_TAG) | cut -c 2- ; else echo $(BUILD_INFO_GIT_TAG) ; fi)

build_info_fields := \
	version=$(BUILD_INFO_VERSION) \
	gitTag=$(BUILD_INFO_GIT_TAG) \
	gitCommit=$(BUILD_INFO_GIT_COMMIT) \
	buildDate=$(BUILD_INFO_BUILD_DATE)
build_info_ld_flags := $(foreach entry,$(build_info_fields),-X github.com/gojek/darkroom/internal/version.$(entry))

LD_FLAGS := -ldflags="-s -w $(build_info_ld_flags)"
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GO_BUILD := GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=1 go build $(LD_FLAGS)
GO_RUN := GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=1 go run $(LD_FLAGS)

all: test-ci

setup:
	go get golang.org/x/lint/golint

run: copy-config
	@$(GO_RUN) main.go server

compile:
	@mkdir -p $(BUILD_DIR)
	@$(GO_BUILD) -o $(APP_EXECUTABLE) main.go

lint:
	@golint ./... | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; }

format:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./... -covermode=count -coverprofile=profile.cov

coverage: goveralls
	@$(GOVERALLS) -coverprofile=profile.cov -service=github

copy-config:
	@cp config.example.yaml config.yaml

docker-image:
	@docker build -t ${USER}/darkroom:latest -f build/Dockerfile .

docker-docs:
	@docker build -t darkroom-docs:latest -f build/Dockerfile.docs .

test-ci: copy-config compile lint format vet test

# find or download goveralls
goveralls:
ifeq (, $(shell which goveralls))
	@{ \
	set -e ;\
	GOVERALLS_TMP_DIR=$$(mktemp -d) ;\
	cd $$GOVERALLS_TMP_DIR ;\
	go mod init tmp ;\
	go install github.com/mattn/goveralls ;\
	rm -rf $$GOVERALLS_TMP_DIR ;\
	}
GOVERALLS=$(GOBIN)/goveralls
else
GOVERALLS=$(shell which goveralls)
endif
