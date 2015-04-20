#!/usr/bin/env bash
mkdocs build --clean
mkdir -p site/code
godoc -html=true . > site/code/index.html
