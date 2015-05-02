#!/bin/bash
# The script does automatic checking on a Go package and its sub-packages, including:
# 1. gofmt         (http://golang.org/cmd/gofmt/)
# 2. goimports     (https://github.com/bradfitz/goimports)
# 3. golint        (https://github.com/golang/lint)
# 4. go vet        (http://golang.org/cmd/vet)
# 5. race detector (http://blog.golang.org/race-detector)
# 6. test coverage (http://blog.golang.org/cover)

set -e

GO="go"
TEST_FLAGS="-tags leveldb"

# Automatic checks
test -z "$(gofmt -l -w .     | tee /dev/stderr)"
#test -z "$(goimports -l -w . | tee /dev/stderr)"
test -z "$(golint .          | tee /dev/stderr)"
#$GO vet ./...
#$GO test -tags leveldb -race ./...

# Run test coverage on each subdirectories and merge the coverage profile.

echo "mode: count" > profile.cov

$GO get -tags "$STORAGE_ENGINE" ./...

# Standard $GO tooling behavior is to ignore dirs with leading underscors
for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path './Godeps/*' -type d);
do
    if ls $dir/*.go &> /dev/null; then
	$GO test $TEST_FLAGS -covermode=count -coverprofile="$dir/profile.tmp" "$dir"
	if [ -f $dir/profile.tmp ]
	then
            cat $dir/profile.tmp | tail -n +2 >> profile.cov
            rm $dir/profile.tmp
	fi
    fi
done

$GO tool cover -func profile.cov

# To submit the test coverage result to coveralls.io,
# use goveralls (https://github.com/mattn/goveralls)
# goveralls -coverprofile=profile.cov -service=travis-ci
