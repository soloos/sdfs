package solodn

import (
	"fmt"
	"soloos/common/iron"
	"soloos/common/log"
	"soloos/common/snet"
	"soloos/common/solofsapi"
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"soloos/solofs/localfs"
	"soloos/solofs/memstg"
)

type Solodn struct {
	*soloosbase.SoloosEnv
	srpcPeer snet.Peer
	webPeer  snet.Peer

	solonnClient   solofsapi.SolonnClient
	solodnClient   solofsapi.SolodnClient
	memBlockDriver *memstg.MemBlockDriver
	netBlockDriver *memstg.NetBlockDriver
	netINodeDriver *memstg.NetINodeDriver

	localFs         localfs.LocalFs
	localFsSNetPeer snet.Peer

	heartBeatServerOptionsArr []snet.HeartBeatServerOptions
	srpcServer                SrpcServer
	webServer                 WebServer
	serverDriver              iron.ServerDriver
}

func (p *Solodn) initLocalFs(options SolodnOptions) error {
	var err error
	err = p.localFs.Init(options.LocalFsRoot)
	if err != nil {
		return err
	}

	p.localFsSNetPeer.ID = snet.MakeSysPeerID(fmt.Sprintf("SOLODN_LOCAL_FS"))
	p.localFsSNetPeer.SetAddress("LocalFs")
	p.localFsSNetPeer.ServiceProtocol = snet.ProtocolLocalFs
	err = p.SNetDriver.RegisterPeer(p.localFsSNetPeer)
	if err != nil {
		return err
	}

	return nil
}

func (p *Solodn) initSNetPeer(options SolodnOptions) error {
	var err error

	p.srpcPeer.ID = options.SrpcPeerID
	p.srpcPeer.SetAddress(options.SrpcServerServeAddr)
	p.srpcPeer.ServiceProtocol = solofstypes.DefaultSolofsRPCProtocol
	err = p.SNetDriver.RegisterPeer(p.srpcPeer)
	if err != nil {
		return err
	}

	p.webPeer.ID = options.WebPeerID
	p.webPeer.SetAddress(options.WebServer.ServeStr)
	p.webPeer.ServiceProtocol = snet.ProtocolWeb
	err = p.SNetDriver.RegisterPeer(p.webPeer)
	if err != nil {
		return err
	}

	return nil
}

func (p *Solodn) initNetBlockDriver() error {
	p.netBlockDriver.SetSolonnClient(&p.solonnClient)
	p.netBlockDriver.SetSolodnClient(&p.solodnClient)
	p.netBlockDriver.SetHelper(
		p.netBlockDriver.PrepareNetBlockMetaData,
		p.localFs.PReadMemBlockWithDisk,
		p.localFs.UploadMemBlockWithDisk,
		nil,
	)
	return nil
}

func (p *Solodn) initNetINodeDriver() error {
	p.netINodeDriver.SetSolonnClient(&p.solonnClient)
	p.netBlockDriver.SetSolodnClient(&p.solodnClient)
	p.netINodeDriver.SetHelper(
		p.netINodeDriver.PrepareNetINodeMetaDataOnlyLoadDB,
		p.netINodeDriver.PrepareNetINodeMetaDataWithStorDB,
		p.netINodeDriver.NetINodeCommitSizeInDB,
	)
	return nil
}

func (p *Solodn) Init(soloosEnv *soloosbase.SoloosEnv,
	options SolodnOptions,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *memstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) error {
	var err error

	p.SoloosEnv = soloosEnv

	err = p.solonnClient.Init(p.SoloosEnv, options.SolonnSrpcPeerID)
	if err != nil {
		log.Warn("Solodn Init solonnClient.Init failed, err:", err)
		return err
	}

	err = p.solodnClient.Init(p.SoloosEnv)
	if err != nil {
		log.Warn("Solodn Init solodnClient.Init failed, err:", err)
		return err
	}

	p.memBlockDriver = memBlockDriver

	p.netBlockDriver = netBlockDriver
	err = p.initNetBlockDriver()
	if err != nil {
		log.Warn("Solodn Init initNetBlockDriver failed, err:", err)
		return err
	}

	p.netINodeDriver = netINodeDriver
	err = p.initNetINodeDriver()
	if err != nil {
		log.Warn("Solodn Init initNetINodeDriver failed, err:", err)
		return err
	}

	err = p.initSNetPeer(options)
	if err != nil {
		log.Warn("Solodn Init initSNetPeer failed, err:", err)
		return err
	}

	err = p.srpcServer.Init(p, options.SrpcServerListenAddr, options.SrpcServerServeAddr)
	if err != nil {
		log.Warn("Solodn Init SrpcServer.Init failed, err:", err)
		return err
	}

	err = p.webServer.Init(p, options.WebServer)
	if err != nil {
		log.Warn("Solodn Init WebServer.Init failed, err:", err)
		return err
	}

	err = p.serverDriver.Init(&p.srpcServer, &p.webServer)
	if err != nil {
		return err
	}

	err = p.initLocalFs(options)
	if err != nil {
		log.Warn("Solodn Init initLocalFs failed, err:", err)
		return err
	}

	return nil
}

func (p *Solodn) Serve() error {
	var err error

	err = p.RegisterInSolonn()
	if err != nil {
		return err
	}

	err = p.StartHeartBeat()
	if err != nil {
		return err
	}

	err = p.serverDriver.Serve()
	if err != nil {
		return err
	}

	return nil
}

func (p *Solodn) Close() error {
	var err error
	err = p.serverDriver.Close()
	if err != nil {
		return err
	}

	return nil
}
