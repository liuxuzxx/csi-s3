package main

import (
	"flag"

	"github.com/liuxuzxx/csi-s3/pkg/driver"
	"k8s.io/klog"
)

var (
	endpoint string
	nodeId   string
)

func main() {
	println("Start S3 CSI...")
	flag.StringVar(&endpoint, "endpoint", "", "CSI Endpoint")
	flag.StringVar(&nodeId, "nodeId", "", "node id")
	klog.InitFlags(nil)
	flag.Parse()
	d := driver.NewDriver(nodeId, endpoint)
	d.Run()
	println("End S3 CSI...")
}
