#!/usr/bin/env bash

go get github.com/jmhodges/levigo
go test -tags leveldb -v .
