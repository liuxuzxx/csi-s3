package driver

import (
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/liuxuzxx/csi-s3/pkg/s3"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
	"k8s.io/mount-utils"
)

type NodeServer struct {
	nodeId string
}

func NewNodeServer(nodeId string) *NodeServer {
	return &NodeServer{
		nodeId: nodeId,
	}
}

// 看到有些实现CSI的插件，NodeStageVolume方法并没有什么动作，有些甚至直接返回了Unimplement
func (n *NodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	klog.V(4).Infof("NodeStageVolume: called with args %+v", *req)
	return &csi.NodeStageVolumeResponse{}, nil
}

func (n *NodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	klog.V(4).Infof("NodeUnstageVolue: called with args %+v", req)
	return &csi.NodeUnstageVolumeResponse{}, nil
}

// 首先调用的是NodeStageVolume方法，执行了一些操作之后，才会接着调用NodePushlishVolume方法
// NodePublishVolume 从全局目录mount到目标目录(后续将映射到Pod中)
func (ns *NodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	klog.V(4).Infof("NodePublishVolume: called with args %+v", *req)
	volumeId := req.GetVolumeId()
	if len(volumeId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume Id not provided")
	}
	target := req.GetTargetPath()
	if len(target) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path not provided")
	}

	notMnt, err := checkMount(target)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !notMnt {
		return &csi.NodePublishVolumeResponse{}, nil
	}

	mnt := s3.NewMounter(req)
	err = mnt.Mount(volumeId, target)
	klog.V(4).Info("Rclone mount command execute finish!")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &csi.NodePublishVolumeResponse{}, nil
}

// 需要处理这个事件，然后在这个回调方法中调用umount的操作，避免过多的rclone和过多的mount一直不能被释放
// todo
func (ns *NodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	klog.V(4).Infof("NodeUnpublishVolume: called with args %+v", *req)

	volumeId := req.GetVolumeId()
	targetPath := req.GetTargetPath()

	if len(volumeId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "NodeUnpublishVolume:Volume ID is missing in request")
	}
	if len(targetPath) == 0 {
		return nil, status.Error(codes.InvalidArgument, "NodeUnpublishVolume:Target Path is missing in request")
	}

	if err := s3.FuseUnmount(targetPath); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.V(4).Infof("s3:volume %s has been unmounted.", volumeId)

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// NodeGetInfo 返回节点信息
func (ns *NodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	klog.V(4).Infof("NodeGetInfo: called with args %+v", *req)

	return &csi.NodeGetInfoResponse{
		NodeId: ns.nodeId,
	}, nil
}

// NodeGetCapabilities 返回节点支持的功能
func (ns *NodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	klog.V(4).Infof("NodeGetCapabilities: called with args %+v", *req)

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
					},
				},
			},
		},
	}, nil
}

func (ns *NodeServer) NodeGetVolumeStats(ctx context.Context, in *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (ns *NodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func checkMount(path string) (bool, error) {
	err := mkDirAll(path)
	if err != nil {
		return false, err
	}
	notMnt, err := mount.New("").IsLikelyNotMountPoint(path)
	if err != nil {
		return false, err
	}
	return notMnt, nil
}

func mkDirAll(path string) error {
	klog.V(4).Infof("Create Pod mount path:%s", path)
	err := os.MkdirAll(path, os.FileMode(0755))
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	return nil
}
