.PHONY: all build server cli check shell bundles

DOCKER_DEVIMAGE = neosearch-dev-env
DOCKER_PATH = /go/src/github.com/NeowayLabs/neosearch

export STORAGE_ENGINE="leveldb"

all: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -v `pwd`:$(DOCKER_PATH) -i -t $(DOCKER_DEVIMAGE) hack/make.sh

server: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -v `pwd`:$(DOCKER_PATH) -i -t $(DOCKER_DEVIMAGE) hack/make.sh server

cli: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -v `pwd`:$(DOCKER_PATH) -i -t $(DOCKER_DEVIMAGE) hack/make.sh cli

check: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -v `pwd`:/go/src/github.com/NeowayLabs/neosearch -i -t neosearch-dev-env hack/check.sh

shell: build
	docker run --rm -v `pwd`:/go/src/github.com/NeowayLabs/neosearch --privileged -i -t neosearch-dev-env bash

build:
	docker build -t $(DOCKER_DEVIMAGE) .
