package driver

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/net/context"
	"k8s.io/klog"
)

type CSIDriver struct {
	Name string
	Version string
}

type IdentityServer struct{
    CSIDriver
}

func NewIdentityServer() *IdentityServer {
	return &IdentityServer{
		Name: DriverName,
		Version: Version,
	}
}

// 返回插件的信息
func (i *IdentityServer) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {

	klog.V(4).Infof("GetPluginInfo: called with args %+v", *req)

	return &csi.GetPluginInfoResponse{
		Name:          i.Name,
		VendorVersion: i.Version,
	}, nil
}

// 返回插件支持的功能信息
func (i *IdentityServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	klog.V(4).Infof("GetPluginCapabilities: called with args %+v", *req)

	resp := &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_GROUP_CONTROLLER_SERVICE
					},
				},
			},
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS
					},
				},
			},
		},
	}

	return resp,nil

}


//插件健康检查
func (i *IdentityServer) Probe(ctx context.Context,req *csi.ProbeRequest) (*csi.ProbeResponse,error) {
	klog.V(4).Infof("Probe: called with args %+v",*req)
	return &csi.ProbeResponse{},nil
}