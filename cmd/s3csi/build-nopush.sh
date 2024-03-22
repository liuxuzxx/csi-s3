#!/bin/bash
tag="v1.0.0"
go build
docker build -f ./Dockerfile -t csi-s3:$tag .
rm -rf s3csi
