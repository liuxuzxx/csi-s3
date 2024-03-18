package driver

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/klog"
)

type Driver struct {
	nodeId   string
	endpoint string
}

const (
	Version       = "v1.0.0"
	DriverName    = "minio.s3.csi.xw.com"
	DevicePathKey = "devicePath"
)

func NewDriver(nodeId, endpoint string) *Driver {
	klog.V(4).Infof("Driver: %v version: %v", DriverName, Version)
	return &Driver{
		nodeId:   nodeId,
		endpoint: endpoint,
	}
}

func (d *Driver) Run() {
	ctl := NewControllerServer()
	identity := NewIdentityServer()
	node := NewNodeServer(d.nodeId)

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(logGRPC),
	}
	srv := grpc.NewServer(opts...)

	csi.RegisterControllerServer(srv, ctl)
	csi.RegisterIdentityServer(srv, identity)
	csi.RegisterNodeServer(srv, node)

	proto, addr, err := ParseEndpoint(d.endpoint)
	klog.V(4).Infof("protocol: %s,addr: %s", proto, addr)
	if err != nil {
		klog.Fatal(err.Error())
	}

	if proto == "unix" {
		addr = "/" + addr
		if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
			klog.Fatalf("Failed to remove %s, error: %s", addr, err.Error())
		}
	}

	listener, err := net.Listen(proto, addr)
	if err != nil {
		klog.Fatalf("Failed to listen: %v", err)
	}

	srv.Serve(listener)

}

func logGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	klog.V(4).Infof("GRPC call: %s", info.FullMethod)
	klog.V(4).Infof("GRPC request: %s", protosanitizer.StripSecrets(req))
	resp, err := handler(ctx, req)
	if err != nil {
		klog.Errorf("GRPC error: %v", err)
	} else {
		klog.V(4).Infof("GRPC response: %s", protosanitizer.StripSecrets(resp))
	}
	return resp, err
}

func ParseEndpoint(ep string) (string, string, error) {
	if strings.HasPrefix(strings.ToLower(ep), "unix://") || strings.HasPrefix(strings.ToLower(ep), "tcp://") {
		s := strings.SplitN(ep, "://", 2)
		if s[1] != "" {
			return s[0], s[1], nil
		}
	}
	return "", "", fmt.Errorf("invalid endpoint: %v", ep)
}
