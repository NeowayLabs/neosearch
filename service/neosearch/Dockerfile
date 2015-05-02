from neowaylabs/neosearch-dev-env:latest

MAINTAINER Tiago Natel de Moura <tiago.natel@neoway.com.br> (@tiago4orion)

RUN cd /tmp && git clone https://github.com/NeowayLabs/neosearch.git && \
    cp -R /tmp/neosearch/* /go/src/github.com/NeowayLabs/neosearch/

WORKDIR /go/src/github.com/NeowayLabs/neosearch

ENV STORAGE_ENGINE leveldb
RUN hack/make.sh server

VOLUME ["/data"]

EXPOSE 9500

CMD ["./bundles/0.1.0/server/neosearch", "-d", "/data"]