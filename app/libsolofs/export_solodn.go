package main

import "C"
import (
	"io"
	"reflect"
	"soloos/common/log"
	"soloos/common/solofsapi"
	"unsafe"
)

//export GoSolofsPappend
func GoSolofsPappend(fdID uint64, buffer unsafe.Pointer, bufferLen, offset int32) (int32, C.int) {
	var (
		fd  = env.solofsClient.FdTableGetFd(fdID)
		err error
	)

	var data = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufferLen),
		Cap:  int(bufferLen),
	}))

	err = env.solofsClient.SimpleWriteWithMem(fd.FsINodeIno, data, uint64(offset))
	if err != nil {
		return 0, solofsapi.CODE_ERR
	}

	return bufferLen, 0
}

//export GoSolofsAppend
func GoSolofsAppend(fdID uint64, buffer unsafe.Pointer, bufferLen int32) (int32, C.int) {
	var (
		fd  = env.solofsClient.FdTableGetFd(fdID)
		err error
	)

	var data = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufferLen),
		Cap:  int(bufferLen),
	}))
	err = env.solofsClient.SimpleWriteWithMem(fd.FsINodeIno, data, fd.AppendPosition)
	if err != nil {
		log.Warn(err)
		return 0, solofsapi.CODE_ERR
	}

	env.solofsClient.FdTableFdAddAppendPosition(fdID, uint64(bufferLen))

	return bufferLen, solofsapi.CODE_OK
}

//export GoSolofsRead
func GoSolofsRead(fdID uint64, buffer unsafe.Pointer, bufferLen int32) (int32, C.int) {
	var (
		fd             = env.solofsClient.FdTableGetFd(fdID)
		readDataLength int
		err            error
	)

	var data = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufferLen),
		Cap:  int(bufferLen),
	}))
	readDataLength, err = env.solofsClient.SimpleReadWithMem(fd.FsINodeIno, data, fd.ReadPosition)
	if err != nil && err != io.EOF {
		log.Warn(err, readDataLength)
		return int32(readDataLength), solofsapi.CODE_ERR
	}

	env.solofsClient.FdTableFdAddReadPosition(fdID, uint64(bufferLen))

	return int32(readDataLength), solofsapi.CODE_OK
}

//export GoSolofsPread
func GoSolofsPread(fdID uint64, buffer unsafe.Pointer, bufferLen int32, position uint64) (int32, C.int) {
	var (
		fd             = env.solofsClient.FdTableGetFd(fdID)
		readDataLength int
		err            error
	)

	var data = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufferLen),
		Cap:  int(bufferLen),
	}))

	readDataLength, err = env.solofsClient.SimpleReadWithMem(fd.FsINodeIno, data, position)
	if err != nil {
		return int32(readDataLength), solofsapi.CODE_ERR
	}

	return int32(readDataLength), solofsapi.CODE_OK
}

//export GoSolofsCloseFile
func GoSolofsCloseFile(fdID uint64) C.int {
	ret := doFlushINode(fdID)
	// env.solofsClient.FdTableReleaseFd(fdID)
	return ret
}

//export GoSolofsFlushFile
func GoSolofsFlushFile(fdID uint64) C.int {
	return doFlushINode(fdID)
}

//export GoSolofsHFlushINode
func GoSolofsHFlushINode(fdID uint64) C.int {
	return doFlushINode(fdID)
}

//export GoSolofsHSyncINode
func GoSolofsHSyncINode(fdID uint64) C.int {
	return doFlushINode(fdID)
}

func doFlushINode(fdID uint64) C.int {
	var (
		fd = env.solofsClient.FdTableGetFd(fdID)
	)

	env.solofsClient.SimpleFlush(fd.FsINodeIno)

	return solofsapi.CODE_OK
}
