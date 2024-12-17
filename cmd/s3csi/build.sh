#!/bin/bash
commitId=$(git rev-parse HEAD)
tag="v1.5.0-"$commitId

go build
docker build -f ./Dockerfile -t xwharbor.wxchina.com/cpaas-dev/component/csi-s3:$tag .
docker push xwharbor.wxchina.com/cpaas-dev/component/csi-s3:$tag
rm -rf s3csi
