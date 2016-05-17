#!/bin/bash

export GOPATH=${PWD%src}

formatCode() {
	gofmt -l -w $1
}

building() {
    printf "start building $1..."
    cd $GOPATH/src/$1

    if $2;then
        printf "format $1 codes..."
        formatCode $GOPATH/src/$1
    fi

    [ -f $1 ] && rm $1
    printf "building $1 ...\n"
    go build -i -gcflags "-N -l"
}

main() {
    echo "formating lib/service/types/tool:"
    formatCode "$GOPATH/src/library"
    formatCode "$GOPATH/src/service"
    building "server" true
}

time main
