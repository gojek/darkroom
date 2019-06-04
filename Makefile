APP=darkroom
APP_EXECUTABLE="./out/$(APP)"

all: update-deps ci

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

copy-config:
	cp application.yaml.example application.yaml

ci: copy-config test format vet compile
