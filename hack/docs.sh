#!/usr/bin/env bash

echo "Generating mkdocs"
mkdocs build --clean

echo "Generating Go code docs"
mkdir -p site/code
godoc -html=true . > site/code/index.html

echo "Generating REST API docs"
swagger="/swagger-codegen/modules/swagger-codegen-cli/target/swagger-codegen-cli.jar"
java -jar $swagger generate -i ./docs/rest/api.json -l html  -o site/rest
