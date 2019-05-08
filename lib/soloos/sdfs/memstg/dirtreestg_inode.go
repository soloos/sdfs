package memstg

import (
	fsapitypes "soloos/common/fsapi/types"
	sdfsapitypes "soloos/common/sdfsapi/types"
	"soloos/sdfs/types"
)

func (p *DirTreeStg) FetchFsINodeByID(pFsINodeMeta *sdfsapitypes.FsINodeMeta,
	fsINodeID sdfsapitypes.FsINodeID) error {
	var (
		uFsINode sdfsapitypes.FsINodeUintptr
		err      error
	)
	uFsINode, err = p.FsINodeDriver.GetFsINodeByID(fsINodeID)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}

	*pFsINodeMeta = uFsINode.Ptr().Meta
	return nil
}

func (p *DirTreeStg) FetchFsINodeByName(pFsINodeMeta *sdfsapitypes.FsINodeMeta,
	parentID types.FsINodeID, fsINodeName string) error {
	var (
		uFsINode sdfsapitypes.FsINodeUintptr
		err      error
	)
	uFsINode, err = p.FsINodeDriver.GetFsINodeByName(parentID, fsINodeName)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}

	*pFsINodeMeta = uFsINode.Ptr().Meta
	return nil
}

func (p *DirTreeStg) FetchFsINodeByIDThroughHardLink(pFsINodeMeta *sdfsapitypes.FsINodeMeta,
	fsINodeID types.FsINodeID) error {
	var (
		uFsINode sdfsapitypes.FsINodeUintptr
		err      error
	)
	uFsINode, err = p.FsINodeDriver.GetFsINodeByIDThroughHardLink(fsINodeID)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}

	*pFsINodeMeta = uFsINode.Ptr().Meta
	return nil
}

func (p *DirTreeStg) createFsINode(pFsINodeMeta *sdfsapitypes.FsINodeMeta,
	fsINodeID *types.FsINodeID,
	netINodeID *types.NetINodeID, parentID types.FsINodeID,
	name string, fsINodeType int, mode uint32,
	uid uint32, gid uint32, rdev uint32,
) error {
	var (
		err error
	)
	err = p.FsINodeDriver.PrepareFsINodeForCreate(pFsINodeMeta,
		fsINodeID, netINodeID, parentID,
		name, fsINodeType, mode,
		uid, gid, rdev)
	if err != nil {
		return err
	}

	err = p.FsINodeDriver.CreateFsINode(pFsINodeMeta)
	if err != nil {
		return err
	}

	return nil
}

func (p *DirTreeStg) SimpleOpen(fsINodeMeta *sdfsapitypes.FsINodeMeta,
	flags uint32, out *fsapitypes.OpenOut) error {
	out.Fh = p.FdTable.AllocFd(fsINodeMeta.Ino)
	out.OpenFlags = flags
	return nil
}

func (p *DirTreeStg) Mknod(input *fsapitypes.MknodIn, name string, out *fsapitypes.EntryOut) fsapitypes.Status {
	if len([]byte(name)) > types.FS_MAX_NAME_LENGTH {
		return types.FS_ENAMETOOLONG
	}

	var (
		parentFsINodeMeta sdfsapitypes.FsINodeMeta
		fsINodeMeta       sdfsapitypes.FsINodeMeta
		fsINodeType       int
		err               error
	)

	fsINodeType = types.FsModeToFsINodeType(input.Mode)
	if fsINodeType == types.FSINODE_TYPE_UNKOWN {
		return fsapitypes.EIO
	}

	err = p.FetchFsINodeByIDThroughHardLink(&parentFsINodeMeta, input.NodeId)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.createFsINode(&fsINodeMeta,
		nil, nil, parentFsINodeMeta.Ino,
		name, fsINodeType, input.Mode,
		input.Uid, input.Gid, input.Rdev)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(input.NodeId)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(out, &fsINodeMeta)

	return fsapitypes.OK
}

func (p *DirTreeStg) Unlink(header *fsapitypes.InHeader, name string) fsapitypes.Status {
	var (
		fsINodeMeta sdfsapitypes.FsINodeMeta
		err         error
	)

	err = p.FetchFsINodeByName(&fsINodeMeta, header.NodeId, name)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.FsINodeDriver.UnlinkFsINode(fsINodeMeta.Ino)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(header.NodeId)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	return fsapitypes.OK
}

