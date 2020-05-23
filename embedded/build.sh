#!/usr/bin/env bash

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)" # get dir containing this script
cd $DIR                                                             # always from from script dir

go build -v -o libfactomd.so -buildmode=c-shared libfactomd.go

BUILD=$HOME/Workspace/factomlib-py/

# KLUDGE: push to working copy of python lib
if [[ -d $BUILD ]]; then
	cp libfactomd.so ${BUILD}/factomlib/libfactomd.so
else
	echo 'skip'
fi
