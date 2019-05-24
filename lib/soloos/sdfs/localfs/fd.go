package localfs

import (
	"os"
	"soloos/common/sdfsapitypes"
	"sync/atomic"
)

type Fd struct {
	accessor  int32
	uNetINode sdfsapitypes.NetINodeUintptr
	filePath  string
	file      *os.File
}

func (p *Fd) Init(uNetINode sdfsapitypes.NetINodeUintptr, filePath string) error {
	var err error
	p.uNetINode = uNetINode
	p.filePath = filePath
	p.file, err = os.OpenFile(p.filePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}

	return nil
}

func (p *Fd) BorrowResource() {
	atomic.AddInt32(&p.accessor, 1)
}

func (p *Fd) ReturnResource() {
	atomic.AddInt32(&p.accessor, -1)
}
