# This file describes the standard way to build neosearch, using docker

FROM ubuntu:14.04
MAINTAINER Tiago Katcipis <tiagokatcipis@gmail.com> (@tiagokatcipis)

# Packaged dependencies
RUN apt-get update && apt-get install -y \
	libleveldb-dev \
	apparmor \
	aufs-tools \
	automake \
	btrfs-tools \
	build-essential \
	curl \
	dpkg-sig \
	git \
	iptables \
	libapparmor-dev \
	libcap-dev \
	libsqlite3-dev \
	mercurial \
	parallel \
	python-mock \
	python-pip \
	python-websocket \
	reprepro \
	ruby1.9.1 \
	ruby1.9.1-dev \
	s3cmd=1.1.0* \
	--no-install-recommends

# Install Go
ENV GO_VERSION 1.4.2
RUN curl -sSL https://golang.org/dl/go${GO_VERSION}.src.tar.gz | tar -v -C /usr/local -xz \
	&& mkdir -p /go/bin
ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go
RUN cd /usr/local/go/src && ./make.bash --no-clean 2>&1

# Grab Go's cover tool for dead-simple code coverage testing
RUN go get golang.org/x/tools/cmd/cover

# Add an unprivileged user to be used for tests which need it
RUN groupadd -r neosearch
RUN useradd --create-home --gid neosearch unprivilegeduser

VOLUME /var/lib/docker
WORKDIR /go/src/github.com/NeowayLabs/neosearch

# Let us use a .bashrc file
RUN ln -sfv $PWD/.bashrc ~/.bashrc

# Upload docker source
COPY . /go/src/github.com/NeowayLabs/neosearch
