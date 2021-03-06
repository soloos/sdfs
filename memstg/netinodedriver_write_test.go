package memstg

import (
	"soloos/common/solofstypes"
	"soloos/common/soloosbase"
	"soloos/common/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetINodeDriverNetINodeWrite(t *testing.T) {
	var (
		soloosEnv         soloosbase.SoloosEnv
		mockServer        MockServer
		mockNetINodeTable MockNetINodeTable
		netBlockDriver    NetBlockDriver
		memBlockDriver    MemBlockDriver
		netINodeDriver    NetINodeDriver
		maxBlocks         int32 = 16
		i                 int32
		netBlockCap       int   = 4
		memBlockCap       int   = 4
		blocksLimit       int32 = 2
		uNetINode         solofstypes.NetINodeUintptr
	)
	util.AssertErrIsNil(soloosEnv.InitWithSNet(""))

	assert.NoError(t, mockNetINodeTable.Init(&soloosEnv))
	MemStgMakeDriversWithMockServerForTest(&soloosEnv, "127.0.0.1:10023", &mockServer,
		&netBlockDriver, &memBlockDriver, &netINodeDriver,
		memBlockCap, blocksLimit)
	uNetINode = mockNetINodeTable.AllocNetINode(netBlockCap, memBlockCap)
	defer mockNetINodeTable.ReleaseNetINode(uNetINode)

	for i = 0; i <= maxBlocks; i++ {
		writeOffset := uint64(uint64(i) * uint64(memBlockCap))

		assert.NoError(t, netINodeDriver.PWriteWithMem(uNetINode, []byte{(byte)(i), (byte)(i * 2)}, writeOffset))

		memBlockIndex := int32(writeOffset / uint64(uNetINode.Ptr().MemBlockCap))
		uMemBlock, _ := memBlockDriver.MustGetMemBlockWithReadAcquire(uNetINode, memBlockIndex)
		memBlockData := *uMemBlock.Ptr().BytesSlice()
		assert.Equal(t, memBlockData[0], (byte)(i))
		assert.Equal(t, memBlockData[1], (byte)(i*2))
		uMemBlock.Ptr().ReadRelease()
	}

	for i = 0; i <= maxBlocks; i++ {
		writeOffset := uint64(uint64(i) * uint64(memBlockCap))
		memBlockIndex := int32(writeOffset / uint64(uNetINode.Ptr().MemBlockCap))
		uMemBlock, _ := memBlockDriver.MustGetMemBlockWithReadAcquire(uNetINode, memBlockIndex)
		util.AssertErrIsNil(netINodeDriver.Sync(uNetINode))
		uMemBlock.Ptr().ReadRelease()
	}

	assert.NoError(t, mockServer.Close())
}
