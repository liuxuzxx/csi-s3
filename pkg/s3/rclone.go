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
	rcloneCommand = "rclone"
)

func NewRclone(req *csi.NodePublishVolumeRequest) *Rclone {
	param := req.GetVolumeContext()
	return &Rclone{
		bucket:    param[Bucket],
		endpoint:  param[Endpoint],
		accessKey: param[AccessKey],
		secretKey: param[SecretKey],
	}
}

type Rclone struct {
	bucket    string
	endpoint  string
	accessKey string
	secretKey string
}

func (r *Rclone) Stage(path string) error {
	klog.V(4).Info("Rclone State method not implements")
	return nil
}

func (r *Rclone) Unstage(path string) error {
	klog.V(4).Info("Rclone Unstage method not implements")
	return nil
}

func (r *Rclone) Mount(source string, target string) error {
	url := r.endpointUrl()

	args := []string{
		"mount",
		"minio:/" + r.bucket + "/" + source,
		target,
		"--config=/home/csi/rclone.conf",
		"--s3-endpoint=" + url,
		"--attr-timeout=5m",
		"--vfs-cache-mode=full",
		"--vfs-cache-max-age=24h",
		"--vfs-cache-max-size=10G",
		"--vfs-read-chunk-size-limit=100M",
		"--buffer-size=100M",
		"--file-perms 0777",
		"--daemon",
	}
	envs := []string{
		"AWS_ACCESS_KEY_ID=" + r.accessKey,
		"AWS_SECRET_ACCESS_KEY=" + r.secretKey,
	}
	cmd := exec.Command(rcloneCommand, args...)
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Environ(), envs...)
	klog.V(4).Infof("Rclone with command:%s and args:%s", rcloneCommand, args)

	if out, err := cmd.Output(); err != nil {
		return fmt.Errorf("error execute rclone mount command:%s args:%s output:%s", rcloneCommand, args, out)
	}
	return nil
}

func (r *Rclone) endpointUrl() string {
	if strings.HasPrefix(r.endpoint, "http://") || strings.HasPrefix(r.endpoint, "https://") {
		return r.endpoint
	}
	return "http://" + r.endpoint
}
