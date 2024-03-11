# 概述
支持S3协议的K8S的CSI插件实现

# 备注
由于使用的是MinIO作为存储，替换掉了NFS，但是从往上找到的一些生成支持S3协议的CSI实现，安装上去之后多多少少都是会出现一些问题，包括如下的：
https://github.com/yandex-cloud/k8s-csi-s3.git
https://github.com/ctrox/csi-s3.git

安装之后各种奇奇怪怪的问题，并且发现作者并没有去关心这些issule，所以为了快速投产，所以也就没有那么多的耐心进行等待了，直接自己开发一个CSI的实现得了

# CSI的流程
1.当我们执行PVC的创建的时候，K8S会嗲用CSI插件(使用driver的名字来区分)的Controller服务的CreateVolume接口，创建Volume(这个时候只是创建了一个Volume对象，然后记录给了K8S).  Volume的创建
2.Volume的使用:当我们使用一个Pod当中的某个容器执行volumeMounts的时候，会调用CSI插件的ControllerPublishVolume接口，将这个存储见挂载到某个主机上