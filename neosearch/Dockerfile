from neowaylabs/neosearch-dev-env:latest

MAINTAINER Tiago Natel de Moura <tiago.natel@neoway.com.br> (@tiago4orion)

WORKDIR /go/src/github.com/NeowayLabs/neosearch

RUN hack/make.sh server

VOLUME ["/data"]

EXPOSE 9500

CMD ["./bundles/0.1.0/server/neosearch", "-d", "/data"]