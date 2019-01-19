package api

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
)

type GetNetINodeWithReadAcquire func(netINodeID types.NetINodeID) (types.NetINodeUintptr, error)
type MustGetNetINodeWithReadAcquire func(netINodeID types.NetINodeID,
	size uint64, netBlockCap int, memBlockCap int) (types.NetINodeUintptr, error)
type GetDataNode func(peerID *snettypes.PeerID) snettypes.PeerUintptr
type ChooseDataNodesForNewNetBlock func(uNetINode types.NetINodeUintptr,
	backends *snettypes.PeerUintptrArray8) error

type PrepareNetINodeMetaDataOnlyLoadDB func(uNetINode types.NetINodeUintptr) error
type PrepareNetINodeMetaDataWithStorDB func(uNetINode types.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int) error

type NetINodeDriverHelper struct {
	NameNodeClient                    *NameNodeClient
	PrepareNetINodeMetaDataOnlyLoadDB PrepareNetINodeMetaDataOnlyLoadDB
	PrepareNetINodeMetaDataWithStorDB PrepareNetINodeMetaDataWithStorDB
}
