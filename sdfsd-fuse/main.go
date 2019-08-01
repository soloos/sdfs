package main

import (
	"os"
	"soloos/common/snettypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/sdfs/sdfssdk"
	"soloos/sdfs/sfuse"
)

func main() {
	optionsFile := os.Args[1]

	options, err := LoadOptionsFile(optionsFile)
	util.AssertErrIsNil(err)

	var (
		soloOSEnv    soloosbase.SoloOSEnv
		clientDriver sdfssdk.ClientDriver
		sfuseServer  sfuse.Server
	)

	err = soloOSEnv.InitWithSNet(options.SNetDriverServeAddr)
	util.AssertErrIsNil(err)

	if options.PProfListenAddr != "" {
		go util.PProfServe(options.PProfListenAddr)
	}

	{
		var nameNodeSRPCPeerID snettypes.PeerID
		nameNodeSRPCPeerID.SetStr(options.NameNodeSRPCPeerID)
		err = clientDriver.Init(&soloOSEnv,
			nameNodeSRPCPeerID,
			options.DBDriver, options.Dsn)
		util.AssertErrIsNil(err)
	}

	err = sfuseServer.Init(options.SFuseOptions,
		options.DefaultNetBlockCap,
		options.DefaultMemBlockCap,
		options.DefaultMemBlocksLimit,
		&clientDriver)
	util.AssertErrIsNil(err)

	err = sfuseServer.Serve()
	util.AssertErrIsNil(err)

	err = sfuseServer.Close()
	util.AssertErrIsNil(err)
}
