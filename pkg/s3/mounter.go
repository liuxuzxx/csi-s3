package s3

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	ps "github.com/mitchellh/go-ps"
	"k8s.io/klog"
	"k8s.io/mount-utils"
)

const (
	AwsAccessKeyId     = "AWS_ACCESS_KEY_ID"
	AwsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	MountConfig        = "MOUNT_CONFIG"

	//定义Mounter挂载器类型
	RcloneMounter       = "rclone"
	MountpointS3Mounter = "mountpoint-s3"
	S3fsMounter         = "s3fs"
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
		return NewRclone(req)
	case MountpointS3Mounter:
		return NewMountpointS3(req)
	case S3fsMounter:
		return NewS3fs(req)
	default:
		return NewRclone(req)
	}
}

// 提供一个做umount操作的函数，负责处理所有的umount操作
func FuseUnmount(path string) error {
	if err := mount.New("").Unmount(path); err != nil {
		return err
	}

	p, err := findFuseMountProcess(path)
	if err != nil {
		klog.V(4).Infof("Error getting PID of fuse mount:%s", err)
		return nil
	}
	if p == nil {
		klog.V(4).Infof("Unable to find PID of fuse mount %s,it must have finished already", path)
		return nil
	}

	klog.V(4).Infof("Found fuse pid %v of mount %s, checking if it still runs", p.Pid, path)
	return waitForProcess(p, 1)
}

func waitForProcess(p *os.Process, backoff int) error {
	if backoff == 20 {
		return fmt.Errorf("timeout waiting for PID %v to end", p.Pid)
	}
	cmdLine, err := getCmdLine(p.Pid)
	if err != nil {
		klog.V(4).Infof("Error checking cmdline of PID %v, assuming it is dead: %s", p.Pid, err)
		return nil
	}
	if cmdLine == "" {
		// ignore defunct processes
		// TODO: debug why this happens in the first place
		// seems to only happen on k8s, not on local docker
		klog.V(4).Info("Fuse process seems dead, returning")
		return nil
	}
	if err := p.Signal(syscall.Signal(0)); err != nil {
		klog.V(4).Infof("Fuse process does not seem active or we are unprivileged: %s", err)
		return nil
	}
	klog.V(4).Infof("Fuse process with PID %v still active, waiting...", p.Pid)
	time.Sleep(time.Duration(backoff*100) * time.Millisecond)
	return waitForProcess(p, backoff+1)
}

func findFuseMountProcess(path string) (*os.Process, error) {
	pses, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	for _, p := range pses {
		cl, err := getCmdLine(p.Pid())
		if err != nil {
			klog.V(4).Infof("Unable to get cmdline of PID:%v :%s", p.Pid(), err)
			continue
		}
		if strings.Contains(cl, path) {
			klog.V(4).Infof("Found matching pid %v on path %s", p.Pid(), path)
			return os.FindProcess(p.Pid())
		}
	}
	return nil, nil
}

func getCmdLine(pid int) (string, error) {
	clf := fmt.Sprintf("/proc/%v/cmdline", pid)
	cl, err := os.ReadFile(clf)
	if err != nil {
		return "", err
	}
	return string(cl), nil
}
