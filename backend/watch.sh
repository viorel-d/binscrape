#!/bin/bash

binaries=("$@")
argsLen=$#

runApp() {
    if (( $argsLen == 0 ));
    then
        echo "No binaries names provided. Stopping..."
        exit -1
    else
        echo "Building..."
        for binary in "${binaries[@]}";
        do
            echo $binary
            ($GOBIN/$binary &)
            echo "$binary pid: $(pidof $binary)"
            pidof $binary > "/tmp/$binary.pid"
        done
    fi
}

installGoBinaries() {
    CGO_ENABLED=0 GOOS=linux go install ./...
}

rerunApp() {
    for binary in "${binaries[@]}";
    do
        pidFile="/tmp/$binary.pid"
        if [ -f $pidFile ]
        then
            echo "Killing $binary process PID: $(cat $pidFile)"
            kill -9 $(cat $pidFile)
            rm -f $pidFile
            rm -f $GOBIN/$binary
        fi
    done

    echo "Building..."
    installGoBinaries
    runApp
}

lockBuild() {
    for binary in "${binaries[@]}";
    do
        lockFile="/tmp/$binary.lock"
        # check lock file existence
        if [ -f $lockFile ]
        then
            # waiting for the file to delete
            inotifywait -e delete $lockFile
        fi
        touch $lockFile
    done
}

unlockBuild() {
    for binary in "${binaries[@]}";
    do
        rm -f "/tmp/$binary.lock"
    done
}

# run the app for the first time
runApp

inotifywait -e modify -r -m /backend/ |
    while read path action file; do
        lockBuild
            echo "File $file"
            filePattern="[^\.]+\.go"
            if [[ $file =~ $filePattern ]];
            then
                echo "File changed $file"
                rerunApp
            fi
        unlockBuild
    done
