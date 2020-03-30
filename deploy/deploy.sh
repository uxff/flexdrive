#!/bin/bash
go build cmd/node/main.go

docker build
# todo get docker image id
#docker run $dockerImageId 