func (p *DirTreeStg) Fsync(input *fsapitypes.FsyncIn) fsapitypes.Status {
	var (
		fsINodeMeta sdfsapitypes.FsINodeMeta
		err         error
	)

	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, input.NodeId)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	// TODO flush metadata

	return fsapitypes.OK
}

func (p *DirTreeStg) Lookup(header *fsapitypes.InHeader, name string, out *fsapitypes.EntryOut) fsapitypes.Status {
	if len([]byte(name)) > types.FS_MAX_NAME_LENGTH {
		return types.FS_ENAMETOOLONG
	}

	var (
		fsINodeMeta sdfsapitypes.FsINodeMeta
		err         error
	)

	err = p.FetchFsINodeByName(&fsINodeMeta, header.NodeId, name)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, fsINodeMeta.Ino)
	if err != nil {
		return types.ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(out, &fsINodeMeta)
	return fsapitypes.OK
}

func (p *DirTreeStg) Access(input *fsapitypes.AccessIn) fsapitypes.Status {
	return fsapitypes.OK
}

func (p *DirTreeStg) Forget(nodeid, nlookup uint64) {
}

func (p *DirTreeStg) Release(input *fsapitypes.ReleaseIn) {
}

func (p *DirTreeStg) CheckPermissionChmod(uid uint32, gid uint32,
	fsINodeMeta *sdfsapitypes.FsINodeMeta) bool {

	if uid == 0 || uid == fsINodeMeta.Uid {
		return true
	}

	return false
}

func (p *DirTreeStg) CheckPermissionRead(uid uint32, gid uint32,
	fsINodeMeta *sdfsapitypes.FsINodeMeta) bool {

	perm := uint32(07777) & fsINodeMeta.Mode
	if uid == fsINodeMeta.Uid {
		if perm&types.FS_PERM_USER_READ != 0 {
			return true
		}
	}

	if gid == fsINodeMeta.Gid {
		if perm&types.FS_PERM_GROUP_READ != 0 {
			return true
		}
	}

	if perm&types.FS_PERM_OTHER_READ != 0 {
		return true
	}

	return false
}

func (p *DirTreeStg) CheckPermissionWrite(uid uint32, gid uint32,
	fsINodeMeta *sdfsapitypes.FsINodeMeta) bool {

	perm := uint32(07777) & fsINodeMeta.Mode
	if uid == fsINodeMeta.Uid {
		if perm&types.FS_PERM_USER_WRITE != 0 {
			return true
		}
	}

	if gid == fsINodeMeta.Gid {
		if perm&types.FS_PERM_GROUP_WRITE != 0 {
			return true
		}
	}

	if perm&types.FS_PERM_OTHER_WRITE != 0 {
		return true
	}

	return false
}

func (p *DirTreeStg) CheckPermissionExecute(uid uint32, gid uint32,
	fsINodeMeta *sdfsapitypes.FsINodeMeta) bool {

	perm := uint32(07777) & fsINodeMeta.Mode
	if uid == fsINodeMeta.Uid {
		if perm&types.FS_PERM_USER_EXECUTE != 0 {
			return true
		}
	}

	if gid == fsINodeMeta.Gid {
		if perm&types.FS_PERM_GROUP_EXECUTE != 0 {
			return true
		}
	}

	if perm&types.FS_PERM_OTHER_EXECUTE != 0 {
		return true
	}

	return false
}

func (p *DirTreeStg) RefreshFsINodeACMtimeByIno(fsINodeID types.FsINodeID) error {
	return p.FsINodeDriver.RefreshFsINodeACMtimeByIno(fsINodeID)
}

func (p *DirTreeStg) TruncateINode(pFsINode *sdfsapitypes.FsINodeMeta, size uint64) error {
	var (
		uFsINode sdfsapitypes.FsINodeUintptr
		err      error
	)
	uFsINode, err = p.FsINodeDriver.GetFsINodeByID(pFsINode.Ino)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}
	if uFsINode.Ptr().UNetINode == 0 {
		return nil
	}

	return p.MemStg.NetINodeDriver.NetINodeTruncate(uFsINode.Ptr().UNetINode, size)
}
