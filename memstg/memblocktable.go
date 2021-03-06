package memstg

import (
	"math"
	"soloos/common/solofstypes"
	"soloos/solodb/offheap"
)

type MemBlockTable struct {
	options MemBlockTableOptions
	driver  *MemBlockDriver

	objectSize       int
	tmpMemBlockTable offheap.HKVTableWithBytes12
	memBlockTable    offheap.HKVTableWithBytes12
}

func (p *MemBlockTable) Init(
	options MemBlockTableOptions,
	driver *MemBlockDriver,
) error {
	var err error

	p.options = options
	p.driver = driver

	objectSize := p.options.ObjectSize
	objectsLimit := p.options.ObjectsLimit

	memBlockTableObjectsLimit := int32(math.Ceil(float64(objectsLimit) * 0.9))
	p.objectSize = objectSize

	err = p.driver.OffheapDriver.InitHKVTableWithBytes12(&p.memBlockTable, "MemBlock",
		int(solofstypes.MemBlockStructSize+uintptr(p.objectSize)),
		memBlockTableObjectsLimit,
		offheap.DefaultKVTableSharedCount,
		p.hkvTableInvokePrepareNewBlock,
		p.hkvTableInvokeBeforeReleaseBlock,
	)
	if err != nil {
		return err
	}

	tmpMemBlockTableObjectsLimit := objectsLimit - memBlockTableObjectsLimit
	if tmpMemBlockTableObjectsLimit == 0 {
		tmpMemBlockTableObjectsLimit = 1
	}
	err = p.driver.OffheapDriver.InitHKVTableWithBytes12(&p.tmpMemBlockTable, "TmpMemBlock",
		int(solofstypes.MemBlockStructSize+uintptr(p.objectSize)),
		tmpMemBlockTableObjectsLimit,
		offheap.DefaultKVTableSharedCount,
		p.hkvTableInvokePrepareNewBlock,
		p.hkvTableInvokeBeforeReleaseTmpBlock,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *MemBlockTable) hkvTableInvokePrepareNewBlock(uMemBlock uintptr) {
	pMemBlock := solofstypes.MemBlockUintptr(uMemBlock).Ptr()
	pMemBlock.Reset()
	pMemBlock.Bytes.Data = uMemBlock + solofstypes.MemBlockStructSize
	pMemBlock.Bytes.Len = p.objectSize
	pMemBlock.Bytes.Cap = pMemBlock.Bytes.Len
	pMemBlock.CompleteInit()
}
