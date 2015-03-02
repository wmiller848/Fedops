# Docker build for Fedops
FROM golang
MAINTAINER W. Chase Miller

USER root

ADD . /go/src/github.com/Fedops

RUN go get -u github.com/codegangsta/cli
RUN go get -u code.google.com/p/gopass
RUN go get -u golang.org/x/crypto/ssh
RUN go get -u github.com/pkg/sftp

RUN go install github.com/Fedops/fedops
RUN go install github.com/Fedops/fedops-warehouse
RUN go install github.com/Fedops/fedops-truck

VOLUME /opt/fedops
EXPOSE 13371
WORKDIR /opt/fedops