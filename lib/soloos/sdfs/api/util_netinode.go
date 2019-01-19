package api

import (
	"soloos/sdfs/protocol"
	snettypes "soloos/snet/types"

	flatbuffers "github.com/google/flatbuffers/go"
)

func SetNetINodeInfoResponseError(protocolBuilder *flatbuffers.Builder, code int, err string) {
	protocolBuilder.Reset()
	var (
		errOff            flatbuffers.UOffsetT
		commonResponseOff flatbuffers.UOffsetT
	)
	errOff = protocolBuilder.CreateString(err)
	protocol.CommonResponseStart(protocolBuilder)
	protocol.CommonResponseAddCode(protocolBuilder, int32(code))
	protocol.CommonResponseAddError(protocolBuilder, errOff)
	commonResponseOff = protocol.CommonResponseEnd(protocolBuilder)

	protocol.NetINodeInfoResponseStart(protocolBuilder)
	protocol.NetINodeInfoResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	protocolBuilder.Finish(protocol.NetINodeInfoResponseEnd(protocolBuilder))
}

func SetNetINodeInfoResponse(protocolBuilder *flatbuffers.Builder,
	size uint64, netBlockCap int32, memBlockCap int32) {
	protocolBuilder.Reset()
	var (
		commonResponseOff flatbuffers.UOffsetT
	)
	protocol.CommonResponseStart(protocolBuilder)
	protocol.CommonResponseAddCode(protocolBuilder, snettypes.CODE_OK)
	commonResponseOff = protocol.CommonResponseEnd(protocolBuilder)

	protocol.NetINodeInfoResponseStart(protocolBuilder)
	protocol.NetINodeInfoResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	protocol.NetINodeInfoResponseAddSize(protocolBuilder, size)
	protocol.NetINodeInfoResponseAddNetBlockCap(protocolBuilder, int32(netBlockCap))
	protocol.NetINodeInfoResponseAddMemBlockCap(protocolBuilder, int32(memBlockCap))
	protocolBuilder.Finish(protocol.NetINodeInfoResponseEnd(protocolBuilder))
}

func SetNetINodePReadResponseError(protocolBuilder *flatbuffers.Builder, code int, err string) {
	protocolBuilder.Reset()
	var (
		errOff            flatbuffers.UOffsetT
		commonResponseOff flatbuffers.UOffsetT
	)
	errOff = protocolBuilder.CreateString(err)
	protocol.CommonResponseStart(protocolBuilder)
	protocol.CommonResponseAddCode(protocolBuilder, int32(code))
	protocol.CommonResponseAddError(protocolBuilder, errOff)
	commonResponseOff = protocol.CommonResponseEnd(protocolBuilder)

	protocol.NetINodePReadResponseStart(protocolBuilder)
	protocol.NetINodePReadResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	protocolBuilder.Finish(protocol.NetINodePReadResponseEnd(protocolBuilder))
}

func SetNetINodePReadResponse(protocolBuilder *flatbuffers.Builder, length int32) {
	protocolBuilder.Reset()
	var (
		commonResponseOff flatbuffers.UOffsetT
	)
	protocol.CommonResponseStart(protocolBuilder)
	protocol.CommonResponseAddCode(protocolBuilder, snettypes.CODE_OK)
	commonResponseOff = protocol.CommonResponseEnd(protocolBuilder)

	protocol.NetINodePReadResponseStart(protocolBuilder)
	protocol.NetINodePReadResponseAddCommonResponse(protocolBuilder, commonResponseOff)
	protocol.NetINodePReadResponseAddLength(protocolBuilder, length)
	protocolBuilder.Finish(protocol.NetINodePReadResponseEnd(protocolBuilder))
}
