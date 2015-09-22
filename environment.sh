#!/bin/sh


export GOPATH=$HOME/.go

go get gopkg.in/mgo.v2

go get -u -v github.com/btcsuite/btcd/...
