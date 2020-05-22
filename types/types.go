package types

// ChainCodes 链码数组
type ChainCodes []ChainCode

// PbChainCodeInfo 查询返回的链码信息
type PbChainCodeInfo struct {
	Name                string
	Version             string
	Path                string
	Input               string
	Escc                string
	Vscc                string
	ID                  []byte //不用
	Policy              []byte //背书策略
	Data                []byte //打包数据
	InstantiationPolicy []byte //实例化策略
	Coll                string //隐私数据策略//暂时不支持
}
