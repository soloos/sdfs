package memstg

import (
	"soloos/common/solodbtypes"
	"soloos/common/solofsprotocol"
	"soloos/common/solofstypes"
)

func (p *NetINodeDriver) doGetNetINodeMetaData(isMustGet bool,
	uNetINode solofstypes.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int,
) error {
	var (
		req  solofsprotocol.NetINodeInfoReq
		resp solofsprotocol.NetINodeInfoResp
		err  error
	)

	req.NetINodeID = uNetINode.Ptr().ID
	req.Size = size
	req.NetBlockCap = int32(netBlockCap)
	req.MemBlockCap = int32(memBlockCap)

	if isMustGet {
		err = p.solonnClient.Dispatch("/NetINode/MustGet", &resp, &req)
	} else {
		err = p.solonnClient.Dispatch("/NetINode/Get", &resp, &req)
	}
	if err != nil {
		return err
	}

	uNetINode.Ptr().Size = resp.Size
	uNetINode.Ptr().NetBlockCap = int(resp.NetBlockCap)
	uNetINode.Ptr().MemBlockCap = int(resp.MemBlockCap)

	return nil
}

func (p *NetINodeDriver) getNetINodeMetaData(uNetINode solofstypes.NetINodeUintptr) error {
	return p.doGetNetINodeMetaData(false, uNetINode, 0, 0, 0)
}

func (p *NetINodeDriver) mustGetNetINodeMetaData(uNetINode solofstypes.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int,
) error {
	return p.doGetNetINodeMetaData(true, uNetINode, size, netBlockCap, memBlockCap)
}

func (p *NetINodeDriver) prepareNetINodeMetaDataCommon(pNetINode *solofstypes.NetINode) {
	pNetINode.MemBlockPlacementPolicy.SetType(solofstypes.BlockPlacementPolicyDefault)
	pNetINode.IsDBMetaDataInited.Store(solodbtypes.MetaDataStateInited)
}

func (p *NetINodeDriver) PrepareNetINodeMetaDataOnlyLoadDB(uNetINode solofstypes.NetINodeUintptr) error {
	var err error

	err = p.getNetINodeMetaData(uNetINode)
	if err != nil {
		return err
	}

	p.prepareNetINodeMetaDataCommon(uNetINode.Ptr())

	return nil
}

func (p *NetINodeDriver) PrepareNetINodeMetaDataWithStorDB(uNetINode solofstypes.NetINodeUintptr,
	size uint64, netBlockCap int, memBlockCap int) error {
	var err error

	err = p.mustGetNetINodeMetaData(uNetINode, size, netBlockCap, memBlockCap)
	if err != nil {
		return err
	}

	p.prepareNetINodeMetaDataCommon(uNetINode.Ptr())

	return nil
}

func (p *NetINodeDriver) NetINodeCommitSizeInDB(uNetINode solofstypes.NetINodeUintptr, size uint64) error {
	var err error
	var req = solofsprotocol.NetINodeCommitSizeInDBReq{
		NetINodeID: uNetINode.Ptr().ID,
		Size:       size,
	}

	err = p.solonnClient.Dispatch("/NetINode/CommitSizeInDB", nil, req)
	if err != nil {
		return err
	}

	uNetINode.Ptr().Size = size
	return nil
}
