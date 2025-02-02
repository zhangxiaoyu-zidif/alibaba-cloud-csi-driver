/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package oss

import (
	log "github.com/Sirupsen/logrus"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
)

const (
	driverName = "ossplugin.csi.alibabacloud.com"
)

var (
	version = "1.0.0"
)

type oss struct {
	driver     *csicommon.CSIDriver
	endpoint   string
	idServer   *csicommon.DefaultIdentityServer
	nodeServer *nodeServer

	cap   []*csi.VolumeCapability_AccessMode
	cscap []*csi.ControllerServiceCapability
}

//NewDriver init oss type of csi driver
func NewDriver(nodeID, endpoint string) *oss {
	log.Infof("Driver: %v version: %v", driverName, version)

	d := &oss{}
	d.endpoint = endpoint

	if nodeID == "" {
		nodeID = GetMetaData(INSTANCE_ID)
		log.Infof("Use node id : %s", nodeID)
	}
	csiDriver := csicommon.NewCSIDriver(driverName, version, nodeID)
	csiDriver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER})
	csiDriver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{csi.ControllerServiceCapability_RPC_UNKNOWN})

	d.driver = csiDriver

	return d
}

//NewDriver init oss type of csi nodeServer
func NewNodeServer(d *oss) *nodeServer {
	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(d.driver),
	}
}

func (d *oss) Run() {
	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(d.endpoint,
		csicommon.NewDefaultIdentityServer(d.driver),
		nil,
		//csicommon.NewDefaultControllerServer(d.driver),
		NewNodeServer(d))
	s.Wait()
}
