package netstg

import (
	"soloos/common/snet"
	"soloos/sdbone/offheap"
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetBlockDriver(t *testing.T) {
	var (
		offheapDriver     = &offheap.DefaultOffheapDriver
		mockNetINodeTable types.MockNetINodeTable
		mockMemBlockTable types.MockMemBlockTable
		snetDriver        snet.NetDriver
		snetClientDriver  snet.ClientDriver
		mockServer        MockServer
		nameNodeClient    api.NameNodeClient
		dataNodeClient    api.DataNodeClient
		netBlockDriver    NetBlockDriver
	)
	mockServerAddr := "127.0.0.1:10021"
	assert.NoError(t, mockNetINodeTable.Init(&offheap.DefaultOffheapDriver))
	assert.NoError(t, mockMemBlockTable.Init(offheapDriver, 1024))
	MakeDriversWithMockServerForTest(&snetDriver, &snetClientDriver,
		mockServerAddr, &mockServer,
		&nameNodeClient, &dataNodeClient,
		&netBlockDriver)

	var uPeer0 = snetDriver.AllocPeer(mockServerAddr, types.DefaultSDFSRPCProtocol)
	var uPeer1 = snetDriver.AllocPeer(mockServerAddr, types.DefaultSDFSRPCProtocol)

	data := make([]byte, 8)
	for i := 0; i < len(data); i++ {
		data[i] = 1
	}

	uNetINode := mockNetINodeTable.AllocNetINode(1024, 128)

	netBlockIndex := int32(10)
	uNetBlock, err := netBlockDriver.MustGetNetBlock(uNetINode, netBlockIndex)
	assert.NoError(t, err)
	uNetBlock.Ptr().StorDataBackends.Append(uPeer0)
	uNetBlock.Ptr().StorDataBackends.Append(uPeer1)
	uMemBlock := mockMemBlockTable.AllocMemBlock()
	memBlockIndex := int32(0)
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, 0, 12))
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, 11, 24))
	assert.NoError(t, netBlockDriver.PWrite(uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, 30, 64))
	assert.NoError(t, netBlockDriver.FlushMemBlock(uNetINode, uNetBlock, uMemBlock))

	assert.NoError(t, mockServer.Close())
}
