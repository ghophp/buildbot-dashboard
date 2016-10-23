.PHONY: all deps test clean build

GO ?= go
BIN_NAME=buildbot-dashboard

all: build test

deps:
	${GO} get github.com/op/go-logging
	${GO} get gopkg.in/check.v1
	${GO} get github.com/ghophp/render
	${GO} get github.com/go-martini/martini
	${GO} get github.com/martini-contrib/staticbin
	${GO} get github.com/beatrichartz/martini-sockets
	${GO} get github.com/jteeuwen/go-bindata/...
	${GOPATH}/bin/go-bindata static/...

build: deps
build:
	${GO} build -o bin/latest/${BIN_NAME}
	echo -n bin/latest/${BIN_NAME} | openssl dgst -sha1 -hmac "key"

test: deps
test:
	${GO} test github.com/ghophp/buildbot-dashboard/cache
	${GO} test github.com/ghophp/buildbot-dashboard/config
	${GO} test github.com/ghophp/buildbot-dashboard/handler

clean:
	rm -rf bin/latest/${BIN_NAME}