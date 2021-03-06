package main

import (
	"soloos/common/snet"
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/solofs/solofssdk"
)

var (
	env Env
)

type Env struct {
	Options            Options
	SoloosEnv          soloosbase.SoloosEnv
	solofsClientDriver solofssdk.ClientDriver
	solofsClient       solofssdk.Client
}

func (p *Env) Init(optionsFile string) {
	var (
		err error
	)

	p.Options, err = LoadOptionsFile(optionsFile)
	util.AssertErrIsNil(err)

	err = p.SoloosEnv.InitWithSNet(p.Options.SNetDriverServeAddr)
	util.AssertErrIsNil(err)

	go func() {
		util.PProfServe(p.Options.PProfListenAddr)
	}()

	var solonnSrpcPeerID snet.PeerID
	solonnSrpcPeerID.SetStr(p.Options.SolonnSrpcPeerID)
	util.AssertErrIsNil(p.solofsClientDriver.Init(&p.SoloosEnv,
		solonnSrpcPeerID,
		p.Options.DBDriver, p.Options.Dsn,
	))

	if p.Options.DefaultNetBlockCap == 0 {
		p.Options.DefaultNetBlockCap = p.Options.DefaultMemBlockCap
	}

	util.AssertErrIsNil(
		p.solofsClientDriver.InitClient(&p.solofsClient,
			solofstypes.NameSpaceID(p.Options.NameSpaceID),
			p.Options.DefaultNetBlockCap,
			p.Options.DefaultMemBlockCap,
			p.Options.DefaultMemBlocksLimit,
		))
}
