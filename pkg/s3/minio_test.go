package s3

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestMinioClient(t *testing.T) {
	endpoint := "cpaas-minio.minio:9000"
	accessKeyID := "admin"
	secretAccessKey := "minioadmin"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	buckets, err := minioClient.ListBuckets(ctx)
	if err != nil {
		fmt.Println("ListBuckets Error!")
		log.Fatalln(err)
	}
	fmt.Println("Print Buckets!")
	for _, bucket := range buckets {
		fmt.Println(bucket.Name)
	}

}
