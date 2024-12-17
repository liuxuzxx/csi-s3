package s3

import (
	"context"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestMinioClient(t *testing.T) {
	t.Log("连接到 MinIO 服务器")
	endpoint := "172.16.84.26:9000"
	accessKeyID := "admin"
	secretAccessKey := "passwords"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		t.Fatalf("无法创建 MinIO 客户端: %v", err)
	}

	t.Log("成功创建 MinIO 客户端")
	ctx := context.Background()
	buckets, err := minioClient.ListBuckets(ctx)
	if err != nil {
		t.Fatalf("无法列出存储桶: %v", err)
	}

	t.Log("列出存储桶:")
	for _, bucket := range buckets {
		t.Logf("存储桶名称: %s", bucket.Name)
	}
}
