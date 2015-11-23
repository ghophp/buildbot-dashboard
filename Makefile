.PHONY: all get test build

GO ?= go

all: get build test

get:
	${GO} get ./...
	go get -u github.com/jteeuwen/go-bindata/...
	exec ${GOPATH}/bin/go-bindata static/...

build:
	${GO} build

test: get
	${GO} test ./...