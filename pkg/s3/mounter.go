package s3

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/liuxuzxx/csi-s3/pkg/s3"
)

const (
	AwsAccessKeyId     = "AWS_ACCESS_KEY_ID"
	AwsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	MountConfig        = "MOUNT_CONFIG"

	//定义Mounter挂载器类型
	RcloneMounter       = "rclone"
	MountpointS3Mounter = "mountpoint-s3"
)

// 定义Mounter接口，定义好对接S3的挂载的接口
type Mounter interface {
	//执行Stage的操作
	Stage(path string) error

	//执行Unstage的操作
	Unstage(path string) error

	//执行mount的挂载操作
	Mount(source string, target string) error
}

func NewMounter(req *csi.NodePublishVolumeRequest) Mounter {
	param := req.GetVolumeContext()
	mounter := param[MounterType]
	switch mounter {
	case RcloneMounter:
		return s3.NewRclone(req)
	case MountpointS3Mounter:
		return s3.NewMountpointS3(req)
	default:
		return s3.NewRclone(req)
	}
}
