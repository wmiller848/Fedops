# Docker build for Fedops
FROM fedora:latest
MAINTAINER W. Chase Miller

USER root

# Install our deps
RUN yum -y install gcc make cmake git hostname glibc-static glibc-devel pcre-devel

# Set up golang
WORKDIR /opt
RUN git clone https://github.com/golang/go
WORKDIR /opt/go/src
RUN git checkout go1.4.2
RUN ./all.bash
ENV PATH="$PATH:/opt/go/bin"
ENV GOPATH="/go"
RUN mkdir $GOPATH
ENV PATH="$PATH:$GOPATH/bin"
RUN go get golang.org/x/tools/cmd/...

WORKDIR /opt

ADD . /go/src/github.com/wmiller848/Fedops

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
