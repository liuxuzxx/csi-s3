package s3

import (
	"bytes"
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"k8s.io/klog"
)

const (
	Bucket    string = "bucket"
	AccessKey string = "access-key"
	SecretKey string = "secret-key"
	Endpoint  string = "endpoint"
)

// 增加Minio的操作封装对象，方便处理一些操作
type MinioClient struct {
	endpoint        string
	accessKeyId     string
	secretAccessKey string
	useSSL          bool
	client          *minio.Client
	bucketName      string
	ctx             context.Context
}

func NewClient(req *csi.CreateVolumeRequest) *MinioClient {
	p := req.GetParameters()

	ak := p[AccessKey]
	sk := p[SecretKey]
	ep := p[Endpoint]
	bk := p[Bucket]
	klog.V(4).Infof("Fetch s3 config:%s %s %s %s", ak, sk, ep, bk)
	return NewMinioClient(ep, ak, sk, bk)
}

func NewMinioClient(endpoint string, accessKey string, secretAccessKey string, bucket string) *MinioClient {

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		klog.V(4).Infof(err.Error())
	}
	return &MinioClient{
		endpoint:        endpoint,
		accessKeyId:     accessKey,
		secretAccessKey: secretAccessKey,
		useSSL:          false,
		client:          client,
		bucketName:      bucket,
		ctx:             context.Background(),
	}
}

// 处理创建文件夹相关的动作
func (m *MinioClient) CreateDir(path string) error {
	_, err := m.client.PutObject(m.ctx, m.bucketName, path+"/", bytes.NewReader([]byte("")), 0, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
