.PHONY: all get test build

GO ?= go
WILDCARD ?= ...

all: get build test

get:
	${GO} get ./${WILDCARD}
	${GO} get -u github.com/jteeuwen/go-bindata/${WILDCARD}
	${GOPATH}/bin/go-bindata static/${WILDCARD}

build:
	${GO} build

test: get
	${GO} test ./${WILDCARD}