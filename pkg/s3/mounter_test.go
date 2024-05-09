package s3

import (
	"fmt"
	"testing"
)

func TestFuseUnmount(t *testing.T) {
	err := FuseUnmount("/media/liuxu/data/temp/data/minio")
	if err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Print("正常umount")
	}
}
