package chaincode

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

// Query 执行链码查询
func (M *ChainCodeManager) Query(cc, user, invokeFunc string, args [][]byte, orgs ...string) (channel.Response, error) {
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

	var additionalHandlers []invoke.Handler
	// if t.validate {
	// 	// Add the validation handlers
	// 	additionalHandlers = append(additionalHandlers,
	// 		invoke.NewEndorsementValidationHandler(
	// 			invoke.NewSignatureValidationHandler(),
	// 		),
	// 	)
	// }

	request := channel.Request{
		ChaincodeID: cc,
		Fcn:         invokeFunc,
		Args:        args,
	}
	response, err := channelClient.InvokeHandler(
		invoke.NewProposalProcessorHandler(
			invoke.NewEndorsementHandler(additionalHandlers...),
		),
		request, opts...)

	return response, err
}
