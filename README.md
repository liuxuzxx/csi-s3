# 1. S3 CSI Driver (Minio/Huawei cloud OBS/Amazon S3)

## 1.1 Overview

The Mountpoint for S3 Container Storage Interface (CSI) Driver allows your Kubernetes applications to access S3 objects through a file system interface.

## 1.2 Fatures

- **Static Provisioning** - Associate an existing S3 bucket with a [PersistentVolume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) (PV) for consumption within Kubernetes.
- **Mount Options** - Mount options can be specified in the PersistentVolume (PV) resource to define how the volume should be mounted. For Mountpoint-specific options

## 1.3 Support S3 Server

|  S3 Server Type  | Supported | Remark |
| :--------------: | :-------: | :----: |
|      MinIO       |    Yes    |   No   |
| Huawei Cloud OBS |    Yes    |   No   |
|    Amazon S3     |    Yes    |   No   |

## 1.4 Support Mounter type

- **rclone** - [Rclone Github Link](https://github.com/rclone/rclone.git)
- **mountpoint-s3** [mountpoint-s3 Github Link](https://github.com/awslabs/mountpoint-s3-csi-driver.git)
- **s3fs** [s3fs Github Link](https://github.com/s3fs-fuse/s3fs-fuse.git)

## 1.5 Container Images

| Driver Version | Image(Docker hub)      |
| -------------- | ---------------------- |
| v1.4.0         | liuxuzxx/csi-s3:v1.4.0 |

<summary>Previous Images</summary>

| Driver Version | Image(Docker hub)      |
| -------------- | ---------------------- |
| v1.3.0         | liuxuzxx/csi-s3:v1.3.0 |
| v1.2.0         | liuxuzxx/csi-s3:v1.2.0 |
| v1.1.0         | liuxuzxx/csi-s3:v1.1.0 |

## 1.6 Install

We support install use Helm

1. [Install Helm](https://helm.sh/docs/intro/install/)
2. Install csi-s3

```bash
linux> git clone https://github.com/liuxuzxx/csi-s3.git
linux> cd csi-s3/deploy/s3-csi
linux> helm install csi-s3 ./ -n xxx
```

## 1.7 Self Build

```bash
linux> git clone https://github.com/liuxuzxx/csi-s3.git
linux> cd csi-s3/cmd/s3csi

#build image
linux> bash build-nopush.sh

#go build
linux> go build
```

# 2. 概述

支持 S3 协议的 K8S 的 CSI 插件实现

# 3. 备注

由于使用的是 MinIO 作为存储，替换掉了 NFS，但是从往上找到的一些生成支持 S3 协议的 CSI 实现，安装上去之后多多少少都是会出现一些问题，包括如下的：
https://github.com/yandex-cloud/k8s-csi-s3.git
https://github.com/ctrox/csi-s3.git

安装之后各种奇奇怪怪的问题，并且发现作者并没有去关心这些 issule，所以为了快速投产，所以也就没有那么多的耐心进行等待了，直接自己开发一个 CSI 的实现得了

# 4. CSI 的流程

1.当我们执行 PVC 的创建的时候，K8S 会调用 CSI 插件(使用 driver 的名字来区分)的 Controller 服务的 CreateVolume 接口，创建 Volume(这个时候只是创建了一个 Volume 对象，然后记录给了 K8S). Volume 的创建
2.Volume 的使用:当我们使用一个 Pod 当中的某个容器执行 volumeMounts 的时候，会调用 CSI 插件的 ControllerPublishVolume 接口，将这个存储见挂载到某个主机上

# 5. 打包流程

1. git clone 仓库的 git 地址
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

# 6. 其他工具使用

## 6.1 Rclone 的使用

### 6.1.1 Rclone 的基本介绍

> Rclone 是一个 rsync 的云存储版本.
> 官方说法: "rsync for cloud storage" - Google Drive, S3, Dropbox, Backblaze B2, One Drive, Swift, Hubic, Wasabi, Google Cloud Storage, Azure Blob, Azure Files, Yandex Files

> 我们主要是使用 Rclone 挂载 minio 到本地，查看 mount 的一些操作和性能

### 6.1.2 Rclone 执行挂载 minio 的操作

```bash
#在这个之前请先去Rclone官网下载对应操作系统的rclone程序，自己安装

```

## 6.2 Minio 的基本使用

### 6.2.1 提供 HTTPS 的 minio 服务

1. 首先是去 minio 的官方网站查找对应的程序包

```bash
https://min.io/open-source/download?platform=linux
```

2. 提供执行程序

```bash
#!/bin/sh

export MINIO_ROOT_USER=admin #设置默认ROOT级别的用户名
export MINIO_ROOT_PASSWORD=password  #设置默认的ROOT级别的密码

#./minio server 本地路径 --console-address ":9001" 使用9001作为控制台的web的端口
nohup ./minio server /media/liuxu/data/component/minio/data --console-address ":9001" > nohup.log 2>&1 &
```

3. 提供 https 访问服务

```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./private.key -out ./public.crt
```

4. 把生成的 private.key 和 public.crt 复制到如下目录

```bash
mv ./private.key ~/.minio/certs/
mv ./public.crt ~/.minio/certs/
```

5. 重启 minio 服务即可
