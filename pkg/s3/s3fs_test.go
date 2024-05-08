package s3

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// / s3fs k8s-dev-sc-130:/jeecgboot-vue3 -o url=http://10.20.121.130:9000 -o use_path_request_style /media/liuxu/data/temp/data/minio
func TestUseS3fsMountMinio(t *testing.T) {

	argument := "-o use_path_request_style"

	cas := strings.Split(argument, " ")

	envs := []string{
		"AWS_ACCESS_KEY_ID=admin",
		"AWS_SECRET_ACCESS_KEY=minioadmin",
	}

	args := []string{
		"k8s-dev-sc-130:/pvc-a23ff9bb-731d-4b16-94c0-68d192ca9ac2",
		"/media/liuxu/data/temp/data/minio",
		"-o",
		"url=http://10.20.121.130:9000",
	}

	args = append(args, cas...)

	cmd := exec.Command("s3fs", args...)
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Environ(), envs...)
	if out, err := cmd.Output(); err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Printf("s3fs mount成功:%s", string(out))
	}
}
