// Copyright CERN.
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package node

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"gitlab.cern.ch/kubernetes/storage/eosxd-csi/internal/mountutils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements csi.NodeServer interface.
type Server struct {
	nodeID string
	caps   []*csi.NodeServiceCapability
}

const (
	eosRoot = "/eos"
)

var (
	_ csi.NodeServer = (*Server)(nil)
)

func New(nodeID string) *Server {
	enabledCaps := []csi.NodeServiceCapability_RPC_Type{
		csi.NodeServiceCapability_RPC_VOLUME_CONDITION,
	}

	var caps []*csi.NodeServiceCapability
	for _, c := range enabledCaps {
		caps = append(caps, &csi.NodeServiceCapability{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: c,
				},
			},
		})
	}

	return &Server{
		nodeID: nodeID,
		caps:   caps,
	}
}

func (srv *Server) NodeGetCapabilities(
	ctx context.Context,
	req *csi.NodeGetCapabilitiesRequest,
) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: srv.caps,
	}, nil
}

func (srv *Server) NodeGetInfo(
	ctx context.Context,
	req *csi.NodeGetInfoRequest,
) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{
		NodeId: srv.nodeID,
	}, nil
}

func (srv *Server) NodePublishVolume(
	ctx context.Context,
	req *csi.NodePublishVolumeRequest,
) (*csi.NodePublishVolumeResponse, error) {
	if err := validateNodePublishVolumeRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	targetPath := req.GetTargetPath()

	if err := os.MkdirAll(targetPath, 0700); err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to create mountpoint directory at %s: %v", targetPath, err)
	}

	mntState, err := mountutils.GetState(targetPath)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to probe mountpoint %s: %v", targetPath, err)
	}

	switch mntState {
	case mountutils.StNotMounted:
		if err := doMount(req); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to bind mount: %v", err)
		}
		fallthrough
	case mountutils.StMounted:
		return &csi.NodePublishVolumeResponse{}, nil
	default:
		return nil, status.Errorf(codes.Internal,
			"unexpected mountpoint state in %s: expected %s or %s, got %s",
			targetPath, mountutils.StNotMounted, mountutils.StMounted, mntState)
	}
}

func doMount(req *csi.NodePublishVolumeRequest) error {
	targetPath := req.GetTargetPath()

	if repository := req.GetVolumeContext()["repository"]; repository != "" {
		// Mount a single repository.
		return bindMount(path.Join(eosRoot, repository), targetPath)
	}

	// Mount the whole autofs-EOS root.
	return slaveRecursiveBind(eosRoot, targetPath)
}

func (srv *Server) NodeUnpublishVolume(
	ctx context.Context,
	req *csi.NodeUnpublishVolumeRequest,
) (*csi.NodeUnpublishVolumeResponse, error) {
	if err := validateNodeUnpublishVolumeRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	targetPath := req.GetTargetPath()

	mntState, err := mountutils.GetState(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &csi.NodeUnpublishVolumeResponse{}, nil
		}

		return nil, status.Errorf(codes.Internal,
			"failed to probe for mountpoint %s: %v", targetPath, err)
	}

	if mntState != mountutils.StNotMounted {
		if err := recursiveUnmount(targetPath); err != nil {
			return nil, status.Errorf(codes.Internal,
				"failed to unmount %s: %v", targetPath, err)
		}
	}

	err = os.Remove(targetPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (srv *Server) NodeStageVolume(
	ctx context.Context,
	req *csi.NodeStageVolumeRequest,
) (*csi.NodeStageVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (srv *Server) NodeUnstageVolume(
	ctx context.Context,
	req *csi.NodeUnstageVolumeRequest,
) (*csi.NodeUnstageVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (srv *Server) NodeGetVolumeStats(
	ctx context.Context,
	req *csi.NodeGetVolumeStatsRequest,
) (*csi.NodeGetVolumeStatsResponse, error) {
	if err := validateNodeGetVolumeStatsRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	volPath := req.GetVolumePath()

	st, err := mountutils.GetState(volPath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get mount state for %s: %v", volPath, err)
	}

	switch st {
	case mountutils.StMounted:
		return &csi.NodeGetVolumeStatsResponse{}, nil
	default:
		return &csi.NodeGetVolumeStatsResponse{
			VolumeCondition: &csi.VolumeCondition{
				Abnormal: true,
				Message:  fmt.Sprintf("mountpoint %s is %s", volPath, st),
			},
		}, nil
	}
}

func (srv *Server) NodeExpandVolume(
	ctx context.Context,
	req *csi.NodeExpandVolumeRequest,
) (*csi.NodeExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func validateNodeGetVolumeStatsRequest(req *csi.NodeGetVolumeStatsRequest) error {
	if req.GetVolumeId() == "" {
		return errors.New("volume ID missing in request")
	}

	if req.GetVolumePath() == "" {
		return errors.New("volume path missing in request")
	}

	return nil
}

func validateNodePublishVolumeRequest(req *csi.NodePublishVolumeRequest) error {
	if req.GetVolumeId() == "" {
		return errors.New("volume ID missing in request")
	}

	if req.GetVolumeCapability() == nil {
		return errors.New("volume capability missing in request")
	}

	if req.GetVolumeCapability().GetBlock() != nil {
		return errors.New("volume access type Block is unsupported")
	}

	if req.GetVolumeCapability().GetMount() == nil {
		return errors.New("volume access type must by Mount")
	}

	if req.GetVolumeCapability().GetAccessMode().GetMode() !=
		csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER {
		return fmt.Errorf("volume access mode must be ReadWriteMany")
	}

	return nil
}

func validateNodeUnpublishVolumeRequest(req *csi.NodeUnpublishVolumeRequest) error {
	if req.GetVolumeId() == "" {
		return errors.New("volume ID missing in request")
	}

	if req.GetTargetPath() == "" {
		return errors.New("target path missing in request")
	}

	return nil
}
