# S3 CSI Driver (Minio/Huawei cloud OBS/Amazon S3)

## Overview
The Mountpoint for S3 Container Storage Interface (CSI) Driver allows your Kubernetes applications to access S3 objects through a file system interface.

## Fatures
* **Static Provisioning** - Associate an existing S3 bucket with a [PersistentVolume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) (PV) for consumption within Kubernetes.
* **Mount Options** - Mount options can be specified in the PersistentVolume (PV) resource to define how the volume should be mounted. For Mountpoint-specific options

## Support S3 Server

|S3 Server Type|Supported|Remark|
|:--:|:--:|:--:|
|MinIO|Yes|No|
|Huawei Cloud OBS|Yes|No|
|Amazon S3|Yes|No|

## Support Mounter type

* **rclone** -  [Rclone Github Link](https://github.com/rclone/rclone.git)
* **mountpoint-s3** [mountpoint-s3 Github Link](https://github.com/awslabs/mountpoint-s3-csi-driver.git)

## Container Images
| Driver Version | Image(Docker hub)|
|----------------|------------------|
| v1.2.0         | liuxuzxx/csi-s3:v1.2.0|

<summary>Previous Images</summary>

| Driver Version | Image(Docker hub) |
|----------------|-------------------|
| v1.1.0         | liuxuzxx/csi/s3:v1.1.0|

## Install

We support install use Helm

1. [Install Helm](https://helm.sh/docs/intro/install/)
2. Install csi-s3
```bash
linux> git clone https://github.com/liuxuzxx/csi-s3.git
linux> cd csi-s3/deploy/s3-csi
linux> helm install csi-s3 ./ -n xxx
```


## Self Build

```bash
linux> git clone https://github.com/liuxuzxx/csi-s3.git
linux> cd csi-s3/cmd/s3csi

#build image
linux> bash build-nopush.sh

#go build
linux> go build
```


# 概述
支持S3协议的K8S的CSI插件实现

# 备注
由于使用的是MinIO作为存储，替换掉了NFS，但是从往上找到的一些生成支持S3协议的CSI实现，安装上去之后多多少少都是会出现一些问题，包括如下的：
https://github.com/yandex-cloud/k8s-csi-s3.git
https://github.com/ctrox/csi-s3.git

安装之后各种奇奇怪怪的问题，并且发现作者并没有去关心这些issule，所以为了快速投产，所以也就没有那么多的耐心进行等待了，直接自己开发一个CSI的实现得了

# CSI的流程
1.当我们执行PVC的创建的时候，K8S会调用CSI插件(使用driver的名字来区分)的Controller服务的CreateVolume接口，创建Volume(这个时候只是创建了一个Volume对象，然后记录给了K8S).  Volume的创建
2.Volume的使用:当我们使用一个Pod当中的某个容器执行volumeMounts的时候，会调用CSI插件的ControllerPublishVolume接口，将这个存储见挂载到某个主机上

# 打包流程
1. git clone 仓库的git地址
2. 打包并且构建镜像并推送
```bash
cd ./cmd/s3csi
bash build-nopush.sh
```
3.执行安装
```bash
cd ./deploy/s3-csi
helm install csi-s3 ./ -n namespace(自定义)
```