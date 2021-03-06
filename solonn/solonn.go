package solonn

import (
	"soloos/common/iron"
	"soloos/common/snet"
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"soloos/solofs/memstg"
	"soloos/solofs/metastg"
)

type Solonn struct {
	*soloosbase.SoloosEnv
	srpcPeer snet.Peer
	webPeer  snet.Peer
	metaStg  *metastg.MetaStg

	memBlockDriver *memstg.MemBlockDriver
	netBlockDriver *memstg.NetBlockDriver
	netINodeDriver *memstg.NetINodeDriver

	heartBeatServerOptionsArr []snet.HeartBeatServerOptions
	serverCount               int
	srpcServer                SrpcServer
	webServer                 WebServer
	serverDriver              iron.ServerDriver
}

func (p *Solonn) initSNetPeer(
	srpcPeerID snet.PeerID, srpcServerServeAddr string,
	webPeerID snet.PeerID, webServerServeAddr string,
) error {
	var err error

	p.srpcPeer.ID = srpcPeerID
	p.srpcPeer.SetAddress(srpcServerServeAddr)
	p.srpcPeer.ServiceProtocol = solofstypes.DefaultSolofsRPCProtocol
	err = p.SNetDriver.RegisterPeer(p.srpcPeer)
	if err != nil {
		return err
	}

	p.webPeer.ID = webPeerID
	p.webPeer.SetAddress(webServerServeAddr)
	p.webPeer.ServiceProtocol = snet.ProtocolWeb
	err = p.SNetDriver.RegisterPeer(p.webPeer)
	if err != nil {
		return err
	}

	return nil
}

func (p *Solonn) Init(soloosEnv *soloosbase.SoloosEnv,
	srpcPeerID snet.PeerID,
	srpcServerListenAddr string,
	srpcServerServeAddr string,
	webPeerID snet.PeerID,
	webServerOptions iron.Options,
	metaStg *metastg.MetaStg,
	memBlockDriver *memstg.MemBlockDriver,
	netBlockDriver *memstg.NetBlockDriver,
	netINodeDriver *memstg.NetINodeDriver,
) error {
	var err error

	p.SoloosEnv = soloosEnv
	p.metaStg = metaStg
	p.memBlockDriver = memBlockDriver
	p.netBlockDriver = netBlockDriver
	p.netINodeDriver = netINodeDriver

	err = p.srpcServer.Init(p, srpcServerListenAddr, srpcServerServeAddr)
	if err != nil {
		return err
	}

	err = p.webServer.Init(p, webServerOptions)
	if err != nil {
		return err
	}

	err = p.serverDriver.Init(&p.srpcServer, &p.webServer)
	if err != nil {
		return err
	}

	err = p.initSNetPeer(srpcPeerID, srpcServerServeAddr, webPeerID, webServerOptions.ServeStr)
	if err != nil {
		return err
	}

	return nil
}

func (p *Solonn) Serve() error {
	var err error

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

func (p *Solonn) Close() error {
	var err error
	err = p.serverDriver.Close()
	if err != nil {
		return err
	}

	return nil
}
