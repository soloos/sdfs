package netstg

import (
	sdbapitypes "soloos/common/sdbapi/types"
	"soloos/sdfs/types"
)

func (p *netBlockDriverUploader) PWrite(uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr, netBlockIndex int32,
	uMemBlock types.MemBlockUintptr, memBlockIndex int32,
	offset, end int) error {

	var (
		isMergeEventHappened    bool
		isMergeWriteMaskSuccess bool = false
		pMemBlock                    = uMemBlock.Ptr()
	)

	if pMemBlock.UploadJob.MetaDataState.Load() == sdbapitypes.MetaDataStateUninited {
		// TODO: refine me
		p.PrepareUploadMemBlockJob(&pMemBlock.UploadJob,
			uNetINode, uNetBlock, netBlockIndex, uMemBlock, memBlockIndex, uNetBlock.Ptr().StorDataBackends)
	}

	for isMergeWriteMaskSuccess == false {
		pMemBlock.UploadJob.UploadMaskMutex.Lock()
		isMergeEventHappened, isMergeWriteMaskSuccess =
			pMemBlock.UploadJob.UploadMaskWaiting.Ptr().MergeIncludeNeighbour(offset, end)
		pMemBlock.UploadJob.UploadMaskMutex.Unlock()

		if isMergeWriteMaskSuccess {
			if isMergeEventHappened == false {
				pMemBlock.UploadJob.UNetINode.Ptr().SyncDataSig.Add(1)
				pMemBlock.UploadJob.SyncDataSig.Add(1)
				p.uploadMemBlockJobChan <- pMemBlock.GetUploadMemBlockJobUintptr()
			}
		}

		if isMergeWriteMaskSuccess == false {
			pMemBlock.UploadJob.SyncDataSig.Wait()
		}
	}

	return nil
}
