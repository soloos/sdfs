package main

import (
	"soloos/sdfs/api"
	"soloos/sdfs/datanode"
	"soloos/sdfs/memstg"
	"soloos/sdfs/metastg"
	"soloos/sdfs/namenode"
	"soloos/sdfs/netstg"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util"
	"soloos/util/offheap"
)

type Env struct {
	options          Options
	offheapDriver    *offheap.OffheapDriver
	SNetDriver       snet.NetDriver
	SNetClientDriver snet.ClientDriver
	MetaStg          metastg.MetaStg
	DataNodeClient   api.DataNodeClient
	MemBlockDriver   memstg.MemBlockDriver
	NetBlockDriver   netstg.NetBlockDriver
	NetINodeDriver   memstg.NetINodeDriver
}

func (p *Env) Init(options Options) {
	p.options = options
	p.offheapDriver = &offheap.DefaultOffheapDriver

	util.AssertErrIsNil(p.SNetDriver.Init(p.offheapDriver))
	util.AssertErrIsNil(p.SNetClientDriver.Init(p.offheapDriver))

	util.AssertErrIsNil(p.MetaStg.Init(p.offheapDriver,
		options.DBDriver, options.Dsn))

	p.DataNodeClient.Init(&p.SNetClientDriver)

	{
		var memBlockDriverOptions = memstg.MemBlockDriverOptions{
			[]memstg.MemBlockPoolOptions{
				memstg.MemBlockPoolOptions{
					p.options.DefaultMemBlockCap,
					p.options.DefaultMemBlocksLimit,
				},
			},
		}
		util.AssertErrIsNil(p.MemBlockDriver.Init(p.offheapDriver, memBlockDriverOptions))
	}
}

func (p *Env) startCommon() {
	if p.options.PProfListenAddr != "" {
		go util.PProfServe(p.options.PProfListenAddr)
	}
}

func (p *Env) startNameNode() {
	var (
		nameNode namenode.NameNode
	)

	util.AssertErrIsNil(p.NetBlockDriver.Init(p.offheapDriver,
		&p.SNetDriver, &p.SNetClientDriver,
		nil, &p.DataNodeClient, p.MetaStg.PrepareNetBlockMetaData))

	util.AssertErrIsNil(p.NetINodeDriver.Init(p.offheapDriver,
		&p.NetBlockDriver,
		&p.MemBlockDriver,
		nil,
		p.MetaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		p.MetaStg.PrepareNetINodeMetaDataWithStorDB,
		p.MetaStg.NetINodeCommitSizeInDB,
	))

	util.AssertErrIsNil(nameNode.Init(p.offheapDriver,
		p.options.ListenAddr,
		&p.MetaStg,
		&p.MemBlockDriver,
		&p.NetBlockDriver,
		&p.NetINodeDriver,
	))

	util.AssertErrIsNil(nameNode.Serve())
	util.AssertErrIsNil(nameNode.Close())
}

func (p *Env) startDataNode() {
	var (
		dataNodePeerID  snettypes.PeerID
		dataNode        datanode.DataNode
		nameNodePeerID  snettypes.PeerID
		dataNodeOptions datanode.DataNodeOptions
	)

	copy(dataNodePeerID[:], []byte(p.options.DataNodePeerIDStr))
	copy(nameNodePeerID[:], []byte(p.options.NameNodePeerIDStr))

	dataNodeOptions = datanode.DataNodeOptions{
		PeerID:               dataNodePeerID,
		SrpcServerListenAddr: p.options.ListenAddr,
		SrpcServerServeAddr:  p.options.ListenAddr,
		LocalFsRoot:          p.options.DataNodeLocalFsRoot,
		NameNodePeerID:       nameNodePeerID,
		NameNodeSRPCServer:   p.options.NameNodeAddr,
	}

	util.AssertErrIsNil(p.NetBlockDriver.Init(p.offheapDriver,
		&p.SNetDriver, &p.SNetClientDriver,
		nil, &p.DataNodeClient, p.MetaStg.PrepareNetBlockMetaData))

	util.AssertErrIsNil(p.NetINodeDriver.Init(p.offheapDriver,
		&p.NetBlockDriver,
		&p.MemBlockDriver,
		nil,
		p.MetaStg.PrepareNetINodeMetaDataOnlyLoadDB,
		p.MetaStg.PrepareNetINodeMetaDataWithStorDB,
		p.MetaStg.NetINodeCommitSizeInDB,
	))

	util.AssertErrIsNil(dataNode.Init(p.offheapDriver, dataNodeOptions,
		&p.SNetDriver, &p.SNetClientDriver,
		&p.MetaStg,
		&p.MemBlockDriver,
		&p.NetBlockDriver,
		&p.NetINodeDriver,
	))
	util.AssertErrIsNil(dataNode.Serve())
	util.AssertErrIsNil(dataNode.Close())
}

func (p *Env) Start() {
	if p.options.Mode == "namenode" {
		p.startCommon()
		p.startNameNode()
	}

	if p.options.Mode == "datanode" {
		p.startCommon()
		p.startDataNode()
	}
}