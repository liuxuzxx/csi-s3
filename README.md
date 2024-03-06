# 概述
支持S3协议的K8S的CSI插件实现

# 备注
由于使用的是MinIO作为存储，替换掉了NFS，但是从往上找到的一些生成支持S3协议的CSI实现，安装上去之后多多少少都是会出现一些问题，包括如下的：
https://github.com/yandex-cloud/k8s-csi-s3.git
https://github.com/ctrox/csi-s3.git

安装之后各种奇奇怪怪的问题，并且发现作者并没有去关心这些issule，所以为了快速投产，所以也就没有那么多的耐心进行等待了，直接自己开发一个CSI的实现得了
