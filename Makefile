export GOPATH=$(HOME)/go
export GOBIN=$(PWD)/bin
export GOOS=linux
export GOARCH=amd64

.PHONY: build

GIT_COMMIT := $(shell git rev-list -1 HEAD)

build: 
	go get .
	@echo "Building topic-state to ./build"
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X main.GitCommit=$(GIT_COMMIT)" -o build/topic-state .

clean:
	rm -rf ./build ./bin
