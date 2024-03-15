FROM xwharbor.wxchina.com/cpaas-dev/ubuntu:23.10
LABEL maintainers="lx"
LABEL description="MinIO S3 CSI Driver"

COPY ./s3csi /home/csi/

WORKDIR /home/csi/

ENTRYPOINT ["/home/csi/s3csi"]