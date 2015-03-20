build:
	docker build -t neosearch .

all: build
	@-docker rm -vf neosearch-docker-build
	docker run --name neosearch-ctn-build -v `pwd`:/go/src/github.com/NeowayLabs/neosearch --privileged -i -t neosearch hack/make.sh

check: build
	@-docker rm -vf neosearch-ctn-check
	docker run --name neosearch-ctn-check -v `pwd`:/go/src/github.com/NeowayLabs/neosearch --privileged -i -t neosearch hack/check.sh

shell: build
	docker run --rm -v `pwd`:/go/src/github.com/NeowayLabs/neosearch --privileged -i -t neosearch bash
