#!/bin/bash

export GOPATH=`pwd`
export GOBIN=`pwd`/bin

go get -d gossiper
go install gossiper
