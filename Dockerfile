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
	bzr \
	mercurial \
	--no-install-recommends

# Install Go
ENV GO_VERSION 1.4.2
RUN curl -sSL https://golang.org/dl/go${GO_VERSION}.src.tar.gz | tar -v -C /usr/local -xz \
	&& mkdir -p /go/bin
ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go
RUN cd /usr/local/go/src && ./make.bash --no-clean 2>&1

# Grab Go test coverage tools
RUN go get golang.org/x/tools/cmd/cover && \
    go get github.com/tools/godep && \
    go get github.com/axw/gocov/gocov && \
    go get golang.org/x/tools/cmd/cover && \
    go get github.com/golang/lint/golint && \
    go get golang.org/x/tools/cmd/goimports && \
    go get golang.org/x/tools/cmd/godoc && \
    go get golang.org/x/tools/cmd/vet

# Install package dependencies
RUN go get -d github.com/extemporalgenome/slug && \
    go get -d golang.org/x/text && \
    go get -d github.com/jmhodges/levigo && \
    go get -d github.com/iNamik/go_lexer && \
    go get -d github.com/iNamik/go_container && \
    go get -d github.com/iNamik/go_pkg && \
    go get -d github.com/jteeuwen/go-pkg-optarg && \
    go get -d github.com/peterh/liner && \
    go get -d gopkg.in/yaml.v2 && \
    go get -d launchpad.net/gommap && \
    go get -d github.com/julienschmidt/httprouter

ENV STORAGE_ENGINE leveldb

WORKDIR /go/src/github.com/NeowayLabs/neosearch
