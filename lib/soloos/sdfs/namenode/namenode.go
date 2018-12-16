package namenode

import (
	"soloos/sdfs/metastg"
	snettypes "soloos/snet/types"
	"soloos/util/offheap"
)

type NameNode struct {
	offheapDriver *offheap.OffheapDriver
	metaStg       *metastg.MetaStg

	SRPCServer NameNodeSRPCServer
}

func (p *NameNode) Init(options NameNodeOptions,
	offheapDriver *offheap.OffheapDriver,
	metaStg *metastg.MetaStg,
) error {
	var err error

	p.offheapDriver = offheapDriver
	p.metaStg = metaStg

	err = p.SRPCServer.Init(p, options.SRPCServer)
	if err != nil {
		return err
	}

	return nil
}

func (p *NameNode) RegisterDataNode(peerID *snettypes.PeerID, serveAddr string) (snettypes.PeerUintptr, error) {
	return p.metaStg.RegisterDataNode(peerID, serveAddr)
}

func (p *NameNode) Serve() error {
	return p.SRPCServer.Serve()
}

func (p *NameNode) Close() error {
	return p.SRPCServer.Close()
}
