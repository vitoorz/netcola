#!/bin/sh

SCRIPT_NAME=$0
SCRIPT_DIR=${SCRIPT_NAME%/*}

PACKAGE_NAME=types
TARGET_PROTO=${PACKAGE_NAME}.proto
TARGET_GO=${PACKAGE_NAME}.pb.go

cd $SCRIPT_DIR

payloadFiles=`ls *.proto`

for payload in $payloadFiles
do
    echo "processing $payload"
    cp $payload $TARGET_PROTO
    protoc --go_out . $TARGET_PROTO
    goStyle=`echo ${payload%.*} | awk '{print tolower($0)}'`
    mv $TARGET_GO ../types/${goStyle}.go
    rm $TARGET_PROTO
done

