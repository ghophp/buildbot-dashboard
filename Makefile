.PHONY: all get test build

GO ?= go

all: get build test

get:
	${GO} get ./...
	${GO} get -u github.com/jteeuwen/go-bindata/...
	${GOPATH}/bin/go-bindata static/...

build:
	${GO} build

test: get
	${GO} test ./...