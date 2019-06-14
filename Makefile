APP=darkroom
APP_EXECUTABLE="./out/$(APP)"

all: ci

setup:
	go get -u golang.org/x/lint/golint
	go get -u github.com/axw/gocov/gocov

compile:
	mkdir -p out
	go build -o $(APP_EXECUTABLE) ./cmd/darkroom/main.go

lint: setup
	golint ./... | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; }

format:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

test-cov:
	gocov test ./... > coverage.json

test-cov-report:
	@echo "\nGENERATING TEST REPORT."
	gocov report coverage.json

copy-config:
	mkdir -p out
	cp config.yaml.example config.yaml

ci: copy-config compile lint format vet test test-cov test-cov-report
