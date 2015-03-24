# Docker build for Fedops
FROM fedora:latest
MAINTAINER W. Chase Miller

USER root

RUN yum -y install glibc-devel gcc make cmake git hostname libgit2-devel pcre-devel

WORKDIR /opt
RUN git clone https://github.com/golang/go
WORKDIR /opt/go/src

RUN git checkout go1.4.1
RUN ./all.bash

ENV PATH="$PATH:/opt/go/bin"

ENV GOPATH="/go"
RUN mkdir $GOPATH

ENV PATH="$PATH:$GOPATH"

RUN go get golang.org/x/tools/cmd/...
ADD . /go/src/github.com/wmiller848/Fedops

RUN go get -u gopkg.in/libgit2/git2go.v22
RUN go get -u github.com/codegangsta/cli
RUN go get -u code.google.com/p/gopass
RUN go get -u golang.org/x/crypto/bcrypt
RUN go get -u golang.org/x/crypto/ssh
RUN go get -u github.com/pkg/sftp

RUN go install github.com/wmiller848/Fedops/fedops
RUN go install github.com/wmiller848/Fedops/fedops-warehouse
RUN go install github.com/wmiller848/Fedops/fedops-truck

VOLUME /opt/fedops
EXPOSE 13371
WORKDIR /opt/fedops
