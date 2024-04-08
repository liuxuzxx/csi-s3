package s3

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog"
)

const (
	mountS3Command = "mount-s3"
)

// 采用Aws开发的一个mountpoint-s3的挂载服务，因为是采用Rust编写的，所以感觉要比golang编写的rclone要性能高
// git地址:https://github.com/awslabs/mountpoint-s3
//但是有个问题，就是如下官方网站说的:
//but probably not the right fit for applications that:

//use file operations that S3 doesn't natively support, like directory renaming or symlinks
//make edits to existing files (don't work on your Git repository or run vim in Mountpoint 😄)

type MountpointS3 struct {
	bucket    string
	endpoint  string
	accessKey string
	secretKey string
}

func NewMountpointS3(req *csi.NodePublishVolumeRequest) *MountpointS3 {
	param := req.GetVolumeContext()
	return &MountpointS3{
		bucket:    param[Bucket],
		endpoint:  param[Endpoint],
		accessKey: param[AccessKey],
		secretKey: param[SecretKey],
	}
}

func (m *MountpointS3) endpointUrl() string {
	if strings.HasPrefix(m.endpoint, "http://") || strings.HasPrefix(m.endpoint, "https://") {
		return m.endpoint
	}
	return "http://" + m.endpoint
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
	url := m.endpointUrl()
	args := []string{
		"--endpoint-url=" + url,
		"--allow-delete",
		m.bucket,
		"--prefix=" + source + "/",
		target,
	}

	cmd := exec.Command(mountS3Command, args...)
	envs := []string{
		"AWS_ACCESS_KEY_ID=" + m.accessKey,
		"AWS_SECRET_ACCESS_KEY=" + m.secretKey,
	}
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Environ(), envs...)
	klog.V(4).Infof("Mount fuse with command:%s and args:%s", mountS3Command, args)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error execute fuse mount command:%s and args:%s", mountS3Command, args)
	}
	return nil
}
