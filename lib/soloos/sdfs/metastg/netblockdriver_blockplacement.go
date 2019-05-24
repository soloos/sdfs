package metastg

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/snettypes"
)

func (p *NetBlockDriver) ChooseDataNodesForNewNetBlock(uNetINode sdfsapitypes.NetINodeUintptr) (snettypes.PeerGroup, error) {
	return p.helper.ChooseDataNodesForNewNetBlock(uNetINode)
}
