package chaincode

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/pkg/errors"
)

// Invoke 执行链码调用
func (M *ChainCodeManager) Invoke(cc, user, invokeFunc string, args [][]byte, orgs ...string) (*channel.Response, error) {
	orgID := M.OrgID()
	if len(orgs) > 0 {
		orgID = orgs[0]
	}

	channelClient, err := M.ChannelClient(user, orgID)

	var opts []channel.RequestOption
	opts = append(opts, channel.WithRetry(retry.DefaultResMgmtOpts))
	// opts = append(opts, channel.WithBeforeRetry(func(err error) {
	// 	t.attempt++
	// }))
	if len(M.Peers()) > 0 {
		opts = append(opts, channel.WithTargets(M.Peers()...))
	}
	response, err := channelClient.Execute(
		channel.Request{
			ChaincodeID: cc,
			Fcn:         invokeFunc,
			Args:        args,
		},
		opts...,
	)
	if err != nil {
		return nil, errors.Errorf("SendTransactionProposal return error:%v", err)
	}

	// txID := string(response.TransactionID)

	return &response, nil
}
