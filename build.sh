#!/bin/bash

set -e

BINARY=check_kube
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
VERSION=$(git describe --tags)

mkdir -p $SCRIPTDIR/build $SCRIPTDIR/release
rm -f $SCRIPTDIR/build/* $SCRIPTDIR/release/*

for os in linux darwin; do
	for arch in 386 amd64; do
		binfile="${BINARY}_${os}_${arch}"

		GOOS=${os} GOARCH=${arch} go build -o $SCRIPTDIR/build/$binfile

		pushd build > /dev/null
		zip -q $SCRIPTDIR/release/$binfile.zip $binfile
		popd > /dev/null
	done
done
