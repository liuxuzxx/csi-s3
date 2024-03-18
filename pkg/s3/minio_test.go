package s3

import (
	"fmt"
	"testing"
)

func TestMinIOCreateDir(t *testing.T) {
	c := NewMinioClient("10.20.121.41:30629", "admin", "minioadmin", "test-bucket")
	err := c.CreateDir("pvc-adfaa")
	if err != nil {
		fmt.Println(err.Error())
	}
}
