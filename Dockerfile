# This file describes the standard way to build neosearch, using docker

FROM ubuntu:14.04
MAINTAINER Tiago Katcipis <tiagokatcipis@gmail.com> (@tiagokatcipis)

# Packaged dependencies
RUN apt-get update && apt-get install -y \
        ca-certificates \
	libleveldb-dev \
	build-essential \
	curl \
	git \
	--no-install-recommends

VOLUME ["/go/src/github.com/NeowayLabs/neosearch"]

# Install Go
ENV GO_VERSION 1.4.2
RUN curl -sSL https://golang.org/dl/go${GO_VERSION}.src.tar.gz | tar -v -C /usr/local -xz \
	&& mkdir -p /go/bin
ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go
ENV CGO_ENABLED 0
RUN cd /usr/local/go/src && CGO_ENABLED=0 ./make.bash --no-clean 2>&1

# Grab Go's cover tool for dead-simple code coverage testing
RUN go get golang.org/x/tools/cmd/cover

WORKDIR /go/src/github.com/NeowayLabs/neosearch

