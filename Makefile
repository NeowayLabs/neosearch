.PHONY: all build server cli check shell docs docs-view docs-shell

DOCKER_DEVIMAGE = neosearch-dev
DOCKER_DOCSIMAGE = neosearch-docs
DEV_WORKDIR = /go/src/github.com/NeowayLabs/neosearch
CURRENT_PATH = $(shell pwd)
MOUNT_DEV_VOLUME = -v $(CURRENT_PATH):$(DEV_WORKDIR)
TEST_DIRECTORY ?= .
SHELL_EXPORT := $(foreach v,$(MAKE_ENV),$(v)='$($(v))')

ifeq ($(STORAGE_ENGINE),)
	export STORAGE_ENGINE=leveldb
else
	export STORAGE_ENGINE
endif

ifeq ($(TEST_DIRECTORY),)
	export TEST_DIRECTORY=.
else
	export TEST_DIRECTORY
endif

all: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -e STORAGE_ENGINE=$(STORAGE_ENGINE) -v `pwd`:$(DEV_WORKDIR) -i -t $(DOCKER_DEVIMAGE) hack/make.sh

server: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -e STORAGE_ENGINE=$(STORAGE_ENGINE) -v $(CURRENT_PATH):$(DEV_WORKDIR) -i -t $(DOCKER_DEVIMAGE) hack/make.sh server

cli: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -e STORAGE_ENGINE=$(STORAGE_ENGINE) -v `pwd`:$(DEV_WORKDIR) -i -t $(DOCKER_DEVIMAGE) hack/make.sh cli

library: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -e STORAGE_ENGINE=$(STORAGE_ENGINE) -v `pwd`:$(DEV_WORKDIR) -i -t $(DOCKER_DEVIMAGE) hack/make.sh library

check: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -e TEST_DIRECTORY=$(TEST_DIRECTORY) -e STORAGE_ENGINE=$(STORAGE_ENGINE) -v `pwd`:$(DEV_WORKDIR) -i  $(DOCKER_DEVIMAGE) hack/check.sh

shell: build
	docker run --rm -e STORAGE_ENGINE=$(STORAGE_ENGINE) -v `pwd`:$(DEV_WORKDIR) --privileged -i -t $(DOCKER_DEVIMAGE) bash

docs: build-docs
	docker run --rm $(MOUNT_DEV_VOLUME) --privileged $(DOCKER_DOCSIMAGE) hack/docs.sh

docs-view: docs
	xdg-open ./site/index.html

docs-shell: build-docs
	docker run --rm $(MOUNT_DEV_VOLUME) --privileged -t -i $(DOCKER_DOCSIMAGE) bash

hack/deps.txt:
	./hack/gendeps.sh "$(STORAGE_ENGINE)"

build: hack/deps.txt
	docker build -t $(DOCKER_DEVIMAGE) .

build-docs: build
	docker build -t $(DOCKER_DOCSIMAGE) -f ./docs/Dockerfile .
