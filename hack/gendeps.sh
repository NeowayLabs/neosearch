#!/usr/bin/env bash

STORAGE_ENGINE="$1"

if [[ -z "$STORAGE_ENGINE" ]]; then
   STORAGE_ENGINE="goleveldb"
fi

echo "Generating dependencies file for storage engine ($STORAGE_ENGINE): hack/deps.txt"
go get -insecure -v -u -tags "$STORAGE_ENGINE" ./... 2>&1 | grep download | grep -v neosearch | sed 's/ (download)//g' > hack/deps.txt

