package s3

const (
	AwsAccessKeyId     = "AWS_ACCESS_KEY_ID"
	AwsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	MountConfig        = "MOUNT_CONFIG"
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
