package s3

import (
	"fmt"
	"os"
	"os/exec"

	"k8s.io/klog"
)

const (
	mountS3Command = "mount-s3"
)

// 采用Aws开发的一个mountpoint-s3的挂载服务，因为是采用Rust编写的，所以感觉要比golang编写的rclone要性能高
// https://github.com/awslabs/mountpoint-s3  git地址

type MountpointS3 struct {
	bucket string
}

func NewMountpointS3(bucket string) *MountpointS3 {
	return &MountpointS3{
		bucket: bucket,
	}
}

func (m *MountpointS3) Stage(path string) error {
	klog.V(4).Info("MountpointS3 Stage method not implements")
	return nil
}

func (m *MountpointS3) Unstage(path string) error {
	klog.V(4).Info("MountpointS3 Unstage method not implements")
	return nil
}

func (m *MountpointS3) Mount(source string, target string) error {
	url := os.Getenv("ENDPOINT_URL")
	args := []string{
		"--endpoint-url=" + url,
		m.bucket,
		"--prefix=" + source + "/",
		target,
	}

	cmd := exec.Command(mountS3Command, args...)
	klog.V(4).Infof("Mount fuse with command:%s and args:%s", mountS3Command, args)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error execute fuse mount command:%s and args:%s", mountS3Command, args)
	}
	return nil
}
