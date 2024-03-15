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

// 格式化磁盘 Mount到全局目录
func (n *NodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	klog.V(4).Infof("NodeStageVolume: called with args %+v", *req)

	volumeId := req.GetVolumeId()
	path := req.GetStagingTargetPath()
	if len(volumeId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume Id missing in request")
	}

	if len(path) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}
	notMnt, err := checkMount(path)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	//如果是已经挂载了，那么直接返回响应即可
	if !notMnt {
		return &csi.NodeStageVolumeResponse{}, nil
	}

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
	mnt := s3.NewMountpointS3("k8s-dev-sc")
	err = mnt.Mount(volumeId, target)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *NodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	klog.V(4).Infof("NodeUnpublishVolume: called with args %+v", *req)

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

// TODO 看到Pod里面其实没有创建/var/lib/kubelet/pods/d51de966-1a5f-4d35-843a-56b5e4cf6ed2/volumes/kubernetes.io~csi/pvc-967ee658-01ec-445c-ac7f-6fb058e22c7b/mount 后面这个路径，问题出现在这个地方
func checkMount(path string) (bool, error) {
	notMnt, err := mount.New("").IsLikelyNotMountPoint(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(path, 0750); err != nil {
				return false, err
			}
			notMnt = true
		} else {
			return false, err
		}
	}
	return notMnt, nil
}
