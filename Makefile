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

run: copy-config
	@$(GO_RUN) main.go server

compile:
	@mkdir -p $(BUILD_DIR)
	@$(GO_BUILD) -o $(APP_EXECUTABLE) main.go

lint: golint
	@$(GOLINT) ./... | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; }

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

### Helper functions

golint:
	$(call install-if-needed,GOLINT,golang.org/x/lint/golint,v0.0.0-20210508222113-6edffad5e616)

goveralls:
	$(call install-if-needed,GOVERALLS,github.com/mattn/goveralls,v0.0.12)

is-available = $(if $(shell command -v $(1) 2> /dev/null),yes,no)

define install-if-needed
	@if [ $(call is-available,$(notdir $(2))) = "no" ] ; then \
	echo "Installing $(2)..." ;\
	set -e ;\
	TMP_DIR=$$(mktemp -d) ;\
	cd $$TMP_DIR ;\
	go mod init tmp ;\
	go install $(2)@$(3) ;\
	rm -rf $$TMP_DIR ;\
	fi
	@$(eval $1 := $(if $(shell command -v $(notdir $(2))),$(shell command -v $(notdir $(2))),$(GOBIN)/$(notdir $(2))))
endef
