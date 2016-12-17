#!/bin/bash

ProjectPath=${PWD%src}
export GOPATH=$ProjectPath

formatCode() {
	gofmt -l -s -w $1
}

building() {
    printf "\nstart building $1..."
    cd $ProjectPath/src/game/$1

    if $2;then
        printf "format $1 codes..."
        formatCode $ProjectPath/src/game/$1
    fi

    [ -f $1 ] && rm $1
    printf "building $1 ..."
    go build -i -gcflags "-N -l"
    if [ $? -eq 0 ];then
        printf "generate to app/..."
        mv $1 "$ProjectPath/app"
    fi
}

main() {
    echo "formating lib/service/types/tool:"
    formatCode "$ProjectPath/src/library"
    formatCode "$ProjectPath/src/service"
    building "server" true
    building "gateway" true
}

time main
