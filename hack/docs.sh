#!/usr/bin/env bash
godoc -html=true . > docs/code.html
mkdocs serve
