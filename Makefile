.PHONY: all deps test clean build

GO ?= go
BIN_NAME=buildbot-dashboard

all: build test

deps:
	${GO} get -u -d github.com/ghophp/buildbot-dashboard
	${GO} get gopkg.in/check.v1
	${GO} get github.com/ghophp/render
	${GO} get github.com/go-martini/martini
	${GO} get github.com/martini-contrib/staticbin
	${GO} get github.com/beatrichartz/martini-sockets
	${GO} get github.com/jteeuwen/go-bindata/...
	${GOPATH}/bin/go-bindata static/...

build: deps
build:
	${GO} build -o ${BIN_NAME}

test: deps
test:
	${GO} test github.com/ghophp/buildbot-dashboard/cache
	${GO} test github.com/ghophp/buildbot-dashboard/config
	${GO} test github.com/ghophp/buildbot-dashboard/container
	${GO} test github.com/ghophp/buildbot-dashboard/handler

clean:
	rm ${BIN_NAME}