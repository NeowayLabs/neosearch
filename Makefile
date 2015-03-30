.PHONY: all build check shell

DOCKER_DEVIMAGE = neosearch-dev-env
DOCKER_PATH = /go/src/github.com/NeowayLabs/neosearch

build:
	docker build -t neosearch-dev-env .

all: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -v `pwd`:$(DOCKER_PATH) -i -t $(DOCKER_DEVIMAGE) hack/make.sh server cli

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
