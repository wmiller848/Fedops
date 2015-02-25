# Docker build for Fedops
FROM golang
MAINTAINER W. Chase Miller

USER root

ADD . /go/src/github.com/FedOps

RUN go get code.google.com/p/go.crypto/ssh
RUN go get github.com/codegangsta/cli
RUN go get code.google.com/p/gopass

RUN go install github.com/FedOps/fedops-truck
RUN go install github.com/FedOps/fedops-warehouse
RUN go install github.com/FedOps/fedops

VOLUME /opt/fedops