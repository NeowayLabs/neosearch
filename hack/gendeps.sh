#!/usr/bin/env bash

STORAGE_ENGINE="$1"

if [[ -z "$STORAGE_ENGINE" ]]; then
   STORAGE_ENGINE="leveldb"
fi

echo "Generating dependencies file for storage engine ($STORAGE_ENGINE): hack/deps.txt"
go get -v -u -tags "$STORAGE_ENGINE" ./... 2>&1 | grep download | grep -v neosearch | sed 's/ (download)//g' > hack/deps.txt

