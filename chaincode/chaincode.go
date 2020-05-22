package chaincode

import "fabric-tools/manager"

// ChainCodeManager 用于通道操作
type ChainCodeManager struct {
	manager.Manager
}

// NewChainCodeManager  返回NewChannelManager
func NewChainCodeManager(m *manager.Manager) (*ChainCodeManager, error) {
	return &ChainCodeManager{*m}, nil
}
