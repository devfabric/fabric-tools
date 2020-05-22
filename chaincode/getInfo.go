package chaincode

import (
	"fmt"

	"code.uni-ledger.com/switch/fabric-tools/types"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/core/common/ccprovider"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"

	"github.com/pkg/errors"
)

const (
	lifecycleSCC = "lscc"

	getCCDataFunc     = "getccdata"
	getCollConfigFunc = "getcollectionsconfig"
)

// GetChainCodeInfo 查询链码信息
func (M *ChainCodeManager) GetChainCodeInfo(user, cc string) (*types.PbChainCodeInfo, error) {
	channelClient, err := M.ChannelClient(user, M.OrgID())
	if err != nil {
		return nil, errors.Errorf("error retrieving channel client: %v", err)
	}

	ccData, err := M.getCCData(channelClient, cc)
	if err != nil {
		return nil, errors.WithMessagef(err, "error querying for chaincode data")
	}

	// collConfig, err := M.getCollConfig(channelClient, cc)
	// if err != nil {
	// 	return errors.WithMessagef(err, "error querying for collection config")
	// }

	ccInfo := types.PbChainCodeInfo{
		Name:                ccData.Name,
		Version:             ccData.Version,
		Escc:                ccData.Escc,
		Vscc:                ccData.Vscc,
		ID:                  ccData.Id,
		Policy:              ccData.Policy,
		Data:                ccData.Data,
		InstantiationPolicy: ccData.InstantiationPolicy,
		// Coll                string //隐私数据策略//暂时不支持
	}

	return &ccInfo, nil
}

func (M *ChainCodeManager) getCCData(channelClient *channel.Client, cc string) (*ccprovider.ChaincodeData, error) {
	var args [][]byte
	args = append(args, []byte(M.Cfg.ChannelName()))
	args = append(args, []byte(cc))

	peer := M.Peers()[0]
	fmt.Printf("querying chaincode info for %s on peer: %s...\n", cc, peer.URL())

	response, err := channelClient.Query(
		channel.Request{ChaincodeID: lifecycleSCC, Fcn: getCCDataFunc, Args: args},
		channel.WithTargetEndpoints(peer.URL()))
	if err != nil {
		return nil, errors.Errorf("error querying for chaincode info: %v", err)
	}

	ccData := &ccprovider.ChaincodeData{}
	err = proto.Unmarshal(response.Payload, ccData)
	if err != nil {
		return nil, errors.Errorf("error unmarshalling chaincode data: %v", err)
	}
	return ccData, nil
}

func (M *ChainCodeManager) getCollConfig(channelClient *channel.Client, cc string) (*common.CollectionConfigPackage, error) {
	var args [][]byte
	args = append(args, []byte(cc))

	peer := M.Peers()[0]
	fmt.Printf("querying collections config for %s on peer: %s...\n", cc, peer.URL())

	response, err := channelClient.Query(
		channel.Request{ChaincodeID: lifecycleSCC, Fcn: getCollConfigFunc, Args: args},
		channel.WithTargetEndpoints(peer.URL()))
	if err != nil {
		return nil, errors.Errorf("error querying for collections config: %v", err)
	}

	collConfig := &common.CollectionConfigPackage{}
	err = proto.Unmarshal(response.Payload, collConfig)
	if err != nil {
		return nil, errors.Errorf("error unmarshalling collections config: %v", err)
	}
	return collConfig, nil
}
