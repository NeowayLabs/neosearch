build:
	docker build -t docker .

shell:
	docker run -v `pwd`:/go/src/github.com/NeowayLabs/neosearch --privileged -i -t docker bash
