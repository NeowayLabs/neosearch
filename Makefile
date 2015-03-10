build:
	docker build -t neosearch .

all: build
	docker run -v `pwd`:/go/src/github.com/NeowayLabs/neosearch --privileged -i -t neosearch hack/make.sh

check: build
	docker run -v `pwd`:/go/src/github.com/NeowayLabs/neosearch --privileged -i -t neosearch hack/check.sh

shell: build
	docker run -v `pwd`:/go/src/github.com/NeowayLabs/neosearch --privileged -i -t neosearch bash
