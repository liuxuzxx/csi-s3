package s3

// 定义Mounter接口，定义好对接S3的挂载的接口
type Mounter interface {
	//执行Stage的操作
	Stage(path string) error

	//执行Unstage的操作
	Unstage(path string) error

	//执行mount的挂载操作
	Mount(source string, target string) error
}
