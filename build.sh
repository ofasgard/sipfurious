#!/bin/bash

export GOPATH=`pwd`
export GOBIN=`pwd`/bin

go get -d sipfurious
go install sipfurious

