package s3

import (
	"context"
	"fmt"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

func TestMinIOCreateDir(t *testing.T) {
	c := NewMinioClient("10.20.121.41:30629", "admin", "minioadmin", "test-bucket")
	err := c.CreateDir("pvc-adfaa")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestLinkMinioClient(t *testing.T) {
	endpoint := "localhost:9002"
	accessKey := "minioadmin"
	secretAccessKey := "minioadmin"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		fmt.Println("connect to minio:", err)
	}
	buckets, err := client.ListBuckets(context.Background())

	if err != nil {
		fmt.Println("do minio operation error:", err)
	}
	logrus.Info("buckets:", buckets)
}
