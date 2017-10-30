#!/bin/bash

ProjectPath=${PWD}
#export GOPATH=$ProjectPath

formatCode() {
	gofmt -l -s -w $1
}

building() {
    printf "\nstart building $1..."
    cd $ProjectPath/$1

    if $2;then
        printf "format $1 codes..."
        formatCode $ProjectPath/$1
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
    formatCode "$ProjectPath/library"
    formatCode "$ProjectPath/module"
    formatCode "$ProjectPath/service"
    formatCode "$ProjectPath/server"
    formatCode "$ProjectPath/types"
    building "server" true
}

time main
