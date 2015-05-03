#!/usr/bin/env bash

DEPS_FILE="$1"

if [[ -z "$DEPS_FILE" ]]; then
    DEPS_FILE="/deps.txt"
fi

cat "$DEPS_FILE" | xargs go get -tags $STORAGE_ENGINE -v -d 2>/dev/null

exit 0
