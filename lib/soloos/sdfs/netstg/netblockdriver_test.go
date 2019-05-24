package netstg

import (
	"soloos/common/sdfsapi"
	"soloos/common/sdfsapitypes"
	"soloos/common/soloosbase"
	"soloos/sdfs/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetBlockDriver(t *testing.T) {
	var (
		soloOSEnv         soloosbase.SoloOSEnv
		mockNetINodeTable types.MockNetINodeTable
		mockMemBlockTable types.MockMemBlockTable
		mockServer        MockServer
		nameNodeClient    sdfsapi.NameNodeClient
		dataNodeClient    sdfsapi.DataNodeClient
		netBlockDriver    NetBlockDriver
	)
	assert.NoError(t, soloOSEnv.Init())
	mockServerAddr := "127.0.0.1:10021"
	assert.NoError(t, mockNetINodeTable.Init(&soloOSEnv))
	assert.NoError(t, mockMemBlockTable.Init(&soloOSEnv, 1024))
	MakeDriversWithMockServerForTest(&soloOSEnv,
		mockServerAddr, &mockServer,
		&nameNodeClient, &dataNodeClient,
		&netBlockDriver)

	var uPeer0, _ = soloOSEnv.SNetDriver.MustGetPeer(nil, mockServerAddr, sdfsapitypes.DefaultSDFSRPCProtocol)
	var uPeer1, _ = soloOSEnv.SNetDriver.MustGetPeer(nil, mockServerAddr, sdfsapitypes.DefaultSDFSRPCProtocol)

	data := make([]byte, 8)
	for i := 0; i < len(data); i++ {
		data[i] = 1
	}

	uNetINode := mockNetINodeTable.AllocNetINode(1024, 128)
	defer mockNetINodeTable.ReleaseNetINode(uNetINode)

	netBlockIndex := int32(10)
	uNetBlock, err := netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
	defer netBlockDriver.ReleaseNetBlock(uNetBlock)
	assert.NoError(t, err)
	uNetBlock.Ptr().StorDataBackends.Append(uPeer0)
	uNetBlock.Ptr().StorDataBackends.Append(uPeer1)
	uMemBlock := mockMemBlockTable.AllocMemBlock()
	memBlockIndex := int32(0)
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, 0, 12))
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, 11, 24))
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, 30, 64))
	uMemBlock.Ptr().UploadJob.SyncDataSig.Wait()

	assert.NoError(t, mockServer.Close())
}
