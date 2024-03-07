FROM xwharbor.wxchina.com/cpaas-pub/golang:1.22.1
LABEL maintainers="lx"
LABEL description="MinIO S3 CSI Driver"

COPY ./s3csi /home/csi/

WORKDIR /home/csi/

ENTRYPOINT ["/home/csi/s3csi"]