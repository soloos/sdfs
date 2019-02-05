package memstg

import (
	"soloos/sdfs/netstg"
	"soloos/sdfs/types"
	"soloos/snet"
	"soloos/util"
	"soloos/util/offheap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetINodeDriverNetINodeWrite(t *testing.T) {
	var (
		mockServer       netstg.MockServer
		mockNetINodePool types.MockNetINodePool
		snetDriver       snet.NetDriver
		netBlockDriver   netstg.NetBlockDriver
		memBlockDriver   MemBlockDriver
		netINodeDriver   NetINodeDriver
		maxBlocks        int32 = 16
		i                int32
		netBlockCap      int   = 4
		memBlockCap      int   = 4
		blockChunksLimit int32 = 2
		uNetINode        types.NetINodeUintptr
	)
	assert.NoError(t, mockNetINodePool.Init(&offheap.DefaultOffheapDriver))
	MakeDriversWithMockServerForTest("127.0.0.1:10023", &mockServer, &snetDriver,
		&netBlockDriver, &memBlockDriver, &netINodeDriver,
		memBlockCap, blockChunksLimit)
	uNetINode = mockNetINodePool.AllocNetINode(netBlockCap, memBlockCap)

	for i = 0; i <= maxBlocks; i++ {
		writeOffset := uint64(uint64(i) * uint64(memBlockCap))

		assert.NoError(t, netINodeDriver.PWriteWithMem(uNetINode, []byte{(byte)(i), (byte)(i * 2)}, writeOffset))

		memBlockIndex := int(writeOffset / uint64(uNetINode.Ptr().MemBlockCap))
		uMemBlock, _ := memBlockDriver.MustGetMemBlockWithReadAcquire(uNetINode, memBlockIndex)
		memBlockData := *uMemBlock.Ptr().BytesSlice()
		assert.Equal(t, memBlockData[0], (byte)(i))
		assert.Equal(t, memBlockData[1], (byte)(i*2))
		uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
	}

	for i = 0; i <= maxBlocks; i++ {
		writeOffset := uint64(uint64(i) * uint64(memBlockCap))
		memBlockIndex := int(writeOffset / uint64(uNetINode.Ptr().MemBlockCap))
		uMemBlock, _ := memBlockDriver.MustGetMemBlockWithReadAcquire(uNetINode, memBlockIndex)
		util.AssertErrIsNil(netINodeDriver.Flush(uNetINode))
		uMemBlock.Ptr().Chunk.Ptr().ReadRelease()
	}

	assert.NoError(t, mockServer.Close())
}
