.PHONY: all deps test clean build

GO ?= go
BIN_NAME=buildbot-dashboard

all: build test

deps:
	${GO} get -u github.com/motain/gocheck
	${GO} get -u github.com/ghophp/buildbot-dashboard
	${GO} get -u github.com/ghophp/render
	${GO} get -u github.com/go-martini/martini
	${GO} get -u github.com/martini-contrib/staticbin
	${GO} get -u github.com/beatrichartz/martini-sockets
	${GO} get -u github.com/jteeuwen/go-bindata

build: deps
build:
	${GOPATH}/bin/go-bindata static/...
	${GO} build -o ${BIN_NAME}

test: deps
	${GO} test -v github.com/ghophp/buildbot-dashboard/cache
	${GO} test -v github.com/ghophp/buildbot-dashboard/config
	${GO} test -v github.com/ghophp/buildbot-dashboard/container
	${GO} test -v github.com/ghophp/buildbot-dashboard/handler

clean:
	rm ${BIN_NAME}