package solodn

import (
	"fmt"
	"soloos/common/iron"
	"soloos/common/snet"
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"soloos/solofs/memstg"
	"soloos/solofs/metastg"
	"soloos/solofs/solonn"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	go util.PProfServe("192.168.56.100:17221")
	var (
		solonnIns                 solonn.Solonn
		solonnSrpcPeerID          = snet.MakeSysPeerID("SolonnSrpcForTest")
		solonnWebPeerID           = snet.MakeSysPeerID("SolonnWebForTest")
		solonnSrpcListenAddr      = "127.0.0.1:10401"
		solonnWebListenAddr       = "127.0.0.1:10402"
		netDriverWebServer        iron.Server
		netDriverServerListenAddr = "127.0.0.1:10403"
		netDriverServerServeAddr  = "http://127.0.0.1:10403"
		metaStgForSolonn          metastg.MetaStg

		solodns               [6]Solodn
		solodnSrpcPeerIDs     [6]snet.PeerID
		solodnSrpcListenAddrs = []string{
			"127.0.0.1:10410",
			"127.0.0.1:10411",
			"127.0.0.1:10412",
			"127.0.0.1:10413",
			"127.0.0.1:10414",
			"127.0.0.1:10415",
		}
	)

	var (
		soloosEnvForClient      soloosbase.SoloosEnv
		memBlockDriverForClient memstg.MemBlockDriver
		netBlockDriverForClient memstg.NetBlockDriver
		netINodeDriverForClient memstg.NetINodeDriver

		soloosEnvForSolonn      soloosbase.SoloosEnv
		memBlockDriverForSolonn memstg.MemBlockDriver
		netBlockDriverForSolonn memstg.NetBlockDriver
		netINodeDriverForSolonn memstg.NetINodeDriver

		soloosEnvForSolodns      [6]soloosbase.SoloosEnv
		memBlockDriverForSolodns [6]memstg.MemBlockDriver
		netBlockDriverForSolodns [6]memstg.NetBlockDriver
		netINodeDriverForSolodns [6]memstg.NetINodeDriver

		netBlockCap int   = 32
		memBlockCap int   = 16
		blocksLimit int32 = 4
		uNetINode   solofstypes.NetINodeUintptr
		i           int
		err         error
	)

	assert.NoError(t, soloosEnvForSolonn.InitWithSNet(""))
	{
		var webServerOptions = iron.Options{
			ListenStr: netDriverServerListenAddr,
			ServeStr:  netDriverServerServeAddr,
		}
		util.AssertErrIsNil(netDriverWebServer.Init(webServerOptions))
	}
	go func() {
		assert.NoError(t, soloosEnvForSolonn.SNetDriver.PrepareServer("",
			&netDriverWebServer,
			nil, nil))
		assert.NoError(t, soloosEnvForSolonn.SNetDriver.ServerServe())
	}()
	time.Sleep(100 * time.Millisecond)
	metastg.MakeMetaStgForTest(&soloosEnvForSolonn, &metaStgForSolonn)

	assert.NoError(t, soloosEnvForClient.InitWithSNet(netDriverServerServeAddr))

	memstg.MemStgMakeDriversForTest(&soloosEnvForClient,
		solonnSrpcListenAddr,
		&memBlockDriverForClient, &netBlockDriverForClient, &netINodeDriverForClient, memBlockCap, blocksLimit)

	memstg.MemStgMakeDriversForTest(&soloosEnvForSolonn,
		solonnSrpcListenAddr,
		&memBlockDriverForSolonn, &netBlockDriverForSolonn, &netINodeDriverForSolonn, memBlockCap, blocksLimit)
	solonn.MakeSolonnForTest(&soloosEnvForSolonn, &solonnIns, &metaStgForSolonn,
		solonnSrpcPeerID, solonnSrpcListenAddr,
		solonnWebPeerID, solonnWebListenAddr,
		&memBlockDriverForSolonn, &netBlockDriverForSolonn, &netINodeDriverForSolonn)

	go func() {
		assert.NoError(t, solonnIns.Serve())
	}()
	time.Sleep(time.Millisecond * 300)

	for i = 0; i < len(solodnSrpcListenAddrs); i++ {
		assert.NoError(t, soloosEnvForSolodns[i].InitWithSNet(netDriverServerServeAddr))
		solodnSrpcPeerIDs[i] = snet.MakeSysPeerID(fmt.Sprintf("SolodnForTest_%v", i))

		memstg.MemStgMakeDriversForTest(&soloosEnvForSolodns[i],
			solonnSrpcListenAddr,
			&memBlockDriverForSolodns[i],
			&netBlockDriverForSolodns[i],
			&netINodeDriverForSolodns[i],
			memBlockCap, blocksLimit)

		MakeSolodnForTest(&soloosEnvForSolodns[i],
			&solodns[i],
			solodnSrpcPeerIDs[i], solodnSrpcListenAddrs[i],
			solonnSrpcPeerID, solonnSrpcListenAddr,
			&memBlockDriverForSolodns[i],
			&netBlockDriverForSolodns[i],
			&netINodeDriverForSolodns[i])
		go func(localI int) {
			assert.NoError(t, solodns[localI].Serve())
		}(i)
	}
	time.Sleep(time.Millisecond * 300)

	var (
		netINodeID solofstypes.NetINodeID
	)
	solofstypes.InitTmpNetINodeID(&netINodeID)
	uNetINode, err = netINodeDriverForClient.MustGetNetINode(netINodeID, 0, netBlockCap, memBlockCap)
	defer netINodeDriverForClient.ReleaseNetINode(uNetINode)
	assert.NoError(t, err)

	writeData := make([]byte, 73)
	writeData[3] = 1
	writeData[7] = 2
	writeData[8] = 3
	writeData[33] = 4
	writeData[60] = 5
	assert.NoError(t, netINodeDriverForClient.PWriteWithMem(uNetINode, writeData, 612))
	assert.NoError(t, netINodeDriverForClient.Sync(uNetINode))

	var readData []byte
	readData = make([]byte, 73)
	_, err = netINodeDriverForClient.PReadWithMem(uNetINode, readData, 612)
	util.AssertErrIsNil(err)
	assert.Equal(t, writeData, readData)

	var maxWriteTimes int = 1
	for i = 0; i < maxWriteTimes; i++ {
		assert.NoError(t, netINodeDriverForClient.PWriteWithMem(uNetINode, writeData, uint64(netBlockCap*600+8*i)))
	}

	readData = make([]byte, 73)
	_, err = netINodeDriverForClient.PReadWithMem(uNetINode, readData, 612)
	assert.NoError(t, err)
	assert.Equal(t, writeData, readData)

	time.Sleep(time.Microsecond * 600)
	for i = 0; i < len(solodnSrpcListenAddrs); i++ {
		assert.NoError(t, solodns[i].Close())
	}
	assert.NoError(t, solonnIns.Close())
}
