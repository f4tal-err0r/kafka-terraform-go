export GOPATH=$(HOME)/go

.PHONY: build

GIT_COMMIT := $(shell git rev-list -1 HEAD)

build: 
	@echo "Building topic-state to ./build"
	go build -ldflags "-X main.GitCommit=$(GIT_COMMIT)" -o build/topic-state .

clean:
	rm -rf ./build ./bin
