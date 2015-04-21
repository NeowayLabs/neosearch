.PHONY: all build server cli check shell bundles

DOCKER_DEVIMAGE = neosearch-dev-env
DOCKER_PATH = /go/src/github.com/NeowayLabs/neosearch
CURRENT_PATH = $(shell pwd)
SHELL_EXPORT := $(foreach v,$(MAKE_ENV),$(v)='$($(v))')

ifeq ($(STORAGE_ENGINE),)
	export STORAGE_ENGINE=leveldb
else
	export STORAGE_ENGINE
endif

all: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -e STORAGE_ENGINE=$(STORAGE_ENGINE) -v `pwd`:$(DOCKER_PATH) -i -t $(DOCKER_DEVIMAGE) hack/make.sh

server: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -e STORAGE_ENGINE=$(STORAGE_ENGINE) -v $(CURRENT_PATH):$(DOCKER_PATH) -i -t $(DOCKER_DEVIMAGE) hack/make.sh server

cli: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -e STORAGE_ENGINE=$(STORAGE_ENGINE) -v `pwd`:$(DOCKER_PATH) -i -t $(DOCKER_DEVIMAGE) hack/make.sh cli

check: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -e STORAGE_ENGINE=$(STORAGE_ENGINE) -v `pwd`:/go/src/github.com/NeowayLabs/neosearch -i -t $(DOCKER_DEVIMAGE) hack/check.sh

shell: build
	docker run --rm -e STORAGE_ENGINE=$(STORAGE_ENGINE) -v `pwd`:/go/src/github.com/NeowayLabs/neosearch --privileged -i -t $(DOCKER_DEVIMAGE) bash

docs: build
	docker run --rm -v `pwd`:/go/src/github.com/NeowayLabs/neosearch -p 8000:8000 $(DOCKER_DEVIMAGE) hack/docs.sh
	xdg-open ./site/index.html

build:
	docker build -t $(DOCKER_DEVIMAGE) .
