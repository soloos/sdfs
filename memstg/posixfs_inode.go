package memstg

import (
	"soloos/common/fsapi"
	"soloos/common/solofstypes"
)

func (p *PosixFs) FetchFsINodeByID(pFsINodeMeta *solofstypes.FsINodeMeta,
	fsINodeIno solofstypes.FsINodeIno) error {
	var (
		uFsINode solofstypes.FsINodeUintptr
		err      error
	)
	uFsINode, err = p.FsINodeDriver.GetFsINodeByID(fsINodeIno)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}

	*pFsINodeMeta = uFsINode.Ptr().Meta
	return nil
}

func (p *PosixFs) FetchFsINodeByName(pFsINodeMeta *solofstypes.FsINodeMeta,
	parentID solofstypes.FsINodeIno, fsINodeName string) error {
	var (
		uFsINode solofstypes.FsINodeUintptr
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

func (p *PosixFs) FetchFsINodeByIDThroughHardLink(pFsINodeMeta *solofstypes.FsINodeMeta,
	fsINodeIno solofstypes.FsINodeIno) error {
	var (
		uFsINode solofstypes.FsINodeUintptr
		err      error
	)
	uFsINode, err = p.FsINodeDriver.GetFsINodeByIDThroughHardLink(fsINodeIno)
	defer p.FsINodeDriver.ReleaseFsINode(uFsINode)
	if err != nil {
		return err
	}

	*pFsINodeMeta = uFsINode.Ptr().Meta
	return nil
}

func (p *PosixFs) createFsINode(pFsINodeMeta *solofstypes.FsINodeMeta,
	fsINodeIno *solofstypes.FsINodeIno,
	netINodeID *solofstypes.NetINodeID, parentID solofstypes.FsINodeIno,
	name string, fsINodeType int, mode uint32,
	uid uint32, gid uint32, rdev uint32,
) error {
	var (
		err error
	)
	err = p.FsINodeDriver.PrepareFsINodeForCreate(pFsINodeMeta,
		fsINodeIno, netINodeID, parentID,
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

func (p *PosixFs) SimpleOpen(fsINodeMeta *solofstypes.FsINodeMeta,
	flags uint32, out *fsapi.OpenOut) error {
	out.Fh = p.FdTable.AllocFd(fsINodeMeta.Ino)
	out.OpenFlags = flags
	return nil
}

func (p *PosixFs) Mknod(input *fsapi.MknodIn, name string, out *fsapi.EntryOut) fsapi.Status {
	if len([]byte(name)) > solofstypes.FS_MAX_NAME_LENGTH {
		return solofstypes.FS_ENAMETOOLONG
	}

	var (
		parentFsINodeMeta solofstypes.FsINodeMeta
		fsINodeMeta       solofstypes.FsINodeMeta
		fsINodeType       int
		err               error
	)

	fsINodeType = FsModeToFsINodeType(input.Mode)
	if fsINodeType == solofstypes.FSINODE_TYPE_UNKOWN {
		return fsapi.EIO
	}

	err = p.FetchFsINodeByIDThroughHardLink(&parentFsINodeMeta, input.NodeId)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	err = p.createFsINode(&fsINodeMeta,
		nil, nil, parentFsINodeMeta.Ino,
		name, fsINodeType, input.Mode,
		input.Uid, input.Gid, input.Rdev)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(input.NodeId)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(out, &fsINodeMeta)

	return fsapi.OK
}

func (p *PosixFs) Unlink(header *fsapi.InHeader, name string) fsapi.Status {
	var (
		fsINodeMeta solofstypes.FsINodeMeta
		err         error
	)

	err = p.FetchFsINodeByName(&fsINodeMeta, header.NodeId, name)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	err = p.FsINodeDriver.UnlinkFsINode(fsINodeMeta.Ino)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	err = p.RefreshFsINodeACMtimeByIno(header.NodeId)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	return fsapi.OK
}

func (p *PosixFs) Fsync(input *fsapi.FsyncIn) fsapi.Status {
	var (
		fsINodeMeta solofstypes.FsINodeMeta
		err         error
	)

	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, input.NodeId)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	// TODO flush metadata

	return fsapi.OK
}

func (p *PosixFs) Lookup(header *fsapi.InHeader, name string, out *fsapi.EntryOut) fsapi.Status {
	if len([]byte(name)) > solofstypes.FS_MAX_NAME_LENGTH {
		return solofstypes.FS_ENAMETOOLONG
	}

	var (
		fsINodeMeta solofstypes.FsINodeMeta
		err         error
	)

	err = p.FetchFsINodeByName(&fsINodeMeta, header.NodeId, name)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	err = p.FetchFsINodeByIDThroughHardLink(&fsINodeMeta, fsINodeMeta.Ino)
	if err != nil {
		return ErrorToFsStatus(err)
	}

	p.SetFsEntryOutByFsINode(out, &fsINodeMeta)
	return fsapi.OK
}

func (p *PosixFs) Access(input *fsapi.AccessIn) fsapi.Status {
	return fsapi.OK
}

func (p *PosixFs) Forget(nodeid, nlookup uint64) {
}

func (p *PosixFs) Release(input *fsapi.ReleaseIn) {
}

func (p *PosixFs) CheckPermissionChmod(uid uint32, gid uint32,
	fsINodeMeta *solofstypes.FsINodeMeta) bool {

	if uid == 0 || uid == fsINodeMeta.Uid {
		return true
	}

	return false
}

func (p *PosixFs) CheckPermissionRead(uid uint32, gid uint32,
	fsINodeMeta *solofstypes.FsINodeMeta) bool {

	perm := uint32(07777) & fsINodeMeta.Mode
	if uid == fsINodeMeta.Uid {
		if perm&solofstypes.FS_PERM_USER_READ != 0 {
			return true
		}
	}

	if gid == fsINodeMeta.Gid {
		if perm&solofstypes.FS_PERM_GROUP_READ != 0 {
			return true
		}
	}

	if perm&solofstypes.FS_PERM_OTHER_READ != 0 {
		return true
	}

	return false
}

func (p *PosixFs) CheckPermissionWrite(uid uint32, gid uint32,
	fsINodeMeta *solofstypes.FsINodeMeta) bool {

	perm := uint32(07777) & fsINodeMeta.Mode
	if uid == fsINodeMeta.Uid {
		if perm&solofstypes.FS_PERM_USER_WRITE != 0 {
			return true
		}
	}

	if gid == fsINodeMeta.Gid {
		if perm&solofstypes.FS_PERM_GROUP_WRITE != 0 {
			return true
		}
	}

	if perm&solofstypes.FS_PERM_OTHER_WRITE != 0 {
		return true
	}

	return false
}

func (p *PosixFs) CheckPermissionExecute(uid uint32, gid uint32,
	fsINodeMeta *solofstypes.FsINodeMeta) bool {

	perm := uint32(07777) & fsINodeMeta.Mode
	if uid == fsINodeMeta.Uid {
		if perm&solofstypes.FS_PERM_USER_EXECUTE != 0 {
			return true
		}
	}

	if gid == fsINodeMeta.Gid {
		if perm&solofstypes.FS_PERM_GROUP_EXECUTE != 0 {
			return true
		}
	}

	if perm&solofstypes.FS_PERM_OTHER_EXECUTE != 0 {
		return true
	}

	return false
}

func (p *PosixFs) RefreshFsINodeACMtimeByIno(fsINodeIno solofstypes.FsINodeIno) error {
	return p.FsINodeDriver.RefreshFsINodeACMtimeByIno(fsINodeIno)
}

func (p *PosixFs) TruncateINode(pFsINode *solofstypes.FsINodeMeta, size uint64) error {
	var (
		uFsINode solofstypes.FsINodeUintptr
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
