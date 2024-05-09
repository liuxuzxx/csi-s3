package s3

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog"
)

/// 实现s3fs的挂载器，之前使用rclone的时候，会出现一些问题，目前无法解决，例如我执行了mount之后，但是有时候Pod发现
/// 找不到某个文件，但是通过cat又是可以看到的
/// 还有就是我执行umount之后，使用top在Pod种查看，还是有这个进程的
/// 下面是一个可以在本地运行的参考脚本
/// #!/bin/sh
/// export AWS_ACCESS_KEY_ID=admin
/// export AWS_SECRET_ACCESS_KEY=minioadmin
/// s3fs -o url=http://10.20.121.130:9000 -o use_path_request_style k8s-dev-sc-130 /media/liuxu/data/temp/data/minio

const (
	s3fsCommand          = "s3fs"
	defaultS3fsArguments = "-o use_path_request_style"
)

func NewS3fs(req *csi.NodePublishVolumeRequest) *S3fs {
	param := req.GetVolumeContext()
	s := &S3fs{
		bucket:    param[Bucket],
		endpoint:  param[Endpoint],
		accessKey: param[AccessKey],
		secretKey: param[SecretKey],
	}

	if v, ok := param[Arguments]; ok {
		s.arguments = v
	} else {
		s.arguments = defaultS3fsArguments
		klog.V(4).Infof("The s3fs config arguments is empth,so we use default arguments:%s", defaultS3fsArguments)
	}

	return s
}

type S3fs struct {
	bucket    string
	endpoint  string
	accessKey string
	secretKey string
	arguments string
}

func (s *S3fs) Stage(path string) error {
	klog.V(4).Info("S3fs Stage method not implements")
	return nil
}

func (s *S3fs) Unstage(path string) error {
	klog.V(4).Info("S3fs Unstage method not implements")
	return nil
}

func (s *S3fs) Mount(source string, target string) error {
	url := s.endpointUrl()
	cas := strings.Split(s.arguments, " ")
	args := []string{
		strings.Join([]string{s.bucket, source}, ":/"),
		target,
		"-o",
		"url=" + url,
	}

	args = append(args, cas...)
	envs := []string{
		AwsAccessKeyId + "=" + s.accessKey,
		AwsSecretAccessKey + "=" + s.secretKey,
	}

	cmd := exec.Command(s3fsCommand, args...)
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Environ(), envs...)
	klog.V(4).Infof("S3fs with command:%s, and args:%s", s3fsCommand, args)

	if out, err := cmd.Output(); err != nil {
		return fmt.Errorf("error execute s3fs mount command:%s args:%s outputL%s", s3fsCommand, args, out)
	}
	return nil
}

func (s *S3fs) endpointUrl() string {
	if strings.HasPrefix(s.endpoint, "http://") || strings.HasPrefix(s.endpoint, "https://") {
		return s.endpoint
	}
	return "http://" + s.endpoint
}
