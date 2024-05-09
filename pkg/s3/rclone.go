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
	rcloneCommand          = "rclone"
	defaultRcloneArguments = "--attr-timeout=5m --vfs-cache-mode=full --vfs-cache-max-age=24h --vfs-cache-max-size=10G --vfs-read-chunk-size-limit=100M --buffer-size=100M --file-perms=0777 --daemon"
)

func NewRclone(req *csi.NodePublishVolumeRequest) *Rclone {
	param := req.GetVolumeContext()
	r := &Rclone{
		bucket:    param[Bucket],
		endpoint:  param[Endpoint],
		accessKey: param[AccessKey],
		secretKey: param[SecretKey],
	}

	if v, ok := param[Arguments]; ok {
		r.arguments = v
	} else {
		r.arguments = defaultRcloneArguments
		klog.V(4).Infof("The config arguments is empty, so we use default arguments:%s", defaultRcloneArguments)
	}
	return r
}

type Rclone struct {
	bucket    string
	endpoint  string
	accessKey string
	secretKey string
	arguments string
}

func (r *Rclone) Stage(path string) error {
	klog.V(4).Info("Rclone Stage method not implements")
	return nil
}

func (r *Rclone) Unstage(path string) error {
	klog.V(4).Info("Rclone Unstage method not implements")
	return nil
}

func (r *Rclone) Mount(source string, target string) error {
	url := r.endpointUrl()

	cas := strings.Split(r.arguments, " ")

	args := []string{
		"mount",
		"minio:/" + r.bucket + "/" + source,
		target,
		"--config=/home/csi/rclone.conf",
		"--s3-endpoint=" + url,
	}
	args = append(args, cas...)
	envs := []string{
		AwsAccessKeyId + "=" + r.accessKey,
		AwsSecretAccessKey + "=" + r.secretKey,
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
