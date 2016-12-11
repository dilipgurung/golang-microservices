#!/usr/bin/env bash

p=`pwd`
for d in $(ls ./services); do
  echo "building services/$d"
  cd $p/services/$d
  CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo
done
cd $p
