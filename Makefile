build:
	docker build -t neosearch-dev-env .

all: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -v `pwd`:/go/src/github.com/NeowayLabs/neosearch -i -t neosearch-dev-env hack/make.sh
	@sudo chown -R $(USER):$(USER) .
	@sudo chmod -R 755 .

check: build
	@-docker rm -vf neosearch-ctn
	docker run --name neosearch-ctn -v `pwd`:/go/src/github.com/NeowayLabs/neosearch -i -t neosearch-dev-env hack/check.sh
	@sudo chown -R $(USER):$(USER) .
	@sudo chmod -R 755 .

shell: build
	docker run --rm -v `pwd`:/go/src/github.com/NeowayLabs/neosearch --privileged -i -t neosearch-dev-env bash
