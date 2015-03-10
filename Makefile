build:
	docker build -t neosearch .

shell:
	docker run -v `pwd`:/go/src/github.com/NeowayLabs/neosearch --privileged -i -t neosearch bash
