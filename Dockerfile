FROM golang:alpine

ADD . /go/src/github.com/ghophp/buildbot-dashboard

RUN cd /go/src/github.com/ghophp/buildbot-dashboard && GO111MODULE=on go build -o buildbot-dashboard .

CMD ["/go/src/github.com/ghophp/buildbot-dashboard/buildbot-dashboard"]

EXPOSE 8000