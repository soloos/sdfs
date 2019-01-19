package metastg

func (p *DirTreeDriver) StatLimits() (uint64, uint64) {
	var (
		capacity uint64 = 1024 * 1024 * 1024 * 1024 * 1024 * 100
		files    uint64 = 1024 * 1024 * 1024 * 100
	)
	return capacity, files
}

func (p *DirTreeDriver) StatFs() (uint64, uint64) {
	// TODO return real result
	var (
		usedSize   uint64 = 1024 * 1024 * 100
		filesCount uint64 = 1000
	)
	return usedSize, filesCount
}

func (p *DirTreeDriver) BlkSize() uint32 {
	// TODO return real result
	var (
		blksize uint32 = 1024 * 4
	)
	return blksize
}