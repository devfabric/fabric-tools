package types

// ChainCodeInfo 链码信息
type ChainCodeInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Path    string `json:"path,omitempty"`
	Input   string `json:"input,omitempty"`
	Escc    string `json:"escc,omitempty"`
	Vscc    string `json:"vscc,omitempty"`
	ID      []byte `json:"id,omitempty"`
}

// ChainCode 链码信息
type ChainCode struct {
	Name    string
	Version string
	Path    string
	Gopath  string

	ChaincodePolicy string
}
