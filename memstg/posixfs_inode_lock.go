package memstg

import (
	"soloos/common/fsapi"
	"soloos/common/solofstypes"
)

func (p *PosixFs) setLKOutByMeta(out *fsapi.LkOut, meta *solofstypes.INodeRWMutexMeta) {
	out.Lk.Start = meta.Start
	out.Lk.End = meta.End
	out.Lk.Typ = meta.Typ
	out.Lk.Pid = meta.Pid
}

func (p *PosixFs) setMetaByLKIn(lkIn *fsapi.LkIn, meta *solofstypes.INodeRWMutexMeta) {
	meta.Start = lkIn.Lk.Start
	meta.End = lkIn.Lk.End
	meta.Typ = lkIn.Lk.Typ
	meta.Pid = lkIn.Lk.Pid
}

// GetLk returns existing lock information for file
func (p *PosixFs) GetLk(input *fsapi.LkIn, out *fsapi.LkOut) fsapi.Status {
	var (
		meta solofstypes.INodeRWMutexMeta
		err  error
	)
	err = p.FsMutexDriver.GetLk(input.NodeId, &meta)
	if err != nil {
		return ErrorToFsStatus(err)
	}
	p.setLKOutByMeta(out, &meta)
	return fsapi.OK
}

// SetLk Sets or clears the lock described by lk on file.
func (p *PosixFs) SetLk(input *fsapi.LkIn) fsapi.Status {
	var (
		meta solofstypes.INodeRWMutexMeta
		err  error
	)
	p.setMetaByLKIn(input, &meta)
	err = p.FsMutexDriver.SetLk(input.NodeId, &meta)
	if err != nil {
		return ErrorToFsStatus(err)
	}
	return fsapi.OK
}

// SetLkw Sets or clears the lock described by lk. This call blocks until the operation can be completed.
func (p *PosixFs) SetLkw(input *fsapi.LkIn) fsapi.Status {
	var (
		meta solofstypes.INodeRWMutexMeta
		err  error
	)
	p.setMetaByLKIn(input, &meta)
	err = p.FsMutexDriver.SetLkw(input.NodeId, &meta)
	if err != nil {
		return ErrorToFsStatus(err)
	}
	return fsapi.OK
}
