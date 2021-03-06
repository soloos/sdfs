package localfs

import (
	"soloos/common/snet"
	"soloos/common/solofstypes"
	"soloos/solodb/offheap"
)

func (p *Fd) Upload(uJob solofstypes.UploadMemBlockJobUintptr) error {
	var (
		req                 snet.SNetReq
		netINodeWriteOffset int64
		memBlockCap         int
		uploadChunkMask     offheap.ChunkMask
		writeData           []byte
		err                 error
	)

	uploadChunkMask = uJob.Ptr().GetProcessingChunkMask()

	req.OffheapBody.OffheapBytes = uJob.Ptr().UMemBlock.Ptr().Bytes.Data
	memBlockCap = uJob.Ptr().UMemBlock.Ptr().Bytes.Len
	for chunkMaskIndex := 0; chunkMaskIndex < uploadChunkMask.MaskArrayLen; chunkMaskIndex++ {
		netINodeWriteOffset = int64(memBlockCap)*int64(uJob.Ptr().MemBlockIndex) +
			int64(uploadChunkMask.MaskArray[chunkMaskIndex].Offset)

		writeData = (*uJob.Ptr().UMemBlock.Ptr().BytesSlice())[uploadChunkMask.MaskArray[chunkMaskIndex].Offset:uploadChunkMask.MaskArray[chunkMaskIndex].End]
		err = p.WriteAt(writeData, netINodeWriteOffset)
		if err != nil {
			goto PWRITE_DONE
		}
	}

PWRITE_DONE:
	return err
}

func (p *Fd) WriteAt(data []byte, netINodeOffset int64) error {
	var (
		off int
		n   int
		err error
	)
	for off = 0; off < len(data); off += n {
		n, err = p.file.WriteAt(data, netINodeOffset+int64(off))
		if err != nil {
			return err
		}
	}
	return nil
}
