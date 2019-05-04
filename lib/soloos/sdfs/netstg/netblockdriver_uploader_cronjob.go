package netstg

import (
	"soloos/common/util"
	"soloos/sdfs/types"
)

func (p *netBlockDriverUploader) cronUpload() error {
	var (
		uJob      types.UploadMemBlockJobUintptr
		pJob      *types.UploadMemBlockJob
		pNetINode *types.NetINode
		pNetBlock *types.NetBlock
		i         int
		ok        bool
		err       error
	)

	for {
		uJob, ok = <-p.uploadMemBlockJobChan
		if ok == false {
			panic("uploadMemBlockJobChan closed")
		}

		pJob = uJob.Ptr()
		pNetINode = pJob.UNetINode.Ptr()
		pNetBlock = pJob.UNetBlock.Ptr()

		// prepare upload job
		pJob.MetaDataStateMutex.Lock()
		if pJob.UploadMaskWaiting.Ptr().MaskArrayLen == 0 {
			// upload done and continue
			pJob.MetaDataStateMutex.Unlock()
			goto ONE_RUN_DONE
		} else {
			// upload job exists
			pJob.UploadMaskSwap()
			pJob.MetaDataStateMutex.Unlock()
		}

		util.AssertTrue(pNetBlock.SyncDataBackends.Len > 0)

		// start upload
		// upload primary backend
		err = p.driver.dataNodeClient.UploadMemBlock(uJob, 0, pNetBlock.SyncDataPrimaryBackendTransferCount)

		// upload other backends
		for i = pNetBlock.SyncDataPrimaryBackendTransferCount + 1; i < pNetBlock.SyncDataBackends.Len; i++ {
			err = p.driver.dataNodeClient.UploadMemBlock(uJob, i, 0)
		}

	ONE_RUN_DONE:
		pJob.SyncDataSig.Done()
		pNetINode.SyncDataSig.Done()

		if err != nil {
			// TODO catch error
			pNetINode.LastSyncDataError = err
		} else {
			pJob.UploadMaskProcessing.Ptr().Reset()
		}
	}

	return nil
}
