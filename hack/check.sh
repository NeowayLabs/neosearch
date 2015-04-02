#!/usr/bin/env bash

godep go test -tags leveldb ./... -bench=store
