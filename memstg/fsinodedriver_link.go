package memstg

import (
	"soloos/common/fsapi"
	"soloos/common/solofstypes"
	"sync/atomic"
)

func (p *FsINodeDriver) Link(srcFsINodeIno solofstypes.FsINodeIno,
	parentID solofstypes.FsINodeIno, filename string,
	retFsINode *solofstypes.FsINodeMeta) error {
	var (
		uSrcFsINode solofstypes.FsINodeUintptr
		pSrcFsINode *solofstypes.FsINode
		err         error
	)

	uSrcFsINode, err = p.GetFsINodeByID(srcFsINodeIno)
	defer p.ReleaseFsINode(uSrcFsINode)
	pSrcFsINode = uSrcFsINode.Ptr()
	if err != nil {
		return err
	}

	err = p.PrepareFsINodeForCreate(retFsINode,
		nil, nil, parentID,
		filename, solofstypes.FSINODE_TYPE_HARD_LINK, pSrcFsINode.Meta.Mode,
		0, 0, solofstypes.FS_RDEV)
	if err != nil {
		return err
	}
	retFsINode.HardLinkIno = pSrcFsINode.Meta.Ino
	err = p.CreateFsINode(retFsINode)
	if err != nil {
		return err
	}

	atomic.AddInt32(&pSrcFsINode.Meta.Nlink, 1)
	err = p.UpdateFsINode(&pSrcFsINode.Meta)
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) Symlink(parentID solofstypes.FsINodeIno, pointedTo string, linkName string,
	retFsINodeMeta *solofstypes.FsINodeMeta) error {
	var (
		err error
	)

	err = p.PrepareFsINodeForCreate(retFsINodeMeta,
		nil, nil, parentID,
		linkName, solofstypes.FSINODE_TYPE_SOFT_LINK, fsapi.S_IFLNK|0777,
		0, 0, solofstypes.FS_RDEV)
	if err != nil {
		return err
	}
	err = p.CreateFsINode(retFsINodeMeta)
	if err != nil {
		return err
	}

	err = p.posixFs.FIXAttrDriver.SetXAttr(retFsINodeMeta.Ino, solofstypes.FS_XATTR_SOFT_LNKMETA_KEY, []byte(pointedTo))
	if err != nil {
		return err
	}

	return nil
}

func (p *FsINodeDriver) Readlink(fsINodeIno solofstypes.FsINodeIno) ([]byte, error) {
	return p.posixFs.FIXAttrDriver.GetXAttrData(fsINodeIno, solofstypes.FS_XATTR_SOFT_LNKMETA_KEY)
}

func (p *FsINodeDriver) deleteFsINodeHardLinked(fsINodeIno solofstypes.FsINodeIno) error {
	var (
		uFsINode                solofstypes.FsINodeUintptr
		isFsINodeDeleted        bool
		isShouldReleaseUFsINode bool
		err                     error
	)
	uFsINode, err = p.GetFsINodeByID(fsINodeIno)
	if err != nil {
		isShouldReleaseUFsINode = true

	} else {
		isFsINodeDeleted, err = p.decreaseFsINodeNLink(uFsINode)
		if err != nil {
			isShouldReleaseUFsINode = true
		}

		if isFsINodeDeleted == false {
			isShouldReleaseUFsINode = true
		}
	}

	if isShouldReleaseUFsINode {
		p.ReleaseFsINode(uFsINode)
	}

	return err
}

// decreaseFsINodeNLink return isFsINodeDeleted, decreaseError
// if FsINodeMeta.Nlink == 0 then delete FsINode else decrease(FsINode.Nlink)
func (p *FsINodeDriver) decreaseFsINodeNLink(uFsINode solofstypes.FsINodeUintptr) (bool, error) {
	var (
		pFsINode = uFsINode.Ptr()
		err      error
	)

	if atomic.AddInt32(&pFsINode.Meta.Nlink, -1) > 0 {
		err = p.UpdateFsINode(&pFsINode.Meta)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	// assert fsINode.Nlink == 0
	if pFsINode.Meta.Type == solofstypes.FSINODE_TYPE_HARD_LINK {
		err = p.deleteFsINodeHardLinked(pFsINode.Meta.HardLinkIno)
		if err != nil {
			return false, err
		}

	}

	err = p.helper.DeleteFsINodeByIDInDB(p.posixFs.NameSpaceID, pFsINode.Meta.Ino)
	if err != nil {
		return false, err
	}

	p.DeleteFsINodeCache(uFsINode, pFsINode.Meta.ParentID, pFsINode.Meta.Name())

	return true, nil
}

func (p *FsINodeDriver) UnlinkFsINode(fsINodeIno solofstypes.FsINodeIno) error {
	var (
		uFsINode                solofstypes.FsINodeUintptr
		pFsINode                *solofstypes.FsINode
		oldFsINodeParentID      solofstypes.FsINodeIno
		isFsINodeDeleted        bool
		isShouldReleaseUFsINode bool
		err                     error
	)

	uFsINode, err = p.GetFsINodeByID(fsINodeIno)
	pFsINode = uFsINode.Ptr()

	if err != nil {
		if err.Error() == solofstypes.ErrObjectNotExists.Error() {
			err = nil
			goto DONE
		} else {
			goto DONE
		}
	}

	oldFsINodeParentID = pFsINode.Meta.ParentID
	isFsINodeDeleted, err = p.decreaseFsINodeNLink(uFsINode)
	if err != nil {
		return err
	}

	if isFsINodeDeleted == false {
		pFsINode.Meta.ParentID = solofstypes.ZombieFsINodeParentID
		err = p.UpdateFsINode(&pFsINode.Meta)
		if err != nil {
			goto DONE
		}
		isShouldReleaseUFsINode = true
	}

DONE:
	if uFsINode != 0 {
		p.CleanFsINodeAssitCache(oldFsINodeParentID, uFsINode.Ptr().Meta.Name())
		if isShouldReleaseUFsINode {
			p.ReleaseFsINode(uFsINode)
		}
	}
	return nil
}
