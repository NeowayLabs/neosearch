# neosearch-cli

Command-line tool for executing low-level commands at the core of NeoSearch.

# Installation

```
go get -tags leveldb github.com/NeowayLabs/neosearch
export CGO_LDFLAGS="-L/usr/lib -lleveldb -lsnappy -lstdc++"
export GO_LDFLAGS="-extld g++ -linkmode external -extldflags -static"
go get -tags leveldb -x -a -ldflags "$GO_LDFLAGS" -v github.com/NeowayLabs/neosearch
$GOPATH/bin/neosearch-cli -d /data
```

# NeoSearch Key/Value Syntax

NeoSearch has its own key/value syntax for commands. It's ridiculous simple:

```
USING <database-name> <command> <key> <value/optional>
```
Examples:
```
USING titie.idx SET neosearch "fast searching with document/indexes joins, spatial index and more"
```
