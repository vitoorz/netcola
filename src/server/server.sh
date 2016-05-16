#!/bin/sh
cd ../../
ProjectPath=`pwd`

export GOPATH=$ProjectPath
echo "*** formating codes"
cd $ProjectPath
# format codes to unify the style
gofmt -l -w src/library
gofmt -l -w src/service
gofmt -l -w src/types

cd src/server
gofmt -l -w .

echo "*** building"
time go build -i -v -gcflags "-N -l"
