#!/usr/bin/env bash

VERSION=0.1.0

DEFAULT_BUNDLES=(
    library
    server
    cli
)

bundle() {
	bundlescript=$1
	bundle=$(basename $bundlescript)
	echo "---> Making bundle: $bundle (in bundles/$VERSION/$bundle)"
	mkdir -p bundles/$VERSION/$bundle
	source "$bundlescript" "$(pwd)/bundles/$VERSION/$bundle"
}

main() {
    # We want this to fail if the bundles already exist and cannot be removed.
    # This is to avoid mixing bundles from different versions of the code.
    mkdir -p bundles
    if [ -e "bundles/$VERSION" ]; then
	echo "bundles/$VERSION already exists. Removing."
	rm -fr bundles/$VERSION && mkdir bundles/$VERSION || exit 1
	echo
    fi
    SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

    if [ $# -lt 1 ]; then
	bundles=(${DEFAULT_BUNDLES[@]})
    else
	bundles=($@)
    fi

    
    for bundle in ${bundles[@]}; do
	bundle $SCRIPTDIR/make/$bundle
	echo
    done
}

main "$@"
