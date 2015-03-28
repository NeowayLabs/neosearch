build:
	docker build -t neosearch-dev-env .

all: build
	@-docker rm -vf neosearch-docker-build
	docker run --name neosearch-ctn-build -v `pwd`:/go/src/github.com/NeowayLabs/neosearch --privileged -i -t neosearch-dev-env hack/make.sh

check: build
	@-docker rm -vf neosearch-ctn-check
	docker run --name neosearch-ctn-check -v `pwd`:/go/src/github.com/NeowayLabs/neosearch --privileged -i -t neosearch-dev-env hack/check.sh

shell: build
	docker run --rm -v `pwd`:/go/src/github.com/NeowayLabs/neosearch --privileged -i -t neosearch-dev-env bash
