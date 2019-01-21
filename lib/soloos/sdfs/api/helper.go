package api

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
)

// DataNode
type GetDataNode func(peerID snettypes.PeerID) snettypes.PeerUintptr
type ChooseDataNodesForNewNetBlock func(uNetINode types.NetINodeUintptr) (snettypes.PeerUintptrArray8, error)

// NetINode
type GetNetINodeWithReadAcquire func(netINodeID types.NetINodeID) (types.NetINodeUintptr, error)
type MustGetNetINodeWithReadAcquire func(netINodeID types.NetINodeID,
	size uint64, netBlockCap int, memBlockCap int) (types.NetINodeUintptr, error)
type PrepareNetINodeMetaDataOnlyLoadDB func(uNetINode types.NetINodeUintptr) error
type PrepareNetINodeMetaDataWithStorDB func(uNetINode types.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int) error

// FsINode
type AllocFsINodeID func() types.FsINodeID
type DeleteFsINodeByIDInDB func(fsINodeID types.FsINodeID) error
type ListFsINodeByParentIDFromDB func(parentID types.FsINodeID,
	isFetchAllCols bool,
	beforeLiteralFunc func(resultCount int) (fetchRowsLimit uint64, fetchRowsOffset uint64),
	literalFunc func(types.FsINode) bool,
) error
type UpdateFsINodeInDB func(fsINode types.FsINode) error
type InsertFsINodeInDB func(fsINode types.FsINode) error
type GetFsINodeByIDFromDB func(fsINodeID types.FsINodeID) (types.FsINode, error)
type GetFsINodeByNameFromDB func(parentID types.FsINodeID, fsINodeName string) (types.FsINode, error)

// FsINodeXAttr
type DeleteFIXAttrInDB func(fsINodeID types.FsINodeID) error
type ReplaceFIXAttrInDB func(fsINodeID types.FsINodeID, xattr types.FsINodeXAttr) error
type GetFIXAttrByInoFromDB func(fsINodeID types.FsINodeID) (types.FsINodeXAttr, error)
