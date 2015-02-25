# Docker build for Fedops
FROM golang
MAINTAINER W. Chase Miller

USER root

ADD . /go/src/github.com/Fedops

RUN go get code.google.com/p/go.crypto/ssh
RUN go get github.com/codegangsta/cli
RUN go get code.google.com/p/gopass

RUN go install github.com/Fedops/fedops-truck
RUN go install github.com/Fedops/fedops-warehouse
RUN go install github.com/Fedops/fedops

VOLUME /opt/fedops