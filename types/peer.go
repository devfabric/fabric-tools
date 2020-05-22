package types

// Peer peer信息
type Peer struct {
	URL         string `json:"url,omitempty"`
	MSP         string `json:"msp,omitempty"`
	BlockHeight uint64 `json:"block_height,omitempty"`
}
