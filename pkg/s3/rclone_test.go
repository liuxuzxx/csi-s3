package s3

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestUseRcloneMountMinio(t *testing.T) {

	argument := "--attr-timeout=5m --vfs-cache-mode=full --vfs-cache-max-age=24h --vfs-cache-max-size=10G --vfs-read-chunk-size-limit=100M --buffer-size=100M --file-perms=0777 --daemon"

	cas := strings.Split(argument, " ")

	envs := []string{
		"AWS_ACCESS_KEY_ID=admin",
		"AWS_SECRET_ACCESS_KEY=minioadmin",
	}

	args := []string{
		"mount",
		"liuxu:/k8s-dev-sc",
		"/media/liuxu/data/temp/data/minio",
		"--config=/media/liuxu/data/temp/data/rclone.conf",
		"--s3-endpoint=http://10.20.121.41:30629",
	}

	args = append(args, cas...)

	cmd := exec.Command("rclone", args...)
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Environ(), envs...)
	if out, err := cmd.Output(); err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Printf("Rclone mount成功:%s", string(out))
	}
}
