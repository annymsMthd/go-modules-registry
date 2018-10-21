VERSION := $(shell gogitver)

.PHONY: build-client build-server test package

build: test build-client build-server

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

build-client:
	go build -o ./artifacts/registry-uploader ./cmd/go-modules-registry-uploader/main.go

build-server:
	docker build . -t annymsmthd/go-modules-registry:test -f ./cmd/go-modules-registry/Dockerfile

package: build
