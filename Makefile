.PHONY: all run.docker run.local deps clean test coverage build.docker build.local build.linux.armv8 build.linux.armv7 build.linux build.osx build.windows

ifneq ("$(wildcard .env)","")
ENV_FILE = .env
else
ENV_FILE = .env.example
endif

GOPKGS     = $(shell go list ./... | grep -v /vendor/)
BIN_OUTPUT = bin/buildbot-dashboard

all: build.local test

default: build.local

run.docker: build.docker
	docker run -p 127.0.0.1:8000:8000/tcp --env-file=${ENV_FILE} ghophp/buildbot-dashboard

run.local: deps
	go run $(GOPKGS) main.go

deps:
	go mod vendor

build.docker:
	docker build -t="ghophp/buildbot-dashboard" .

build.local:
	GO111MODULE=on go build -o $(BIN_OUTPUT) .

build.linux.armv8:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GO111MODULE=on go build -o $(BIN_OUTPUT) .

build.linux.armv7:
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 GO111MODULE=on go build -o $(BIN_OUTPUT) .

build.linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o $(BIN_OUTPUT) .

build.osx:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o $(BIN_OUTPUT) .

build.windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build -o $(BIN_OUTPUT) .

test:
	go test -v $(GOPKGS)

coverage:
	go get -u github.com/mattn/goveralls
	./script/coverage.sh --coveralls

clean:
	rm -rf bin
	rm -rf vendor
