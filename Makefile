APP=darkroom
APP_EXECUTABLE="./out/$(APP)"

all: build test

build-deps:
	go install

update-deps: build-deps

compile:
	mkdir -p out
	go build -o $(APP_EXECUTABLE)

format:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...
