#!/bin/bash
tag="v1.4.0"

go build

docker build -f ./Dockerfile -t xwharbor.wxchina.com/cpaas-dev/component/csi-s3:$tag .

docker push xwharbor.wxchina.com/cpaas-dev/component/csi-s3:$tag

rm -rf s3csi
