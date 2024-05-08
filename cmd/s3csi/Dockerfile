FROM xwharbor.wxchina.com/cpaas-dev/ubuntu:23.10
LABEL maintainers="lx"
LABEL description="MinIO S3 CSI Driver"

COPY ./s3csi /home/csi/
COPY ./mount-s3.deb /tmp/
RUN apt-get update && apt-get install -y /tmp/mount-s3.deb && apt-get install -y fuse3
COPY ./rclone.conf /home/csi/
COPY ./rclone /usr/bin/
RUN chmod 755 /usr/bin/rclone
RUN apt install -y s3fs

WORKDIR /home/csi/

ENTRYPOINT ["/home/csi/s3csi"]