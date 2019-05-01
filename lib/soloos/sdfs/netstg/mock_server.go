package netstg

import (
	"soloos/common/snet"
	"soloos/common/snet/srpc"
	snettypes "soloos/common/snet/types"
	"soloos/common/util"
	"soloos/sdfs/api"
	"soloos/sdfs/protocol"
	"soloos/sdfs/types"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
)

const (
	DefaultMockServerAddr = "127.0.0.1:10020"
)

type MockServer struct {
	snetDriver    *snet.NetDriver
	network       string
	addr          string
	srpcServer    srpc.Server
	dataNodePeers []snettypes.PeerUintptr
}

func (p *MockServer) SetDataNodePeers(dataNodePeers []snettypes.PeerUintptr) {
	p.dataNodePeers = dataNodePeers
}

func (p *MockServer) Init(snetDriver *snet.NetDriver, network string, addr string) error {
	var err error
	p.snetDriver = snetDriver
	p.network = network
	p.addr = addr
	err = p.srpcServer.Init(p.network, p.addr)
	if err != nil {
		return err
	}

	p.srpcServer.RegisterService("/DataNode/Register", p.DataNodeRegister)
	p.srpcServer.RegisterService("/NetINode/MustGet", p.NetINodeMustGet)
	p.srpcServer.RegisterService("/NetINode/PWrite", p.NetINodePWrite)
	p.srpcServer.RegisterService("/NetINode/PRead", p.NetINodePRead)
	p.srpcServer.RegisterService("/NetINode/CommitSizeInDB", p.NetINodeCommitSizeInDB)
	p.srpcServer.RegisterService("/NetBlock/PrepareMetaData", p.NetBlockPrepareMetaData)
	p.dataNodePeers = make([]snettypes.PeerUintptr, 3)
	for i := 0; i < len(p.dataNodePeers); i++ {
		p.dataNodePeers[i] = p.snetDriver.AllocPeer(p.addr, types.DefaultSDFSRPCProtocol)
	}

	return nil
}

func (p *MockServer) DataNodeRegister(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {

	var param = make([]byte, reqBodySize)
	util.AssertErrIsNil(conn.ReadAll(param))

	var protocolBuilder flatbuffers.Builder
	api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	util.AssertErrIsNil(conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():]))

	return nil
}

func (p *MockServer) NetINodeMustGet(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {

	var blockData = make([]byte, reqBodySize)
	util.AssertErrIsNil(conn.ReadAll(blockData))

	// request
	var req protocol.NetINodeInfoRequest
	req.Init(blockData[:reqParamSize], flatbuffers.GetUOffsetT(blockData[:reqParamSize]))

	// response
	var protocolBuilder flatbuffers.Builder
	api.SetNetINodeInfoResponse(&protocolBuilder, req.Size(), req.NetBlockCap(), req.MemBlockCap())
	util.AssertErrIsNil(conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():]))

	return nil
}

func (p *MockServer) NetINodePWrite(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {

	var reqBody = make([]byte, reqBodySize)
	util.AssertErrIsNil(conn.ReadAll(reqBody))

	var req protocol.NetINodePWriteRequest
	req.Init(reqBody[:reqParamSize], flatbuffers.GetUOffsetT(reqBody[:reqParamSize]))
	var backends = make([]protocol.SNetPeer, req.TransferBackendsLength())
	for i := 0; i < len(backends); i++ {
		req.TransferBackends(&backends[i], i)
	}

	var protocolBuilder flatbuffers.Builder
	api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	util.AssertErrIsNil(conn.SimpleResponse(reqID, respBody))

	return nil
}

func (p *MockServer) NetINodePRead(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {

	var reqData = make([]byte, reqBodySize)
	util.AssertErrIsNil(conn.ReadAll(reqData))

	var req protocol.NetINodePReadRequest
	req.Init(reqData[:reqParamSize], flatbuffers.GetUOffsetT(reqData[:reqParamSize]))

	var protocolBuilder flatbuffers.Builder
	api.SetNetINodePReadResponse(&protocolBuilder, req.Length())

	respBody := protocolBuilder.Bytes[protocolBuilder.Head():]
	util.AssertErrIsNil(conn.Response(reqID, respBody, make([]byte, req.Length())))
	return nil
}

func (p *MockServer) NetINodeCommitSizeInDB(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {

	var reqData = make([]byte, reqBodySize)
	util.AssertErrIsNil(conn.ReadAll(reqData))

	// response
	var protocolBuilder flatbuffers.Builder
	api.SetCommonResponseCode(&protocolBuilder, snettypes.CODE_OK)
	conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():])

	return nil
}

func (p *MockServer) NetBlockPrepareMetaData(reqID uint64,
	reqBodySize, reqParamSize uint32,
	conn *snettypes.Connection) error {

	var blockData = make([]byte, reqBodySize)
	util.AssertErrIsNil(conn.ReadAll(blockData))

	// request
	var req protocol.NetINodeNetBlockInfoRequest
	req.Init(blockData[:reqParamSize], flatbuffers.GetUOffsetT(blockData[:reqParamSize]))

	// response
	var protocolBuilder flatbuffers.Builder
	api.SetNetINodeNetBlockInfoResponse(&protocolBuilder, p.dataNodePeers[:], req.Cap(), req.Cap())
	util.AssertErrIsNil(conn.SimpleResponse(reqID, protocolBuilder.Bytes[protocolBuilder.Head():]))

	return nil
}

func (p *MockServer) Serve() error {
	return p.srpcServer.Serve()
}

func (p *MockServer) Close() error {
	var err error
	err = p.srpcServer.Close()
	time.Sleep(time.Second)
	return err
}
