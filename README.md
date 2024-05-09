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
* **s3fs** [s3fs Github Link](https://github.com/s3fs-fuse/s3fs-fuse.git)

## Container Images
| Driver Version | Image(Docker hub)|
|----------------|------------------|
| v1.4.0         | liuxuzxx/csi-s3:v1.4.0|

<summary>Previous Images</summary>

| Driver Version | Image(Docker hub) |
|----------------|-------------------|
| v1.3.0         | liuxuzxx/csi-s3:v1.3.0|
| v1.2.0         | liuxuzxx/csi-s3:v1.2.0|
| v1.1.0         | liuxuzxx/csi-s3:v1.1.0|

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

# 版本功能计划

## v1.4.0
1. 修复监听umount事件,当运行一段时间后，发现如下的场景
```bash
PID USER      PR  NI    VIRT    RES    SHR S  %CPU  %MEM     TIME+ COMMAND                                                                                                                                     
998 root      20   0 1290216  76812  39040 S   7.6   0.1   0:02.47 rclone                                                                                                                                      
983 root      20   0 1289704  63384  38656 S   1.3   0.1   0:00.70 rclone                                                                                                                                      
1 root      20   0 2200952  25660  12032 S   0.0   0.0   0:02.51 s3csi                                                                                                                                       
38 root      20   0       0      0      0 Z   0.0   0.0   0:17.22 rclone                                                                                                                                      
79 root      20   0 1290728  79052  41216 S   0.0   0.1   0:05.75 rclone                                                                                                                                      
116 root      20   0       0      0      0 Z   0.0   0.0   0:14.53 rclone                                                                                                                                      
173 root      20   0       0      0      0 Z   0.0   0.0   0:28.67 rclone                                                                                                                                      
220 root      20   0       0      0      0 Z   0.0   0.0   0:11.42 rclone                                                                                                                                      
261 root      20   0       0      0      0 Z   0.0   0.0   0:10.74 rclone                                                                                                                                      
315 root      20   0       0      0      0 Z   0.0   0.0   0:38.89 rclone                                                                                                                                      
357 root      20   0 1291112  83376  40960 S   0.0   0.1   0:07.49 rclone                                                                                                                                      
407 root      20   0 1291176  82480  40960 S   0.0   0.1   0:05.95 rclone                                                                                                                                      
504 root      20   0       0      0      0 Z   0.0   0.0   0:18.58 rclone                                                                                                                                      
551 root      20   0       0      0      0 Z   0.0   0.0   0:07.68 rclone                                                                                                                                      
594 root      20   0 1291432  82124  40832 S   0.0   0.1   0:05.39 rclone                                                                                                                                      
636 root      20   0       0      0      0 Z   0.0   0.0   0:12.27 rclone                                                                                                                                      
681 root      20   0       0      0      0 Z   0.0   0.0   0:09.34 rclone                                                                                                                                      
729 root      20   0       0      0      0 Z   0.0   0.0   0:10.71 rclone                                                                                                                                      
779 root      20   0       0      0      0 Z   0.0   0.0   0:08.45 rclone                                                                                                                                      
825 root      20   0       0      0      0 Z   0.0   0.0   0:05.69 rclone                                                                                                                                      
868 root      20   0       0      0      0 Z   0.0   0.0   0:05.94 rclone                                                                                                                                      
958 root      20   0 1290856  74544  40320 S   0.0   0.1   0:05.21 rclone                                                                                                                                      
1025 root      20   0    2776   1536   1536 S   0.0   0.0   0:00.00 sh                                                                                                                                          
1031 root      20   0    8788   4864   2816 R   0.0   0.0   0:00.01 top 
```
看到的情况是：很多的rclone进程，问题产生的原因是：没有监听处理umount事件，导致很多的rclone的daemon程序，占据很多内存，出现错误