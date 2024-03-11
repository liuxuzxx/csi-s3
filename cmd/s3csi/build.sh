#!/bin/bash

go build

docker build -f ./Dockerfile -t xwharbor.wxchina.com/cpaas-dev/component/csi-s3:release-v1.0.0 .

docker push xwharbor.wxchina.com/cpaas-dev/component/csi-s3:release-v1.0.0

rm -rf s3csi