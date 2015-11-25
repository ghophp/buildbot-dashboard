.PHONY: all deps special-deps test semaphore-test clean build semaphore-build

GO ?= go
BIN_NAME=buildbot-dashboard

all: build test

special-deps:
	${GO} get github.com/ghophp/buildbot-dashboard

deps:
	${GO} get gopkg.in/check.v1
	${GO} get github.com/ghophp/render
	${GO} get github.com/go-martini/martini
	${GO} get github.com/martini-contrib/staticbin
	${GO} get github.com/beatrichartz/martini-sockets
	${GO} get github.com/jteeuwen/go-bindata/...
	${GOPATH}/bin/go-bindata static/...

semaphore-build: deps
semaphore-build:
	${GO} build -o ${BIN_NAME}

semaphore-test: deps
semaphore-test: deps
	${GO} test github.com/ghophp/buildbot-dashboard/cache
	${GO} test github.com/ghophp/buildbot-dashboard/config
	${GO} test github.com/ghophp/buildbot-dashboard/container
	${GO} test github.com/ghophp/buildbot-dashboard/handler

build: special-deps semaphore-build

test: special-deps semaphore-test

clean:
	rm ${BIN_NAME}