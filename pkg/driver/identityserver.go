package driver

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/protobuf/ptypes/wrappers"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

type IdentityServer struct {
	Name    string
	Version string
}

func NewIdentityServer() *IdentityServer {
	return &IdentityServer{
		Name:    DriverName,
		Version: Version,
	}
}

// 返回插件的信息
func (i *IdentityServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {

	klog.V(4).Infof("GetPluginInfo: called with args %+v", *req)
	if i.Name == "" {
		return nil, status.Error(codes.Unavailable, "Driver name not configured,can not be null or empty")
	}

	if i.Version == "" {
		return nil, status.Error(codes.Unavailable, "Driver mis missing version")
	}

	return &csi.GetPluginInfoResponse{
		Name:          i.Name,
		VendorVersion: i.Version,
	}, nil
}

// 返回插件支持的功能信息
func (i *IdentityServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	klog.V(4).Infof("Use CONTROLLER_SERVICE GetPluginCapabilities: called with args %+v", *req)

	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
		},
	}, nil

}

// 插件健康检查
func (i *IdentityServer) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	klog.V(4).Infof("IdentityServer Probe: called with args %+v", *req)
	return &csi.ProbeResponse{Ready: &wrappers.BoolValue{Value: true}}, nil
}
