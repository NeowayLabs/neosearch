#!/usr/bin/env bash

./hack/deps.sh

go build -tags leveldb -v
cd neosearch && go build -tags leveldb -v 
