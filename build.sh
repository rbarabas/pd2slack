#!/usr/bin/env bash

rm -f ./pd2slack
echo "Building pd2slack binary"
go build .

echo "Building docker container"
cd docker
docker build .

echo "done"
