package memstg

import (
	fsapitypes "soloos/common/fsapi/types"
	"soloos/sdfs/types"
)

func (p *DirTreeStg) SimpleFlush(fsINodeID types.FsINodeID) error {
	var (
		uFsINode types.FsINodeUintptr
		pFsINode *types.FsINode
		err      error
	)

	uFsINode, err = p.FsINodeDriver.GetFsINodeByIDThroughHardLink(fsINodeID)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	pFsINode = uFsINode.Ptr()
	if err != nil {
		return err
	}

	if pFsINode.UNetINode != 0 {
		err = p.MemStg.NetINodeDriver.Sync(pFsINode.UNetINode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *DirTreeStg) Flush(input *fsapitypes.FlushIn) fsapitypes.Status {
	var err = p.SimpleFlush(input.NodeId)
	return types.ErrorToFsStatus(err)
}
