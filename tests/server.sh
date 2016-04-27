#! /bin/bash

# Settings
REPO=github.com/pilebones/backdoorGolang/
BIN=backdoorGolang
PORT=${PORT:-"1234"}
HOSTS=${HOSTS:-"localhost 127.0.0.0 ::1 0.0.0.0"}

# Check before execute
which go > /dev/null || exit "Golang required, please install this package"
which netcat > /dev/null || exit "Netcat required, please install this package"
[ -z $GOPATH ] && exit "GOPATH must be defined, please configure your golang environment"

# Init testing env
export PATH=$PATH:$GOPATH
go get $REPO
which $BIN > /dev/null || exit "$BIN must be in your PATH"

# Run server for testing
for HOST in $HOSTS; do
	echo "- Test server on $HOST"
	# $BIN -l $HOST -p $PORT 2&> /dev/null &
	$BIN -l $HOST -p $PORT 1> /dev/null &
	PID=$!
	echo "Server running..."
	echo "*PID : $PID"
	kill -INT $PID
done
