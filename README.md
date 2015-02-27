[![Build Status](https://travis-ci.org/NeowayLabs/neosearch.svg?branch=master)](https://travis-ci.org/NeowayLabs/neosearch) [![Build Status](https://drone.io/github.com/NeowayLabs/neosearch/status.png)](https://drone.io/github.com/NeowayLabs/neosearch/latest)

NeoSearch - Neoway Full Text Search Index
==========================================

NeoSearch is a feature-limited full-text-search library with focus on indices relationships, its main goal is provide very fast JOIN operations between information stored on different indices.

It's not a complete FTS (Full Text Search) engine, in the common sense, but aims to solve very specific problems of FTS. At the moment, NeoSearch is a laboratory for research, not recommended for production usage, here we will test various technology for fast storage and search algorithms. In the future, maybe, we can proud of a very nice tech for solve search in big data companies.

NeoSearch is like a Lucene library but without all of the complexities of a complete FTS engine, written in Go, focusing on high performance search with data relationships.

It's not yet complete, still in active development, then stay tuned for updates.

# Dependencies

* leveldb >= 1.15
* snappy (opcional, only required for compressed data)
* Go > 1.3

# Install

```bash
export CGO_CFLAGS='-I <path/to/leveldb/include>'
export CGO_LDFLAGS='-L /home/secplus/projects/3rdparty/leveldb/'
go get -u -v github.com/NeowayLabs/neosearch

cd $GOPATH/src/github/NeowayLabs/neosearch
go test -tags leveldb -v .
```
